package main

import (
	"sync"
)

type Store interface {
	Inc(ip string) (int, error)
	UniqueCount() (int, error)
}

type MemoryStore struct {
	reg map[string]int
	mu  sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		reg: make(map[string]int),
	}
}

func (ms *MemoryStore) Inc(ip string) (int, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.reg[ip]++
	return ms.reg[ip], nil
}

func (ms *MemoryStore) UniqueCount() (int, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return len(ms.reg), nil
}
