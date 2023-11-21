package web

import (
	"fmt"
	"net"
	"os"
)

func SendResponse(statusCode int, status, message string, conn net.Conn) {
	resp := fmt.Sprintf("HTTP/1.1 %d %s%s%s%s%s", statusCode, status, end, end, message, end)
	conn.Write([]byte(resp))
}

func GetFile(file string) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
