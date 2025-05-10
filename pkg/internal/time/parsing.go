package time

import (
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"idfm/pkg/data"
	"idfm/pkg/internal/utils"
	"math"
	"strings"
	"time"
)

// Result represents a transport timing result
type Result struct {
	Dest string `json:"dest"`
	Time string `json:"time"`
}

// FindResults processes entries and requests to find matching results
func FindResults(entries []MonitoredStopVisit, lineId string, stopIds []string, stopName string, destination string) []Result {
	results := make([]Result, 0)

	for _, requestedStopId := range stopIds {
		for _, entry := range entries {
			// Check LineRef
			lineRefValue := entry.MonitoredVehicleJourney.LineRef.Value
			if lineRefValue != fmt.Sprintf("STIF:Line::%s:", lineId) {
				continue
			}

			// Check Direction
			dirRefValue := entry.MonitoredVehicleJourney.DirectionRef.Value
			var dir string
			if dirRefValue == "Aller" {
				dir = "A"
			} else if dirRefValue == "Retour" {
				dir = "R"
			} else if strings.HasSuffix(dirRefValue, ":A") {
				dir = "A"
			} else if strings.HasSuffix(dirRefValue, ":R") {
				dir = "R"
			} else {
				continue
			}

			if destination != dir {
				continue
			}

			// Check stop ID
			stopIDFull := entry.MonitoringRef.Value
			stopID := utils.OnlyNumberRegex.FindString(stopIDFull)

			if stopID != requestedStopId {
				continue
			}

			// Get destination name
			if len(entry.MonitoredVehicleJourney.DestinationName) == 0 {
				continue
			}
			destName := entry.MonitoredVehicleJourney.DestinationName[0].Value

			// Calculate remaining time
			var remainingTime string
			if entry.MonitoredVehicleJourney.MonitoredCall.VehicleAtStop {
				remainingTime = "Arrêt"
			} else {
				expectedTime := entry.MonitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime
				upcoming, err := time.Parse(time.RFC3339, expectedTime)
				if err != nil {
					continue
				}

				remaining := int(math.Max(0, math.Floor(upcoming.Sub(time.Now()).Minutes())))
				remainingTime = fmt.Sprintf("%d mn", remaining)
			}

			// Store result
			results = append(results, Result{
				Dest: destName,
				Time: remainingTime,
			})

			// Update cache
			data.StopIdForDirectionCache.Set(lineId+"-"+stopName+"-"+destination, requestedStopId, ttlcache.DefaultTTL)
		}
	}

	return results
}
