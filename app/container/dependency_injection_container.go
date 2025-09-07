package container

import (
	"log"

	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
	"github.com/mytheresa/go-hiring-challenge/app/repository"
)

type DependencyInjectionContainer struct {
	DB                 *database.Client
	CategoryRepository interfaces.CategoriesRepository
	ProductRepository  interfaces.ProductsRepository
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

	categoryRepository := repository.NewCategoryRepository(db)
	productsRepository := repository.NewProductsRepository(db)

	Container = &DependencyInjectionContainer{
		DB:                 db,
		CategoryRepository: categoryRepository,
		ProductRepository:  productsRepository,
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
