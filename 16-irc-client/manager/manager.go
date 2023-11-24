package manager

import (
	"bufio"
	"fmt"
	"log/slog"
	"mirc/client"
	"mirc/logs"
	"os"
	"strings"
)

type Manager struct {
	cli     *client.Client
	logger  *slog.Logger
	nick    string
	channel string
}

func New(logAll bool) *Manager {
	logger := logs.New(logAll)

	return &Manager{
		logger: logger,
	}
}

func (m *Manager) Start(ssl bool, addr, port, user, nick, pass string) error {
	m.nick = nick
	m.cli = client.New(m.rxProcessor)
	if err := m.cli.Connect(addr, port, ssl); err != nil {
		m.logger.Error("failed to connect to the irc server", slog.Any("err", err))
		return err
	}

	m.logger.Info(
		"connected to the irc server",
		slog.String("address", addr),
		slog.String("port", port),
		slog.Bool("secure", ssl),
	)

	if pass != "" {
		m.txProcessor("PASS", pass)
	}
	m.txProcessor("NICK", nick)
	m.txProcessor("USER", user)

	// set up a goroutine to read commands from stdin
	in := make(chan string, 4)
	go func() {
		con := bufio.NewReader(os.Stdin)
		for {
			s, err := con.ReadString('\n')
			if err != nil {
				// wha?, maybe ctrl-D...
				close(in)
				break
			}
			// no point in sending empty lines down the channel
			if len(s) > 2 {
				in <- s[0 : len(s)-1]
			}
		}
	}()

	// set up a goroutine to do parsey things with the stuff from stdin
	go func() {
		for cmd := range in {
			cmd = cutNewLines(cmd)
			m.logger.Debug("command detected", slog.String("cmd", cmd))
			if cmd[0] == ':' {
				switch idx := strings.Index(cmd, " "); {
				case cmd[1] == 'p':
					m.part()

				case cmd[1] == 'l':
					m.listNames()

				case cmd[1] == 'q':
					var msg string
					if idx > -1 {
						msg = cmd[idx+1:]
					}
					m.Stop(msg)

				case idx == -1:
					m.logger.Warn("invalid command")
					continue

				case cmd[1] == 'n':
					m.nickChange(strings.TrimSpace(cmd[idx+1:]))

				case cmd[1] == 'j':
					m.join(strings.TrimSpace(cmd[idx+1:]))

				default:
					m.logger.Warn("unknown command")
				}
			} else {
				m.send(cmd)
			}
		}
	}()

	return nil
}

func (m *Manager) Stop(msg string) error {
	if msg == "" {
		msg = "Bye"
	}
	m.txProcessor("QUIT", msg)
	if err := m.cli.Disconnect(); err != nil {
		m.logger.Error("failed to disconnect from the irc server", slog.Any("err", err))
		return err
	}

	m.logger.Info("disconnected from the irc server")
	return nil
}

func (m *Manager) join(ch string) {
	if ch == "" {
		m.logger.Warn("you must inform the the name of the channel")
		return
	}

	m.part()

	if ch[0] != '#' {
		ch = "#" + ch
	}
	m.logger.Info("entering channel", slog.String("channel", ch))
	m.txProcessor("JOIN", ch)
	m.channel = ch
}

func (m *Manager) listNames() {
	if m.channel != "" {
		m.logger.Info("listing names from channel", slog.String("channel", m.channel))
		m.txProcessor("NAMES", m.channel)
	}
}

func (m *Manager) part() {
	if m.channel != "" {
		m.logger.Info("leaving channel", slog.String("channel", m.channel))
		m.txProcessor("PART", m.channel)
		m.channel = ""
	}
}

func (m *Manager) send(msg string) {
	if m.channel == "" {
		m.logger.Warn("you can not send messages outside of a channel")
		return
	}
	if msg == "" {
		m.logger.Warn("you must type any message to send it")
		return
	}

	m.txProcessor("PRIVMSG", m.channel, ":"+msg)
}

func (m *Manager) nickChange(nick string) {
	m.logger.Info("changing nickname", slog.String("old nickname", m.nick), slog.String("new nickname", nick))
	m.nick = nick
	m.txProcessor("NICK", nick)
}

func (m *Manager) rxProcessor(msg string) {
	m.logger.Debug("message received", slog.String("message", msg))

	switch {
	case strings.HasPrefix(msg, "PING"):
		m.txProcessor("PONG", msg[5:])

	case msg == client.Disconnected:
		m.logger.Warn("irc client disconnected from the server")
		os.Exit(1)

	default:
		msg = ""
	}
}

func (m *Manager) txProcessor(cmd string, params ...string) {
	param := strings.Join(params, " ")
	if param != "" {
		param = " " + param
	}
	msg := fmt.Sprintf("%s%s", cmd, param)

	size, err := m.cli.SendMessage(msg + "\r\n")
	if err != nil {
		m.logger.Error("failed to send message", slog.String("message", msg), slog.Any("err", err))
		return
	}

	m.logger.Debug("message sent", slog.String("message", msg), slog.Int("size", size))
}

// cutNewLines() pares down a string to the part before the first "\r" or "\n".
func cutNewLines(s string) string {
	r := strings.SplitN(s, "\r", 2)
	r = strings.SplitN(r[0], "\n", 2)
	return strings.TrimSpace(r[0])
}
