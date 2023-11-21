package redis

import (
	"errors"
	"fmt"
	"io"
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

	addr := conn.RemoteAddr().String()

	if r.log {
		fmt.Println("new client connected:", addr)
	}

	for {
		buffer := make([]byte, 1024)
		size, err := conn.Read(buffer)
		if err != nil {
			if r.log {
				if errors.Is(err, io.EOF) {
					fmt.Println("client disconnected", addr)
				} else {
					fmt.Printf("error reading: %#v\n", err)
				}
			}

			return
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
	if cmd.Type == serde.Array && len(cmd.Array) > 0 {
		ca := cmd.Array[0]
		if ca.Type == serde.BulkString {
			switch ca.Value {
			case "ping":
				r.handlePing(conn)
				return

			case "echo":
				r.handleEcho(cmd, conn)
				return
			}
		}
	}

	conn.Write(sendError("err", "unknown command"))
}

func (r *Redis) handlePing(conn net.Conn) {
	c := serde.Command{
		Type:  serde.SimpleString,
		Value: "PONG",
	}

	sendResponse(c, conn)
}

func (r *Redis) handleEcho(cmd serde.Command, conn net.Conn) {
	if len(cmd.Array) == 2 {
		ca := cmd.Array[1]
		c := serde.Command{
			Type:  serde.SimpleString,
			Value: ca.Value,
		}

		sendResponse(c, conn)
	}

	conn.Write(sendError("err", "bad echo request"))
}

func sendError(err string, msg string) []byte {
	err = strings.ToUpper(err)
	return []byte(fmt.Sprintf("-%s %s", err, msg))
}

func sendResponse(cmd serde.Command, conn net.Conn) {
	msg, err := serde.Serialize(cmd)
	if err != nil {
		conn.Write(sendError("err", "internal server error"))
		return
	}

	conn.Write([]byte(msg))
}
