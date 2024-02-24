package user_handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
	"gorm.io/gorm"
)

func GatherProductsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/products/:slug", ProductHandler)
	user_api_group.GET("/products/:slug", ProductApiHandler)
}

func ProductHandler(c echo.Context) error {
	var product *models.Product
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&product).Error; err != nil {
		return err
	}

	var cart_product *models.CartProduct
	if err := storage.GormStorageInstance.DB.Where("product_id = ? AND cart_id = ?", product.ID, utils.GetCartFromContext(c.Request().Context()).ID).First(&cart_product).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error(err)
			return err
		}

		cart_product = &models.CartProduct{
			Quantity: 0,
		}
	}

	return utils.Render(c, user_templates.Product(product, cart_product))
}

func ProductApiHandler(c echo.Context) error {
	var product *models.Product
	if err := storage.GormStorageInstance.DB.Where("slug = ?", c.Param("slug")).First(&product).Error; err != nil {
		return err
	}

	var cart_product *models.CartProduct
	if err := storage.GormStorageInstance.DB.Where("product_id = ? AND cart_id = ?", product.ID, utils.GetCartFromContext(c.Request().Context()).ID).First(&cart_product).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error(err)
			return err
		}

		cart_product = &models.CartProduct{
			Quantity: 0,
		}
	}

	return utils.Render(c, user_templates.ProductApi(product, cart_product))
}
