package manager

import (
	"redis-go-clone/internal/model"
	"sync"
	"testing"
	"time"
)

func TestStartBackgroundExpiryManager(t *testing.T) {
	testCases := []struct {
		name     string
		data     map[string]model.StoredData
		interval time.Duration
		wait     time.Duration
		want     int
	}{
		{
			name: "Expires single key",
			data: map[string]model.StoredData{
				"key1": {Value: "value1", ExpiryDate: time.Now().Add(100 * time.Millisecond).Unix()},
			},
			interval: 50 * time.Millisecond,
			wait:     200 * time.Millisecond,
			want:     0,
		},
		{
			name: "Multiple keys with different expiry",
			data: map[string]model.StoredData{
				"key1": {Value: "value1", ExpiryDate: time.Now().Add(100 * time.Millisecond).Unix()},
				"key2": {Value: "value2", ExpiryDate: time.Now().Add(500 * time.Millisecond).Unix()},
				"key3": {Value: "value3", ExpiryDate: 0}, // No expiry
			},
			interval: 50 * time.Millisecond,
			wait:     200 * time.Millisecond,
			want:     2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storedData := tc.data
			mu := &sync.RWMutex{}

			StartBackgroundExpiryManager(storedData, mu, tc.interval)
			time.Sleep(tc.wait)

			mu.RLock()
			if len(storedData) != tc.want {
				t.Errorf("Expected %d keys, got %d", tc.want, len(storedData))
			}
			mu.RUnlock()
		})
	}
}
