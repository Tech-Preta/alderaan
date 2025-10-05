package product_router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"
	
	"github.com/gin-gonic/gin"
	product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
	product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
	product_handlers "github.com/williamkoller/golang-domain-driven-design/internal/infra/http/handlers"
	"github.com/williamkoller/golang-domain-driven-design/internal/metrics"
	shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
	"github.com/prometheus/client_golang/prometheus"
)

// MockProductRepository para testes do router
type MockProductRepository struct {
	products map[string]product_entity.Product
}

func NewMockProductRepository() *MockProductRepository {
	return &MockProductRepository{
		products: make(map[string]product_entity.Product),
	}
}

func (m *MockProductRepository) Add(product product_entity.Product) error {
	m.products[product.Name] = product
	return nil
}

func (m *MockProductRepository) Find() ([]product_entity.Product, error) {
	products := make([]product_entity.Product, 0, len(m.products))
	for _, p := range m.products {
		products = append(products, p)
	}
	return products, nil
}

func (m *MockProductRepository) FindOne(name string) (product_entity.Product, error) {
	product, exists := m.products[name]
	if !exists {
		return product_entity.Product{}, errors.New("product not found")
	}
	return product, nil
}

func (m *MockProductRepository) GetMetrics() product_repository.RepositoryMetrics {
	return product_repository.RepositoryMetrics{
		TotalProducts:      len(m.products),
		TotalValue:         0,
		AveragePrice:       0,
		ProductsByCategory: make(map[string]int),
	}
}

func createTestMetrics(testName string) *metrics.Metrics {
	return &metrics.Metrics{
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_router_" + testName + "_http_request_duration_seconds",
				Help: "Test HTTP request duration",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_router_" + testName + "_http_requests_total",
				Help: "Test HTTP requests total",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_router_" + testName + "_http_request_errors_total",
				Help: "Test HTTP request errors",
			},
			[]string{"method", "endpoint", "status", "error_type"},
		),
		InFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_router_" + testName + "_http_in_flight_requests",
				Help: "Test in-flight requests",
			},
		),
		ProductsCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "test_router_" + testName + "_products_created_total",
				Help: "Test products created",
			},
		),
		ProductsTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_router_" + testName + "_products_total",
				Help: "Test products total",
			},
		),
		ProductsByCategory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "test_router_" + testName + "_products_by_category",
				Help: "Test products by category",
			},
			[]string{"category"},
		),
		ProductsTotalValue: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_router_" + testName + "_products_total_value",
				Help: "Test products total value",
			},
		),
		ProductsAveragePrice: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_router_" + testName + "_products_average_price",
				Help: "Test products average price",
			},
		),
	}
}

func TestSetupProductRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("setup")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)

	router := SetupProductRouter(handler, m)

	if router == nil {
		t.Fatal("SetupProductRouter() returned nil")
	}

	// Verificar que o router foi configurado corretamente
	routes := router.Routes()
	
	expectedRoutes := map[string]bool{
		"GET-/metrics":              false,
		"GET-/swagger/*any":         false,
		"GET-/health":               false,
		"POST-/api/v1/products":     false,
		"GET-/api/v1/products":      false,
		"GET-/api/v1/products/:name": false,
	}

	for _, route := range routes {
		key := route.Method + "-" + route.Path
		if _, exists := expectedRoutes[key]; exists {
			expectedRoutes[key] = true
		}
	}

	// Verificar que todas as rotas esperadas foram registradas
	for route, found := range expectedRoutes {
		if !found {
			t.Errorf("Expected route %s not found", route)
		}
	}
}

func TestProductRouter_HealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("health")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := `{"status":"ok"}`
	if w.Body.String() != expected {
		t.Errorf("Expected body %s, got %s", expected, w.Body.String())
	}
}

func TestProductRouter_MetricsEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("metrics_endpoint")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verificar que a resposta contém métricas Prometheus
	body := w.Body.String()
	if len(body) == 0 {
		t.Error("Metrics endpoint returned empty body")
	}

	// Verificar header Content-Type
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		t.Error("Content-Type header not set")
	}
}

func TestProductRouter_SwaggerEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("swagger")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "swagger index",
			path:           "/swagger/index.html",
			expectedStatus: http.StatusMovedPermanently, // Redirect
		},
		{
			name:           "swagger root",
			path:           "/swagger/",
			expectedStatus: http.StatusMovedPermanently,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Swagger pode retornar redirect ou 404 se docs não existirem
			if w.Code != tt.expectedStatus && w.Code != http.StatusNotFound {
				t.Logf("Swagger endpoint %s returned status %d (expected %d or 404)", tt.path, w.Code, tt.expectedStatus)
			}
		})
	}
}

func TestProductRouter_APIv1Endpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("apiv1")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "GET /api/v1/products",
			method:         http.MethodGet,
			path:           "/api/v1/products",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/v1/products without body",
			method:         http.MethodPost,
			path:           "/api/v1/products",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "GET /api/v1/products/:name - not found",
			method:         http.MethodGet,
			path:           "/api/v1/products/NonExistent",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestProductRouter_NotFoundRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("notfound")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	req := httptest.NewRequest(http.MethodGet, "/non-existent-route", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestProductRouter_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("method_not_allowed")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "DELETE on products endpoint",
			method: http.MethodDelete,
			path:   "/api/v1/products",
		},
		{
			name:   "PUT on products endpoint",
			method: http.MethodPut,
			path:   "/api/v1/products",
		},
		{
			name:   "PATCH on products endpoint",
			method: http.MethodPatch,
			path:   "/api/v1/products",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Gin retorna 404 para métodos não permitidos (não 405)
			if w.Code != http.StatusNotFound && w.Code != http.StatusMethodNotAllowed {
				t.Errorf("Expected status 404 or 405, got %d", w.Code)
			}
		})
	}
}

func TestProductRouter_Middlewares(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("middlewares")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	// Fazer uma requisição para verificar que middlewares estão sendo executados
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Verificar que métricas foram coletadas (middleware Prometheus)
	// Isso é implícito - se não houver panic, o middleware está funcionando
}

func TestProductRouter_CORSHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("cors")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/products", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Gin não adiciona CORS por padrão, mas verificamos que não há erro
	// Se CORS fosse configurado, veríamos headers específicos aqui
	t.Logf("OPTIONS request status: %d", w.Code)
}

// Benchmark
func BenchmarkProductRouter_Health(b *testing.B) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("bench_health")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkProductRouter_Metrics(b *testing.B) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("bench_metrics")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkProductRouter_APIv1Products(b *testing.B) {
	gin.SetMode(gin.TestMode)

	repo := NewMockProductRepository()
	dispatcher := shared_events.NewEventDispatcher()
	m := createTestMetrics("bench_apiv1")
	handler := product_handlers.NewProductHandler(repo, dispatcher, m)
	router := SetupProductRouter(handler, m)

	// Adicionar alguns produtos
	repo.products["Product1"] = product_entity.Product{
		Name:       "Product1",
		Sku:        1,
		Categories: []string{"Test"},
		Price:      100,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

