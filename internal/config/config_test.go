package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadConfigSuccess(t *testing.T) {
	os.Clearenv()
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_DB", "test_db")
	os.Setenv("MYSQL_USER", "tester")
	os.Setenv("MYSQL_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "supersecret")
	os.Setenv("JWT_EXPIRES_IN", "24h")
	os.Setenv("APP_PORT", "8000") // override default

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Expected nil err, got %v", err)
	}

	if cfg.AppPort != "8000" {
		t.Errorf("Expected APP_PORT 8000, got %s", cfg.AppPort)
	}

	if cfg.JWTExpiresIn != 24*time.Hour {
		t.Errorf("Expected 24h duration, got %v", cfg.JWTExpiresIn)
	}

	if cfg.AppEnv != "development" {
		t.Errorf("Expected fallback APP_ENV development, got %s", cfg.AppEnv)
	}
}

func TestLoadConfigMissingRequired(t *testing.T) {
	os.Clearenv()
	// Missing MYSQL_HOST intentionally
	os.Setenv("MYSQL_DB", "test_db")
	os.Setenv("MYSQL_USER", "tester")
	os.Setenv("MYSQL_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "supersecret")
	os.Setenv("JWT_EXPIRES_IN", "2h")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error due to missing MYSQL_HOST, got successful config load")
	}
}

func TestInvalidDuration(t *testing.T) {
	os.Clearenv()
	os.Setenv("MYSQL_HOST", "127.0.0.1")
	os.Setenv("MYSQL_DB", "test_db")
	os.Setenv("MYSQL_USER", "tester")
	os.Setenv("MYSQL_PASSWORD", "testpass")
	os.Setenv("JWT_SECRET", "supersecret")
	os.Setenv("JWT_EXPIRES_IN", "invalid_time_format")

	_, err := Load()
	if err == nil {
		t.Fatal("Expected error due to invalid duration, got successful config load")
	}
}
