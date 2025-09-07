package models

// Category represents a category in the catalog.
// It includes a unique code and a category name.
type Category struct {
	ID   int    `gorm:"primaryKey"`
	Code string `gorm:"uniqueIndex;size:50;not null;"`
	Name string `gorm:"index:size:255;not null"`

	Products []Product `gorm:"foreignKey:CategoryID"`
}

func (p *Category) TableName() string {
	return "categories"
}
