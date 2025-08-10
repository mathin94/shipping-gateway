package model

type Metrics struct {
	Length int `json:"length"` // Length of the item in centimeters
	Width  int `json:"width"`  // Width of the item in centimeters
	Height int `json:"height"` // Height of the item in centimeters
	Weight int `json:"weight"` // Weight of the item in grams
}
type ItemRequest struct {
	Name        string `json:"name"`        // Name of the item
	Description string `json:"description"` // Description of the item
	Price       int    `json:"price"`       // Price of the item in cents
	Metrics
	Quantity int `json:"quantity"` // Quantity of the item
}

type CourierRateRequest struct {
	OriginSubdistrictID      string        `json:"origin_subdistrict_id"`
	DestinationSubdistrictID string        `json:"destination_subdistrict_id"`
	OriginQuery              string        `json:"origin_query"`            // Optional query for origin, e.g., postal code or subdistrict name
	DestinationQuery         string        `json:"destination_query"`       // Optional query for destination, e.g., postal code or subdistrict name
	OriginPostalCode         string        `json:"origin_postal_code"`      // Postal code of the origin area
	DestinationPostalCode    string        `json:"destination_postal_code"` // Postal code of the destination area
	CourierCode              string        `json:"courier_code"`            // Code of the courier service
	Items                    []ItemRequest `json:"items"`                   // List of items to be shipped
}

type CourierPrice struct {
	CourierCode string `json:"courier_code"` // Code of the courier service
	CourierName string `json:"courier_name"` // Name of the courier service
	ServiceType string `json:"service_type"` // Type of service provided by the courier
	ServiceName string `json:"service_name"` // Name of the service provided by the courier
	ServiceCode string `json:"service_code"` // Code of the service provided by the courier
	Price       int    `json:"price"`        // Price of the service in cents
	ETD         string `json:"etd"`          // Estimated time of delivery
}

type CourierRate struct {
	Origin      LocationResponse `json:"origin"`
	Destination LocationResponse `json:"destination"`
	Prices      []CourierPrice   `json:"prices"` // List of prices for the courier services
}

type CourierRateResponse struct {
	Response
	Data CourierRate `json:"data"` // Data containing the origin, destination, and prices for the courier services
}
