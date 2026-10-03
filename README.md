# Servidor HTTP com Domain Driven Design, Gin e Shutdown Controlado

Este projeto demonstra como criar um servidor HTTP completo com Domain-Driven Design (DDD), usando Gin e implementando um graceful shutdown.

## 📋 Estrutura do Projeto

```
.
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── domain
│   │   └── product
│   │       ├── entity
│   │       │   └── product.go
│   │       ├── events
│   │       │   └── product_created_event.go
│   │       └── repository
│   │           └── product_repository.go
│   ├── infra
│   │   └── http
│   │       ├── handlers
│   │       │   └── product_handler.go
│   │       └── router
│   │           └── product_router.go
│   └── shared
│       └── domain
│           └── events
│               └── events_handler.go
└── README.md
```

## 🎯 Características

- **Domínio puro e desacoplado**: Entidades de domínio com validação integrada
- **Repositório thread-safe**: Usando `sync.RWMutex` para acesso concorrente seguro
- **Dispatcher de eventos de domínio**: Sistema de eventos para extensibilidade
- **Graceful shutdown**: Encerramento controlado respeitando requisições em andamento
- **Arquitetura limpa**: Separação clara entre camadas de domínio, aplicação e infraestrutura
- **Flyway Migrations**: Versionamento automático e rastreável do schema do banco

## 🚀 Como Executar

### Pré-requisitos

- Go 1.23 ou superior

### Instalação

1. Clone o repositório
2. Instale as dependências:

```bash
go mod download
```

3. **Configurar banco de dados:**

```bash
# Copiar arquivo de configuração
cp config.env.example config.env

# Editar se necessário (valores padrão já funcionam)
# vim config.env

# Iniciar PostgreSQL
make db-up

# Aguardar ~10 segundos para migrations rodarem automaticamente
```

4. **Execute a API:**

**Opção 1:** Com Makefile (recomendado):

```bash
make run
```

**Opção 2:** Diretamente:

```bash
go run cmd/main.go
```

O servidor estará disponível em `http://localhost:8080`

### Comandos Úteis

```bash
make help           # Ver todos os comandos disponíveis
make swagger        # Gerar documentação Swagger
make build          # Compilar binário
make run            # Executar com auto-geração de docs
make dev            # Executar sem regenerar docs
make clean          # Limpar binários
make test           # Executar testes

# Banco de dados
make db-up          # Iniciar PostgreSQL
make db-down        # Parar PostgreSQL
make db-connect     # Conectar ao banco
make db-seed        # Popular com dados
make db-clean       # Limpar dados
make db-reset       # Resetar banco

# Plataforma completa
make platform-up    # Inicia tudo (DB + API + Monitoring)
make platform-down  # Para tudo
make platform-logs  # Ver logs de tudo
make platform-status # Status dos serviços
```

### 🐳 Execução via Docker (GitHub Container Registry)

A imagem multi-arquitetura (`linux/amd64`, `linux/arm64`) está disponível no GitHub Container Registry:

```bash
# Baixar a versão mais recente
docker pull ghcr.io/tech-preta/alderaan-api:latest

# Ou uma versão específica
docker pull ghcr.io/tech-preta/alderaan-api:1.0.0

# Executar container apontando para seu PostgreSQL
docker run -d \
  --name alderaan-api \
  -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=alderaan \
  -e DB_PASSWORD=alderaan123 \
  -e DB_NAME=alderaan_db \
  ghcr.io/tech-preta/alderaan-api:latest
```

### 📊 Instalação via Helm Chart (OCI Registry)

O Helm Chart oficial é publicado no GitHub Packages (OCI Registry):

```bash
# Instalação direta com PostgreSQL embutido
helm install alderaan oci://ghcr.io/tech-preta/helm-charts/alderaan

# Instalação de versão específica com namespace dedicado
helm install alderaan oci://ghcr.io/tech-preta/helm-charts/alderaan \
  --version 1.0.0 \
  --namespace alderaan \
  --create-namespace
```

📖 [**Guia completo do Helm Chart →**](charts/README.md)

## 📚 Documentação Interativa (Swagger)

Acesse a documentação interativa da API em:

```
http://localhost:8080/swagger/index.html
```

Você pode testar todos os endpoints diretamente pelo navegador! 🎯

## 📡 Endpoints da API

### Criar Produto

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notebook",
    "sku": 12345,
    "categories": ["Eletrônicos", "Computadores"],
    "price": 3500
  }'
