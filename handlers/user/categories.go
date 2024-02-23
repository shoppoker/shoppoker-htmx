package user_handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherCategoriesRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/categories/:slug", CategoryHandler)
	user_api_group.GET("/categories/:slug", CategoryApiHandler)
}

func CategoryApiHandler(c echo.Context) error {
	var category *models.Category
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&category).Error; err != nil {
		return err
	}

	if err := storage.GormStorageInstance.DB.Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.CategoryApi(category))
}

func CategoryHandler(c echo.Context) error {
	var category *models.Category
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&category).Error; err != nil {
		return err
	}

	if err := storage.GormStorageInstance.DB.Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.Category(category))
}
