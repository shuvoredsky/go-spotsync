package middlewares

import (
	"net/http"
	"spotsync/internal/auth"
	"spotsync/internal/httpresponse"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, httpresponse.NewError("Missing authorization header", nil))
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, httpresponse.NewError("Invalid authorization header format", nil))
			}

			claims, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, httpresponse.NewError("Invalid or expired token", nil))
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, _ := c.Get("user_role").(string)
		if role != "admin" {
			return c.JSON(http.StatusForbidden, httpresponse.NewError("Access denied: admin only", nil))
		}
		return next(c)
	}
}
