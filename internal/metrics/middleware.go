package metrics

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware é um middleware Gin que coleta métricas de requisições HTTP
func PrometheusMiddleware(metrics *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Incrementa requisições em andamento (Saturation)
		metrics.InFlightRequests.Inc()
		defer metrics.InFlightRequests.Dec()

		// Captura o tempo de início
		start := time.Now()

		// Processa a requisição
		c.Next()

		// Calcula a duração
		duration := time.Since(start).Seconds()

		// Extrai informações da requisição
		method := c.Request.Method
		endpoint := normalizeEndpoint(c.FullPath())
		status := strconv.Itoa(c.Writer.Status())
		statusCode := c.Writer.Status()

		// Determina se é um erro
		isError := statusCode >= 400

		// Registra as métricas
		metrics.RecordHTTPRequest(method, endpoint, status, duration, isError)
	}
}

// normalizeEndpoint normaliza o endpoint para evitar cardinalidade alta
// Exemplo: /api/v1/products/:name -> /api/v1/products/{name}
func normalizeEndpoint(path string) string {
	if path == "" {
		return "unknown"
	}

	// Substitui :param por {param} para consistência
	normalized := strings.ReplaceAll(path, ":", "{")
	normalized = strings.ReplaceAll(normalized, "}", "")

	// Se contém {, adiciona } no final do parâmetro
	parts := strings.Split(normalized, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, "{") && !strings.HasSuffix(part, "}") {
			parts[i] = part + "}"
		}
	}

	return strings.Join(parts, "/")
}
