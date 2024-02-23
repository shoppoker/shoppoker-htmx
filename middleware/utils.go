package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
)

func UseUrl(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path == "/metrics" || c.Request().URL.Path == "/metrics/*" {
			return next(c)
		}
		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "url", c.Request().URL.String())))
		return next(c)
	}
}

func UseCategories(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path == "/metrics" || c.Request().URL.Path == "/metrics/*" {
			return next(c)
		}

		var categories []*models.Category
		if err := storage.GormStorageInstance.DB.Where("parent_id = ? and is_enabled = ?", 0, true).Find(&categories).Error; err != nil {
			log.Error(err)
			return err
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "categories", categories)))
		return next(c)
	}
}
