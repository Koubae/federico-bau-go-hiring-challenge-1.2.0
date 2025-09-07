package database

import (
	"database/sql"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Client struct {
	Config *DatabaseConfig
	DB     *gorm.DB
}

func (c *Client) String() string {
	return fmt.Sprintf("PostgreSQL -- %v", c.Config)
}

func (c *Client) GetConnection() *sql.DB {
	database, err := c.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection, error: %v\n", err)
	}
	return database
}

func (c *Client) Shutdown() {
	database := c.GetConnection()
	if err := database.Close(); err != nil {
		log.Printf("Failed to shutdown PostgreSQL: %v\n", err.Error())
	}
}

func (c *Client) Ping() {
	database := c.GetConnection()
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v\n", err.Error())
	}
}
