package catalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/app/dto"
	"github.com/mytheresa/go-hiring-challenge/app/interfaces"
	"github.com/mytheresa/go-hiring-challenge/app/models"
	"github.com/shopspring/decimal"
)

// mockProductsRepository is a mock implementation of interfaces.ProductsRepository
type mockProductsRepository struct {
	products       []models.Product
	count          int64
	getError       error
	countError     error
	getByCodeError error
}

func (m *mockProductsRepository) GetProductByCode(code string) (*models.Product, error) {
	if m.getByCodeError != nil {
		return nil, m.getByCodeError
	}

	for _, product := range m.products {
		if product.Code == code {
			return &product, nil
		}
	}

	return nil, nil
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
			expectedBody: ListCatalogResponse{
				Total: 2,
				Products: []dto.Product{
					{
						Code:  "PROD001",
						Price: 100.50,
						Category: dto.Category{
							Code: "CLOTHING",
							Name: "Clothing",
						},
					},
					{
						Code:  "PROD002",
						Price: 200.75,
						Category: dto.Category{
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
			expectedBody: ListCatalogResponse{
				Total: 2,
				Products: []dto.Product{
					{
						Code:  "PROD002",
						Price: 200.75,
						Category: dto.Category{
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
			expectedBody: ListCatalogResponse{
				Total:    0,
				Products: []dto.Product{},
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
				handler := NewCatalogHandler(tt.repository)
				req := httptest.NewRequest("GET", "/catalog"+tt.queryParams, nil)
				rr := httptest.NewRecorder()
				handler.ListCatalog(rr, req)

				if rr.Code != tt.expectedStatus {
					t.Errorf("ListCatalog() status = %v, want %v", rr.Code, tt.expectedStatus)
				}

				if tt.checkBody && tt.expectedStatus == http.StatusOK {
					var response ListCatalogResponse
					err := json.Unmarshal(rr.Body.Bytes(), &response)
					if err != nil {
						t.Fatalf("Failed to unmarshal response: %v", err)
					}

					expectedResponse := tt.expectedBody.(ListCatalogResponse)
					if !reflect.DeepEqual(response, expectedResponse) {
						t.Errorf("ListCatalog() response = %v, want %v", response, expectedResponse)
					}
				}

				if tt.expectedStatus != http.StatusOK && tt.checkBody == false {
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

			var response ListCatalogResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

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

				var response ListCatalogResponse
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

			var response ListCatalogResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if len(response.Products) != 2 {
				t.Errorf("Expected 2 products, got %d", len(response.Products))
			}
		},
	)
}

func TestCatalogHandler_ListCatalog_PriceFilter(t *testing.T) {
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

			var response ListCatalogResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if len(response.Products) != 2 {
				t.Errorf("Expected 2 products, got %d", len(response.Products))
			}
		},
	)
}

func TestCatalogHandler_GetProductDetails(t *testing.T) {
	// Mock data setup
	mockProductWithVariants := models.Product{
		ID:    1,
		Code:  "PROD001",
		Price: decimal.NewFromFloat(99.99),
		Category: models.Category{
			Code: "CLOTHING",
			Name: "Clothing",
		},
		Variants: []models.Variant{
			{
				ID:    1,
				SKU:   "SKU001A",
				Name:  "Variant A",
				Price: decimal.NewFromFloat(89.99),
			},
			{
				ID:    2,
				SKU:   "SKU001B",
				Name:  "Variant B",
				Price: decimal.Zero, // Should inherit base price
			},
		},
	}

	mockProductNoVariants := models.Product{
		ID:    2,
		Code:  "PROD002",
		Price: decimal.NewFromFloat(150.00),
		Category: models.Category{
			Code: "SHOES",
			Name: "Shoes",
		},
		Variants: []models.Variant{},
	}

	tests := []struct {
		name           string
		productCode    string
		repository     *mockProductsRepository
		expectedStatus int
		expectedBody   interface{}
		checkBody      bool
	}{
		{
			name:        "successful get product with variants",
			productCode: "PROD001",
			repository: &mockProductsRepository{
				products: []models.Product{mockProductWithVariants},
			},
			expectedStatus: http.StatusOK,
			expectedBody: GetProductDetailsResponse{
				dto.ProductWithDetails{
					ID: 1,
					Product: dto.Product{
						Code:  "PROD001",
						Price: 99.99,
						Category: dto.Category{
							Code: "CLOTHING",
							Name: "Clothing",
						},
					},
					ProductVariant: []dto.ProductVariant{
						{
							ID:    1,
							SKU:   "SKU001A",
							Name:  "Variant A",
							Price: 89.99,
						},
						{
							ID:    2,
							SKU:   "SKU001B",
							Name:  "Variant B",
							Price: 99.99, // Inherited from base product
						},
					},
				},
			},
			checkBody: true,
		},
		{
			name:        "successful get product without variants",
			productCode: "PROD002",
			repository: &mockProductsRepository{
				products: []models.Product{mockProductNoVariants},
			},
			expectedStatus: http.StatusOK,
			expectedBody: GetProductDetailsResponse{
				dto.ProductWithDetails{
					ID: 2,
					Product: dto.Product{
						Code:  "PROD002",
						Price: 150.00,
						Category: dto.Category{
							Code: "SHOES",
							Name: "Shoes",
						},
					},
					ProductVariant: []dto.ProductVariant{},
				},
			},
			checkBody: true,
		},
		{
			name:        "product not found",
			productCode: "NONEXISTENT",
			repository: &mockProductsRepository{
				products: []models.Product{mockProductWithVariants},
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: GetProductDetailsResponse{
				dto.ProductWithDetails{
					ID: 0,
					Product: dto.Product{
						Code:  "",
						Price: 0,
						Category: dto.Category{
							Code: "",
							Name: "",
						},
					},
					ProductVariant: []dto.ProductVariant{},
				},
			},
			checkBody: true,
		},
		{
			name:        "repository error",
			productCode: "PROD001",
			repository: &mockProductsRepository{
				getByCodeError: errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			checkBody:      false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				handler := NewCatalogHandler(tt.repository)
				req := httptest.NewRequest("GET", "/catalog/"+tt.productCode, nil)
				req.SetPathValue("code", tt.productCode)
				rr := httptest.NewRecorder()

				handler.GetProductDetails(rr, req)

				if rr.Code != tt.expectedStatus {
					t.Errorf("GetProductDetails() status = %v, want %v", rr.Code, tt.expectedStatus)
				}

				if tt.checkBody && tt.expectedStatus == http.StatusOK {
					var response GetProductDetailsResponse
					err := json.Unmarshal(rr.Body.Bytes(), &response)
					if err != nil {
						t.Fatalf("Failed to unmarshal response: %v", err)
					}

					expectedResponse := tt.expectedBody.(GetProductDetailsResponse)
					if !reflect.DeepEqual(response, expectedResponse) {
						t.Errorf("GetProductDetails() response = %v, want %v", response, expectedResponse)
					}
				}

				if tt.expectedStatus != http.StatusOK && tt.checkBody == false {
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

func TestCatalogHandler_GetProductDetails_Integration(t *testing.T) {
	mockProduct := models.Product{
		ID:    1,
		Code:  "INTEGRATION001",
		Price: decimal.NewFromFloat(199.99),
		Category: models.Category{
			Code: "INTEGRATION",
			Name: "Integration Test Category",
		},
		Variants: []models.Variant{
			{
				ID:    10,
				SKU:   "INT001A",
				Name:  "Integration Variant A",
				Price: decimal.NewFromFloat(179.99),
			},
			{
				ID:    11,
				SKU:   "INT001B",
				Name:  "Integration Variant B",
				Price: decimal.Zero, // Should inherit base price
			},
		},
	}

	mockRepo := &mockProductsRepository{
		products: []models.Product{mockProduct},
	}

	handler := NewCatalogHandler(mockRepo)

	t.Run(
		"correct content type", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/catalog/INTEGRATION001", nil)
			req.SetPathValue("code", "INTEGRATION001")
			rr := httptest.NewRecorder()

			handler.GetProductDetails(rr, req)

			contentType := rr.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}
		},
	)

	t.Run(
		"response structure validation", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/catalog/INTEGRATION001", nil)
			req.SetPathValue("code", "INTEGRATION001")
			rr := httptest.NewRecorder()

			handler.GetProductDetails(rr, req)

			if rr.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", rr.Code)
			}

			var response GetProductDetailsResponse
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Validate product details
			if response.ID != 1 {
				t.Errorf("Expected product ID 1, got %d", response.ID)
			}

			if response.Code != "INTEGRATION001" {
				t.Errorf("Expected product code INTEGRATION001, got %s", response.Code)
			}

			if response.Price != 199.99 {
				t.Errorf("Expected product price 199.99, got %f", response.Price)
			}

			if response.Category.Code != "INTEGRATION" {
				t.Errorf("Expected category code INTEGRATION, got %s", response.Category.Code)
			}

			if response.Category.Name != "Integration Test Category" {
				t.Errorf("Expected category name 'Integration Test Category', got %s", response.Category.Name)
			}

			// Validate variants
			if len(response.ProductVariant) != 2 {
				t.Errorf("Expected 2 variants, got %d", len(response.ProductVariant))
			}

			// First variant should have its own price
			variant1 := response.ProductVariant[0]
			if variant1.ID != 10 {
				t.Errorf("Expected variant ID 10, got %d", variant1.ID)
			}
			if variant1.SKU != "INT001A" {
				t.Errorf("Expected variant SKU INT001A, got %s", variant1.SKU)
			}
			if variant1.Name != "Integration Variant A" {
				t.Errorf("Expected variant name 'Integration Variant A', got %s", variant1.Name)
			}
			if variant1.Price != 179.99 {
				t.Errorf("Expected variant price 179.99, got %f", variant1.Price)
			}

			// Second variant should inherit base price
			variant2 := response.ProductVariant[1]
			if variant2.ID != 11 {
				t.Errorf("Expected variant ID 11, got %d", variant2.ID)
			}
			if variant2.SKU != "INT001B" {
				t.Errorf("Expected variant SKU INT001B, got %s", variant2.SKU)
			}
			if variant2.Name != "Integration Variant B" {
				t.Errorf("Expected variant name 'Integration Variant B', got %s", variant2.Name)
			}
			if variant2.Price != 199.99 {
				t.Errorf("Expected variant price 199.99 (inherited from base), got %f", variant2.Price)
			}
		},
	)
}

func TestCatalogHandler_GetProductDetails_PriceInheritance(t *testing.T) {
	tests := []struct {
		name          string
		basePrice     float64
		variantPrice  float64
		expectedPrice float64
		description   string
	}{
		{
			name:          "variant with zero price inherits base price",
			basePrice:     100.00,
			variantPrice:  0.00,
			expectedPrice: 100.00,
			description:   "When variant price is zero, it should inherit the base product price",
		},
		{
			name:          "variant with specific price keeps its own price",
			basePrice:     100.00,
			variantPrice:  75.50,
			expectedPrice: 75.50,
			description:   "When variant has its own price, it should keep that price",
		},
		{
			name:          "variant price can be higher than base price",
			basePrice:     50.00,
			variantPrice:  150.00,
			expectedPrice: 150.00,
			description:   "Variant price can be higher than base price",
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				mockProduct := models.Product{
					ID:    1,
					Code:  "PRICE_TEST",
					Price: decimal.NewFromFloat(tt.basePrice),
					Category: models.Category{
						Code: "TEST",
						Name: "Test Category",
					},
					Variants: []models.Variant{
						{
							ID:    1,
							SKU:   "TEST_SKU",
							Name:  "Test Variant",
							Price: decimal.NewFromFloat(tt.variantPrice),
						},
					},
				}

				mockRepo := &mockProductsRepository{
					products: []models.Product{mockProduct},
				}

				handler := NewCatalogHandler(mockRepo)
				req := httptest.NewRequest("GET", "/catalog/PRICE_TEST", nil)
				req.SetPathValue("code", "PRICE_TEST")
				rr := httptest.NewRecorder()

				handler.GetProductDetails(rr, req)

				if rr.Code != http.StatusOK {
					t.Fatalf("Expected status 200, got %d", rr.Code)
				}

				var response GetProductDetailsResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if len(response.ProductVariant) != 1 {
					t.Fatalf("Expected 1 variant, got %d", len(response.ProductVariant))
				}

				actualPrice := response.ProductVariant[0].Price
				if actualPrice != tt.expectedPrice {
					t.Errorf("%s: Expected variant price %f, got %f", tt.description, tt.expectedPrice, actualPrice)
				}
			},
		)
	}
}
