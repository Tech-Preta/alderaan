.PHONY: help swagger build run clean test

# Cores para output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

help: ## Mostra esta mensagem de ajuda
	@echo ''
	@echo 'Uso:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "  ${YELLOW}%-15s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${WHITE}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

swagger: ## Gera a documentação Swagger
	@echo "${GREEN}Gerando documentação Swagger...${RESET}"
	@swag init -g cmd/main.go -o docs
	@echo "${GREEN}✅ Documentação gerada em docs/${RESET}"

build: swagger ## Compila o binário do servidor
	@echo "${GREEN}Compilando servidor...${RESET}"
	@go build -o bin/server cmd/main.go
	@echo "${GREEN}✅ Binário compilado em bin/server${RESET}"

run: swagger ## Gera documentação e executa o servidor
	@echo "${GREEN}Iniciando servidor...${RESET}"
	@go run cmd/main.go

dev: ## Executa em modo desenvolvimento (sem regenerar swagger)
	@echo "${GREEN}Iniciando servidor em modo dev...${RESET}"
	@go run cmd/main.go

clean: ## Remove binários e arquivos temporários
	@echo "${YELLOW}Limpando arquivos...${RESET}"
	@rm -rf bin/
	@rm -rf tmp/
	@echo "${GREEN}✅ Limpeza concluída${RESET}"

test: ## Executa os testes
	@echo "${GREEN}Executando testes...${RESET}"
	@go test -v ./...

test-coverage: ## Executa os testes com cobertura
	@echo "${GREEN}Executando testes com cobertura...${RESET}"
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@echo "${GREEN}Cobertura geral:${RESET}"
	@go tool cover -func=coverage.out | grep total

test-coverage-html: ## Gera relatório HTML de cobertura
	@echo "${GREEN}Gerando relatório HTML de cobertura...${RESET}"
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✅ Relatório gerado em coverage.html${RESET}"

test-unit: ## Executa apenas testes unitários (rápidos)
	@echo "${GREEN}Executando testes unitários...${RESET}"
	@go test -short -v ./internal/domain/... ./internal/metrics/... ./internal/shared/...

deps: ## Instala as dependências
	@echo "${GREEN}Instalando dependências...${RESET}"
	@go mod download
	@go mod tidy
	@echo "${GREEN}✅ Dependências instaladas${RESET}"

install-swag: ## Instala o CLI do swag
	@echo "${GREEN}Instalando swag CLI...${RESET}"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "${GREEN}✅ Swag CLI instalado${RESET}"

docker-build: ## Constrói a imagem Docker
	@echo "${GREEN}Construindo imagem Docker...${RESET}"
	@docker build -t alderaan-api:latest .
	@echo "${GREEN}✅ Imagem Docker construída${RESET}"

docker-run: ## Executa o container Docker
	@echo "${GREEN}Executando container...${RESET}"
	@docker run -p 8080:8080 alderaan-api:latest

monitoring-up: ## Inicia Prometheus, Alertmanager e Grafana
	@echo "${GREEN}Iniciando stack de monitoramento...${RESET}"
	@docker-compose up -d prometheus alertmanager grafana
	@echo "${GREEN}✅ Prometheus: http://localhost:9090${RESET}"
	@echo "${GREEN}✅ Alertmanager: http://localhost:9093${RESET}"
	@echo "${GREEN}✅ Grafana: http://localhost:3000 (admin/admin)${RESET}"

monitoring-down: ## Para Prometheus, Alertmanager e Grafana
	@echo "${YELLOW}Parando stack de monitoramento...${RESET}"
	@docker-compose stop prometheus alertmanager grafana
	@echo "${GREEN}✅ Monitoramento parado${RESET}"

monitoring-logs: ## Ver logs do monitoramento
	@docker-compose logs -f prometheus alertmanager grafana

check-metrics: ## Verifica se as métricas estão acessíveis
	@echo "${GREEN}Verificando endpoint /metrics...${RESET}"
	@curl -s http://localhost:8080/metrics | head -20
	@echo "\n${GREEN}✅ Métricas disponíveis!${RESET}"

## Database commands
db-up: ## Inicia PostgreSQL com Docker
	@echo "${GREEN}Iniciando PostgreSQL...${RESET}"
	@docker-compose up -d postgres
	@echo "${GREEN}✅ PostgreSQL rodando em localhost:5432${RESET}"

