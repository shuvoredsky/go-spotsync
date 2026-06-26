package parkingzone

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/parking_zone/dto"
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

func (h *handler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid request payload", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Validation failed", err.Error()))
	}

	zone, err := h.service.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to create zone", err.Error()))
	}

	return c.JSON(http.StatusCreated, httpresponse.NewSuccess("Parking zone created successfully", dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt.String(),
	}))
}

func (h *handler) GetAllZones(c echo.Context) error {
	zones, err := h.service.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to get zones", err.Error()))
	}
	return c.JSON(http.StatusOK, httpresponse.NewSuccess("Parking zones retrieved successfully", zones))
}

func (h *handler) GetZoneByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid zone ID", err.Error()))
	}

	zone, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.NewError("Zone not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to get zone", err.Error()))
	}
	return c.JSON(http.StatusOK, httpresponse.NewSuccess("Parking zone retrieved successfully", zone))
}

func (h *handler) UpdateZone(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid zone ID", err.Error()))
	}

	var req dto.UpdateZoneRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid request payload", err.Error()))
	}

	zone, err := h.service.UpdateZone(uint(id), req)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.NewError("Zone not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to update zone", err.Error()))
	}

	return c.JSON(http.StatusOK, httpresponse.NewSuccess("Parking zone updated successfully", dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt.String(),
	}))
}

func (h *handler) DeleteZone(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid zone ID", err.Error()))
	}

	if err := h.service.DeleteZone(uint(id)); err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.NewError("Zone not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to delete zone", err.Error()))
	}
	return c.JSON(http.StatusOK, httpresponse.NewSuccess("Parking zone deleted successfully", nil))
}
