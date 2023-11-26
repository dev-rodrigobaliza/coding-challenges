package engine

import (
	"fmt"
	"strconv"
	"strings"
)

type Command struct {
	Name   string `json:"name"`
	Topic  string `json:"topic,omitempty"`
	Number int    `json:"number,omitempty"`
}

func newCommand(req []byte) *Command {
	var (
		name   string
		topic  string
		number int
	)

	headers := strings.Split(string(req), " ")
	if len(headers) < 1 {
		return nil
	}

	name = strings.ToLower(headers[0])
	switch name {
	case "connect":
		if headers[1][0] != '{' || headers[1][len(headers[1])-1] != '}' {
			return nil
		}

	case "ping":
		break

	case "sub", "unsub":
		topic, number = parseHeader(headers)
		if topic == "" {
			return nil
		}

	case "pub":
		topic, number = parseHeader(headers)
		if topic == "" || number == 0 {
			return nil
		}

	default:
		return nil
	}

	cmd := Command{
		Name:   name,
		Topic:  topic,
		Number: number,
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

	return s
}

func parseHeader(headers []string) (t string, n int) {
	if len(headers) < 2 {
		return
	}

	if len(headers) > 1 {
		t = headers[1]
	}

	if len(headers) > 2 {
		var err error
		n, err = strconv.Atoi(headers[2])
		if err != nil {
			return "", 0
		}
	}

	return
}
