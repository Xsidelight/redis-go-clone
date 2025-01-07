# Redis Go Clone

This project is a Redis clone implemented in Go. It was created as a solution to the [Build Your Own Redis Challenge](https://codingchallenges.fyi/challenges/challenge-redis). The goal of this project is to learn about networking protocols and in-memory databases by building a basic version of Redis.

## Features

- **Core Commands:**
  - `SET`: Set the string value of a key.
  - `GET`: Get the value of a key.
  - `DEL`: Delete one or more keys.
  - `EXIST`: Check if a key exists.
  - `LPUSH`: Prepend one or multiple values to a list.
  - `RPUSH`: Append one or multiple values to a list.
  - `INCR`: Increment the integer value of a key by one.
  - `DECR`: Decrement the integer value of a key by one.
  - `SAVE`: Persist the current database state to disk.

- **Persistence:**
  - **SAVE:** Save the in-memory database state to a JSON file (`data.json`).
  - **LOAD:** Automatically load the database state from `data.json` on startup.

- **Concurrency:**
  - Thread-safe operations using `sync.RWMutex`.

- **Testing:**
  - Comprehensive unit tests for all commands and functionalities.

## Installation

### Prerequisites

- **Go:** Version 1.16 or higher.

### Steps

1. **Clone the Repository:**

    ```sh
    git clone https://github.com/yourusername/redis-go-clone.git
    cd redis-go-clone
    ```

2. **Build the Project:**

    ```sh
    go build
    ```

3. **Run the Server:**

    ```sh
    ./redis-go-clone
    ```

    The server listens on port `6379` by default.

## Usage

Connect to the server using a Redis client, such as `redis-cli`:

```sh
redis-cli -h localhost -p 6379
```

## License

This project is licensed under the MIT License.
