package storage

import (
	"strconv"

	"github.com/w1png/go-htmx-ecommerce-template/config"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormStorage struct {
	DB *gorm.DB
}

var GormStorageInstance *GormStorage

func InitStorage() error {
	storage := &GormStorage{}

	var err error
	switch config.ConfigInstance.StorageType {
	case "sqlite":
		storage.DB, err = gorm.Open(sqlite.Open(config.ConfigInstance.SqlitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return err
		}
	case "postgres":
		storage.DB, err = gorm.Open(
			postgres.Open(
				"host="+config.ConfigInstance.PostgresHost+
					" port="+strconv.Itoa(config.ConfigInstance.PostgresPort)+
					" user="+config.ConfigInstance.PostgresUser+
					" password="+config.ConfigInstance.PostgresPassword+
					" dbname="+config.ConfigInstance.PostgresDatabase+
					" sslmode=disable",
			),
			&gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			},
		)
		if err != nil {
			return err
		}
	default:
		return errors.NewUnknownDatabaseTypeError(config.ConfigInstance.StorageType)
	}

	GormStorageInstance = storage

	return storage.DB.AutoMigrate(
		&models.Order{},
		&models.User{},
		&models.Cart{},
		&models.CartProduct{},
		&models.Category{},
		&models.Product{},
		&models.OrderProduct{},
	)
}
