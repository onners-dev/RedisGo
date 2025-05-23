# RedisGo

A minimal, educational Redis clone in Go — featuring in-memory, TCP-accessible storage with basic Redis-like commands and expiry. Perfect for learning, contributing, or hacking on new database features!

---

## Features

- Persistent, thread-safe in-memory string storage
- Redis-style commands: `SET`, `GET`, `DEL`, `EXPIRE`, `TTL`, `INCR`, `DECR`, `KEYS`, `DUMPALL`
- Time-to-live (TTL) and key expiry
- Atomic integer operations via `INCR`/`DECR`
- Simple, readable codebase and extensive unit tests

---

## Installation

1. **Clone the repo**
   ```sh
   git clone https://github.com/yourusername/RedisGo.git
   cd RedisGo
    ```

2. **Install Go if you don't have it**
    https://go.dev/doc/install 

3. **Build and run**
    ```sh
    go run .
    ```


## Usage

```sh
telnet localhost 6379
```

or use nc
```sh
nc localhost 6379
```

### Commands

SET foo bar
+OK

GET foo
$3
bar

EXPIRE foo 10
:1

TTL foo
:10

INCR counter
:1

DECR counter
:0

DEL foo
:1

KEYS
*1
$7
counter

DUMPALL
*1
$7
counter
$1
0

### Coming Soon/To-Do
- Persistence:
-  Append-Only File (AOF) logging
- Snapshotting (RDB-like periodic dump)
- Data Structures:
- Lists (LPUSH, RPOP, etc.)
- Sets (SADD, SREM, SMEMBERS, etc.)
- Hashes (HSET, HGET, etc.)
- Sorted Sets (ZADD, ZRANGE)

### Running Tests
```sh
go test
```
