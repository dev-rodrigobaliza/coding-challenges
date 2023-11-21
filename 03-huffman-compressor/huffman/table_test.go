package huffman

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHuffmanTable_SerializeAndDeserialize(t *testing.T) {
	f, err := os.Open("test/test_data.txt")
	require.Nil(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.Nil(t, err)

	// Get the HuffmanTable
	freq := CountFrequencies(data)
	tree := NewHuffmanTree(freq)
	table := NewHuffmanEncTable(tree)
	// serialization
	ser, err := table.Serialize()
	require.Nil(t, err)

	tableF, err := os.Create("test/table.bin")
	require.Nil(t, err)
	defer os.Remove("test/table.bin")
	defer tableF.Close()

	n, err := tableF.Write(ser)
	require.Nil(t, err)
	require.EqualValues(t, n, len(ser))

	// Deserialization
	deTable, err := DeserializeHuffmanEncTable(ser)
	require.Nil(t, err)
	require.Equal(t, deTable.ItemNum(), table.ItemNum())
	require.True(t, deTable.Equals(table))
	require.True(t, table.Equals(deTable))
}
