package repository

import (
	"errors"
	"gorm.io/gorm"
	"shipping-gateway/internal/entity"
	"shipping-gateway/internal/model"
)

type AreaRepository struct {
	Repository[entity.Area]
}

func NewAreaRepository() *AreaRepository {
	return &AreaRepository{}
}

func (r *AreaRepository) SaveOrCreate(db *gorm.DB, area *entity.Area) error {
	var existingArea entity.Area

	// Check if the area already exists
	err := db.Model(&existingArea).Where("external_id = ?", area.ExternalID).First(&existingArea).Error
	if err == nil {
		// Area exists, update it
		area.ID = existingArea.ID // Preserve the ID for update
		return db.Save(area).Error
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// An error occurred other than "not found"
		return err
	}

	// Area does not exist, create a new one
	return db.Create(area).Error
}

func (r *AreaRepository) FindByParameter(db *gorm.DB, params model.FindAreaParams) (*entity.Area, error) {
	var area entity.Area

	query := db.Model(&area)
	if params.OriginalSubdistrictID != "" {
		query = query.Where("original_subdistrict_id = ?", params.OriginalSubdistrictID)
	}

	if params.OriginalPostalCode != "" {
		query = query.Where("original_postal_code = ?", params.OriginalPostalCode)
	}

	if params.ExternalSource != "" {
		query = query.Where("external_source = ?", params.ExternalSource)
	}

	if params.ExternalID != "" {
		query = query.Where("external_id = ?", params.ExternalID)
	}

	err := query.First(&area).Error
	return &area, err
}
