package service

import (
	"context"

	portout "be/internal/core/port/out"
)

type HealthService struct {
	db portout.DBPinger
}

func NewHealthService(db portout.DBPinger) *HealthService {
	return &HealthService{db: db}
}

func (s *HealthService) Check(ctx context.Context) error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Ping(ctx)
}
