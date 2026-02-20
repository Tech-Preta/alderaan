package product_repository

import (
	"sync"
	"testing"

	product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
)

func TestNewRepository(t *testing.T) {
	repo := NewRepository()

	if repo == nil {
		t.Fatal("NewRepository() returned nil")
	}

	if repo.data == nil {
		t.Error("NewRepository() data map is nil")
	}

	if len(repo.data) != 0 {
		t.Errorf("NewRepository() data map not empty, got %d items", len(repo.data))
	}
}

func TestProductRepository_Add(t *testing.T) {
	tests := []struct {
		name    string
		product product_entity.Product
		wantErr bool
		errMsg  string
		setup   func(*ProductRepository)
	}{
		{
			name: "add new product successfully",
			product: product_entity.Product{
				Name:       "Notebook",
				Sku:        123,
				Categories: []string{"Electronics"},
				Price:      3500,
			},
			wantErr: false,
		},
		{
			name: "add duplicate product",
			product: product_entity.Product{
				Name:       "Notebook",
				Sku:        123,
				Categories: []string{"Electronics"},
				Price:      3500,
			},
			wantErr: true,
			errMsg:  "product already exists",
			setup: func(r *ProductRepository) {
				r.data["Notebook"] = product_entity.Product{
					Name:       "Notebook",
					Sku:        456,
					Categories: []string{"Electronics"},
					Price:      2500,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewRepository()

			if tt.setup != nil {
				tt.setup(repo)
			}

			err := repo.Add(tt.product)

			if tt.wantErr {
				if err == nil {
					t.Error("Add() expected error, got nil")
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Add() error = %v, want %v", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Add() unexpected error = %v", err)
					return
				}

				// Verificar se o produto foi adicionado
				products, _ := repo.Find()
				if len(products) != 1 {
					t.Errorf("Add() products count = %d, want 1", len(products))
				}
			}
		})
	}
}

func TestProductRepository_Find(t *testing.T) {
	repo := NewRepository()

	// Repositório vazio
	t.Run("empty repository", func(t *testing.T) {
		products, err := repo.Find()
		if err != nil {
			t.Errorf("Find() unexpected error = %v", err)
		}
		if len(products) != 0 {
			t.Errorf("Find() expected 0 products, got %d", len(products))
		}
	})

	// Adicionar produtos
	products := []product_entity.Product{
		{Name: "Product1", Sku: 1, Categories: []string{"Cat1"}, Price: 100},
		{Name: "Product2", Sku: 2, Categories: []string{"Cat2"}, Price: 200},
		{Name: "Product3", Sku: 3, Categories: []string{"Cat3"}, Price: 300},
	}

	for _, p := range products {
		_ = repo.Add(p)
	}

	// Repositório com produtos
	t.Run("repository with products", func(t *testing.T) {
		found, err := repo.Find()
		if err != nil {
			t.Errorf("Find() unexpected error = %v", err)
		}
		if len(found) != 3 {
			t.Errorf("Find() expected 3 products, got %d", len(found))
		}
	})
}

func TestProductRepository_FindOne(t *testing.T) {
	repo := NewRepository()

	// Adicionar um produto
	product := product_entity.Product{
		Name:       "Test Product",
		Sku:        999,
		Categories: []string{"Test"},
		Price:      500,
	}
	_ = repo.Add(product)

	tests := []struct {
		name     string
		findName string
		wantErr  bool
		wantSku  int
	}{
		{
			name:     "find existing product",
			findName: "Test Product",
			wantErr:  false,
			wantSku:  999,
		},
		{
			name:     "find non-existing product",
			findName: "Non Existing",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := repo.FindOne(tt.findName)

			if tt.wantErr {
				if err == nil {
					t.Error("FindOne() expected error, got nil")
				}
				if err != nil && err.Error() != "product not found" {
					t.Errorf("FindOne() error = %v, want 'product not found'", err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("FindOne() unexpected error = %v", err)
					return
				}
				if found.Sku != tt.wantSku {
					t.Errorf("FindOne() Sku = %v, want %v", found.Sku, tt.wantSku)
				}
			}
		})
	}
}

