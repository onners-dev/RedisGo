package main

import (
	"log"
)

func main() {
	store := NewStore()
	server := NewServer(store)
	if err := server.Listen(":6379"); err != nil {
		log.Fatal(err)
	}
}
