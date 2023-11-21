package huffman

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func constructPrerequisite(t *testing.T, data []byte) ([]byte, uint64, HuffmanDecTable) {
	// Counts the frequency with which each byte occurs
	freq := CountFrequencies(data)
	// Build a Huffman tree
	tree := NewHuffmanTree(freq)
	// Get the Huffman coding table
	table := NewHuffmanEncTable(tree)

	denseBytes, totalBits, err := compressBytesWith(data, table)
	require.Nil(t, err)

	ser, err := table.Serialize()
	require.Nil(t, err)

	decTable, err := DeserializeHuffmanDecTable(ser)
	require.Nil(t, err)

	return denseBytes, totalBits, decTable
}

func TestReader_ReadByte(t *testing.T) {
	f, err := os.Open("test/test_data.txt")
	require.Nil(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.Nil(t, err)

	denseBytes, totalBits, decTable := constructPrerequisite(t, data)

	// Read the compressed content
	reader := NewBitsReader(denseBytes, totalBits, decTable)
	firstByte, err := reader.ReadByte()
	require.Nil(t, err)
	require.Equal(t, firstByte, data[0])

	secondByte, err := reader.ReadByte()
	require.Nil(t, err)
	require.Equal(t, secondByte, data[1])
}

func TestBitsReader_ReadAll(t *testing.T) {
	f, err := os.Open("test/test_data.txt")
	require.Nil(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.Nil(t, err)

	denseBytes, totalBits, decTable := constructPrerequisite(t, data)

	// Read the compressed content
	reader := NewBitsReader(denseBytes, totalBits, decTable)
	dataRecovered, err := reader.ReadAll()
	require.Nil(t, err)
	require.Equal(t, len(data), len(dataRecovered))
	require.EqualValues(t, data, dataRecovered)
	fmt.Printf("len(data)=%d, len(dataComplete)=%d\n", len(data), len(dataRecovered))
}
