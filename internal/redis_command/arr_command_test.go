package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
	"testing"
)

func TestLPush(t *testing.T) {
	var mu sync.RWMutex
	storedData := make(map[string]model.StoredData)

	tests := []struct {
		cmdArray  []any
		expected  string
		finalList []any
	}{
		{[]any{"LPUSH", "mylist", "world"}, "+OK\r\n", []any{"world"}},
		{[]any{"LPUSH", "mylist", "hello"}, "+OK\r\n", []any{"hello", "world"}},
		{[]any{"LPUSH", "mylist"}, "-ERR missing argument for LPUSH\r\n", []any{"hello", "world"}},
		{[]any{"LPUSH", 123, "hello"}, "-ERR invalid argument for LPUSH\r\n", []any{"hello", "world"}},
	}

	for _, tt := range tests {
		result := LPush(tt.cmdArray, storedData, &mu)
		if result != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, result)
		}
		if tt.expected == "+OK\r\n" {
			list := storedData["mylist"].Value.([]any)
			for i, v := range tt.finalList {
				if list[i] != v {
					t.Errorf("expected list %v, got %v", tt.finalList, list)
				}
			}
		}
	}
}

func TestRPush(t *testing.T) {
	var mu sync.RWMutex
	storedData := make(map[string]model.StoredData)

	tests := []struct {
		cmdArray  []any
		expected  string
		finalList []any
	}{
		{[]any{"RPUSH", "mylist", "hello"}, "+OK\r\n", []any{"hello"}},
		{[]any{"RPUSH", "mylist", "world"}, "+OK\r\n", []any{"hello", "world"}},
		{[]any{"RPUSH", "mylist"}, "-ERR missing argument for RPUSH\r\n", []any{"hello", "world"}},
		{[]any{"RPUSH", 123, "world"}, "-ERR invalid argument for RPUSH\r\n", []any{"hello", "world"}},
	}

	for _, tt := range tests {
		result := RPush(tt.cmdArray, storedData, &mu)
		if result != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, result)
		}
		if tt.expected == "+OK\r\n" {
			list := storedData["mylist"].Value.([]any)
			for i, v := range tt.finalList {
				if list[i] != v {
					t.Errorf("expected list %v, got %v", tt.finalList, list)
				}
			}
		}
	}
}
