package product_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
	product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
	"github.com/williamkoller/golang-domain-driven-design/internal/metrics"
	shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
)

type ProductHandler struct {
	repo       product_repository.IProductRepository
	dispatcher *shared_events.EventDispatcher
	metrics    *metrics.Metrics
}

func NewProductHandler(repo product_repository.IProductRepository, dispatcher *shared_events.EventDispatcher, m *metrics.Metrics) *ProductHandler {
	return &ProductHandler{repo, dispatcher, m}
}

// CreateProductInput representa os dados de entrada para criar um produto
type CreateProductInput struct {
	Name       string   `json:"name" binding:"required" example:"Notebook"`
	Sku        int      `json:"sku" binding:"required" example:"12345"`
	Categories []string `json:"categories" binding:"required" example:"Eletrônicos,Computadores"`
	Price      int      `json:"price" binding:"required" example:"3500"`
}

// ErrorResponse representa uma resposta de erro
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// Create godoc
// @Summary      Criar um novo produto
// @Description  Cria um novo produto com nome, SKU, categorias e preço
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      CreateProductInput  true  "Dados do produto"
// @Success      201      {object}  product_entity.Product
// @Failure      400      {object}  ErrorResponse
// @Failure      409      {object}  ErrorResponse
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var input CreateProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, _, err := product_entity.NewProduct(input.Name, input.Sku, input.Categories, input.Price, h.dispatcher)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Add(*product); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Atualizar métricas de negócio
	h.metrics.IncrementProductsCreated()
	h.updateBusinessMetrics()

	c.JSON(http.StatusCreated, product)
}

// FindAll godoc
// @Summary      Listar todos os produtos
// @Description  Retorna uma lista com todos os produtos cadastrados
// @Tags         products
// @Produce      json
// @Success      200  {array}   product_entity.Product
// @Router       /products [get]
func (h *ProductHandler) FindAll(c *gin.Context) {
	products, _ := h.repo.Find()
	c.JSON(http.StatusOK, products)
}

// FindOne godoc
// @Summary      Buscar produto por nome
// @Description  Retorna um produto específico pelo nome
// @Tags         products
// @Produce      json
// @Param        name  path      string  true  "Nome do produto"
// @Success      200   {object}  product_entity.Product
// @Failure      404   {object}  ErrorResponse
// @Router       /products/{name} [get]
func (h *ProductHandler) FindOne(c *gin.Context) {
	name := c.Param("name")
	product, err := h.repo.FindOne(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// updateBusinessMetrics atualiza todas as métricas de negócio
func (h *ProductHandler) updateBusinessMetrics() {
	repoMetrics := h.repo.GetMetrics()

	// Atualizar total de produtos
	h.metrics.ProductsTotal.Set(float64(repoMetrics.TotalProducts))

	// Atualizar valor total
	h.metrics.UpdateProductsTotalValue(repoMetrics.TotalValue)

	// Atualizar preço médio
	h.metrics.UpdateProductsAveragePrice(repoMetrics.AveragePrice)

	// Atualizar produtos por categoria
	h.metrics.ResetProductsByCategory()
	for category, count := range repoMetrics.ProductsByCategory {
		h.metrics.UpdateProductsByCategory(category, float64(count))
	}
}
