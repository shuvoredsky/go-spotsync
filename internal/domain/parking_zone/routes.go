package parkingzone

import (
	"spotsync/internal/auth"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// public routes
	api := e.Group("/api/v1/zones")
	api.GET("", handler.GetAllZones)
	api.GET("/:id", handler.GetZoneByID)

	// admin only routes
	adminApi := e.Group("/api/v1/zones", middlewares.AuthMiddleware(jwtService), middlewares.AdminOnly)
	adminApi.POST("", handler.CreateZone)
	adminApi.PUT("/:id", handler.UpdateZone)
	adminApi.DELETE("/:id", handler.DeleteZone)
}
