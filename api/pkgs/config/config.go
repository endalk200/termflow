package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

// AppConfig holds all configuration for the application
type AppConfig struct {
	Port            int    `env:"PORT" default:"8080"`
	ApplicationPort int    `env:"APPLICATION_PORT" default:"8080"`
	DbHost          string `env:"DB_HOST" required:"true"`
	DbPort          int    `env:"DB_PORT" default:"5432"`
	DbUser          string `env:"DB_USER" required:"true"`
	DbPass          string `env:"DB_PASS" required:"true"`
	DbName          string `env:"DB_NAME" required:"true"`
	LogLevel        string `env:"LOG_LEVEL" default:"INFO"`
}

// LoadConfig dynamically loads environment variables into the config struct
func LoadConfig(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structField := t.Field(i)

		// Get the environment variable name from the "env" tag
		envName := structField.Tag.Get("env")
		if envName == "" {
			return fmt.Errorf("missing 'env' tag for field %s", structField.Name)
		}

		// Get the value of the environment variable
		envValue := os.Getenv(envName)

		// If the environment variable is not set, check for default value
		if envValue == "" {
			if defaultValue, ok := structField.Tag.Lookup("default"); ok {
				envValue = defaultValue
			}
		}

		// If the environment variable is required but not set, return an error
		if envValue == "" && structField.Tag.Get("required") == "true" {
			return fmt.Errorf("required environment variable %s is not set", envName)
		}

		// Set the field based on its type
		if err := setFieldValue(field, envValue); err != nil {
			return fmt.Errorf("error setting field %s: %v", structField.Name, err)
		}
	}

	return nil
}

// setFieldValue sets a field to a given string value, converting types as necessary
func setFieldValue(field reflect.Value, value string) error {
	if value == "" {
		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value: %v", err)
		}
		field.SetInt(int64(intValue))
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value: %v", err)
		}
		field.SetBool(boolValue)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}
