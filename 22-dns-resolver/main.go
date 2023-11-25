package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	id     = 22
	server = "8.8.8.8:53"
)

type address struct {
	Type string
	Host string
}

func main() {
	host := flag.String("host", "dns.google.com", "host name to query ip address")

	flag.Parse()

	if *host == "" {
		fmt.Println("please specify a valid host name")
		os.Exit(1)
	}

	req := messageBuilder(*host)
	res := askDNS(req)
	addrs := parseResponse(req, res)

	fmt.Printf("%q ips:\n", *host)
	for _, a := range addrs {
		fmt.Printf("\t%+v\n", a)
	}
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

func parseResponse(req, res []byte) []address {
	if !validateID(req, res) {
		fmt.Printf("failed to validate id from the response, want: %x, got: %x\n", req[:4], res[:4])
		os.Exit(1)
	}
	// get count of responses
	count := binary.BigEndian.Uint16(res[6:10])
	// remove header
	res = res[12:]
	// remove question
	pos := bytes.Index(res, []byte{0, 0})
	pos += 5
	// get each response
	addrs := make([]address, 0, int(count))
	for i := 0; i < int(count); i++ {
		r := res[pos:]
		// skip pointer
		p := 2
		// get type
		t := r[p : p+2]
		p += 2
		// get class
		//c := r[p : p+2]
		p += 2
		// get ttl
		//ttl := r[p : p+4]
		p += 4
		// get size
		s := r[p : p+2]
		p += 2
		// get addr
		size := int(binary.BigEndian.Uint16(s))
		raw := r[p : p+size]

		addrs = append(addrs, parseAddress(t, raw))
		// adjust pos
		pos += p + size
	}

	return addrs
}

func validateID(req, res []byte) bool {
	for i := 0; i < 2; i++ {
		if req[i] != res[i] {
			return false
		}
	}

	return true
}

func parseAddress(t, raw []byte) address {
	var a address
	switch t[1] {
	case 1:
		a.Type = "A record"
		a.Host = parseIP(raw)

	case 2:
		a.Type = "name server"

	case 5:
		a.Type = "CNAME"
		a.Host = parseHost(raw)

	case 0xf:
		a.Type = "mail server"

	default:
		a.Type = "unknown"
	}

	return a
}

func parseIP(raw []byte) string {
	switch len(raw) {
	case 4:
		return fmt.Sprintf("%d.%d.%d.%d", raw[0], raw[1], raw[2], raw[3])

	case 8:
		return fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x", raw[0:2], raw[2:4], raw[4:6], raw[6:8], raw[8:10], raw[10:12], raw[12:14], raw[14:16])

	default:
		return "invalid ip address"
	}
}

func parseHost(raw []byte) string {
	b := strings.Builder{}
	pos := 0
	for {
		if raw[pos] == 0 {
			break
		}
		if pos > 0 {
			b.WriteString(".")
		}

		s := int(raw[pos])
		pos++
		b.Write(raw[pos : pos+s])
		pos += s
	}

	return b.String()
	//07636e6e2d746c73 036d6170 06666173746c79 036e6574 00
}
