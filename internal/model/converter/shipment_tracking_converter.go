package converter

import (
	"encoding/json"
	"fmt"
	"shipping-gateway/external/biteship"
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

	// get last status from history
	if len(s.History) > 0 {
		log.TrackingStatus = s.History[len(s.History)-1].Status
	}
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
	resp.LastTrackedAt = s.LastTrackedAt
	
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

func BiteshipTrackingToResponse(t *biteship.TrackingResponse) (*model.ShipmentTrackingResponse, error) {
	if t == nil {
		return nil, fmt.Errorf("biteship tracking response is nil")
	}

	var resp model.ShipmentTrackingResponse
	resp.CourierCode = t.Courier.Company
	resp.Waybill = t.WaybillID
	resp.TrackingURL = t.Link
	resp.OriginInfo = model.ShipmentAddressInfo{
		ContactName: t.Origin.ContactName,
		Address:     t.Origin.Address,
	}
	resp.DestinationInfo = model.ShipmentAddressInfo{
		ContactName: t.Destination.ContactName,
		Address:     t.Destination.Address,
	}

	resp.History = make([]model.ShipmentHistoryItem, 0, len(t.History))
	for _, item := range t.History {
		resp.History = append(resp.History, model.ShipmentHistoryItem{
			Note:      item.Note,
			UpdatedAt: item.UpdatedAt,
			Status:    item.Status.ToString(),
			Message:   item.Status.GetMessage(),
		})
	}

	// 2021-03-16T18:17:00+07:00
	format := "2006-01-02T15:04:05-07:00"
	lastTrackedAt := t.History[len(t.History)-1].UpdatedAt // Assuming the last item is the most recent
	if lastTrackedAt != "" {
		parsedTime, err := time.Parse(format, lastTrackedAt)
		if err == nil {
			resp.LastTrackedAt = parsedTime
		} else {
			resp.LastTrackedAt = time.Now() // Fallback to current time if parsing fails
		}
	} else {
		resp.LastTrackedAt = time.Now() // Fallback to current time if no history
	}

	return &resp, nil
}
