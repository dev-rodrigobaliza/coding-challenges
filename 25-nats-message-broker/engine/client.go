package engine

import "net"

type Client struct {
	Conn      net.Conn
	Connected bool
	Command   *Command
}

func newClient(conn net.Conn) *Client {
	return &Client{
		Conn: conn,
	}
}

func (c *Client) Addr() string {
	return c.Conn.RemoteAddr().String()
}