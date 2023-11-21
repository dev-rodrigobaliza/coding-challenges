package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func main() {
	var (
		cmd     string
		file    string
		scanner *bufio.Scanner
		b       bool
		l       bool
		w       bool
		c       bool
	)

	switch len(os.Args) {
    case 1:
        cmd = ""
        file = ""

	case 2:
		if strings.HasPrefix(os.Args[1], "-") {
			cmd = os.Args[1]
		} else {
			file = os.Args[1]
		}

	case 3:
		cmd = os.Args[1]
		file = os.Args[2]

	default:
		println("invalid usage!")
		return
	}

	if file != "" {
		fileExists(file)

		f, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Printf("open file error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()

		reader := bufio.NewReader(f)
		scanner = bufio.NewScanner(reader)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	switch cmd {
	case "-b":
		b = true

	case "-l":
		l = true

	case "-w":
		w = true

	case "-c":
		c = true

	default:
		b = true
		l = true
		w = true
		c = true
	}

	count(scanner, b, l, w, c)

	fmt.Printf("%s\n", file)
}

func count(scanner *bufio.Scanner, b, l, w, c bool) {
	var (
		s     string
		bSize int
		lSize int
		wSize int
		cSize int
	)

	for scanner.Scan() {
        buf := scanner.Bytes()

		if l {
			lSize += 1
		}

		s += string(buf)
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if b {
		bSize = len(s)
		fmt.Printf("%d ", bSize)
	}

	if l {
		fmt.Printf("%d ", lSize)
	}

	if w {
        wSize = len(strings.Fields(s))
		fmt.Printf("%d ", wSize)
	}

	if c {
        cSize += len([]rune(s))
		fmt.Printf("%d ", cSize)
	}
}

func fileExists(file string) fs.FileInfo {
	fs, err := os.Stat(file)
	if os.IsNotExist(err) {
		println("file do not exists")
		os.Exit(1)
	}

	return fs
}
