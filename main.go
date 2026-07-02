package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

// context allowing Redis client to manage cancellation/timeouts/lifetime
var ctx = context.Background()

// authenticating api key
/*func apiKeyAuth(config *Config) gin.HandlerFunc {
	// valid keys
	validKeys := config.APIKeys

	// checking if they key is valid
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")

		if !validKeys[apiKey] {
			incrementUnauthorizedRequests()

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized: invalid or missing API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
} */

// authenticating tokens
func jwtAuthMiddleware(config *Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		scheme, jwtToken, found := strings.Cut(authHeader, " ")

		if !found || scheme != "Bearer" || jwtToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid || token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
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
			incrementRedisErrors()

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
			incrementRateLimitedRequests()

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
		incrementTotalRequests()

		start := time.Now()

		c.Next()
		status := c.Writer.Status()

		if status >= 200 && status < 300 {
			incrementSuccessfulRequests()
		}

		latency := time.Since(start)
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
	router.GET("/users", requestLogger(), jwtAuthMiddleware(config), redisRateLimiter(config, redisClient), func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"totalRequests":        metrics.TotalRequests,
			"successfulRequests":   metrics.SuccessfulRequests,
			"unauthorizedRequests": metrics.UnauthorizedRequests,
			"rateLimitedRequests":  metrics.RateLimitedRequests,
			"redisErrors":          metrics.RedisErrors,
		})
	})

	router.GET("/health", func(c *gin.Context) {
		// check Redis
		_, redisErr := redisClient.Ping(ctx).Result()

		// check backend
		resp, backendErr := http.Get(config.BackendURL + "/users")
		if resp != nil {
			defer resp.Body.Close()
		}

		if redisErr != nil || backendErr != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "unhealthy",
				"redis":   redisErr == nil,
				"backend": backendErr == nil,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"redis":   true,
			"backend": true,
		})
	})

	type Login struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}

	router.POST("/login", func(c *gin.Context) {
		var login Login

		validUsername := "jason"
		validPassword := "password123"

		if err := c.ShouldBindJSON(&login); err != nil { // 2. Pass the pointer to it
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if login.Username == validUsername && login.Password == validPassword {
			token, tokenErr := generateJWT(config, login.Username)
			if tokenErr == nil {
				c.JSON(http.StatusOK, gin.H{
					"token": token,
				})
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid username or password",
		})

	})

	// start server
	router.Run()
}
