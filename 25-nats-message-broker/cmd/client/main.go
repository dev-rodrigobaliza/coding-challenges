package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	end = "\r\n"
)

func main() {
	// dial
	conn, err := net.Dial("tcp", ":4222")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	go handleResponses(conn)
	time.Sleep(time.Second)

	// connect
	if _, err := conn.Write([]byte("connect {}" + end)); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	// ping
	if _, err := conn.Write([]byte("ping" + end)); err != nil {
		panic(err)
	}
	time.Sleep(time.Second)

	// // sub
	// if _, err := conn.Write([]byte("sub foo 1" + end)); err != nil {
	// 	panic(err)
	// }

	// // pub
	// if _, err := conn.Write([]byte("pub foo 3" + end + "bar" + end)); err != nil {
	// 	panic(err)
	// }
}

func handleResponses(conn net.Conn) {
	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		res := scanner.Text()
		msgs := strings.Split(res, end)
		for _, msg := range msgs {
			if msg != "" {
				fmt.Println(scanner.Text())
			}
		}
	}
	if err := scanner.Err(); err != nil {
		if !errors.Is(err, io.EOF) {
			fmt.Println("failed to read bytes from the connection")
			panic(err)
		}
	}
}
