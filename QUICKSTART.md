# 🚀 Quickstart - Alderaan API

Guia rápido para começar a usar a API com PostgreSQL.

## ⚡ Início Rápido (2 minutos)

### **Opção 1: Stack Completa com Docker (Mais Fácil) 🚀**

```bash
# Clone o repositório
git clone <url-do-repo>
cd alderaan

# Inicia TUDO de uma vez (PostgreSQL + API + Prometheus + Grafana)
docker-compose up -d

# Ou use o Makefile
make platform-up

# Aguardar ~30 segundos para tudo iniciar
```

Você terá acesso a:
- **API**: `http://localhost:8080`
- **Swagger**: `http://localhost:8080/swagger/index.html`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)
- **PostgreSQL**: `localhost:5432`

### **Opção 2: Desenvolvimento Local**

```bash
# Copie o arquivo de configuração
cp config.env.example config.env

# Instale dependências
go mod download

# Iniciar PostgreSQL
make db-up

# Aguardar ~10 segundos (migrations automáticas)

# Executar a API
make run
```

Você verá:
```
✅ Conectado ao banco de dados PostgreSQL
📊 Usando repositório PostgreSQL

🚀 Server running on http://localhost:8080
📊 Metrics available at http://localhost:8080/metrics
📚 Swagger UI at http://localhost:8080/swagger/index.html
🗄️  Database: alderaan@localhost:5432/alderaan_db
```

### **4. Testar a API**

#### **Criar um produto:**

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "sku": 99999,
    "categories": ["Eletrônicos", "Computadores"],
    "price": 15000
  }'
```

#### **Listar produtos:**

```bash
curl http://localhost:8080/api/v1/products
```

#### **Buscar produto específico:**

```bash
curl http://localhost:8080/api/v1/products/MacBook%20Pro
```

### **5. Verificar no Banco**

```bash
# Conectar ao PostgreSQL
make db-connect

# No psql, execute:
SELECT * FROM products;
SELECT * FROM categories;
SELECT p.name, c.name as category
FROM products p
JOIN product_categories pc ON p.id = pc.product_id
JOIN categories c ON pc.category_id = c.id;

# Sair: \q
```

## 🎯 Dados de Exemplo

O banco já vem com dados de seed. Para popular:

```bash
# Popular com 8 produtos de exemplo
make db-seed
```

Produtos incluídos:
- Notebook Dell Inspiron
- Mouse Gamer RGB
- Teclado Mecânico
- Monitor 4K LG
- Webcam HD Logitech
- Headset HyperX
- SSD Samsung 1TB
- Memória RAM 16GB

## 📊 Interfaces Web Disponíveis

### **Swagger UI - Documentação da API**
```
http://localhost:8080/swagger/index.html
```
- ✅ Ver todos os endpoints
- ✅ Testar diretamente pelo navegador
- ✅ Ver exemplos de request/response
- ✅ Validar schemas

### **Prometheus - Métricas e Alertas**
```
http://localhost:9090
```
- ✅ Consultar métricas com PromQL
- ✅ Ver alertas ativos
- ✅ Explorar targets e status

### **Grafana - Dashboards Visuais**
```
http://localhost:3000
```
- 👤 Login: admin/admin
- ✅ Datasource Prometheus já configurado
- ✅ Criar dashboards personalizados
- ✅ Visualizar métricas em tempo real

## 🧹 Limpeza

```bash
# Parar PostgreSQL
make db-down

# Parar monitoramento
make monitoring-down

# Limpar binários
make clean

# Resetar banco (apaga tudo e recria)
make db-reset
```

## 🔧 Comandos Úteis

```bash
# Desenvolvimento
make dev            # Rodar sem regenerar Swagger
make swagger        # Apenas gerar Swagger
make build          # Compilar binário

# Banco de dados
make db-logs        # Ver logs do PostgreSQL
make db-connect     # Conectar no banco
make db-clean       # Limpar dados
make db-migrate     # Rodar migrations

# Monitoramento
make monitoring-logs  # Logs do Prometheus/Grafana
make check-metrics    # Verificar endpoint /metrics

# Ajuda
make help           # Ver todos os comandos
```

## 🐳 Docker Compose Unificado

Um único `docker-compose.yml` na raiz gerencia toda a plataforma:

```bash
# Iniciar toda a stack
docker-compose up -d

# Iniciar serviços específicos
docker-compose up -d postgres        # Apenas banco
docker-compose up -d postgres api    # Banco + API
docker-compose up -d prometheus grafana  # Apenas monitoramento

# Ver status
docker-compose ps

# Ver logs
docker-compose logs -f

# Parar tudo
docker-compose down

# Ou use Makefile
make platform-up        # Iniciar tudo
make platform-down      # Parar tudo
make platform-logs      # Ver logs
make platform-status    # Ver status
```

Serviços incluídos:
- ✅ **PostgreSQL** (localhost:5432) - Banco de dados
- ✅ **API** (localhost:8080) - Aplicação Go
- ✅ **Prometheus** (localhost:9090) - Métricas
- ✅ **Grafana** (localhost:3000) - Visualização

## 🔍 Troubleshooting

### **Erro: "connection refused"**

```bash
# Verificar se PostgreSQL está rodando
docker-compose ps postgres

# Ver logs
docker-compose logs postgres

# Reiniciar
make db-down
make db-up
```

### **Erro: "role does not exist"**

```bash
# Resetar banco completamente
make db-down
docker-compose down -v  # Remove volumes
make db-up
```

### **Erro: "port 5432 already in use"**

```bash
# Parar PostgreSQL local
brew services stop postgresql

# Ou mudar porta no config.env
DB_PORT=5433
```

### **Erro: "relation does not exist"**

```bash
# Rodar migrations manualmente
make db-migrate
```

## 📚 Próximos Passos

1. ✅ Leia a [documentação completa](docs/README.md)
2. ✅ Explore o [schema do banco](db/README.md)
3. ✅ Configure [monitoramento](monitoring/README.md)
4. ✅ Estude os [padrões de arquitetura](docs/02-clean-architecture.md)

## 🎓 Estrutura do Projeto

```
.
├── cmd/main.go              # Entrada da aplicação
├── internal/
│   ├── domain/              # Lógica de negócio
│   ├── infra/               # Infraestrutura (HTTP, persistência)
│   └── shared/              # Código compartilhado
├── db/
│   ├── migrations/          # Migrations SQL
│   └── seed.sql            # Dados iniciais
├── monitoring/              # Prometheus + Grafana
└── docs/                   # Documentação

```

## 💡 Dicas

- 🔥 Use `make help` para ver todos os comandos disponíveis
- 🎯 Teste a API pelo Swagger UI (mais fácil que curl)
- 📊 Configure dashboards no Grafana para visualizar métricas
- 🗄️ Use `make db-connect` para explorar o banco diretamente
- 🧪 Use `make db-reset` para recomeçar com dados limpos

---

**🚀 Pronto! Agora você tem uma API completa com DDD, PostgreSQL e monitoramento!**
