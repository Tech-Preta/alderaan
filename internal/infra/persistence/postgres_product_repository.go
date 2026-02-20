package persistence

import (
	"database/sql"
	"errors"
	"fmt"

	product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
	product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

// Add adiciona um novo produto ao banco de dados
func (r *PostgresProductRepository) Add(product product_entity.Product) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Inserir produto
	var productID int
	err = tx.QueryRow(`
		INSERT INTO products (name, sku, price)
		VALUES ($1, $2, $3)
		RETURNING id
	`, product.Name, product.Sku, product.Price).Scan(&productID)

	if err != nil {
		return fmt.Errorf("erro ao inserir produto: %w", err)
	}

	// Inserir categorias e relacionamentos
	for _, categoryName := range product.Categories {
		var categoryID int

		// Inserir categoria se não existir (ou pegar ID se já existe)
		err = tx.QueryRow(`
			INSERT INTO categories (name)
			VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, categoryName).Scan(&categoryID)

		if err != nil {
			return fmt.Errorf("erro ao inserir categoria: %w", err)
		}

		// Criar relacionamento produto-categoria
		_, err = tx.Exec(`
			INSERT INTO product_categories (product_id, category_id)
			VALUES ($1, $2)
		`, productID, categoryID)

		if err != nil {
			return fmt.Errorf("erro ao associar categoria ao produto: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("erro ao commitar transação: %w", err)
	}

	return nil
}

// Find retorna todos os produtos
func (r *PostgresProductRepository) Find() ([]product_entity.Product, error) {
	rows, err := r.db.Query(`
		SELECT DISTINCT p.id, p.name, p.sku, p.price
		FROM products p
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %w", err)
	}
	defer rows.Close()

	var products []product_entity.Product

	for rows.Next() {
		var (
			id    int
			name  string
			sku   int
			price int
		)

		if err := rows.Scan(&id, &name, &sku, &price); err != nil {
			return nil, fmt.Errorf("erro ao escanear produto: %w", err)
		}

		// Buscar categorias do produto
		categories, err := r.getProductCategories(id)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar categorias do produto %d: %w", id, err)
		}

		products = append(products, product_entity.Product{
			Name:       name,
			Sku:        sku,
			Categories: categories,
			Price:      price,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar produtos: %w", err)
	}

	return products, nil
}

// FindOne busca um produto pelo nome
func (r *PostgresProductRepository) FindOne(name string) (product_entity.Product, error) {
	var (
		id    int
		sku   int
		price int
	)

	err := r.db.QueryRow(`
		SELECT id, name, sku, price
		FROM products
		WHERE name = $1
	`, name).Scan(&id, &name, &sku, &price)

	if err == sql.ErrNoRows {
		return product_entity.Product{}, errors.New("product not found")
	}
	if err != nil {
		return product_entity.Product{}, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	// Buscar categorias do produto
	categories, err := r.getProductCategories(id)
	if err != nil {
		return product_entity.Product{}, fmt.Errorf("erro ao buscar categorias: %w", err)
	}

	return product_entity.Product{
		Name:       name,
		Sku:        sku,
		Categories: categories,
		Price:      price,
	}, nil
}

// GetMetrics retorna métricas do repositório
func (r *PostgresProductRepository) GetMetrics() product_repository.RepositoryMetrics {
	metrics := product_repository.RepositoryMetrics{
		ProductsByCategory: make(map[string]int),
	}

	// Total de produtos
	_ = r.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&metrics.TotalProducts)

	// Valor total e preço médio
	var totalValue sql.NullFloat64
	_ = r.db.QueryRow(`
		SELECT 
			SUM(price * stock), 
			AVG(price) 
		FROM products
	`).Scan(&totalValue, &metrics.AveragePrice)

	if totalValue.Valid {
		metrics.TotalValue = totalValue.Float64
	}

	// Produtos por categoria
	rows, err := r.db.Query(`
		SELECT c.name, COUNT(pc.product_id)
		FROM categories c
		LEFT JOIN product_categories pc ON c.id = pc.category_id
		GROUP BY c.name
		ORDER BY COUNT(pc.product_id) DESC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var category string
			var count int
			if err := rows.Scan(&category, &count); err == nil {
				metrics.ProductsByCategory[category] = count
			}
		}
	}

	return metrics
}

// getProductCategories retorna as categorias de um produto
func (r *PostgresProductRepository) getProductCategories(productID int) ([]string, error) {
	rows, err := r.db.Query(`
		SELECT c.name
		FROM categories c
		INNER JOIN product_categories pc ON c.id = pc.category_id
		WHERE pc.product_id = $1
		ORDER BY c.name
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}
