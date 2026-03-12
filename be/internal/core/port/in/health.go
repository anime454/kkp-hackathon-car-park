package in

import "context"

type HealthUseCase interface {
	Check(ctx context.Context) error
}
