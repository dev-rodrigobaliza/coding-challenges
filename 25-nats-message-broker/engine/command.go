package engine

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type Command struct {
	Name   string `json:"name"`
	Topic  string `json:"topic,omitempty"`
	Number int    `json:"number,omitempty"`
	Data   []byte `json:"data,omitempty"`
}

func newCommand(req []byte) *Command {
	parts := bytes.Split(req, []byte(newLine))
	parts = parts[:len(parts)-1]
	if len(parts) < 1 || len(parts) > 2 {
		return nil
	}

	var (
		name   string
		topic  string
		number int
		data   []byte
	)

	headers := strings.Split(string(parts[0]), " ")
	if len(headers) < 1 {
		return nil
	}

	name = strings.ToLower(headers[0])
	switch name {
	case "connect":
		break

	case "ping":
		break

	case "pong":
		break

	case "sub":
		topic, number := parseHeader(headers)
		if topic == "" || number == 0 {
			return nil
		}

	case "pub":
		topic, number := parseHeader(headers)
		if topic == "" || number == 0 {
			return nil
		}
		if len(parts) < 2 {
			return nil
		}

		data = parts[1]

	default:
		return nil
	}

	cmd := Command{
		Name:   name,
		Topic:  topic,
		Number: number,
		Data:   data,
	}

	return &cmd
}

func (c *Command) String() string {
	s := fmt.Sprintf("[name: %s]", c.Name)
	if c.Topic != "" {
		s = fmt.Sprintf("%s[topic: %s]", s, c.Topic)
	}
	if c.Number > 0 {
		s = fmt.Sprintf("%s[number: %d]", s, c.Number)
	}
	if len(c.Data) > 0 {
		s = fmt.Sprintf("%s[data: %d]", s, len(c.Data))
	}

	return s
}

func parseHeader(headers []string) (string, int) {
	if len(headers) != 3 {
		return "", 0
	}

	t := headers[1]
	n, err := strconv.Atoi(headers[2])
	if err != nil {
		return "", 0
	}

	return t, n
}
