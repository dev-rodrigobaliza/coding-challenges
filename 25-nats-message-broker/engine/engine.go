package engine

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"nats/logs"
	"nats/safemap"
	"nats/store"
	"net"
)

const (
	network = "tcp"
	newLine = "\r\n"

	ok             = "+OK"
	pong           = "PONG"
	info           = "INFO"
	message        = "MSG"
	errPrefix      = "-ERR"
	unknownCommand = "Unknown Protocol Operation"
	unexpectedData = "Unexpected Data"
)

var (
	ErrTopicNotFound = errors.New("topic not found")
)

type Engine struct {
	addr     string
	logger   *slog.Logger
	topics   store.Store[map[string]struct{}]
	clients  store.Store[*Client]
	listener net.Listener
}

func New(addr, dsn string, logAll bool) *Engine {
	logger := logs.New(logAll)
	topics, err := getStore[map[string]struct{}](dsn)
	if err != nil {
		logger.Error("failed to create the topics store", slog.String("dsn", dsn))
		panic(err)
	}
	clients, err := getStore[*Client](dsn)
	if err != nil {
		logger.Error("failed to create the clients store", slog.String("dsn", dsn))
		panic(err)
	}

	return &Engine{
		addr:    addr,
		logger:  logger,
		topics:  topics,
		clients: clients,
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

		c := newClient(conn)
		e.clients.Set(c.Addr(), c)
		e.sendInfo(conn)
		go e.handleRequests(conn)

		e.logger.Debug("new client connected", slog.String("addr", c.Addr()), slog.Int("clients connected", e.clients.Len()))
	}
}

func (e *Engine) handleRequests(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	client := e.clients.Get(addr)
	if client.Addr() == "" {
		e.logger.Warn("received request from an unknown client", slog.String("addr", addr))
		conn.Close()
		return
	}

	reader := bufio.NewReader(conn)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		go e.handleMessage(conn, scanner.Bytes())
	}
	if err := scanner.Err(); err != nil {
		if !errors.Is(err, io.EOF) {
			e.logger.Error("failed to read bytes from the connection", slog.Any("err", err))
		}
	}

	if err := conn.Close(); err != nil {
		e.logger.Error("failed to close the client connection", slog.String("addr", addr), slog.Any("err", err))
	}

	e.clients.Delete(addr)
	e.logger.Error("client disconnected", slog.String("addr", addr), slog.Int("clients connected", e.clients.Len()))
}

func (e *Engine) handleMessage(conn net.Conn, message []byte) {
	addr := conn.RemoteAddr().String()
	client := e.clients.Get(addr)
	if client == nil || client.Addr() == "" {
		e.logger.Warn("received message from an unknown client", slog.String("addr", addr))
		return
	}

	if client.Command != nil {
		e.logger.Debug("data received", slog.String("client", addr), slog.Int("data size", len(message)))
		e.handleData(conn, message)
		return
	}

	cmd := newCommand(message)
	e.handleCommand(conn, cmd)
}

func (e *Engine) handleCommand(conn net.Conn, cmd *Command) {
	addr := conn.RemoteAddr().String()
	client := e.clients.Get(addr)
	if client.Addr() == "" {
		e.logger.Warn("received command from an unknown client", slog.String("addr", addr))
		return
	}

	if !client.Connected && (cmd == nil || cmd.Name != "connect") {
		conn.Close()
		return
	}

	switch cmd.Name {
	case "connect":
		client.Connected = true
		e.sendCommand(conn, ok)

	case "ping":
		e.sendCommand(conn, pong)

	case "pub":
		client.Command = cmd

	case "sub":
		clients := e.topics.Get(cmd.Topic)
		if clients == nil {
			clients = make(map[string]struct{})
			e.topics.Set(cmd.Topic, clients)
		}

		clients[addr] = struct{}{}
		e.topics.Set(cmd.Topic, clients)
		e.sendCommand(conn, ok)
		e.logger.Debug("new topic subscription", slog.String("topic", cmd.Topic), slog.String("client", addr), slog.Int("subscribers", len(clients)))

	case "unsub":
		clients := e.topics.Get(cmd.Topic)
		if clients != nil {
			delete(clients, addr)
			e.topics.Set(cmd.Topic, clients)
		}

		e.sendCommand(conn, ok)
		e.logger.Debug("topic unsubscription", slog.String("topic", cmd.Topic), slog.String("client", addr), slog.Int("subscribers", len(clients)))

	default:
		e.logger.Warn("unknown command received", slog.String("client", addr), slog.String("command", cmd.String()))
		if !client.Connected {
			conn.Close()
			return
		}

		e.sendError(conn, unknownCommand)
	}
}

func (e *Engine) handleData(conn net.Conn, message []byte) {
	addr := conn.RemoteAddr().String()
	client := e.clients.Get(addr)
	if client.Addr() == "" {
		e.logger.Warn("received command from an unknown client", slog.String("addr", addr))
		return
	}

	switch client.Command.Name {
	case "pub":
		e.sendCommand(conn, ok)
		e.sendTopic(client.Command.Topic, string(message))

	default:
		e.logger.Warn("unexpected data received", slog.String("client", addr), slog.Int("data size", len(message)))
		e.sendError(conn, unknownCommand)
	}

	client.Command = nil
}

func (e *Engine) sendError(conn net.Conn, msg string) {
	e.sendMessage(conn, fmt.Sprintf("%s '%s'", errPrefix, msg))
}

func (e *Engine) sendCommand(conn net.Conn, msg string) {
	e.sendMessage(conn, msg)
}

func (e *Engine) sendTopic(topic, msg string) {
	clients := e.topics.Get(topic)

	for addr := range clients {
		client := e.clients.Get(addr)
		if client == nil || client.Conn == nil {
			e.clients.Delete(addr)
			delete(clients, addr)
			e.topics.Set(topic, clients)
			continue
		}

		e.sendCommand(client.Conn, fmt.Sprintf("%s %s", message, topic))
		e.sendCommand(client.Conn, msg)
	}
}

func (e *Engine) sendBroadcast(msg string) {
	addrs := e.clients.GetAllKeys()
	for _, addr := range addrs {
		client := e.clients.Get(addr)
		if client == nil || client.Conn == nil {
			e.clients.Delete(addr)
			continue
		}

		e.sendCommand(client.Conn, msg)
	}
}

func (e *Engine) sendInfo(conn net.Conn) {
	msg := fmt.Sprintf(`INFO {"server_addr": "%s", "client_addr": "%s"}`, e.addr, conn.RemoteAddr().String())
	e.sendMessage(conn, msg, newLine)
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

func getStore[D any](dsn string) (store.Store[D], error) {
	if dsn == ":memory:" {
		data := safemap.New[D]()
		return data, nil
	}
	return nil, errors.New("unknown store")
}
