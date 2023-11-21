package huffman

import (
	// "fmt"
	"testing"
	// "github.com/stretchr/testify/require"
)

// Test creating a Huffman tree
func TestConstructHuffmanTree(t *testing.T) {
	testCases := []*struct {
		freq   Frequencies
		expect map[byte]string
	}{
		{
			freq: Frequencies{'a': 2, 'b': 2, 'c': 2, 'e': 2, 'f': 1},
			// Huffman trees are not unique (when the same weight, different nodes as left and right subtrees will cause inconsistent results)
			// Although the structure of Huffman trees is not unique, their weighted path sum is the same
			expect: map[byte]string{'a': "00", 'b': "111", 'c': "10", 'e': "01", 'f': "110"},
		},
		{
			freq:   Frequencies{'a': 4, 'b': 1, 'c': 6, 'd': 8, 'e': 3},
			expect: map[byte]string{'d': "0", 'e': "1101", 'a': "111", 'b': "1100", 'c': "10"},
		},
		{
			freq: Frequencies{' ': 20, 'a': 40, 'm': 10, 'l': 7, 'f': 8, 't': 15},
		},
		{
			freq: Frequencies{'i': 20},
		},
	}

	for _, tc := range testCases {
		_, _ = ConstructHuffmanTree(tc.freq)
		// for _, leaf := range leaves {
		// require.EqualValues(t, tc.expect[leaf.Byte], leaf.Code.String())
		// }
		// fmt.Println()
	}
}

func TestConstructHuffmanTree_AllBytes256(t *testing.T) {
	// Generate all possible bytes (0~255)
	// tests := make([]byte, 0, 256)
	var freq Frequencies = make(Frequencies, 256)
	var total uint64 = 0

	var comple uint64 = 100000
	for i := 0; i <= 255; i++ {
		cnt := uint64(i) + 1
		if i == 255 {
			cnt += comple
		}
		total += cnt
		freq[(byte)(i)] = cnt
	}

	ConstructHuffmanTree(freq)
	// for _, leaf := range leaves {
	// 	fmt.Printf("%d: %s, %d\n", leaf.Byte, leaf.Code.String(), len(leaf.Code.String()))
	// }
	// fmt.Println()
}

func TestConstructHuffmanTreeFreq_AllFreqEqual(t *testing.T) {
	// Generate all possible bytes (0~255)
	// tests := make([]byte, 0, 256)
	var freq Frequencies = make(Frequencies, 256)
	var total uint64 = 0
	for i := 0; i <= 255; i++ {
		var cnt uint64 = 1
		total += cnt
		freq[(byte)(i)] = cnt
	}

	ConstructHuffmanTree(freq)
	// for _, _ := range leaves {
	// 	// fmt.Printf("%d: %s, %d\n", leaf.Byte, leaf.Code.String(), len(leaf.Code.String()))
	// }
	// fmt.Println()
}
