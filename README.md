# Go API Gateway

A containerized API Gateway built in Go that handles authentication, rate limiting, request logging, metrics, health checks, and reverse proxying to backend services.

## Features

- Reverse proxy to backend services
- JWT authentication with signed tokens
- API key authentication support
- Redis-backed rate limiting
- Request logging with method, path, status, and latency
- Metrics endpoint for gateway activity
- Health check endpoint for gateway dependencies
- Environment-based configuration
- Docker Compose multi-service setup

## Tech Stack

- Go
- Gin
- Redis
- Docker
- Docker Compose
- JWT

## Architecture

                     Client
                        │
                        ▼
                 API Gateway
                        │
       ┌────────┬─────────────┬─────────┐
       ▼        ▼             ▼         ▼
    JWT Auth  Rate Limit   Logging   Metrics
                  │
                  ▼
               Redis
                  │
                  ▼
         Round Robin Load Balancer
             │                │
             ▼                ▼
       Backend 1         Backend 2
## Getting Started
   ### Prerequisites
    - Docker Desktop
    - Git
  ### Clone Repository
    - clone https://github.com/jasons1142/api-gateway.git
    - cd api-gateway

  ### Build and Start
    - docker compose up --build
  
  ### Environment Variables
    BACKEND_URLS=http://backend-service-1:8081,http://backend-service-2:8081
    
    REDIS_ADDR=redis:6379
    
    RATE_LIMIT=5
    RATE_LIMIT_WINDOW_SECONDS=60
    
    VALID_API_KEYS=test-key
    
    JWT_SECRET=my-super-secret-key
    JWT_EXPIRATION_MINUTES=30
## Example Requests
  ### Generate Token
    - echo $TOKEN
  ### Access a protected endpoint
    - curl \ -H "Authorization: Bearer <TOKEN>" \http://localhost:8080/users
  ### Metrics
    - curl http://localhost:8080/metrics
  ### Health
    - curl http://localhost:8080/health
  ### Stop the project
    - Docker compose donw
    
## What I Learned