func TestProductRepository_GetMetrics(t *testing.T) {
	repo := NewRepository()

	t.Run("empty repository metrics", func(t *testing.T) {
		metrics := repo.GetMetrics()

		if metrics.TotalProducts != 0 {
			t.Errorf("GetMetrics() TotalProducts = %d, want 0", metrics.TotalProducts)
		}
		if metrics.TotalValue != 0 {
			t.Errorf("GetMetrics() TotalValue = %f, want 0", metrics.TotalValue)
		}
		if metrics.AveragePrice != 0 {
			t.Errorf("GetMetrics() AveragePrice = %f, want 0", metrics.AveragePrice)
		}
	})

	// Adicionar produtos
	products := []product_entity.Product{
		{Name: "Product1", Sku: 1, Categories: []string{"Electronics", "Computers"}, Price: 1000},
		{Name: "Product2", Sku: 2, Categories: []string{"Electronics"}, Price: 2000},
		{Name: "Product3", Sku: 3, Categories: []string{"Books"}, Price: 500},
	}

	for _, p := range products {
		_ = repo.Add(p)
	}

	t.Run("repository with products metrics", func(t *testing.T) {
		metrics := repo.GetMetrics()

		if metrics.TotalProducts != 3 {
			t.Errorf("GetMetrics() TotalProducts = %d, want 3", metrics.TotalProducts)
		}

		expectedTotal := 3500.0
		if metrics.TotalValue != expectedTotal {
			t.Errorf("GetMetrics() TotalValue = %f, want %f", metrics.TotalValue, expectedTotal)
		}

		expectedAvg := 3500.0 / 3
		if metrics.AveragePrice != expectedAvg {
			t.Errorf("GetMetrics() AveragePrice = %f, want %f", metrics.AveragePrice, expectedAvg)
		}

		// Verificar contagem por categoria
		if metrics.ProductsByCategory["Electronics"] != 2 {
			t.Errorf("GetMetrics() Electronics count = %d, want 2", metrics.ProductsByCategory["Electronics"])
		}
		if metrics.ProductsByCategory["Computers"] != 1 {
			t.Errorf("GetMetrics() Computers count = %d, want 1", metrics.ProductsByCategory["Computers"])
		}
		if metrics.ProductsByCategory["Books"] != 1 {
			t.Errorf("GetMetrics() Books count = %d, want 1", metrics.ProductsByCategory["Books"])
		}
	})
}

// Teste de concorrência
func TestProductRepository_ConcurrentAccess(t *testing.T) {
	repo := NewRepository()
	var wg sync.WaitGroup

	// Adicionar produtos concorrentemente
	numGoroutines := 100
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			product := product_entity.Product{
				Name:       "Product" + string(rune(id)),
				Sku:        id,
				Categories: []string{"Test"},
				Price:      100,
			}
			_ = repo.Add(product)
		}(i)
	}

	wg.Wait()

	// Ler produtos concorrentemente
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			_, _ = repo.Find()
		}()
	}

	wg.Wait()

	// Verificar que não houve race conditions
	products, _ := repo.Find()
	if len(products) == 0 {
		t.Error("ConcurrentAccess() no products added")
	}
}

func TestProductRepository_ConcurrentFindAndAdd(t *testing.T) {
	repo := NewRepository()
	var wg sync.WaitGroup

	// Mix de operações concorrentes
	numOperations := 50
	wg.Add(numOperations * 2)

	// Adicionar
	for i := 0; i < numOperations; i++ {
		go func(id int) {
			defer wg.Done()
			product := product_entity.Product{
				Name:       "ConcurrentProduct" + string(rune(id)),
				Sku:        id + 1000,
				Categories: []string{"Concurrent"},
				Price:      id * 10,
			}
			_ = repo.Add(product)
		}(i)
	}

	// Ler
	for i := 0; i < numOperations; i++ {
		go func() {
			defer wg.Done()
			_ = repo.GetMetrics()
		}()
	}

	wg.Wait()
}

// Benchmarks
func BenchmarkProductRepository_Add(b *testing.B) {
	product := product_entity.Product{
		Name:       "Benchmark Product",
		Sku:        12345,
		Categories: []string{"Benchmark"},
		Price:      999,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := NewRepository() // Reset para evitar duplicatas
		_ = repo.Add(product)
	}
}

func BenchmarkProductRepository_Find(b *testing.B) {
	repo := NewRepository()

	// Adicionar alguns produtos primeiro
	for i := 0; i < 100; i++ {
		product := product_entity.Product{
			Name:       "Product" + string(rune(i)),
			Sku:        i,
			Categories: []string{"Test"},
			Price:      100,
		}
		_ = repo.Add(product)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Find()
	}
}

func BenchmarkProductRepository_GetMetrics(b *testing.B) {
	repo := NewRepository()

	// Adicionar produtos
	for i := 0; i < 100; i++ {
		product := product_entity.Product{
			Name:       "Product" + string(rune(i)),
			Sku:        i,
			Categories: []string{"Cat1", "Cat2"},
			Price:      i * 10,
		}
		_ = repo.Add(product)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.GetMetrics()
	}
}
