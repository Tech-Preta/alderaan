# 🔄 Flyway - Database Migrations

Este guia explica como gerenciamos migrações de banco de dados usando o **Flyway**, uma ferramenta profissional e robusta para versionamento de schemas.

## 📚 O que é Flyway?

**Flyway** é uma ferramenta open-source para migração e versionamento de banco de dados. Ele permite:

- ✅ **Versionamento**: Controle de versão do schema do banco
- ✅ **Rastreabilidade**: Histórico completo de todas as mudanças
- ✅ **Automação**: Migrations automáticas no deploy
- ✅ **Segurança**: Validação e checksums de migrations
- ✅ **Rollback**: Suporte a undo migrations (U migrations)
- ✅ **Repeatable**: Migrations que rodam sempre (views, procedures)

## 🗂️ Estrutura de Migrations

```
db/migrations/
├── conf/
│   └── flyway.conf           # Configuração do Flyway
├── V1__create_products_tables.sql   # Migration versioned
├── U1__rollback_products_tables.sql # Undo migration
└── R__seed_data.sql          # Repeatable migration
```

### **Tipos de Migrations**

1. **Versioned Migrations (V)**: Executam uma única vez, em ordem
   - Formato: `V{version}__{description}.sql`
   - Exemplo: `V1__create_products_tables.sql`
   - Uso: Criar tabelas, alterar schemas, adicionar colunas

2. **Undo Migrations (U)**: Revertem migrations versionadas
   - Formato: `U{version}__{description}.sql`
   - Exemplo: `U1__rollback_products_tables.sql`
   - Uso: Rollback de mudanças específicas

3. **Repeatable Migrations (R)**: Executam toda vez que mudam
   - Formato: `R__{description}.sql`
   - Exemplo: `R__seed_data.sql`
   - Uso: Views, stored procedures, functions, seed data

## 🚀 Como Funciona

### **1. Primeira Execução**

```bash
# Flyway cria tabela de controle
CREATE TABLE flyway_schema_history (
    installed_rank INT NOT NULL,
    version VARCHAR(50),
    description VARCHAR(200) NOT NULL,
    type VARCHAR(20) NOT NULL,
    script VARCHAR(1000) NOT NULL,
    checksum INT,
    installed_by VARCHAR(100) NOT NULL,
    installed_on TIMESTAMP NOT NULL,
    execution_time INT NOT NULL,
    success BOOLEAN NOT NULL
);
```

### **2. Execução de Migrations**

```
┌─────────────────────────────────────────┐
│ 1. Flyway lê arquivos de migration     │
└─────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────┐
│ 2. Verifica flyway_schema_history      │
└─────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────┐
│ 3. Identifica migrations pendentes     │
└─────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────┐
│ 4. Executa migrations em ordem         │
└─────────────────────────────────────────┘
            ↓
┌─────────────────────────────────────────┐
│ 5. Registra resultado no histórico     │
└─────────────────────────────────────────┘
```

### **3. Validação**

O Flyway valida:
- ✅ Ordem das migrations
- ✅ Checksums (detecta modificações)
- ✅ Migrations pendentes
- ✅ Migrations aplicadas fora de ordem

## 🛠️ Comandos Principais

### **Migrate**

Aplica todas as migrations pendentes:

```bash
# Via Makefile
make db-migrate

# Via Docker Compose direto
docker-compose up flyway

# Resultado esperado
Flyway Community Edition 10.0.0
Database: jdbc:postgresql://postgres:5432/alderaan_db (PostgreSQL 16.0)
Successfully validated 2 migrations (execution time 00:00.012s)
Creating Schema History table "public"."flyway_schema_history" ...
Current version of schema "public": << Empty Schema >>
Migrating schema "public" to version "1 - create products tables"
Successfully applied 1 migration to schema "public", now at version v1 (execution time 00:00.087s)
```

### **Info**

Mostra status de todas as migrations:

```bash
# Via Makefile
make db-migrate-info

# Resultado esperado
+-----------+---------+-------------------------+------+---------------------+---------+
| Category  | Version | Description             | Type | Installed On        | State   |
+-----------+---------+-------------------------+------+---------------------+---------+
| Versioned | 1       | create products tables  | SQL  | 2024-10-04 22:00:00 | Success |
| Repeatable|         | seed data              | SQL  | 2024-10-04 22:00:05 | Success |
+-----------+---------+-------------------------+------+---------------------+---------+
```

