package commands

import (
	"fmt"
	"net"
)

func (exe *CommandExecutor) StartReplication(conn net.Conn) {
	ping := "*1\r\n$4\r\nping\r\n"
	replconfListeningPort := fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%s\r\n", exe.cfg.Port)
	replconfCapa := "*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"
	psync := "*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"
	conn.Write([]byte(ping))
	conn.Write([]byte(replconfListeningPort))
	conn.Write([]byte(replconfCapa))
	conn.Write([]byte(psync))
}
