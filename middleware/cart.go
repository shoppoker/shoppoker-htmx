package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"gorm.io/gorm"
)

func UseCart(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path == "/metrics" || c.Request().URL.Path == "/metrics/*" {
			return next(c)
		}

		var uuid_ uuid.UUID
		cart_uuid, err := c.Cookie("cart_uuid")
		uuid_value := ""
		if err != nil && err != http.ErrNoCookie {
			log.Error(err)
			return err
		} else if err == nil {
			uuid_value = cart_uuid.Value
		}

		uuid_, _ = uuid.Parse(uuid_value)
		var cart *models.Cart
		if err := storage.GormStorageInstance.DB.Where("uuid = ?", uuid_).First(&cart).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error(err)
				return err
			}

			cart = models.NewCart()
			if err := storage.GormStorageInstance.DB.Create(&cart).Error; err != nil {
				log.Error(err)
				return err
			}

			c.SetCookie(&http.Cookie{
				Name:  "cart_uuid",
				Path:  "/",
				Value: cart.UUID.String(),
			})
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "cart", cart)))
		return next(c)
	}
}
