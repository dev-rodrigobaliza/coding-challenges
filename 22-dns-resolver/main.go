package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

const (
	id     = 22
	server = "8.8.8.8:53"
)

func main() {
	req := messageBuilder("dns.google.com")
	res := askDNS(req)

	if !validateID(req, res) {
		fmt.Printf("failed to validate id from the response, want: %x, got: %x\n", req[:4], res[:4])
		return
	}

	fmt.Println("ok")
	fmt.Printf("%x", res)
}

func messageBuilder(hostnames ...string) []byte {
	if len(hostnames) == 0 {
		return nil
	}

	b := new(bytes.Buffer)
	// id
	writeInt16(b, id)
	// flags (fixed for now)
	writeInt16(b, 256)
	// number of questions
	writeInt16(b, len(hostnames))
	// number of answer resources (fixed for now)
	writeInt16(b, 0)
	// number of authority resource records (fixed for now)
	writeInt16(b, 0)
	// number of additional resource records (fixed for now)
	writeInt16(b, 0)
	// encodedQuestion
	for _, h := range hostnames {
		parts := strings.Split(h, ".")
		if len(parts) < 2 {
			return nil
		}
		for _, p := range parts {
			writeInt8(b, len(p))
			b.WriteString(p)
		}
		writeInt8(b, 0)
	}
	// query type (fixed for now)
	writeInt16(b, 1)
	// query class (fixed for now)
	writeInt16(b, 1)

	return b.Bytes()
}

func writeInt8(b *bytes.Buffer, val int) {
	binary.Write(b, binary.BigEndian, uint8(val))
}

func writeInt16(b *bytes.Buffer, val int) {
	binary.Write(b, binary.BigEndian, uint16(val))
}

func askDNS(req []byte) []byte {
	conn, err := net.Dial("udp", server)
	if err != nil {
		fmt.Printf("failed to dial: %v\n", err)
		return nil
	}
	defer conn.Close()

	if _, err := conn.Write(req); err != nil {
		fmt.Printf("failed to send req: %v\n", err)
		return nil
	}

	res := make([]byte, 2048)
	n, err := conn.Read(res)
	if err != nil {
		fmt.Printf("failed to read resp: %v\n", err)
		return nil
	}

	return res[:n]
}

func validateID(req, res []byte) bool {
	for i := 0; i < 2; i++ {
		if req[i] != res[i] {
			return false
		}
	}

	return true
}