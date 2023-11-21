package main

import (
	"fmt"
	"jp/parser"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("invalid usage, call jp <filename.json>")
		os.Exit(2)
	}

	buf, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("failed to read file %q: %v", os.Args[1], err)
		os.Exit(2)
	}

	if !parser.IsValid(string(buf), true) {
		fmt.Println("invalid json file")
		os.Exit(1)
	}

	fmt.Println("valid json file")
}
