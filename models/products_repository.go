package models

import (
	"github.com/mytheresa/go-hiring-challenge/app/database"
)

type ProductsRepository struct {
	db *database.Client
}

func NewProductsRepository(db *database.Client) *ProductsRepository {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.DB.Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
