package parkingzone

import (
	"errors"
	"spotsync/internal/domain/parking_zone/dto"
	"time"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("parking zone not found")

type Service interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	DeleteZone(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}
	if err := s.repo.CreateZone(zone); err != nil {
		return nil, err
	}
	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity, // Brand new zone has 0 reservations
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (s *service) GetAllZones() ([]dto.ZoneResponse, error) {
	return s.repo.GetAllZones()
}

func (s *service) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetZoneByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}
	return zone, nil
}

func (s *service) UpdateZone(id uint, req dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.repo.GetZoneRaw(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrZoneNotFound
		}
		return nil, err
	}

	if req.Name != "" {
		zone.Name = req.Name
	}
	if req.Type != "" {
		zone.Type = req.Type
	}
	if req.TotalCapacity > 0 {
		zone.TotalCapacity = req.TotalCapacity
	}
	if req.PricePerHour > 0 {
		zone.PricePerHour = req.PricePerHour
	}

	if err := s.repo.UpdateZone(zone); err != nil {
		return nil, err
	}

	// Fetch refreshed zone with dynamically computed available spots
	return s.repo.GetZoneByID(id)
}


func (s *service) DeleteZone(id uint) error {
	_, err := s.repo.GetZoneRaw(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrZoneNotFound
		}
		return err
	}
	return s.repo.DeleteZone(id)
}
