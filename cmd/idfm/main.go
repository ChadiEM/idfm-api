package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"idfm/pkg/cache"
	"idfm/pkg/env"
	"idfm/pkg/handlers"
	"log"
	"net/http"
)

// Define a rate limiter with 5 requests per second and a burst of 10.
var limiter = rate.NewLimiter(5, 10)

// Middleware to check the rate limit.
func rateLimiter(c *gin.Context) {
	if !limiter.Allow() {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	if env.IDFM_API_KEY == "" {
		log.Fatal("IDFM_API_KEY not defined. Please create an API key and try again.")
	}

	// Initialize the CSV cache system first
	cache.InitializeCaches()

	r := gin.Default()

	r.Use(rateLimiter)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	r.GET("/api/idfm/lines/:type/:id", handlers.IDFMLineHandler())
	r.GET("/api/idfm/timings/:type/:id/:stop/:dir", handlers.IDFMTimeHandler())

	r.Run()
}
