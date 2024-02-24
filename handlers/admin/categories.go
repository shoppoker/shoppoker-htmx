package admin_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"github.com/w1png/go-htmx-ecommerce-template/utils"

	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
)

func GatherCategoriesRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/categories", CategoriesIndexHandler)
	admin_api_group.GET("/categories", CategoriesIndexApiHandler)
	admin_api_group.GET("/categories/page/:page", CategoriesPageHandler)
	admin_api_group.POST("/categories/search", CategoriesSearchHandler)
	admin_api_group.GET("/categories/add", AddCategoryModalHandler)
	admin_api_group.POST("/categories", PostCategoryHandler)
	admin_api_group.GET("/categories/:id/delete", DeleteCategoryModalHandler)
	admin_api_group.DELETE("/categories/:id", DeleteCategoryHandler)
	admin_api_group.GET("/categories/:id/edit", EditCategoryModalHandler)
	admin_api_group.PUT("/categories/:id", PutCategoryHandler)
}

func CategoriesIndexHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.CATEGORIES_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}
	var categories []*models.Category
	if err := query.Find(&categories).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Categories(categories, search))
}

func CategoriesIndexApiHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.CATEGORIES_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}
	var categories []*models.Category
	if err := query.Find(&categories).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.CategoriesApi(categories, search))
}

func CategoriesPageHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB
	var categories []*models.Category
	if search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%")
	}
	if err := query.Limit(models.CATEGORIES_PER_PAGE).Offset((page - 1) * models.CATEGORIES_PER_PAGE).Find(&categories).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.CategoriesList(categories, page+1, search))
}

func CategoriesSearchHandler(c echo.Context) error {
	search := c.FormValue("search")
	query := storage.GormStorageInstance.DB.Limit(models.CATEGORIES_PER_PAGE)
	var categories []*models.Category
	if search == "" {
		c.Response().Header().Set("HX-Replace-Url", "/admin/categories")
		if err := storage.GormStorageInstance.DB.Find(&categories).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	} else {
		c.Response().Header().Set("HX-Replace-Url", "/admin/categories?search="+search)
		if err := query.Where("LOWER(name) LIKE LOWER(?)", "%"+search+"%").Find(&categories).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	}

	return utils.Render(c, admin_templates.CategoriesList(categories, 2, search))
}

func AddCategoryModalHandler(c echo.Context) error {
	var categories []*models.Category
	if err := storage.GormStorageInstance.DB.Where("parent_id = ?", 0).Find(&categories).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.AddCategoryModal(categories))
}

func PostCategoryHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Имя не может быть пустым")
	}

	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}

	if storage.GormStorageInstance.DB.Where("slug = ?", slug).First(&models.Category{}).Error == nil {
		return c.String(http.StatusBadRequest, "Категория с такой ссылкой уже существует")
	}

	tags := c.FormValue("tags")
	if tags == "" {
		return c.String(http.StatusBadRequest, "Теги не могут быть пустыми")
	}

	parent_id, err := strconv.ParseUint(c.FormValue("parent_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))

	category := models.NewCategory(name, slug, tags, uint(parent_id))
	category.IsEnabled = is_enabled
	if err := storage.GormStorageInstance.DB.Create(&category).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Category(category))
}

func DeleteCategoryModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var category *models.Category
	if err := storage.GormStorageInstance.DB.First(&category, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.DeleteCategoryModal(category))
}

func DeleteCategoryHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.GormStorageInstance.DB.Delete(&models.Category{}, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return c.NoContent(http.StatusOK)
}

func EditCategoryModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var category *models.Category
	if err := storage.GormStorageInstance.DB.First(&category, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	var categories []*models.Category
	if err := storage.GormStorageInstance.DB.Where("parent_id = ?", 0).Find(&categories).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.EditCategoryModal(category, categories))
}

func PutCategoryHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Имя не может быть пустым")
	}

	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}

	var slug_category *models.Category
	if storage.GormStorageInstance.DB.Where("slug = ?", slug).First(&slug_category).Error == nil && slug_category.ID != uint(id) {
		return c.String(http.StatusBadRequest, "Категория с такой ссылкой уже существует")
	}

	tags := c.FormValue("tags")
	if tags == "" {
		return c.String(http.StatusBadRequest, "Теги не могут быть пустыми")
	}

	parent_id, err := strconv.ParseUint(c.FormValue("parent_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))

	var category *models.Category
	if err := storage.GormStorageInstance.DB.First(&category, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	category.Name = name
	category.Slug = slug
	category.Tags = tags
	category.ParentId = uint(parent_id)
	category.IsEnabled = is_enabled

	if err := storage.GormStorageInstance.DB.Save(&category).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Category(category))
}
