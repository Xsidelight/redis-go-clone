package server

import (
	"fmt"
	"log"
	"net"
	"redis-go-clone/internal/handler"
)

func StartServer() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Failed to close server: %v", err)
		}
	}(listener)

	fmt.Println("Redis Lite server listening on port 6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handler.HandleClient(conn)
	}
}
