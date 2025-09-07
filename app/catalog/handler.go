package catalog

import (
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
)

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

type CatalogHandler struct {
	repository interfaces.ProductsRepository
}

func NewCatalogHandler(r interfaces.ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repository: r,
	}
}

type ListCatalogResponse struct {
	Total    int64     `json:"total"`
	Products []Product `json:"products"`
}

func (h *CatalogHandler) ListCatalog(w http.ResponseWriter, r *http.Request) {
	pagination, err := api.PaginationRequest(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Pagination validation failure: "+err.Error())
		return
	}

	q := r.URL.Query()

	var category *string
	if cat := q.Get("category"); cat != "" {
		category = &cat
	}
	var priceLessThen *float64
	if p := q.Get("priceLessThen"); p != "" {
		if val, err := strconv.ParseFloat(p, 64); err == nil {
			priceLessThen = &val
		}
	}

	res, err := h.repository.GetAllProductsWithPagination(category, priceLessThen, pagination.Limit, pagination.Offset)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred, please try again later.")
		return
	}

	total, err := h.repository.Count()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred, please try again later.")
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: Category{
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
		}
	}

	response := ListCatalogResponse{
		Total:    *total,
		Products: products,
	}
	api.OKResponse(w, response)

}

type GetProductDetailsResponse struct {
	ProductWithDetails
}

func (h *CatalogHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	res, err := h.repository.GetProductByCode(code)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred, please try again later.")
		return
	}
	if res == nil {
		api.ErrorResponse(w, http.StatusNotFound, "Product not found")
		return
	}

	basePrice := res.Price.InexactFloat64()
	variants := make([]ProductVariant, len(res.Variants))
	for i, v := range res.Variants {
		variantPrice := v.Price.InexactFloat64()
		if variantPrice == 0 {
			variantPrice = basePrice
		}

		variants[i] = ProductVariant{
			ID:    v.ID,
			SKU:   v.SKU,
			Name:  v.Name,
			Price: variantPrice,
		}
	}

	response := GetProductDetailsResponse{
		ProductWithDetails{
			ID: res.ID,
			Product: Product{
				Code:  res.Code,
				Price: res.Price.InexactFloat64(),
				Category: Category{
					Code: res.Category.Code,
					Name: res.Category.Name,
				},
			},
			ProductVariant: variants,
		},
	}
	api.OKResponse(w, response)

}
