package reservation

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateReservation(c echo.Context) error {
	userID, _ := c.Get("user_id").(uint)

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid request payload", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Validation failed", err.Error()))
	}

	response, err := h.service.CreateReservation(userID, req)
	if err != nil {
		if errors.Is(err, ErrZoneFull) {
			return c.JSON(http.StatusConflict, httpresponse.NewError("Parking zone is full", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to create reservation", err.Error()))
	}

	return c.JSON(http.StatusCreated, httpresponse.NewSuccess("Reservation confirmed successfully", response))
}

func (h *handler) GetMyReservations(c echo.Context) error {
	userID, _ := c.Get("user_id").(uint)

	reservations, err := h.service.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to get reservations", err.Error()))
	}

	return c.JSON(http.StatusOK, httpresponse.NewSuccess("My reservations retrieved successfully", reservations))
}

func (h *handler) GetAllReservations(c echo.Context) error {
	reservations, err := h.service.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to get reservations", err.Error()))
	}

	return c.JSON(http.StatusOK, httpresponse.NewSuccess("All reservations retrieved successfully", reservations))
}

func (h *handler) CancelReservation(c echo.Context) error {
	userID, _ := c.Get("user_id").(uint)
	role, _ := c.Get("user_role").(string)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid reservation ID", err.Error()))
	}

	if err := h.service.CancelReservation(uint(id), userID, role); err != nil {
		if errors.Is(err, ErrNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.NewError("Reservation not found", nil))
		}
		if errors.Is(err, ErrNotOwner) {
			return c.JSON(http.StatusForbidden, httpresponse.NewError("You cannot cancel this reservation", nil))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to cancel reservation", err.Error()))
	}

	return c.JSON(http.StatusOK, httpresponse.NewSuccess("Reservation cancelled successfully", nil))
}
