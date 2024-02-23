package middleware

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/config"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
)

func UseAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth_token, err := c.Cookie("auth_token")
		if err != nil {
			return next(c)
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(auth_token.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.ConfigInstance.JWTSecret), nil
		})
		if err != nil {
			return next(c)
		}

		if !token.Valid {
			return next(c)
		}

		if claims["user_id"] == nil {
			return next(c)
		}

		user_id := uint(claims["user_id"].(float64))
		var user *models.User
		if err := storage.GormStorageInstance.DB.First(&user, user_id).Error; err != nil {
			return next(c)
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "user", user)))

		return next(c)
	}
}

func UseAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Robots-Tag", "noindex, nofollow")

		user_context := c.Request().Context().Value("user")
		if user_context == nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		user := user_context.(*models.User)
		if user.IsAdmin {
			return next(c)
		}

		return c.NoContent(http.StatusForbidden)
	}
}
