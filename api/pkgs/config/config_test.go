package config

import (
	"os"
	"testing"
)

type testCfg struct {
	ApplicationPort int    `env:"APPLICATION_PORT" default:"8080"`
	DbHost          string `env:"DB_HOST" required:"true"`
	DbPort          int    `env:"DB_PORT" default:"5432"`
	DbUser          string `env:"DB_USER" required:"true"`
	DbPass          string `env:"DB_PASS" required:"true"`
	DbName          string `env:"DB_NAME" required:"true"`
	LogLevel        string `env:"LOG_LEVEL" default:"INFO"`
}

// TestLoadConfig_SuccessWithEnvVariables tests LoadConfig with all environment variables set
func TestLoadConfig_SuccessWithEnvVariables(t *testing.T) {
	os.Setenv("APPLICATION_PORT", "8081")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "3306")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "password")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("LOG_LEVEL", "DEBUG")

	var cfg testCfg

	// Call LoadConfig
	err := LoadConfig(&cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.ApplicationPort != 8081 {
		t.Errorf("Expected ApplicationPort to be 8081, got %d", cfg.ApplicationPort)
	}
	if cfg.DbHost != "localhost" {
		t.Errorf("Expected DbHost to be 'localhost', got %s", cfg.DbHost)
	}
	if cfg.DbPort != 3306 {
		t.Errorf("Expected DbPort to be 3306, got %d", cfg.DbPort)
	}
	if cfg.DbUser != "root" {
		t.Errorf("Expected DbUser to be 'root', got %s", cfg.DbUser)
	}
	if cfg.DbPass != "password" {
		t.Errorf("Expected DbPass to be 'password', got %s", cfg.DbPass)
	}
	if cfg.DbName != "testdb" {
		t.Errorf("Expected DbName to be 'testdb', got %s", cfg.DbName)
	}
	if cfg.LogLevel != "DEBUG" {
		t.Errorf("Expected LogLevel to be 'DEBUG', got %s", cfg.LogLevel)
	}
}

// TestLoadConfig_UseDefaultValues tests LoadConfig when some environment variables are missing, relying on default values
func TestLoadConfig_UseDefaultValues(t *testing.T) {
	// Unset all environment variables except required ones
	os.Clearenv()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "password")
	os.Setenv("DB_NAME", "testdb")

	var cfg testCfg

	// Call LoadConfig
	err := LoadConfig(&cfg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if cfg.ApplicationPort != 8080 {
		t.Errorf("Expected default ApplicationPort to be 8080, got %d", cfg.ApplicationPort)
	}
	if cfg.DbPort != 5432 {
		t.Errorf("Expected default DbPort to be 5432, got %d", cfg.DbPort)
	}
	if cfg.LogLevel != "INFO" {
		t.Errorf("Expected default LogLevel to be 'INFO', got %s", cfg.LogLevel)
	}
}

// TestLoadConfig_MissingRequiredEnv tests LoadConfig when required environment variables are missing
func TestLoadConfig_MissingRequiredEnv(t *testing.T) {
	// Unset all environment variables
	os.Clearenv()

	var cfg AppConfig

	// Call LoadConfig and expect an error
	err := LoadConfig(&cfg)
	if err == nil {
		t.Fatal("Expected error due to missing required environment variables, got nil")
	}

	// Ensure the error mentions the missing environment variable
	expectedError := "required environment variable DB_HOST is not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}

// TestLoadConfig_InvalidIntEnvValue tests LoadConfig when invalid integer values are provided for int fields
func TestLoadConfig_InvalidIntEnvValue(t *testing.T) {
	// Set invalid integer values for Port
	os.Setenv("APPLICATION_PORT", "invalid_int")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "root")
	os.Setenv("DB_PASS", "password")
	os.Setenv("DB_NAME", "testdb")

	var cfg testCfg

	// Call LoadConfig and expect an error
	err := LoadConfig(&cfg)
	if err == nil {
		t.Fatal("Expected error due to invalid integer environment variable, got nil")
	}

	// Ensure the error mentions the invalid integer value
	expectedError := "error setting field ApplicationPort: invalid integer value: strconv.Atoi: parsing \"invalid_int\": invalid syntax"
	if err.Error() != expectedError {
		t.Errorf("Expected error %s, got %s", expectedError, err.Error())
	}
}
