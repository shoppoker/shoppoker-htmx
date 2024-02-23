package admin_handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherCustomChipsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/custom_chip_bases", CustomChipBasesHandler)
	admin_api_group.GET("/custom_chip_bases", CustomChipBasesApiHandler)
	admin_api_group.POST("/custom_chip_bases/search", SearchCustomChipBasesHandler)
	admin_api_group.GET("/custom_chip_bases/add", AddCustomChipBaseModalHandler)
	admin_api_group.POST("/custom_chip_bases", PostCustomChipBaseHandler)
	admin_api_group.GET("/custom_chip_bases/:id/edit", EditCustomChipBaseModalHandler)
	admin_api_group.PUT("/custom_chip_bases/:id", PutCustomChipBaseHandler)
	admin_api_group.GET("/custom_chip_bases/:id/delete", DeleteCustomChipBaseModalHandler)
	admin_api_group.DELETE("/custom_chip_bases/:id", DeleteCustomChipBaseHandler)
}

func CustomChipBasesHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}
	var custom_chip_bases []*models.CustomChipBase
	if err := query.Find(&custom_chip_bases).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутреняя ошибка сервера")
	}

	return utils.Render(c, admin_templates.CustomChipBases(custom_chip_bases, search))
}

func CustomChipBasesApiHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}
	var custom_chip_bases []*models.CustomChipBase
	if err := query.Find(&custom_chip_bases).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Внутреняя ошибка сервера")
	}

	return utils.Render(c, admin_templates.CustomChipBasesApi(custom_chip_bases, search))
}

func SearchCustomChipBasesHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	search := c.FormValue("search")
	var custom_chip_bases []*models.CustomChipBase
	if search == "" {
		c.Response().Header().Set("HX-Replace-Url", "/admin/custom_chip_bases")
		if err := storage.GormStorageInstance.DB.Find(&custom_chip_bases).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	} else {
		c.Response().Header().Set("HX-Replace-Url", "/admin/custom_chip_bases?search="+search)
		if err := storage.GormStorageInstance.DB.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%").Find(&custom_chip_bases).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	}

	return utils.Render(c, admin_templates.CustomChipBasesList(custom_chip_bases))
}

func AddCustomChipBaseModalHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.AddCustomChipBaseModal())
}

func PostCustomChipBaseHandler(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(100 * 1024 * 1024); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}

	sticker_scale, err := strconv.ParseFloat(c.FormValue("sticker_scale"), 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))

	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}
	if storage.GormStorageInstance.DB.Where("slug = ?", slug).First(&models.CustomChipBase{}).Error == nil {
		return c.String(http.StatusBadRequest, "Товар с такой ссылкой уже существует")
	}

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil || price < 0 {
		return c.String(http.StatusBadRequest, "Неверная цена")
	}
	discount_price, err := strconv.Atoi(c.FormValue("discount_price"))
	if err != nil || discount_price < 0 {
		return c.String(http.StatusBadRequest, "Неверная цена со скидкой")
	}

	price_foil, err := strconv.Atoi(c.FormValue("price_foil"))
	if err != nil || price_foil < 0 {
		return c.String(http.StatusBadRequest, "Неверная цена для голографического ламината")
	}
	discount_price_foil, err := strconv.Atoi(c.FormValue("discount_price_foil"))
	if err != nil || discount_price_foil < 0 {
		return c.String(http.StatusBadRequest, "Неверная цена со скидкой для голографического ламината")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "Необходимо выбрать изображение")
	}

	image, err := file.Open()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	var image_bytes []byte
	if image_bytes, err = io.ReadAll(image); err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	image_id, err := file_storage.FileStorageInstance.UploadFile(image_bytes)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	custom_chip_base := models.NewCustomChipBase(
		title,
		slug,
		description,
		sticker_scale,
		image_id,
		is_enabled,
	)

	if err := storage.GormStorageInstance.DB.Create(custom_chip_base).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.CustomChipBase(custom_chip_base))
}

func DeleteCustomChipBaseModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var custom_chip_base *models.CustomChipBase
	if err := storage.GormStorageInstance.DB.First(&custom_chip_base, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.DeleteCustomChipBaseModal(custom_chip_base))
}

func DeleteCustomChipBaseHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.GormStorageInstance.DB.Delete(&models.CustomChipBase{}, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return c.NoContent(http.StatusOK)
}

func EditCustomChipBaseModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var custom_chip_base *models.CustomChipBase
	if err := storage.GormStorageInstance.DB.First(&custom_chip_base, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.EditCustomChipBaseModal(custom_chip_base))
}

func PutCustomChipBaseHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}

	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}

	sticker_scale, err := strconv.ParseFloat(c.FormValue("sticker_scale"), 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверная шкала")
	}

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	discount_price, err := strconv.Atoi(c.FormValue("discount_price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))

	foil_price, err := strconv.Atoi(c.FormValue("foil_price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	foil_discount_price, err := strconv.Atoi(c.FormValue("foil_discount_price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	var custom_chip_base *models.CustomChipBase
	if err := storage.GormStorageInstance.DB.First(&custom_chip_base, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	custom_chip_base.Title = title
	custom_chip_base.Slug = slug
	custom_chip_base.Description = description
	custom_chip_base.StickerScale = sticker_scale
	custom_chip_base.Price = price
	custom_chip_base.DiscountPrice = discount_price
	custom_chip_base.IsEnabled = is_enabled
	custom_chip_base.PriceFoil = foil_price
	custom_chip_base.DiscountPriceFoil = foil_discount_price

	if err := storage.GormStorageInstance.DB.Save(&custom_chip_base).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return c.NoContent(http.StatusOK)
}
