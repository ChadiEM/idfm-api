package data

import (
	"github.com/jellydator/ttlcache/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

var (
	TypeAndNumberToLineNameCache = ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](12*time.Hour),
		ttlcache.WithCapacity[string, string](100),
	)
	StopIdForDirectionCache = ttlcache.New[string, string](
		ttlcache.WithTTL[string, string](12*time.Hour),
		ttlcache.WithCapacity[string, string](1000),
	)
)

func InitCache() {
	go TypeAndNumberToLineNameCache.Start()
	go StopIdForDirectionCache.Start()

	// Prometheus metrics
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "idfm",
		Name:      "cache_size",
		Help:      "Size of the cache",
		ConstLabels: prometheus.Labels{
			"type": "lines",
		},
	}, func() float64 {
		return float64(TypeAndNumberToLineNameCache.Len())
	})

	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Namespace: "idfm",
		Name:      "cache_size",
		Help:      "Size of the cache",
		ConstLabels: prometheus.Labels{
			"type": "stops",
		},
	}, func() float64 {
		return float64(StopIdForDirectionCache.Len())
	})
}
