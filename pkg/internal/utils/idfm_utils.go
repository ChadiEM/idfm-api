package utils

import (
	"fmt"
	"idfm/pkg/cache"
	"idfm/pkg/internal/types"
	"math"
	"regexp"
	"strings"
	"time"
)

var (
	METRO = types.TransportType{API: "metro"}
	BUS   = types.TransportType{API: "bus"}
	RER   = types.TransportType{API: "rail"}
	TRAM  = types.TransportType{API: "tram"}
)

var onlyNumberRegex = regexp.MustCompile(`[0-9]+`)

// GetStopIDs retrieves stop IDs for the given stop
func GetStopIDs(lineId string, stopName string) ([]string, error) {
	stopIDs := make([]string, 0)
	stopsAtRoute := make(map[string][]string)

	routeIDWithPrefix := fmt.Sprintf("IDFM:%s", lineId)

	err := cache.ProcessStopCache(func(row cache.StopCacheData) (bool, error) {
		// Early filter - only process rows for this route
		if row.LineID != routeIDWithPrefix {
			return true, nil
		}

		// Track all stops for this route
		if stopsAtRoute[lineId] == nil {
			stopsAtRoute[lineId] = make([]string, 0)
		}
		stopsAtRoute[lineId] = append(stopsAtRoute[lineId], row.StopName)

		// Only process further if this is the stop we're looking for
		if row.StopName == stopName {
			// Extract numeric stop ID
			curStopID := onlyNumberRegex.FindString(row.StopID)

			stopIDs = append(stopIDs, curStopID)
		}

		// We can have at most two stops with the same name (A & R) on the same route
		if len(stopIDs) == 2 {
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return stopIDs, nil
}

// FindResults processes entries and requests to find matching results
func FindResults(entries []map[string]interface{}, transport types.Transport, lineId string, stopIds []string) []types.Result {
	requests := createRequests(transport, lineId, stopIds)

	results := make([]types.Result, 0)

	for _, request := range requests {
		for _, entry := range entries {
			// Extract DirectionName
			mvj, ok := entry["MonitoredVehicleJourney"].(map[string]interface{})
			if !ok {
				continue
			}

			dirRefMap, ok := mvj["DirectionRef"].(map[string]interface{})
			if !ok {
				continue
			}

			dirRefValue, ok := dirRefMap["value"].(string)
			if !ok {
				continue
			}

			var dir string
			if dirRefValue == "Aller" {
				dir = "A"
			} else if dirRefValue == "Retour" {
				dir = "R"
			} else {
				continue
			}

			directionNames, ok := mvj["DirectionName"].([]interface{})
			if !ok || len(directionNames) == 0 {
				continue
			}

			dirNameMap, ok := directionNames[0].(map[string]interface{})
			if !ok {
				continue
			}

			dirNameValue, ok := dirNameMap["value"].(string)
			if !ok {
				continue
			}

			// Process direction name
			dirName := regexp.MustCompile(`[^a-z]`).ReplaceAllString(strings.ToLower(dirNameValue), "")

			// Extract route and stop IDs
			lineRef, ok := mvj["LineRef"].(map[string]interface{})
			if !ok {
				continue
			}

			routeIDFull, ok := lineRef["value"].(string)
			if !ok {
				continue
			}
			routeID := strings.TrimSuffix(strings.TrimPrefix(routeIDFull, "STIF:Line::"), ":")

			monitoringRef, ok := entry["MonitoringRef"].(map[string]interface{})
			if !ok {
				continue
			}

			stopIDFull, ok := monitoringRef["value"].(string)
			if !ok {
				continue
			}
			stopID := onlyNumberRegex.FindString(stopIDFull)

			// Check if this entry matches our query
			queriedDirName := regexp.MustCompile(`[^a-z]`).ReplaceAllString(
				strings.ToLower(request.Direction), "")
			if request.RouteID != routeID || request.StopID != stopID ||
				!(request.Direction == dir || queriedDirName == dirName) {
				continue
			}

			// Extract destination name
			destNames, ok := mvj["DestinationName"].([]interface{})
			if !ok || len(destNames) == 0 {
				continue
			}

			destNameMap, ok := destNames[0].(map[string]interface{})
			if !ok {
				continue
			}

			destName, ok := destNameMap["value"].(string)
			if !ok {
				continue
			}

			// Calculate remaining time
			var remainingTime string
			monitoredCall, ok := mvj["MonitoredCall"].(map[string]interface{})
			if !ok {
				continue
			}

			atStop, ok := monitoredCall["VehicleAtStop"].(bool)
			if ok && atStop {
				remainingTime = "Arrêt"
			} else {
				expectedTime, ok := monitoredCall["ExpectedDepartureTime"].(string)
				if !ok {
					continue
				}

				upcoming, err := time.Parse(time.RFC3339, expectedTime)
				if err != nil {
					continue
				}

				remaining := int(math.Max(0, math.Floor(upcoming.Sub(time.Now()).Minutes())))
				remainingTime = fmt.Sprintf("%d mn", remaining)
			}

			// Store result
			if results == nil {
				results = make([]types.Result, 0)
			}

			results = append(results, types.Result{
				Dest: destName,
				Time: remainingTime,
			})

			// Update cache
			cache.StopIdForDirectionCacheLock.Lock()
			cache.StopIdForDirectionCache[request.RouteID+"-"+request.StopName+"-"+request.Direction] = stopID
			cache.StopIdForDirectionCacheLock.Unlock()
		}
	}

	return results
}

// createRequests creates a map of transport requests from the provided transports, line IDs, and stop name map
func createRequests(transport types.Transport, lineID string, stopIds []string) []types.Request {
	requests := make([]types.Request, 0)

	for _, stopID := range stopIds {
		requests = append(requests, types.Request{
			RouteID:   lineID,
			StopID:    stopID,
			StopName:  transport.Stop,
			Direction: transport.Destination,
		})
	}

	return requests
}

// ConfigurationError represents configuration-related errors
type ConfigurationError struct {
	Message string
}

func (e *ConfigurationError) Error() string {
	return e.Message
}
