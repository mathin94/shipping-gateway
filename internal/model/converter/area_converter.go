package converter

import (
	"encoding/json"
	"shipping-gateway/internal/entity"
)

func AreaFromJSONString(jsonString string) (*entity.Area, error) {
	if jsonString == "" {
		return nil, nil // or return an error if you prefer
	}

	var area entity.Area
	if err := json.Unmarshal([]byte(jsonString), &area); err != nil {
		return nil, err // return the error if unmarshalling fails
	}

	return &area, nil
}
