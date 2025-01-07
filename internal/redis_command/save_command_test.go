package redis_command

import (
	"os"
	"redis-go-clone/internal/model"
	"sync"
	"testing"
)

func TestSave(t *testing.T) {
	storedData := map[string]model.StoredData{
		"key1": {Value: "value1", ExpiryDate: 0},
		"key2": {Value: 123, ExpiryDate: 0},
	}

	mu := &sync.RWMutex{}

	result := Save([]any{"SAVE"}, storedData, mu)
	if result != "+OK\r\n" {
		t.Errorf("expected +OK\r\n, got %v", result)
	}

	if _, err := os.Stat("data.json"); os.IsNotExist(err) {
		t.Error("data.json file was not created")
	}

	os.Remove("data.json")
}
