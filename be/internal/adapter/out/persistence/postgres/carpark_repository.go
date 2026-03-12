package postgres

import (
	"context"

	"be/internal/core/domain"
	"gorm.io/gorm"
)

type CarParkRepository struct {
	db *gorm.DB
}

func NewCarParkRepository(db *gorm.DB) *CarParkRepository {
	return &CarParkRepository{db: db}
}

func (r *CarParkRepository) Create(ctx context.Context, carPark *domain.CarPark) error {
	return r.db.WithContext(ctx).Create(carPark).Error
}

func (r *CarParkRepository) FindByID(ctx context.Context, id string) (*domain.CarPark, error) {
	var cp domain.CarPark
	if err := r.db.WithContext(ctx).First(&cp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cp, nil
}
