package user

import (
	"errors"
	"net/http"
	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid request payload", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Validation failed", err.Error()))
	}

	response, err := h.service.Register(req)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return c.JSON(http.StatusBadRequest, httpresponse.NewError("Email already registered", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to register user", err.Error()))
	}

	return c.JSON(http.StatusCreated, httpresponse.NewSuccess("User registered successfully", response))
}

func (h *handler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid request payload", err.Error()))
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Validation failed", err.Error()))
	}

	response, err := h.service.Login(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.NewError("Invalid email or password", err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to login", err.Error()))
	}

	return c.JSON(http.StatusOK, httpresponse.NewSuccess("Login successful", response))
}
