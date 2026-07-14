package config

import (
	"os"
	"reflect"
	"strconv"

	"github.com/golobby/dotenv"
)

type AppConfig struct {
	App struct {
		Name    string `env:"APP_NAME"`
		Version string `env:"APP_VERSION"`
	}
	Server struct {
		Host string `env:"HTTP_HOST"`
		Port string `env:"HTTP_PORT"`
	}
	Logger struct {
		Level string `env:"LOG_LEVEL"`
	}
	Postgres struct {
		ConnectionName string `env:"POSTGRES_CONNECTION_NAME"`
		Host           string `env:"POSTGRES_HOST"`
		Port           int    `env:"POSTGRES_PORT"`
		Username       string `env:"POSTGRES_USERNAME"`
		Password       string `env:"POSTGRES_PASSWORD"`
		DBName         string `env:"POSTGRES_DBNAME"`
		SSL            bool   `env:"POSTGRES_SSL"`
		MaxIdleConns   int    `env:"POSTGRES_MAX_IDLE_CONNS"`
		MaxOpenConns   int    `env:"POSTGRES_MAX_OPEN_CONNS"`
	}
}

func NewAppConfig() (*AppConfig, error) {
	config := &AppConfig{}

	file, err := os.Open(".env")
	if err == nil {
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				_ = closeErr
			}
		}()
		if err := dotenv.NewDecoder(file).Decode(config); err != nil {
			return config, err
		}
	}

	if err := loadFromEnv(config); err != nil {
		return config, err
	}

	return config, nil
}

func loadFromEnv(config *AppConfig) error {
	v := reflect.ValueOf(config).Elem()
	t := reflect.TypeOf(config).Elem()

	return processStruct(v, t)
}

func processStruct(v reflect.Value, t reflect.Type) error {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.Struct {
			if err := processStruct(field, fieldType.Type); err != nil {
				return err
			}
			continue
		}

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}
		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}
		if err := setFieldValue(field, envValue); err != nil {
			return err
		}
	}
	return nil
}

func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return nil
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return nil
		}
		field.SetBool(boolVal)
	}

	return nil
}
