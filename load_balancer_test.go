package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadBalancerRoundRobin(t *testing.T) {
	backend1, _ := url.Parse("http://backend-service-1:8081")
	backend2, _ := url.Parse("http://backend-service-2:8081")

	lb := &LoadBalancer{
		backends: []*url.URL{
			backend1,
			backend2,
		},
		Current: 0,
	}

	check1 := lb.NextBackend()
	assert.Equal(t, check1, backend1)

	check2 := lb.NextBackend()
	assert.Equal(t, check2, backend2)

	check3 := lb.NextBackend()
	assert.Equal(t, backend1, check3)
}
