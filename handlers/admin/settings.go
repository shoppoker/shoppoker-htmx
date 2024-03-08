package admin_handlers

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
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
	if c.Request().ParseMultipartForm(200*1024*1024) != nil {
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

	wholesale := c.FormValue("wholesale")
	if wholesale == "" {
		return c.String(http.StatusBadRequest, "Текст для оптовой продажи не может быть путой")
	}

	wholesale_file, err := c.FormFile("wholesale_file")
	if err == nil {
		wholesale_file_opened, err := wholesale_file.Open()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}

		var content bytes.Buffer
		_, err = io.Copy(&content, wholesale_file_opened)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}

		wholesale_file_id, err := file_storage.FileStorageInstance.UploadFile(content.Bytes())
		if err != nil {
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}

		s := strings.Split(wholesale_file.Filename, ".")
		extension := s[len(s)-1]
		settings.SettingsInstance.WholeSaleFileExtension = extension

		settings.SettingsInstance.WholeSaleFile = wholesale_file_id
	}

	settings.SettingsInstance.WholeSale = wholesale
	settings.SettingsInstance.PhoneNumber = phone_number
	settings.SettingsInstance.Email = email
	settings.SettingsInstance.WhatsappUrl = whatsapp_url
	settings.SettingsInstance.TelegramUrl = telegram_url

	if err := settings.SettingsInstance.Update(); err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return c.String(http.StatusOK, "Настройки обновлены")
}
