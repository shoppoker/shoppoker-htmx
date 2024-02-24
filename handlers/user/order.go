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

func GatherOrdersRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/checkout", CheckoutHandler)
	user_api_group.GET("/checkout", CheckoutApiHandler)

	user_api_group.GET("/checkout/delivery_type_form", GetDeliveryTypeForm)
	user_api_group.POST("/checkout", PostOrderHandler)
}

func CheckoutHandler(c echo.Context) error {
	return utils.Render(c, user_templates.Checkout())
}

func CheckoutApiHandler(c echo.Context) error {
	return utils.Render(c, user_templates.CheckoutApi())
}

func GetDeliveryTypeForm(c echo.Context) error {
	dt, err := strconv.Atoi(c.QueryParam("delivery_type"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	delivery_type := models.DeliveryType(dt)

	return utils.Render(c, user_templates.GetForm(delivery_type))
}

func PostOrderHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Имя не может быть пустым")
	}

	message := c.FormValue("name")
	phone_number := c.FormValue("phone_number")
	if !utils.ValidatePhoneNumber(phone_number) {
		return c.String(http.StatusBadRequest, "Неправильный формат номера телефона")
	}
	email := c.FormValue("email")
	if !utils.ValidateEmail(email) {
		return c.String(http.StatusBadRequest, "Неправильный формат адреса электронной почты")
	}

	delivery_type_int, err := strconv.Atoi(c.FormValue("delivery_type"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	delivery_type := models.DeliveryType(delivery_type_int)

	var city, adress string
	if delivery_type == models.DeliveryTypeDelivery {
		city = c.FormValue("city")
		adress = c.FormValue("adress")

		if city == "" || adress == "" {
			return c.String(http.StatusBadRequest, "Город и адрес не могут быть пустыми")
		}
	}

	cart := utils.GetCartFromContext(c.Request().Context())
	if len(cart.Products) == 0 {
		return c.String(http.StatusBadRequest, "Корзина пуста")
	}

	order := models.NewOrder(
		name,
		phone_number,
		email,
		message,
		delivery_type,
		adress,
		city,
	)

	if err := storage.GormStorageInstance.DB.Create(order).Error; err != nil {
		log.Error(err)
		return err
	}

	var order_products []*models.OrderProduct
	for _, cart_product := range cart.Products {
		if cart_product.Quantity == 0 {
			continue
		}

		order_product := models.NewOrderProduct(
			cart_product.ProductId,
			order.ID,
			cart_product.Slug,
			cart_product.Title,
			cart_product.Price,
			cart_product.DiscountPrice,
			cart_product.Quantity,
		)

		if err := storage.GormStorageInstance.DB.Create(order_product).Error; err != nil {
			log.Error(err)
			return err
		}

		cart_product.Quantity = 0
		if err := storage.GormStorageInstance.DB.Save(cart_product).Error; err != nil {
			log.Error(err)
			return err
		}

		order_products = append(order_products, order_product)
	}

	order.Products = order_products

	c.Response().Header().Del("HX-Reswap")

	return utils.Render(c, user_templates.CheckoutComplete(order))
}
