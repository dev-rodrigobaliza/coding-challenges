package ws

import (
	"errors"
	"fmt"
	"io"
	"net"
	"path"
	"strings"
	"ws/web"
)

const (
	network = "tcp"
)

type WebServer struct {
	addr string
	dir  string
	log  bool
}

func New(addr, dir string, log bool) *WebServer {
	return &WebServer{
		addr: addr,
		dir:  dir,
		log:  log,
	}
}

func (ws *WebServer) Start() error {
	listen, err := net.Listen(network, ws.addr)
	if err != nil {
		return fmt.Errorf("failed to listen to addrres %q: %v", ws.addr, err)
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return fmt.Errorf("failed to accept a new conneciton: %v", err)
		}

		go ws.handleRequest(conn)
	}

	// ToDo (@rodrigo) create gracefull shutdown with channel
}

func (ws *WebServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr().String()

	if ws.log {
		fmt.Println("new client connected:", addr)
	}

	buffer := make([]byte, 1024)
	size, err := conn.Read(buffer)
	if err != nil {
		if ws.log {
			if errors.Is(err, io.EOF) {
				fmt.Println("client disconnected", addr)
			} else {
				fmt.Printf("error reading: %#v\n", err)
			}
		}

		return
	}

	msg := string(buffer[:size])
	if ws.log {
		fmt.Println("message received:", size, msg)
	}

	ws.handleCommand(buffer[:size], conn)
}

func (ws *WebServer) handleCommand(rawReq []byte, conn net.Conn) {
	req, err := web.NewFromRaw(string(rawReq))
	if err != nil {
		web.SendResponse(400, "ERROR", "bad request", conn)
		return
	}

	switch strings.ToUpper(req.Method) {
	case "GET":
		ws.handleGet(req, conn)
		return

	default:
		web.SendResponse(400, "ERROR", "bad request", conn)
	}
}

func (ws *WebServer) handleGet(req *web.Request, conn net.Conn) {
	if req.Path == "/" {
		req.Path = "/index.html"
	}

	file := path.Join(ws.dir, req.Path[1:])
	msg, err := web.GetFile(file)
	if err != nil {
		web.SendResponse(404, "ERROR", "not found", conn)
		return
	}

	web.SendResponse(200, "OK", msg, conn)
}