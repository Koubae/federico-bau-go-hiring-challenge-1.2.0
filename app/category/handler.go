package category

import (
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
)

type Handler struct {
	repository interfaces.CategoriesRepository
}

func NewCategoryHandler(r interfaces.CategoriesRepository) *Handler {
	return &Handler{
		repository: r,
	}
}

type ListCategoriesResponse struct {
	Total      int64          `json:"total"`
	Categories []dto.Category `json:"categories"`
}

func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	pagination, err := api.PaginationRequest(r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Pagination validation failure: "+err.Error())
		return
	}

	res, err := h.repository.GetAll(pagination.Limit, pagination.Offset)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred, please try again later.")
		return
	}

	total, err := h.repository.Count()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred, please try again later.")
		return
	}

	categories := make([]dto.Category, len(res))
	for i, v := range res {
		categories[i] = dto.Category{
			Code: v.Code,
			Name: v.Name,
		}
	}

	response := ListCategoriesResponse{
		Total:      *total,
		Categories: categories,
	}
	api.OKResponse(w, response)
}
