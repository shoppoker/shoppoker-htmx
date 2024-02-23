package admin_handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherExamplesRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/examples", ExamplesHandler)
	admin_api_group.GET("/examples", ExamplesApiHandler)
	admin_api_group.GET("/examples/page/:page", ExamplesPageHandler)
	admin_api_group.POST("/examples/search", SearchExamplesHandler)
	admin_api_group.GET("/examples/add", AddExampleModalHandler)
	admin_api_group.POST("/examples", PostExampleHandler)
	admin_api_group.GET("/examples/:id/edit", EditExampleModalHandler)
	admin_api_group.PUT("/examples/:id", PutExampleHandler)
	admin_api_group.GET("/examples/:id/delete", DeleteExampleModalHandler)
	admin_api_group.DELETE("/examples/:id", DeleteExampleHandler)
}

func ExamplesHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var examples []*models.Example
	if err := query.Find(&examples).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.Examples(examples, search))
}

func ExamplesApiHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var examples []*models.Example
	if err := query.Find(&examples).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.ExamplesApi(examples, search))
}

func ExamplesPageHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE).Offset((page - 1) * models.EXAMPLES_PER_PAGE)
	var examples []*models.Example
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}
	if err := query.Find(&examples).Error; err != nil {
		return err
	}
	return utils.Render(c, admin_templates.ExamplesList(examples, search, page+1))
}

func SearchExamplesHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	search := c.FormValue("search")
	var examples []*models.Example
	query := storage.GormStorageInstance.DB.Limit(models.EXAMPLES_PER_PAGE)
	if search == "" {
		c.Response().Header().Set("HX-Replace-Url", "/admin/examples")
		if err := storage.GormStorageInstance.DB.Find(&examples).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	} else {
		c.Response().Header().Set("HX-Replace-Url", "/admin/examples?search="+search)
		if err := query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%").Find(&examples).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	}

	return utils.Render(c, admin_templates.ExamplesList(examples, search, 2))
}

func AddExampleModalHandler(c echo.Context) error {
	var custom_chip_bases []*models.CustomChipBase
	if err := storage.GormStorageInstance.DB.Find(&custom_chip_bases).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.AddExampleModal(custom_chip_bases))
}

func PostExampleHandler(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(100 * 1024 * 1024); err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	tags_split := c.FormValue("tags")
	if tags_split == "" {
		return c.String(http.StatusBadRequest, "Теги не могут быть пустыми")
	}
	tags := strings.Split(tags_split, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	custom_chip_base_id, err := strconv.ParseUint(c.FormValue("custom_chip_base_id"), 10, 64)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.GormStorageInstance.DB.First(&models.CustomChipBase{}, custom_chip_base_id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	files := form.File["images"]

	var images, thumbnails []string
	for _, file := range files {
		image_file, err := file.Open()
		if err != nil {
			return c.String(http.StatusBadRequest, "Неправильный запрос")
		}
		image, thumbnail, err := utils.ProcessImage(image_file)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}

		image_id, err := file_storage.FileStorageInstance.UploadFile(image)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}

		thumbnail_id, err := file_storage.FileStorageInstance.UploadFile(thumbnail)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}

		images = append(images, string(image_id))
		thumbnails = append(thumbnails, string(thumbnail_id))
	}

	example := models.NewExample(
		title,
		tags,
		images,
		thumbnails,
		uint(custom_chip_base_id),
	)

	if err := storage.GormStorageInstance.DB.Create(&example).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Example(example))
}

func DeleteExampleModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var example *models.Example
	if err := storage.GormStorageInstance.DB.First(&example, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.DeleteExampleModal(example))
}

func DeleteExampleHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}
	if err := storage.GormStorageInstance.DB.Delete(&models.Example{}, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}
	return c.NoContent(http.StatusOK)
}

func EditExampleModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var example *models.Example
	if err := storage.GormStorageInstance.DB.First(&example, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	var custom_chip_bases []*models.CustomChipBase
	if err := storage.GormStorageInstance.DB.Find(&custom_chip_bases).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.EditExampleModal(example, custom_chip_bases))
}

func PutExampleHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	tags_split := c.FormValue("tags")
	if tags_split == "" {
		return c.String(http.StatusBadRequest, "Теги не могут быть пустыми")
	}
	tags := strings.Split(tags_split, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}

	var example *models.Example
	if err := storage.GormStorageInstance.DB.First(&example, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	custom_chip_base_id, err := strconv.ParseUint(c.FormValue("custom_chip_base_id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.GormStorageInstance.DB.First(&models.CustomChipBase{}, custom_chip_base_id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	example.Title = title
	example.Tags = tags
	example.CustomChipBaseID = uint(custom_chip_base_id)

	if err := storage.GormStorageInstance.DB.Save(&example).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Example(example))
}
