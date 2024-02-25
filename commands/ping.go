package commands

import (
	"fmt"
	"net"
)

func (exe *CommandExecutor) pong(conn net.Conn) error {
	n, err := conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		return fmt.Errorf("Failed to write data: %v", err)
	}
	fmt.Println("Sent", n, "bytes")
	return nil
}
