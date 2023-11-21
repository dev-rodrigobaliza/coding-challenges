package main

import (
	"compressor/huffman"
	"flag"
	"fmt"
	"os"
)

func main() {
	performCompress := flag.Bool("compress", false, "compress given file")
	performDecompress := flag.Bool("decompress", false, "decompress given file")
	inputFile := flag.String("input", "", "input filename")
	outputFile := flag.String("output", "", "output filename")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("please specify the input filename")
		os.Exit(0)
	}
	if *outputFile == "" {
		fmt.Println("please specify the output filename")
		os.Exit(0)
	}
	if *performCompress && *performDecompress || !*performCompress && !*performDecompress {
		fmt.Println("compress flag or decompress should one set to true")
		os.Exit(0)
	}

	if *performCompress {
		fmt.Println("performing compression...")
		err := huffman.CompressFile(*inputFile, *outputFile)
		if err != nil {
			fmt.Printf("compression failed: %v\n", err)
		} else {
			fmt.Println("compression ok")
		}
	}

	if *performDecompress {
		fmt.Println("performing decompression...")
		err := huffman.DecompressFile(*inputFile, *outputFile)
		if err != nil {
			fmt.Printf("decompression failed: %v\n", err)
		} else {
			fmt.Println("decompression ok")
		}
	}
}
