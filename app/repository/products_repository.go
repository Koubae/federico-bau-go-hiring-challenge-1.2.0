package repository

import (
	"errors"
	"log"

	"github.com/jellydator/ttlcache/v3"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"gorm.io/gorm"
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

func (r *SQLProductsRepository) GetProductByCode(code string) (*models.Product, error) {
	var product models.Product

	dbWithCtx := r.db.DB.Model(&models.Product{})

	if err := dbWithCtx.
		Preload("Category").
		Preload("Variants").
		Where("code = ?", code).
		First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.New("error while Get Product By Code")

	}
	return &product, nil
}

func (r *SQLProductsRepository) GetAllProductsWithPagination(
	category *string,
	priceLessThen *float64,
	limit int,
	offset int,
) ([]models.Product, error) {
	var products []models.Product

	dbWithCtx := r.db.DB.Model(&models.Product{})

	if category != nil && *category != "" {
		dbWithCtx = dbWithCtx.Joins("Category").Where("\"Category\".name = ?", *category)
	}
	if priceLessThen != nil {
		dbWithCtx = dbWithCtx.Where("price <= ?", *priceLessThen)
	}

	if err := dbWithCtx.
		Preload("Category").
		Preload("Variants").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&products).
		Error; err != nil {
		return nil, errors.New("error while listing Products Catalogs")
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
		return nil, errors.New("error while counting Products")
	}

	r.cache.Set(CacheProductCountKey, &count, CacheTTL)
	return &count, nil
}
