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

| Command                           | Description                                   | Example usage                  | Example response             |
|------------------------------------|-----------------------------------------------|-------------------------------|------------------------------|
| `SET key value`                   | Set string value for a key                    | `SET foo bar`                  | `+OK`                        |
| `GET key`                         | Get string value for a key                    | `GET foo`                      | `$3`<br>`bar`                |
| `DEL key [key ...]`               | Delete one or more keys                       | `DEL foo`                      | `:1` (number deleted)        |
| `EXPIRE key seconds`              | Set expiry in seconds for a key               | `EXPIRE foo 10`                | `:1` (success)               |
| `TTL key`                         | Get time-to-live in seconds                   | `TTL foo`                      | `:9`                         |
| `INCR key`                        | Increment integer value                       | `INCR counter`                 | `:1`                         |
| `DECR key`                        | Decrement integer value                       | `DECR counter`                 | `:0`                         |
| `KEYS`                            | List all non-expired keys                     | `KEYS`                         | `*1`<br>`$7`<br>`counter`    |
| `DUMPALL`                         | Get all string keys and values                | `DUMPALL`                      | `*1 ...`                     |
| `MSET key value [key value ...]`  | Set multiple string keys at once              | `MSET a 1 b 2`                 | `+OK`                        |
| `MGET key [key ...]`              | Get multiple string values                    | `MGET a b missing`             | `*3 ...`                     |
| `LPUSH key value [value ...]`     | Prepend one/more items to a list              | `LPUSH list a b c`             | `:3`                         |
| `RPOP key`                        | Remove and return the last item from a list   | `RPOP list`                    | `$1`<br>`a`                  |
| `LLEN key`                        | Get the number of items in a list             | `LLEN list`                    | `:2`                         |
| `SADD key member [member ...]`    | Add one/more items to a set                   | `SADD myset x y y`             | `:2` (added, unique)         |
| `SREM key member [member ...]`    | Remove one/more items from a set              | `SREM myset x`                 | `:1` (removed count)         |
| `SMEMBERS key`                    | Get all members of a set                      | `SMEMBERS myset`               | `*1`<br>`$1`<br>`y`          |
| `PING`                            | Test connection                               | `PING`                         | `PONG`                       |
| `HSET key field value`       | Set field in hash              | `HSET h foo bar`          | `:1`      |
| `HGET key field`             | Get field from hash            | `HGET h foo`              | `$3`<br>`bar` |
| `HDEL key field [field ...]` | Delete field(s) in hash        | `HDEL h foo`              | `:1`      |
| `HGETALL key`                | Get all fields/values in hash  | `HGETALL h`               | `*2 ...`  |


### Coming Soon/To-Do
- Persistence:
- Append-Only File (AOF) logging
- Snapshotting (RDB-like periodic dump)
- Webapp to showcase features

### Running Tests
```sh
go test
```


### You can now use redis-cli or library clients:

```sh
redis-cli -p 6379
> SET foo bar
OK
> GET foo
"bar"
> COMMANDS
1) "PING"
2) "ECHO message"
3) "SET key value"
...etc
