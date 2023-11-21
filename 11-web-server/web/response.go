package web

import (
	"fmt"
	"net"
)

func SendResponse(statusCode int, status, message string, conn net.Conn) {
	resp := fmt.Sprintf("HTTP/1.1 %d %s%s%s%s%s", statusCode, status, end, end, message, end)
	conn.Write([]byte(resp))
}
