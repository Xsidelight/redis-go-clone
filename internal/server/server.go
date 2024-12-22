package server

import (
	"fmt"
	"log"
	"net"
)

func StartServer() {
	//err := http.ListenAndServe(":6379", nil)
	//if err != nil {
	//	panic("error starting the server")
	//}
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

		go HandleClient(conn)
	}
}
