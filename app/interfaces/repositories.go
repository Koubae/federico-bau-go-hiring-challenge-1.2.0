package interfaces

import (
	"github.com/mytheresa/go-hiring-challenge/models"
)

type ProductsRepository interface {
	GetAllProducts() ([]models.Product, error)
}
