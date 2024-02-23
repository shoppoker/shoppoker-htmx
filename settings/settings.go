package settings

import (
	"reflect"

	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"gorm.io/gorm"
)

var SettingsInstance *Settings

type Settings struct {
	gorm.Model

	PhoneNumber string `default:"+7 (999) 999-99-99"`
	Email       string `default:"admin@website.com"`
}

func InitSettings() error {
	SettingsInstance = &Settings{}

	storage.GormStorageInstance.DB.AutoMigrate(SettingsInstance)

	if err := storage.GormStorageInstance.DB.Where("id = ?", 1).First(SettingsInstance).Error; err == nil {
		return nil
	}

	for i := 0; i < reflect.TypeOf(Settings{}).NumField(); i++ {
		field := reflect.TypeOf(Settings{}).Field(i)

		if field.Tag.Get("default") != "" {
			fieldValue := reflect.ValueOf(SettingsInstance).Elem().FieldByName(field.Name)
			fieldValue.SetString(field.Tag.Get("default"))
		}
	}

	return storage.GormStorageInstance.DB.Create(SettingsInstance).Error
}

func (s *Settings) Update() error {
	return storage.GormStorageInstance.DB.Save(s).Error
}
