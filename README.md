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
## Environment Variables
## API Endpoints
## Example Requests
## What I Learned
