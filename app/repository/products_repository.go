package repository

import (
	"log"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

const (
	CacheTTL             = 3 * time.Minute
	CacheProductCountKey = "products__count"
)

type SQLProductsRepository struct {
	db    *database.Client
	cache *ttlcache.Cache[string, *int64] // cache key to count
}

func NewProductsRepository(db *database.Client) *SQLProductsRepository {
	cache := ttlcache.New[string, *int64](
		ttlcache.WithTTL[string, *int64](CacheTTL),
	)
	go cache.Start()

	return &SQLProductsRepository{
		db:    db,
		cache: cache,
	}
}

func (r *SQLProductsRepository) GetAllProductsWithPagination(limit int, offset int) ([]models.Product, error) {
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

func (r *SQLProductsRepository) Count() (*int64, error) {
	if item := r.cache.Get(CacheProductCountKey); item != nil {
		return item.Value(), nil
	}
	log.Println("No product count in cache, fetching from database")

	var count int64
	if err := r.db.DB.Model(&models.Product{}).Count(&count).Error; err != nil {
		return nil, err
	}

	r.cache.Set(CacheProductCountKey, &count, CacheTTL)
	return &count, nil
}
