package product_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
	product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
	"github.com/williamkoller/golang-domain-driven-design/internal/metrics"
	shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
)

// createTestMetrics cria métricas isoladas para testes
func createTestMetrics(testName string) *metrics.Metrics {
	return &metrics.Metrics{
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_" + testName + "_http_request_duration_seconds",
				Help: "Test HTTP request duration",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_" + testName + "_http_requests_total",
				Help: "Test HTTP requests total",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_" + testName + "_http_request_errors_total",
				Help: "Test HTTP request errors",
			},
			[]string{"method", "endpoint", "status", "error_type"},
		),
		InFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_" + testName + "_http_in_flight_requests",
				Help: "Test in-flight requests",
			},
		),
		ProductsCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "test_" + testName + "_products_created_total",
				Help: "Test products created",
			},
		),
		ProductsTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_" + testName + "_products_total",
				Help: "Test products total",
			},
		),
		ProductsByCategory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "test_" + testName + "_products_by_category",
				Help: "Test products by category",
			},
			[]string{"category"},
		),
		ProductsTotalValue: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_" + testName + "_products_total_value",
				Help: "Test products total value",
			},
		),
		ProductsAveragePrice: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_" + testName + "_products_average_price",
				Help: "Test products average price",
			},
		),
	}
}

// MockProductRepository implementa IProductRepository para testes
type MockProductRepository struct {
	products        map[string]product_entity.Product
	addError        error
	findError       error
	findOneError    error
	metricsToReturn product_repository.RepositoryMetrics
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[string]product_entity.Product),
		metricsToReturn: product_repository.RepositoryMetrics{
			TotalProducts:      0,
			TotalValue:         0,
			AveragePrice:       0,
			ProductsByCategory: make(map[string]int),
		},
	}
}

func (m *MockProductRepository) Add(product product_entity.Product) error {
	if m.addError != nil {
		return m.addError
	}
	m.products[product.Name] = product
	return nil
}

func (m *MockProductRepository) Find() ([]product_entity.Product, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	products := make([]product_entity.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MockProductRepository) FindOne(name string) (product_entity.Product, error) {
	if m.findOneError != nil {
		return product_entity.Product{}, m.findOneError
	}
	product, exists := m.products[name]
	if !exists {
		return product_entity.Product{}, errors.New("product not found")
	}
	return product, nil
}

func (m *MockProductRepository) GetMetrics() product_repository.RepositoryMetrics {
	// Calcular métricas reais baseadas nos produtos mock
	totalValue := 0
	productsByCategory := make(map[string]int)

	for _, product := range m.products {
		totalValue += product.Price
		for _, category := range product.Categories {
			productsByCategory[category]++
		}
	}

	avgPrice := 0.0
	if len(m.products) > 0 {
		avgPrice = float64(totalValue) / float64(len(m.products))
	}

	return product_repository.RepositoryMetrics{
		TotalProducts:      len(m.products),
		TotalValue:         float64(totalValue),
		AveragePrice:       avgPrice,
		ProductsByCategory: productsByCategory,
	}
}

func setupTestRouter(handler *ProductHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/products", handler.Create)
		v1.GET("/products", handler.FindAll)
		v1.GET("/products/:name", handler.FindOne)
	}

	return router
}

func TestNewProductHandler(t *testing.T) {
	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("new_handler")

	handler := NewProductHandler(repo, dispatcher, m)

	if handler == nil {
		t.Fatal("NewProductHandler() returned nil")
	}

	if handler.repo == nil {
		t.Error("handler.repo is nil")
	}

	if handler.dispatcher == nil {
		t.Error("handler.dispatcher is nil")
	}

	if handler.metrics == nil {
		t.Error("handler.metrics is nil")
	}
}

func TestProductHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		setupMock      func(*MockProductRepository)
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "create product successfully",
			requestBody: CreateProductInput{
				Name:       "Notebook",
				Sku:        12345,
				Categories: []string{"Electronics", "Computers"},
				Price:      3500,
			},
			expectedStatus: http.StatusCreated,
			setupMock:      func(m *MockProductRepository) {},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response product_entity.Product
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.Name != "Notebook" {
					t.Errorf("Expected name Notebook, got %s", response.Name)
				}
			},
		},
		{
			name:           "invalid JSON body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockProductRepository) {},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
			},
		},
		{
			name: "missing required field - name",
			requestBody: CreateProductInput{
				Name:       "",
				Sku:        12345,
				Categories: []string{"Electronics"},
				Price:      3500,
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockProductRepository) {},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
			},
		},
		{
			name: "invalid sku (zero)",
			requestBody: CreateProductInput{
				Name:       "Test Product",
				Sku:        0,
				Categories: []string{"Test"},
				Price:      100,
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockProductRepository) {},
			checkResponse:  func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name: "invalid price (zero)",
			requestBody: CreateProductInput{
				Name:       "Test Product",
				Sku:        123,
				Categories: []string{"Test"},
				Price:      0,
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockProductRepository) {},
			checkResponse:  func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name: "empty categories",
			requestBody: CreateProductInput{
				Name:       "Test Product",
				Sku:        123,
				Categories: []string{},
				Price:      100,
			},
			expectedStatus: http.StatusBadRequest,
			setupMock:      func(m *MockProductRepository) {},
			checkResponse:  func(t *testing.T, w *httptest.ResponseRecorder) {},
		},
		{
			name: "duplicate product",
			requestBody: CreateProductInput{
				Name:       "Existing Product",
				Sku:        999,
				Categories: []string{"Test"},
				Price:      500,
			},
			expectedStatus: http.StatusConflict,
			setupMock: func(m *MockProductRepository) {
				m.addError = errors.New("product already exists")
			},
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
			},
		},
	}

	testID := 0
	for _, tt := range tests {
		testID++
		currentTestID := testID // Capturar para closure
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockProductRepository()
			dispatcher := shared_events.NewEventDispatcher()
			m := createTestMetrics(fmt.Sprintf("create_%d", currentTestID))

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			handler := NewProductHandler(mockRepo, dispatcher, m)
			router := setupTestRouter(handler)

			// Criar request
			var body []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestProductHandler_FindAll(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockProductRepository)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "find all products - empty repository",
			setupMock: func(m *MockProductRepository) {
				// Repository vazio
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "find all products - with products",
			setupMock: func(m *MockProductRepository) {
				m.products["Product1"] = product_entity.Product{
					Name:       "Product1",
					Sku:        1,
					Categories: []string{"Cat1"},
					Price:      100,
				}
				m.products["Product2"] = product_entity.Product{
					Name:       "Product2",
					Sku:        2,
					Categories: []string{"Cat2"},
					Price:      200,
				}
				m.products["Product3"] = product_entity.Product{
					Name:       "Product3",
					Sku:        3,
					Categories: []string{"Cat3"},
					Price:      300,
				}
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name: "find all products - repository error",
			setupMock: func(m *MockProductRepository) {
				m.findError = errors.New("database error")
			},
			expectedStatus: http.StatusOK, // Handler não retorna erro, retorna array vazio
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockProductRepository()
			dispatcher := shared_events.NewEventDispatcher()
			m := createTestMetrics("findall_" + tt.name)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			handler := NewProductHandler(mockRepo, dispatcher, m)
			router := setupTestRouter(handler)

			// Criar request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var products []product_entity.Product
			if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil {
				t.Errorf("Failed to unmarshal response: %v", err)
			}

			if len(products) != tt.expectedCount {
				t.Errorf("Expected %d products, got %d", tt.expectedCount, len(products))
			}
		})
	}
}

