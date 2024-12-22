# Redis Go Clone

A Redis clone implemented in Go as a solution to the [Coding Challenges Redis Challenge](https://codingchallenges.fyi/challenges/challenge-redis).

## Features

- In-memory key-value store
- String data type support
- Basic Redis commands (GET, SET, DEL)
- TTL support for keys
- Simple networking protocol
- RESP (Redis Serialization Protocol) implementation

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/redis-go-clone.git

# Navigate to project directory
cd redis-go-clone

# Build the project
go build
```

## Usage

### Starting the Server

```bash
./redis-go-clone
```

The server will start listening on port 6379 by default.

### Example Commands

Connect to the server using any Redis client, or use telnet:

```bash
telnet localhost 6379
```

Basic commands:

```
SET key value
GET key
DEL key
EXPIRE key seconds
```

## Development

### Requirements

- Go 1.16 or higher

### Running Tests

```bash
go test ./...
```

## About

This project was created as a solution to the [Build Your Own Redis Challenge](https://codingchallenges.fyi/challenges/challenge-redis), where the goal is to implement a basic Redis clone to learn about networking protocols and in-memory databases.
