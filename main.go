package main

import (
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	// target backend service
	backendURL, _ := url.Parse("http://localhost:8081")

	// reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(backendURL)

	// create web server
	router := gin.Default()

	// if we get ping request run function c
	router.GET("/users", func(c *gin.Context) {
		// send json reponse
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// start server
	router.Run()
}
