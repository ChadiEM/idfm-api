package handlers

import (
	"idfm/pkg/internal/types"
	"idfm/pkg/internal/utils"
)

func parseTransportType(transportType string) (types.TransportType, error) {
	var t types.TransportType

	if transportType == "bus" {
		t = utils.BUS
	} else if transportType == "metro" {
		t = utils.METRO
	} else if transportType == "rer" {
		t = utils.RER
	} else if transportType == "tram" {
		t = utils.TRAM
	} else {
		return t, &utils.ConfigurationError{Message: "Invalid transport type " + transportType}
	}

	return t, nil
}
