package localcache

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

type LocalCache struct {
	kvs map[string]string
}

func New() *LocalCache {
	kvs := make(map[string]string)
	return &LocalCache{kvs}
}

func (l *LocalCache) HealthCheck() error {
	return nil
}

func (l *LocalCache) Read(key string) (string, error) {
	mu.Lock()
	value, ok := l.kvs[key]
	mu.Unlock()
	if !ok {
		return "", fmt.Errorf("no such key")
	}
	return value, nil
}

func (l *LocalCache) Write(key, value string) error {
	mu.Lock()
	l.kvs[key] = value
	mu.Unlock()
	return nil
}

func (l *LocalCache) Remove(key string) error {
	mu.Lock()
	delete(l.kvs, key)
	mu.Unlock()
	return nil
}
