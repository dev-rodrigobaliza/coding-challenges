package huffman

import (
	"fmt"
	"strings"
)

type Frequencies map[byte]uint64

func (f Frequencies) Increment(key byte) {
	f[key]++
}

type HuffmanCodeInterface interface {
	fmt.Stringer
	BitLen() int
	AppendOne()
	AppendZero()
	ReverseNew() HuffmanCodeInterface
	Clone() HuffmanCodeInterface
}

const (
	MaxHuffmanCodeBitLen = 24 // The longest encoded bits are allowed up to 24 bits
)

// A binary-encoded format with a uint32 type representing a huffman code
// The high 8-bit is the length and the low 24-bit is the code
// The highest bit of the bit itself is placed in the highest bit of the lower 24 bits in Uint32
type HuffmanCode struct {
	bits uint32
}

// NewHuffmanCodeFromString Create an object from the 01 string
func NewHuffmanCodeFromString(s string) *HuffmanCode {
	code := &HuffmanCode{}

	for _, ch := range s {
		if ch == '0' {
			code.AppendZero()
		} else if ch == '1' {
			code.AppendOne()
		} else {
			panic(fmt.Sprintf("character must be either '0' or '1', but found %c", ch))
		}
	}

	return code
}

// BitLen Returns the bit length
func (h *HuffmanCode) BitLen() int {
	// High 8th
	return int(uint32(h.bits&0xFF000000) >> 24)
}

func (h *HuffmanCode) setBitLen(l uint8) {
	high := uint32(l) << 24
	h.bits &= 0x00FFFFFF
	h.bits |= high
}

func (h *HuffmanCode) BitsLow16() uint16 {
	return uint16(h.bits & 0x0000FFFF)
}

// Bits Returns the bit itself in the form of uint32
// Discard the length of the 8 bits higher and move the highest displacement of the significant bit to start from the highest bit of uint32
func (h *HuffmanCode) Bits() uint32 {
	return (h.bits & 0x00FFFFFF) << 8
}

// Bits Returns the bit itself in the form of uint32
// However, the displacement is not carried out
func (h *HuffmanCode) BitsUntouched() uint32 {
	return (h.bits & 0x00FFFFFF)
}

// AllBits Returns the bit bit with the bitlen in the form of uint32
func (h *HuffmanCode) AllBits() uint32 {
	return h.bits
}

// Implement FMT. Stringer interface
func (h *HuffmanCode) String() string {
	bitLen := h.BitLen()
	res := strings.Builder{}
	res.Grow(bitLen)
	bits := h.Bits()
	var mask uint32 = 0x80000000

	for i := 0; i < bitLen; i++ {
		if (mask>>i)&bits == 0 {
			res.WriteByte('0')
		} else {
			res.WriteByte('1')
		}
	}

	return res.String()
}

// AppendOne Add 1 after the bit
func (h *HuffmanCode) AppendOne() {
	oldBitLen := h.BitLen()
	if oldBitLen >= MaxHuffmanCodeBitLen {
		fmt.Printf("someone try to append more bit: %s\n", h.String())
		return
	}

	shift := 23 - oldBitLen
	h.bits |= uint32(1 << shift)

	h.setBitLen(uint8(oldBitLen) + 1)
}

// AppendOne Add 0 after the bit
func (h *HuffmanCode) AppendZero() {
	oldBitLen := h.BitLen()
	if oldBitLen >= MaxHuffmanCodeBitLen {
		return
	}

	h.setBitLen(uint8(oldBitLen) + 1)
}

// ReverseNew Reverse the bits and return a new object
func (h *HuffmanCode) ReverseNew() *HuffmanCode {
	clone := &HuffmanCode{}
	bitLen := h.BitLen()
	if bitLen == 0 {
		return clone
	}

	bits := h.BitsUntouched()
	bits >>= (24 - bitLen)
	var mask uint32 = 0x1

	for i := 0; i < bitLen; i++ {
		if (bits>>i)&mask == 0 {
			clone.AppendZero()
		} else {
			clone.AppendOne()
		}
	}

	return clone
}

// Clone Returns a copy of the object
func (h *HuffmanCode) Clone() *HuffmanCode {
	return &HuffmanCode{
		bits: h.bits,
	}
}
