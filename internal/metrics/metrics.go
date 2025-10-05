package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics contém todas as métricas da aplicação
type Metrics struct {
	// Golden Signals - Latency
	HTTPRequestDuration *prometheus.HistogramVec

	// Golden Signals - Traffic
	HTTPRequestsTotal *prometheus.CounterVec

	// Golden Signals - Errors
	HTTPRequestErrors *prometheus.CounterVec

	// Golden Signals - Saturation
	InFlightRequests prometheus.Gauge

	// Métricas de Negócio
	ProductsCreated      prometheus.Counter
	ProductsTotal        prometheus.Gauge
	ProductsByCategory   *prometheus.GaugeVec
	ProductsTotalValue   prometheus.Gauge
	ProductsAveragePrice prometheus.Gauge
}

// NewMetrics cria e registra todas as métricas
func NewMetrics() *Metrics {
	return &Metrics{
		// Golden Signals - Latency
		// Mede o tempo de resposta das requisições HTTP
		HTTPRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "Duração das requisições HTTP em segundos (Latency)",
				Buckets: []float64{
					0.001, // 1ms
					0.005, // 5ms
					0.01,  // 10ms
					0.025, // 25ms
					0.05,  // 50ms
					0.1,   // 100ms
					0.25,  // 250ms
					0.5,   // 500ms
					1.0,   // 1s
					2.5,   // 2.5s
					5.0,   // 5s
					10.0,  // 10s
				},
			},
			[]string{"method", "endpoint", "status"},
		),

		// Golden Signals - Traffic
		// Conta o número total de requisições
		HTTPRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total de requisições HTTP recebidas (Traffic)",
			},
			[]string{"method", "endpoint", "status"},
		),

		// Golden Signals - Errors
		// Conta erros HTTP (status >= 400)
		HTTPRequestErrors: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_request_errors_total",
				Help: "Total de erros em requisições HTTP (Errors)",
			},
			[]string{"method", "endpoint", "status", "error_type"},
		),

		// Golden Signals - Saturation
		// Mede requisições sendo processadas simultaneamente
		InFlightRequests: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_in_flight_requests",
				Help: "Número de requisições HTTP sendo processadas no momento (Saturation)",
			},
		),

		// Métricas de Negócio - Produtos Criados
		ProductsCreated: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "products_created_total",
				Help: "Total de produtos criados desde o início da aplicação",
			},
		),

		// Métricas de Negócio - Total de Produtos
		ProductsTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "products_total",
				Help: "Número atual de produtos cadastrados",
			},
		),

		// Métricas de Negócio - Produtos por Categoria
		ProductsByCategory: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "products_by_category",
				Help: "Número de produtos por categoria",
			},
			[]string{"category"},
		),

		// Métricas de Negócio - Valor Total dos Produtos
		ProductsTotalValue: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "products_total_value",
				Help: "Valor total de todos os produtos em estoque",
			},
		),

		// Métricas de Negócio - Preço Médio dos Produtos
		ProductsAveragePrice: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "products_average_price",
				Help: "Preço médio dos produtos",
			},
		),
	}
}

// RecordHTTPRequest registra uma requisição HTTP completa
func (m *Metrics) RecordHTTPRequest(method, endpoint, status string, duration float64, isError bool) {
	// Latency
	m.HTTPRequestDuration.WithLabelValues(method, endpoint, status).Observe(duration)

	// Traffic
	m.HTTPRequestsTotal.WithLabelValues(method, endpoint, status).Inc()

	// Errors
	if isError {
		errorType := "client_error"
		if status[0] == '5' {
			errorType = "server_error"
		}
		m.HTTPRequestErrors.WithLabelValues(method, endpoint, status, errorType).Inc()
	}
}

// IncrementProductsCreated incrementa o contador de produtos criados
func (m *Metrics) IncrementProductsCreated() {
	m.ProductsCreated.Inc()
	m.ProductsTotal.Inc()
}

// UpdateProductsByCategory atualiza a contagem de produtos por categoria
func (m *Metrics) UpdateProductsByCategory(category string, count float64) {
	m.ProductsByCategory.WithLabelValues(category).Set(count)
}

// UpdateProductsTotalValue atualiza o valor total dos produtos
func (m *Metrics) UpdateProductsTotalValue(totalValue float64) {
	m.ProductsTotalValue.Set(totalValue)
}

// UpdateProductsAveragePrice atualiza o preço médio dos produtos
func (m *Metrics) UpdateProductsAveragePrice(avgPrice float64) {
	m.ProductsAveragePrice.Set(avgPrice)
}

// ResetProductsByCategory limpa as métricas de categoria (útil antes de recalcular)
func (m *Metrics) ResetProductsByCategory() {
	m.ProductsByCategory.Reset()
}
