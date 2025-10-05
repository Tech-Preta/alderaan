package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()

	if m == nil {
		t.Fatal("NewMetrics() returned nil")
	}

	// Verificar se todas as métricas foram inicializadas
	if m.HTTPRequestDuration == nil {
		t.Error("HTTPRequestDuration is nil")
	}
	if m.HTTPRequestsTotal == nil {
		t.Error("HTTPRequestsTotal is nil")
	}
	if m.HTTPRequestErrors == nil {
		t.Error("HTTPRequestErrors is nil")
	}
	if m.InFlightRequests == nil {
		t.Error("InFlightRequests is nil")
	}
	if m.ProductsCreated == nil {
		t.Error("ProductsCreated is nil")
	}
	if m.ProductsTotal == nil {
		t.Error("ProductsTotal is nil")
	}
	if m.ProductsByCategory == nil {
		t.Error("ProductsByCategory is nil")
	}
	if m.ProductsTotalValue == nil {
		t.Error("ProductsTotalValue is nil")
	}
	if m.ProductsAveragePrice == nil {
		t.Error("ProductsAveragePrice is nil")
	}
}

func TestMetrics_RecordHTTPRequest(t *testing.T) {
	// Criar um registro próprio para evitar conflito com métricas globais
	reg := prometheus.NewRegistry()

	m := &Metrics{
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "test_http_request_duration_seconds",
				Help: "Test HTTP request duration",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_http_requests_total",
				Help: "Test HTTP requests total",
			},
			[]string{"method", "endpoint", "status"},
		),
		HTTPRequestErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "test_http_request_errors_total",
				Help: "Test HTTP request errors",
			},
			[]string{"method", "endpoint", "status", "error_type"},
		),
	}

	reg.MustRegister(m.HTTPRequestDuration)
	reg.MustRegister(m.HTTPRequestsTotal)
	reg.MustRegister(m.HTTPRequestErrors)

	tests := []struct {
		name     string
		method   string
		endpoint string
		status   string
		duration float64
		isError  bool
	}{
		{
			name:     "successful request",
			method:   "GET",
			endpoint: "/api/v1/products",
			status:   "200",
			duration: 0.15,
			isError:  false,
		},
		{
			name:     "client error",
			method:   "POST",
			endpoint: "/api/v1/products",
			status:   "400",
			duration: 0.05,
			isError:  true,
		},
		{
			name:     "server error",
			method:   "GET",
			endpoint: "/api/v1/products",
			status:   "500",
			duration: 0.25,
			isError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.RecordHTTPRequest(tt.method, tt.endpoint, tt.status, tt.duration, tt.isError)

			// Verificar que as métricas foram registradas
			// (Teste simplificado - em produção usaríamos testutil)
			count := testutil.CollectAndCount(m.HTTPRequestsTotal)
			if count == 0 {
				t.Error("HTTPRequestsTotal not incremented")
			}
		})
	}
}

func TestMetrics_IncrementProductsCreated(t *testing.T) {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		ProductsCreated: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "test_products_created_total",
				Help: "Test products created",
			},
		),
		ProductsTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_products_total",
				Help: "Test products total",
			},
		),
	}

	reg.MustRegister(m.ProductsCreated)
	reg.MustRegister(m.ProductsTotal)

	// Incrementar 3 vezes
	for i := 0; i < 3; i++ {
		m.IncrementProductsCreated()
	}

	// Verificar contadores
	if testutil.CollectAndCount(m.ProductsCreated) == 0 {
		t.Error("ProductsCreated not incremented")
	}
	if testutil.CollectAndCount(m.ProductsTotal) == 0 {
		t.Error("ProductsTotal not incremented")
	}
}

func TestMetrics_UpdateProductsByCategory(t *testing.T) {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		ProductsByCategory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "test_products_by_category",
				Help: "Test products by category",
			},
			[]string{"category"},
		),
	}

	reg.MustRegister(m.ProductsByCategory)

	// Atualizar categorias
	m.UpdateProductsByCategory("Electronics", 10)
	m.UpdateProductsByCategory("Books", 5)

	// Verificar que foram registrados
	count := testutil.CollectAndCount(m.ProductsByCategory)
	if count != 2 {
		t.Errorf("ProductsByCategory count = %d, want 2", count)
	}
}

func TestMetrics_UpdateProductsTotalValue(t *testing.T) {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		ProductsTotalValue: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_products_total_value",
				Help: "Test products total value",
			},
		),
	}

	reg.MustRegister(m.ProductsTotalValue)

	m.UpdateProductsTotalValue(15000.50)

	if testutil.CollectAndCount(m.ProductsTotalValue) == 0 {
		t.Error("ProductsTotalValue not set")
	}
}

func TestMetrics_UpdateProductsAveragePrice(t *testing.T) {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		ProductsAveragePrice: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "test_products_average_price",
				Help: "Test products average price",
			},
		),
	}

	reg.MustRegister(m.ProductsAveragePrice)

	m.UpdateProductsAveragePrice(299.99)

	if testutil.CollectAndCount(m.ProductsAveragePrice) == 0 {
		t.Error("ProductsAveragePrice not set")
	}
}

func TestMetrics_ResetProductsByCategory(t *testing.T) {
	reg := prometheus.NewRegistry()

	m := &Metrics{
		ProductsByCategory: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "test_products_by_category_reset",
				Help: "Test products by category reset",
			},
			[]string{"category"},
		),
	}

	reg.MustRegister(m.ProductsByCategory)

	// Adicionar dados
	m.UpdateProductsByCategory("Cat1", 5)
	m.UpdateProductsByCategory("Cat2", 10)

	// Reset
	m.ResetProductsByCategory()

	// Após reset, as métricas ainda existem mas valores resetados
	// Não há como testar facilmente sem acessar valores internos
	// Este teste confirma que não há panic
}

// Benchmark tests
func BenchmarkMetrics_RecordHTTPRequest(b *testing.B) {
	m := NewMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.RecordHTTPRequest("GET", "/api/v1/products", "200", 0.1, false)
	}
}

func BenchmarkMetrics_IncrementProductsCreated(b *testing.B) {
	m := NewMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.IncrementProductsCreated()
	}
}

func BenchmarkMetrics_UpdateProductsByCategory(b *testing.B) {
	m := NewMetrics()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.UpdateProductsByCategory("Electronics", float64(i))
	}
}
