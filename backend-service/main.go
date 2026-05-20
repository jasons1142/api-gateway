package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/users", func(c *gin.Context) {
		// send json reponse
		c.JSON(http.StatusOK, gin.H{
			"users": []string{"Jason", "Alice", "Bob"},
		})
	})

	router.Run(":8081")
}
