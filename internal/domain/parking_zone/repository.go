package parkingzone

import (
	"spotsync/internal/domain/parking_zone/dto"

	"gorm.io/gorm"
)

type Repository interface {
	CreateZone(zone *ParkingZone) error
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(zone *ParkingZone) error
	DeleteZone(id uint) error
	GetZoneRaw(id uint) (*ParkingZone, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateZone(zone *ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *repository) GetAllZones() ([]dto.ZoneResponse, error) {
	var zones []ParkingZone
	if err := r.db.Find(&zones).Error; err != nil {
		return nil, err
	}

	var responses []dto.ZoneResponse
	for _, zone := range zones {
		// count active reservations
		var activeCount int64
		r.db.Table("reservations").
			Where("zone_id = ? AND status = ? AND deleted_at IS NULL", zone.ID, "active").
			Count(&activeCount)

		responses = append(responses, dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: zone.TotalCapacity - int(activeCount),
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt.String(),
		})
	}
	return responses, nil
}

func (r *repository) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	var zone ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, err
	}

	var activeCount int64
	r.db.Table("reservations").
		Where("zone_id = ? AND status = ? AND deleted_at IS NULL", zone.ID, "active").
		Count(&activeCount)

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity - int(activeCount),
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt.String(),
	}, nil
}

func (r *repository) UpdateZone(zone *ParkingZone) error {
	return r.db.Save(zone).Error
}

func (r *repository) DeleteZone(id uint) error {
	return r.db.Delete(&ParkingZone{}, id).Error
}

func (r *repository) GetZoneRaw(id uint) (*ParkingZone, error) {
	var zone ParkingZone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}
