package category

import (
	"net/http"
	"strings"

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

type CreateCategoryResponse struct {
	ID int `json:"id"`
	dto.Category
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	payload, err := api.ParseRequestBodyJSON[dto.Category](r)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if strings.TrimSpace(payload.Name) == "" || strings.TrimSpace(payload.Code) == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid payload, missing required fields")
		return
	}

	category, err := h.repository.GetByCode(payload.Code)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred, please try again later.")
		return
	} else if category != nil {
		api.ErrorResponse(w, http.StatusConflict, "Category already exists")
		return
	}

	res, err := h.repository.Create(payload)
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Unexpected error occurred, please try again later.")
		return
	}

	response := CreateCategoryResponse{
		ID: res.ID,
		Category: dto.Category{
			Name: res.Name,
			Code: res.Code,
		},
	}

	api.OKResponseWithStatus(w, response, http.StatusCreated)
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
