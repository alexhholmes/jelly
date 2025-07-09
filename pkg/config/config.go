package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Photo struct {
		MaxFileSizeMB int `yaml:"max_file_size_mb" env:"PHOTO_MAX_FILE_SIZE_MB"`
	} `yaml:"photo"`
	Database struct {
		Host     string `yaml:"host" env:"DB_HOST"`
		Port     int    `yaml:"port" env:"DB_PORT"`
		Name     string `yaml:"name" env:"DB_NAME"`
		User     string `yaml:"user" env:"DB_USER"`
		Password string `yaml:"password" env:"DB_PASSWORD"`
	} `yaml:"database"`
	Server struct {
		Port         int    `yaml:"port" env:"SERVER_PORT"`
		ReadTimeout  string `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT"`
		WriteTimeout string `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT"`
	} `yaml:"server"`
	Storage struct {
		Type      string `yaml:"type" env:"STORAGE_TYPE"`
		LocalPath string `yaml:"local_path" env:"STORAGE_LOCAL_PATH"`
		S3Bucket  string `yaml:"s3_bucket" env:"STORAGE_S3_BUCKET"`
		S3Region  string `yaml:"s3_region" env:"STORAGE_S3_REGION"`
	} `yaml:"storage"`
}

// Load reads configuration from config.yaml and sets environment variables
// Environment variables override config file values
func Load() (*Config, error) {
	config := &Config{}

	// Try to find config.yaml in current directory or config/ directory
	configPaths := []string{
		"config.yaml",
		"config/config.yaml",
		"cfg/config.yaml",
	}

	var configFile string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFile = path
			break
		}
	}

	// If config file exists, load it
	if configFile != "" {
		data, err := os.ReadFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file %s: %w", configFile, err)
		}
	}

	// Set environment variables from config, but don't override existing ones
	setEnvFromConfig(config)

	return config, nil
}

// setEnvFromConfig sets environment variables from config values
// Only sets if the environment variable doesn't already exist
func setEnvFromConfig(config *Config) {
	setEnvFromStruct(reflect.ValueOf(config).Elem())
}

// setEnvFromStruct recursively processes struct fields and sets environment variables
func setEnvFromStruct(v reflect.Value) {
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Get the env tag
		envTag := fieldType.Tag.Get("env")

		// If field is a struct, recurse into it
		if field.Kind() == reflect.Struct {
			setEnvFromStruct(field)
			continue
		}

		// Skip if no env tag or field is not settable
		if envTag == "" || !field.CanInterface() {
			continue
		}

		// Convert field value to string based on type
		var value string
		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Int:
			if field.Int() != 0 {
				value = strconv.FormatInt(field.Int(), 10)
			}
		}

		// User set env var override config file
		if value != "" && os.Getenv(envTag) == "" {
			os.Setenv(envTag, value)
		}
	}
}

// GetPhotoMaxFileSizeBytes returns the maximum file size in bytes from environment variable
func GetPhotoMaxFileSizeBytes() int64 {
	valueStr := os.Getenv("PHOTO_MAX_FILE_SIZE_MB")
	if valueStr == "" {
		// Default to 10MB if not set
		valueStr = "10"
	}

	maxSizeMB, err := strconv.Atoi(valueStr)
	if err != nil {
		// If conversion fails, log error and return default size
		fmt.Printf("Invalid PHOTO_MAX_FILE_SIZE_MB value: %s, using default 10MB\n", valueStr)
		maxSizeMB = 10
	}

	return int64(maxSizeMB) << 20 // Convert MB to bytes
}
