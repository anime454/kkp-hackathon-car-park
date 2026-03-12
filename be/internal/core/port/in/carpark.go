package in

import (
	"context"

	"be/internal/core/domain"
)

type CarParkUseCase interface {
	Create(ctx context.Context, name string) (*domain.CarPark, error)
	Get(ctx context.Context, id string) (*domain.CarPark, error)
}
