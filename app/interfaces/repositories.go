package interfaces

import (
	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductsRepository interface {
	GetProductByCode(code string) (*models.Product, error)
	GetAllProductsWithPagination(
		category *string,
		priceLessThen *float64,
		limit int,
		offset int,
	) (
		[]models.Product,
		error,
	)
	Count() (*int64, error)
}

type CategoriesRepository interface {
	Create(data dto.Category) (*models.Category, error)
	GetByCode(code string) (*models.Category, error)
	GetAll(limit int, offset int) ([]models.Category, error)
	Count() (*int64, error)
}
