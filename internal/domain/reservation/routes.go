package reservation

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

	// authenticated routes
	api := e.Group("/api/v1/reservations", middlewares.AuthMiddleware(jwtService))
	api.POST("", handler.CreateReservation)
	api.GET("/my-reservations", handler.GetMyReservations)
	api.DELETE("/:id", handler.CancelReservation)

	// admin only
	adminApi := e.Group("/api/v1/reservations", middlewares.AuthMiddleware(jwtService), middlewares.AdminOnly)
	adminApi.GET("", handler.GetAllReservations)
}
