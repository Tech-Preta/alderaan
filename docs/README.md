# 📚 Documentação Técnica

Bem-vindo à documentação técnica do projeto! Aqui você encontrará guias detalhados sobre os principais conceitos e padrões utilizados.

## 📖 Índice

### 1. [Domain-Driven Design (DDD)](01-domain-driven-design.md)
Aprenda sobre Domain-Driven Design e como aplicar conceitos como Entidades, Value Objects, Agregados, Repositórios e Eventos de Domínio.

**Tópicos abordados:**
- O que é DDD
- Conceitos fundamentais (Entidades, Value Objects, Agregados)
- Repositórios e Serviços de Domínio
- Eventos de Domínio
- Benefícios e quando usar
- Exemplo prático no projeto

---

### 2. [Arquitetura Limpa (Clean Architecture)](02-clean-architecture.md)
Entenda os princípios da Clean Architecture e como organizar seu código em camadas desacopladas e testáveis.

**Tópicos abordados:**
- Princípios fundamentais
- Regra de dependência
- Inversão de dependência
- Camadas (Entities, Use Cases, Interface Adapters, Frameworks)
- Fluxo de dados
- Benefícios e anti-padrões
- Princípios SOLID

---

### 3. [Event Dispatcher (Despachador de Eventos)](03-event-dispatcher.md)
Descubra como implementar comunicação desacoplada entre componentes usando eventos de domínio.

**Tópicos abordados:**
- O que é Event Dispatcher
- Problema que resolve
- Componentes (Event, Handler, Dispatcher)
- Execução assíncrona
- Thread-safety com sync.RWMutex
- Exemplos práticos
- Boas práticas

---

### 4. [Graceful Shutdown (Desligamento Controlado)](04-graceful-shutdown.md)
Aprenda a implementar shutdown controlado para garantir que sua aplicação encerre de forma segura.

**Tópicos abordados:**
- O que é Graceful Shutdown
- Captura de sinais (SIGINT, SIGTERM)
- Context com timeout
- Implementação completa
- Cenários avançados (múltiplos recursos, workers)
- Integração com Docker e Kubernetes
- Boas práticas

---

### 5. [RESTful API com Gin](05-restful-api-gin.md)
Domine o framework Gin e aprenda a criar APIs RESTful robustas e performáticas.

**Tópicos abordados:**
- O que é Gin e REST
- Métodos HTTP
- Roteamento e grupos de rotas
- Binding e validação de dados
- Middleware (autenticação, CORS, logging)
- Implementação completa de CRUD
- Boas práticas RESTful
- Performance

---

### 6. [Swagger / OpenAPI Documentation](06-swagger-documentation.md)
Aprenda a documentar sua API automaticamente usando Swagger/OpenAPI com swaggo.

**Tópicos abordados:**
- O que é Swagger e OpenAPI
- Configuração do swaggo/swag
- Anotações Swagger no código
- Gerando documentação automaticamente
- Interface Swagger UI interativa
- Testando endpoints pelo navegador
- Boas práticas de documentação
- Swagger em produção

---

### 7. [Prometheus Monitoring - Golden Signals](07-prometheus-monitoring.md)
Aprenda a instrumentar e monitorar sua aplicação com Prometheus baseado em Golden Signals e métricas de negócio.

**Tópicos abordados:**
- O que é Prometheus e Golden Signals
- Latency, Traffic, Errors, Saturation (Golden Signals)
- Métricas de negócio
- Tipos de métricas (Counter, Gauge, Histogram, Summary)
- Configuração do Prometheus
- Queries PromQL essenciais
- Dashboards no Grafana
- Alertas e boas práticas

---

### 8. [Docker & Deployment](08-docker-deployment.md)
Aprenda a containerizar e fazer deploy da aplicação usando Docker e Docker Compose.

**Tópicos abordados:**
- Dockerfile multi-stage
- Docker Compose
- Build e otimizações
- Deploy em produção
- CI/CD Pipeline
- Segurança

---

### 9. [Flyway Migrations](09-flyway-migrations.md)
Gerencie o schema do banco de dados com migrações versionadas usando Flyway.

**Tópicos abordados:**
- Versionamento de schema
- Migrações e rollbacks
- Seed data
- Comandos e boas práticas

---

### 10. [Prometheus Queries](10-prometheus-queries.md)
Consultas PromQL avançadas para análise de métricas.

**Tópicos abordados:**
- Queries PromQL
- Análise de latência
- Detecção de erros
- Performance tuning

---

### 11. [Refactoring Metrics](11-refactoring-metrics.md)
Refatoração e otimização de coleta de métricas.

**Tópicos abordados:**
- Organização de métricas
- Performance
- Boas práticas

---

### 12. [🚀 Releases, Tags e CHANGELOG Automáticos](12-automated-releases.md)
Aprenda como funciona o sistema de releases automáticos do projeto.

**Tópicos abordados:**
- Conventional Commits
- Semantic Versioning
- Release automático via GitHub Actions
- Geração automática de CHANGELOG
- Tags Git automáticas
- Comandos Make para releases
- Exemplos práticos
- Troubleshooting

---

## 🎯 Fluxo de Leitura Recomendado

Se você é novo no projeto, recomendamos ler nesta ordem:

1. **Iniciantes:**
   - [RESTful API com Gin](05-restful-api-gin.md) - Entenda a camada HTTP
   - [Swagger Documentation](06-swagger-documentation.md) - Documente sua API
   - [Graceful Shutdown](04-graceful-shutdown.md) - Aprenda sobre shutdown seguro
   - [Prometheus Queries](10-prometheus-queries.md) - Consulte métricas com PromQL

2. **Intermediário:**
   - [Domain-Driven Design](01-domain-driven-design.md) - Conceitos de domínio
   - [Event Dispatcher](03-event-dispatcher.md) - Comunicação desacoplada
   - [Flyway Migrations](09-flyway-migrations.md) - Gerenciamento de schema
   - [Prometheus Monitoring](07-prometheus-monitoring.md) - Monitoramento e observabilidade

3. **Avançado:**
   - [Arquitetura Limpa](02-clean-architecture.md) - Visão arquitetural completa
   - [Docker & Deployment](08-docker-deployment.md) - Deploy com containers

## 🔗 Links Úteis

- [README Principal](../README.md) - Voltar ao README do projeto
- [Código do Projeto](../internal/) - Explorar implementação
- [Artigo Original](https://williamkoller.substack.com) - Artigo que inspirou o projeto

## 💡 Como Contribuir com a Documentação

Encontrou um erro ou quer melhorar a documentação? Contribuições são bem-vindas!

1. Faça um fork do projeto
2. Edite os arquivos markdown na pasta `docs/`
3. Envie um pull request

## 📞 Dúvidas?

Se tiver dúvidas ou sugestões sobre a documentação, abra uma issue no projeto.

---

**Feito com ❤️ para a comunidade Go**
