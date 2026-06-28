package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// context allowing Redis client to manage cancellation/timeouts/lifetime
var ctx = context.Background()

// authenticating api key
func apiKeyAuth(config *Config) gin.HandlerFunc {
	// valid keys
	validKeys := config.APIKeys

	// checking if they key is valid
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")

		if !validKeys[apiKey] {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// limiting large number of requests
func redisRateLimiter(config *Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")

		redisKey := "rate_limit:" + apiKey

		count, err := redisClient.Incr(ctx, redisKey).Result()

		if err != nil {
			fmt.Println("Redis error:", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to communicate with Redis",
			})
			c.Abort()
			return
		}

		if count == 1 {
			redisClient.Expire(ctx, redisKey, config.RateLimitWindow)
		}

		if int64(config.RateLimit) < count {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// logging requests made
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		apiKey := c.GetHeader("x-api-key")
		method := c.Request.Method
		path := c.Request.URL.Path
		timestamp := start.Format("2006-01-02 15:04:05")

		fmt.Printf("[%s] | Latency: %s | Status: %d | API Key: %s | Method: %s | Path: %s\n", timestamp, latency, status, apiKey, method, path)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// loading in config struct
	config, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// connection to redis
	var redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	// target backend service
	backendURL, _ := url.Parse(config.BackendURL)

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// create web server
	router := gin.Default()

	// if we get ping request and key is valid run function c
	router.GET("/users", requestLogger(), apiKeyAuth(config), redisRateLimiter(config, redisClient), func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// start server
	router.Run()
}
