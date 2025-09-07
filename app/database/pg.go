package database

import (
	"errors"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var client *Client

func New(prepareStmt bool) (*Client, error) {
	var err error
	var db *gorm.DB
	var config *DatabaseConfig

	config = NewDatabaseConfig()

	db, err = gorm.Open(
		postgres.Open(config.Dns()), &gorm.Config{
			Logger:      logger.Default.LogMode(config.LogLevel),
			PrepareStmt: prepareStmt,
		},
	)
	if err != nil {
		return nil, errors.New("Failed to connect database, error: " + err.Error())
	}

	_, err = db.DB()
	if err != nil {
		return nil, errors.New("Failed to get database connection, error: " + err.Error())
	}

	client = &Client{
		Config: config,
		DB:     db,
	}
	client.Ping()

	log.Printf("Database connected 🌐🛜✅ %v\n", client)
	return client, nil
}

func GetClient() *Client {
	if client == nil {
		panic("PostgreSQL Client is not initialized!")
	}
	return client
}
