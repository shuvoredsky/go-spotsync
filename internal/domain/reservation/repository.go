package reservation

import (
	"errors"
	"spotsync/internal/domain/reservation/dto"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull = errors.New("parking zone is full")
	ErrNotFound = errors.New("reservation not found")
	ErrNotOwner = errors.New("you are not allowed to cancel this reservation")
)

type Repository interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*Reservation, error)
	GetMyReservations(userID uint) ([]Reservation, error)
	GetAllReservations() ([]Reservation, error)
	CancelReservation(reservationID uint, userID uint, role string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateReservation(userID uint, req dto.CreateReservationRequest) (*Reservation, error) {
	var newReservation *Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock the parking zone row
		var zone struct {
			ID            uint
			TotalCapacity int
		}
		if err := tx.Table("parking_zones").
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("id, total_capacity").
			Where("id = ? AND deleted_at IS NULL", req.ZoneID).
			First(&zone).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("parking zone not found")
			}
			return err
		}

		// 2. Count active reservations
		var activeCount int64
		if err := tx.Model(&Reservation{}).
			Where("zone_id = ? AND status = ? AND deleted_at IS NULL", req.ZoneID, "active").
			Count(&activeCount).Error; err != nil {
			return err
		}

		// 3. Check capacity
		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		// 4. Create reservation
		newReservation = &Reservation{
			UserID:       userID,
			ZoneID:       req.ZoneID,
			LicensePlate: req.LicensePlate,
			Status:       "active",
		}
		return tx.Create(newReservation).Error
	})

	if err != nil {
		return nil, err
	}
	return newReservation, nil
}

func (r *repository) GetMyReservations(userID uint) ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Preload("Zone").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) GetAllReservations() ([]Reservation, error) {
	var reservations []Reservation
	if err := r.db.Preload("Zone").Preload("User").
		Find(&reservations).Error; err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *repository) CancelReservation(reservationID uint, userID uint, role string) error {
	var reservation Reservation
	if err := r.db.First(&reservation, reservationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}

	// driver শুধু নিজেরটা cancel করতে পারবে
	if role != "admin" && reservation.UserID != userID {
		return ErrNotOwner
	}

	return r.db.Model(&reservation).Update("status", "cancelled").Error
}
