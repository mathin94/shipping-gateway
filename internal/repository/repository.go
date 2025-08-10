package repository

import "gorm.io/gorm"

type Repository[T any] struct {
	DB *gorm.DB
}

func (r *Repository[T]) Create(db *gorm.DB, entity *T) error {
	return db.Create(entity).Error
}

func (r *Repository[T]) Update(db *gorm.DB, entity *T) error {
	return db.Save(entity).Error
}

func (r *Repository[T]) Delete(db *gorm.DB, entity *T) error {
	return db.Delete(entity).Error
}

func (r *Repository[T]) FindByID(db *gorm.DB, id any) (*T, error) {
	var entity T
	err := db.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
