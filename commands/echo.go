package commands

import (
	"fmt"
	"net"
)

func (exe *CommandExecutor) echo(conn net.Conn, message string) error {
	msg := fmt.Sprintf("$%d\r\n%s\r\n", len(message), message)
	n, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("Failed to write data", err)
		return err
	}
	fmt.Println("Sent", n, "bytes")
	return nil
}
