package stop

import (
	"encoding/json"
	"fmt"
	"idfm/pkg/data"
	"idfm/pkg/internal/utils"
	"io"
	"net/http"
	"net/url"
)

const (
	stopRecordsEndpoint = "https://data.iledefrance-mobilites.fr/api/explore/v2.1/catalog/datasets/arrets-lignes/records"
)

type stopIdsAPIResponse struct {
	TotalCount int `json:"total_count"`
	Results    []struct {
		StopID string `json:"stop_id"`
	} `json:"results"`
}

type stopNamesAPIResponse struct {
	TotalCount int `json:"total_count"`
	Results    []struct {
		StopName string `json:"stop_name"`
	} `json:"results"`
}

// GetCachedStopIDsForDirection retrieves stop IDs for the given stop and direction from the cache
func GetCachedStopIDsForDirection(lineId string, stopName string, direction string) (string, bool) {
	stopCacheKey := lineId + "-" + stopName + "-" + direction

	cacheItem := data.StopIdForDirectionCache.Get(stopCacheKey)
	if cacheItem != nil && !cacheItem.IsExpired() {
		return cacheItem.Value(), true
	}

	return "", false
}

// GetStopIDs retrieves stop IDs for the given stop from IDFM API
func GetStopIDs(lineId string, stopName string) ([]string, error) {
	stopIdsResponse, err := requestStopIds(lineId, stopName)
	if err != nil {
		return nil, err
	}

	if stopIdsResponse.TotalCount > 0 {
		stopIDs := make([]string, stopIdsResponse.TotalCount)

		for index, result := range stopIdsResponse.Results {
			stopIDs[index] = utils.OnlyNumberRegex.FindString(result.StopID)
		}

		return stopIDs, nil
	} else {
		// Help the user by providing stop names
		allStopNamesResponse, err := requestAllStopNames(lineId)
		if err != nil {
			return []string{}, err
		}

		stopNames := make([]string, allStopNamesResponse.TotalCount)

		for index, result := range allStopNamesResponse.Results {
			stopNames[index] = result.StopName
		}

		marshal, err := json.Marshal(stopNames)
		if err != nil {
			return []string{}, err
		}
		return nil, &utils.RequestError{Message: fmt.Sprintf("Stop \"%s\" not found. Available stops: %s", stopName, marshal)}
	}
}

func requestStopIds(lineId string, stopName string) (stopIdsAPIResponse, error) {
	// Prepare query parameters
	params := url.Values{}
	params.Add("select", "stop_id")
	params.Add("where", fmt.Sprintf("id=\"IDFM:%s\" AND stop_name=\"%s\"", lineId, stopName))

	resp, err := http.Get(stopRecordsEndpoint + "?" + params.Encode())
	if err != nil {
		return stopIdsAPIResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return stopIdsAPIResponse{}, err
	}

	var apiResp stopIdsAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return stopIdsAPIResponse{}, err
	}
	return apiResp, nil
}

func requestAllStopNames(lineId string) (stopNamesAPIResponse, error) {
	// Prepare query parameters
	params := url.Values{}
	params.Add("select", "stop_name")
	params.Add("where", fmt.Sprintf("id=\"IDFM:%s\"", lineId))
	params.Add("limit", "100")
	resp, err := http.Get(stopRecordsEndpoint + "?" + params.Encode())
	if err != nil {
		return stopNamesAPIResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return stopNamesAPIResponse{}, err
	}

	var apiResp stopNamesAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return stopNamesAPIResponse{}, err
	}
	return apiResp, nil
}
