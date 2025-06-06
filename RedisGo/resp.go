package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parse a RESP (REdis Serialization Protocol) array and return a slice of strings.
func parseRESP(reader *bufio.Reader) ([]string, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	if b != '*' {
		return nil, errors.New("expected '*' for RESP array")
	}
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(line))
	if err != nil {
		return nil, err
	}
	parts := make([]string, count)
	for i := 0; i < count; i++ {
		prefix, err := reader.ReadByte()
		if err != nil {
			return nil, err
		}
		if prefix != '$' {
			return nil, errors.New("expected '$' for RESP bulk string")
		}
		slen, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		strlen, err := strconv.Atoi(strings.TrimSpace(slen))
		if err != nil {
			return nil, err
		}
		buf := make([]byte, strlen+2) // +2 for \r\n
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return nil, err
		}
		parts[i] = string(buf[:strlen])
	}
	return parts, nil
}

// RESP response helpers
func respSimple(msg string) string { return "+" + msg + "\r\n" }
func respError(msg string) string  { return "-ERR " + msg + "\r\n" }
func respInt(n int) string         { return ":" + strconv.Itoa(n) + "\r\n" }
func respBulk(msg string) string   { return fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg) }
func respNullBulk() string         { return "$-1\r\n" }
func respArray(arr []string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("*%d\r\n", len(arr)))
	for _, item := range arr {
		buf.WriteString(respBulk(item))
	}
	return buf.String()
}
