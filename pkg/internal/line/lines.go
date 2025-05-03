package line

import (
	"fmt"
	"idfm/pkg/cache"
	"idfm/pkg/internal/types"
	"idfm/pkg/internal/utils"
)

// GetLineDetails retrieves line details from the API
func GetLineDetails(line types.Line) (string, error) {
	var id string

	err := cache.ProcessLineCache(func(data cache.LineCacheData) (bool, error) {
		if data.TransportMode == line.Type.API &&
			data.ShortNameLine == line.ID &&
			(data.OperatorName == "RATP" || data.OperatorName == "SNCF") {
			id = data.IDLine
			return false, nil

		}
		return true, nil
	})

	if err != nil {
		return "", err
	}

	if id != "" {
		return id, nil
	}

	return "", &utils.ConfigurationError{Message: fmt.Sprintf("Cannot find line %v", line)}
}
