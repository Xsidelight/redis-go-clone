package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
	"testing"
)

func TestIncr(t *testing.T) {
	tests := []struct {
		name       string
		cmdArray   []any
		storedData map[string]model.StoredData
		expected   string
	}{
		{
			name:       "missing argument",
			cmdArray:   []any{"INCR"},
			storedData: map[string]model.StoredData{},
			expected:   "-ERR missing argument for INCR\r\n",
		},
		{
			name:       "invalid argument type",
			cmdArray:   []any{"INCR", 123},
			storedData: map[string]model.StoredData{},
			expected:   "-ERR invalid argument for INCR\r\n",
		},
		{
			name:       "key does not exist",
			cmdArray:   []any{"INCR", "counter"},
			storedData: map[string]model.StoredData{},
			expected:   "-ERR key does not exist\r\n",
		},
		{
			name:     "value is not int",
			cmdArray: []any{"INCR", "counter"},
			storedData: map[string]model.StoredData{
				"counter": {Value: "not an int"},
			},
			expected: "-ERR value is not type of int\r\n",
		},
		{
			name:     "successful increment",
			cmdArray: []any{"INCR", "counter"},
			storedData: map[string]model.StoredData{
				"counter": {Value: 1},
			},
			expected: "+OK\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &sync.RWMutex{}
			result := Incr(tt.cmdArray, tt.storedData, mu)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestDecr(t *testing.T) {
	tests := []struct {
		name       string
		cmdArray   []any
		storedData map[string]model.StoredData
		expected   string
	}{
		{
			name:     "missing argument",
			cmdArray: []any{"DECR"},
			expected: "-ERR missing argument for DECR\r\n",
		},
		{
			name:     "invalid argument type",
			cmdArray: []any{"DECR", 123},
			expected: "-ERR invalid argument for DECR\r\n",
		},
		{
			name:       "key does not exist",
			cmdArray:   []any{"DECR", "nonexistent"},
			storedData: map[string]model.StoredData{},
			expected:   "-ERR key does not exist\r\n",
		},
		{
			name:     "value is not int",
			cmdArray: []any{"DECR", "key"},
			storedData: map[string]model.StoredData{
				"key": {Value: "string"},
			},
			expected: "-ERR value is not type of int\r\n",
		},
		{
			name:     "successful decrement",
			cmdArray: []any{"DECR", "key"},
			storedData: map[string]model.StoredData{
				"key": {Value: 10},
			},
			expected: "+OK\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &sync.RWMutex{}
			result := Decr(tt.cmdArray, tt.storedData, mu)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
