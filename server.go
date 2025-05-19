package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
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
		case "DEL":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'DEL'\r\n")
				continue
			}
			count := 0
			for _, key := range parts[1:] {
				if s.store.Del(key) {
					count++
				}
			}
			fmt.Fprintf(conn, ":%d\r\n", count)
		case "EXPIRE":
			if len(parts) != 3 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'EXPIRE'\r\n")
				continue
			}
			key := parts[1]
			secs, err := strconv.Atoi(parts[2])
			if err != nil || secs < 0 {
				fmt.Fprintf(conn, "-ERR invalid expire time\r\n")
				continue
			}
			if s.store.Expire(key, secs) {
				fmt.Fprintf(conn, ":1\r\n")
			} else {
				fmt.Fprintf(conn, ":0\r\n")
			}
		case "KEYS":
			keys := s.store.Keys()
			fmt.Fprintf(conn, "*%d\r\n", len(keys))
			for _, key := range keys {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(key), key)
			}
		case "DUMPALL":
			all := s.store.DumpAll()
			fmt.Fprintf(conn, "*%d\r\n", len(all))
			for k, v := range all {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n$%d\r\n%s\r\n", len(k), k, len(v), v)
			}
		case "TTL":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'TTL'\r\n")
				continue
			}
			key := parts[1]
			ttl := s.store.TTL(key)
			fmt.Fprintf(conn, ":%d\r\n", ttl)
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
