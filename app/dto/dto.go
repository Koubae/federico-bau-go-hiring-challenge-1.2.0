package dto

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Product struct {
	Code     string   `json:"code"`
	Price    float64  `json:"price"`
	Category Category `json:"category"`
}

type ProductVariant struct {
	ID    uint    `json:"id"`
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductWithDetails struct {
	ID uint `json:"id"`
	Product
	ProductVariant []ProductVariant `json:"variants"`
}
