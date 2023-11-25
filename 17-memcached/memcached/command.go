package memcached

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type command struct {
	Name      string `json:"name"`
	Key       string `json:"key"`
	Flags     uint16 `json:"flags"`
	ExpTime   uint32 `json:"exptime"`
	Size      uint32 `json:"size"`
	NoReply   bool   `json:"noreply"`
	Timestamp int64  `json:"timestamp"`
	Data      []byte `json:"data,omitempty"`
}

func newCommand(msgs []string) *command {
	parts := strings.Split(msgs[0], " ")
	if len(parts) < 2 {
		return nil
	}

	var (
		f   uint64
		e   uint64
		s   uint64
		t   int64
		d   []byte
		err error
	)

	name := strings.ToLower(parts[0])
	switch name {
	case "get":
		if len(parts) > 2 || len(msgs) > 1 {
			return nil
		}

	case "set":
		f, err = strconv.ParseUint(parts[2], 10, 16)
		if err != nil {
			return nil
		}

		e, err = strconv.ParseUint(parts[3], 10, 32)
		if err != nil {
			return nil
		}

		s, err = strconv.ParseUint(parts[4], 10, 32)
		if err != nil {
			return nil
		}

		if len(parts) == 6 && parts[6] != "noreply" {
			return nil
		}

		if e > 0 {
			t = time.Now().Add(time.Duration(e * uint64(time.Second))).UnixMicro()
		}

		if len(msgs) != 2 {
			return nil
		}

		d = []byte(msgs[1])
		if s == 0 || len(d) == 0 || int(s) != len(d) {
			return nil
		}

	default:
		return nil
	}

	c := command{
		Name:      parts[0],
		Key:       parts[1],
		Flags:     uint16(f),
		ExpTime:   uint32(e),
		Size:      uint32(s),
		Timestamp: t,
		NoReply:   len(parts) == 6,
		Data:      d,
	}

	return &c
}

func (c *command) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}
