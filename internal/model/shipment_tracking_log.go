package model

import "time"

type ShipmentHistoryItem struct {
	Note      string `json:"note"`
	Message   string `json:"message"`
	UpdatedAt string `json:"updated_at"`
	Status    string `json:"status"`
}

type ShipmentAddressInfo struct {
	ContactName string `json:"contact_name"`
	Address     string `json:"address"`
}

type ShipmentTrackingResp struct {
	Response
	Data ShipmentTrackingResponse `json:"data"`
}

type ShipmentTrackingResponse struct {
	CourierCode     string                `json:"courier_code"`
	Waybill         string                `json:"waybill"`
	TrackingURL     string                `json:"tracking_url"`
	OriginInfo      ShipmentAddressInfo   `json:"origin_info"`
	DestinationInfo ShipmentAddressInfo   `json:"destination_info"`
	History         []ShipmentHistoryItem `json:"history"`
	LastTrackedAt   time.Time             `json:"last_tracked_at"`
}
