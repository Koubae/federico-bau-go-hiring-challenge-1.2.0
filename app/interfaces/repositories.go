package interfaces

import (
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductsRepository interface {
	GetAllProductsWithPagination(limit int, offset int) ([]models.Product, error)
	Count() (*int64, error)
}
