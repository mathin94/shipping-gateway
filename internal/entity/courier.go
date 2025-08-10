package entity

import (
	"gorm.io/gorm"
	"time"
)

type Courier struct {
	ID        uint           `gorm:"primaryKey"`
	Code      string         `gorm:"type:varchar(50);not null;unique"` // Unique code for the courier
	Name      string         `gorm:"type:varchar(100);not null"`       // Name of the courier
	CreatedAt time.Time      `gorm:"autoCreateTime"`                   // Timestamp when the record was created
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`                   // Timestamp
	DeletedAt gorm.DeletedAt `gorm:"index"`                            // Soft delete field
}

// TableName returns the name of the table in the database
func (Courier) TableName() string {
	return "couriers"
}
