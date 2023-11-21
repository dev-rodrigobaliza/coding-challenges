package huffman

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountFrequencies(t *testing.T) {

	testCases := []*struct {
		Data   []byte
		Expect map[byte]uint64
	}{
		{[]byte("abcdef"), map[byte]uint64{'a': 1, 'b': 1, 'c': 1, 'd': 1, 'e': 1, 'f': 1}},
		{[]byte("aabbcceef"), map[byte]uint64{'a': 2, 'b': 2, 'c': 2, 'e': 2, 'f': 1}},
	}

	for _, tc := range testCases {
		have := CountFrequencies(tc.Data)
		require.EqualValues(t, tc.Expect, have)
	}
}
