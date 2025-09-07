package category

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
	"github.com/mytheresa/go-hiring-challenge/app/models"
)

// mockCategoriesRepository is a mock implementation of interfaces.CategoriesRepository
type mockCategoriesRepository struct {
	categories      []models.Category
	count           int64
	createError     error
	getByCodeError  error
	getAllError     error
	countError      error
	categoryByCode  map[string]*models.Category
	shouldReturnNil bool
}

func (m *mockCategoriesRepository) Create(data dto.Category) (*models.Category, error) {
	if m.createError != nil {
		return nil, m.createError
	}

	// Simulate ID generation
	newID := len(m.categories) + 1
	category := &models.Category{
		ID:   newID,
		Code: data.Code,
		Name: data.Name,
	}

	m.categories = append(m.categories, *category)
	return category, nil
}

func (m *mockCategoriesRepository) GetByCode(code string) (*models.Category, error) {
	if m.getByCodeError != nil {
		return nil, m.getByCodeError
	}

	if m.shouldReturnNil {
		return nil, nil
	}

	if m.categoryByCode != nil {
		return m.categoryByCode[code], nil
	}

	for _, category := range m.categories {
		if category.Code == code {
			return &category, nil
		}
	}

	return nil, nil
}

func (m *mockCategoriesRepository) GetAll(limit int, offset int) ([]models.Category, error) {
	if m.getAllError != nil {
		return nil, m.getAllError
	}

	start := offset
	end := offset + limit

	if start >= len(m.categories) {
		return []models.Category{}, nil
	}

	if end > len(m.categories) {
		end = len(m.categories)
	}

	return m.categories[start:end], nil
}

func (m *mockCategoriesRepository) Count() (*int64, error) {
	if m.countError != nil {
		return nil, m.countError
	}
	return &m.count, nil
}