db-down: ## Para PostgreSQL
	@echo "${YELLOW}Parando PostgreSQL...${RESET}"
	@docker-compose stop postgres
	@echo "${GREEN}✅ PostgreSQL parado${RESET}"

db-logs: ## Ver logs do PostgreSQL
	@docker-compose logs -f postgres

db-connect: ## Conectar ao PostgreSQL
	@docker-compose exec postgres psql -U alderaan -d alderaan_db

## Flyway Migration commands
db-migrate: ## Rodar migrations com Flyway
	@echo "${GREEN}🚀 Executando migrations com Flyway...${RESET}"
	@docker-compose up flyway
	@echo "${GREEN}✅ Migrations executadas${RESET}"

db-migrate-info: ## Ver status das migrations
	@echo "${GREEN}📊 Status das migrations:${RESET}"
	@docker-compose run --rm flyway info

db-migrate-validate: ## Validar migrations
	@echo "${GREEN}✓ Validando migrations...${RESET}"
	@docker-compose run --rm flyway validate
	@echo "${GREEN}✅ Migrations válidas${RESET}"

db-migrate-repair: ## Reparar histórico de migrations
	@echo "${YELLOW}🔧 Reparando histórico de migrations...${RESET}"
	@docker-compose run --rm flyway repair
	@echo "${GREEN}✅ Histórico reparado${RESET}"

db-migrate-baseline: ## Criar baseline das migrations
	@echo "${GREEN}📌 Criando baseline...${RESET}"
	@docker-compose run --rm flyway baseline
	@echo "${GREEN}✅ Baseline criado${RESET}"

db-seed: ## Popular banco com dados de exemplo
	@echo "${GREEN}🌱 Populando banco com dados...${RESET}"
	@docker-compose exec postgres psql -U alderaan -d alderaan_db -f /flyway/sql/seed.sql
	@echo "${GREEN}✅ Dados inseridos${RESET}"

db-clean: ## Limpar todos os dados (CUIDADO!)
	@echo "${YELLOW}⚠️  Limpando banco de dados...${RESET}"
	@docker-compose exec postgres psql -U alderaan -d alderaan_db -c "TRUNCATE products, categories, product_categories RESTART IDENTITY CASCADE;"
	@echo "${GREEN}✅ Banco limpo${RESET}"

db-clean-all: ## Limpar banco e histórico Flyway (CUIDADO!)
	@echo "${YELLOW}⚠️  Limpando banco e histórico Flyway...${RESET}"
	@docker-compose run --rm flyway clean
	@echo "${GREEN}✅ Banco e histórico limpos${RESET}"

db-reset: ## Recriar banco do zero (CUIDADO!)
	@echo "${YELLOW}⚠️  Resetando banco de dados...${RESET}"
	@make db-clean-all
	@make db-migrate
	@echo "${GREEN}✅ Banco resetado${RESET}"

## Platform commands
platform-up: ## Inicia toda a plataforma (DB + API + Monitoring)
	@echo "${GREEN}🚀 Iniciando plataforma completa...${RESET}"
	@docker-compose up -d
	@echo ""
	@echo "${GREEN}✅ Plataforma iniciada com sucesso!${RESET}"
	@echo ""
	@echo "Serviços disponíveis:"
	@echo "  🗄️  PostgreSQL:   http://localhost:5432"
	@echo "  🚀 API:           http://localhost:8080"
	@echo "  📚 Swagger:       http://localhost:8080/swagger/index.html"
	@echo "  📊 Prometheus:    http://localhost:9090"
	@echo "  🚨 Alertmanager:  http://localhost:9093"
	@echo "  📈 Grafana:       http://localhost:3000 (admin/admin)"
	@echo ""

platform-down: ## Para toda a plataforma
	@echo "${YELLOW}Parando plataforma...${RESET}"
	@docker-compose down
	@echo "${GREEN}✅ Plataforma parada${RESET}"

platform-logs: ## Ver logs de todos os serviços
	@docker-compose logs -f

platform-status: ## Ver status de todos os serviços
	@docker-compose ps

platform-restart: ## Reinicia toda a plataforma
	@echo "${YELLOW}Reiniciando plataforma...${RESET}"
	@docker-compose restart
	@echo "${GREEN}✅ Plataforma reiniciada${RESET}"
