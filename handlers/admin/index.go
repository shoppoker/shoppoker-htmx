package admin_handlers

import (
	"github.com/labstack/echo/v4"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherIndexRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("", AdminIndexHandler)
	admin_api_group.GET("/index", AdminIndexApiHandler)
}

func AdminIndexHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.Index())
}

func AdminIndexApiHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.IndexApi())
}
