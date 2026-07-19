// Copyright (c) BeduSec. All rights reserved.
package limiter

import (
	"math"
	"time"

	"github.com/bedusec/mago/internal/store"
)

type RateLimiter struct {
	store  store.Store
	rate   float64
	burst  int
}

func NewRateLimiter(st store.Store, rate float64, burst int) *RateLimiter {
	return &RateLimiter{
		store: st,
		rate:  rate,
		burst: burst,
	}
}

func (rl *RateLimiter) Allow(key string) (bool, int, time.Duration, error) {
	state, err := rl.store.GetBucket(key)
	if err != nil {
		return false, 0, 0, err
	}

	now := time.Now()
	elapsed := now.Sub(state.LastRefill).Seconds()
	newTokens := elapsed * rl.rate
	tokens := math.Min(state.Tokens+newTokens, float64(rl.burst))
	if tokens < 1 {
		waitTime := time.Duration((1-tokens)/rl.rate * float64(time.Second))
		return false, 0, waitTime, nil
	}
	tokens--
	state.Tokens = tokens
	state.LastRefill = now

	if err := rl.store.UpdateBucket(key, state); err != nil {
		return false, 0, 0, err
	}

	remaining := int(math.Floor(tokens))
	return true, remaining, 0, nil
}