package handlers

import (
	"github.com/gin-gonic/gin"
	"idfm/pkg/cache"
	"idfm/pkg/internal/line"
	"idfm/pkg/internal/time"
	"idfm/pkg/internal/types"
	"idfm/pkg/internal/utils"
	"log"
)

func IDFMTimeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		transportType, err := parseTransportType(c.Param("type"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		transportId := c.Param("id")
		stopName := c.Param("stop")
		dir := c.Param("dir")

		lineRequest := types.Line{Type: transportType, ID: transportId}

		lineCacheKey := transportType.API + "-" + transportId
		var lineID string
		if cache.TypeAndNumberToLineNameCache[lineCacheKey] != "" {
			lineID = cache.TypeAndNumberToLineNameCache[lineCacheKey]
		} else {
			resLineId, err := line.GetLineDetails(lineRequest)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			lineID = resLineId
			cache.TypeAndNumberToLineNameCacheLock.Lock()
			cache.TypeAndNumberToLineNameCache[lineCacheKey] = resLineId
			cache.TypeAndNumberToLineNameCacheLock.Unlock()
		}

		stopCacheKey := lineID + "-" + stopName + "-" + dir
		var stopIDs []string
		if cachedStopIDs, exists := cache.StopIdForDirectionCache[stopCacheKey]; exists {
			stopIDs = []string{cachedStopIDs}
		} else {
			curStopIDs, err := utils.GetStopIDs(lineID, stopName)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			stopIDs = curStopIDs
		}

		allTimings, err := time.GetAllTimings(stopIDs)
		if err != nil {
			log.Printf("Error getting timings: %v", err)
		}

		transport := types.Transport{Type: transportType, Number: transportId, Stop: stopName, Destination: dir}
		results := utils.FindResults(allTimings, transport, lineID, stopIDs)

		c.JSON(200, results)
	}
}
