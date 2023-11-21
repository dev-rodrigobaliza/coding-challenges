package huffman

import (
	"fmt"
	"hash/crc32"
	"log"
	"strings"
)

type huffmanDeserializer interface {
	ItemNum() int
}

// Huffman Coded Table
type HuffmanEncTable map[byte]*HuffmanCode

// Huffman Decoder Table
type HuffmanDecTable map[HuffmanCode]byte

const (
	HuffmanEncTableSerStartFlag uint32 = 0x48464553 // "HFES"
	HuffmanEncTableSerEndFlag   uint32 = 0x48464545 // "HFEE"
	ChecksumPoly                       = 0xD5828281
)

const (
	MetaSize               = 4            // flag or numbder size
	MinHuffmanTableSerSize = 4 * MetaSize // bytes

	TableItemSize = 5 // bytes
)

var (
	ErrInvalidSize        = fmt.Errorf(fmt.Sprintf("len of data is small then %d", MinHuffmanTableSerSize))
	ErrInvalidStartFlag   = fmt.Errorf("start flag invalid")
	ErrInvalidEndFlag     = fmt.Errorf("end flag invalid")
	ErrCursorOverflow     = fmt.Errorf("cursor overflow")
	ErrChecksumNotMatched = fmt.Errorf("checksum not matched")
	ErrDeserialize        = fmt.Errorf("parse error")
)

var (
	crc32q = crc32.MakeTable(ChecksumPoly)
)

func NewHuffmanEncTable(tree *HuffmanTree) HuffmanEncTable {
	table := make(HuffmanEncTable, len(tree.Leaves))
	for _, leaf := range tree.Leaves {
		table[leaf.Byte] = leaf.Code
	}

	return table
}

// Get Get the encoding in the HuffmanEncTable
func (h HuffmanEncTable) Get(key byte) *HuffmanCode {
	if node, ok := h[key]; ok {
		return node
	}
	return nil
}

// ItemNum Returns the number of entries in the HuffmanEncTable
func (h HuffmanEncTable) ItemNum() int {
	return len(h)
}

// Equals Check whether the two HuffmanEncTables are the same
func (h HuffmanEncTable) Equals(other HuffmanEncTable) bool {
	if len(h) != len(other) {
		return false
	}
	for k, v := range h {
		ov, ok := other[k]
		if !ok {
			return false
		}
		if ov.Bits() != v.Bits() {
			return false
		}
	}
	return true
}

func (h HuffmanEncTable) PrettyString() string {
	prettyStringBuilder := strings.Builder{}
	prettyStringBuilder.Grow(1024)
	for k, v := range h {
		prettyStringBuilder.WriteString(fmt.Sprintf("%d(%#X)[%c]: %s(len=%d)\n", k, k, k, v.String(), len(v.String())))
	}

	return prettyStringBuilder.String()
}

func NewHuffmanDecTable(n int) HuffmanDecTable {
	return make(HuffmanDecTable, n)
}

// Get Gets the byte corresponding to an encoding
func (h HuffmanDecTable) Get(key HuffmanCode) (byte, bool) {
	v, ok := h[key]
	return v, ok
}

// ItemNum Returns the number of entries in the HuffmanDecTable
func (h HuffmanDecTable) ItemNum() int {
	return len(h)
}

// Serialize Serialize the HuffmanEncTable into a byte slice
// The serialization format is as follows (big-endian: high bits are placed at the low address, and the low bits are placed at the high address)
// START_FLAG				4 bytes
// NUMBER OF TABLE ITEMS	4 bytes (uint32)
// TABLE_ITEM_1(BYTE+CODE)	1+4=5 bytes
// TABLE_ITEM_2(BYTE+CODE)	1+4=5 bytes
// ...
// TABLE_ITEM_N(BYTE+CODE)	1+4=5 bytes
// CRC32					4 bytes
// END_FLAG					4 bytes
func (h HuffmanEncTable) Serialize() ([]byte, error) {
	n := len(h)
	size := MinHuffmanTableSerSize + 5*n
	ser := make([]byte, 0, size)

	// Write the start flag
	ser = writeUint32ToBytes(HuffmanEncTableSerStartFlag, ser)
	// Number of writes
	ser = writeUint32ToBytes(uint32(n), ser)
	// Write entries to the table
	for key, code := range h {
		ser = append(ser, key)
		ser = writeUint32ToBytes(code.AllBits(), ser)
	}
	// Write the checksum of the previous content
	checksum := crc32.Checksum(ser, crc32q)
	ser = writeUint32ToBytes(checksum, ser)

	// Write the end flag
	ser = writeUint32ToBytes(HuffmanEncTableSerEndFlag, ser)
	return ser, nil
}

func parseFlag(data []byte, cursor int, flag uint32) (int, error) {
	// Read the logo
	got, err := readNextUint32(data, cursor)
	if err != nil {
		return 0, err
	}
	if got != flag {
		return 0, fmt.Errorf("wrong flag")
	}
	cursor += MetaSize

	return cursor, nil
}

