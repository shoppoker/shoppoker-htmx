package user_handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

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
	user_api_group.GET("/cart/clear", ClearCartHandler)
	user_api_group.GET("/cart/buttons/:id", GetCartButtonsHandler)
	user_api_group.GET("/cart/products/amount", GetCartProductsAmountHandler)
}

func ClearCartHandler(c echo.Context) error {
	var cart_products []*models.CartProduct
	for _, cart_product := range utils.GetCartFromContext(c.Request().Context()).Products {
		if cart_product.Quantity != 0 {
			cart_products = append(cart_products, cart_product)
		}
	}

	wg := sync.WaitGroup{}
	for _, cart_product := range cart_products {
		wg.Add(1)
		go func(cart_product *models.CartProduct, wg *sync.WaitGroup) {
			cart_product.Quantity = 0
			if err := storage.GormStorageInstance.DB.Save(cart_product).Error; err != nil {
				log.Error(err)
			}
			wg.Done()
		}(cart_product, &wg)
	}
	wg.Wait()

	c.Response().Header().Set("HX-Trigger", "cart_updated")

	return utils.Render(c, user_templates.CartProducts(cart_products))
}

func GetCartHandler(c echo.Context) error {
	var cart_products []*models.CartProduct
	for _, cart_product := range utils.GetCartFromContext(c.Request().Context()).Products {
		if cart_product.Quantity != 0 {
			cart_products = append(cart_products, cart_product)
		}
	}

	wg := sync.WaitGroup{}
	for _, cart_product := range cart_products {
		wg.Add(1)
		go func(cart_product *models.CartProduct, wg *sync.WaitGroup) {
			var product *models.Product
			if err := storage.GormStorageInstance.DB.Where("id = ?", cart_product.ProductId).First(&product).Error; err != nil {
				wg.Done()
				return
			}
			if product.StockType == models.StockTypeOutOfStock {
				cart_product.Quantity = 0
			}
			if err := storage.GormStorageInstance.DB.Save(cart_product).Error; err != nil {
				log.Error(err)
			}
			wg.Done()
		}(cart_product, &wg)
	}
	wg.Wait()

	return utils.Render(c, user_templates.CartProducts(cart_products))
}

func GetCartProductsAmountHandler(c echo.Context) error {
	return c.String(http.StatusOK, fmt.Sprintf("%d", utils.GetCartFromContext(c.Request().Context()).GetProductAmount()))
}

func GetCartButtonsHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusOK)
	}

	var cart_product *models.CartProduct
	if err := storage.GormStorageInstance.DB.Where("product_id = ? AND cart_id = ?", id, utils.GetCartFromContext(c.Request().Context()).ID).First(&cart_product).Error; err != nil {
		log.Error(err)
		return c.NoContent(http.StatusOK)
	}

	return utils.Render(c, components.AddToCartButton(uint(id), cart_product.Quantity))
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

	c.Response().Header().Set("HX-Trigger", "cart_updated")
	return utils.Render(c, components.AddToCartButton(cart_product.ProductId, cart_product.Quantity))
}
