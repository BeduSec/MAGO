// Copyright (c) BeduSec. All rights reserved.
package store

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(redisURL string) (*RedisStore, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("redis parse url: %w", err)
	}
	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}
	return &RedisStore{client: client}, nil
}

func (r *RedisStore) GetBucket(key string) (*TokenBucketState, error) {
	ctx := context.Background()
	pipe := r.client.Pipeline()
	tokenCmd := pipe.Get(ctx, "bucket:"+key+":tokens")
	refillCmd := pipe.Get(ctx, "bucket:"+key+":last_refill")
	_, err := pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	tokens, _ := strconv.ParseFloat(tokenCmd.Val(), 64)
	lastRefill, _ := time.Parse(time.RFC3339Nano, refillCmd.Val())
	return &TokenBucketState{Tokens: tokens, LastRefill: lastRefill}, nil
}

func (r *RedisStore) UpdateBucket(key string, state *TokenBucketState) error {
	ctx := context.Background()
	pipe := r.client.Pipeline()
	pipe.Set(ctx, "bucket:"+key+":tokens", fmt.Sprintf("%f", state.Tokens), 0)
	pipe.Set(ctx, "bucket:"+key+":last_refill", state.LastRefill.Format(time.RFC3339Nano), 0)
	_, err := pipe.Exec(ctx)
	return err
}