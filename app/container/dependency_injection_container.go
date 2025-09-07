package container

import (
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"log"
)

type DependencyInjectionContainer struct {
	DB *database.Client
}

var Container *DependencyInjectionContainer

func CreateDIContainer() {
	if Container != nil {
		return
	}

	db, err := database.New()
	if err != nil {
		log.Fatal(err)
	}

	Container = &DependencyInjectionContainer{
		DB: db,
	}
}

func ShutDown() {
	if Container == nil {
		log.Println("DependencyInjectionContainer is not initialized, skipping shutdown")
		return
	}
	Container.Shutdown()
}

func (c *DependencyInjectionContainer) Shutdown() {
	log.Println("Shutting down DependencyInjectionContainer and all its resources")

	c.DB.Shutdown()
	log.Println("MySQL database disconnected")

}
