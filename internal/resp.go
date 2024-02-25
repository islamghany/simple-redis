package internal

import (
	"strconv"
	"strings"
)

const (
	Terminator = "\r\n"
)
const (
	SimpleString = '+'
	Error        = '-'
	Integer      = ':'
	BulkString   = '$'
	Array        = '*'
)

func ToSimpleString(data []byte) string {
	return string(SimpleString) + string(data) + Terminator
}
func FromSimpleString(raw string) []byte {
	if raw[0] != SimpleString {
		return nil
	}
	return []byte(raw[1 : len(raw)-2])
}

// ToSimpleString converts a string to a RESP bulk string
func ToBulkString(data []byte) string {
	l := len(data)
	return "$" + strconv.Itoa(l) + Terminator + string(data) + Terminator
}

// FromBulkString extracts the data from a RESP bulk string
func FromBulkString(raw string) []byte {
	if len(raw) == 0 {
		return nil
	}
	if raw[0] != BulkString {
		return nil
	}
	parts := strings.Split(raw, Terminator)
	if len(parts) != 3 {
		return nil
	}
	return []byte(parts[1])
}

// FromBulkStringArray extracts elements from a RESP array.
// Each element is a bulk string.
func FromBulkStringArray(raw string) []string {
	if len(raw) == 0 {
		return nil
	}
	if raw[0] != Array {
		return nil
	}
	num, i := 0, 1
	// get the number of bulk strings in the array
	for raw[i] != '\r' {
		num = num*10 + int(raw[i]-'0')
		i++
	}
	datas := raw[i+2:]
	// cut the string into bulk strings
	var result []string
	p := 0
	for num > 0 {
		if datas[p] == BulkString {
			k := p + 1
			n := 0
			for datas[k] != '\r' {
				n = n*10 + int(datas[k]-'0')
				k++
			}
			nextP := k + 2 + n + 2
			result = append(result, datas[p:nextP])
			p = nextP
			num--
		}
	}
	return result
}
func ToBulkStringArray(data []string) string {
	l := len(data)
	res := "*" + strconv.Itoa(l) + Terminator
	for _, v := range data {
		res += ToBulkString([]byte(v))
	}
	return res
}