### **Validate**

Valida migrations aplicadas vs arquivos:

```bash
# Via Makefile
make db-migrate-validate

# Valida:
# - Checksums não mudaram
# - Ordem está correta
# - Nenhuma migration faltando
```

### **Repair**

Repara histórico de migrations (útil após falhas):

```bash
# Via Makefile
make db-migrate-repair

# Remove entradas falhadas
# Recalcula checksums
```

### **Baseline**

Cria um baseline (ponto de partida) para banco existente:

```bash
# Via Makefile
make db-migrate-baseline

# Marca versão atual como baseline
# Migrations anteriores são ignoradas
```

### **Clean**

**⚠️ CUIDADO!** Remove todos os objetos do schema:

```bash
# Via Makefile
make db-clean-all

# Apaga:
# - Todas as tabelas
# - Views, procedures, functions
# - Histórico do Flyway
```

## 📝 Criando Migrations

### **1. Versioned Migration**

```bash
# Criar arquivo
cd db/migrations
touch V2__add_product_images.sql
```

```sql
-- V2__add_product_images.sql
-- Adiciona suporte a imagens de produtos

ALTER TABLE products
ADD COLUMN image_url VARCHAR(500);

CREATE TABLE product_images (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_product_images_product_id ON product_images(product_id);

COMMENT ON TABLE product_images IS 'Múltiplas imagens para produtos';
```

### **2. Undo Migration**

```bash
# Criar undo correspondente
touch U2__rollback_product_images.sql
```

```sql
-- U2__rollback_product_images.sql
-- Reverte adição de suporte a imagens

DROP TABLE IF EXISTS product_images;
ALTER TABLE products DROP COLUMN IF EXISTS image_url;
```

### **3. Repeatable Migration**

```bash
# Criar migration repetível
touch R__update_product_stats_view.sql
```

```sql
-- R__update_product_stats_view.sql
-- View de estatísticas de produtos (atualiza sempre que mudar)

CREATE OR REPLACE VIEW product_stats AS
SELECT
    COUNT(*) as total_products,
    SUM(price) as total_value,
    AVG(price) as avg_price,
    MIN(price) as min_price,
    MAX(price) as max_price
FROM products;

COMMENT ON VIEW product_stats IS 'Estatísticas agregadas de produtos';
```

## 🎯 Boas Práticas

### **1. Nomenclatura**

✅ **BOM:**
```
V1__create_products_tables.sql
V2__add_product_categories.sql
V3__add_user_authentication.sql
```

❌ **RUIM:**
```
001.sql
migration.sql
fix_stuff.sql
```

### **2. Uma Mudança por Migration**

✅ **BOM:**
```
V2__add_email_to_users.sql
V3__create_orders_table.sql
```

❌ **RUIM:**
```
V2__add_email_and_create_orders_and_fix_products.sql
```

### **3. Nunca Modificar Migrations Aplicadas**

```bash
# ❌ NUNCA faça isso se a migration já foi aplicada
vim V1__create_products_tables.sql  # Editar migration existente

# ✅ Crie uma nova migration
touch V4__modify_products_table.sql
```

### **4. Testar Migrations Localmente**

```bash
# 1. Limpar banco local
make db-clean-all

# 2. Aplicar migrations
make db-migrate

# 3. Validar
make db-migrate-validate

# 4. Testar rollback (se aplicável)
docker-compose run --rm flyway undo
```

### **5. Usar Transações**

```sql
-- Migration será toda revertida se algo falhar
BEGIN;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
);

ALTER TABLE products
ADD COLUMN user_id INTEGER REFERENCES users(id);

COMMIT;
```

### **6. Adicionar Rollback Safety**

```sql
-- Verificar se coluna já existe antes de adicionar
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name='products' AND column_name='image_url'
    ) THEN
        ALTER TABLE products ADD COLUMN image_url VARCHAR(500);
    END IF;
END $$;
```

## 🔄 Workflow de Desenvolvimento

### **Desenvolvimento Local**

