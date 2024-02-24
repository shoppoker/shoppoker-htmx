package user_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"github.com/w1png/go-htmx-ecommerce-template/templates/components"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
	"gorm.io/gorm"
)

func GatherCartRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_api_group.GET("/cart", GetCartHandler)
	user_api_group.PUT("/cart/change_quantity/:product_id", ChangeCartProductQuantityHandler)
}

func GetCartHandler(c echo.Context) error {
	var cart_products []*models.CartProduct
	for _, cart_product := range utils.GetCartFromContext(c.Request().Context()).Products {
		if cart_product.Quantity != 0 {
			cart_products = append(cart_products, cart_product)
		}
	}

	// for _, cart_product := range cart_products {
	// 	if cart_product.Product.StockType == models.StockTypeOutOfStock {
	// 		cart_product.Quantity = 0
	// 	}
	// }

	return utils.Render(c, user_templates.CartProducts(cart_products))
}

func ChangeCartProductQuantityHandler(c echo.Context) error {
	should_decrease := c.QueryParam("decrease") == "true"

	product_id, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		return err
	}

	var product *models.Product
	if err := storage.GormStorageInstance.DB.Where("id = ?", product_id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.String(http.StatusNotFound, "Товар не найден")
		}
		log.Error(err)
		return err
	}

	if product.StockType == models.StockTypeOutOfStock {
		return c.String(http.StatusBadRequest, "Товара нет в наличии")
	}

	cart := utils.GetCartFromContext(c.Request().Context())

	var cart_product *models.CartProduct
	if err := storage.GormStorageInstance.DB.Where("product_id = ? AND cart_id = ?", product_id, cart.ID).First(&cart_product).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error(err)
			return err
		}

		cart_product = models.NewCartProduct(
			product.ID,
			cart.ID,
			product.Slug,
			product.Title,
			product.Price,
			product.DiscountPrice,
			product.Thumbnails[0],
			0,
		)
	}

	if should_decrease {
		cart_product.Quantity--
	} else {
		cart_product.Quantity++
	}

	if cart_product.Quantity < 0 {
		cart_product.Quantity = 0
	}

	if err = storage.GormStorageInstance.DB.Save(&cart_product).Error; err != nil {
		log.Error(err)
		return err
	}

	return utils.Render(c, components.AddToCartButton(cart_product.ProductId, cart_product.Quantity))
}
