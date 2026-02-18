package main

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Store interface {
	Inc(ip string) (int64, error)
	UniqueCount() (int64, error)
}

type RedisStore struct {
	client *redis.Client
}

type MemoryStore struct {
	reg map[string]int64
	mu  sync.RWMutex
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (c *RedisStore) Inc(ip string) (int64, error) {
	newValue, err := c.client.HIncrBy(context.TODO(), "ip", ip, 1).Result()
	if err != nil {
		return 0, err
	}
	return newValue, nil
}

func (c *RedisStore) UniqueCount() (int64, error) {
	return c.client.HLen(context.TODO(), "ip").Result()
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		reg: make(map[string]int64),
	}
}

func (ms *MemoryStore) Inc(ip string) (int64, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.reg[ip]++
	return ms.reg[ip], nil
}

func (ms *MemoryStore) UniqueCount() (int64, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return int64(len(ms.reg)), nil
}
