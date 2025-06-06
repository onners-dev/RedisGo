package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
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
		conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
		peek, err := reader.Peek(1)
		if err != nil {
			return
		}
		var parts []string
		if peek[0] == '*' {
			// RESP
			parts, err = parseRESP(reader)
			if err != nil {
				conn.Write([]byte(respError("Protocol error: " + err.Error())))
				continue
			}
		} else {
			// Plain line (telnet, nc)
			line, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			parts = strings.Fields(strings.TrimSpace(line))
			if len(parts) == 0 {
				continue
			}
		}

		cmd := strings.ToUpper(parts[0])
		args := parts[1:]

		switch cmd {
		// ---------- Meta ----------
		case "PING":
			if len(args) == 0 {
				conn.Write([]byte(respSimple("PONG")))
			} else {
				conn.Write([]byte(respBulk(args[0])))
			}
		case "ECHO":
			if len(args) != 1 {
				conn.Write([]byte(respError("Wrong number of arguments for 'ECHO'")))
			} else {
				conn.Write([]byte(respBulk(args[0])))
			}
		// ---------- String Commands ----------
		case "SET":
			if len(args) != 2 {
				conn.Write([]byte(respError("Wrong number of arguments for 'SET'")))
				continue
			}
			s.store.Set(args[0], args[1])
			conn.Write([]byte(respSimple("OK")))
		case "GET":
			if len(args) != 1 {
				conn.Write([]byte(respError("Wrong number of arguments for 'GET'")))
				continue
			}
			val, ok := s.store.Get(args[0])
			if !ok {
				conn.Write([]byte(respNullBulk()))
			} else {
				conn.Write([]byte(respBulk(val)))
			}
		case "DEL":
			if len(args) < 1 {
				conn.Write([]byte(respError("Wrong number of arguments for 'DEL'")))
				continue
			}
			deleted := 0
			for _, k := range args {
				if s.store.Del(k) {
					deleted++
				}
			}
			conn.Write([]byte(respInt(deleted)))
		case "INCR":
			if len(args) != 1 {
				conn.Write([]byte(respError("Wrong number of arguments for 'INCR'")))
				continue
			}
			val, err := s.store.Incr(args[0])
			if err != nil {
				conn.Write([]byte(respError(err.Error())))
				continue
			}
			conn.Write([]byte(respInt(val)))
		case "DECR":
			if len(args) != 1 {
				conn.Write([]byte(respError("Wrong number of arguments for 'DECR'")))
				continue
			}
			val, err := s.store.Decr(args[0])
			if err != nil {
				conn.Write([]byte(respError(err.Error())))
				continue
			}
			conn.Write([]byte(respInt(val)))
		case "MSET":
			if len(args) < 2 || len(args)%2 != 0 {
				conn.Write([]byte(respError("MSET requires an even number of arguments (key value ...)")))
				continue
			}
			err := s.store.MSet(args...)
			if err != nil {
				conn.Write([]byte(respError(err.Error())))
			} else {
				conn.Write([]byte(respSimple("OK")))
			}
		case "MGET":
			if len(args) < 1 {
				conn.Write([]byte(respError("Wrong number of arguments for 'MGET'")))
				continue
			}
			vals := s.store.MGet(args...)
			// Convert to RESP array, treating empty string as nil
			items := make([]string, len(vals))
			for i, v := range vals {
				if v == "" {
					items[i] = respNullBulk()
				} else {
					items[i] = respBulk(v)
				}
			}
			// Merge as a prebuilt RESP array:
			resp := "*" + strconv.Itoa(len(items)) + "\r\n" + strings.Join(items, "")
			conn.Write([]byte(resp))
		// ---------- List Commands ----------
		case "LPUSH":
			if len(args) < 2 {
				conn.Write([]byte(respError("LPUSH requires a key and at least one value")))
				continue
			}
			newLen := s.store.LPush(args[0], args[1:]...)
			conn.Write([]byte(respInt(newLen)))
		case "RPOP":
			if len(args) != 1 {
				conn.Write([]byte(respError("RPOP requires a key")))
				continue
			}
			val, err := s.store.RPop(args[0])
			if err != nil {
				conn.Write([]byte(respNullBulk()))
			} else {
				conn.Write([]byte(respBulk(val)))
			}
		case "LLEN":
			if len(args) != 1 {
				conn.Write([]byte(respError("LLEN requires a key")))
				continue
			}
			length := s.store.LLen(args[0])
			conn.Write([]byte(respInt(length)))
		// ---------- Set Commands ----------
		case "SADD":
			if len(args) < 2 {
				conn.Write([]byte(respError("SADD requires a key and at least one value")))
				continue
			}
			n := s.store.SAdd(args[0], args[1:]...)
			conn.Write([]byte(respInt(n)))
		case "SREM":
			if len(args) < 2 {
				conn.Write([]byte(respError("SREM requires a key and at least one value")))
				continue
			}
			n := s.store.SRem(args[0], args[1:]...)
			conn.Write([]byte(respInt(n)))
		case "SMEMBERS":
			if len(args) != 1 {
				conn.Write([]byte(respError("SMEMBERS requires a key")))
				continue
			}
			members, err := s.store.SMembers(args[0])
			if err != nil || len(members) == 0 {
				conn.Write([]byte("*0\r\n"))
			} else {
				conn.Write([]byte(respArray(members)))
			}
		// ---------- Hash Commands ----------
		case "HSET":
			if len(args) != 3 {
				conn.Write([]byte(respError("HSET requires a key, field, value")))
				continue
			}
			added := s.store.HSet(args[0], args[1], args[2])
			conn.Write([]byte(respInt(added)))
		case "HGET":
			if len(args) != 2 {
				conn.Write([]byte(respError("HGET requires a key and field")))
				continue
			}
			val, ok := s.store.HGet(args[0], args[1])
			if !ok {
				conn.Write([]byte(respNullBulk()))
			} else {
				conn.Write([]byte(respBulk(val)))
			}
		case "HDEL":
			if len(args) < 2 {
				conn.Write([]byte(respError("HDEL requires a key and at least one field")))
				continue
			}
			num := s.store.HDel(args[0], args[1:]...)
			conn.Write([]byte(respInt(num)))
		case "HGETALL":
			if len(args) != 1 {
				conn.Write([]byte(respError("HGETALL requires a key")))
				continue
			}
			m, err := s.store.HGetAll(args[0])
			if err != nil || len(m) == 0 {
				conn.Write([]byte("*0\r\n"))
			} else {
				arr := []string{}
				for k, v := range m {
					arr = append(arr, k, v)
				}
				conn.Write([]byte(respArray(arr)))
			}
		// ---------- ZSet Commands ----------
		case "ZADD":
			if len(args) != 3 {
				conn.Write([]byte(respError("ZADD requires a key, score, member")))
				continue
			}
			score, err := strconv.ParseFloat(args[1], 64)
			if err != nil {
				conn.Write([]byte(respError("Invalid score for ZADD")))
				continue
			}
			n := s.store.ZAdd(args[0], score, args[2])
			conn.Write([]byte(respInt(n)))
		case "ZREM":
			if len(args) != 2 {
				conn.Write([]byte(respError("ZREM requires a key and member")))
				continue
			}
			n := s.store.ZRem(args[0], args[1])
			conn.Write([]byte(respInt(n)))
		case "ZRANGE":
			if len(args) != 3 {
				conn.Write([]byte(respError("ZRANGE requires key start stop")))
				continue
			}
			start, err1 := strconv.Atoi(args[1])
			stop, err2 := strconv.Atoi(args[2])
			if err1 != nil || err2 != nil {
				conn.Write([]byte(respError("invalid integer range for ZRANGE")))
				continue
			}
			members, err := s.store.ZRange(args[0], start, stop)
			if err != nil || len(members) == 0 {
				conn.Write([]byte("*0\r\n"))
			} else {
				conn.Write([]byte(respArray(members)))
			}
		// ---------- Key Management ----------
		case "EXPIRE":
			if len(args) != 2 {
				conn.Write([]byte(respError("Wrong number of arguments for 'EXPIRE'")))
				continue
			}
			secs, err := strconv.Atoi(args[1])
			if err != nil {
				conn.Write([]byte(respError("Invalid seconds for 'EXPIRE'")))
				continue
			}
			ok := s.store.Expire(args[0], secs)
			if ok {
				conn.Write([]byte(respInt(1)))
			} else {
				conn.Write([]byte(respInt(0)))
			}
		case "TTL":
			if len(args) != 1 {
				conn.Write([]byte(respError("Wrong number of arguments for 'TTL'")))
				continue
			}
			ttl := s.store.TTL(args[0])
			conn.Write([]byte(respInt(ttl)))
		case "KEYS":
			keys := s.store.Keys()
			conn.Write([]byte(respArray(keys)))
		case "DUMPALL":
			kv := s.store.DumpAll()
			arr := []string{}
			for k, v := range kv {
				arr = append(arr, k, v)
			}
			conn.Write([]byte(respArray(arr)))
		// ---------- HELP / COMMANDS ----------
		case "COMMANDS", "HELP":
			commands := []string{
				"PING", "ECHO message", "SET key value", "GET key", "DEL key [key ...]", "EXPIRE key seconds", "TTL key",
				"INCR key", "DECR key", "MSET key value [key value ...]", "MGET key [key ...]",
				"LPUSH key value [value ...]", "RPOP key", "LLEN key",
				"SADD key member [member ...]", "SREM key member [member ...]", "SMEMBERS key",
				"HSET key field value", "HGET key field", "HDEL key field [field ...]", "HGETALL key",
				"ZADD key score member", "ZREM key member", "ZRANGE key start stop",
				"DUMPALL", "KEYS",
			}
			conn.Write([]byte(respArray(commands)))
		default:
			conn.Write([]byte(respError("unknown command `" + cmd + "`")))
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
