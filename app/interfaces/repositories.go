package interfaces

import (
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type ProductsRepository interface {
	GetAllProducts() ([]models.Product, error)
}
