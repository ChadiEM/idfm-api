package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
	"idfm/pkg/cache"
	"idfm/pkg/env"
	"idfm/pkg/handlers"
	"log"
	"net/http"
	"time"
)

// Define a rate limiter with 5 requests per second and a burst of 10.
var limiter = rate.NewLimiter(5, 10)

// Prometheus metrics
var (
	cacheLastUpdateTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "idfm_cache_last_update_timestamp",
		Help: "Timestamp of the last cache update",
	}, []string{"cache_type"})

	cacheSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "idfm_cache_size",
		Help: "Size of the cache",
	}, []string{"cache_type"})

	cacheFetchTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "idfm_cache_fetch_time",
		Help: "Time it took to fetch the cache",
	}, []string{"cache_type"})
)

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

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// API group
	idfm := r.Group("/api/idfm")
	{
		idfm.GET("/lines/:type/:id", handlers.IDFMLineHandler())
		idfm.GET("/timings/:type/:id/:stop/:dir", handlers.IDFMTimeHandler())
	}

	// Start updating cache metrics in the background
	go updateCacheMetrics()

	r.Run()
}

func updateCacheMetrics() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C

		// Get cache stats
		linesCacheLastUpdate, linesCacheSize, linesCacheFetchTime := cache.GetLinesCacheStats()
		stopsCacheLastUpdate, stopsCacheSize, stopsCacheFetchTime := cache.GetStopsCacheStats()

		// Update cache metrics
		cacheLastUpdateTime.WithLabelValues("lines").Set(float64(linesCacheLastUpdate.Unix()))
		cacheLastUpdateTime.WithLabelValues("stops").Set(float64(stopsCacheLastUpdate.Unix()))

		// Update cache size metrics
		cacheSize.WithLabelValues("lines").Set(float64(linesCacheSize))
		cacheSize.WithLabelValues("stops").Set(float64(stopsCacheSize))

		// Update cache fetch time metrics
		cacheFetchTime.WithLabelValues("line").Set(linesCacheFetchTime.Seconds())
		cacheFetchTime.WithLabelValues("stops").Set(stopsCacheFetchTime.Seconds())

	}
}
