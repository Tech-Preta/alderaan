-- Migration Rollback: Remover tabela de produtos

DROP TRIGGER IF EXISTS update_products_updated_at ON products;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_product_categories_category_id;
DROP INDEX IF EXISTS idx_product_categories_product_id;
DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_products_created_at;
DROP INDEX IF EXISTS idx_products_sku;
DROP INDEX IF EXISTS idx_products_name;

DROP TABLE IF EXISTS product_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS products;

