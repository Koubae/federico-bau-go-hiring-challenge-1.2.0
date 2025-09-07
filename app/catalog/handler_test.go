package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/shopspring/decimal"
)

// mockProductsRepository is a mock implementation of interfaces.ProductsRepository
type mockProductsRepository struct {
	products   []models.Product
	count      int64
	getError   error
	countError error
}

func (m *mockProductsRepository) GetAllProductsWithPagination(
	category *string,
	priceLessThen *float64,
	limit int,
	offset int,
) ([]models.Product, error) {
	if m.getError != nil {
		return nil, m.getError
	}

	start := offset
	end := offset + limit

	if start >= len(m.products) {
		return []models.Product{}, nil
	}

	if end > len(m.products) {
		end = len(m.products)
	}

	return m.products[start:end], nil
}

func (m *mockProductsRepository) Count() (*int64, error) {
	if m.countError != nil {
		return nil, m.countError
	}
	return &m.count, nil
}

func TestNewCatalogHandler(t *testing.T) {
	type args struct {
		r interfaces.ProductsRepository
	}

	mockRepo := &mockProductsRepository{}

	tests := []struct {
		name string
		args args
		want *CatalogHandler
	}{
		{
			name: "creates new catalog handler with valid repository",
			args: args{r: mockRepo},
			want: &CatalogHandler{repository: mockRepo},
		},
		{
			name: "creates new catalog handler with nil repository",
			args: args{r: nil},
			want: &CatalogHandler{repository: nil},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewCatalogHandler(tt.args.r)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewCatalogHandler() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}

func TestCatalogHandler_ListCatalog(t *testing.T) {
	// Test data setup
	mockProducts := []models.Product{
		{
			Code:  "PROD001",
			Price: decimal.NewFromFloat(100.50),
			Category: models.Category{
				Code: "CLOTHING",
				Name: "Clothing",
			},
		},
		{
			Code:  "PROD002",
			Price: decimal.NewFromFloat(200.75),
			Category: models.Category{
				Code: "SHOES",
				Name: "Shoes",
			},
		},
	}

	tests := []struct {
		name           string
		queryParams    string
		repository     *mockProductsRepository
		expectedStatus int
		expectedBody   interface{}
		checkBody      bool
	}{
		{
			name:        "successful get with default pagination",
			queryParams: "",
			repository: &mockProductsRepository{
				products: mockProducts,
				count:    2,
			},
			expectedStatus: http.StatusOK,
			expectedBody: Response{
				Total: 2,
				Products: []Product{
					{
						Code:  "PROD001",
						Price: 100.50,
						Category: Category{
							Code: "CLOTHING",
							Name: "Clothing",
						},
					},
					{
						Code:  "PROD002",
						Price: 200.75,
						Category: Category{
							Code: "SHOES",
							Name: "Shoes",
						},
					},
				},
			},
			checkBody: true,
		},
		{
			name:        "successful get with custom pagination",
			queryParams: "?limit=1&offset=1",
			repository: &mockProductsRepository{
				products: mockProducts,
				count:    2,
			},
			expectedStatus: http.StatusOK,
			expectedBody: Response{
				Total: 2,
				Products: []Product{
					{
						Code:  "PROD002",
						Price: 200.75,
						Category: Category{
							Code: "SHOES",
							Name: "Shoes",
						},
					},
				},
			},
			checkBody: true,
		},
		{
			name:           "invalid limit parameter",
			queryParams:    "?limit=invalid",
			repository:     &mockProductsRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "invalid offset parameter",
			queryParams:    "?offset=invalid",
			repository:     &mockProductsRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:        "repository get error",
			queryParams: "",
			repository: &mockProductsRepository{
				getError: errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody:      false,
		},
		{
			name:        "repository count error",
			queryParams: "",
			repository: &mockProductsRepository{
				products:   mockProducts,
				countError: errors.New("count error"),
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody:      false,
		},
		{
			name:        "empty product list",
			queryParams: "",
			repository: &mockProductsRepository{
				products: []models.Product{},
				count:    0,
			},
			expectedStatus: http.StatusOK,
			expectedBody: Response{
				Total:    0,
				Products: []Product{},
			},
			checkBody: true,
		},
		{
			name:           "limit exceeds maximum",
			queryParams:    "?limit=101",
			repository:     &mockProductsRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "zero limit",
			queryParams:    "?limit=0",
			repository:     &mockProductsRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
		{
			name:           "negative offset",
			queryParams:    "?offset=-1",
			repository:     &mockProductsRepository{},
			expectedStatus: http.StatusBadRequest,
			checkBody:      false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// Create handler with mock repository
				handler := NewCatalogHandler(tt.repository)

				// Create request
				req := httptest.NewRequest("GET", "/catalog"+tt.queryParams, nil)

				// Create response recorder
				rr := httptest.NewRecorder()

				// Execute handler
				handler.ListCatalog(rr, req)

				// Check status code
				if rr.Code != tt.expectedStatus {
					t.Errorf("ListCatalog() status = %v, want %v", rr.Code, tt.expectedStatus)
				}

				// Check response body if needed
				if tt.checkBody && tt.expectedStatus == http.StatusOK {
					var response Response
					err := json.Unmarshal(rr.Body.Bytes(), &response)
					if err != nil {
						t.Fatalf("Failed to unmarshal response: %v", err)
					}

					expectedResponse := tt.expectedBody.(Response)
					if !reflect.DeepEqual(response, expectedResponse) {
						t.Errorf("ListCatalog() response = %v, want %v", response, expectedResponse)
					}
				}

				// For error cases, check that response contains error structure
				if tt.expectedStatus != http.StatusOK && tt.checkBody == false {
					// Verify that response body contains error information
					var errorResponse map[string]interface{}
					err := json.Unmarshal(rr.Body.Bytes(), &errorResponse)
					if err != nil {
						t.Errorf("Expected JSON error response but got: %s", rr.Body.String())
					}
				}
			},
		)
	}
}

func TestCatalogHandler_ListCatalog_Integration(t *testing.T) {
	// This test focuses on the integration aspects and edge cases
	mockRepo := &mockProductsRepository{
		products: []models.Product{
			{
				Code:  "TEST001",
				Price: decimal.NewFromFloat(99.99),
				Category: models.Category{
					Code: "TEST",
					Name: "Test Category",
				},
			},
		},
		count: 1,
	}

	handler := NewCatalogHandler(mockRepo)

	t.Run(
		"correct content type", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/catalog", nil)
			rr := httptest.NewRecorder()

			handler.ListCatalog(rr, req)

			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}
		},
	)

	t.Run(
		"response structure validation", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/catalog", nil)
			rr := httptest.NewRecorder()

			handler.ListCatalog(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", rr.Code)
			}

			var response Response
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Validate response structure
			if response.Total != 1 {
				t.Errorf("Expected total 1, got %d", response.Total)
			}

			if len(response.Products) != 1 {
				t.Errorf("Expected 1 product, got %d", len(response.Products))
			}

			product := response.Products[0]
			if product.Code != "TEST001" {
				t.Errorf("Expected product code TEST001, got %s", product.Code)
			}

			if product.Price != 99.99 {
				t.Errorf("Expected product price 99.99, got %f", product.Price)
			}

			if product.Category.Code != "TEST" {
				t.Errorf("Expected category code TEST, got %s", product.Category.Code)
			}

			if product.Category.Name != "Test Category" {
				t.Errorf("Expected category name 'Test Category', got %s", product.Category.Name)
			}
		},
	)

}

func TestCatalogHandler_ListCatalog_PaginationLogic(t *testing.T) {
	// Create a larger dataset to test pagination properly
	var mockProducts []models.Product
	for i := 1; i <= 25; i++ {
		mockProducts = append(
			mockProducts, models.Product{
				Code:  fmt.Sprintf("PROD%03d", i),
				Price: decimal.NewFromFloat(float64(i * 10)),
				Category: models.Category{
					Code: "TEST",
					Name: "Test Category",
				},
			},
		)
	}

	mockRepo := &mockProductsRepository{
		products: mockProducts,
		count:    25,
	}

	handler := NewCatalogHandler(mockRepo)

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
			expectedCount: 10,
			expectedFirst: "PROD001",
			expectedLast:  "PROD010",
		},
		{
			name:          "second page default limit",
			queryParams:   "?offset=10",
			expectedCount: 10,
			expectedFirst: "PROD011",
			expectedLast:  "PROD020",
		},
		{
			name:          "custom limit",
			queryParams:   "?limit=5",
			expectedCount: 5,
			expectedFirst: "PROD001",
			expectedLast:  "PROD005",
		},
		{
			name:          "custom limit and offset",
			queryParams:   "?limit=3&offset=20",
			expectedCount: 3,
			expectedFirst: "PROD021",
			expectedLast:  "PROD023",
		},
		{
			name:          "offset beyond data",
			queryParams:   "?offset=30",
			expectedCount: 0,
			expectedFirst: "",
			expectedLast:  "",
		},
		{
			name:          "partial last page",
			queryParams:   "?limit=10&offset=20",
			expectedCount: 5,
			expectedFirst: "PROD021",
			expectedLast:  "PROD025",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				req := httptest.NewRequest("GET", "/catalog"+tt.queryParams, nil)
				rr := httptest.NewRecorder()

				handler.ListCatalog(rr, req)

				if rr.Code != http.StatusOK {
					t.Fatalf("Expected status 200, got %d", rr.Code)
				}

				var response Response
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if len(response.Products) != tt.expectedCount {
					t.Errorf("Expected %d products, got %d", tt.expectedCount, len(response.Products))
				}

				if tt.expectedCount > 0 {
					if response.Products[0].Code != tt.expectedFirst {
						t.Errorf("Expected first product %s, got %s", tt.expectedFirst, response.Products[0].Code)
					}

					if response.Products[len(response.Products)-1].Code != tt.expectedLast {
						t.Errorf(
							"Expected last product %s, got %s",
							tt.expectedLast,
							response.Products[len(response.Products)-1].Code,
						)
					}
				}

				if response.Total != 25 {
					t.Errorf("Expected total 25, got %d", response.Total)
				}
			},
		)
	}
}

