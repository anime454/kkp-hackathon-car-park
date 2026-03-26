package domain

import "time"

const (
	SlotStatusFree   = "free"
	SlotStatusParked = "parked"
	SlotStatusClose  = "close"
)

const (
	SlotTypeVIP    = "VIP"
	SlotTypeNormal = "normal"
)

type Slot struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	Label     string     `json:"label"`
	Type      string     `json:"type"`
	Status    string     `json:"status"`
	ParkedAt  *time.Time `json:"parkedAt"`
	CarParkID string     `json:"carParkId"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}
