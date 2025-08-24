package entity

import (
	"gorm.io/gorm"
	"time"
)

type Area struct {
	ID                    uint      `json:"id" gorm:"primaryKey"`
	OriginalSubdistrictID uint      `json:"original_subdistrict_id" gorm:"not null"`               // ID from the original subdistrict
	OriginalPostalCode    string    `json:"original_postal_code" gorm:"type:varchar(20);not null"` // Postal code from the original subdistrict
	Description           string    `json:"description" gorm:"not null"`
	ExternalSource        string    `json:"external_source" gorm:"type:varchar(100);not null"` // Name of the external service providing the area data
	ExternalID            string    `json:"external_id" gorm:"type:varchar(200);unique"`       // ID from the external service
	ExternalInfo          string    `json:"external_info" gorm:"type:text"`                    // JSON string containing additional info from the external service
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime"`                  // Timestamp when the record was created
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`                  // Timestamp
}

// TableName overrides the table name used by Area to `areas`
func (area *Area) TableName() string {
	return "areas"
}

// BeforeCreate is a GORM hook that sets the CreatedAt timestamp
func (area *Area) BeforeCreate(tx *gorm.DB) error {
	area.CreatedAt = time.Now()
	area.UpdatedAt = time.Now()
	return nil
}
