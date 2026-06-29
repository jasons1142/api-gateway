package main

import "sync/atomic"

type Metrics struct {
	TotalRequests        int64
	SuccessfulRequests   int64
	UnauthorizedRequests int64
	RateLimitedRequests  int64
	RedisErrors          int64
}

var metrics = &Metrics{}

func incrementTotalRequests() {
	atomic.AddInt64(&metrics.TotalRequests, 1)
}

func incrementSuccessfulRequests() {
	atomic.AddInt64(&metrics.SuccessfulRequests, 1)
}

func incrementUnauthorizedRequests() {
	atomic.AddInt64(&metrics.UnauthorizedRequests, 1)
}

func incrementRateLimitedRequests() {
	atomic.AddInt64(&metrics.RateLimitedRequests, 1)
}

func incrementRedisErrors() {
	atomic.AddInt64(&metrics.RedisErrors, 1)
}
