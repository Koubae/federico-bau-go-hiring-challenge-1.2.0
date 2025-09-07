package repository

import (
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type SQLProductsRepository struct {
	db *database.Client
}

func NewProductsRepository(db *database.Client) *SQLProductsRepository {
	return &SQLProductsRepository{
		db: db,
	}
}

func (r *SQLProductsRepository) GetAllProducts(limit int, offset int) ([]models.Product, error) {
	var products []models.Product
	if err := r.db.DB.
		Preload("Category").
		Preload("Variants").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error; err != nil {
		return nil, err
	}
	return products, nil
}
