package memcached

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mc/logs"
	"mc/store"
	"net"
	"strings"
	"time"
)

const (
	network    = "tcp"
	endMessage = "\r\n"
	notStored  = "NOT_STORED"
	stored     = "STORED"
	end        = "END"
	value      = "VALUE"
)

type Memcached struct {
	addr     string
	logger   *slog.Logger
	data     store.Store
	listener net.Listener
}

func New(addr string, data store.Store, logAll bool) *Memcached {
	logger := logs.New(logAll)

	return &Memcached{
		addr:   addr,
		logger: logger,
		data:   data,
	}
}

func (m *Memcached) Start() error {
	m.logger.Info("starting the server", slog.String("addr", m.addr))

	listen, err := net.Listen(network, m.addr)
	if err != nil {
		m.logger.Error("failed to start the server", slog.String("addr", m.addr), slog.Any("err", err))
		return err
	}

	m.listener = listen

	go m.handleConnections()

	return nil
}

func (m *Memcached) Stop() error {
	m.logger.Info("stoping the server", slog.String("addr", m.addr))

	return m.listener.Close()
}

func (m *Memcached) handleConnections() {
	m.logger.Info("the server is accepting new connections", slog.String("addr", m.addr))

	for {
		conn, err := m.listener.Accept()
		if err != nil {
			m.logger.Error("failed to accept a new connection", slog.Any("err", err))
			continue
		}

		addr := conn.RemoteAddr().String()
		m.logger.Debug("new client connected", slog.String("addr", addr))

		go m.handleRequest(conn)
	}
}

func (m *Memcached) handleRequest(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		size, err := conn.Read(buffer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				m.logger.Error("failed to read bytes from the connection", slog.Any("err", err))
			}

			break
		}

		if size > 0 {
			msg := string(buffer[:size])
			go m.handleMessage(conn, msg)
		}
	}

	conn.Close()
	m.logger.Error("client disconnected")
}

func (m *Memcached) handleMessage(conn net.Conn, msg string) {
	addr := conn.RemoteAddr().String()

	parts := strings.Split(msg, "\r\n")
	if len(parts) > 1 {
		parts = parts[:len(parts)-1] // remove the last one (always empty)
	}

	if len(parts) == 0 {
		m.logger.Warn("invalid message received", slog.String("client", addr))
		m.sendError(conn, end)
		return
	}

	cmd := newCommand(parts)
	if cmd == nil {
		m.logger.Warn("invalid command received", slog.String("client", addr), slog.String("message", parts[0]))
		m.sendError(conn, end)
		return
	}

	m.logger.Debug("message received", slog.String("client", addr), slog.String("command", parts[0]), slog.Int("data size", len(cmd.Data)))
	m.handleCommand(conn, cmd)
}

func (m *Memcached) handleCommand(conn net.Conn, cmd *command) {
	switch cmd.Name {
	case "set":
		m.cmdSet(conn, cmd)

	case "get":
		m.cmdGet(conn, cmd)

	default:
		m.logger.Warn("unknown command received", slog.String("command", cmd.String()))
		if !cmd.NoReply {
			m.sendError(conn, end)
		}
	}
}

func (m *Memcached) cmdSet(conn net.Conn, cmd *command) {
	// if m.data.Get(cmd.Key) != nil {
	// 	m.logger.Warn("failed to execute command set", slog.String("status", "key already exists"), slog.String("key", cmd.Key))
	// 	if !cmd.NoReply {
	// 		m.sendError(conn, notStored)
	// 	}

	// 	return
	// }

	val, err := json.Marshal(cmd)
	if err != nil {
		m.logger.Error("failed to execute command set", slog.String("status", "failed to marshal data"), slog.Any("err", err))
		if !cmd.NoReply {
			m.sendError(conn, notStored)
		}

		return
	}

	m.data.Set(cmd.Key, val)
	if !cmd.NoReply {
		m.sendMessage(conn, stored)
	}
}

func (m *Memcached) cmdGet(conn net.Conn, cmd *command) {
	val := m.data.Get(cmd.Key)
	if val == nil {
		if !cmd.NoReply {
			m.sendMessage(conn, end)
		}

		return
	}

	var c command
	if err := json.Unmarshal(val, &c); err != nil {
		m.logger.Error("failed to execute command get", slog.String("status", "failed to unmarshal data"), slog.Any("err", err))
		if !cmd.NoReply {
			m.sendError(conn, notStored)
		}

		return
	}

	if c.Timestamp > 0 && c.Timestamp < time.Now().UnixMicro() {
		if !cmd.NoReply {
			m.sendMessage(conn, end)
		}

		return
	}

	m.sendCommand(conn, &c, value)
}

func (m *Memcached) sendError(conn net.Conn, msg string) {
	m.sendMessage(conn, msg)
}

func (m *Memcached) sendCommand(conn net.Conn, cmd *command, msg string) {
	msg = fmt.Sprintf("%s %s %d %d", msg, cmd.Key, cmd.Flags, cmd.Size)
	m.sendMessage(conn, msg, string(cmd.Data), end)
}

func (m *Memcached) sendMessage(conn net.Conn, msgs ...string) {
	total := 0
	for _, msg := range msgs {
		msg = msg + endMessage
		n, err := conn.Write([]byte(msg))
		if err != nil {
			m.logger.Error("failed to send message", slog.Any("err", err))
			return
		}

		total += n
	}

	addr := conn.RemoteAddr().String()
	m.logger.Debug("message sent", slog.String("client", addr), slog.Int("message size", total))
}
