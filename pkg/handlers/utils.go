package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"idfm/pkg/internal/utils"
	"net/http"
	"slices"
)

func validateTransportType(transportType string) (string, error) {
	if slices.Contains(utils.AllowedTransportTypes, transportType) {
		return transportType, nil
	}
	return "", &utils.RequestError{Message: fmt.Sprintf("Invalid transport type: %s. Valid types: %s", transportType, utils.AllowedTransportTypes)}
}

func handleGinError(c *gin.Context, err error) {
	var requestError *utils.RequestError
	if errors.As(err, &requestError) {
		c.JSON(http.StatusBadRequest, gin.H{"request error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	return
}
