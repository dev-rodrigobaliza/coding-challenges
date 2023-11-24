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
	cli    *client.Client
	logger *slog.Logger
}

func New(logAll bool) *Manager {
	logger := logs.New(logAll)

	return &Manager{
		logger: logger,
	}
}

func (m *Manager) Start(ssl bool, addr, port, user, nick, pass string) error {
	m.cli = client.New(m.rxProcessor)
	if err := m.cli.Connect(addr, port, ssl); err != nil {
		m.logger.Error("failed toconnect to the irc server", slog.Any("err", err))
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
			if cmd[0] == ':' {
				m.logger.Info("command detected", slog.String("cmd", cmd))
				// switch idx := strings.Index(cmd, " "); {
				// case cmd[1] == 'd':
				// 	fmt.Printf(c.String())
				// case cmd[1] == 'n':
				// 	parts := strings.Split(cmd, " ")
				// 	username := strings.TrimSpace(parts[1])
				// 	channelname := strings.TrimSpace(parts[2])
				// 	_, userIsOn := c.StateTracker().IsOn(channelname, username)
				// 	fmt.Printf("Checking if %s is in %s Online: %t\n", username, channelname, userIsOn)
				// case cmd[1] == 'f':
				// 	if len(cmd) > 2 && cmd[2] == 'e' {
				// 		// enable flooding
				// 		c.Config().Flood = true
				// 	} else if len(cmd) > 2 && cmd[2] == 'd' {
				// 		// disable flooding
				// 		c.Config().Flood = false
				// 	}
				// 	for i := 0; i < 20; i++ {
				// 		c.Privmsg("#", "flood test!")
				// 	}
				// case idx == -1:
				// 	continue
				// case cmd[1] == 'q':
				// 	c.Quit(cmd[idx+1 : len(cmd)])
				// case cmd[1] == 's':
				// 	c.Close()
				// case cmd[1] == 'j':
				// 	c.Join(cmd[idx+1 : len(cmd)])
				// case cmd[1] == 'p':
				// 	c.Part(cmd[idx+1 : len(cmd)])
				// }
			} else {
				msg := cutNewLines(cmd)
				m.txProcessor(msg)
			}
		}
	}()

	return nil
}

func (m *Manager) Stop() error {
	m.txProcessor("QUIT Bye")
	if err := m.cli.Disconnect(); err != nil {
		m.logger.Error("failed to disconnect from the irc server", slog.Any("err", err))
		return err
	}

	m.logger.Info("disconnected from the irc server")
	return nil
}

func (m *Manager) rxProcessor(msg string) {
	m.logger.Info("message received", slog.String("message", msg))

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

	m.logger.Info("message sent", slog.String("message", msg), slog.Int("size", size))
}

// cutNewLines() pares down a string to the part before the first "\r" or "\n".
func cutNewLines(s string) string {
	r := strings.SplitN(s, "\r", 2)
	r = strings.SplitN(r[0], "\n", 2)
	return r[0]
}
