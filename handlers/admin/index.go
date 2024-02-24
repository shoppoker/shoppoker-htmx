package admin_handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherIndexRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("", AdminIndexHandler)
	admin_api_group.GET("/index", AdminIndexApiHandler)
	admin_api_group.GET("/orders/:id/info", OpenOrderInfoModal)
	admin_api_group.GET("/orders/:id", UpdateOrderStatusHandler)
}

func OpenOrderInfoModal(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	var order *models.Order
	if err := storage.GormStorageInstance.DB.First(&order, id).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.OrderInfoModal(order))
}

func AdminIndexHandler(c echo.Context) error {
	var orders []*models.Order
	if err := storage.GormStorageInstance.DB.Find(&orders).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.Index(orders))
}

func AdminIndexApiHandler(c echo.Context) error {
	var orders []*models.Order
	if err := storage.GormStorageInstance.DB.Find(&orders).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.IndexApi(orders))
}

func UpdateOrderStatusHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	statusRaw, err := strconv.Atoi(c.QueryParam("status"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	status := models.OrderStatus(statusRaw)
	if status == models.OrderStatusAny {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	fmt.Printf("status: %v\n", status)

	var order *models.Order
	if err := storage.GormStorageInstance.DB.First(&order, id).Error; err != nil {
		return err
	}

	order.Status = status

	if err := storage.GormStorageInstance.DB.Save(&order).Error; err != nil {
		return err
	}

	c.Response().Header().Set("HX-Trigger", fmt.Sprintf("update_status_%d", order.ID))
	return c.NoContent(http.StatusOK)
}
