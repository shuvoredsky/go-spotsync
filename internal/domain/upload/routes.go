package upload

import (
	"spotsync/internal/auth"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, jwtService auth.JWTService) {
	handler := NewUploadHandler()

	uploadGroup := e.Group("/api/v1/upload")
	uploadGroup.POST("", handler.UploadImage)
	uploadGroup.POST("/image", handler.UploadImage)
}
