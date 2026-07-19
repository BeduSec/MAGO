// Copyright (c) BeduSec. All rights reserved.
package store

import (
	"sync"
	"time"
)

type MemoryStore struct {
	mu    sync.Mutex
	items map[string]*TokenBucketState
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]*TokenBucketState),
	}
}

func (m *MemoryStore) GetBucket(key string) (*TokenBucketState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	state, ok := m.items[key]
	if !ok {
		return &TokenBucketState{Tokens: 0, LastRefill: time.Now()}, nil
	}
	copy := *state
	return &copy, nil
}

func (m *MemoryStore) UpdateBucket(key string, state *TokenBucketState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = state
	return nil
}