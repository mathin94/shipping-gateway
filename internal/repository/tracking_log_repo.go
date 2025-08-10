package repository

import (
	"errors"
	"gorm.io/gorm"
	"shipping-gateway/internal/entity"
)

type TrackingLogRepository struct {
	Repository[entity.ShipmentTrackingLog]
}

func NewTrackingLogRepository() *TrackingLogRepository {
	return &TrackingLogRepository{}
}

func (r *TrackingLogRepository) FindByWaybillAndCourierCode(db *gorm.DB, waybill, courierCode string) (*entity.ShipmentTrackingLog, error) {
	var log entity.ShipmentTrackingLog
	if err := db.Where("waybill = ?", waybill).Where("courier_code = ?", courierCode).First(&log).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No record found
		}
		return nil, err // Other errors
	}
	return &log, nil
}
