package repository

import "time"

const (
	CacheTTL              = 3 * time.Minute
	CacheProductCountKey  = "products__count"
	CacheCategoryCountKey = "categories__count"
)
