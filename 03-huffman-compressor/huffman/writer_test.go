package huffman

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBitsWriter_Uint32_Primary(t *testing.T) {

	s1 := "01010101"
	s2 := "101111"
	s3 := "1111001111"
	s4 := "111"

	code1 := NewHuffmanCodeFromString(s1)
	code2 := NewHuffmanCodeFromString(s2)
	code3 := NewHuffmanCodeFromString(s3)
	code4 := NewHuffmanCodeFromString(s4)

	w := NewBitsWriter()
	w.WriteUint32(code1.Bits(), uint8(code1.BitLen()))
	w.WriteUint32(code2.Bits(), uint8(code2.BitLen()))
	w.WriteUint32(code3.Bits(), uint8(code3.BitLen()))
	w.WriteUint32(code4.Bits(), uint8(code4.BitLen()))

	b := w.Buf()
	// 01010101 10111111 11001111 11100000
	// 0x55		0xbf	 0xcf	  0xe0
	s := BytesToString(b, code1.BitLen()+code2.BitLen()+code3.BitLen()+code4.BitLen())
	req := s1 + s2 + s3 + s4
	fmt.Println("got :", s)
	fmt.Println("want:", req)
	require.EqualValues(t, s, req)

}

func TestBitsWriter_Uint32(t *testing.T) {
	tests := []struct {
		Codes []string
	}{
		{Codes: []string{"01010101", "11111111", "00000", "111101"}},
		{Codes: []string{"10101010", "00001111", "1111111111111110", "10101010", "0101010111110000"}},
		{Codes: []string{"1010101010001101", "10101010", "0000000011111111"}},
		{Codes: []string{"1", "0", "0", "1", "0", "1", "0", "1", "010111"}},
		{Codes: []string{"010101", "01010101", "1111110111010101", "0101"}},
		{Codes: []string{"1010101010", "0101010101", "01001100"}},
		{Codes: []string{"0000", "000000001", "111010", "01011010101010"}},
		{Codes: []string{"111111111100", "000", "00000000", "1111111111000000"}},
		{Codes: []string{"1", "0", "1", "0", "0101010101", "0001111"}},
		{Codes: []string{"10101010010111111000110", "101010101011111110000", "10101010111100001010", "101010101010110000101110", "101010101010101010111111", "0001111"}},
	}

	for _, tc := range tests {
		w := NewBitsWriter()
		expect := strings.Builder{}
		var total int = 0
		for _, s := range tc.Codes {
			expect.WriteString(s)
			huff := NewHuffmanCodeFromString(s)
			bitlen := uint8(huff.BitLen())
			total += int(bitlen)
			w.WriteUint32(huff.Bits(), bitlen)
		}
		b := w.Buf()

		require.EqualValues(t, expect.String(), BytesToString(b, total))
	}
}

func TestBitsWriter_Uint32_WithError(t *testing.T) {
	w := NewBitsWriter()
	require.NotNil(t, w.WriteUint32(0x9600, 40))
}
