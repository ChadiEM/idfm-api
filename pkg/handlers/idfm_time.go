package handlers

import (
	"github.com/gin-gonic/gin"
	"idfm/pkg/internal/line"
	"idfm/pkg/internal/stop"
	"idfm/pkg/internal/time"
	"idfm/pkg/internal/utils"
	"net/http"
)

func IDFMTimeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		transportType, err := validateTransportType(c.Param("type"))
		if err != nil {
			handleGinError(c, err)
			return
		}
		transportId := c.Param("id")
		stopName := c.Param("stop")

		dir := c.Query("direction")
		platform := c.Query("platform")

		lineID, err := line.GetLineDetailsOrCache(transportType, transportId)
		if err != nil {
			handleGinError(c, err)
			return
		}

		var stopIDs []utils.StopId
		stopID, exists := stop.GetCachedStopIDsForDirection(lineID, stopName, dir, platform)
		if exists {
			stopIDs = []utils.StopId{stopID}
		} else {
			stopIDs, err = stop.GetStopIDs(lineID, stopName)
			if err != nil {
				handleGinError(c, err)
				return
			}
		}

		allTimings, err := time.GetAllTimings(stopIDs)
		if err != nil {
			handleGinError(c, err)
			return
		}

		results := time.FindResults(allTimings, lineID, stopIDs, stopName, dir, platform)

		c.JSON(http.StatusOK, results)
	}
}
