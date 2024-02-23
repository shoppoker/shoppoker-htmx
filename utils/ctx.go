package utils

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/w1png/go-htmx-ecommerce-template/models"
)

type ResponseData struct {
	User *models.User
	Data interface{}
}

func MarshalResponse(c echo.Context, data interface{}) *ResponseData {
	var user *models.User
	userAny := c.Request().Context().Value("user")
	if userAny == nil {
		user = nil
	} else {
		user = userAny.(*models.User)
	}

	return &ResponseData{
		User: user,
		Data: data,
	}
}

func GetUserFromContext(ctx context.Context) *models.User {
	userAny := ctx.Value("user")
	if userAny == nil {
		return nil
	}
	return userAny.(*models.User)
}

func GetUrlFromContext(ctx context.Context) string {
	urlAny := ctx.Value("url")
	if urlAny == nil {
		return ""
	}
	return urlAny.(string)
}

func GetCategoriesFromContext(ctx context.Context) []*models.Category {
	categoriesAny := ctx.Value("categories")
	if categoriesAny == nil {
		return nil
	}
	return categoriesAny.([]*models.Category)
}

func GetCartFromContext(ctx context.Context) *models.Cart {
	cartAny := ctx.Value("cart")
	if cartAny == nil {
		return nil
	}
	return cartAny.(*models.Cart)
}
