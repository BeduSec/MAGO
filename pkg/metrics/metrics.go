// Copyright (c) BeduSec. All rights reserved.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mago_requests_total",
		Help: "Total number of requests processed.",
	}, []string{"method", "path", "status"})
	BlockedRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mago_blocked_requests_total",
		Help: "Total number of requests blocked by WAF.",
	}, []string{"rule_id"})
	RateLimitedRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mago_ratelimited_requests_total",
		Help: "Total number of requests rate limited.",
	})
	RequestLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "mago_request_duration_seconds",
		Help:    "Request latency in seconds.",
		Buckets: prometheus.DefBuckets,
	})
)