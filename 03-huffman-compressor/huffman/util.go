package huffman

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	Uint16ByteSize = 2 // bytes
	Uint32ByteSize = 4 // bytes
	Uint64ByteSize = 8 // bytes
)

// CountFrequencies Count the frequency of occurrence of each byte in the input byte slice
func CountFrequencies(data []byte) Frequencies {
	frequency := make(Frequencies)
	for _, d := range data {
		frequency[d] += 1
	}

	return frequency
}

// BytesToString Convert the total bitsize bits in data to 01 strings
func BytesToString(data []byte, bitsize int) string {
	builder := strings.Builder{}
	builder.Grow(bitsize)
	var scannedLen int = 0
	var mask byte = 1 << 7

	for i := 0; i < len(data); i++ {
		var curSlot int = 0
		bits := data[i]
		for curSlot < 8 && scannedLen < bitsize {
			if (bits & mask) == 0 {
				builder.WriteByte('0')
			} else {
				builder.WriteByte('1')
			}
			bits <<= 1
			curSlot++
			scannedLen++
		}
	}

	return builder.String()
}

// writeUint32ToBytes writes 4 bytes into the buf
func writeUint32ToBytes(in uint32, buf []byte) []byte {
	// Big-endian
	buf = append(buf, byte((in&0xFF000000)>>24))
	buf = append(buf, byte((in&0x00FF0000)>>16))
	buf = append(buf, byte((in&0x0000FF00)>>8))
	buf = append(buf, byte(in&0x000000FF))

	return buf
}

// writeUint16ToBytes writes 2 bytes into the buf
func writeUint16ToBytes(in uint16, buf []byte) []byte {
	// Big-endian
	buf = append(buf, byte((in&0xFF00)>>8))
	buf = append(buf, byte((in & 0x00FF)))

	return buf
}

// writeUint64ToBytes writes 8 bytes into the buf
func writeUint64ToBytes(in uint64, buf []byte) []byte {
	// Big-endian
	buf = append(buf, byte((in&0xFF00000000000000)>>56))
	buf = append(buf, byte((in&0x00FF000000000000)>>48))
	buf = append(buf, byte((in&0x0000FF0000000000)>>40))
	buf = append(buf, byte((in&0x000000FF00000000)>>32))
	buf = append(buf, byte((in&0x00000000FF000000)>>24))
	buf = append(buf, byte((in&0x0000000000FF0000)>>16))
	buf = append(buf, byte((in&0x000000000000FF00)>>8))
	buf = append(buf, byte(in&0x00000000000000FF))

	return buf
}

// readNextUint32 Read the next UINT32 from the byte slice buf, starting from the index start of the byte slice
func readNextUint32(buf []byte, start int) (uint32, error) {
	n := len(buf)
	if n < Uint32ByteSize {
		return 0, ErrInvalidSize
	}
	if start+Uint32ByteSize-1 > n {
		return 0, ErrCursorOverflow
	}

	var ans uint32 = 0

	byte0 := buf[start]
	byte1 := buf[start+1]
	byte2 := buf[start+2]
	byte3 := buf[start+3]

	ans |= uint32(byte0) << 24
	ans |= uint32(byte1) << 16
	ans |= uint32(byte2) << 8
	ans |= uint32(byte3)

	return ans, nil
}

// readNextUint16 Read the next uint16 from the byte slice, starting from the index start of the byte slice
func readNextUint16(buf []byte, start int) (uint16, error) {
	n := len(buf)
	if n < Uint16ByteSize {
		return 0, ErrInvalidSize
	}
	if start+Uint16ByteSize-1 > n {
		return 0, ErrCursorOverflow
	}

	var ans uint16 = 0

	byte0 := buf[start]
	byte1 := buf[start+1]

	ans |= uint16(byte0) << 8
	ans |= uint16(byte1)

	return ans, nil
}

// readNextUint64 Read the next UINT64 from the byte slice buf, starting from the index start of the byte slice
func readNextUint64(buf []byte, start int) (uint64, error) {
	n := len(buf)
	if n < Uint64ByteSize {
		return 0, ErrInvalidSize
	}
	if start+Uint64ByteSize-1 > n {
		return 0, ErrCursorOverflow
	}

	var ans uint64 = 0

	byte0 := buf[start]
	byte1 := buf[start+1]
	byte2 := buf[start+2]
	byte3 := buf[start+3]
	byte4 := buf[start+4]
	byte5 := buf[start+5]
	byte6 := buf[start+6]
	byte7 := buf[start+7]

	ans |= uint64(byte0) << 56
	ans |= uint64(byte1) << 48
	ans |= uint64(byte2) << 40
	ans |= uint64(byte3) << 32
	ans |= uint64(byte4) << 24
	ans |= uint64(byte5) << 16
	ans |= uint64(byte6) << 8
	ans |= uint64(byte7)

	return ans, nil
}

// Sha256SumFile Calculate the SHA256 of a file
func Sha256SumFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
