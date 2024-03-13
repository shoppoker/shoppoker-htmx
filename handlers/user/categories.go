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

func GatherCategoriesRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/categories/:slug", CategoryHandler)
	user_api_group.GET("/categories/:slug", CategoryApiHandler)
	user_api_group.GET("/categories/:slug/products/page/:page", CategoryProductPageHandler)
}

func CategoryApiHandler(c echo.Context) error {
	query := storage.GormStorageInstance.DB

	sort := c.QueryParam("sort")
	if sort != "" && sort != "recommended" {
		query = query.Order("CASE WHEN discount_price > 0 THEN discount_price ELSE price END " + sort)
	}

	var category *models.Category
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&category).Error; err != nil {
		return err
	}

	if err := query.Order("priority DESC").Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
		return err
	}

	storage.GormStorageInstance.DB.First(&category.Parent, category.ParentId)

	return utils.Render(c, user_templates.CategoryApi(category, sort))
}

func CategoryHandler(c echo.Context) error {
	query := storage.GormStorageInstance.DB

	sort := c.QueryParam("sort")
	if sort != "" && sort != "recommended" {
		query = query.Order("CASE WHEN discount_price > 0 THEN discount_price ELSE price END " + sort)
	}

	var category *models.Category
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&category).Error; err != nil {
		return err
	}

	storage.GormStorageInstance.DB.First(&category.Parent, category.ParentId)

	if err := query.Order("priority DESC").Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.Category(category, sort))
}

func CategoryProductPageHandler(c echo.Context) error {
	sort := c.QueryParam("sort")

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	slug := c.Param("slug")
	var category *models.Category
	if err := storage.GormStorageInstance.DB.Where("slug = ?", slug).First(&category).Error; err != nil {
		return err
	}

	if sort != "" && sort != "recommended" {
		c.Response().Header().Set("HX-Replace-Url", "/categories/"+slug+"?sort="+sort)
	} else {
		c.Response().Header().Set("HX-Replace-Url", "/categories/"+slug)
	}
	var products []*models.Product
	if err := storage.GormStorageInstance.DB.Order("priority DESC").Limit(models.PRODUCTS_PER_PAGE).Offset((page-1)*models.PRODUCTS_PER_PAGE).Where("category_id = ?", category.ID).Find(&products).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.ProductList(products, slug, page+1, sort))
}
