package user_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherExamplesRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/examples", Examples)
	user_api_group.GET("/examples", ExamplesApi)
	user_api_group.GET("/examples/page/:page", ExamplesPage)
}

func Examples(c echo.Context) error {
	var examples []*models.Example

	if err := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE).Find(&examples).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.Examples(examples))
}

func ExamplesApi(c echo.Context) error {
	var examples []*models.Example

	if err := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE).Find(&examples).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.ExamplesApi(examples))
}

func ExamplesPage(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var examples []*models.Example
	if err := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE).Offset((page - 1) * models.EXAMPLES_PER_PAGE).Find(&examples).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.ExamplesList(examples, page+1))
}
