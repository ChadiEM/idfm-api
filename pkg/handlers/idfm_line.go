package handlers

import (
	"github.com/gin-gonic/gin"
	"idfm/pkg/internal/line"
	"idfm/pkg/internal/types"
)

func IDFMLineHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		transportType, err := parseTransportType(c.Param("type"))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		transportId := c.Param("id")

		lineRequest := types.Line{Type: transportType, ID: transportId}

		lineId, err := line.GetLineDetails(lineRequest)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"id": lineId})
	}
}
