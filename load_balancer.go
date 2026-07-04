package main

import (
	"net/url"
	"sync"
)

type LoadBalancer struct {
	backends []*url.URL
	Current  int
	mu       sync.Mutex
}

func (lb *LoadBalancer) NextBackend() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	temp := lb.Current

	lb.Current++

	if lb.Current == len(lb.backends) {
		lb.Current = 0
	}

	return lb.backends[temp]
}
