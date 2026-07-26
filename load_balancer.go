package main

import (
	"sync"
	"time"
)

type LoadBalancer struct {
	backends []*Backend
	Current  int
	mu       sync.Mutex
}

func (lb *LoadBalancer) NextBackend(cooldown time.Duration) *Backend {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	for i := 0; i < len(lb.backends); i++ {
		backend := lb.backends[lb.Current]

		lb.Current++

		if lb.Current >= len(lb.backends) {
			lb.Current = 0
		}

		if backend.IsAvailable(cooldown) {
			return backend
		}
	}

	return nil
}