func parseStartFlag(data []byte, cursor int) (int, error) {
	// Read the start flag
	cursor, err := parseFlag(data, cursor, HuffmanEncTableSerStartFlag)
	if err != nil {
		return 0, ErrInvalidStartFlag
	}

	return cursor, nil
}

func parseEndFlag(data []byte, cursor int) (int, error) {
	// Read the start flag
	cursor, err := parseFlag(data, cursor, HuffmanEncTableSerEndFlag)
	if err != nil {
		return 0, ErrInvalidEndFlag
	}

	return cursor, nil
}

func parseItemNum(data []byte, cursor int) (int, int, error) {
	// The number of table entries read
	itemNum, err := readNextUint32(data, cursor)
	if err != nil {
		return 0, 0, err
	}
	cursor += MetaSize

	return int(itemNum), cursor, nil
}

func parseEncTable(data []byte, cursor int, itemNum int) (HuffmanEncTable, int, error) {
	retHuff := make(HuffmanEncTable, itemNum)

	for i := 0; i < int(itemNum); i++ {
		key, code, err := readNextTableItem(data, cursor)
		if err != nil {
			return nil, 0, err
		}
		cursor += TableItemSize
		retHuff[key] = &HuffmanCode{bits: code}
	}

	return retHuff, cursor, nil
}

func validateChecksum(data []byte, cursor int) (int, error) {
	expectedChecksum, err := readNextUint32(data, cursor) // The existing checksum in the data
	if err != nil {
		return 0, err
	}
	// Calculate the checksum before the cursor
	calChecksum := crc32.Checksum(data[0:cursor], crc32q)
	if expectedChecksum != calChecksum {
		log.Printf("expected checksum is %x, but got %x\n", expectedChecksum, calChecksum)
		return 0, ErrChecksumNotMatched
	}
	cursor += MetaSize

	return cursor, nil
}

func parseDecTable(data []byte, cursor int, itemNum int) (HuffmanDecTable, int, error) {
	retHuff := make(HuffmanDecTable, itemNum)

	for i := 0; i < int(itemNum); i++ {
		key, code, err := readNextTableItem(data, cursor)
		if err != nil {
			return nil, 0, err
		}
		cursor += TableItemSize
		retHuff[HuffmanCode{bits: code}] = key
	}

	return retHuff, cursor, nil
}

type huffTableItemParser func(data []byte, cursor int, itemNum int) (huffmanDeserializer, int, error)

// Deserialize byte slices
func deserialize(data []byte, parser huffTableItemParser) (huffmanDeserializer, error) {
	n := len(data)
	if n < MinHuffmanTableSerSize {
		return nil, ErrInvalidSize
	}

	// Read the start flag
	cursor := 0
	cursor, err := parseStartFlag(data, cursor)
	if err != nil {
		return nil, err
	}

	// The number of table entries read
	itemNum, cursor, err := parseItemNum(data, cursor)
	if err != nil {
		return nil, err
	}

	// Parse the entries in the table
	huffTable, cursor, err := parser(data, cursor, itemNum)
	if err != nil {
		return nil, err
	}

	// Check inspection and
	cursor, err = validateChecksum(data, cursor)
	if err != nil {
		return nil, err
	}

	// Check the end sign
	_, err = parseEndFlag(data, cursor)
	if err != nil {
		return nil, err
	}

	return huffTable, nil
}

// DeserializeHuffmanEncTable Desequence the byte slices back to the HuffmanEncTable
func DeserializeHuffmanEncTable(data []byte) (HuffmanEncTable, error) {
	huffTable, err := deserialize(data, func(data []byte, cursor, itemNum int) (huffmanDeserializer, int, error) {
		return parseEncTable(data, cursor, itemNum)
	})

	if err != nil {
		return nil, err
	}

	if encTable, ok := huffTable.(HuffmanEncTable); ok {
		return encTable, nil
	}
	return nil, ErrDeserialize
}

// DeserializeHuffmanDecTable Reverse the byte slice back to HuffmanDecTable
func DeserializeHuffmanDecTable(data []byte) (HuffmanDecTable, error) {
	huffTable, err := deserialize(data, func(data []byte, cursor, itemNum int) (huffmanDeserializer, int, error) {
		return parseDecTable(data, cursor, itemNum)
	})

	if err != nil {
		return nil, err
	}

	if decTable, ok := huffTable.(HuffmanDecTable); ok {
		return decTable, nil
	}
	return nil, ErrDeserialize
}

// readNextTableItem Read the next table entry
func readNextTableItem(buf []byte, start int) (byte, uint32, error) {
	n := len(buf)
	if n < TableItemSize {
		return 0, 0, ErrInvalidSize
	}
	if start+TableItemSize-1 > n {
		return 0, 0, ErrCursorOverflow
	}

	key := buf[start]
	code, err := readNextUint32(buf, start+1)
	if err != nil {
		return 0, 0, err
	}
	return key, code, nil
}
