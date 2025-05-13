package stop

import (
	"encoding/json"
	"fmt"
	"idfm/pkg/data"
	"idfm/pkg/internal/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
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
func GetCachedStopIDsForDirection(lineId string, stopName string, direction string, platform string) (utils.StopId, bool) {
	stopCacheKey := data.StopCacheKey{
		LineId:    lineId,
		StopName:  stopName,
		Direction: direction,
		Platform:  platform,
	}

	cacheItem := data.StopIdForDirectionCache.Get(stopCacheKey)
	if cacheItem != nil && !cacheItem.IsExpired() {
		return cacheItem.Value(), true
	}

	return utils.StopId{}, false
}

// GetStopIDs retrieves stop IDs for the given stop from IDFM API
func GetStopIDs(lineId string, stopName string) ([]utils.StopId, error) {
	stopIdsResponse, err := requestStopIds(lineId, stopName)
	if err != nil {
		return nil, err
	}

	if stopIdsResponse.TotalCount > 0 {
		stopIDs := make([]utils.StopId, stopIdsResponse.TotalCount)

		for index, result := range stopIdsResponse.Results {
			stopId := result.StopID
			numericPart := utils.OnlyNumberRegex.FindString(stopId)

			if strings.Contains(stopId, "monomodalStopPlace") {
				// monomodal means that we should query the area instead of the stop...
				stopIDs[index] = utils.StopId{
					Id:   numericPart,
					Type: utils.Area,
				}
			} else {
				stopIDs[index] = utils.StopId{
					Id:   numericPart,
					Type: utils.Point,
				}
			}
		}

		return stopIDs, nil
	} else {
		// Help the user by providing stop names
		allStopNamesResponse, err := requestAllStopNames(lineId)
		if err != nil {
			return nil, err
		}

		stopNames := make([]string, len(allStopNamesResponse.Results))

		for index, result := range allStopNamesResponse.Results {
			stopNames[index] = result.StopName
		}

		marshal, err := json.Marshal(stopNames)
		if err != nil {
			return nil, err
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
		return apiResp, err
	}
	return apiResp, nil
}
