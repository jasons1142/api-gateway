package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimitInfo struct {
	Count       int
	WindowStart time.Time
}

var rateLimits = map[string]RateLimitInfo{}

// authenticating api key
func apiKeyAuth() gin.HandlerFunc {
	// valid keys
	validKeys := map[string]bool{
		"abc123":        true,
		"my-secret-key": true,
	}

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
func rateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")

		limit := 5
		window := time.Minute
		now := time.Now()

		info, exists := rateLimits[apiKey]

		if !exists || now.Sub(info.WindowStart) > window {
			info = RateLimitInfo{
				Count:       1,
				WindowStart: now,
			}
		} else {
			info.Count++
		}

		rateLimits[apiKey] = info

		if info.Count > limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too Many Requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	// target backend service
	backendURL, _ := url.Parse("http://localhost:8081")

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// create web server
	router := gin.Default()

	// if we get ping request and key is valid run function c
	router.GET("/users", apiKeyAuth(), rateLimiter(), func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// start server
	router.Run()
}
