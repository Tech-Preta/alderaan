-- Seed: Dados iniciais para desenvolvimento e testes

-- Inserir categorias iniciais
INSERT INTO categories (name) VALUES
    ('Eletrônicos'),
    ('Computadores'),
    ('Smartphones'),
    ('Periféricos'),
    ('Gaming'),
    ('Livros'),
    ('Tecnologia')
ON CONFLICT (name) DO NOTHING;

-- Inserir produtos de exemplo
INSERT INTO products (name, sku, price) VALUES
    ('Notebook Dell Inspiron', 12345, 350000),  -- R$ 3.500,00
    ('Mouse Gamer RGB', 23456, 25000),          -- R$ 250,00
    ('Teclado Mecânico', 34567, 45000),         -- R$ 450,00
    ('Monitor 4K LG', 45678, 180000),           -- R$ 1.800,00
    ('Webcam HD Logitech', 56789, 35000),       -- R$ 350,00
    ('Headset HyperX', 67890, 55000),           -- R$ 550,00
    ('SSD Samsung 1TB', 78901, 65000),          -- R$ 650,00
    ('Memória RAM 16GB', 89012, 40000)          -- R$ 400,00
ON CONFLICT (name) DO NOTHING;

-- Associar produtos às categorias
-- Notebook
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'Notebook Dell Inspiron' AND c.name IN ('Eletrônicos', 'Computadores')
ON CONFLICT DO NOTHING;

-- Mouse Gamer
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'Mouse Gamer RGB' AND c.name IN ('Periféricos', 'Gaming')
ON CONFLICT DO NOTHING;

-- Teclado Mecânico
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'Teclado Mecânico' AND c.name IN ('Periféricos', 'Gaming')
ON CONFLICT DO NOTHING;

-- Monitor
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'Monitor 4K LG' AND c.name IN ('Eletrônicos')
ON CONFLICT DO NOTHING;

-- Webcam
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'Webcam HD Logitech' AND c.name IN ('Periféricos')
ON CONFLICT DO NOTHING;

-- Headset
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'Headset HyperX' AND c.name IN ('Periféricos', 'Gaming')
ON CONFLICT DO NOTHING;

-- SSD
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'SSD Samsung 1TB' AND c.name IN ('Eletrônicos', 'Computadores')
ON CONFLICT DO NOTHING;

-- RAM
INSERT INTO product_categories (product_id, category_id)
SELECT p.id, c.id FROM products p, categories c
WHERE p.name = 'Memória RAM 16GB' AND c.name IN ('Computadores')
ON CONFLICT DO NOTHING;
