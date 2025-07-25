package handlers

import (
	"github.com/gin-gonic/gin"
	"idfm/pkg/internal/line"
	"net/http"
)

func IDFMLineHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		transportType, err := validateTransportType(c.Param("type"))
		if err != nil {
			handleGinError(c, err)
			return
		}
		transportId := c.Param("id")

		operator := c.Query("operator")

		lineID, err := line.GetLineDetailsOrCache(transportType, transportId, operator)
		if err != nil {
			handleGinError(c, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": lineID})
	}
}
