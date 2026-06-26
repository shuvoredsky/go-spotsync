package server

import (
	"spotsync/internal/auth"
	"spotsync/internal/config"

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}))

	// validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// jwt service
	jwtService := auth.NewJWTService(cfg.JwtSecret)

	// register routes (add korbo ekta ekta kore)
	registerRoutes(e, db, jwtService)

	// start server
	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
