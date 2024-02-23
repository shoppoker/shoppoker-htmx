package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
)

func GatherFilesHandler(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/file/:file_type/:file_extension/:object_storage_id", S3FilesHandler)
}

func S3FilesHandler(c echo.Context) error {
	object_storage_id := c.Param("object_storage_id")
	if object_storage_id == "" {
		return c.NoContent(http.StatusNotFound)
	}

	file_bytes, err := file_storage.FileStorageInstance.DownloadFile(file_storage.ObjectStorageId(object_storage_id))
	if err != nil {
		log.Error(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.Blob(http.StatusOK, fmt.Sprintf("%s/%s", c.Param("file_type"), c.Param("file_extension")), file_bytes)
}
