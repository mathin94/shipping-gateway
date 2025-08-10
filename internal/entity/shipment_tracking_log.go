package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ShipmentTrackingLog struct {
	ID              string `gorm:"primaryKey;type:varchar(36)"`
	CourierCode     string `gorm:"type:varchar(50);not null"`
	Waybill         string `gorm:"type:varchar(255);not null"`
	OriginInfo      string `gorm:"type:text;not null"`
	DestinationInfo string `gorm:"type:text;not null"`
	TrackingURL     string `gorm:"type:varchar(255);"`
	TrackingStatus  string `gorm:"type:varchar(50);not null"`
	ShipmentHistory string `gorm:"type:text;"`
	LastTrackedAt   time.Time
	CreatedAt       time.Time `gorm:"autoCreateTime"` // Timestamp when the record was created
	UpdatedAt       time.Time `gorm:"autoUpdateTime"` // Timestamp
}

// TableName returns the name of the table in the database
func (s *ShipmentTrackingLog) TableName() string {
	return "shipment_tracking_logs"
}

// BeforeCreate is a GORM hook that sets the ID and timestamps before creating a new record
func (s *ShipmentTrackingLog) BeforeCreate(tx *gorm.DB) (err error) {
	id, _ := uuid.NewV7()
	s.ID = id.String()
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return nil
}
