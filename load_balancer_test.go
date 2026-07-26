package main

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadBalancerRoundRobin(t *testing.T) {
	backend1URL, _ := url.Parse("http://backend-service-1:8081")
	backend2URL, _ := url.Parse("http://backend-service-2:8081")

	backend1 := &Backend{
		URL:   backend1URL,
		State: CircuitClosed,
	}

	backend2 := &Backend{
		URL:   backend2URL,
		State: CircuitClosed,
	}

	lb := &LoadBalancer{
		backends: []*Backend{
			backend1,
			backend2,
		},
		Current: 0,
	}

	cooldown := 30 * time.Second

	check1 := lb.NextBackend(cooldown)
	assert.Equal(t, backend1, check1)

	check2 := lb.NextBackend(cooldown)
	assert.Equal(t, backend2, check2)

	check3 := lb.NextBackend(cooldown)
	assert.Equal(t, backend1, check3)
}
