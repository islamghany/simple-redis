package internal

import (
	"fmt"
	"strconv"
	"strings"
)

func fakeparser(data []byte) ([]string, error) {
	msg := string(data)
	if len(msg) == 0 {
		return nil, fmt.Errorf("invalid request")
	}
	cmds := strings.Split(msg, "\r\n")
	var args []string
	args = append(args, cmds[0])
	for i := int(1); i < len(cmds)-1; i += 2 {
		_, err := strconv.Atoi(cmds[i])
		if err != nil {
			return nil, fmt.Errorf("invalid request")
		}
		args = append(args, cmds[i+1])
	}
	return args, nil
}

func CMDParser(data []byte) ([]string, error) {
	// return fakeparser(data)
	msg := string(data)
	if len(msg) == 0 {
		return nil, fmt.Errorf("invalid request")
	}
	cmds := strings.Split(msg, "\r\n")
	var args []string
	for i := int(1); i < len(cmds)-1; i += 2 {
		length, err := strconv.Atoi(cmds[i][1:])
		if err != nil {
			return nil, fmt.Errorf("invalid request")
		}
		args = append(args, cmds[i+1][:length])
	}
	return args, nil
}

func Parser(raw []byte) ([]string, error) {
	cmds := FromBulkStringArray(string(raw))
	if len(cmds) < 1 {
		return nil, fmt.Errorf("invalid request")
	}
	cmd := strings.ToLower(string(FromBulkString(cmds[0])))
	var args []string
	args = append(args, cmd)
	for i := 1; i < len(cmds); i++ {
		args = append(args, string(FromBulkString(cmds[i])))
	}
	fmt.Println("args:", args)
	return args, nil
}
