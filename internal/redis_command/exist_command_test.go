package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
	"testing"
)

func TestExist(t *testing.T) {
	tests := []struct {
		name       string
		cmdArray   []any
		storedData map[string]model.StoredData
		expected   string
	}{
		{
			name:       "missing argument",
			cmdArray:   []any{"EXIST"},
			storedData: map[string]model.StoredData{},
			expected:   "-ERR missing argument for EXIST\r\n",
		},
		{
			name:       "invalid argument type",
			cmdArray:   []any{"EXIST", 123},
			storedData: map[string]model.StoredData{},
			expected:   "-ERR invalid argument for EXIST\r\n",
		},
		{
			name:       "key does not exist",
			cmdArray:   []any{"EXIST", "key1"},
			storedData: map[string]model.StoredData{},
			expected:   ":0\r\n",
		},
		{
			name:     "key exists",
			cmdArray: []any{"EXIST", "key1"},
			storedData: map[string]model.StoredData{
				"key1": {Value: "value1"},
			},
			expected: ":1\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &sync.RWMutex{}
			result := Exist(tt.cmdArray, tt.storedData, mu)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}