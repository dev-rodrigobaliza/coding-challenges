package memcached

import (
	"errors"
	"io"
	"log/slog"
	"mc/logs"
	"mc/store"
	"net"
	"strings"
)

const (
	network = "tcp"
)

type Memcached struct {
	addr     string
	logger   *slog.Logger
	data     store.Store
	listener net.Listener
	clients  map[string]*client
}

func New(addr string, data store.Store, logAll bool) *Memcached {
	logger := logs.New(logAll)

	return &Memcached{
		addr:    addr,
		logger:  logger,
		data:    data,
		clients: make(map[string]*client),
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
		m.clients[addr] = newClient(conn)
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
	client, ok := m.clients[addr]
	if !ok {
		m.logger.Error("message received from a unknown client", slog.String("addr", addr), slog.Int("size", len(msg)))
		client.conn.Close()
		return
	}

	parts := strings.Split(msg, "\r\n")
	if len(parts) < 1 || len(parts) > 2 {
		m.logger.Warn("invalid message received", slog.String("client", client.String()), slog.Int("size", len(msg)))
	}

	cmd := newCommand(parts[0])
	if cmd == nil {
		m.logger.Warn("invalid command received", slog.String("client", client.String()), slog.String("message", parts[0]))
		return
	}

	if len(parts) == 2 {
		cmd.Data = []byte(parts[1])
	}

	m.logger.Debug("message received", slog.String("client", client.String()), slog.String("command", cmd.String()), slog.Int("data size", len(cmd.Data)))
	m.handleCommand(client, cmd)
}

func (m *Memcached) handleCommand(client *client, cmd *command) {
	// ToDo (@rodrigo) process the command
}
