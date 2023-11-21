package huffman

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompressBytes(t *testing.T) {
	f, err := os.Open("test/test_data.txt")
	require.Nil(t, err)
	defer f.Close()

	data, err := io.ReadAll(f)
	require.Nil(t, err)

	denseData, _, err := CompressBytes(data)
	require.Nil(t, err)
	fmt.Printf("data size = %d, denseData size = %d\n", len(data), len(denseData))
}

var (
	s1 = "test/test_data.txt"
	s2 = "test/test_data.bin"
	s3 = "test/test_data_recover.txt"
)

func TestSaveHuffmanTablePretty(t *testing.T) {
	srcF, err := os.Open(s1)
	require.Nil(t, err)
	defer srcF.Close()

	allSrcBytes, err := io.ReadAll(srcF)
	require.Nil(t, err)

	freq := CountFrequencies(allSrcBytes)
	// fmt.Printf("len(freq) = %d\n", len(freq))
	tree := NewHuffmanTree(freq)

	encTable := NewHuffmanEncTable(tree)

	prettyEncString := encTable.PrettyString()
	outname := "test/test_data-huffman.txt"
	psF, err := os.Create(outname)
	require.Nil(t, err)
	io.Copy(psF, strings.NewReader(prettyEncString))
	psF.Close()
	os.Remove(outname)
}

func TestCompressFile(t *testing.T) {
	err := CompressFile(s1, s2)
	if err != nil {
		log.Println(err)
	}
}

func TestDecompressFile(t *testing.T) {
	err := DecompressFile(s2, s3)
	if err != nil {
		log.Println(err)
	}
}

func TestCompressAndDecompress(t *testing.T) {
	filenames := []string{
		"test/test_data.txt",
		"test/test_data2.txt",
	}
	for _, filename := range filenames {
		log.Printf("performing compressing %s...\n", filename)
		binname := filename + ".bin"
		err := CompressFile(filename, binname)
		require.Nil(t, err)

		log.Printf("performing decompressing %s...\n", binname)
		recovername := filename + ".recover"
		err = DecompressFile(binname, recovername)
		require.Nil(t, err)

		// Check whether the md5 of the source file is the same as that of the file that is compressed and then decompressed
		originalHash, err := Sha256SumFile(filename)
		require.Nil(t, err)
		afterHash, err := Sha256SumFile(recovername)
		require.Nil(t, err)
		require.Equal(t, originalHash, afterHash)
		os.Remove(binname)
		os.Remove(recovername)
	}
}
