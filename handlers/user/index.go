package user_handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherIndexHandlers(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/", IndexHandler)
	user_api_group.GET("/index", IndexApiHandler)
}

func IndexApiHandler(c echo.Context) error {
	var featured_products []*models.Product
	if err := storage.GormStorageInstance.DB.Where("is_featured = ?", true).Find(&featured_products).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, user_templates.IndexApi(featured_products))
}

func IndexHandler(c echo.Context) error {
	var featured_products []*models.Product
	if err := storage.GormStorageInstance.DB.Where("is_featured = ? and is_enabled = ?", true, true).Find(&featured_products).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, user_templates.Index(featured_products))
}
