package time

import (
	"encoding/json"
	"fmt"
	"idfm/pkg/env"
	"log"
	"net/http"
	"sync"
	"time"
)

// GetAllTimings retrieves all timings for the given stop IDs
func GetAllTimings(stopIDs []string) ([]map[string]interface{}, error) {
	var allTimings []map[string]interface{}
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for _, stopID := range stopIDs {
		wg.Add(1)
		go func(sid string) {
			defer wg.Done()

			if data, err := requestInfo(sid); err == nil {
				mutex.Lock()
				allTimings = append(allTimings, data...)
				mutex.Unlock()
			} else {
				log.Printf("ERROR: IDFM: Unable to read timings for %s: %v", sid, err)
			}
		}(stopID)
	}

	wg.Wait()
	return allTimings, nil
}

// requestInfo fetches information for a specific stop ID
func requestInfo(stopID string) ([]map[string]interface{}, error) {
	urlStr := fmt.Sprintf("https://prim.iledefrance-mobilites.fr/marketplace/stop-monitoring?MonitoringRef=STIF:StopPoint:Q:%s:", stopID)

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	req, err := http.NewRequest("GET", urlStr, nil)
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

	var result struct {
		Siri struct {
			ServiceDelivery struct {
				StopMonitoringDelivery []struct {
					MonitoredStopVisit []map[string]interface{} `json:"MonitoredStopVisit"`
				} `json:"StopMonitoringDelivery"`
			} `json:"ServiceDelivery"`
		} `json:"Siri"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Siri.ServiceDelivery.StopMonitoringDelivery) == 0 {
		return nil, fmt.Errorf("no stop monitoring delivery data found")
	}

	return result.Siri.ServiceDelivery.StopMonitoringDelivery[0].MonitoredStopVisit, nil
}
