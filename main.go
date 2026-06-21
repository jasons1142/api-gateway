package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

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

func main() {
	// target backend service
	backendURL, _ := url.Parse("http://localhost:8081")

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// create web server
	router := gin.Default()

	// if we get ping request and key is valid run function c
	router.GET("/users", apiKeyAuth(), func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// start server
	router.Run()
}
