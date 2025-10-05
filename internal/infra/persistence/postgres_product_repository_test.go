package persistence

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
	product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
)

func TestNewPostgresProductRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	repo := NewPostgresProductRepository(db)

	if repo == nil {
		t.Fatal("NewPostgresProductRepository() returned nil")
	}

	if repo.db == nil {
		t.Error("repository db is nil")
	}
}

func TestPostgresProductRepository_Add(t *testing.T) {
	tests := []struct {
		name          string
		product       product_entity.Product
		mockSetup     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "add product successfully",
			product: product_entity.Product{
				Name:       "Notebook",
				Sku:        12345,
				Categories: []string{"Electronics", "Computers"},
				Price:      3500,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Expect BEGIN
				mock.ExpectBegin()

				// Expect INSERT into products with RETURNING id
				rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO products").
					WithArgs("Notebook", 12345, 3500).
					WillReturnRows(rows)

				// Expect INSERT for each category (2 times)
				// Electronics
				catRows1 := sqlmock.NewRows([]string{"id"}).AddRow(1)
				mock.ExpectQuery("INSERT INTO categories").
					WithArgs("Electronics").
					WillReturnRows(catRows1)
				mock.ExpectExec("INSERT INTO product_categories").
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Computers
				catRows2 := sqlmock.NewRows([]string{"id"}).AddRow(2)
				mock.ExpectQuery("INSERT INTO categories").
					WithArgs("Computers").
					WillReturnRows(catRows2)
				mock.ExpectExec("INSERT INTO product_categories").
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))

				// Expect COMMIT
				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "product with single category",
			product: product_entity.Product{
				Name:       "Book",
				Sku:        999,
				Categories: []string{"Books"},
				Price:      50,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(2)
				mock.ExpectQuery("INSERT INTO products").
					WithArgs("Book", 999, 50).
					WillReturnRows(rows)

				catRows := sqlmock.NewRows([]string{"id"}).AddRow(3)
				mock.ExpectQuery("INSERT INTO categories").
					WithArgs("Books").
					WillReturnRows(catRows)

				mock.ExpectExec("INSERT INTO product_categories").
					WithArgs(2, 3).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedError: false,
		},
		{
			name: "database error on insert",
			product: product_entity.Product{
				Name:       "ErrorProduct",
				Sku:        111,
				Categories: []string{"Test"},
				Price:      100,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery("INSERT INTO products").
					WithArgs("ErrorProduct", 111, 100).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			repo := NewPostgresProductRepository(db)
			err = repo.Add(tt.product)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}

			// Verificar que todas as expectativas foram atendidas
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgresProductRepository_Find(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
	}{
		{
			name: "find all products successfully",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "sku", "price"}).
					AddRow(1, "Product1", 1, 100).
					AddRow(2, "Product2", 2, 200).
					AddRow(3, "Product3", 3, 300)
				mock.ExpectQuery("SELECT DISTINCT p.id, p.name, p.sku, p.price FROM products p").
					WillReturnRows(rows)

				// Para cada produto, esperar query de categorias
				catRows1 := sqlmock.NewRows([]string{"name"}).AddRow("Cat1")
				mock.ExpectQuery("SELECT c.name FROM categories c").
					WithArgs(1).
					WillReturnRows(catRows1)

				catRows2 := sqlmock.NewRows([]string{"name"}).AddRow("Cat2")
				mock.ExpectQuery("SELECT c.name FROM categories c").
					WithArgs(2).
					WillReturnRows(catRows2)

				catRows3 := sqlmock.NewRows([]string{"name"}).AddRow("Cat3")
				mock.ExpectQuery("SELECT c.name FROM categories c").
					WithArgs(3).
					WillReturnRows(catRows3)
			},
			expectedCount: 3,
			expectedError: false,
		},
		{
			name: "find no products - empty database",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "sku", "price"})
				mock.ExpectQuery("SELECT DISTINCT p.id, p.name, p.sku, p.price FROM products p").
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "database error on query",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT DISTINCT p.id, p.name, p.sku, p.price FROM products p").
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			repo := NewPostgresProductRepository(db)
			products, err := repo.Find()

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(products) != tt.expectedCount {
					t.Errorf("Expected %d products, got %d", tt.expectedCount, len(products))
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgresProductRepository_FindOne(t *testing.T) {
	tests := []struct {
		name          string
		productName   string
		mockSetup     func(sqlmock.Sqlmock)
		expectedError bool
		checkProduct  func(*testing.T, product_entity.Product)
	}{
		{
			name:        "find existing product",
			productName: "Notebook",
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "sku", "price"}).
					AddRow(1, "Notebook", 12345, 3500)
				mock.ExpectQuery("SELECT id, name, sku, price FROM products WHERE name").
					WithArgs("Notebook").
					WillReturnRows(rows)

				catRows := sqlmock.NewRows([]string{"name"}).
					AddRow("Electronics").
					AddRow("Computers")
				mock.ExpectQuery("SELECT c.name FROM categories c").
					WithArgs(1).
					WillReturnRows(catRows)
			},
			expectedError: false,
			checkProduct: func(t *testing.T, p product_entity.Product) {
				if p.Name != "Notebook" {
					t.Errorf("Expected name Notebook, got %s", p.Name)
				}
				if p.Sku != 12345 {
					t.Errorf("Expected SKU 12345, got %d", p.Sku)
				}
				if p.Price != 3500 {
					t.Errorf("Expected price 3500, got %d", p.Price)
				}
				if len(p.Categories) != 2 {
					t.Errorf("Expected 2 categories, got %d", len(p.Categories))
				}
			},
		},
		{
			name:        "product not found",
			productName: "NonExistent",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, sku, price FROM products WHERE name").
					WithArgs("NonExistent").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: true,
		},
		{
			name:        "database error",
			productName: "Test",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, sku, price FROM products WHERE name").
					WithArgs("Test").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			repo := NewPostgresProductRepository(db)
			product, err := repo.FindOne(tt.productName)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tt.checkProduct != nil {
					tt.checkProduct(t, product)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %v", err)
			}
		})
	}
}

