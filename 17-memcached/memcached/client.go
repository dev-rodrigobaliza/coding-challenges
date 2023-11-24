package memcached

import "net"

type client struct {
	conn net.Conn
}

func newClient(conn net.Conn) *client {
	return &client{conn: conn}
}

func (c *client) String() string {
	return c.conn.RemoteAddr().String()
}
