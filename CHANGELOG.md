# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [1.0.0] - 2024-10-04

### 🎉 Lançamento Inicial

Primeira versão completa da API Alderaan com Domain-Driven Design, Clean Architecture, e observabilidade completa.

### ✨ Adicionado

#### **Arquitetura e Design**
- Implementação completa de Domain-Driven Design (DDD)
- Clean Architecture com separação clara de camadas (Domain, Infrastructure, Shared)
- Event Dispatcher para comunicação desacoplada via eventos de domínio
- Repository Pattern com thread-safety usando `sync.RWMutex`
- Graceful Shutdown para encerramento controlado do servidor

#### **API RESTful**
- Framework Gin para alta performance HTTP
- Endpoints RESTful para gerenciamento de produtos:
  - `POST /api/v1/products` - Criar produto
  - `GET /api/v1/products` - Listar todos os produtos
  - `GET /api/v1/products/:name` - Buscar produto por nome
- Health check endpoint `/health`
- Validação de entrada de dados
- Tratamento de erros padronizado

#### **Documentação Interativa**
- Swagger/OpenAPI integrado via `swaggo/swag`
- Documentação automática da API
- Interface Swagger UI em `/swagger/index.html`
- Anotações completas nos handlers HTTP
- Exemplos de request/response

#### **Observabilidade e Monitoramento**
- **Prometheus** para coleta de métricas
  - Golden Signals (Latency, Traffic, Errors, Saturation)
  - Métricas de negócio (produtos criados, total, por categoria, valor total, preço médio)
  - Endpoint `/metrics` para scraping
  - Middleware automático para coleta de métricas HTTP
- **Grafana** para visualização de dados
  - Dashboard "Alderaan API Overview" pré-configurado
  - 13 painéis cobrindo Golden Signals e métricas de negócio
  - Datasource Prometheus pré-provisionado
- **Alertmanager** para gerenciamento de alertas
  - 15+ regras de alertas configuradas
  - Rotas para alertas críticos, warnings e informativos
  - Inibição de alertas duplicados
  - Webhook receivers configuráveis

#### **Persistência de Dados**
- **PostgreSQL** como banco de dados principal
- Schema com 3 tabelas:
  - `products` - Dados dos produtos
  - `categories` - Categorias disponíveis
  - `product_categories` - Relacionamento N:N
- Repositório PostgreSQL implementando `IProductRepository`
- Pool de conexões configurável
- Triggers automáticos para `updated_at`

#### **Migrações de Banco de Dados**
- **Flyway** para versionamento do schema
- Migrações versionadas (`V1__create_products_tables.sql`)
- Migrações de undo (`U1__rollback_products_tables.sql`)
- Migrações repetíveis para seed data (`R__seed_data.sql`)
- Comandos Makefile para gerenciamento de migrações

#### **Containerização**
- **Dockerfile multi-stage** para build otimizado
  - Stage 1: Builder (compilação)
  - Stage 2: Runtime (Alpine Linux)
  - Usuário não-root para segurança
  - Health check integrado
  - Imagem final < 30MB
- **Docker Compose** unificado para toda a plataforma
  - PostgreSQL
  - API Go
  - Prometheus
  - Grafana
  - Alertmanager
  - Flyway (migrações)
  - Network compartilhada
  - Volumes persistentes

#### **Testes**
- **Suíte completa de testes unitários**
  - 7 arquivos de teste criados
  - 1.790+ linhas de código de teste
  - 52 testes unitários
  - 12 benchmarks
  - 6 pacotes com 100% de cobertura
- **Cobertura de Testes:**
  - `domain/product/entity`: 100%
  - `domain/product/events`: 100%
  - `domain/product/repository`: 100%
  - `metrics`: 100%
  - `shared/config`: 100%
  - `shared/domain/events`: 100%
- **Técnicas de Teste:**
  - Table-driven tests
  - Subtests para organização hierárquica
  - Testes de concorrência e thread-safety
  - Race condition detection
  - Benchmarks para tracking de performance
  - Mocks simples para eventos

#### **Automação e Tooling**
- **Makefile** com 30+ comandos úteis:
  - Build e execução
  - Testes (test, test-coverage, test-coverage-html, test-unit)
  - Geração de Swagger
  - Gerenciamento Docker Compose (platform-up, platform-down, platform-logs)
  - Gerenciamento PostgreSQL (db-up, db-down, db-migrate, db-seed, db-reset)
  - Monitoramento (monitoring-up, monitoring-down, check-metrics)
  - Migrações Flyway (db-migrate-info, db-migrate-validate)
- **Scripts de inicialização** simplificados

#### **Configuração**
- Gerenciamento via variáveis de ambiente
- Arquivo `config.env.example` como template
- Valores padrão sensatos
- Suporte a diferentes ambientes (dev, prod)
- Configuração centralizada em `internal/shared/config`

#### **Documentação Completa**
- **12 documentos técnicos** em `docs/`:
  1. Domain-Driven Design (DDD)
  2. Clean Architecture
  3. Event Dispatcher
  4. Graceful Shutdown
  5. RESTful API com Gin
  6. Swagger Documentation
  7. Prometheus Monitoring
  8. Docker Deployment
  9. Flyway Migrations
  10. Prometheus Queries (PromQL)
  11. Refatoração: Desacoplamento do Pacote Metrics
  12. Estratégia de Testes
- **README.md** principal com quickstart
- **QUICKSTART.md** detalhado
- **TESTING-SUMMARY.md** com resumo de testes
- **REFACTORING-SUMMARY.md** com histórico de refatorações

### 🏗️ Estrutura do Projeto

