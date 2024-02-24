package admin_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherProductsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/products", ProductsHandler)
	admin_api_group.GET("/products", ProductsApiHandler)
	admin_api_group.GET("/products/page/:page", ProductsPageHandler)
	admin_api_group.POST("/products/search", ProductsSearchHandler)
	admin_api_group.GET("/products/add", AddProductModalHandler)
	admin_api_group.POST("/products", PostProductHandler)
	admin_api_group.GET("/products/:id/edit", EditProductModalHandler)
	admin_api_group.PUT("/products/:id", PutProductHandler)
	admin_api_group.GET("/products/:id/delete", DeleteProductModalHandler)
	admin_api_group.DELETE("/products/:id", DeleteProductHandler)
}

func ProductsHandler(c echo.Context) error {
	search := c.QueryParam("search")

	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var products []*models.Product
	if err := query.Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Products(products, search))
}

func ProductsApiHandler(c echo.Context) error {
	search := c.QueryParam("search")

	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var products []*models.Product
	if err := query.Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.ProductsApi(products, search))
}

func ProductsPageHandler(c echo.Context) error {
	search := c.QueryParam("search")

	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var products []*models.Product
	if err := query.Offset((page - 1) * models.PRODUCTS_PER_PAGE).Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.ProductsList(products, page+1, search))
}

func ProductsSearchHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	search := c.FormValue("search")
	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)
	var products []*models.Product
	if search == "" {
		c.Response().Header().Set("HX-Replace-Url", "/admin/products")
		if err := storage.GormStorageInstance.DB.Find(&products).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	} else {
		c.Response().Header().Set("HX-Replace-Url", "/admin/products?search="+search)
		if err := query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%").Find(&products).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	}

	return utils.Render(c, admin_templates.ProductsList(products, 2, search))
}

func AddProductModalHandler(c echo.Context) error {
	var categories []*models.Category
	if err := storage.GormStorageInstance.DB.Find(&categories).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.AddProductModal(categories))
}

func PostProductHandler(c echo.Context) error {
	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}
	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}
	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}

	if storage.GormStorageInstance.DB.Where("slug = ?", slug).First(&models.Product{}).Error == nil {
		return c.String(http.StatusBadRequest, "Товар с такой ссылкой уже существует")
	}

	tags := c.FormValue("tags")

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}
	discount_price, err := strconv.Atoi(c.FormValue("discount_price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	category_id, err := strconv.ParseUint(c.FormValue("category"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	if category_id != 0 {
		if err := storage.GormStorageInstance.DB.First(&models.Category{}, category_id).Error; err != nil {
			return c.String(http.StatusBadRequest, "Неправильный запрос")
		}
	}

	stock_type_int, err := strconv.Atoi(c.FormValue("stock_type"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	stock_type := models.StockType(stock_type_int)
	if stock_type.ToString() == "" {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))
	is_featured, _ := strconv.ParseBool(c.FormValue("is_featured"))

	form, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
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

	product := models.NewProduct(
		slug,
		title,
		description,
		price,
		discount_price,
		stock_type,
		tags,
		uint(category_id),
		images,
		thumbnails,
		is_enabled,
		is_featured,
	)

	if err := storage.GormStorageInstance.DB.Create(&product).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Product(product))
}

func EditProductModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	var product *models.Product
	if err := storage.GormStorageInstance.DB.First(&product, id).Error; err != nil {
		return err
	}

	var categories []*models.Category
	if err := storage.GormStorageInstance.DB.Find(&categories).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.EditProductModal(product, categories))
}

func PutProductHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var product *models.Product
	if err := storage.GormStorageInstance.DB.First(&product, id).Error; err != nil {
		return err
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}
	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}
	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}

	var slug_product *models.Product
	if storage.GormStorageInstance.DB.Where("slug = ?", slug).First(&slug_product).Error == nil && slug_product.ID != uint(id) {
		return c.String(http.StatusBadRequest, "Товар с такой ссылкой уже существует")
	}

	tags := c.FormValue("tags")

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}
	discount_price, err := strconv.Atoi(c.FormValue("discount_price"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	category_id, err := strconv.ParseUint(c.FormValue("category"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	if category_id == 0 {
		if err := storage.GormStorageInstance.DB.First(&models.Category{}, category_id).Error; err != nil {
			return c.String(http.StatusBadRequest, "Неправильный запрос")
		}
	}

	stock_type_int, err := strconv.Atoi(c.FormValue("stock_type"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	stock_type := models.StockType(stock_type_int)
	if stock_type.ToString() == "" {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))
	is_featured, _ := strconv.ParseBool(c.FormValue("is_featured"))

	form, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
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

	product.Title = title
	product.Description = description
	product.Slug = slug
	product.Tags = tags
	product.Price = price
	product.DiscountPrice = discount_price
	product.StockType = stock_type
	product.CategoryId = uint(category_id)
	if len(thumbnails) != 0 {
		product.Images = images
		product.Thumbnails = thumbnails
	}
	product.IsEnabled = is_enabled
	product.IsFeatured = is_featured

	if err := storage.GormStorageInstance.DB.Save(&product).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Product(product))
}

func DeleteProductModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	var product *models.Product
	if err := storage.GormStorageInstance.DB.First(&product, id).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.DeleteProductModal(product))
}

func DeleteProductHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	if err := storage.GormStorageInstance.DB.Delete(&models.Product{}, id).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
