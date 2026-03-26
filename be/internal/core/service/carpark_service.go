package service

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"be/internal/config"
	"be/internal/core/domain"
	portin "be/internal/core/port/in"
	portout "be/internal/core/port/out"
	"github.com/google/uuid"
)

type CarParkService struct {
	repo       portout.CarParkRepository
	slotRepo   portout.SlotRepository
	adminUser  string
	adminPass  string
	tokens     map[string]time.Time
	tokenMutex sync.RWMutex
}

var (
	errInvalidSlotStatus = errors.New("invalid slot status")
	errInvalidToken      = errors.New("invalid token")
	errUnauthorized      = errors.New("unauthorized")
	errSlotNeverParked   = errors.New("slot has never been parked")
)

func NewCarParkService(repo portout.CarParkRepository, slotRepo portout.SlotRepository, adminCfg config.AdminConfig) *CarParkService {
	return &CarParkService{
		repo:      repo,
		slotRepo:  slotRepo,
		adminUser: adminCfg.Username,
		adminPass: adminCfg.Password,
		tokens:    make(map[string]time.Time),
	}
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

func (s *CarParkService) EnsureSeedData(ctx context.Context) error {
	carParks, err := s.repo.List(ctx)
	if err != nil {
		return err
	}

	carParkID := ""
	if len(carParks) == 0 {
		cp := &domain.CarPark{
			ID:   uuid.NewString(),
			Name: "Main Car Park",
		}
		if err := s.repo.Create(ctx, cp); err != nil {
			return err
		}
		carParkID = cp.ID
	} else {
		carParkID = carParks[0].ID
	}

	slots, err := s.slotRepo.List(ctx)
	if err != nil {
		return err
	}
	if len(slots) > 0 {
		return nil
	}

	for i := 1; i <= 20; i++ {
		slotType := domain.SlotTypeNormal
		if i <= 4 {
			slotType = domain.SlotTypeVIP
		}
		slot := &domain.Slot{
			ID:        uuid.NewString(),
			Label:     "A" + strconv.Itoa(i),
			Type:      slotType,
			Status:    domain.SlotStatusFree,
			CarParkID: carParkID,
		}
		if err := s.slotRepo.Create(ctx, slot); err != nil {
			return err
		}
	}

	return nil
}

func (s *CarParkService) ListSlots(ctx context.Context) ([]domain.Slot, error) {
	return s.slotRepo.List(ctx)
}

func (s *CarParkService) GetDashboard(ctx context.Context) (*portin.Dashboard, error) {
	slots, err := s.slotRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	dashboard := &portin.Dashboard{
		Slots:       slots,
		Total:       len(slots),
		UpdatedAt:   time.Now().UTC(),
		PriceNormal: 10,
		PriceVIP:    20,
	}

	for _, slot := range slots {
		switch slot.Status {
		case domain.SlotStatusFree:
			dashboard.Free++
		case domain.SlotStatusParked:
			dashboard.Parked++
		case domain.SlotStatusClose:
			dashboard.Closed++
		}
	}

	return dashboard, nil
}

func isValidStatus(status string) bool {
	return status == domain.SlotStatusFree || status == domain.SlotStatusParked || status == domain.SlotStatusClose
}

func (s *CarParkService) UpdateSlotStatus(ctx context.Context, slotID, status string) (*portin.SlotUpdateResult, error) {
	nextStatus := strings.TrimSpace(strings.ToLower(status))
	if !isValidStatus(nextStatus) {
		return nil, errInvalidSlotStatus
	}

	slot, err := s.slotRepo.FindByID(ctx, slotID)
	if err != nil {
		return nil, err
	}

	currentStatus := slot.Status
	var parkingFee *portin.ParkingFee

	if nextStatus == domain.SlotStatusParked && currentStatus != domain.SlotStatusParked {
		now := time.Now().UTC()
		slot.ParkedAt = &now
	}

	if nextStatus == domain.SlotStatusFree && currentStatus == domain.SlotStatusParked {
		fee, err := s.CalculateFee(ctx, slot.ID, time.Now().UTC())
		if err == nil {
			parkingFee = fee
		}
		slot.ParkedAt = nil
	}

	if nextStatus == domain.SlotStatusClose {
		slot.ParkedAt = nil
	}

	slot.Status = nextStatus
	if err := s.slotRepo.Update(ctx, slot); err != nil {
		return nil, err
	}

	return &portin.SlotUpdateResult{Slot: *slot, Fee: parkingFee}, nil
}

func parkingRates(slotType string) (base int, additional int) {
	if slotType == domain.SlotTypeVIP {
		return 20, 10
	}
	return 10, 5
}

func (s *CarParkService) CalculateFee(ctx context.Context, slotID string, exitAt time.Time) (*portin.ParkingFee, error) {
	slot, err := s.slotRepo.FindByID(ctx, slotID)
	if err != nil {
		return nil, err
	}
	if slot.ParkedAt == nil {
		return nil, errSlotNeverParked
	}

	duration := exitAt.Sub(*slot.ParkedAt)
	if duration < 0 {
		duration = 0
	}
	hours := int(math.Ceil(duration.Hours()))
	if hours < 1 {
		hours = 1
	}

	base, additional := parkingRates(slot.Type)
	amount := base
	if hours > 1 {
		amount += (hours - 1) * additional
	}

	return &portin.ParkingFee{
		SlotID:   slot.ID,
		SlotType: slot.Type,
		Hours:    hours,
		Amount:   amount,
	}, nil
}

func (s *CarParkService) AdminLogin(_ context.Context, username, password string) (string, error) {
	if username != s.adminUser || password != s.adminPass {
		return "", errUnauthorized
	}
	token := uuid.NewString()
	s.tokenMutex.Lock()
	s.tokens[token] = time.Now().UTC().Add(8 * time.Hour)
	s.tokenMutex.Unlock()
	return token, nil
}

func (s *CarParkService) Authorize(_ context.Context, token string) error {
	t := strings.TrimSpace(token)
	if t == "" {
		return errInvalidToken
	}

	s.tokenMutex.RLock()
	expiresAt, ok := s.tokens[t]
	s.tokenMutex.RUnlock()
	if !ok || time.Now().UTC().After(expiresAt) {
		return errInvalidToken
	}
	return nil
}
