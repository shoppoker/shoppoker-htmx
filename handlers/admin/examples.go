package admin_handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherExamplesRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/examples", ExamplesHandler)
	admin_api_group.GET("/examples", ExamplesApiHandler)
}

func ExamplesHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var examples []*models.Example
	if err := query.Find(&examples).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.Examples(examples, search))
}

func ExamplesApiHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var examples []*models.Example
	if err := query.Find(&examples).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.ExamplesApi(examples, search))
}
