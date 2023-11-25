package main

import (
	"fmt"
	"net"
	"time"
)

const (
	end = "\r\n"
)

func main() {
	conn, err := net.Dial("tcp", ":11211")
	if err != nil {
		panic(err)
	}

	// set test
	if _, err := conn.Write([]byte("set test 0 1 4" + end + "azul" + end)); err != nil {
		panic(err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buffer[:n]))

	time.Sleep(time.Second*2)

	// get test
	if _, err := conn.Write([]byte("get test" + end)); err != nil {
		panic(err)
	}

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buffer[:n]))

	// set test
	if _, err := conn.Write([]byte("set test 0 100 4" + end + "anil" + end)); err != nil {
		panic(err)
	}

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buffer[:n]))

	time.Sleep(time.Second*2)

	// get test
	if _, err := conn.Write([]byte("get test" + end)); err != nil {
		panic(err)
	}

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buffer[:n]))

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buffer[:n]))

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(buffer[:n]))
}