```
alderaan/
├── cmd/
│   └── main.go                          # Entry point da aplicação
├── internal/
│   ├── domain/                          # Camada de Domínio (DDD)
│   │   └── product/
│   │       ├── entity/                  # Entidades
│   │       ├── events/                  # Eventos de domínio
│   │       └── repository/              # Interfaces de repositório
│   ├── infra/                           # Camada de Infraestrutura
│   │   ├── http/                        # HTTP handlers e routers
│   │   └── persistence/                 # Implementação PostgreSQL
│   ├── metrics/                         # Observabilidade (Prometheus)
│   └── shared/                          # Componentes compartilhados
│       ├── config/                      # Configuração
│       ├── database/                    # Cliente de banco de dados
│       └── domain/events/               # Event dispatcher
├── db/
│   └── migrations/                      # Migrações Flyway
├── monitoring/
│   ├── prometheus.yml                   # Config Prometheus
│   ├── alerts.yml                       # Regras de alertas
│   ├── alertmanager.yml                 # Config Alertmanager
│   └── grafana/                         # Dashboards e provisioning
├── docs/                                # Documentação técnica (12 docs)
├── docker-compose.yml                   # Stack completa
├── Dockerfile                           # Multi-stage build
├── Makefile                             # Automação (30+ comandos)
└── go.mod                               # Dependências Go
```

### 🔧 Tecnologias Utilizadas

| Categoria | Tecnologia | Versão |
|-----------|-----------|---------|
| **Linguagem** | Go | 1.24+ |
| **Web Framework** | Gin | v1.10.0 |
| **Banco de Dados** | PostgreSQL | 15-alpine |
| **Migrações** | Flyway | 10-alpine |
| **Monitoramento** | Prometheus | v2.54.1 |
| **Visualização** | Grafana | 11.2.0 |
| **Alertas** | Alertmanager | v0.27.0 |
| **Documentação** | Swagger/OpenAPI | swaggo/swag |
| **Métricas** | Prometheus Client | v1.23.2 |
| **Containerização** | Docker & Docker Compose | v3.8 |

### 📈 Métricas de Qualidade

- **Cobertura de Testes**: 38.6% geral, 100% na camada de domínio
- **Linhas de Código de Teste**: 1.790+
- **Zero Race Conditions**: Validado com `go test -race`
- **Performance**: Testes executam em < 2s
- **Documentação**: 12 documentos técnicos (5.000+ linhas)

### 🚀 Comandos Principais

```bash
# Iniciar toda a plataforma
make platform-up

# Rodar testes
make test-coverage

# Gerar documentação Swagger
make swagger

# Build da aplicação
make build

# Migrações de banco
make db-migrate

# Ver logs
make platform-logs
```

### 🔒 Segurança

- Dockerfile executa como usuário não-root
- Variáveis sensíveis via environment variables
- SSL/TLS configurável para PostgreSQL
- Validação de input em todos os endpoints
- Health checks em todos os serviços

### 📊 Observabilidade

#### **Golden Signals Implementados**
- ✅ **Latency**: Histograma de duração de requisições
- ✅ **Traffic**: Contador total de requisições
- ✅ **Errors**: Contador de erros (4xx, 5xx)
- ✅ **Saturation**: Gauge de requisições em andamento

#### **Métricas de Negócio**
- ✅ Total de produtos criados
- ✅ Total de produtos ativos
- ✅ Produtos por categoria
- ✅ Valor total do inventário
- ✅ Preço médio dos produtos

### 🎯 Próximas Versões (Roadmap)

#### **v1.1.0**
- [ ] Autenticação e autorização (JWT)
- [ ] Paginação nos endpoints de listagem
- [ ] Filtros e ordenação
- [ ] Soft delete de produtos

#### **v1.2.0**
- [ ] Testes de integração HTTP
- [ ] Testes de integração PostgreSQL
- [ ] CI/CD com GitHub Actions
- [ ] Testes E2E

#### **v1.3.0**
- [ ] Tracing distribuído (Jaeger/Tempo)
- [ ] Structured logging
- [ ] Rate limiting
- [ ] Caching (Redis)

### 🐛 Correções

Nenhuma correção necessária nesta versão inicial.

### ⚡ Melhorias de Performance

- Repository thread-safe com `sync.RWMutex`
- Pool de conexões PostgreSQL otimizado
- Dockerfile multi-stage reduz imagem para < 30MB
- Middleware de métricas com overhead mínimo

### 🔄 Refatorações

#### **Desacoplamento do Pacote Metrics**
- Movido de `internal/shared/metrics/` para `internal/metrics/`
- Melhora a separação de responsabilidades
- Facilita crescimento futuro (tracing, logging)
- 3 imports atualizados
- Documentação completa da refatoração

### 📝 Notas de Atualização

Para atualizar de uma versão anterior (se existisse):

1. **Banco de Dados**: Execute migrações Flyway
   ```bash
   make db-migrate
   ```

2. **Configuração**: Atualize `config.env` com novas variáveis

3. **Docker**: Recrie os containers
   ```bash
   make platform-down
   make platform-up
   ```

### 👥 Contribuidores

- **Natália Granato** - Implementação inicial
- **William Koller** - Arquitetura base (referência)

### 📄 Licença

Este projeto está licenciado sob a Licença MIT.

### 🔗 Links Úteis

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Alertmanager**: http://localhost:9093
- **Métricas**: http://localhost:8080/metrics
- **Health Check**: http://localhost:8080/health

---

## [Não Lançado]

### Em Desenvolvimento
- Autenticação JWT
- Testes de integração

---

**Formato**: [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/)
**Versionamento**: [Semantic Versioning](https://semver.org/lang/pt-BR/)

[1.0.0]: https://github.com/williamkoller/golang-domain-driven-design/releases/tag/v1.0.0
