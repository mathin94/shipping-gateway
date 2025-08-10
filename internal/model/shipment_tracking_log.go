package model

import "time"

type ShipmentHistoryItem struct {
	Note      string `json:"note"`
	UpdatedAt string `json:"updated_at"`
	Status    string `json:"status"`
}

type ShipmentAddressInfo struct {
	ContactName string `json:"contact_name"`
	Address     string `json:"address"`
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
