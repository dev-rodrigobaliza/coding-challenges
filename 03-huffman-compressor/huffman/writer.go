package huffman

import (
	"errors"
)

const (
	MaxUint16Len = 16
	MaxUint32Len = 32
)

var (
	ErrMaxLenExceeded = errors.New("maximum length exceeded")
)

// BitsWriter Defines how bits are written
type BitsWriter struct {
	// The buffer where the bitstream is stored
	buf []byte
	// Current byte index
	idx uint64
	// Current bit index
	slot uint64
}

func NewBitsWriter() *BitsWriter {
	return &BitsWriter{
		buf:  make([]byte, 1),
		idx:  0,
		slot: 0,
	}
}

func (w *BitsWriter) updateCursor() {
	w.slot = (w.slot + 1) % 8
	if w.slot == 0 {
		w.idx++
		w.buf = append(w.buf, 0)
	}
}

func (w *BitsWriter) appendOne() {
	mask := uint8(1 << (7 - w.slot))
	w.buf[w.idx] |= mask
	w.updateCursor()
}

func (w *BitsWriter) appendZero() {
	w.updateCursor()
}

// WriteUint16 Write n-bit bits, starting from the high bit, up to a maximum of 16 bits
func (w *BitsWriter) WriteUint16(a uint16, n uint8) error {
	if n > MaxUint16Len {
		return ErrMaxLenExceeded
	}

	// Special Circumstances Handling
	if w.slot == 0 && (n == 8 || n == 16) {
		high := byte((a & 0xFF00) >> 8)
		if n == 8 {
			w.buf = append(w.buf, 0, 0)
			w.buf[w.idx] = high
			w.idx += 1
		} else if n == 16 {
			low := byte(a & 0x00FF)
			w.buf = append(w.buf, 0, 0, 0)
			w.buf[w.idx] = high
			w.idx += 1
			w.buf[w.idx] = low
			w.idx += 1
		}
		return nil
	}

	// slot != 0 or n % 8 != 0
	var mask uint16 = 0x8000
	for i := 0; i < int(n); i++ {
		if (a<<i)&mask == 0 {
			w.appendZero()
		} else {
			w.appendOne()
		}
	}

	return nil
}

// WriteUint32 Write n-bit bits, starting from the high bit, up to a maximum of 32 bits
func (w *BitsWriter) WriteUint32(a uint32, n uint8) error {
	if n > MaxUint32Len {
		return ErrMaxLenExceeded
	}

	// In special cases, there may be exactly 1, 2, 3, or 4 bytes
	if w.slot == 0 && (n == 8 || n == 16 || n == 24 || n == 32) {
		high := byte((a & 0xFF000000) >> 24) // 32 bits in the upper eighth
		if n == 8 {
			w.buf = append(w.buf, 0, 0)
			w.buf[w.idx] = high
			w.idx += 1
		} else if n == 16 {
			middleLeft := byte((a & 0x00FF0000) >> 16)
			w.buf = append(w.buf, 0, 0, 0)
			w.buf[w.idx] = high
			w.idx += 1
			w.buf[w.idx] = middleLeft
			w.idx += 1
		} else if n == 24 {
			middleLeft := byte((a & 0x00FF0000) >> 16)
			middleRight := byte((a & 0x0000FF00) >> 8)
			w.buf = append(w.buf, 0, 0, 0, 0)
			w.buf[w.idx] = high
			w.idx += 1
			w.buf[w.idx] = middleLeft
			w.idx += 1
			w.buf[w.idx] = middleRight
			w.idx += 1
		} else { // n == 32
			middleLeft := byte((a & 0x00FF0000) >> 16)
			middleRight := byte((a & 0x0000FF00) >> 8)
			low := byte(a & 0x000000FF)
			w.buf = append(w.buf, 0, 0, 0, 0, 0)
			w.buf[w.idx] = high
			w.idx += 1
			w.buf[w.idx] = middleLeft
			w.idx += 1
			w.buf[w.idx] = middleRight
			w.idx += 1
			w.buf[w.idx] = low
			w.idx += 1
		}
		return nil
	}

	// In general, n bits are written
	var mask uint32 = 0x80000000
	for i := 0; i < int(n); i++ {
		if (a<<i)&mask == 0 {
			w.appendZero()
		} else {
			w.appendOne()
		}
	}

	return nil
}

// Buf Returns a copy of the underlying bit buffer
func (w *BitsWriter) Buf() []byte {
	// Only valid data needs to be copied, not the entire w.buf buffer, otherwise the data will contain a large number of invalid zeros
	var validByteLen int
	if w.slot == 0 {
		validByteLen = int(w.idx)
	} else {
		validByteLen = int(w.idx) + 1
	}
	cp := make([]byte, validByteLen)
	copy(cp, w.buf)

	return cp
}
