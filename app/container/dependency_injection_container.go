package container

import (
	"log"

	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type DependencyInjectionContainer struct {
	DB                *database.Client
	ProductRepository interfaces.ProductsRepository
}

var Container *DependencyInjectionContainer

func CreateDIContainer() {
	if Container != nil {
		return
	}

	db, err := database.New(true)
	if err != nil {
		log.Fatal(err)
	}

	productsRepository := models.NewProductsRepository(db)

	Container = &DependencyInjectionContainer{
		DB:                db,
		ProductRepository: productsRepository,
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
