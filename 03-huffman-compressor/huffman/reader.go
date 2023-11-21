package huffman

import "fmt"

// BitsReader Define how bits are read
type BitsReader struct {
	buf    []byte
	table  HuffmanDecTable
	index  int
	cursor uint64
	remain uint64
}

var (
	ErrBitCodeNotFound = fmt.Errorf("bitcode not found")
	ErrBitsExhausted   = fmt.Errorf("bit exhausted")
)

func NewBitsReader(buf []byte, bitLen uint64, decodeTable HuffmanDecTable) *BitsReader {
	return &BitsReader{
		buf:    buf,
		table:  decodeTable,
		index:  0,
		cursor: 0,
		remain: bitLen,
	}
}

// ReadByte A byte is parsed from the bit
func (r *BitsReader) ReadByte() (byte, error) {
	if r.remain == 0 {
		return 0, ErrBitsExhausted
	}

	// Start reading from the cursor bit of the first byte of the index
	parsedCode := HuffmanCode{}
	var ret byte = 0

	i := 0
	foundOne := false

	// Maximum read at a time MaxHuffmanCodeBitLen bit
	for ; i < MaxHuffmanCodeBitLen && r.remain > 0; i++ {
		if r.nextBit() {
			parsedCode.AppendOne()
		} else {
			parsedCode.AppendZero()
		}

		r.cursor = (r.cursor + 1) % 8
		r.remain--
		if r.cursor == 0 {
			r.index++
		}

		// Look for the encoding in the decode table
		key, ok := r.table.Get(parsedCode)
		if ok {
			ret = key
			foundOne = true
			break
		}
	}

	if i == MaxHuffmanCodeBitLen {
		// A non-existent bit encoding was discovered
		return 0, ErrBitCodeNotFound
	}

	if r.remain == 0 && !foundOne {
		return 0, ErrBitsExhausted
	}

	return ret, nil
}

// ReadAll Parse all bits
func (r *BitsReader) ReadAll() ([]byte, error) {
	approxLen := r.remain / 8
	ret := make([]byte, 0, approxLen)

	for r.remain > 0 {
		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		ret = append(ret, b)
	}

	return ret, nil
}

// Determine the next bit (the first cursor bit of the index byte
// Returning true indicates bit 1 and false indicates bit 0
func (r *BitsReader) nextBit() bool {
	mask := byte(0x80 >> r.cursor)
	return (r.buf[r.index] & mask) == mask
}
