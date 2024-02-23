package user_handlers

import "github.com/labstack/echo/v4"

func HealthHandler(c echo.Context) error {
	return c.String(200, "OK")
}
