package huffman

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func randInt32(min, max int32) int32 {
	return rand.Int31n(max-min) + min
}

// randomOneZeroString Randomly generated strings containing only 0s and 1s
func randomOneZeroString() string {
	rand.Seed(time.Now().Unix())

	size := randInt32(1, 24)
	var builder strings.Builder
	builder.Grow(int(size))

	for i := 0; i < int(size); i++ {
		if rand.Intn(2) == 0 {
			builder.WriteByte('0')
		} else {
			builder.WriteByte('1')
		}
	}

	return builder.String()
}

func TestHuffmanCode_Primary(t *testing.T) {
	var s = "0110101001010000001"
	require.EqualValues(t, s, NewHuffmanCodeFromString(s).String())
	require.EqualValues(t, len(s), NewHuffmanCodeFromString(s).BitLen())
}

func TestHuffmanCode(t *testing.T) {
	var testCases []string
	for i := 0; i < int(randInt32(10, 100)); i++ {
		testCases = append(testCases, randomOneZeroString())
	}
	for _, tc := range testCases {
		require.EqualValues(t, tc, NewHuffmanCodeFromString(tc).String())
		require.EqualValues(t, len(tc), NewHuffmanCodeFromString(tc).BitLen())
	}

	// Test for illegal strings
	require.Panics(t, func() { NewHuffmanCodeFromString("012301010") })
}

func reverseString(s string) string {
	b := strings.Builder{}
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}

func TestHuffmanCode_ReversePrimary(t *testing.T) {
	ss := []string{
		"100101001110110",
		"100110010100111011010101",
		"10010",
		"00101101",
		"100100010111",
	}
	for _, s := range ss {
		require.EqualValues(t, reverseString(s), NewHuffmanCodeFromString(s).ReverseNew().String())
		require.EqualValues(t, len(s), NewHuffmanCodeFromString(s).ReverseNew().BitLen())
	}
}

func TestHuffmanCode_Reverse(t *testing.T) {
	var testCases []string
	for i := 0; i < int(randInt32(1, 100)); i++ {
		testCases = append(testCases, randomOneZeroString())
	}
	for _, tc := range testCases {
		require.EqualValues(t, reverseString(tc), NewHuffmanCodeFromString(tc).ReverseNew().String())
		require.EqualValues(t, len(reverseString(tc)), NewHuffmanCodeFromString(tc).ReverseNew().BitLen())
	}

	h := &HuffmanCode{}
	require.EqualValues(t, h.BitLen(), 0)
	require.EqualValues(t, h.BitLen(), h.ReverseNew().BitLen())
}

func TestHuffmanCode_Reverse16(t *testing.T) {
	h := NewHuffmanCodeFromString("1011011101110100")
	hr := h.ReverseNew()
	require.EqualValues(t, reverseString(h.String()), hr.String())
	require.Equal(t, h.BitLen(), 16)
	require.Equal(t, hr.BitLen(), 16)

	h = NewHuffmanCodeFromString("0010111011101101")
	hr = h.ReverseNew()
	require.EqualValues(t, reverseString(h.String()), hr.String())
	require.Equal(t, h.BitLen(), 16)
	require.Equal(t, hr.BitLen(), 16)
}

func TestHuffmanClone(t *testing.T) {
	code := NewHuffmanCodeFromString("10101010101111")
	clone := code.Clone()
	require.EqualValues(t, code.Bits(), clone.Bits())
}
