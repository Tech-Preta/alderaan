# 🗄️ Database - PostgreSQL + Flyway

Este guia explica como o banco de dados PostgreSQL é configurado e gerenciado com o Flyway.

## 📊 Schema do Banco de Dados

### **Tabelas**

#### **1. products** (Tabela principal)
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    sku INTEGER NOT NULL UNIQUE,
    price INTEGER NOT NULL CHECK (price > 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Campos:**
- `id`: ID autoincrementado
- `name`: Nome do produto (único)
- `sku`: Stock Keeping Unit - código único
- `price`: Preço em centavos (evita problemas com ponto flutuante)
- `created_at` / `updated_at`: Timestamps automáticos

#### **2. categories** (Categorias)
```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### **3. product_categories** (Relacionamento many-to-many)
```sql
CREATE TABLE product_categories (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    PRIMARY KEY (product_id, category_id)
);
```

### **Índices para Performance**

```sql
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_sku ON products(sku);
CREATE INDEX idx_products_created_at ON products(created_at);
CREATE INDEX idx_categories_name ON categories(name);
CREATE INDEX idx_product_categories_product_id ON product_categories(product_id);
CREATE INDEX idx_product_categories_category_id ON product_categories(category_id);
```

### **Triggers**

```sql
-- Atualiza updated_at automaticamente ao modificar produto
CREATE TRIGGER update_products_updated_at 
BEFORE UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## 🔄 Flyway Migrations

### **Estrutura de Migrations**

```
db/migrations/
├── conf/
│   └── flyway.conf                       # Configuração do Flyway
├── V1__create_products_tables.sql        # Migration versionada
├── U1__rollback_products_tables.sql      # Undo migration
└── R__seed_data.sql                      # Repeatable migration (seed)
```

### **Status das Migrations**

```bash
# Ver histórico de migrations
make db-migrate-info

# Resultado:
# +-----------+---------+------------------------+------+---------------------+---------+
# | Category  | Version | Description            | Type | Installed On        | State   |
# +-----------+---------+------------------------+------+---------------------+---------+
# | Versioned | 1       | create products tables | SQL  | 2024-10-04 22:00:00 | Success |
# | Repeatable|         | seed data             | SQL  | 2024-10-04 22:00:05 | Success |
# +-----------+---------+------------------------+------+---------------------+---------+
```

### **Aplicar Migrations**

```bash
# Aplicar todas as migrations pendentes
make db-migrate

# Validar migrations
make db-migrate-validate

# Ver detalhes
docker-compose logs flyway
```

## 🌱 Seed Data

O arquivo `R__seed_data.sql` popula o banco com dados de exemplo:

- **9 produtos** em diversas categorias
- **5 categorias**: Electronics, Books, Clothing, Home, Sports

```bash
# Dados são inseridos automaticamente via Flyway (repeatable migration)
# Mas você pode executar manualmente se necessário:
make db-seed
```

## 🔌 Conexão

### **Variáveis de Ambiente**

```bash
DB_HOST=localhost      # ou 'postgres' dentro do Docker
DB_PORT=5432
DB_USER=alderaan
DB_PASSWORD=alderaan123
DB_NAME=alderaan_db
DB_SSLMODE=disable
```

### **String de Conexão**

```
postgresql://alderaan:alderaan123@localhost:5432/alderaan_db?sslmode=disable
```

### **Conectar via CLI**

```bash
# Via Docker Compose
make db-connect

# Via psql direto
psql -h localhost -p 5432 -U alderaan -d alderaan_db

# Senha: alderaan123
```

## 📝 Queries Úteis

### **Ver Todos os Produtos**

```sql
SELECT 
    p.id,
    p.name,
    p.sku,
    p.price / 100.0 as price_dollars,
    ARRAY_AGG(c.name) as categories,
    p.created_at
FROM products p
LEFT JOIN product_categories pc ON p.id = pc.product_id
LEFT JOIN categories c ON pc.category_id = c.id
GROUP BY p.id, p.name, p.sku, p.price, p.created_at
ORDER BY p.created_at DESC;
```

### **Produtos por Categoria**

```sql
SELECT 
    c.name as category,
    COUNT(p.id) as total_products,
    AVG(p.price / 100.0) as avg_price_dollars
FROM categories c
LEFT JOIN product_categories pc ON c.id = pc.category_id
LEFT JOIN products p ON pc.product_id = p.id
GROUP BY c.name
ORDER BY total_products DESC;
```

### **Produtos Mais Caros**

```sql
SELECT 
    name,
    sku,
    price / 100.0 as price_dollars
FROM products
ORDER BY price DESC
LIMIT 10;
```

### **Estatísticas Gerais**

```sql
SELECT 
    COUNT(*) as total_products,
    SUM(price) / 100.0 as total_value_dollars,
    AVG(price) / 100.0 as avg_price_dollars,
    MIN(price) / 100.0 as min_price_dollars,
    MAX(price) / 100.0 as max_price_dollars
FROM products;
```

### **Histórico de Migrations do Flyway**

```sql
SELECT 
    installed_rank,
    version,
    description,
    type,
    script,
    installed_on,
    execution_time,
    success
FROM flyway_schema_history
ORDER BY installed_rank;
```

## 🛠️ Comandos do Makefile

### **Gerenciamento do PostgreSQL**

```bash
# Iniciar PostgreSQL
make db-up

# Parar PostgreSQL
make db-down

# Ver logs
make db-logs

# Conectar ao banco
make db-connect
```

### **Flyway Migrations**

```bash
# Aplicar migrations
make db-migrate

# Ver status das migrations
make db-migrate-info

# Validar migrations
make db-migrate-validate

# Reparar histórico (se necessário)
make db-migrate-repair

# Criar baseline
make db-migrate-baseline
```

### **Gerenciamento de Dados**

```bash
# Popular com seed data (já é feito automaticamente)
make db-seed

# Limpar apenas dados (mantém estrutura)
make db-clean

# Limpar tudo incluindo histórico Flyway
make db-clean-all

# Reset completo (limpa e recria tudo)
make db-reset
```

## 🔒 Segurança

### **Credenciais em Produção**

```bash
# ❌ NUNCA faça isso em produção
DB_PASSWORD=alderaan123

# ✅ Use secrets management
# - AWS Secrets Manager
# - HashiCorp Vault
# - Kubernetes Secrets
# - Docker Secrets
```

### **Backup**

```bash
# Backup completo
docker exec alderaan-postgres pg_dump -U alderaan alderaan_db > backup.sql

# Backup com timestamp
docker exec alderaan-postgres pg_dump -U alderaan alderaan_db > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore
docker exec -i alderaan-postgres psql -U alderaan -d alderaan_db < backup.sql
```

### **SSL/TLS**

```bash
# Em produção, sempre use SSL
DB_SSLMODE=require

# Com certificado específico
DB_SSLMODE=verify-full
DB_SSLROOTCERT=/path/to/ca-cert.pem
```

## 🐳 Docker Volumes

### **Persistência de Dados**

```yaml
volumes:
  postgres_data:
    driver: local
```

Os dados são persistidos em:
- **macOS/Linux**: `/var/lib/docker/volumes/alderaan_postgres_data/_data`
- **Windows**: `C:\ProgramData\Docker\volumes\alderaan_postgres_data\_data`

### **Limpar Volumes**

```bash
# Para a plataforma
make platform-down

# Remove volumes (⚠️ PERDE TODOS OS DADOS)
docker volume rm alderaan_postgres_data

# Reinicia do zero
make platform-up
```

## 🔍 Troubleshooting

### **Erro: "relation does not exist"**

```bash
# Verificar se migrations foram aplicadas
make db-migrate-info

# Se não, aplicar migrations
make db-migrate
```

### **Erro: "password authentication failed"**

```bash
# Verificar credenciais no docker-compose.yml
cat docker-compose.yml | grep POSTGRES_

# Resetar senha
docker-compose down
docker volume rm alderaan_postgres_data
docker-compose up -d postgres
```

### **Banco lento**

```sql
-- Ver queries lentas
SELECT 
    pid,
    now() - query_start as duration,
    query
FROM pg_stat_activity
WHERE state = 'active'
  AND now() - query_start > interval '5 seconds'
ORDER BY duration DESC;

-- Ver índices não utilizados
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY schemaname, tablename;
```

### **Espaço em disco**

```sql
-- Ver tamanho das tabelas
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Ver tamanho do banco
SELECT pg_size_pretty(pg_database_size('alderaan_db'));
```

## 📚 Recursos

- **PostgreSQL Docs**: https://www.postgresql.org/docs/
- **Flyway Docs**: https://flywaydb.org/documentation/
- **Flyway Migrations**: [docs/09-flyway-migrations.md](../docs/09-flyway-migrations.md)
- **Docker PostgreSQL**: https://hub.docker.com/_/postgres

---

**Voltar para:** [README Principal](../README.md) | [Documentação](../docs/README.md)
