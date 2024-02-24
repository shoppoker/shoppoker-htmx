package user_handlers

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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

func getCategoryChildrenProducts(products []*models.Product, children []*models.Category, sort string) []*models.Product {
	for _, child := range children {
		if err := storage.GormStorageInstance.DB.Where("category_id = ?", child.ID).Find(&child.Products).Error; err != nil {
			log.Error(err)
			return []*models.Product{}
		}
		products = append(products, child.Products...)
	}

	sort_func := func(a, b *models.Product) int {
		return 0
	}

	if sort == "asc" {
		sort_func = func(a, b *models.Product) int {
			price_1 := a.Price
			if a.DiscountPrice > 0 {
				price_1 = a.DiscountPrice
			}

			price_2 := b.Price
			if b.DiscountPrice > 0 {
				price_2 = b.DiscountPrice
			}

			return price_1 - price_2
		}
	} else if sort == "desc" {
		sort_func = func(a, b *models.Product) int {
			price_1 := a.Price
			if a.DiscountPrice > 0 {
				price_1 = a.DiscountPrice
			}

			price_2 := b.Price
			if b.DiscountPrice > 0 {
				price_2 = b.DiscountPrice
			}

			return price_2 - price_1
		}
	}
	slices.SortFunc(products, sort_func)

	return products
}

func CategoryApiHandler(c echo.Context) error {
	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)

	sort := c.QueryParam("sort")
	if sort != "" && sort != "recommended" {
		query = query.Order("CASE WHEN discount_price > 0 THEN discount_price ELSE price END " + sort)
	}

	var category *models.Category
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&category).Error; err != nil {
		return err
	}

	if err := query.Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
		return err
	}

	storage.GormStorageInstance.DB.First(&category.Parent, category.ParentId)
	category.Products = getCategoryChildrenProducts(category.Products, category.Children, sort)

	return utils.Render(c, user_templates.CategoryApi(category, sort))
}

func CategoryHandler(c echo.Context) error {
	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)

	sort := c.QueryParam("sort")
	if sort != "" && sort != "recommended" {
		query = query.Order("CASE WHEN discount_price > 0 THEN discount_price ELSE price END " + sort)
	}

	var category *models.Category
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&category).Error; err != nil {
		return err
	}

	storage.GormStorageInstance.DB.First(&category.Parent, category.ParentId)

	if err := query.Where("category_id = ?", category.ID).Find(&category.Products).Error; err != nil {
		return err
	}

	category.Products = getCategoryChildrenProducts(category.Products, category.Children, sort)

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
	if err := storage.GormStorageInstance.DB.Where("category_id = ?", category.ID).Find(&products).Error; err != nil {
		return err
	}

	products = getCategoryChildrenProducts(products, category.Children, sort)

	return utils.Render(c, user_templates.ProductList(products, slug, page+1, sort))
}
