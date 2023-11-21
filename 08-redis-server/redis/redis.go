package redis

import (
	"errors"
	"fmt"
	"net"
	"rs/serde"
	"strings"
)

const (
	network = "tcp"
)

type Redis struct {
	addr string
	log  bool
}

var (
	ErrServerStopped = errors.New("server is not started")
)

func New(addr string, log bool) *Redis {
	return &Redis{
		addr: addr,
		log:  log,
	}
}

func (r *Redis) Start() error {
	listen, err := net.Listen(network, r.addr)
	if err != nil {
		return fmt.Errorf("failed to listen to addrres %q: %v", r.addr, err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept a new conneciton: %v", err)
		}

		go r.handleRequest(conn)
	}

	// ToDo (@rodrigo) create gracefull shutdown with channel
}

func (r *Redis) handleRequest(conn net.Conn) {
	defer conn.Close()

	if r.log {
		fmt.Println("new client connected:", conn.RemoteAddr().String())
	}

	for {
		buffer := make([]byte, 1024)
		size, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}

		msg := string(buffer[:size])
		if r.log {
			fmt.Println("message received:", size, msg)
		}

		cmd, err := serde.Deserialize(msg)
		if err != nil {
			conn.Write(sendError("err", err.Error()))
			continue
		}

		r.handleCommand(cmd, conn)
	}
}

func (r *Redis) handleCommand(cmd serde.Command, conn net.Conn) {
	conn.Write([]byte{00})
}

func sendError(err string, msg string) []byte {
	err = strings.ToUpper(err)
	return []byte(fmt.Sprintf("-%s %s", err, msg))
}
