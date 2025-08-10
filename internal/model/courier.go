package model

type CourierResponse struct {
	ID   uint   `json:"id"`
	Code string `json:"code"` // Unique code for the courier
	Name string `json:"name"` // Name of the courier
}

// SearchCourierRequest represents the request structure for searching couriers
type SearchCourierRequest struct {
	Name string `json:"name" form:"name"` // Name of the courier to search for
}

// SearchCourierResponse represents the response structure for searching couriers
type SearchCourierResponse struct {
	Response
	Data []CourierResponse  `json:"data"` // List of couriers matching the search criteria
	Meta PaginationResponse `json:"meta"` // Pagination metadata
}
