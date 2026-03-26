package out

import (
	"context"

	"be/internal/core/domain"
)

type CarParkRepository interface {
	Create(ctx context.Context, carPark *domain.CarPark) error
	FindByID(ctx context.Context, id string) (*domain.CarPark, error)
	List(ctx context.Context) ([]domain.CarPark, error)
}

type SlotRepository interface {
	Create(ctx context.Context, slot *domain.Slot) error
	List(ctx context.Context) ([]domain.Slot, error)
	FindByID(ctx context.Context, id string) (*domain.Slot, error)
	Update(ctx context.Context, slot *domain.Slot) error
}
