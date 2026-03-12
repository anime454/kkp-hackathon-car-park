package out

import (
	"context"

	"be/internal/core/domain"
)

type CarParkRepository interface {
	Create(ctx context.Context, carPark *domain.CarPark) error
	FindByID(ctx context.Context, id string) (*domain.CarPark, error)
}
