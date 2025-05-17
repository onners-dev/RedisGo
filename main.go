package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		cmd := strings.TrimSpace(line)
		switch strings.ToUpper(cmd) {
		case "PING":
			fmt.Fprintf(conn, "PONG\r\n")
		default:
			fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", cmd)
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("RedisGo server started on :6379")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
