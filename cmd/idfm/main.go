package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"idfm/pkg/data"
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

	r := gin.Default()

	r.Use(rateLimiter)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// API group
	idfm := r.Group("/api/idfm")
	{
		idfm.GET("/lines/:type/:id", handlers.IDFMLineHandler())
		idfm.GET("/timings/:type/:id/:stop", handlers.IDFMTimeHandler())
	}

	data.InitCache()

	r.Run()
}
