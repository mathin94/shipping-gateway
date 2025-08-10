package model

type FindAreaParams struct {
	OriginalSubdistrictID string `json:"original_subdistrict_id"` // ID from the original subdistrict
	OriginalPostalCode    string `json:"original_postal_code"`    // Postal code from the original subdistrict
	ExternalSource        string `json:"external_source"`         // Name of
	ExternalID            string `json:"external_id"`             // ID from the external service
}

type LocationResponse struct {
	Country     string `json:"country"`      // Country code of the area
	CountryCode string `json:"country_code"` // Country code of the area
	Province    string `json:"province"`     // Province of the area
	City        string `json:"city"`         // City of the area
	District    string `json:"district"`     // District of the area
	Address     string `json:"address"`      // Address of the area
	PostalCode  string `json:"postal_code"`  // Postal code of the area
}
