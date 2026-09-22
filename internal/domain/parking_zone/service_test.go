package parkingzone

import (
	"spotsync/internal/domain/parking_zone/dto"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
)

type mockParkingZoneRepo struct {
	zone                 *ParkingZone
	activeCount          int64
	updateCalled         bool
	updatedTotalCapacity int
}

func (m *mockParkingZoneRepo) CreateZone(zone *ParkingZone) error {
	return nil
}

func (m *mockParkingZoneRepo) GetAllZones() ([]dto.ZoneResponse, error) {
	return nil, nil
}

func (m *mockParkingZoneRepo) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	if m.zone == nil {
		return nil, ErrZoneNotFound
	}
	avail := m.zone.TotalCapacity - int(m.activeCount)
	return &dto.ZoneResponse{
		ID:             m.zone.ID,
		Name:           m.zone.Name,
		Type:           m.zone.Type,
		TotalCapacity:  m.zone.TotalCapacity,
		AvailableSpots: avail,
		PricePerHour:   m.zone.PricePerHour,
		CreatedAt:      time.Now().Format(time.RFC3339),
	}, nil
}

func (m *mockParkingZoneRepo) GetZoneRaw(id uint) (*ParkingZone, error) {
	if m.zone == nil {
		return nil, ErrZoneNotFound
	}
	return m.zone, nil
}

func (m *mockParkingZoneRepo) UpdateZone(zone *ParkingZone) error {
	m.updateCalled = true
	m.updatedTotalCapacity = zone.TotalCapacity
	return nil
}

func (m *mockParkingZoneRepo) DeleteZone(id uint) error {
	return nil
}

func (m *mockParkingZoneRepo) GetActiveReservationCount(zoneID uint) (int64, error) {
	return m.activeCount, nil
}

func TestUpdateZone_RejectCapacityBelowActiveReservations(t *testing.T) {
	zone := &ParkingZone{
		Name:          "North Zone",
		Type:          "Standard",
		TotalCapacity: 10,
		PricePerHour:  25.0,
	}
	zone.ID = 1

	mockRepo := &mockParkingZoneRepo{
		zone:        zone,
		activeCount: 5,
	}

	svc := NewService(mockRepo)

	// Attempt to reduce capacity to 3 (below 5 active reservations)
	req := dto.UpdateZoneRequest{
		TotalCapacity: 3,
	}

	_, err := svc.UpdateZone(1, req)
	if err == nil {
		t.Fatalf("expected error when reducing capacity below active reservations, got nil")
	}

	expectedMsg := "cannot reduce capacity below 5 active reservations"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("expected error message to contain '%s', got '%v'", expectedMsg, err.Error())
	}

	if mockRepo.updateCalled {
		t.Errorf("repo.UpdateZone should not have been called when validation fails")
	}
}

func TestUpdateZone_AllowCapacityEqualToActiveReservations(t *testing.T) {
	zone := &ParkingZone{
		Name:          "North Zone",
		Type:          "Standard",
		TotalCapacity: 10,
		PricePerHour:  25.0,
	}
	zone.ID = 1

	mockRepo := &mockParkingZoneRepo{
		zone:        zone,
		activeCount: 5,
	}

	svc := NewService(mockRepo)

	// Reduce capacity to exactly 5 (equal to active reservations)
	req := dto.UpdateZoneRequest{
		TotalCapacity: 5,
	}

	res, err := svc.UpdateZone(1, req)
	if err != nil {
		t.Fatalf("unexpected error when reducing capacity to active reservations: %v", err)
	}

	if !mockRepo.updateCalled {
		t.Errorf("expected repo.UpdateZone to be called")
	}

	if res.TotalCapacity != 5 {
		t.Errorf("expected TotalCapacity to be 5, got %d", res.TotalCapacity)
	}

	if res.AvailableSpots != 0 {
		t.Errorf("expected AvailableSpots to be 0, got %d", res.AvailableSpots)
	}
}

func TestUpdateZone_AllowCapacityAboveActiveReservations(t *testing.T) {
	zone := &ParkingZone{
		Model:         gorm.Model{ID: 1},
		Name:          "North Zone",
		Type:          "Standard",
		TotalCapacity: 10,
		PricePerHour:  25.0,
	}

	mockRepo := &mockParkingZoneRepo{
		zone:        zone,
		activeCount: 5,
	}

	svc := NewService(mockRepo)

	// Update capacity to 8 (above 5 active reservations)
	req := dto.UpdateZoneRequest{
		TotalCapacity: 8,
	}

	res, err := svc.UpdateZone(1, req)
	if err != nil {
		t.Fatalf("unexpected error when setting capacity above active reservations: %v", err)
	}

	if !mockRepo.updateCalled {
		t.Errorf("expected repo.UpdateZone to be called")
	}

	if res.TotalCapacity != 8 {
		t.Errorf("expected TotalCapacity to be 8, got %d", res.TotalCapacity)
	}

	if res.AvailableSpots != 3 {
		t.Errorf("expected AvailableSpots to be 3, got %d", res.AvailableSpots)
	}
}
