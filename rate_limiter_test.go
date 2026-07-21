package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRateLimiter(t *testing.T) {
	s := miniredis.RunT(t)
	defer s.Close()

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	_, redisErr := redisClient.Ping(ctx).Result()

	if redisErr != nil {
		t.Fatalf("failed to connect to test redis: %v", redisErr)
	}

	config := &Config{
		RateLimit:       3,
		RateLimitWindow: 60,
	}

	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.GET(
		"/test",
		redisRateLimiter(config, redisClient),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "request allowed",
			})
		},
	)

	for i := 1; i <= 4; i++ {
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		req.Header.Set("x-api-key", "test-api-key")

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if i <= config.RateLimit {
			assert.Equal(t, http.StatusOK, rr.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}
}

func TestRateLimiterIndependentClients(t *testing.T) {
	s := miniredis.RunT(t)

	redisClient := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	ctx := context.Background()

	err := redisClient.Ping(ctx).Err()
	require.NoError(t, err)

	config := &Config{
		RateLimit:       3,
		RateLimitWindow: 60,
	}

	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.GET(
		"/test",
		redisRateLimiter(config, redisClient),
		func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "request allowed",
			})
		},
	)

	for i := 1; i <= 4; i++ {
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		req.Header.Set("x-api-key", "client-a")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if i <= config.RateLimit {
			assert.Equal(
				t,
				http.StatusOK,
				rr.Code,
				"client A request %d should be allowed",
				i,
			)
		} else {
			assert.Equal(
				t,
				http.StatusTooManyRequests,
				rr.Code,
				"client A request %d should be blocked",
				i,
			)
		}
	}

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	require.NoError(t, err)

	req.Header.Set("x-api-key", "client-b")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(
		t,
		http.StatusOK,
		rr.Code,
		"client B should not be blocked by client A's limit",
	)
}
