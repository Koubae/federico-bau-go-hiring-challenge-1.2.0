package repository

import (
	"errors"
	"log"

	"github.com/jellydator/ttlcache/v3"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

type SQLCategoryRepository struct {
	db    *database.Client
	cache *ttlcache.Cache[string, *int64] // cache key to count
}

func NewCategoryRepository(db *database.Client) *SQLCategoryRepository {
	cache := ttlcache.New[string, *int64](
		ttlcache.WithTTL[string, *int64](CacheTTL),
	)
	go cache.Start()

	return &SQLCategoryRepository{
		db:    db,
		cache: cache,
	}
}

func (r *SQLCategoryRepository) Create(data dto.Category) (*models.Category, error) {
	category := &models.Category{
		Code: data.Code,
		Name: data.Name,
	}
	err := r.db.DB.Create(&category).Error
	if err != nil {
		log.Printf("Error while creating Category %v+, error: %v\n", category, err)
		return nil, errors.New("error while creating Category")
	}

	log.Printf("New Category created: %v+\n", category)
	return category, nil

}

func (r *SQLCategoryRepository) GetAll(limit int, offset int) ([]models.Category, error) {
	var categories []models.Category

	dbWithCtx := r.db.DB.Model(&models.Category{})
	if err := dbWithCtx.
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&categories).
		Error; err != nil {
		return nil, errors.New("error while listing Categories")
	}
	return categories, nil
}

func (r *SQLCategoryRepository) Count() (*int64, error) {
	if item := r.cache.Get(CacheCategoryCountKey); item != nil {
		return item.Value(), nil
	}
	log.Println("No Category count in cache, fetching from database")

	var count int64
	if err := r.db.DB.Model(&models.Category{}).Count(&count).Error; err != nil {
		return nil, errors.New("error while counting Categories")
	}

	r.cache.Set(CacheCategoryCountKey, &count, CacheTTL)
	return &count, nil
}
