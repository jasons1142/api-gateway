package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BackendURLs     []string
	RedisAddr       string
	RateLimit       int
	RateLimitWindow time.Duration
	APIKeys         map[string]bool

	JWTSecret            string
	JWTExpirationMinutes int
}

func LoadConfig() (*Config, error) {
	backendURLs := os.Getenv("BACKEND_URLS")
	if len(backendURLs) == 0 {
		return nil, errors.New("BACKEND_URLs are required")
	}

	backendURLList := strings.Split(backendURLs, ",")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, errors.New("REDIS_ADDR is required")
	}

	rateLimit := os.Getenv("RATE_LIMIT")
	if rateLimit == "" {
		return nil, errors.New("RATE_LIMIT is required")
	}

	rateLimitWindow := os.Getenv("RATE_LIMIT_WINDOW_SECONDS")
	if rateLimitWindow == "" {
		return nil, errors.New("RATE_LIMIT_WINDOW_SECONDS is required")
	}

	apiKeys := os.Getenv("VALID_API_KEYS")
	if apiKeys == "" {
		return nil, errors.New("VALID_API_KEYS is required")
	}

	rateLimitInt, err := strconv.Atoi(rateLimit)
	if err != nil {
		return nil, errors.New("RATE_LIMIT must be a number")
	}

	rateLimitWindowInt, err := strconv.Atoi(rateLimitWindow)
	if err != nil {
		return nil, errors.New("RATE_LIMIT_WINDOW_SECONDS must be a number")
	}

	apiKeyMap := make(map[string]bool)

	for _, key := range strings.Split(apiKeys, ",") {
		apiKeyMap[strings.TrimSpace(key)] = true
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET is required")
	}

	jwtExpirationMinutes := os.Getenv("JWT_EXPIRATION_MINUTES")
	if jwtExpirationMinutes == "" {
		return nil, errors.New("JWT_EXPIRATION_MINUTES is required")
	}

	jwtExpirationMinuteInt, err := strconv.Atoi(jwtExpirationMinutes)
	if err != nil {
		return nil, errors.New("JWT_EXPIRATION_MINUTES must be a number")
	}

	return &Config{
		BackendURLs:          backendURLList,
		RedisAddr:            redisAddr,
		RateLimit:            rateLimitInt,
		APIKeys:              apiKeyMap,
		RateLimitWindow:      time.Duration(rateLimitWindowInt) * time.Second,
		JWTSecret:            jwtSecret,
		JWTExpirationMinutes: jwtExpirationMinuteInt,
	}, nil
}
