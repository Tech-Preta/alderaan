# 🪐 Alderaan API

API RESTful em Go construída com **Domain-Driven Design (DDD)**, **Clean Architecture**, observabilidade completa (**Prometheus** e **Grafana**) e persistência relacional com **PostgreSQL**.

---

## 🎯 Principais Características

- **Clean Architecture & DDD**: Separação rigorosa entre domínio, aplicação e infraestrutura com entidades ricas e validação de regras de negócio.
- **Event Dispatcher**: Sistema desacoplado para disparo e tratamento de eventos de domínio.
- **Concorrência Segura & Graceful Shutdown**: Repositórios com controle thread-safe e encerramento controlado de conexões e requisições HTTP.
- **Banco de Dados & Migrations**: PostgreSQL com versionamento de schema automatizado via **Flyway**.
- **Observabilidade Completa**: Métricas nativas de Golden Signals (latência, tráfego, erros, saturação) e métricas de negócio expostas no endpoint `/metrics` com dashboards provisionados no **Grafana**.
- **Documentação Interativa**: Swagger / OpenAPI gerado automaticamente.
- **Cloud Native & Distribuição**:
  - Imagem Docker multi-arquitetura (`linux/amd64`, `linux/arm64`) no **GitHub Container Registry (GHCR)**.
  - **Helm Chart** empacotado e publicado como artefato OCI no GitHub Packages.
  - Pipeline de CI/CD com versionamento semântico automático (**Semantic Release**).

---

## 🚀 Como Usar

### Pré-requisitos

- [Go](https://golang.org/) 1.26+
- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/) (opcional, para automação de tarefas)

---

### Opção 1: Stack Completa com Docker Compose (Recomendado)

Inicia PostgreSQL, Flyway migrations, API, Prometheus e Grafana em uma única linha:

```bash
docker-compose up -d
# ou: make platform-up
```

**Painéis e serviços disponíveis:**

| Serviço | URL | Credenciais |
| :--- | :--- | :--- |
| **API REST** | `http://localhost:8080/api/v1/products` | — |
| **Swagger UI** | `http://localhost:8080/swagger/index.html` | — |
| **Métricas** | `http://localhost:8080/metrics` | — |
| **Prometheus** | `http://localhost:9090` | — |
| **Grafana** | `http://localhost:3000` | `admin` / `admin` |
| **PostgreSQL** | `localhost:5432` | `alderaan` / `alderaan123` |

Para encerrar:
```bash
docker-compose down
# ou: make platform-down
```

---

### Opção 2: Desenvolvimento Local (Go)

1. **Configurar variáveis de ambiente:**
   ```bash
   cp config.env.example config.env
   ```

2. **Iniciar o banco de dados e aplicar migrations:**
   ```bash
   make db-up
   ```

3. **Executar a API:**
   ```bash
   make run
   # ou: go run cmd/main.go
   ```

---

### Opção 3: Executar via Imagem Docker (GHCR)

```bash
# Baixar imagem oficial
docker pull ghcr.io/tech-preta/alderaan-api:latest

# Executar container
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

---

### Opção 4: Deploy no Kubernetes com Helm (OCI Registry)

```bash
# Instalação direta a partir do registro OCI oficial
helm install alderaan oci://ghcr.io/tech-preta/helm-charts/alderaan --version 1.0.0
```

Para customizações de produção, consulte o [Guia do Helm Chart](charts/README.md).

---

## 🛠️ Comandos Mais Usados (`Makefile`)

```bash
make help           # Lista todos os comandos disponíveis
make test           # Executa a suíte de testes com detecção de race condition
make run            # Atualiza documentação Swagger e executa a aplicação
make build          # Compila o binário otimizado para produção
make db-up          # Inicia PostgreSQL e roda migrations do Flyway
make db-down        # Encerra o banco de dados
make db-seed        # Popula o banco com dados de teste
make platform-up    # Sobe todo o ecossistema (App + DB + Observabilidade)
make platform-down  # Desliga todo o ecossistema
```

---

## 📚 Documentação Técnica Aprofundada

Para guias detalhados sobre decisões de design, padrões de arquitetura e infraestrutura, consulte o diretório [`docs/`](docs/README.md):

| Tópico | Documento |
| :--- | :--- |
| **Arquitetura & Design** | [Domain-Driven Design (DDD)](docs/01-domain-driven-design.md) • [Clean Architecture](docs/02-clean-architecture.md) • [Event Dispatcher](docs/03-event-dispatcher.md) |
| **API & Resiliência** | [RESTful API com Gin](docs/05-restful-api-gin.md) • [Graceful Shutdown](docs/04-graceful-shutdown.md) • [Swagger / OpenAPI](docs/06-swagger-documentation.md) • [Exemplos de Chamadas](docs/api-examples.md) |
| **Observabilidade** | [Prometheus Monitoring](docs/07-prometheus-monitoring.md) • [Guia de PromQL](docs/10-prometheus-queries.md) • [Refatoração de Métricas](docs/11-refactoring-metrics.md) |
| **Deploy & CI/CD** | [Docker Deployment](docs/08-docker-deployment.md) • [Helm Charts](charts/README.md) • [Releases e Versionamento](docs/12-automated-releases.md) |
| **Banco de Dados** | [Flyway Migrations](docs/09-flyway-migrations.md) • [Documentação do PostgreSQL](db/README.md) |

---

## 🤝 Contribuição e Licença

- **Commits**: Seguem o padrão [Conventional Commits](https://www.conventionalcommits.org/pt-br/) (`feat:`, `fix:`, `chore:`, etc.) para automação de releases.
- **Licença**: Distribuído sob a licença [MIT](LICENSE).
- **Créditos**: Baseado nos conceitos do artigo de [William Koller](https://williamkoller.substack.com).
