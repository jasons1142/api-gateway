package main

import (
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
