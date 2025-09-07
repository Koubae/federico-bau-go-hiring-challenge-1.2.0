package database

import (
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var client *Client

func New(user, password, dbname, port string) (*Client, error) {
	var err error
	var db *gorm.DB

	// TODO : config
	// TODO : host!
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, password, port, dbname)

	db, err = gorm.Open(
		postgres.Open(dsn), &gorm.Config{
			Logger:      logger.Default.LogMode(logger.Info),
			PrepareStmt: true,
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
		DB: db,
	}
	client.Ping()

	log.Println("PostgreSQL database connected")
	return client, nil
}

func GetClient() *Client {
	if client == nil {
		panic("PostgreSQL Client is not initialized!")
	}
	return client
}
