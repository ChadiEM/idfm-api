package data

import (
	"github.com/jellydator/ttlcache/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"idfm/pkg/internal/utils"
	"time"
)

type LineCacheKey struct {
	LineType string
	LineId   string
}

type StopCacheKey struct {
	LineId    string
	StopName  string
	Direction string
	Platform  string
}

var (
	TypeAndNumberToLineNameCache = ttlcache.New[LineCacheKey, string](
		ttlcache.WithTTL[LineCacheKey, string](12*time.Hour),
		ttlcache.WithCapacity[LineCacheKey, string](100),
	)
	StopIdForDirectionCache = ttlcache.New[StopCacheKey, utils.StopId](
		ttlcache.WithTTL[StopCacheKey, utils.StopId](12*time.Hour),
		ttlcache.WithCapacity[StopCacheKey, utils.StopId](1000),
	)
)

func registerCacheSizeMetric[K comparable, V any](cacheType string, cache *ttlcache.Cache[K, V]) {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "idfm",
		Name:      "cache_size",
		Help:      "Size of the cache",
		ConstLabels: prometheus.Labels{
			"type": cacheType,
		},
	}, func() float64 {
		return float64(cache.Len())
	})
}

func registerCacheHitMetric[K comparable, V any](cacheType string, cache *ttlcache.Cache[K, V]) {
	promauto.NewCounterFunc(prometheus.CounterOpts{
		Namespace: "idfm",
		Name:      "cache_hits",
		Help:      "Cache hits",
		ConstLabels: prometheus.Labels{
			"type": cacheType,
		},
	}, func() float64 {
		return float64(cache.Metrics().Hits)
	})
}

func registerCacheMissMetric[K comparable, V any](cacheType string, cache *ttlcache.Cache[K, V]) {
	promauto.NewCounterFunc(prometheus.CounterOpts{
		Namespace: "idfm",
		Name:      "cache_misses",
		Help:      "Cache misses",
		ConstLabels: prometheus.Labels{
			"type": cacheType,
		},
	}, func() float64 {
		return float64(cache.Metrics().Misses)
	})
}

func registerCacheInsertionsMetric[K comparable, V any](cacheType string, cache *ttlcache.Cache[K, V]) {
	promauto.NewCounterFunc(prometheus.CounterOpts{
		Namespace: "idfm",
		Name:      "cache_insertions",
		Help:      "Cache insertions",
		ConstLabels: prometheus.Labels{
			"type": cacheType,
		},
	}, func() float64 {
		return float64(cache.Metrics().Insertions)
	})
}

func registerCacheEvictionsMetric[K comparable, V any](cacheType string, cache *ttlcache.Cache[K, V]) {
	promauto.NewCounterFunc(prometheus.CounterOpts{
		Namespace: "idfm",
		Name:      "cache_evictions",
		Help:      "Cache evictions",
		ConstLabels: prometheus.Labels{
			"type": cacheType,
		},
	}, func() float64 {
		return float64(cache.Metrics().Evictions)
	})
}

func InitCache() {
	go TypeAndNumberToLineNameCache.Start()
	go StopIdForDirectionCache.Start()

	// Prometheus metrics
	registerCacheSizeMetric("stops", StopIdForDirectionCache)
	registerCacheSizeMetric("lines", TypeAndNumberToLineNameCache)

	registerCacheHitMetric("stops", StopIdForDirectionCache)
	registerCacheHitMetric("lines", TypeAndNumberToLineNameCache)

	registerCacheMissMetric("stops", StopIdForDirectionCache)
	registerCacheMissMetric("lines", TypeAndNumberToLineNameCache)

	registerCacheInsertionsMetric("stops", StopIdForDirectionCache)
	registerCacheInsertionsMetric("lines", TypeAndNumberToLineNameCache)

	registerCacheEvictionsMetric("stops", StopIdForDirectionCache)
	registerCacheEvictionsMetric("lines", TypeAndNumberToLineNameCache)
}
