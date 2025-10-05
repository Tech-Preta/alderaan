-- Migration: Criar tabela de produtos
-- Autor: Sistema Alderaan
-- Data: 2024-10-04

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    sku INTEGER NOT NULL UNIQUE,
    price INTEGER NOT NULL CHECK (price > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabela para categorias de produtos (relacionamento many-to-many)
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de junção entre produtos e categorias
CREATE TABLE IF NOT EXISTS product_categories (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);

-- Índices para melhor performance
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_product_categories_product_id ON product_categories(product_id);
CREATE INDEX idx_product_categories_category_id ON product_categories(category_id);

-- Função para atualizar updated_at automaticamente
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger para atualizar updated_at
CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comentários nas tabelas
COMMENT ON TABLE products IS 'Tabela principal de produtos';
COMMENT ON TABLE categories IS 'Categorias disponíveis para produtos';
COMMENT ON TABLE product_categories IS 'Relacionamento many-to-many entre produtos e categorias';

COMMENT ON COLUMN products.name IS 'Nome do produto (único)';
COMMENT ON COLUMN products.sku IS 'Stock Keeping Unit - código único do produto';
COMMENT ON COLUMN products.price IS 'Preço em centavos (evita problemas com ponto flutuante)';

