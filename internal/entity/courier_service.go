package entity

import "time"

type CourierService struct {
	ID          uint      `gorm:"primaryKey"`                      // Unique identifier for the courier service
	CourierID   uint      `gorm:"not null"`                        // ID of the associated courier
	CourierCode string    `gorm:"type:varchar(50);not null"`       // Code of the associated courier
	CourierName string    `gorm:"type:varchar(100);not null"`      // Name of the associated courier
	ServiceCode string    `gorm:"type:varchar(50);not null"`       // Code of the courier service
	ServiceName string    `gorm:"type:varchar(100);not null"`      // Name of the courier service
	Tier        string    `gorm:"type:varchar(50);not null"`       // Tier of the courier service (e.g., standard, express)
	Description string    `gorm:"type:varchar(200)"`               // Description of the courier service
	ServiceType string    `gorm:"type:varchar(50);not null;index"` // Type of the courier service (e.g., delivery, pickup)
	ETD         string    `gorm:"type:varchar(50);not null"`       // Estimated Time of Delivery for the service
	CreatedAt   time.Time `gorm:"autoCreateTime"`                  // Timestamp when the record was created
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`                  // Timestamp
}

// TableName returns the name of the table in the database
func (CourierService) TableName() string {
	return "courier_services"
}
