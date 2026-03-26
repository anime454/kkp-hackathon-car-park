package in

import (
	"context"
	"time"

	"be/internal/core/domain"
)

type Dashboard struct {
	Slots       []domain.Slot `json:"slots"`
	Total       int           `json:"total"`
	Free        int           `json:"free"`
	Parked      int           `json:"parked"`
	Closed      int           `json:"closed"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	PriceNormal int           `json:"priceNormal"`
	PriceVIP    int           `json:"priceVip"`
}

type ParkingFee struct {
	SlotID   string `json:"slotId"`
	SlotType string `json:"slotType"`
	Hours    int    `json:"hours"`
	Amount   int    `json:"amount"`
}

type SlotUpdateResult struct {
	Slot domain.Slot `json:"slot"`
	Fee  *ParkingFee `json:"fee,omitempty"`
}

type CarParkUseCase interface {
	Create(ctx context.Context, name string) (*domain.CarPark, error)
	Get(ctx context.Context, id string) (*domain.CarPark, error)
	GetDashboard(ctx context.Context) (*Dashboard, error)
	ListSlots(ctx context.Context) ([]domain.Slot, error)
	UpdateSlotStatus(ctx context.Context, slotID, status string) (*SlotUpdateResult, error)
	CalculateFee(ctx context.Context, slotID string, exitAt time.Time) (*ParkingFee, error)
	AdminLogin(ctx context.Context, username, password string) (string, error)
	Authorize(ctx context.Context, token string) error
}