```

### Listar Todos os Produtos

```bash
curl http://localhost:8080/api/v1/products
```

### Buscar Produto por Nome

```bash
curl http://localhost:8080/api/v1/products/Notebook
```

## 🏗️ Arquitetura

### Camada de Domínio

- **Entidade Product**: Representa um produto com validações de negócio
- **ProductCreatedEvent**: Evento disparado quando um produto é criado
- **ProductRepository**: Interface e implementação para persistência de produtos

### Camada de Infraestrutura

- **ProductHandler**: Handlers HTTP para operações com produtos
- **Router**: Configuração de rotas da API

### Camada Compartilhada

- **EventDispatcher**: Sistema de eventos para comunicação entre componentes

## 🛡️ Graceful Shutdown

O servidor implementa um shutdown controlado que:

1. Captura sinais de terminação (SIGINT, SIGTERM)
2. Aguarda até 5 segundos para finalizar requisições em andamento
3. Encerra o servidor de forma limpa

Para testar, execute o servidor e pressione `Ctrl+C`. Você verá:

```
Received termination signal. Shutting down server...
Server shut down gracefully.
```

## 🎯 Características implementadas

- ✅ **Arquitetura Limpa** - Separação clara entre camadas
- ✅ **DDD** - Domain-Driven Design com entidades ricas
- ✅ **Persistência PostgreSQL** - Dados persistidos em banco relacional
- ✅ **Migrations** - Controle de versão do schema do banco
- ✅ **Thread-Safe** - Repositório com proteção para concorrência
- ✅ **Event Dispatcher** - Sistema de eventos de domínio
- ✅ **Graceful Shutdown** - Encerramento controlado com fechamento do banco
- ✅ **RESTful API** - Endpoints bem definidos com Gin
- ✅ **Swagger/OpenAPI** - Documentação interativa automática
- ✅ **Prometheus Monitoring** - Golden Signals + métricas de negócio

## 📊 Monitoramento (Prometheus + Grafana)

### Métricas Disponíveis

**Golden Signals (Google SRE):**
- 🕐 **Latency**: Tempo de resposta das requisições
- 📈 **Traffic**: Taxa de requisições por segundo
- ❌ **Errors**: Taxa de erros (4xx, 5xx)
- 📊 **Saturation**: Requisições simultâneas

**Métricas de Negócio:**
- Produtos criados (total e taxa)
- Total de produtos em estoque
- Produtos por categoria
- Valor total do inventário
- Preço médio dos produtos

### Acessar Métricas

```
http://localhost:8080/metrics  # Endpoint Prometheus
```

### Iniciar Stack Completa (Tudo de Uma Vez)

```bash
# Inicia PostgreSQL + API + Prometheus + Grafana
docker-compose up -d

# Ou use o Makefile
make platform-up
```

Acesse:
- **API**: `http://localhost:8080`
- **Swagger**: `http://localhost:8080/swagger/index.html`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)
- **PostgreSQL**: `localhost:5432`

📖 [**Guia completo de monitoramento →**](monitoring/README.md)

## 📚 Tecnologias Utilizadas

- **Go 1.24**: Linguagem de programação
- **Gin**: Framework web de alta performance
- **PostgreSQL**: Banco de dados relacional
- **Prometheus**: Sistema de monitoramento e métricas
- **Swagger/OpenAPI**: Documentação interativa da API
- **Flyway**: Migrations de banco de dados versionadas
- **DDD**: Domain-Driven Design para arquitetura limpa
- **Docker**: Containerização com multi-stage build (~110MB)

## 📚 Documentação Detalhada

Para entender os conceitos e padrões utilizados neste projeto, consulte a documentação técnica:

- **[Domain-Driven Design (DDD)](docs/01-domain-driven-design.md)** - Entenda como organizamos o domínio
- **[Arquitetura Limpa](docs/02-clean-architecture.md)** - Conheça a estrutura de camadas
- **[Event Dispatcher](docs/03-event-dispatcher.md)** - Aprenda sobre eventos de domínio
- **[Graceful Shutdown](docs/04-graceful-shutdown.md)** - Veja como implementamos shutdown controlado
- **[RESTful API com Gin](docs/05-restful-api-gin.md)** - Domine o framework Gin
- **[Swagger / OpenAPI](docs/06-swagger-documentation.md)** - Documentação interativa da API
- **[Prometheus Monitoring](docs/07-prometheus-monitoring.md)** - Golden Signals e métricas de negócio
- **[Docker & Deployment](docs/08-docker-deployment.md)** - Multi-stage build e containerização
- **[Flyway Migrations](docs/09-flyway-migrations.md)** - Gerenciamento profissional de migrations
- **[Prometheus Queries (PromQL)](docs/10-prometheus-queries.md)** - Guia completo de queries para métricas
- **[Releases Automáticos](docs/12-automated-releases.md)** - Sistema de releases, tags e CHANGELOG automáticos
- **[Database PostgreSQL](db/README.md)** - Schema SQL, migrations e persistência

📖 [**Ver toda a documentação →**](docs/README.md)

## 🚀 Releases e Versionamento

Este projeto utiliza **Semantic Versioning** e **Conventional Commits** para releases automáticos:

- ✅ **Versioning automático** baseado em commits
- ✅ **CHANGELOG.md** gerado automaticamente
- ✅ **Tags Git** criadas automaticamente
- ✅ **GitHub Releases** publicados automaticamente

### Como funciona

Ao fazer push para `main` com commits no formato [Conventional Commits](https://www.conventionalcommits.org/pt-br/):

```bash
# Nova funcionalidade (minor version)
git commit -m "feat: adiciona autenticação JWT"

# Correção de bug (patch version)
git commit -m "fix: corrige timeout em requisições"

# Breaking change (major version)
git commit -m "feat!: redesenha API v2"
```

O sistema automaticamente:
1. Analisa os commits
2. Determina a nova versão
3. Atualiza o CHANGELOG.md
4. Cria tag Git
5. Publica release no GitHub

📖 [**Guia completo de releases →**](docs/12-automated-releases.md)

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

**Baseado no artigo de [William Koller](https://williamkoller.substack.com)**
