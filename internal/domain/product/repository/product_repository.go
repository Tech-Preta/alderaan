package product_repository

import (
	"errors"
	"sync"

	product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
)

type IProductRepository interface {
	Add(product product_entity.Product) error
	Find() ([]product_entity.Product, error)
	FindOne(name string) (product_entity.Product, error)
	GetMetrics() RepositoryMetrics
}

// RepositoryMetrics contém métricas calculadas do repositório
type RepositoryMetrics struct {
	TotalProducts      int
	TotalValue         float64
	AveragePrice       float64
	ProductsByCategory map[string]int
}

type ProductRepository struct {
	data map[string]product_entity.Product
	mu   sync.RWMutex
}

func NewRepository() *ProductRepository {
	return &ProductRepository{
		data: make(map[string]product_entity.Product),
	}
}

func (r *ProductRepository) Add(product product_entity.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[product.Name]; exists {
		return errors.New("product already exists")
	}

	r.data[product.Name] = product

	return nil
}

func (r *ProductRepository) Find() ([]product_entity.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	products := make([]product_entity.Product, 0, len(r.data))

	for _, p := range r.data {
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) FindOne(name string) (product_entity.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	product, exists := r.data[name]

	if !exists {
		return product_entity.Product{}, errors.New("product not found")
	}

	return product, nil
}

// GetMetrics calcula e retorna métricas do repositório
func (r *ProductRepository) GetMetrics() RepositoryMetrics {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics := RepositoryMetrics{
		TotalProducts:      len(r.data),
		ProductsByCategory: make(map[string]int),
	}

	totalValue := 0
	for _, product := range r.data {
		// Valor total
		totalValue += product.Price

		// Produtos por categoria
		for _, category := range product.Categories {
			metrics.ProductsByCategory[category]++
		}
	}

	metrics.TotalValue = float64(totalValue)
	if metrics.TotalProducts > 0 {
		metrics.AveragePrice = metrics.TotalValue / float64(metrics.TotalProducts)
	}

	return metrics
}
