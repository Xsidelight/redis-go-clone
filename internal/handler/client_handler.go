package handler

import (
	"log"
	"net"
	"redis-go-clone/pkg/resp"
	"strings"
	"sync"
)

var storedData = make(map[string]any)
var mu sync.RWMutex

func HandleClient(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 4096) // Buffer for reading client data
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from client: %v", err)
			return
		}

		// Deserialize the RESP message
		input := string(buffer[:n])
		command, err := resp.DeserializeRESP(input)
		if err != nil {
			log.Printf("Invalid RESP message: %v", err)
			conn.Write([]byte("-ERR invalid RESP message\r\n"))
			continue
		}

		// Process the command
		response := processCommand(command)
		conn.Write([]byte(response))
	}
}

func processCommand(command any) string {
	// Ensure the command is an array
	cmdArray, ok := command.([]any)
	if !ok || len(cmdArray) == 0 {
		return "-ERR invalid command\r\n"
	}

	// First element is the command name
	cmdName, ok := cmdArray[0].(string)
	if !ok {
		return "-ERR invalid command name\r\n"
	}
	cmdName = strings.ToUpper(cmdName)

	// Handle supported commands
	switch cmdName {
	case "SET":
		if len(cmdArray) < 3 {
			return "-ERR missing argument for SET\r\n"
		}
		key, ok := cmdArray[1].(string)
		if !ok {
			return "-ERR invalid argument for SET\r\n"
		}
		mu.Lock()
		defer mu.Unlock()
		value := cmdArray[2]
		storedData[key] = value

		return "+OK\r\n"
	case "GET":
		if len(cmdArray) < 2 {
			return "-ERR missing argument for GET\r\n"
		}
		key, ok := cmdArray[1].(string)
		if !ok {
			return "-ERR invalid argument for GET\r\n"
		}

		mu.Lock()
		defer mu.Unlock()
		value, ok := storedData[key]

		if !ok {
			return "$-1\r\n"
		}
		return resp.SerializeRESP(value)
	default:
		return "-ERR unknown command\r\n"
	}
}
