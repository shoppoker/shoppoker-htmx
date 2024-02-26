package admin_handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/settings"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherSettingsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/settings", SettingsHandler)
	admin_api_group.GET("/settings", SettingsApiHandler)
	admin_api_group.PUT("/settings", PutSettingsHandler)
}

func SettingsHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.Settings())
}

func SettingsApiHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.SettingsApi())
}

func PutSettingsHandler(c echo.Context) error {
	if c.Request().ParseForm() != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	phone_number := c.FormValue("phone_number")
	if phone_number == "" {
		return c.String(http.StatusBadRequest, "Номер телефона не может быть путым")
	}

	email := c.FormValue("email")
	if email == "" {
		return c.String(http.StatusBadRequest, "Электронная почта не может быть путой")
	}

	whatsapp_url := c.FormValue("whatsapp_url")
	if whatsapp_url == "" {
		return c.String(http.StatusBadRequest, "Ссылка на WhatsApp не может быть путой")
	}

	telegram_url := c.FormValue("telegram_url")
	if telegram_url == "" {
		return c.String(http.StatusBadRequest, "Ссылка на Telegram не может быть путой")
	}

	settings.SettingsInstance.PhoneNumber = phone_number
	settings.SettingsInstance.Email = email
	settings.SettingsInstance.WhatsappUrl = whatsapp_url
	settings.SettingsInstance.TelegramUrl = telegram_url

	if err := settings.SettingsInstance.Update(); err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return c.String(http.StatusOK, "Настройки обновлены")
}
