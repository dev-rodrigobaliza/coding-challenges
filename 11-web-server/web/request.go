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
	Headers   []string
}

func NewFromRaw(rawReq string) (*Request, error) {
	parts := strings.Split(rawReq, end)

	inner := strings.Split(parts[0], " ")
	if len(inner) != 3 {
		return nil, ErrInvalidRawReq
	}

	req := Request{
		Method:  inner[0],
		Path:    inner[1],
		Version: inner[2],
	}

	for i, item := range parts {
		if i == 0 {
			continue
		}

		getHeader(&req, item)
	}

	return &req, nil
}

func getHeader(req *Request, item string) {
	if strings.HasPrefix(item, "Host:") {
		req.Host = item[6:]
	}
	if strings.HasPrefix(item, "User-Agent:") {
		req.UserAgent = item[12:]
	}
	if strings.HasPrefix(item, "Accept:") {
		req.Accept = item[8:]
	}

	req.Headers = append(req.Headers, item)
}
