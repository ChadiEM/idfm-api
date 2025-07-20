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
	Dest     string `json:"dest"`
	Time     string `json:"time"`
	Status   string `json:"status"`
	Platform string `json:"platform,omitempty"`
}

// FindResults processes entries and requests to find matching results
func FindResults(entries []MonitoredStopVisit, lineId string, stopIds []utils.StopId, stopName string, destination string, platform string) []Result {
	results := make([]Result, 0)

	for _, requestedStopId := range stopIds {
		for _, entry := range entries {
			// Check LineRef
			lineRefValue := entry.MonitoredVehicleJourney.LineRef.Value
			if lineRefValue != fmt.Sprintf("STIF:Line::%s:", lineId) {
				continue
			}

			// Check Platform
			if platform != "" && platform != entry.MonitoredVehicleJourney.MonitoredCall.ArrivalPlatformName.Value {
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
			} else if len(entry.MonitoredVehicleJourney.DirectionName) > 0 && entry.MonitoredVehicleJourney.DirectionName[0].Value == "Aller" {
				dir = "A"
			} else if len(entry.MonitoredVehicleJourney.DirectionName) > 0 && entry.MonitoredVehicleJourney.DirectionName[0].Value == "Retour" {
				dir = "R"
			}

			if destination != "" && destination != dir {
				continue
			}

			// Check stop ID
			stopIDFull := entry.MonitoringRef.Value
			stopID := utils.OnlyNumberRegex.FindString(stopIDFull)

			// there's a bug where the returned stop id is different from the requested one...
			// workaround this bug by assuming it is the same if the number of requested stops is 1
			// example: /api/idfm/timings/rail/A/Auber?direction=A
			if len(stopIds) > 1 {
				if stopID != requestedStopId.Id {
					continue
				}
			}

			// Calculate remaining time
			var remainingTime string
			if entry.MonitoredVehicleJourney.MonitoredCall.VehicleAtStop {
				remainingTime = "onStop"
			} else {
				upcoming := entry.MonitoredVehicleJourney.MonitoredCall.ExpectedDepartureTime
				remaining := int(math.Max(0, math.Floor(upcoming.Sub(time.Now()).Minutes())))
				remainingTime = fmt.Sprintf("%d mn", remaining)
			}

			// Store result
			results = append(results, Result{
				Dest:     entry.MonitoredVehicleJourney.DestinationName[0].Value,
				Time:     remainingTime,
				Status:   entry.MonitoredVehicleJourney.MonitoredCall.DepartureStatus,
				Platform: entry.MonitoredVehicleJourney.MonitoredCall.ArrivalPlatformName.Value,
			})

			// Update cache
			if destination != "" || platform != "" {
				stopCacheKey := data.StopCacheKey{
					LineId:    lineId,
					StopName:  stopName,
					Direction: destination,
					Platform:  platform,
				}
				data.StopIdForDirectionCache.Set(stopCacheKey, requestedStopId, ttlcache.DefaultTTL)
			}
		}
	}

	return results
}
