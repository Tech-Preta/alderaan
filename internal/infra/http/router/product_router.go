package product_router

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	product_handlers "github.com/williamkoller/golang-domain-driven-design/internal/infra/http/handlers"
	"github.com/williamkoller/golang-domain-driven-design/internal/metrics"
)

func SetupProductRouter(productHandler *product_handlers.ProductHandler, m *metrics.Metrics) *gin.Engine {
	r := gin.New()

	// Middleware padrão do Gin
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Middleware de métricas Prometheus (Golden Signals)
	r.Use(metrics.PrometheusMiddleware(m))

	// Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/products", productHandler.Create)
		v1.GET("/products", productHandler.FindAll)
		v1.GET("/products/:name", productHandler.FindOne)
	}

	return r
}
