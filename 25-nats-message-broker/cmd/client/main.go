package main

import (
	"bytes"
	"fmt"
	"net"
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

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	log(buffer[:n])

	// connect
	if _, err := conn.Write([]byte("connect {}" + end)); err != nil {
		panic(err)
	}

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	log(buffer[:n])

	// ping
	if _, err := conn.Write([]byte("ping" + end)); err != nil {
		panic(err)
	}

	n, err = conn.Read(buffer)
	if err != nil {
		panic(err)
	}

	log(buffer[:n])

	// // sub
	// if _, err := conn.Write([]byte("sub foo 1" + end)); err != nil {
	// 	panic(err)
	// }

	// n, err = conn.Read(buffer)
	// if err != nil {
	// 	panic(err)
	// }

	// log(buffer[:n])

	// // pub
	// if _, err := conn.Write([]byte("pub foo 3" + end + "bar" + end)); err != nil {
	// 	panic(err)
	// }

	// n, err = conn.Read(buffer)
	// if err != nil {
	// 	panic(err)
	// }

	// log(buffer[:n])

	// // get last pub
	// n, err = conn.Read(buffer)
	// if err != nil {
	// 	panic(err)
	// }

	// log(buffer[:n])
}

func log(msg []byte) {
	m := bytes.Split(msg, []byte{13, 10})

	for _, m := range m {
		if len(m) > 0 {
			fmt.Println(string(m))
		}
	}

	time.Sleep(time.Millisecond * 500)
}
