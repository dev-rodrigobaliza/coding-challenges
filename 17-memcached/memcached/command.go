package memcached

import (
	"encoding/json"
	"strconv"
	"strings"
)

type command struct {
	Name    string `json:"name"`
	Key     string `json:"key"`
	Flags   uint16 `json:"flags"`
	ExpTime uint32 `json:"exptime"`
	Size    uint32 `json:"size"`
	NoReply bool   `json:"noreply"`
	Data    []byte `json:"data,omitempty"`
}

func newCommand(msg string) *command {
	parts := strings.Split(msg, " ")
	if len(parts) < 5 || len(parts) > 6 {
		return nil
	}

	f, err := strconv.ParseUint(parts[2], 10, 16)
	if err != nil {
		return nil
	}

	e, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil
	}

	s, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return nil
	}

	if len(parts) == 6 && parts[6] != "noreply" {
		return nil
	}

	c := command{
		Name:    parts[0],
		Key:     parts[1],
		Flags:   uint16(f),
		ExpTime: uint32(e),
		Size:    uint32(s),
		NoReply: len(parts) == 6,
	}

	return &c
}

func (c *command) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}
