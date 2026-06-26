package server

import (
	"spotsync/internal/auth"
	parkingzone "spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func registerRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	user.RegisterRoutes(e, db, jwtService)
	parkingzone.RegisterRoutes(e, db, jwtService)
	reservation.RegisterRoutes(e, db, jwtService)
}
