package database

import (
	"fmt"
	"go_be_enrollment/internal/config"
	"go_be_enrollment/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMySQL(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDB,
		cfg.MySQLParams,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.Error("Failed to connect to MySQL database: " + err.Error())
		return nil, err
	}

	logger.Log.Info("Connected to MySQL database successfully")
	return db, nil
}
