package user

import (
	"spotsync/internal/auth"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	repo := NewRepository(db)
	service := NewService(repo, jwtService)
	handler := NewHandler(service)

	api := e.Group("/api/v1/auth")
	api.POST("/register", handler.Register)
	api.POST("/login", handler.Login)
}
