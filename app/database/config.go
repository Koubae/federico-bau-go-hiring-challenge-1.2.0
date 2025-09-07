package database

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	User     string
	DBName   string
	Host     string
	Port     string
	SSLMode  string
	LogLevel logger.LogLevel

	password string
}

func (c *DatabaseConfig) String() string {
	return fmt.Sprintf("Database Connected @%s:%s as user '%s' on db '%s'", c.Host, c.Port, c.User, c.DBName)
}

func (c *DatabaseConfig) Dns() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.password, c.Host, c.Port, c.DBName, c.SSLMode,
	)
}

var databaseConfig *DatabaseConfig

func NewDatabaseConfig() *DatabaseConfig {
	// TODO: better way to handle how we load env variables!"
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	sslMode := os.Getenv("POSTGRES_SSL_MODE")

	logLevel := os.Getenv("POSTGRES_LOG_LEVEL")

	databaseConfig = &DatabaseConfig{
		User:     user,
		DBName:   dbname,
		Host:     host,
		Port:     port,
		SSLMode:  sslMode,
		LogLevel: stringToGormLogLevel(logLevel),
		password: password,
	}
	return databaseConfig
}

func stringToGormLogLevel(level string) logger.LogLevel {
	switch strings.ToLower(level) {
	case "silent":
		return logger.Silent
	case "info":
		return logger.Info
	case "warn", "warning":
		return logger.Warn
	case "error":
		return logger.Error
	default:
		return logger.Warn
	}
}