```bash
# 1. Criar nova migration
vim db/migrations/V2__add_feature.sql

# 2. Aplicar localmente
make db-migrate

# 3. Verificar
make db-migrate-info

# 4. Testar API
make run

# 5. Commit
git add db/migrations/V2__add_feature.sql
git commit -m "feat: adiciona feature X"
```

### **CI/CD Pipeline**

```yaml
# .github/workflows/deploy.yml
steps:
  - name: Run Flyway Migrations
    run: |
      docker-compose up -d postgres
      docker-compose up flyway

  - name: Validate Migrations
    run: docker-compose run --rm flyway validate

  - name: Deploy Application
    run: docker-compose up -d api
```

### **Produção**

```bash
# 1. Backup antes de migrations
pg_dump alderaan_db > backup_$(date +%Y%m%d).sql

# 2. Aplicar migrations
docker-compose up flyway

# 3. Verificar status
make db-migrate-info

# 4. Se algo der errado, restore backup
psql alderaan_db < backup_20241004.sql
```

## 📊 Monitoramento

### **Verificar Histórico**

```sql
-- Ver todas as migrations aplicadas
SELECT
    installed_rank,
    version,
    description,
    type,
    installed_on,
    execution_time,
    success
FROM flyway_schema_history
ORDER BY installed_rank;
```

### **Detectar Migrations Pendentes**

```bash
# Listar migrations não aplicadas
make db-migrate-info | grep Pending
```

### **Verificar Integridade**

```bash
# Validar checksums
make db-migrate-validate

# Se checksum falhou:
# 1. Verificar se alguém modificou migration aplicada
# 2. Se foi modificação legítima, usar repair:
make db-migrate-repair
```

## 🐛 Troubleshooting

### **Erro: "Detected applied migration not resolved locally"**

```bash
# Alguém aplicou migration que não existe mais nos arquivos

# Solução 1: Recuperar arquivo de migration
git checkout origin/main -- db/migrations/V2__missing.sql

# Solução 2: Remover entrada do histórico (CUIDADO!)
make db-migrate-repair
```

### **Erro: "Validate failed: Migration checksum mismatch"**

```bash
# Migration foi modificada após ser aplicada

# Ver qual migration:
make db-migrate-info

# Solução 1: Reverter modificação
git checkout HEAD -- db/migrations/V2__modified.sql

# Solução 2: Aceitar mudança (atualizar checksum)
make db-migrate-repair
```

### **Erro: "Found non-empty schema(s) without schema history table"**

```bash
# Banco já tem objetos mas sem histórico Flyway

# Criar baseline
make db-migrate-baseline

# Agora pode aplicar novas migrations
make db-migrate
```

### **Erro: Migration falhou no meio**

```bash
# Ver logs
docker-compose logs flyway

# Ver último estado
make db-migrate-info

# Reparar histórico
make db-migrate-repair

# Corrigir migration
vim db/migrations/V3__failed.sql

# Tentar novamente
make db-migrate
```

## 🔐 Segurança

### **Não Commitar Senhas**

```bash
# ❌ Nunca faça isso
flyway.password=super_secret_password

# ✅ Use variáveis de ambiente
flyway.password=${DB_PASSWORD}
```

### **Backup Antes de Migrations**

```bash
# Script de backup automático
#!/bin/bash
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
docker-compose exec postgres pg_dump -U alderaan alderaan_db \
  > backups/backup_${TIMESTAMP}.sql

# Aplicar migrations
docker-compose up flyway
```

### **Testar em Staging Primeiro**

```bash
# 1. Deploy em staging
ENVIRONMENT=staging make db-migrate

# 2. Testes automatizados
make test

# 3. Se OK, deploy em produção
ENVIRONMENT=production make db-migrate
```

## 📚 Recursos

- **Documentação Oficial**: https://flywaydb.org/documentation/
- **Conceitos**: https://flywaydb.org/documentation/concepts/migrations
- **Comandos**: https://flywaydb.org/documentation/commands/
- **Configuração**: https://flywaydb.org/documentation/configuration/parameters

---

**Anterior:** [Docker & Deployment](08-docker-deployment.md) | **Próximo:** [Database](../db/README.md) | **Voltar ao início:** [README](README.md)
