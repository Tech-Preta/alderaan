# Multi-stage build para otimizar tamanho da imagem

# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Instalar dependências necessárias para compilação
RUN apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar go.mod e go.sum primeiro (melhor cache de layers)
COPY go.mod go.sum ./

# Download de dependências (será cacheado se não houver mudanças)
RUN go mod download

# Copiar código fonte
COPY . .

# Gerar documentação Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go -o docs

# Build do binário com otimizações
# CGO_ENABLED=0 para criar binário estático
# -ldflags para reduzir tamanho (remover debug info)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server \
    cmd/main.go

# Stage 2: Runtime
FROM alpine:3.19

# Instalar certificados CA para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root para segurança
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /app

# Copiar binário do stage de build
COPY --from=builder /app/server /app/server

# Copiar documentação Swagger gerada
COPY --from=builder /app/docs /app/docs

# Copiar certificados e timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Mudar ownership para usuário não-root
RUN chown -R appuser:appgroup /app

# Mudar para usuário não-root
USER appuser

# Expor porta da aplicação
EXPOSE 8080

# Variáveis de ambiente padrão (podem ser sobrescritas no docker-compose)
ENV DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=alderaan \
    DB_PASSWORD=alderaan123 \
    DB_NAME=alderaan_db \
    DB_SSLMODE=disable \
    SERVER_PORT=8080 \
    GIN_MODE=release

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Executar aplicação
CMD ["/app/server"]

