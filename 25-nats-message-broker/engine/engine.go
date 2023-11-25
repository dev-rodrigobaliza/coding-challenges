package engine

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"nats/logs"
	"nats/store"
	"net"
)

const (
	network = "tcp"
	newLine = "\r\n"
)

type Engine struct {
	addr     string
	logger   *slog.Logger
	data     store.Store
	listener net.Listener
}

func New(addr string, data store.Store, logAll bool) *Engine {
	logger := logs.New(logAll)

	return &Engine{
		addr:   addr,
		logger: logger,
		data:   data,
	}
}

func (e *Engine) Start() error {
	e.logger.Info("starting the server", slog.String("addr", e.addr))

	listen, err := net.Listen(network, e.addr)
	if err != nil {
		e.logger.Error("failed to start the server", slog.String("addr", e.addr), slog.Any("err", err))
		return err
	}

	e.listener = listen

	go e.handleConnections()

	return nil
}

func (e *Engine) Stop() error {
	e.logger.Info("stoping the server", slog.String("addr", e.addr))

	return e.listener.Close()
}

func (e *Engine) handleConnections() {
	e.logger.Info("the server is accepting new connections", slog.String("addr", e.addr))

	for {
		conn, err := e.listener.Accept()
		if err != nil {
			e.logger.Error("failed to accept a new connection", slog.Any("err", err))
			continue
		}

		addr := conn.RemoteAddr().String()
		e.logger.Debug("new client connected", slog.String("addr", addr))

		go e.handleRequest(conn)
	}
}

func (e *Engine) handleRequest(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		size, err := conn.Read(buffer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				e.logger.Error("failed to read bytes from the connection", slog.Any("err", err))
			}

			break
		}

		if size > 0 {
			go e.handleMessage(conn, buffer[:size])
		}
	}

	conn.Close()
	e.logger.Error("client disconnected")
}

func (e *Engine) handleMessage(conn net.Conn, message []byte) {
	addr := conn.RemoteAddr().String()
	cmd := newCommand(message)
	if cmd == nil {
		e.logger.Warn("invalid command received", slog.String("client", addr), slog.String("message", fmt.Sprintf("%x", message)))
		e.sendError(conn)
		return
	}

	e.logger.Debug("message received", slog.String("client", addr), slog.String("command", cmd.String()))
	e.handleCommand(conn, cmd)
}

func (e *Engine) handleCommand(conn net.Conn, cmd *Command) {
	switch cmd.Name {
	case "pub":
		e.logger.Debug("pub cmd")

	case "sub":
		e.logger.Debug("sub cmd")

	default:
		e.logger.Warn("unknown command received", slog.String("command", cmd.String()))
		e.sendError(conn)
	}
}

func (e *Engine) sendError(conn net.Conn) {
	e.sendMessage(conn, "NOK")
}

func (e *Engine) sendCommand(conn net.Conn, cmd *Command, msg string) {
	// msg = fmt.Sprintf("%s %s %d %d", msg, cmd.Key, cmd.Flags, cmd.Size)
	// e.sendMessage(conn, msg, string(cmd.Data), end)
}

func (e *Engine) sendMessage(conn net.Conn, msgs ...string) {
	total := 0
	for _, msg := range msgs {
		msg = msg + newLine
		n, err := conn.Write([]byte(msg))
		if err != nil {
			e.logger.Error("failed to send message", slog.Any("err", err))
			return
		}

		total += n
	}

	addr := conn.RemoteAddr().String()
	e.logger.Debug("message sent", slog.String("client", addr), slog.Int("message size", total))
}
