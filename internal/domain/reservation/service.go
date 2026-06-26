package reservation

import (
	"spotsync/internal/domain/reservation/dto"
)

type Service interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	GetAllReservations() ([]dto.AdminReservationResponse, error)
	CancelReservation(reservationID uint, userID uint, role string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	reservation, err := s.repo.CreateReservation(userID, req)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt.String(),
		UpdatedAt:    reservation.UpdatedAt.String(),
	}, nil
}

func (s *service) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.repo.GetMyReservations(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.MyReservationResponse
	for _, r := range reservations {
		responses = append(responses, dto.MyReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			CreatedAt: r.CreatedAt.String(),
		})
	}
	return responses, nil
}

func (s *service) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.repo.GetAllReservations()
	if err != nil {
		return nil, err
	}

	var responses []dto.AdminReservationResponse
	for _, r := range reservations {
		responses = append(responses, dto.AdminReservationResponse{
			ID:           r.ID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			Zone: dto.ZoneInfo{
				ID:   r.Zone.ID,
				Name: r.Zone.Name,
				Type: r.Zone.Type,
			},
			User: dto.UserInfo{
				ID:    r.User.ID,
				Name:  r.User.Name,
				Email: r.User.Email,
			},
			CreatedAt: r.CreatedAt.String(),
		})
	}
	return responses, nil
}

func (s *service) CancelReservation(reservationID uint, userID uint, role string) error {
	return s.repo.CancelReservation(reservationID, userID, role)
}
