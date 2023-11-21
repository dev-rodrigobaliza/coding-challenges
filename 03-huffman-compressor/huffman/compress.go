package huffman

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"path"
)

const (
	CompressedFileStartFlag uint16 = 0x5259
	CompressedFileEndFlag   uint16 = 0x414E
)

var (
	ErrCanNotParseFileHeader = fmt.Errorf("can not parse file header")
)

// compressBytesWith Byte slices are compressed using a given Huffman encoding table
// Returns the compressed byte slice and the number of valid bits after compression
func compressBytesWith(data []byte, table HuffmanEncTable) ([]byte, uint64, error) {
	w := NewBitsWriter()
	var totalBits uint64 = 0

	// Traverse the data, encoding each byte
	for i, b := range data {
		code := table.Get(b)
		if code == nil {
			return nil, 0, fmt.Errorf("code for %b(%c at %d) not found", b, b, i)
		}
		bitlen := code.BitLen()
		totalBits += uint64(bitlen)
		err := w.WriteUint32(code.Bits(), uint8(bitlen))
		if err != nil {
			return nil, 0, fmt.Errorf(err.Error())
		}
	}
	return w.Buf(), totalBits, nil
}

// CompressBytes Compress a slice of bytes
// Returns the compressed byte slice and the number of valid bits after compression
func CompressBytes(data []byte) ([]byte, uint64, error) {
	// Counts the frequency with which each byte occurs
	freq := CountFrequencies(data)
	// Build a Huffman tree
	tree := NewHuffmanTree(freq)
	// Get the Huffman coding table
	table := NewHuffmanEncTable(tree)

	return compressBytesWith(data, table)
}

// CompressFile Compress a file
// Compress the src file and write it to a DST file
//
// The compressed file format is as follows: (big-endian)
// HEADER
//   - START_FLAG						2 bytes (uint16)
//   - SRC_FILENAME_LEN					2 bytes (uint16)
//   - BYTE SIZE BEFORE COMPRESSION		4 bytes (uint32)
//   - BYTE SIZE AFTER COMPRESSION		4 bytes (uint32)
//   - SRC_FILENAME						n bytes
//
// DATA
//   - HUFFMAN TABLE
//     -- HUFFMAN TABLE SIZE 	4 bytes (uint32)
//     -- HUFFMAN TABLE DATA
//   - COMPRESSED DATA
//     -- VALID BIT LEN			4 bytes (uint32) + 1 bytes = 5 bytes
//     -- COMPRESSED BIT
//
// TAIL
//   - CRC32 CHECKSUM	  	4 bytes (uint32)
//   - END_FLAG				2 bytes (uint16)
func CompressFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	// Read in the source file content
	allSrcBytes, err := io.ReadAll(srcF)
	if err != nil {
		return err
	}

	// Perform source file compression
	freq := CountFrequencies(allSrcBytes)
	tree := NewHuffmanTree(freq)
	encTable := NewHuffmanEncTable(tree)

	compressedBytes, bitLen, err := compressBytesWith(allSrcBytes, encTable)
	if err != nil {
		return err
	}

	// Huffman Computer
	encTableSer, err := encTable.Serialize()
	if err != nil {
		return err
	}

	// Prepare to write to the target file
	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()

	originalSize := uint32(len(allSrcBytes))
	compressedSize := uint32(len(compressedBytes))
	filenameNoDir := path.Base(src)
	filenameNoDirSize := uint16(len(filenameNoDir))

	// 27 bytes is a fixed cost for the file
	var alloc uint32 = 27 + compressedSize + uint32(len(encTableSer))
	dstBytes := make([]byte, 0, alloc)
	// Write to the file header
	dstBytes = writeUint16ToBytes(CompressedFileStartFlag, dstBytes) // The file begins to be identified
	dstBytes = writeUint16ToBytes(filenameNoDirSize, dstBytes)       // The length of the file name
	dstBytes = writeUint32ToBytes(originalSize, dstBytes)            // Byte size before compression
	dstBytes = writeUint32ToBytes(compressedSize, dstBytes)          // Byte size after compression
	dstBytes = append(dstBytes, []byte(filenameNoDir)...)            // The name of the source file

	// Write to the datazone
	dstBytes = writeUint32ToBytes(uint32(len(encTableSer)), dstBytes) // Huffman computer size
	dstBytes = append(dstBytes, encTableSer...)                       // Huffman Computer

	// Calculate how many bytes will be used after compression based on the actual bit length
	bytesNeededAfterCompressed := bitLen / 8
	slot := bitLen % 8
	if slot != 0 {
		bytesNeededAfterCompressed += 1
	}

	// Use 5 bytes to record bitLen:
	// bytesNeededAfterCompressed with 4 bytes
	// slot with 1 byte
	dstBytes = writeUint32ToBytes(uint32(bytesNeededAfterCompressed), dstBytes)
	dstBytes = append(dstBytes, byte(slot))
	// dstBytes = writeUint64ToBytes(bitLen, dstBytes) // bitlen
	dstBytes = append(dstBytes, compressedBytes...) // The compressed data itself

	// Write to the end of the file
	checksum := crc32.Checksum(dstBytes, crc32q)
	dstBytes = writeUint32ToBytes(checksum, dstBytes)              // checksum
	dstBytes = writeUint16ToBytes(CompressedFileEndFlag, dstBytes) // End tag

	// Write files all at once
	n, err := dstF.Write(dstBytes)
	if err != nil {
		return err
	}

	log.Printf("successfully written %d bytes into %s\n", n, dst)

	return nil
}