func TestHandler_CreateCategory(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		repository     *mockCategoriesRepository
		expectedStatus int
		expectedBody   interface{}
		checkBody      bool
	}{
		{
			name:        "successful category creation",
			requestBody: `{"code":"CLOTHING","name":"Clothing"}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true, // No existing category
			},
			expectedStatus: http.StatusCreated,
			expectedBody: CreateCategoryResponse{
				ID: 1,
				Category: dto.Category{
					Code: "CLOTHING",
					Name: "Clothing",
				},
			},
			checkBody: true,
		},
		{
			name:        "invalid JSON format",
			requestBody: `{"code":"CLOTHING","name":}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true,
			},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "missing required field - code",
			requestBody: `{"name":"Clothing"}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true,
			},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "missing required field - name",
			requestBody: `{"code":"CLOTHING"}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true,
			},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "empty code field",
			requestBody: `{"code":"","name":"Clothing"}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true,
			},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "empty name field",
			requestBody: `{"code":"CLOTHING","name":""}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true,
			},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "whitespace only fields",
			requestBody: `{"code":"   ","name":"   "}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true,
			},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "category already exists",
			requestBody: `{"code":"EXISTING","name":"Existing Category"}`,
			repository: &mockCategoriesRepository{
				categoryByCode: map[string]*models.Category{
					"EXISTING": {
						ID:   1,
						Code: "EXISTING",
						Name: "Existing Category",
					},
				},
			},
			expectedStatus: http.StatusConflict,
			checkBody:      false,
		},
		{
			name:        "repository GetByCode error",
			requestBody: `{"code":"CLOTHING","name":"Clothing"}`,
			repository: &mockCategoriesRepository{
				getByCodeError: errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody:      false,
		},
		{
			name:        "repository Create error",
			requestBody: `{"code":"CLOTHING","name":"Clothing"}`,
			repository: &mockCategoriesRepository{
				shouldReturnNil: true,
				createError:     errors.New("create error"),
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody:      false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handler := &Handler{
					repository: tt.repository,
				}

				req := httptest.NewRequest("POST", "/categories", strings.NewReader(tt.requestBody))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				handler.CreateCategory(rr, req)

				if rr.Code != tt.expectedStatus {
					t.Errorf("CreateCategory() status = %v, want %v", rr.Code, tt.expectedStatus)
				}

				if tt.checkBody && tt.expectedStatus == http.StatusCreated {
					var response CreateCategoryResponse
					err := json.Unmarshal(rr.Body.Bytes(), &response)
					if err != nil {
						t.Fatalf("Failed to unmarshal response: %v", err)
					}

					expectedResponse := tt.expectedBody.(CreateCategoryResponse)
					if !reflect.DeepEqual(response, expectedResponse) {
						t.Errorf("CreateCategory() response = %v, want %v", response, expectedResponse)
					}
				}

				if tt.expectedStatus != http.StatusCreated && !tt.checkBody {
					var errorResponse map[string]interface{}
					err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
					if err != nil {
						t.Errorf("Expected JSON error response but got: %s", rr.Body.String())
					}
				}

				// Verify Content-Type header
				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", contentType)
				}
			},
		)
	}
}

func TestHandler_ListCategories(t *testing.T) {
	mockCategories := []models.Category{
		{
			ID:   1,
			Code: "CLOTHING",
			Name: "Clothing",
		},
		{
			ID:   2,
			Code: "SHOES",
			Name: "Shoes",
		},
		{
			ID:   3,
			Code: "ACCESSORIES",
			Name: "Accessories",
		},
	}

	tests := []struct {
		name           string
		queryParams    string
		repository     *mockCategoriesRepository
		expectedStatus int
		expectedBody   interface{}
		checkBody      bool
	}{
		{
			name:        "successful list with default pagination",
			queryParams: "",
			repository: &mockCategoriesRepository{
				categories: mockCategories,
				count:      3,
			},
			expectedStatus: http.StatusOK,
			expectedBody: ListCategoriesResponse{
				Total: 3,
				Categories: []dto.Category{
					{Code: "CLOTHING", Name: "Clothing"},
					{Code: "SHOES", Name: "Shoes"},
					{Code: "ACCESSORIES", Name: "Accessories"},
				},
			},
			checkBody: true,
		},
		{
			name:        "successful list with custom pagination",
			queryParams: "?limit=2&offset=1",
			repository: &mockCategoriesRepository{
				categories: mockCategories,
				count:      3,
			},
			expectedStatus: http.StatusOK,
			expectedBody: ListCategoriesResponse{
				Total: 3,
				Categories: []dto.Category{
					{Code: "SHOES", Name: "Shoes"},
					{Code: "ACCESSORIES", Name: "Accessories"},
				},
			},
			checkBody: true,
		},
		{
			name:        "empty categories list",
			queryParams: "",
			repository: &mockCategoriesRepository{
				categories: []models.Category{},
				count:      0,
			},
			expectedStatus: http.StatusOK,
			expectedBody: ListCategoriesResponse{
				Total:      0,
				Categories: []dto.Category{},
			},
			checkBody: true,
		},
		{
			name:           "invalid limit parameter",
			queryParams:    "?limit=invalid",
			repository:     &mockCategoriesRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "invalid offset parameter",
			queryParams:    "?offset=invalid",
			repository:     &mockCategoriesRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "limit exceeds maximum",
			queryParams:    "?limit=101",
			repository:     &mockCategoriesRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "zero limit",
			queryParams:    "?limit=0",
			repository:     &mockCategoriesRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "negative offset",
			queryParams:    "?offset=-1",
			repository:     &mockCategoriesRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "repository GetAll error",
			queryParams: "",
			repository: &mockCategoriesRepository{
				getAllError: errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody:      false,
		},
		{
			name:        "repository Count error",
			queryParams: "",
			repository: &mockCategoriesRepository{
				categories: mockCategories,
				countError: errors.New("count error"),
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody:      false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handler := &Handler{
					repository: tt.repository,
				}

				req := httptest.NewRequest("GET", "/categories"+tt.queryParams, nil)
				rr := httptest.NewRecorder()

				handler.ListCategories(rr, req)

				if rr.Code != tt.expectedStatus {
					t.Errorf("ListCategories() status = %v, want %v", rr.Code, tt.expectedStatus)
				}

				if tt.checkBody && tt.expectedStatus == http.StatusOK {
					var response ListCategoriesResponse
					err := json.Unmarshal(rr.Body.Bytes(), &response)
					if err != nil {
						t.Fatalf("Failed to unmarshal response: %v", err)
					}

					expectedResponse := tt.expectedBody.(ListCategoriesResponse)
					if !reflect.DeepEqual(response, expectedResponse) {
						t.Errorf("ListCategories() response = %v, want %v", response, expectedResponse)
					}
				}

				if tt.expectedStatus != http.StatusOK && !tt.checkBody {
					var errorResponse map[string]interface{}
					err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
					if err != nil {
						t.Errorf("Expected JSON error response but got: %s", rr.Body.String())
					}
				}

				// Verify Content-Type header
				contentType := rr.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", contentType)
				}
			},
		)
	}
}

func TestNewCategoryHandler(t *testing.T) {
	tests := []struct {
		name string
		args interfaces.CategoriesRepository
		want *Handler
	}{
		{
			name: "creates new category handler with valid repository",
			args: &mockCategoriesRepository{},
			want: &Handler{repository: &mockCategoriesRepository{}},
		},
		{
			name: "creates new category handler with nil repository",
			args: nil,
			want: &Handler{repository: nil},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewCategoryHandler(tt.args)
				if !reflect.DeepEqual(got.repository, tt.want.repository) {
					t.Errorf("NewCategoryHandler() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

// Integration tests
func TestHandler_CreateCategory_Integration(t *testing.T) {
	t.Run(
		"complete workflow - create and verify response structure", func(t *testing.T) {
			mockRepo := &mockCategoriesRepository{
				shouldReturnNil: true,
			}

			handler := &Handler{repository: mockRepo}

			requestBody := `{"code":"INTEGRATION","name":"Integration Test Category"}`
			req := httptest.NewRequest("POST", "/categories", strings.NewReader(requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.CreateCategory(rr, req)

			if rr.Code != http.StatusCreated {
				t.Fatalf("Expected status 201, got %d", rr.Code)
			}

			var response CreateCategoryResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.ID != 1 {
				t.Errorf("Expected ID 1, got %d", response.ID)
			}

			if response.Code != "INTEGRATION" {
				t.Errorf("Expected code INTEGRATION, got %s", response.Code)
			}

			if response.Name != "Integration Test Category" {
				t.Errorf("Expected name 'Integration Test Category', got %s", response.Name)
			}
		},
	)
}

func TestHandler_ListCategories_Integration(t *testing.T) {
	mockCategories := []models.Category{
		{ID: 1, Code: "CAT1", Name: "Category 1"},
		{ID: 2, Code: "CAT2", Name: "Category 2"},
		{ID: 3, Code: "CAT3", Name: "Category 3"},
		{ID: 4, Code: "CAT4", Name: "Category 4"},
		{ID: 5, Code: "CAT5", Name: "Category 5"},
	}

	mockRepo := &mockCategoriesRepository{
		categories: mockCategories,
		count:      5,
	}

	handler := &Handler{repository: mockRepo}

	t.Run(
		"pagination logic verification", func(t *testing.T) {
			tests := []struct {
				name          string
				queryParams   string
				expectedCount int
				expectedFirst string
				expectedLast  string
			}{
				{
					name:          "first page default limit",
					queryParams:   "",
					expectedCount: 5,
					expectedFirst: "CAT1",
					expectedLast:  "CAT5",
				},
				{
					name:          "custom limit",
					queryParams:   "?limit=3",
					expectedCount: 3,
					expectedFirst: "CAT1",
					expectedLast:  "CAT3",
				},
				{
					name:          "custom offset",
					queryParams:   "?offset=2",
					expectedCount: 3,
					expectedFirst: "CAT3",
					expectedLast:  "CAT5",
				},
				{
					name:          "custom limit and offset",
					queryParams:   "?limit=2&offset=3",
					expectedCount: 2,
					expectedFirst: "CAT4",
					expectedLast:  "CAT5",
				},
				{
					name:          "offset beyond data",
					queryParams:   "?offset=10",
					expectedCount: 0,
					expectedFirst: "",
					expectedLast:  "",
				},
			}

			for _, tt := range tests {
				t.Run(
					tt.name, func(t *testing.T) {
						req := httptest.NewRequest("GET", "/categories"+tt.queryParams, nil)
						rr := httptest.NewRecorder()

						handler.ListCategories(rr, req)

						if rr.Code != http.StatusOK {
							t.Fatalf("Expected status 200, got %d", rr.Code)
						}

						var response ListCategoriesResponse
						err := json.Unmarshal(rr.Body.Bytes(), &response)
						if err != nil {
							t.Fatalf("Failed to unmarshal response: %v", err)
						}

						if len(response.Categories) != tt.expectedCount {
							t.Errorf("Expected %d categories, got %d", tt.expectedCount, len(response.Categories))
						}

						if tt.expectedCount > 0 {
							if response.Categories[0].Code != tt.expectedFirst {
								t.Errorf(
									"Expected first category %s, got %s",
									tt.expectedFirst,
									response.Categories[0].Code,
								)
							}

							if response.Categories[len(response.Categories)-1].Code != tt.expectedLast {
								t.Errorf(
									"Expected last category %s, got %s",
									tt.expectedLast,
									response.Categories[len(response.Categories)-1].Code,
								)
							}
						}

						if response.Total != 5 {
							t.Errorf("Expected total 5, got %d", response.Total)
						}
					},
				)
			}
		},
	)
}

func TestHandler_CreateCategory_EdgeCases(t *testing.T) {
	t.Run(
		"valid category with special characters", func(t *testing.T) {
			mockRepo := &mockCategoriesRepository{
				shouldReturnNil: true,
			}

			handler := &Handler{repository: mockRepo}

			requestBody := `{"code":"SPECIAL-CODE_123","name":"Special Category & More!"}`
			req := httptest.NewRequest("POST", "/categories", bytes.NewReader([]byte(requestBody)))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.CreateCategory(rr, req)

			if rr.Code != http.StatusCreated {
				t.Fatalf("Expected status 201, got %d", rr.Code)
			}

			var response CreateCategoryResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Code != "SPECIAL-CODE_123" {
				t.Errorf("Expected code 'SPECIAL-CODE_123', got %s", response.Code)
			}

			if response.Name != "Special Category & More!" {
				t.Errorf("Expected name 'Special Category & More!', got %s", response.Name)
			}
		},
	)

	t.Run(
		"empty request body", func(t *testing.T) {
			mockRepo := &mockCategoriesRepository{
				shouldReturnNil: true,
			}

			handler := &Handler{repository: mockRepo}

			req := httptest.NewRequest("POST", "/categories", strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.CreateCategory(rr, req)

			if rr.Code != http.StatusBadRequest {
				t.Fatalf("Expected status 400, got %d", rr.Code)
			}
		},
	)
}
