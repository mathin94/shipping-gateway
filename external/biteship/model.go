package biteship

import (
	"encoding/json"
	"shipping-gateway/internal/model"
	"strconv"
	"time"
)

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int    `json:"value"`
	Weight      int    `json:"weight"`
	Length      int    `json:"length"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Quantity    int    `json:"quantity"`
}

type RateRequest struct {
	OriginLatitude        float64 `json:"origin_latitude,omitempty"`
	OriginLongitude       float64 `json:"origin_longitude,omitempty"`
	OriginPostalCode      int     `json:"origin_postal_code,omitempty"`
	OriginAreaID          string  `json:"origin_area_id,omitempty"`
	DestinationPostalCode int     `json:"destination_postal_code,omitempty"`
	DestinationAreaID     string  `json:"destination_area_id,omitempty"`
	DestinationLatitude   float64 `json:"destination_latitude,omitempty"`
	DestinationLongitude  float64 `json:"destination_longitude,omitempty"`
	Couriers              string  `json:"couriers"`
	Items                 []Item  `json:"items"`
}

type Area struct {
	ID                               string  `json:"id"`
	Name                             string  `json:"name,omitempty"`
	Latitude                         float64 `json:"latitude,omitempty"`
	Longitude                        float64 `json:"longitude,omitempty"`
	PostalCode                       int     `json:"postal_code"`
	CountryName                      string  `json:"country_name"`
	CountryCode                      string  `json:"country_code"`
	AdministrativeDivisionLevel1Name string  `json:"administrative_division_level_1_name"`
	AdministrativeDivisionLevel1Type string  `json:"administrative_division_level_1_type"`
	AdministrativeDivisionLevel2Name string  `json:"administrative_division_level_2_name"`
	AdministrativeDivisionLevel2Type string  `json:"administrative_division_level_2_type"`
	AdministrativeDivisionLevel3Name string  `json:"administrative_division_level_3_name"`
	AdministrativeDivisionLevel3Type string  `json:"administrative_division_level_3_type"`
	Address                          string  `json:"address,omitempty"`
}

type AreaPricingWrapper struct {
	LocationID                       string  `json:"location_id"`
	Latitude                         float64 `json:"latitude,omitempty"`
	Longitude                        float64 `json:"longitude,omitempty"`
	PostalCode                       int     `json:"postal_code"`
	CountryName                      string  `json:"country_name"`
	CountryCode                      string  `json:"country_code"`
	AdministrativeDivisionLevel1Name string  `json:"administrative_division_level_1_name"`
	AdministrativeDivisionLevel1Type string  `json:"administrative_division_level_1_type"`
	AdministrativeDivisionLevel2Name string  `json:"administrative_division_level_2_name"`
	AdministrativeDivisionLevel2Type string  `json:"administrative_division_level_2_type"`
	AdministrativeDivisionLevel3Name string  `json:"administrative_division_level_3_name"`
	AdministrativeDivisionLevel3Type string  `json:"administrative_division_level_3_type"`
	AdministrativeDivisionLevel4Name string  `json:"administrative_division_level_4_name,omitempty"`
	AdministrativeDivisionLevel4Type string  `json:"administrative_division_level_4_type,omitempty"`
	Address                          string  `json:"address,omitempty"`
}

func (area AreaPricingWrapper) ToLocationResponse() model.LocationResponse {
	return model.LocationResponse{
		Country:     area.CountryName,
		CountryCode: area.CountryCode,
		Province:    area.AdministrativeDivisionLevel1Name,
		City:        area.AdministrativeDivisionLevel2Name,
		District:    area.AdministrativeDivisionLevel3Name,
		Address:     area.Address,
		PostalCode:  strconv.Itoa(area.PostalCode),
	}
}

type AreaResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Object  string `json:"object"`
	Areas   []Area `json:"areas"`
}

type Pricing struct {
	AvailableCollectionMethod    []string      `json:"available_collection_method"`
	AvailableForCashOnDelivery   bool          `json:"available_for_cash_on_delivery"`
	AvailableForProofOfDelivery  bool          `json:"available_for_proof_of_delivery"`
	AvailableForInstantWaybillID bool          `json:"available_for_instant_waybill_id"`
	AvailableForInsurance        bool          `json:"available_for_insurance"`
	Company                      string        `json:"company"`
	CourierName                  string        `json:"courier_name"`
	CourierCode                  string        `json:"courier_code"`
	CourierServiceName           string        `json:"courier_service_name"`
	CourierServiceCode           string        `json:"courier_service_code"`
	Currency                     string        `json:"currency"`
	Description                  string        `json:"description"`
	Duration                     string        `json:"duration"`
	ShipmentDurationRange        string        `json:"shipment_duration_range"`
	ShipmentDurationUnit         string        `json:"shipment_duration_unit"`
	Price                        int           `json:"price"`
	TaxLines                     []interface{} `json:"tax_lines,omitempty"`
	Type                         string        `json:"type"`
}

type RateResponse struct {
	Success     bool               `json:"success"`
	Message     string             `json:"message"`
	Object      string             `json:"object"`
	Stops       []interface{}      `json:"stops"`
	Origin      AreaPricingWrapper `json:"origin"`
	Destination AreaPricingWrapper `json:"destination"`
	Pricing     []Pricing          `json:"pricing"`
}

// ToCourierRateResponse func
func (r RateResponse) ToCourierRateResponse() model.CourierRateResponse {
	var resp model.CourierRateResponse

	resp.Response.Status = "success"
	resp.Response.Code = 20000
	resp.Response.Message = r.Message
	resp.Data.Origin = r.Origin.ToLocationResponse()
	resp.Data.Destination = r.Destination.ToLocationResponse()

	prices := make([]model.CourierPrice, 0)
	for _, pricing := range r.Pricing {
		price := model.CourierPrice{
			CourierCode: pricing.CourierCode,
			CourierName: pricing.CourierName,
			ServiceType: pricing.Type,
			ServiceName: pricing.CourierServiceName,
			ServiceCode: pricing.CourierServiceCode,
			Price:       pricing.Price,
			ETD:         pricing.Duration,
		}
		prices = append(prices, price)
	}

	resp.Data.Prices = prices
	return resp
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
}

func (e ErrorResponse) GetHTTPStatusCode() int {
	// Get Substring first 3 digit 
	code := strconv.Itoa(e.Code)
	if len(code) < 3 {
		return 500 // Default to Internal Server Error if code is not valid
	}

	iCode, _ := strconv.Atoi(code[:3])
	return iCode
}

func (e ErrorResponse) IsEmptyData() bool {
	return e.Code == ErrRateNoCourierAvailable
}

func NewErrorResponse(code int, err string) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Code:    code,
		Error:   err,
	}
}

func ErrorResponseFromBytes(data []byte) *ErrorResponse {
	var errResp ErrorResponse
	if err := json.Unmarshal(data, &errResp); err != nil {
		return NewErrorResponse(ErrInvalidParsingResponse, "Failed to parse error response")
	}

	return &errResp
}

// CourierTracking struct
type CourierTracking struct {
	Company           string `json:"company"`
	DriverName        string `json:"driver_name"`
	DriverPhone       string `json:"driver_phone"`
	DriverPhotoURL    string `json:"driver_photo_url"`
	DriverPlateNumber string `json:"driver_plate_number"`
}

// TrackingAdddress struct
type TrackingAdddress struct {
	ContactName string `json:"contact_name"`
	Address     string `json:"address"`
}

// TrackingHistoryItem struct
type TrackingHistoryItem struct {
	Note      string `json:"note"`
	UpdatedAt string `json:"updated_at"`
	Status    string `json:"status"`
}

// TrackingResponse struct
type TrackingResponse struct {
	Success     bool                  `json:"success"`
	Message     string                `json:"message"`
	Object      string                `json:"object"`
	ID          string                `json:"id"`
	WaybillID   string                `json:"waybill_id"`
	Courier     CourierTracking       `json:"courier"`
	Origin      TrackingAdddress      `json:"origin"`
	Destination TrackingAdddress      `json:"destination"`
	History     []TrackingHistoryItem `json:"history"`
	Link        string                `json:"link"`
	OrderID     string                `json:"order_id,omitempty"`
	Status      string                `json:"status"`
}

func (t TrackingResponse) ToShipmentTrackingResponse() model.ShipmentTrackingResponse {
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
			Status:    item.Status,
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

	return resp
}
