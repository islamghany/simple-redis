package commands

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func (exe *CommandExecutor) pxSet(conn net.Conn, key string, value string, ttlMS int) error {
	exe.mux.Lock()
	ttl := time.Duration(ttlMS) * time.Millisecond
	exe.store.KeyValueStore.SetWithTTL(key, value, ttl)
	exe.mux.Unlock()
	// if exe.cfg.Role == "slave" {
	// 	fmt.Print("Replicating to master", key, value)
	// 	return nil
	// }
	n, err := conn.Write([]byte("+OK\r\n"))
	if err != nil {
		return fmt.Errorf("Failed to write data: %v", err)
	}
	fmt.Println("Sent", n, "bytes")
	return nil
}

func (exe *CommandExecutor) infinitSet(conn net.Conn, key, value string) error {
	fmt.Println("Setting key", key, "value", value)
	exe.mux.Lock()
	exe.store.KeyValueStore.Set(key, value)
	exe.mux.Unlock()
	value, ok := exe.store.KeyValueStore.Get(key)
	fmt.Println("Set key", key, "value", value, "ok", ok)
	// if exe.cfg.Role == "slave" {
	// 	fmt.Print("Replicating to master", key, value)
	// 	return nil
	// }
	n, err := conn.Write([]byte("+OK\r\n"))
	if err != nil {
		return fmt.Errorf("Failed to write data: %v", err)
	}
	fmt.Println("Sent", n, "bytes")
	return nil
}

func (exe *CommandExecutor) set(conn net.Conn, args []string) error {
	if len(args) == 2 {
		return exe.infinitSet(conn, args[0], args[1])
	} else if len(args) == 4 {
		ttl, err := strconv.Atoi(args[3])
		if err != nil {
			return fmt.Errorf("ttl must be an integer")
		}
		return exe.pxSet(conn, args[0], args[1], ttl)
	} else {
		return fmt.Errorf("Invalid number of arguments")
	}

}

func (exe *CommandExecutor) get(conn net.Conn, key string) error {
	value, ok := exe.store.KeyValueStore.Get(key)
	fmt.Println("Getting value for key", key, "value", value, "ok", ok)

	if !ok {
		n, err := conn.Write([]byte("$-1\r\n"))
		if err != nil {
			return fmt.Errorf("Failed to write data: %v", err)
		}
		fmt.Println("Sent", n, "bytes")
		return nil
	}
	msg := fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)
	n, err := conn.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("Failed to write data: %v", err)
	}
	fmt.Println("Sent", n, "bytes")
	return nil
}
