package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "backend-service"
	}

	router.GET("/users", func(c *gin.Context) {
		// send json reponse
		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"users":   []string{"Jason", "Alice", "Bob"},
		})
	})

	router.Run(":8081")
}
