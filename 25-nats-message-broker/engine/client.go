package engine

type Client struct {
	Addr      string
	Connected bool
	Command   *Command
}

func newClient(addr string) *Client {
	return &Client{
		Addr: addr,
	}
}
