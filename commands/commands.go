package commands

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync"

	"github.com/islamghany/simple-redis/internal"
	"github.com/islamghany/simple-redis/store"
)

type CommandExecutor struct {
	store    *store.Store
	mux      *sync.Mutex
	cfg      *internal.Config
	replicas map[string]net.Conn
}

func NewCommandExecutor(store *store.Store, mux *sync.Mutex, cfg *internal.Config, replicas map[string]net.Conn) *CommandExecutor {
	return &CommandExecutor{
		store:    store,
		mux:      mux,
		cfg:      cfg,
		replicas: replicas,
	}
}

func (exe *CommandExecutor) Execute(conn net.Conn, args []string, rawMessage []byte) error {
	if len(args) == 0 {
		return fmt.Errorf("invalid command")
	}
	cmd := args[0]
	switch cmd {
	case "echo", "ECHO":
		return exe.echo(conn, args[1])
	case "ping", "PING":
		return exe.pong(conn)
	case "set", "SET":
		err := exe.set(conn, args[1:])
		if err != nil {
			return nil
		}
		for _, replica := range exe.replicas {
			_, err := replica.Write([]byte(rawMessage))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		return nil
	case "get", "GET":
		return exe.get(conn, args[1])
	case "info", "INFO":
		exe.info(conn, args[1:])
		return nil
	case "REPLCONF", "replconf":
		exe.replicas[conn.RemoteAddr().String()] = conn
		_, err := conn.Write([]byte("+OK\r\n"))
		return err
	case "PSYNC", "psync":
		_, err := conn.Write([]byte(fmt.Sprintf("+FULLRESYNC %s 0\r\n", exe.cfg.Master_ID)))
		RDBContent, _ := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
		conn.Write([]byte(fmt.Sprintf("$%v\r\n%v", len(string(RDBContent)), string(RDBContent))))
		return err
	default:
		return fmt.Errorf("unknown command")
	}
}