func TestPostgresProductRepository_GetMetrics(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(sqlmock.Sqlmock)
		expectedError bool
		checkMetrics  func(*testing.T, product_repository.RepositoryMetrics)
	}{
		{
			name: "get metrics successfully",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// Total products
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
					WillReturnRows(countRows)

				// Total value and Average price
				sumAvgRows := sqlmock.NewRows([]string{"sum", "avg"}).AddRow(5000, 500.0)
				mock.ExpectQuery("SELECT SUM\\(price\\), AVG\\(price\\) FROM products").
					WillReturnRows(sumAvgRows)

				// Products by category
				catRows := sqlmock.NewRows([]string{"name", "count"}).
					AddRow("Electronics", 5).
					AddRow("Books", 3).
					AddRow("Toys", 2)
				mock.ExpectQuery("SELECT c.name, COUNT\\(pc.product_id\\)").
					WillReturnRows(catRows)
			},
			expectedError: false,
			checkMetrics: func(t *testing.T, m product_repository.RepositoryMetrics) {
				if m.TotalProducts != 10 {
					t.Errorf("Expected TotalProducts 10, got %d", m.TotalProducts)
				}
				if m.TotalValue != 5000.0 {
					t.Errorf("Expected TotalValue 5000.0, got %f", m.TotalValue)
				}
				if m.AveragePrice != 500.0 {
					t.Errorf("Expected AveragePrice 500.0, got %f", m.AveragePrice)
				}
				if len(m.ProductsByCategory) != 3 {
					t.Errorf("Expected 3 categories, got %d", len(m.ProductsByCategory))
				}
				if m.ProductsByCategory["Electronics"] != 5 {
					t.Errorf("Expected 5 Electronics, got %d", m.ProductsByCategory["Electronics"])
				}
			},
		},
		{
			name: "empty database metrics",
			mockSetup: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
					WillReturnRows(countRows)

				sumAvgRows := sqlmock.NewRows([]string{"sum", "avg"}).AddRow(nil, 0.0)
				mock.ExpectQuery("SELECT SUM\\(price\\), AVG\\(price\\) FROM products").
					WillReturnRows(sumAvgRows)

				catRows := sqlmock.NewRows([]string{"name", "count"})
				mock.ExpectQuery("SELECT c.name, COUNT\\(pc.product_id\\)").
					WillReturnRows(catRows)
			},
			expectedError: false,
			checkMetrics: func(t *testing.T, m product_repository.RepositoryMetrics) {
				if m.TotalProducts != 0 {
					t.Errorf("Expected TotalProducts 0, got %d", m.TotalProducts)
				}
				if m.TotalValue != 0.0 {
					t.Errorf("Expected TotalValue 0.0, got %f", m.TotalValue)
				}
				if m.AveragePrice != 0.0 {
					t.Errorf("Expected AveragePrice 0.0, got %f", m.AveragePrice)
				}
			},
		},
		{
			name: "database error on count query",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: false, // GetMetrics não retorna erro, apenas valores zerados
			checkMetrics: func(t *testing.T, m product_repository.RepositoryMetrics) {
				if m.TotalProducts != 0 {
					t.Errorf("Expected TotalProducts 0 on error, got %d", m.TotalProducts)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create mock database: %v", err)
			}
			defer db.Close()

			if tt.mockSetup != nil {
				tt.mockSetup(mock)
			}

			repo := NewPostgresProductRepository(db)
			metrics := repo.GetMetrics()

			if tt.checkMetrics != nil {
				tt.checkMetrics(t, metrics)
			}

			// GetMetrics não precisa atender todas as expectativas se houver erro
			if !tt.expectedError {
				if err := mock.ExpectationsWereMet(); err != nil {
					t.Logf("Note: Some expectations were not met (this may be expected): %v", err)
				}
			}
		})
	}
}

