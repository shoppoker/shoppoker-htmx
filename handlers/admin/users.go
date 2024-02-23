package admin_handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherUsersRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/users", UsersIndexHandler)
	admin_api_group.GET("/users", UsersIndexApiHandler)
	admin_api_group.GET("/users/page/:page", UsersPageHandler)
	admin_api_group.GET("/users/add", AddAdminModalHandler)
	admin_api_group.POST("/users", PostUserHandler)
	admin_api_group.POST("/users/search", UsersSearchHandler)
	admin_api_group.GET("/users/:id/delete", DeleteAdminModalHandler)
	admin_api_group.DELETE("/users/:id", DeleteAdminHandler)
	admin_api_group.GET("/users/:id/edit", EditAdminModalHandler)
	admin_api_group.PUT("/users/:id", PutUserHandler)
}

func UsersIndexHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.USERS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(username) LIKE LOWER(?)", "%"+search+"%")
	}
	var users []*models.User
	if err := query.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		fmt.Printf("User: %s\n\n\n", user.Username)
	}

	return utils.Render(c, admin_templates.Users(users, search))
}

func UsersIndexApiHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.USERS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(username) LIKE LOWER(?)", "%"+search+"%")
	}

	var users []*models.User
	if err := query.Find(&users).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.UsersApi(users, search))
}

func UsersPageHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB
	var users []*models.User
	if search != "" {
		query = query.Where("LOWER(username) LIKE LOWER(?)", "%"+search+"%")
	}
	if err := query.Limit(models.USERS_PER_PAGE).Offset((page - 1) * models.USERS_PER_PAGE).Find(&users).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.UsersList(users, page+1, search))
}

func UsersSearchHandler(c echo.Context) error {
	search := c.FormValue("search")
	query := storage.GormStorageInstance.DB.Limit(models.USERS_PER_PAGE)
	var users []*models.User
	if search == "" {
		c.Response().Header().Set("HX-Replace-Url", "/admin/users")
		if err := storage.GormStorageInstance.DB.Find(&users).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	} else {
		c.Response().Header().Set("HX-Replace-Url", "/admin/users?search="+search)
		if err := query.Where("LOWER(username) LIKE LOWER(?)", "%"+search+"%").Find(&users).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	}

	return utils.Render(c, admin_templates.UsersList(users, 2, search))
}

func DeleteAdminModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var user *models.User
	if err := storage.GormStorageInstance.DB.First(&user, id).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.DeleteAdminModal(user))
}

func DeleteAdminHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.GormStorageInstance.DB.Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func AddAdminModalHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.AddAdminModal())
}

func PostUserHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	username := c.FormValue("username")
	if !models.IsUsernameValid(username) {
		return c.String(http.StatusBadRequest, models.GetUsernameRules())
	}

	if err := storage.GormStorageInstance.DB.Where("username = ?", username).First(&models.User{}).Error; err == nil {
		return c.String(http.StatusBadRequest, "Пользователь с таким именем уже существует")
	}

	password := c.FormValue("password")
	if !models.IsPasswordValid(password) {
		return c.String(http.StatusBadRequest, models.GetPasswordRules())
	}

	password_repeat := c.FormValue("password_repeat")
	if password != password_repeat {
		return c.String(http.StatusBadRequest, "Пароли не совпадают")
	}

	user, err := models.NewUser(username, password, true)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}
	if err := storage.GormStorageInstance.DB.Create(&user).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.User(user))
}

func EditAdminModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var user *models.User
	if err := storage.GormStorageInstance.DB.First(&user, id).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.EditAdminModal(user))
}

func PutUserHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	username := c.FormValue("username")
	if !models.IsUsernameValid(username) {
		return c.String(http.StatusBadRequest, models.GetUsernameRules())
	}

	password := c.FormValue("password")
	if password != "" {

		if !models.IsPasswordValid(password) {
			return c.String(http.StatusBadRequest, models.GetPasswordRules())
		}

		password_repeat := c.FormValue("password_repeat")
		if password != password_repeat {
			return c.String(http.StatusBadRequest, "Пароли не совпадают")
		}
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var user *models.User
	if err := storage.GormStorageInstance.DB.First(&user, id).Error; err != nil {
		return err
	}

	user.Username = username
	if password != "" {
		user.PasswordHash, err = user.HashPassword(password)
	}
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}
	if err := storage.GormStorageInstance.DB.Save(&user).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.User(user))
}
