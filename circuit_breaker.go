package main

import (
	"net/url"
	"sync"
	"time"
)

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type Backend struct {
	URL *url.URL

	FailureCount int
	State        CircuitState
	OpenedAt     time.Time

	mu sync.Mutex
}

func (b *Backend) IsAvailable(coolDown time.Duration) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.State == CircuitClosed {
		return true
	}

	if b.State == CircuitOpen {
		if time.Since(b.OpenedAt) < coolDown {
			return false
		}

		b.State = CircuitHalfOpen
		return true

	}

	if b.State == CircuitHalfOpen {
		return false
	}

	return false

}

func (b *Backend) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.FailureCount = 0

	b.State = CircuitClosed

	b.OpenedAt = time.Time{}
}

func (b *Backend) RecordFailure(failureThreshold int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.State == CircuitOpen {
		return
	}

	b.FailureCount++

	if b.FailureCount >= failureThreshold {
		b.State = CircuitOpen
		b.OpenedAt = time.Now()
	}
}
