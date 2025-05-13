package line

import (
	"encoding/json"
	"fmt"
	"github.com/jellydator/ttlcache/v3"
	"idfm/pkg/data"
	"idfm/pkg/internal/utils"
	"io"
	"net/http"
	"net/url"
)

const (
	lineRecordsEndpoint = "https://data.iledefrance-mobilites.fr/api/explore/v2.1/catalog/datasets/referentiel-des-lignes/records"
)

type linesAPIResponse struct {
	TotalCount int `json:"total_count"`
	Results    []struct {
		IDLine string `json:"id_line"`
	} `json:"results"`
}

type allLinesAPIResponse struct {
	TotalCount int `json:"total_count"`
	Results    []struct {
		NameLine string `json:"name_line"`
	} `json:"results"`
}

// GetLineDetailsOrCache retrieves line details from the cache/API
func GetLineDetailsOrCache(lineType string, lineId string) (string, error) {
	lineCacheKey := data.LineCacheKey{LineType: lineType, LineId: lineId}
	cacheItem := data.TypeAndNumberToLineNameCache.Get(lineCacheKey)
	if cacheItem != nil && !cacheItem.IsExpired() {
		return cacheItem.Value(), nil
	}

	// Prepare query parameters
	params := url.Values{}
	params.Add("select", "id_line")
	params.Add("where",
		fmt.Sprintf("transportmode=\"%s\" AND name_line=\"%s\" AND (operatorname=\"RATP\" OR operatorname=\"SNCF\")", lineType, lineId))

	resp, err := http.Get(lineRecordsEndpoint + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResp linesAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", err
	}

	if apiResp.TotalCount == 0 {
		allLines, err := getAllLines(lineType)
		if err != nil {
			return "", err
		}

		lineIds := make([]string, len(allLines.Results))

		for index, result := range allLines.Results {
			lineIds[index] = result.NameLine
		}

		marshal, err := json.Marshal(lineIds)
		if err != nil {
			return "", err
		}
		return "", &utils.RequestError{Message: fmt.Sprintf("%s \"%s\" not found. Available lines: %s", lineType, lineId, marshal)}
	} else if apiResp.TotalCount == 1 {
		resLineId := apiResp.Results[0].IDLine
		data.TypeAndNumberToLineNameCache.Set(lineCacheKey, resLineId, ttlcache.DefaultTTL)
		return resLineId, nil
	}

	return "", &utils.RequestError{Message: fmt.Sprintf("Invalid transport %s %s", lineType, lineId)}
}

// getAllLines retrieves all lines for that type
func getAllLines(lineType string) (allLinesAPIResponse, error) {
	// Prepare query parameters
	params := url.Values{}
	params.Add("select", "name_line")
	params.Add("where",
		fmt.Sprintf("transportmode=\"%s\" AND (operatorname=\"RATP\" OR operatorname=\"SNCF\")", lineType))
	params.Add("limit", "100")

	resp, err := http.Get(lineRecordsEndpoint + "?" + params.Encode())
	if err != nil {
		return allLinesAPIResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return allLinesAPIResponse{}, err
	}

	var apiResp allLinesAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return apiResp, err
	}

	return apiResp, nil
}
