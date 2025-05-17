package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	store *Store
}

func NewServer(store *Store) *Server {
	return &Server{store: store}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		cmdLine := strings.TrimSpace(line)
		parts := strings.Fields(cmdLine)
		if len(parts) == 0 {
			fmt.Fprintf(conn, "-ERR empty command\r\n")
			continue
		}
		switch strings.ToUpper(parts[0]) {
		case "PING":
			fmt.Fprintf(conn, "PONG\r\n")
		case "SET":
			if len(parts) < 3 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'SET'\r\n")
				continue
			}
			key := parts[1]
			value := strings.Join(parts[2:], " ")
			s.store.Set(key, value)
			fmt.Fprintf(conn, "+OK\r\n")
		case "GET":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'GET'\r\n")
				continue
			}
			key := parts[1]
			val, ok := s.store.Get(key)
			if !ok {
				fmt.Fprintf(conn, "$-1\r\n")
			} else {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(val), val)
			}
		default:
			fmt.Fprintf(conn, "-ERR unknown command '%s'\r\n", parts[0])
		}
	}
}

func (s *Server) Listen(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Println("RedisGo server started on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go s.handleConnection(conn)
	}
}
