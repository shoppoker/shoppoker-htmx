package user_handlers

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/config"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
	"gorm.io/gorm"
)

func GatherLoginRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/admin_login", LoginPageHandler)
	user_api_group.GET("/admin_login", LoginPageApiHandler)

	user_page_group.GET("/logout", LogoutHandler)

	user_api_group.POST("/admin_login", PostLoginHandler)
}

func LoginPageApiHandler(c echo.Context) error {
	if c.Request().Context().Value("user") != nil {
		c.Response().Header().Set("HX-Redirect", "/admin")
		c.Response().Header().Set("HX-Replace-Url", "/admin")
		return c.Redirect(http.StatusFound, "/admin")
	}

	return utils.Render(c, user_templates.LoginApi())
}

func LoginPageHandler(c echo.Context) error {
	if c.Request().Context().Value("user") != nil {
		c.Response().Header().Set("HX-Redirect", "/admin")
		c.Response().Header().Set("HX-Replace-Url", "/admin")
		return c.Redirect(http.StatusFound, "/admin")
	}

	return utils.Render(c, user_templates.Login())
}

func PostLoginHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" {
		return c.String(http.StatusBadRequest, "Имя пользователя не может быть пустым")
	}

	if password == "" {
		return c.String(http.StatusBadRequest, "Пароль не может быть пустым")
	}

	var user *models.User
	if err := storage.GormStorageInstance.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.String(http.StatusBadRequest, "Неправильный логин или пароль")
		}
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	if !user.ComparePassword(password) {
		return c.String(http.StatusBadRequest, "Неправильный логин или пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})

	tokenString, err := token.SignedString([]byte(config.ConfigInstance.JWTSecret))
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	http.SetCookie(c.Response(), &http.Cookie{
		Name:  "auth_token",
		Value: tokenString,
		Path:  "/",
	})

	c.Response().Header().Set("HX-Redirect", "/admin")
	c.Response().Header().Set("HX-Replace-Url", "/admin")
	return c.Redirect(http.StatusFound, "/admin")
}

func LogoutHandler(c echo.Context) error {
	http.SetCookie(c.Response(), &http.Cookie{
		Name:   "auth_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	c.Response().Header().Set("HX-Redirect", "/")
	c.Response().Header().Set("HX-Replace-Url", "/")
	return c.Redirect(http.StatusFound, "/")

}