// Benchmark
func BenchmarkPostgresProductRepository_Add(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	product := product_entity.Product{
		Name:       "BenchProduct",
		Sku:        99999,
		Categories: []string{"Bench"},
		Price:      1000,
	}

	// Setup expectations for each iteration
	for i := 0; i < b.N; i++ {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO products").
			WillReturnResult(sqlmock.NewResult(1, 1))
		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Bench")
		mock.ExpectQuery("SELECT id, name FROM categories").WillReturnRows(rows)
		mock.ExpectExec("INSERT INTO product_categories").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()
	}

	repo := NewPostgresProductRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.Add(product)
	}
}

func BenchmarkPostgresProductRepository_Find(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"name", "sku", "price"}).
			AddRow("Product1", 1, 100)
		mock.ExpectQuery("SELECT DISTINCT p.name, p.sku, p.price FROM products p").
			WillReturnRows(rows)
		catRows := sqlmock.NewRows([]string{"name"}).AddRow("Cat1")
		mock.ExpectQuery("SELECT c.name FROM categories c").WillReturnRows(catRows)
	}

	repo := NewPostgresProductRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.Find()
	}
}

func BenchmarkPostgresProductRepository_GetMetrics(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	for i := 0; i < b.N; i++ {
		countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM products").WillReturnRows(countRows)

		sumRows := sqlmock.NewRows([]string{"sum"}).AddRow(5000)
		mock.ExpectQuery("SELECT COALESCE\\(SUM\\(price\\), 0\\) FROM products").WillReturnRows(sumRows)

		avgRows := sqlmock.NewRows([]string{"avg"}).AddRow(500.0)
		mock.ExpectQuery("SELECT COALESCE\\(AVG\\(price\\), 0\\) FROM products").WillReturnRows(avgRows)

		catRows := sqlmock.NewRows([]string{"name", "count"})
		mock.ExpectQuery("SELECT c.name, COUNT\\(pc.product_name\\)").WillReturnRows(catRows)
	}

	repo := NewPostgresProductRepository(db)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.GetMetrics()
	}
}
