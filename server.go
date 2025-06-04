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
		case "INCR":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'INCR'\r\n")
				continue
			}
			key := parts[1]
			val, err := s.store.Incr(key)
			if err != nil {
				fmt.Fprintf(conn, "-ERR %s\r\n", err.Error())
			} else {
				fmt.Fprintf(conn, ":%d\r\n", val)
			}
		case "DECR":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'DECR'\r\n")
				continue
			} else {
				key := parts[1]
				val, err := s.store.Decr(key)
				if err != nil {
					fmt.Fprintf(conn, "-ERR %s\r\n", err.Error())
				} else {
					fmt.Fprintf(conn, ":%d\r\n", val)
				}
			}
		case "MSET":
			if len(parts) < 3 || len(parts[1:])%2 != 0 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'MSET'\r\n")
				continue
			}
			err := s.store.MSet(parts[1:]...)
			if err != nil {
				fmt.Fprintf(conn, "-ERR %s\r\n", err.Error())
			} else {
				fmt.Fprintf(conn, "+OK\r\n")
			}
		case "MGET":
			if len(parts) < 2 {
				fmt.Fprintf(conn, "-ERR wrong number of arguments for 'MGET'\r\n")
				continue
			}
			values := s.store.MGet(parts[1:]...)
			fmt.Fprintf(conn, "*%d\r\n", len(values))
			for _, v := range values {
				if v == "" {
					fmt.Fprintf(conn, "$-1\r\n")
				} else {
					fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(v), v)
				}
			}
		case "LPUSH":
			if len(parts) < 3 {
				fmt.Fprintf(conn, "-ERR usage: LPUSH key value [value ...]\r\n")
				continue
			}
			key := parts[1]
			values := parts[2:]
			n := s.store.LPush(key, values...)
			fmt.Fprintf(conn, ":%d\r\n", n)

		case "RPOP":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR usage: RPOP key\r\n")
				continue
			}
			key := parts[1]
			val, err := s.store.RPop(key)
			if err != nil {
				fmt.Fprintf(conn, "$-1\r\n")
			} else {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(val), val)
			}

		case "LLEN":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR usage: LLEN key\r\n")
				continue
			}
			key := parts[1]
			n := s.store.LLen(key)
			fmt.Fprintf(conn, ":%d\r\n", n)

		case "SADD":
			if len(parts) < 3 {
				fmt.Fprintf(conn, "-ERR usage: SADD key member [member ...]\r\n")
				continue
			}
			key := parts[1]
			members := parts[2:]
			n := s.store.SAdd(key, members...)
			fmt.Fprintf(conn, ":%d\r\n", n)

		case "SREM":
			if len(parts) < 3 {
				fmt.Fprintf(conn, "-ERR usage: SREM key member [member ...]\r\n")
				continue
			}
			key := parts[1]
			members := parts[2:]
			n := s.store.SRem(key, members...)
			fmt.Fprintf(conn, ":%d\r\n", n)

		case "SMEMBERS":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR usage: SMEMBERS key\r\n")
				continue
			}
			key := parts[1]
			members, err := s.store.SMembers(key)
			if err != nil {
				fmt.Fprintf(conn, "*0\r\n")
			} else {
				fmt.Fprintf(conn, "*%d\r\n", len(members))
				for _, m := range members {
					fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(m), m)
				}
			}
		case "HSET":
			if len(parts) != 4 {
				fmt.Fprintf(conn, "-ERR usage: HSET key field value\r\n")
				continue
			}
			key, field, value := parts[1], parts[2], parts[3]
			n := s.store.HSet(key, field, value)
			fmt.Fprintf(conn, ":%d\r\n", n)
		case "HGET":
			if len(parts) != 3 {
				fmt.Fprintf(conn, "-ERR usage: HGET key field\r\n")
				continue
			}
			key, field := parts[1], parts[2]
			val, ok := s.store.HGet(key, field)
			if !ok {
				fmt.Fprintf(conn, "$-1\r\n")
			} else {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(val), val)
			}
		case "HDEL":
			if len(parts) < 3 {
				fmt.Fprintf(conn, "-ERR usage: HDEL key field [field ...]\r\n")
				continue
			}
			key := parts[1]
			n := s.store.HDel(key, parts[2:]...)
			fmt.Fprintf(conn, ":%d\r\n", n)
		case "HGETALL":
			if len(parts) != 2 {
				fmt.Fprintf(conn, "-ERR usage: HGETALL key\r\n")
				continue
			}
			key := parts[1]
			vals, err := s.store.HGetAll(key)
			if err != nil {
				fmt.Fprintf(conn, "*0\r\n")
				continue
			}
			fmt.Fprintf(conn, "*%d\r\n", len(vals)*2)
			for f, v := range vals {
				fmt.Fprintf(conn, "$%d\r\n%s\r\n$%d\r\n%s\r\n", len(f), f, len(v), v)
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
