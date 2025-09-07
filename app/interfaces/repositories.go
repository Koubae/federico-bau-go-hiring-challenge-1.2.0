package interfaces

import (
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
