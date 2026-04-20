package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName           string
	AppEnv            string
	AppPort           string
	AppTimezone       string
	MySQLHost         string
	MySQLPort         string
	MySQLDB           string
	MySQLUser         string
	MySQLPassword     string
	MySQLParams       string
	JWTSecret         string
	JWTExpiresIn      time.Duration
	CORSAllowOrigins  string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2Endpoint        string
	R2PublicBaseURL   string
}

// Load reads config from environment variables and validates required fields.
func Load() (*Config, error) {
	_ = godotenv.Load() // Ignore error on missing .env file

	cfg := &Config{
		AppName:           getEnv("APP_NAME", "Enrollment System"),
		AppEnv:            getEnv("APP_ENV", "development"),
		AppPort:           getEnv("APP_PORT", "8080"),
		AppTimezone:       getEnv("APP_TIMEZONE", "Asia/Ho_Chi_Minh"),
		MySQLPort:         getEnv("MYSQL_PORT", "3306"),
		MySQLParams:       getEnv("MYSQL_PARAMS", "charset=utf8mb4&parseTime=True&loc=Local"),
		CORSAllowOrigins:  getEnv("CORS_ALLOW_ORIGINS", "*"),
		R2AccountID:       getEnv("R2_ACCOUNT_ID", ""),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
		R2BucketName:      getEnv("R2_BUCKET_NAME", ""),
		R2Endpoint:        getEnv("R2_ENDPOINT", ""),
		R2PublicBaseURL:   getEnv("R2_PUBLIC_BASE_URL", ""),
	}

	var err error

	cfg.MySQLHost, err = getEnvRequired("MYSQL_HOST")
	if err != nil {
		return nil, err
	}

	cfg.MySQLDB, err = getEnvRequired("MYSQL_DB")
	if err != nil {
		return nil, err
	}

	cfg.MySQLUser, err = getEnvRequired("MYSQL_USER")
	if err != nil {
		return nil, err
	}

	cfg.MySQLPassword, err = getEnvRequired("MYSQL_PASSWORD")
	if err != nil {
		return nil, err
	}

	cfg.JWTSecret, err = getEnvRequired("JWT_SECRET")
	if err != nil {
		return nil, err
	}

	jwtExpiresStr, err := getEnvRequired("JWT_EXPIRES_IN")
	if err != nil {
		return nil, err
	}

	cfg.JWTExpiresIn, err = parseDuration(jwtExpiresStr)
	if err != nil {
		return nil, fmt.Errorf("invalid duration for JWT_EXPIRES_IN: %w", err)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvRequired(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return val, nil
}

func parseDuration(val string) (time.Duration, error) {
	return time.ParseDuration(val)
}
