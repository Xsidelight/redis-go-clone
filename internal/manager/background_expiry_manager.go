package manager

import (
	"redis-go-clone/internal/model"
	"sync"
	"time"
)

func StartBackgroundExpiryManager(storedData map[string]model.StoredData, mu *sync.RWMutex, interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			mu.Lock()
			now := time.Now().Unix()
			for key, value := range storedData {
				if value.ExpiryDate > 0 && value.ExpiryDate <= now {
					delete(storedData, key)
				}
			}
			mu.Unlock()
		}
	}()
}
