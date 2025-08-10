package converter

import (
	"encoding/json"
	"fmt"
	"shipping-gateway/internal/entity"
	"shipping-gateway/internal/model"
	"time"
)

func ShipmentTrackingLogResponseToEntity(s *model.ShipmentTrackingResponse) (*entity.ShipmentTrackingLog, error) {
	var log entity.ShipmentTrackingLog
	log.CourierCode = s.CourierCode
	log.Waybill = s.Waybill
	log.TrackingURL = s.TrackingURL

	// Marshal OriginInfo and DestinationInfo to JSON strings
	originInfoJSON, err := json.Marshal(s.OriginInfo)
	if err != nil {
		return nil, fmt.Errorf("error marshalling origin info: %w", err)
	}
	log.OriginInfo = string(originInfoJSON)

	destinationInfoJSON, err := json.Marshal(s.DestinationInfo)
	if err != nil {
		return nil, fmt.Errorf("error marshalling destination info: %w", err)
	}
	log.DestinationInfo = string(destinationInfoJSON)

	// Marshal History to JSON string
	historyJSON, err := json.Marshal(s.History)
	if err != nil {
		return nil, fmt.Errorf("error marshalling shipment history: %w", err)
	}
	log.ShipmentHistory = string(historyJSON)

	log.LastTrackedAt = time.Now()
	return &log, nil
}

func ShipmentTrackingResponseFromJSONString(jsonStr string) (*model.ShipmentTrackingResponse, error) {
	if jsonStr == "" {
		return nil, fmt.Errorf("err parse response: json string cannot be empty")
	}
	var resp model.ShipmentTrackingResponse
	if err := json.Unmarshal([]byte(jsonStr), &resp); err != nil {
		return nil, fmt.Errorf("error unmarshalling json string: %w", err)
	}
	return &resp, nil
}

func ShipmentTrackingLogToResponse(s *entity.ShipmentTrackingLog) (*model.ShipmentTrackingResponse, error) {
	var resp model.ShipmentTrackingResponse
	resp.CourierCode = s.CourierCode
	resp.Waybill = s.Waybill
	resp.TrackingURL = s.TrackingURL

	// Unmarshal OriginInfo and DestinationInfo from JSON strings
	if err := json.Unmarshal([]byte(s.OriginInfo), &resp.OriginInfo); err != nil {
		return nil, fmt.Errorf("error unmarshalling origin info: %w", err)
	}
	if err := json.Unmarshal([]byte(s.DestinationInfo), &resp.DestinationInfo); err != nil {
		return nil, fmt.Errorf("error unmarshalling destination info: %w", err)
	}

	// Unmarshal ShipmentHistory from JSON string
	if err := json.Unmarshal([]byte(s.ShipmentHistory), &resp.History); err != nil {
		return nil, fmt.Errorf("error unmarshalling shipment history: %w", err)
	}

	return &resp, nil
}
