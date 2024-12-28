# Redis Go Clone

This project is a Redis clone implemented in Go. It was created as a solution to the [Build Your Own Redis Challenge](https://codingchallenges.fyi/challenges/challenge-redis). The goal of this project is to learn about networking protocols and in-memory databases by building a basic version of Redis.

## Features

- Basic Redis commands (GET, SET, DEL, etc.)
- In-memory data storage
- Simple networking protocol implementation

## Getting Started

### Prerequisites

- Go 1.16 or higher

### Installation

Clone the repository:

```sh
git clone https://github.com/yourusername/redis-go-clone.git
cd redis-go-clone
```

Build the project:

```sh
go build
```

Run the server:

```sh
./redis-go-clone
```

### Usage

Connect to the server using a Redis client:

```sh
redis-cli -h localhost -p 6379
```

## License

This project is licensed under the MIT License.
