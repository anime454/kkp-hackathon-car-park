package service

import (
	"context"

	"be/internal/core/domain"
	portout "be/internal/core/port/out"
	"github.com/google/uuid"
)

type CarParkService struct {
	repo portout.CarParkRepository
}

func NewCarParkService(repo portout.CarParkRepository) *CarParkService {
	return &CarParkService{repo: repo}
}

func (s *CarParkService) Create(ctx context.Context, name string) (*domain.CarPark, error) {
	cp := &domain.CarPark{
		ID:   uuid.NewString(),
		Name: name,
	}
	if err := s.repo.Create(ctx, cp); err != nil {
		return nil, err
	}
	return cp, nil
}

func (s *CarParkService) Get(ctx context.Context, id string) (*domain.CarPark, error) {
	return s.repo.FindByID(ctx, id)
}