func TestCatalogHandler_ListCatalog_CategoryFilter(t *testing.T) {
	// Test data with different categories
	mockProducts := []models.Product{
		{
			Code:  "CLOTHING001",
			Price: decimal.NewFromFloat(50.00),
			Category: models.Category{
				Code: "CLOTHING",
				Name: "Clothing",
			},
		},
		{
			Code:  "SHOES001",
			Price: decimal.NewFromFloat(100.00),
			Category: models.Category{
				Code: "SHOES",
				Name: "Shoes",
			},
		},
	}

	mockRepo := &mockProductsRepository{
		products: mockProducts,
		count:    2,
	}

	handler := NewCatalogHandler(mockRepo)

	t.Run(
		"filter by category", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/catalog?category=Clothing", nil)
			rr := httptest.NewRecorder()

			handler.ListCatalog(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", rr.Code)
			}

			var response Response
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Since mock doesn't filter, we just verify the request was processed correctly
			if len(response.Products) != 2 {
				t.Errorf("Expected 2 products, got %d", len(response.Products))
			}
		},
	)
}

func TestCatalogHandler_ListCatalog_PriceFilter(t *testing.T) {
	// Test data with different prices
	mockProducts := []models.Product{
		{
			Code:  "CHEAP001",
			Price: decimal.NewFromFloat(25.00),
			Category: models.Category{
				Code: "TEST",
				Name: "Test",
			},
		},
		{
			Code:  "EXPENSIVE001",
			Price: decimal.NewFromFloat(150.00),
			Category: models.Category{
				Code: "TEST",
				Name: "Test",
			},
		},
	}

	mockRepo := &mockProductsRepository{
		products: mockProducts,
		count:    2,
	}

	handler := NewCatalogHandler(mockRepo)

	t.Run(
		"filter by price less than", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/catalog?priceLessThen=100.00", nil)
			rr := httptest.NewRecorder()

			handler.ListCatalog(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", rr.Code)
			}

			var response Response
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Since mock doesn't filter, we just verify the request was processed correctly
			if len(response.Products) != 2 {
				t.Errorf("Expected 2 products, got %d", len(response.Products))
			}
		},
	)
}
