package server

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Start(db *gorm.DB, cfg *config.Config) {
	e := echo.New()

	// middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS allow origins from config (supports comma-separated list or single origin)
	var allowOrigins []string
	if cfg.AllowedOrigin != "" {
		for _, o := range strings.Split(cfg.AllowedOrigin, ",") {
			trimmed := strings.TrimSpace(o)
			if trimmed != "" {
				allowOrigins = append(allowOrigins, trimmed)
			}
		}
	}
	if len(allowOrigins) == 0 {
		allowOrigins = []string{"http://localhost:3000"}
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH", "HEAD"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
	}))

	// validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// jwt service
	jwtService := auth.NewJWTService(cfg.JwtSecret)

	// static files
	e.Static("/uploads", "uploads")

	// register routes (add korbo ekta ekta kore)
	registerRoutes(e, db, jwtService)

	// start server
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
