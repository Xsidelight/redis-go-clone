package server

import (
	"fmt"
	"log"
	"net"
	"redis-go-clone/cmd/config"
	"redis-go-clone/internal/handler"
	"redis-go-clone/internal/manager"
	"time"
)

func StartServer() {
	config := config.NewConfig()

	manager.LoadData(config.DB, config.Lock)

	manager.StartBackgroundExpiryManager(config.DB, config.Lock, 1*time.Second)

	h := handler.NewClientHandler(config)

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

		go h.HandleClient(conn)
	}
}
