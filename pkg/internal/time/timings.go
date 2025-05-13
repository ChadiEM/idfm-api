package time

import (
	"encoding/json"
	"fmt"
	"idfm/pkg/env"
	"idfm/pkg/internal/utils"
	"net/http"
	"net/url"
	"time"
)

const (
	stopMonitoringEndpoint = "https://prim.iledefrance-mobilites.fr/marketplace/stop-monitoring"
)

// StopMonitoringAPIResponse represents the structure of the API response
type StopMonitoringAPIResponse struct {
	Siri Siri `json:"Siri"`
}

// GetAllTimings retrieves all timings for the given stop IDs with typed data
func GetAllTimings(stopIDs []utils.StopId) ([]MonitoredStopVisit, error) {
	var allTimings []MonitoredStopVisit

	for _, stopID := range stopIDs {
		if data, err := requestInfo(stopID); err == nil {
			allTimings = append(allTimings, data...)
		} else {
			return nil, err
		}
	}

	return allTimings, nil
}

// requestInfo fetches information for a specific stop ID
func requestInfo(stopID utils.StopId) ([]MonitoredStopVisit, error) {
	params := url.Values{}
	if stopID.Type == utils.Area {
		params.Add("MonitoringRef", fmt.Sprintf("STIF:StopArea:SP:%s:", stopID.Id))
	} else if stopID.Type == utils.Point {
		params.Add("MonitoringRef", fmt.Sprintf("STIF:StopPoint:Q:%s:", stopID.Id))
	} else {
		return nil, fmt.Errorf("invalid stop ID type: %s", stopID.Type)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	req, err := http.NewRequest("GET", stopMonitoringEndpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("apiKey", env.IDFM_API_KEY)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result StopMonitoringAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Siri.ServiceDelivery.StopMonitoringDelivery) == 0 {
		return nil, fmt.Errorf("no stop monitoring delivery data found")
	}

	return result.Siri.ServiceDelivery.StopMonitoringDelivery[0].MonitoredStopVisit, nil
}
