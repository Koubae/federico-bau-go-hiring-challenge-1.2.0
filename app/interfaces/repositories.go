package interfaces

import (
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductsRepository interface {
	GetAllProducts(limit int, offset int) ([]models.Product, error)
}
