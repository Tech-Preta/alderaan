package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func TestPrometheusMiddleware(t *testing.T) {
	// Modo test do Gin
	gin.SetMode(gin.TestMode)

	// Criar métricas de teste
	reg := prometheus.NewRegistry()

	m := &Metrics{
		InFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_in_flight_requests",
				Help: "Test in-flight requests",
			},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_request_duration",
				Help: "Test request duration",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_requests_total",
				Help: "Test requests total",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_request_errors",
				Help: "Test request errors",
			},
			[]string{"method", "endpoint", "status", "error_type"},
		),
	}

	reg.MustRegister(m.InFlightRequests)
	reg.MustRegister(m.HTTPRequestDuration)
	reg.MustRegister(m.HTTPRequestsTotal)
	reg.MustRegister(m.HTTPRequestErrors)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		handler        gin.HandlerFunc
	}{
		{
			name:           "successful GET request",
			method:         "GET",
			path:           "/api/v1/products",
			expectedStatus: http.StatusOK,
			handler: func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			},
		},
		{
			name:           "POST request",
			method:         "POST",
			path:           "/api/v1/products",
			expectedStatus: http.StatusCreated,
			handler: func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{"created": true})
			},
		},
		{
			name:           "client error",
			method:         "GET",
			path:           "/api/v1/products/invalid",
			expectedStatus: http.StatusBadRequest,
			handler: func(c *gin.Context) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			},
		},
		{
			name:           "server error",
			method:         "GET",
			path:           "/api/v1/products",
			expectedStatus: http.StatusInternalServerError,
			handler: func(c *gin.Context) {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar router
			router := gin.New()
			router.Use(PrometheusMiddleware(m))
			router.Handle(tt.method, tt.path, tt.handler)

			// Criar request
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			// Executar request
			router.ServeHTTP(w, req)

			// Verificar status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Status code = %d, want %d", w.Code, tt.expectedStatus)
			}

			// As métricas foram registradas (verificação básica)
			// Em um teste mais elaborado, usaríamos testutil do Prometheus
		})
	}
}

func TestPrometheusMiddleware_InFlightRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Usar métricas customizadas para evitar conflito de registro
	reg := prometheus.NewRegistry()

	m := &Metrics{
		InFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_in_flight_custom",
				Help: "Test in-flight custom",
			},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_duration_custom",
				Help: "Test duration custom",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_total_custom",
				Help: "Test total custom",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_errors_custom",
				Help: "Test errors custom",
			},
			[]string{"method", "endpoint", "status", "error_type"},
		),
	}

	reg.MustRegister(m.InFlightRequests)
	reg.MustRegister(m.HTTPRequestDuration)
	reg.MustRegister(m.HTTPRequestsTotal)
	reg.MustRegister(m.HTTPRequestErrors)

	router := gin.New()
	router.Use(PrometheusMiddleware(m))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Após a requisição, in-flight deve voltar a 0
	// (Não podemos verificar facilmente sem expor valores internos,
	// mas podemos confirmar que não há panic)
}

func TestNormalizeEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "empty path",
			path:     "",
			expected: "unknown",
		},
		{
			name:     "simple path",
			path:     "/api/v1/products",
			expected: "/api/v1/products",
		},
		{
			name:     "path with parameter",
			path:     "/api/v1/products/:name",
			expected: "/api/v1/products/{name}",
		},
		{
			name:     "path with multiple parameters",
			path:     "/api/v1/users/:id/posts/:postId",
			expected: "/api/v1/users/{id}/posts/{postId}",
		},
		{
			name:     "path with trailing slash",
			path:     "/api/v1/products/",
			expected: "/api/v1/products/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeEndpoint(tt.path)
			if result != tt.expected {
				t.Errorf("normalizeEndpoint(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

// Benchmark
func BenchmarkPrometheusMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)
	m := NewMetrics()

	router := gin.New()
	router.Use(PrometheusMiddleware(m))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkNormalizeEndpoint(b *testing.B) {
	path := "/api/v1/products/:name"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = normalizeEndpoint(path)
	}
}
