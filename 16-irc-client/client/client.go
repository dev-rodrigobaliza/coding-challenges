package client

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"syscall"
)

var (
	ErrInvalidAddr      = errors.New("invalid irc server address")
	ErrInvalidPort      = errors.New("invalid irc server port")
	ErrAlreadyConnected = errors.New("already connected to a server")
	ErrConnectFailed    = errors.New("failed to connect to a server")

	Disconnected = "-=-disconnected-=-"
)

type MessageProcessor func(message string)

type Client struct {
	addr         string
	connected    bool
	conn         net.Conn
	msgProcessor MessageProcessor
}

func New(msgProcessor MessageProcessor) *Client {
	return &Client{msgProcessor: msgProcessor}
}

func (c *Client) Connect(addr, port string, ssl bool) error {
	if c.connected {
		return ErrAlreadyConnected
	}
	if addr == "" {
		return ErrInvalidAddr
	}
	if port == "" {
		return ErrInvalidPort
	}

	server := net.JoinHostPort(addr, port)
	conn, err := net.Dial("tcp", server)
	if err != nil {
		return errors.Join(err, ErrConnectFailed)
	}

	if ssl {
		s := tls.Client(conn, &tls.Config{
			ServerName: "irc.freenode.net",
		})
		if err := s.Handshake(); err != nil {
			return err
		}
		conn = s
	}

	c.addr = addr
	c.conn = conn

	go c.receiveMessages()

	return nil
}

func (c *Client) Disconnect() error {
	if c.connected {
		c.connected = false
		return c.conn.Close()
	}

	return nil
}

func (c *Client) SendMessage(msg string) (int, error) {
	return fmt.Fprintf(c.conn, "%s\r\n", msg)
}

func (c *Client) receiveMessages() {
	tp := textproto.NewReader(bufio.NewReader(c.conn))
	for {
		msg, err := tp.ReadLine()
		if err != nil {
			if isNetConnClosedErr(err) {
				c.connected = false
				c.msgProcessor(Disconnected)
			}
		}

		c.msgProcessor(msg)
	}
}

func isNetConnClosedErr(err error) bool {
	switch {
	case
		errors.Is(err, net.ErrClosed),
		errors.Is(err, io.EOF),
		errors.Is(err, syscall.EPIPE):
		return true
	default:
		return false
	}
}
