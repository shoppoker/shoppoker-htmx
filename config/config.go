package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
)

var ConfigInstance *Config

type Config struct {
	Port      string `env:"PORT"`
	JWTSecret string `env:"JWT_SECRET"`

	StorageType string `env:"STORAGE_TYPE"`

	SqlitePath       string `env:"SQLITE_PATH" default:"data.db"`
	PostgresHost     string `env:"POSTGRES_HOST" default:"localhost"`
	PostgresPort     int    `env:"POSTGRES_PORT" default:"5432"`
	PostgresUser     string `env:"POSTGRES_USER" default:"postgres"`
	PostgresPassword string `env:"POSTGRES_PASSWORD" default:"postgres"`
	PostgresDatabase string `env:"POSTGRES_DATABASE" default:"postgres"`

	ObjectStorageBucketName string `env:"OBJECT_STORAGE_BUCKET_NAME"`
}

func InitConfig() error {
	ConfigInstance = &Config{}

	for i := 0; i < reflect.TypeOf(Config{}).NumField(); i++ {
		field := reflect.TypeOf(Config{}).Field(i)

		envName := field.Tag.Get("env")

		envValue, ok := os.LookupEnv(envName)
		if !ok {
			if field.Tag.Get("default") != "" {
				envValue = field.Tag.Get("default")
				log.Warn(fmt.Sprintf("Environment variable %s not set. Using default value %s", envName, envValue))
			} else {
				return errors.NewEnvironmentVariableNotFoundError(envName)
			}
		}

		if field.Type.Kind() == reflect.Int {
			envValueInt, err := strconv.Atoi(envValue)
			if err != nil {
				return errors.NewEnvironmentVariableNotFoundError(envName)
			}
			fieldValue := reflect.ValueOf(ConfigInstance).Elem().FieldByName(field.Name)
			fieldValue.SetInt(int64(envValueInt))
			continue
		}
		fieldValue := reflect.ValueOf(ConfigInstance).Elem().FieldByName(field.Name)
		fieldValue.SetString(envValue)
	}

	return nil
}
