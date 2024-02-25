package commands

import (
	"fmt"
	"net"
)

func (exe *CommandExecutor) info(conn net.Conn, args []string) {
	fmt.Println("Info command")
	res := fmt.Sprintf("role:%s\r\nmaster_replid:%s\r\nmaster_repl_offset:%d\r\n", exe.cfg.Role, exe.cfg.Master_ID, exe.cfg.MasterReplOffset)
	bulkString := fmt.Sprintf("$%d\r\n%s\r\n", len(res), res)
	conn.Write([]byte(bulkString))
	fmt.Println("Response:", res)
}
