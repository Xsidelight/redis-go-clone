package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
	"testing"
)

func TestDel(t *testing.T) {
	tests := []struct {
		name       string
		cmdArray   []any
		storedData map[string]model.StoredData
		expected   string
	}{
		{
			name:     "missing argument",
			cmdArray: []any{"DEL"},
			expected: "-ERR missing argument for DEL\r\n",
		},
		{
			name:     "invalid argument",
			cmdArray: []any{"DEL", 123},
			expected: "-ERR invalid argument for DEL\r\n",
		},
		{
			name:     "delete existing key",
			cmdArray: []any{"DEL", "key1"},
			storedData: map[string]model.StoredData{
				"key1": {Value: "value1"},
			},
			expected: ":1\r\n",
		},
		{
			name:     "delete non-existing key",
			cmdArray: []any{"DEL", "key2"},
			storedData: map[string]model.StoredData{
				"key1": {Value: "value1"},
			},
			expected: ":0\r\n",
		},
		{
			name:     "delete multiple keys",
			cmdArray: []any{"DEL", "key1", "key2"},
			storedData: map[string]model.StoredData{
				"key1": {Value: "value1"},
				"key2": {Value: "value2"},
				"key3": {Value: "value3"},
			},
			expected: ":2\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &sync.RWMutex{}
			result := Del(tt.cmdArray, tt.storedData, mu)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
