package catalog

import (
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
)

type Response struct {
	Total    int64     `json:"total"`
	Products []Product `json:"products"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Product struct {
	Code     string   `json:"code"`
	Price    float64  `json:"price"`
	Category Category `json:"category"`
}

type CatalogHandler struct {
	repository interfaces.ProductsRepository
}

func NewCatalogHandler(r interfaces.ProductsRepository) *CatalogHandler {
	return &CatalogHandler{
		repository: r,
	}
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

	response := Response{
		Total:    *total,
		Products: products,
	}
	api.OKResponse(w, response)

}
