package config

import (
	"redis-go-clone/internal/model"
	"sync"
)

type Config struct {
	DB   map[string]model.StoredData
	Lock *sync.RWMutex
}

func NewConfig() *Config {
	return &Config{
		DB:   make(map[string]model.StoredData),
		Lock: &sync.RWMutex{},
	}
}
