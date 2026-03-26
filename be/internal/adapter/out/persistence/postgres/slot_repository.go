package postgres

import (
	"context"

	"be/internal/core/domain"
	"gorm.io/gorm"
)

type SlotRepository struct {
	db *gorm.DB
}

func NewSlotRepository(db *gorm.DB) *SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) Create(ctx context.Context, slot *domain.Slot) error {
	return r.db.WithContext(ctx).Create(slot).Error
}

func (r *SlotRepository) List(ctx context.Context) ([]domain.Slot, error) {
	var slots []domain.Slot
	if err := r.db.WithContext(ctx).Order("label asc").Find(&slots).Error; err != nil {
		return nil, err
	}
	return slots, nil
}

func (r *SlotRepository) FindByID(ctx context.Context, id string) (*domain.Slot, error) {
	var slot domain.Slot
	if err := r.db.WithContext(ctx).First(&slot, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *SlotRepository) Update(ctx context.Context, slot *domain.Slot) error {
	return r.db.WithContext(ctx).Save(slot).Error
}