func TestProductHandler_FindOne(t *testing.T) {
	tests := []struct {
		name           string
		productName    string
		setupMock      func(*MockProductRepository)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "find existing product",
			productName: "Notebook",
			setupMock: func(m *MockProductRepository) {
				m.products["Notebook"] = product_entity.Product{
					Name:       "Notebook",
					Sku:        12345,
					Categories: []string{"Electronics"},
					Price:      3500,
				}
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var product product_entity.Product
				if err := json.Unmarshal(w.Body.Bytes(), &product); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if product.Name != "Notebook" {
					t.Errorf("Expected name Notebook, got %s", product.Name)
				}
				if product.Sku != 12345 {
					t.Errorf("Expected SKU 12345, got %d", product.Sku)
				}
			},
		},
		{
			name:        "product not found",
			productName: "NonExistent",
			setupMock: func(m *MockProductRepository) {
				// Repository vazio
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if response.Error != "product not found" {
					t.Errorf("Expected error 'product not found', got %s", response.Error)
				}
			},
		},
		{
			name:        "repository error",
			productName: "Test",
			setupMock: func(m *MockProductRepository) {
				m.findOneError = errors.New("database error")
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ErrorResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockProductRepository()
			dispatcher := shared_events.NewEventDispatcher()
			m := createTestMetrics("findone_" + tt.name)

			if tt.setupMock != nil {
				tt.setupMock(mockRepo)
			}

			handler := NewProductHandler(mockRepo, dispatcher, m)
			router := setupTestRouter(handler)

			// Criar request
			req := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+tt.productName, nil)
			w := httptest.NewRecorder()

			// Execute
			router.ServeHTTP(w, req)

			// Assert
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestProductHandler_Integration(t *testing.T) {
	t.Run("create and retrieve product", func(t *testing.T) {
		// Setup
		mockRepo := NewMockProductRepository()
		dispatcher := shared_events.NewEventDispatcher()
		m := createTestMetrics("integration")
		handler := NewProductHandler(mockRepo, dispatcher, m)
		router := setupTestRouter(handler)

		// 1. Criar produto
		createInput := CreateProductInput{
			Name:       "IntegrationTestProduct",
			Sku:        99999,
			Categories: []string{"TestCategory"},
			Price:      1000,
		}
		body, _ := json.Marshal(createInput)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Failed to create product: status %d", w.Code)
		}

		// 2. Buscar produto criado
		req = httptest.NewRequest(http.MethodGet, "/api/v1/products/IntegrationTestProduct", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Failed to find product: status %d", w.Code)
		}

		var product product_entity.Product
		if err := json.Unmarshal(w.Body.Bytes(), &product); err != nil {
			t.Fatalf("Failed to unmarshal product: %v", err)
		}

		if product.Name != "IntegrationTestProduct" {
			t.Errorf("Expected name 'IntegrationTestProduct', got %s", product.Name)
		}

		// 3. Listar todos os produtos
		req = httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Failed to list products: status %d", w.Code)
		}

		var products []product_entity.Product
		if err := json.Unmarshal(w.Body.Bytes(), &products); err != nil {
			t.Fatalf("Failed to unmarshal products: %v", err)
		}

		if len(products) != 1 {
			t.Errorf("Expected 1 product, got %d", len(products))
		}
	})
}

// Benchmark
func BenchmarkProductHandler_Create(b *testing.B) {
	mockRepo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("bench_create")
	handler := NewProductHandler(mockRepo, dispatcher, m)
	router := setupTestRouter(handler)

	input := CreateProductInput{
		Name:       "Benchmark Product",
		Sku:        12345,
		Categories: []string{"Benchmark"},
		Price:      999,
	}
	body, _ := json.Marshal(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkProductHandler_FindAll(b *testing.B) {
	mockRepo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("bench_findall")

	// Adicionar alguns produtos
	for i := 0; i < 100; i++ {
		mockRepo.products["Product"+string(rune(i))] = product_entity.Product{
			Name:       "Product" + string(rune(i)),
			Sku:        i,
			Categories: []string{"Test"},
			Price:      100,
		}
	}

	handler := NewProductHandler(mockRepo, dispatcher, m)
	router := setupTestRouter(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkProductHandler_FindOne(b *testing.B) {
	mockRepo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("bench_findone")

	mockRepo.products["TestProduct"] = product_entity.Product{
		Name:       "TestProduct",
		Sku:        123,
		Categories: []string{"Test"},
		Price:      100,
	}

	handler := NewProductHandler(mockRepo, dispatcher, m)
	router := setupTestRouter(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/products/TestProduct", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
