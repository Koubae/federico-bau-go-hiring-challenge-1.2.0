package models

import (
	"github.com/mytheresa/go-hiring-challenge/app/database"
)

type SQLProductsRepository struct {
	db *database.Client
}

func NewProductsRepository(db *database.Client) *SQLProductsRepository {
	return &SQLProductsRepository{
		db: db,
	}
}

func (r *SQLProductsRepository) GetAllProducts() ([]Product, error) {
	var products []Product
	if err := r.db.DB.Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
