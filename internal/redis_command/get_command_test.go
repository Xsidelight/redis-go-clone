package redis_command

import (
	"redis-go-clone/internal/model"
	"sync"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		cmdArray   []any
		storedData map[string]model.StoredData
		want       string
	}{
		{
			name:     "missing argument",
			cmdArray: []any{"GET"},
			want:     "-ERR missing argument for GET\r\n",
		},
		{
			name:     "invalid argument type",
			cmdArray: []any{"GET", 123},
			want:     "-ERR invalid argument for GET\r\n",
		},
		{
			name:     "key not found",
			cmdArray: []any{"GET", "nonexistent"},
			want:     "$-1\r\n",
		},
		{
			name:     "expired key",
			cmdArray: []any{"GET", "expired"},
			storedData: map[string]model.StoredData{
				"expired": {
					Value:      "value",
					ExpiryDate: time.Now().Add(-1 * time.Hour).Unix(),
				},
			},
			want: "$-1\r\n",
		},
		{
			name:     "get existing string",
			cmdArray: []any{"GET", "key1"},
			storedData: map[string]model.StoredData{
				"key1": {Value: "something like this one"},
			},
			want: "$23\r\nsomething like this one\r\n",
		},
		{
			name:     "non-expired key",
			cmdArray: []any{"GET", "valid"},
			storedData: map[string]model.StoredData{
				"valid": {
					Value:      "value",
					ExpiryDate: time.Now().Add(1 * time.Hour).Unix(),
				},
			},
			want: "$5\r\nvalue\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mu := &sync.RWMutex{}
			if tt.storedData == nil {
				tt.storedData = make(map[string]model.StoredData)
			}
			if got := Get(tt.cmdArray, tt.storedData, mu); got != tt.want {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
