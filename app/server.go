package main

import (
	"flag"
	"fmt"
	"io"
	"sync"

	// Uncomment this block to pass the first stage
	"net"
	"os"

	"github.com/islamghany/simple-redis/commands"
	"github.com/islamghany/simple-redis/internal"
	"github.com/islamghany/simple-redis/store"
	"github.com/islamghany/simple-redis/utils"
)

func main() {

	var cfg internal.Config
	replicas := make(map[string]net.Conn)
	var isSlave bool
	flag.StringVar(&cfg.Port, "port", "6379", "Port to bind to")
	flag.BoolVar(&isSlave, "replicaof", false, "Start as a replica")
	flag.Parse()
	args := flag.Args()
	if isSlave {
		cfg.Role = "slave"
		cfg.MasterHost = args[0]
		cfg.MasterPort = args[1]
	} else {
		cfg.Role = "master"
		cfg.Master_ID = utils.RandomID(40)
		cfg.MasterReplOffset = 0
	}

	//Uncomment this block to pass the first stage
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", cfg.Port))
	if err != nil {
		fmt.Printf("Failed to bind to port %s", cfg.Port)
		os.Exit(1)
	}
	defer l.Close()
	store := store.NewStore()
	executor := commands.NewCommandExecutor(store, &sync.Mutex{}, &cfg, replicas)
	if isSlave {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", cfg.MasterHost, cfg.MasterPort))
		if err != nil {
			fmt.Println("Failed to connect to master", err)
			os.Exit(1)
		}
		defer conn.Close()
		go executor.StartReplication(conn)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
			continue
		}
		go handleClient(conn, executor)
	}

}

func handleClient(conn net.Conn, exe *commands.CommandExecutor) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Failed to read data", err)
			return
		}
		data := buffer[:n]
		fmt.Println("Received data:", string(data))
		args, err := internal.Parser(data)
		fmt.Println("Received data:", args, "err:", err)
		if err != nil {
			fmt.Println("Failed to parse data", err)
			return
		}
		err = exe.Execute(conn, args, data)
		if err != nil {
			fmt.Println("Failed to execute command", err)
			return
		}
	}
}
