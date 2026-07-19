// Copyright (c) BeduSec. All rights reserved.
package store

import (
	"fmt"
	"time"

	"github.com/bedusec/mago/pkg/config"
)

type TokenBucketState struct {
	Tokens     float64
	LastRefill time.Time
}

type Store interface {
	GetBucket(key string) (*TokenBucketState, error)
	UpdateBucket(key string, state *TokenBucketState) error
}

func New(cfg config.StoreConfig) (Store, error) {
	switch cfg.Type {
	case "memory":
		return NewMemoryStore(), nil
	case "redis":
		return NewRedisStore(cfg.RedisURL)
	default:
		return nil, fmt.Errorf("unknown store type: %s", cfg.Type)
	}
}