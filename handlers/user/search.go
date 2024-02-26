package user_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherSearchRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/search", SearchHandler)
	user_api_group.GET("/search", SearchApiHandler)
	user_api_group.POST("/search", PostSearchHandler)
}

func SearchHandler(c echo.Context) error {
	search := c.QueryParam("search")
	if search == "" {
		return utils.Render(c, user_templates.Search([]*models.Product{}, ""))
	}
	var products []*models.Product
	if err := storage.GormStorageInstance.DB.
		Order("created_at DESC").
		Where("LOWER(title) LIKE LOWER(?) OR LOWER(tags) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%").
		Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, user_templates.Search(products, search))
}

func SearchApiHandler(c echo.Context) error {
	search := c.QueryParam("search")
	if search == "" {
		return utils.Render(c, user_templates.SearchApi([]*models.Product{}, ""))
	}
	var products []*models.Product
	if err := storage.GormStorageInstance.DB.
		Order("created_at DESC").
		Where("LOWER(title) LIKE LOWER(?) OR LOWER(tags) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%").
		Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, user_templates.SearchApi(products, search))
}

func PostSearchHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	search := c.Request().FormValue("search")
	if search == "" {
		return c.NoContent(http.StatusOK)
	}

	var products []*models.Product
	if err := storage.GormStorageInstance.DB.Order("created_at DESC").Where("LOWER(title) LIKE LOWER(?) OR LOWER(tags) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%").Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	c.Response().Header().Set("HX-Replace-URL", "/search?search="+search)
	return utils.Render(c, user_templates.SearchList(products, search, 2))
}

func SearchPageHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	search := c.QueryParam("search")
	var products []*models.Product
	if err := storage.GormStorageInstance.DB.Order("created_at DESC").Limit(models.PRODUCTS_PER_PAGE).Offset((page-1)*models.PRODUCTS_PER_PAGE).Where("LOWER(title) LIKE LOWER(?) OR LOWER(tags) LIKE LOWER(?)", "%"+search+"%", "%"+search+"%").Where("is_enabled = ?", true).Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, user_templates.SearchList(products, search, page+1))
}