func decompressBytesWith(data []byte, bitLen uint64, table HuffmanDecTable) ([]byte, error) {
	reader := NewBitsReader(data, bitLen, table)

	recovery, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return recovery, nil

}

// DecompressBytes Decompress a slice of bytes
// Input parameters include the compressed byte slice itself, the number of significant bits in the byte slice, and the Huffman decoding table
func DecompressBytes(data []byte, bitLen uint64, table HuffmanDecTable) ([]byte, error) {
	return decompressBytesWith(data, bitLen, table)
}

// DecompressFile Decompress a file
// Extract the src file and write it to a dst file
func DecompressFile(src, dst string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	srcBytes, err := io.ReadAll(srcF)
	if err != nil {
		return err
	}

	// Parse the compressed bytes of the source file
	cursor := 0
	// File header
	cursor, err = parseFileHeader(srcBytes, cursor)
	if err != nil {
		return fmt.Errorf("can not parse file header: %v", err)
	}

	// Data area
	decompressedBytes, cursor, err := parseCompressedDataArea(srcBytes, cursor)
	if err != nil {
		return fmt.Errorf("can not parse file data area: %v", err)
	}

	// Tail of file
	// Check whether the data is correct
	_, err = parseFileTail(srcBytes, cursor)
	if err != nil {
		return fmt.Errorf("can not parse file tail: %v", err)
	}

	// Create an object file to prepare for writeback
	dstF, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstF.Close()

	n, err := dstF.Write(decompressedBytes)
	if err != nil {
		return err
	}
	log.Printf("successfully written %d bytes into destination: %s\n", n, dst)

	return nil
}

// Parse the compressed file header
func parseFileHeader(srcBytes []byte, cursor int) (newCursor int, err error) {
	defer func() {
		if p := recover(); p != nil {
			// Here's a snapshot of possible slice access caused by panic caused by out-of-bounds
			newCursor = 0
			err = fmt.Errorf("%v", p)
		}
	}()

	// The file starts marking
	gotStartFlag, err := readNextUint16(srcBytes, cursor)
	if err != nil {
		return 0, err
	}
	if gotStartFlag != CompressedFileStartFlag {
		return 0, ErrInvalidStartFlag
	}
	cursor += Uint16ByteSize

	// The length of the file name before compression
	beforeFilenameLen, err := readNextUint16(srcBytes, cursor)
	if err != nil {
		return 0, err
	}
	cursor += Uint16ByteSize

	// 32-bit pre-compression file size
	_, err = readNextUint32(srcBytes, cursor)
	if err != nil {
		return 0, err
	}
	cursor += Uint32ByteSize

	// 32-bit compressed file size
	_, err = readNextUint32(srcBytes, cursor)
	if err != nil {
		return 0, err
	}
	cursor += Uint32ByteSize

	// The name of the source file
	end := cursor + int(beforeFilenameLen)
	if end > len(srcBytes) {
		return 0, ErrCursorOverflow
	}
	_ = srcBytes[cursor:end]
	cursor = end

	return cursor, nil
}

// Parse the compressed file data area
func parseCompressedDataArea(srcBytes []byte, cursor int) (data []byte, newCursor int, err error) {
	defer func() {
		if p := recover(); p != nil {
			// Here's a snapshot of possible slice access caused by panic caused by out-of-bounds
			data = nil
			newCursor = 0
			err = fmt.Errorf("%v", p)
		}
	}()

	// Huffman Computer
	huffTableLen, err := readNextUint32(srcBytes, cursor)
	if err != nil {
		return nil, 0, err
	}
	cursor += Uint32ByteSize

	decTable, err := DeserializeHuffmanDecTable(srcBytes[cursor : cursor+int(huffTableLen)])
	if err != nil {
		return nil, 0, err
	}
	cursor += int(huffTableLen)

	// Compress data parsing
	compressedBytesLen, err := readNextUint32(srcBytes, cursor)
	if err != nil {
		return nil, 0, err
	}
	cursor += Uint32ByteSize
	slot := uint8(srcBytes[cursor])
	cursor += 1

	var validBitLen uint64
	if slot == 0 {
		validBitLen = uint64(compressedBytesLen * 8)
	} else {
		validBitLen = uint64((compressedBytesLen-1)*8 + uint32(slot))
	}

	decompressedBytes, err := decompressBytesWith(srcBytes[cursor:], validBitLen, decTable)
	if err != nil {
		return nil, 0, err
	}
	cursor += int(compressedBytesLen)

	return decompressedBytes, cursor, nil
}

// Parse the end of the compressed file
func parseFileTail(srcBytes []byte, cursor int) (newCursor int, err error) {
	defer func() {
		if p := recover(); p != nil {
			newCursor = 0
			err = fmt.Errorf("%v", p)
		}
	}()

	expectedChecksum, err := readNextUint32(srcBytes, cursor)
	if err != nil {
		return 0, err
	}
	curChecksum := crc32.Checksum(srcBytes[0:cursor], crc32q)
	if curChecksum != expectedChecksum {
		return 0, ErrChecksumNotMatched
	}
	cursor += Uint32ByteSize

	// End of file marker
	gotEndFlag, err := readNextUint16(srcBytes, cursor)
	if err != nil {
		return 0, err
	}
	if gotEndFlag != CompressedFileEndFlag {
		return 0, ErrInvalidEndFlag
	}
	cursor += Uint16ByteSize

	return cursor, nil
}
