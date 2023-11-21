package web

import (
	"errors"
	"strings"
)

var (
	ErrInvalidRawReq = errors.New("invalid raw request")
)

type Request struct {
	Method    string
	Path      string
	Version   string
	Host      string
	UserAgent string
	Accept    string
	Data      string
}

func NewFromRaw(rawReq string) (*Request, error) {
	parts := strings.Split(rawReq, end)
	if len(parts) != 6 {
		return nil, ErrInvalidRawReq
	}

	inner := strings.Split(parts[0], " ")
	if len(inner) != 3 {
		return nil, ErrInvalidRawReq
	}

	req := Request{
		Method:    inner[0],
		Path:      inner[1],
		Version:   inner[2],
		Host:      parts[1],
		UserAgent: parts[2],
		Accept:    parts[3],
		Data:      parts[4],
	}

	return &req, nil
}
