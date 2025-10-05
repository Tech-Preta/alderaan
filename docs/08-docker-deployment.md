# 🐳 Docker & Deployment

Este guia explica como construir, executar e fazer deploy da aplicação usando Docker.

## 📦 Dockerfile Multi-Stage Build

O projeto usa **multi-stage build** para otimizar o tamanho da imagem final.

### **Stage 1: Builder (golang:1.24-alpine)**

Responsável por:
- ✅ Download de dependências Go
- ✅ Geração da documentação Swagger
- ✅ Compilação do binário estático

### **Stage 2: Runtime (alpine:3.19)**

Responsável por:
- ✅ Executar apenas o binário compilado
- ✅ Imagem mínima (~20MB vs ~800MB do builder)
- ✅ Segurança com usuário não-root
- ✅ Health check configurado

## 🏗️ Construindo a Imagem

### **Build Manual**

```bash
# Build da imagem
docker build -t alderaan-api:latest .

# Build com tag de versão
docker build -t alderaan-api:1.0.0 .

# Build sem cache
docker build --no-cache -t alderaan-api:latest .

# Ou use Makefile
make docker-build
```

### **Build via Docker Compose**

```bash
# Build automático ao iniciar
docker-compose up -d --build

# Build sem iniciar
docker-compose build
```

## 🚀 Executando a Aplicação

### **Container Individual**

```bash
# Executar com variáveis de ambiente
docker run -d \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5432 \
  -e DB_USER=alderaan \
  -e DB_PASSWORD=alderaan123 \
  -e DB_NAME=alderaan_db \
  --name alderaan-api \
  alderaan-api:latest

# Executar com arquivo .env
docker run -d \
  -p 8080:8080 \
  --env-file config.env \
  --name alderaan-api \
  alderaan-api:latest
```

### **Stack Completa (Recomendado)**

```bash
# Usar docker-compose
docker-compose up -d

# Ou Makefile
make platform-up
```

## 🔍 Inspecionando a Imagem

### **Tamanho da Imagem**

```bash
# Ver tamanho
docker images alderaan-api:latest

# Comparar stages
docker images | grep golang  # ~800MB
docker images | grep alderaan # ~20MB
```

### **Camadas da Imagem**

```bash
# Ver histórico de camadas
docker history alderaan-api:latest

# Análise detalhada com dive
docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  wagoodman/dive:latest alderaan-api:latest
```

### **Inspecionar Container**

```bash
# Ver logs
docker logs alderaan-api -f

# Executar shell no container
docker exec -it alderaan-api sh

# Ver processos
docker top alderaan-api

# Ver recursos consumidos
docker stats alderaan-api
```

## 🧪 Testando a Imagem

### **Health Check**

```bash
# Verificar saúde do container
docker inspect --format='{{json .State.Health}}' alderaan-api | jq

# Aguardar ficar healthy
until [ "$(docker inspect -f {{.State.Health.Status}} alderaan-api)" == "healthy" ]; do
    sleep 1
done
echo "Container is healthy!"
```

### **Teste de Conectividade**

```bash
# Teste do health endpoint
curl http://localhost:8080/health

# Teste da API
curl http://localhost:8080/api/v1/products

# Teste das métricas
curl http://localhost:8080/metrics
```

## 🔐 Segurança

### **Usuário Não-Root**

A aplicação roda como `appuser` (UID 1001) ao invés de root:

```dockerfile
USER appuser
```

**Verificar:**
```bash
docker exec alderaan-api whoami
# Output: appuser

docker exec alderaan-api id
# Output: uid=1001(appuser) gid=1001(appgroup)
```

### **Scan de Vulnerabilidades**

```bash
# Usando Docker Scout
docker scout cves alderaan-api:latest

# Usando Trivy
trivy image alderaan-api:latest

# Usando Snyk
snyk container test alderaan-api:latest
```

### **Boas Práticas Implementadas**

- ✅ **Imagem base oficial e mínima** (alpine)
- ✅ **Multi-stage build** (reduz tamanho e superfície de ataque)
- ✅ **Usuário não-root**
- ✅ **Binário estático** (sem dependências externas)
- ✅ **Health check** configurado
- ✅ **Certificados CA** incluídos para HTTPS
- ✅ **.dockerignore** otimizado

## 📊 Otimizações

### **Cache de Layers**

As camadas são organizadas para maximizar cache:

1. `COPY go.mod go.sum` - Muda raramente
2. `RUN go mod download` - Usa cache se go.mod não mudou
3. `COPY . .` - Código fonte (muda frequentemente)
4. `RUN go build` - Só recompila se código mudou

### **Tamanho da Imagem**

```
golang:1.24-alpine (builder)  ~800MB  ❌ Não vai para produção
alpine:3.19 (runtime)         ~7MB   ✅
+ binário Go                  ~15MB  ✅
+ certificados + timezone     ~2MB   ✅
= Imagem final                ~24MB  ✅
```

### **Build Time**

```bash
# Build completo (sem cache)
time docker build -t alderaan-api:latest .
# ~60-90 segundos

# Build com cache (sem mudanças)
time docker build -t alderaan-api:latest .
# ~5 segundos

# Build com cache (apenas código mudou)
time docker build -t alderaan-api:latest .
# ~20-30 segundos
```

## 🌐 Deploy em Produção

### **Kubernetes**

Arquivo de deployment em `charts/`:

```bash
# Deploy usando Helm
helm install alderaan ./charts/alderaan

# Ou kubectl direto
kubectl apply -f k8s/
```

### **Docker Swarm**

```bash
# Inicializar swarm
docker swarm init

# Deploy da stack
docker stack deploy -c docker-compose.yml alderaan

# Ver serviços
docker service ls

# Escalar API
docker service scale alderaan_api=3
```

### **Railway / Render / Fly.io**

Esses serviços detectam o Dockerfile automaticamente:

```bash
# Railway
railway up

# Render
# Conecte o GitHub repo e configure:
# - Build Command: docker build
# - Start Command: automático

# Fly.io
fly deploy
```

## 🔄 CI/CD Pipeline

### **GitHub Actions**

Exemplo de workflow:

```yaml
name: Docker Build & Push

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build Docker image
        run: docker build -t alderaan-api:${{ github.sha }} .

      - name: Run tests
        run: docker run alderaan-api:${{ github.sha }} go test ./...

      - name: Login to DockerHub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Push to DockerHub
        run: |
          docker tag alderaan-api:${{ github.sha }} user/alderaan-api:latest
          docker push user/alderaan-api:latest
```

## 🐛 Troubleshooting

### **Erro: "failed to solve: the Dockerfile cannot be empty"**

```bash
# Verificar se Dockerfile existe e não está vazio
cat Dockerfile

# Rebuild
docker-compose build --no-cache
```

### **Erro: "bind: address already in use"**

```bash
# Verificar porta 8080 em uso
lsof -i :8080

# Parar container antigo
docker stop alderaan-api
docker rm alderaan-api
```

### **Container reinicia constantemente**

```bash
# Ver logs
docker logs alderaan-api --tail 100

# Ver motivo do restart
docker inspect alderaan-api | jq '.[0].State'

# Verificar health check
docker inspect alderaan-api | jq '.[0].State.Health'
```

### **Erro de conexão com banco**

```bash
# Verificar se postgres está rodando
docker ps | grep postgres

# Verificar rede
docker network inspect alderaan-network

# Testar conexão do container
docker exec alderaan-api wget -O- postgres:5432
```

### **Imagem muito grande**

```bash
# Verificar o que está ocupando espaço
docker history alderaan-api:latest

# Usar dive para análise
docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  wagoodman/dive:latest alderaan-api:latest

# Verificar .dockerignore
cat .dockerignore
```

## 📝 Variáveis de Ambiente

### **Banco de Dados**

```bash
DB_HOST=postgres          # Host do PostgreSQL
DB_PORT=5432             # Porta do PostgreSQL
DB_USER=alderaan         # Usuário do banco
DB_PASSWORD=alderaan123  # Senha do banco
DB_NAME=alderaan_db      # Nome do banco
DB_SSLMODE=disable       # Modo SSL (disable/require)
```

### **Servidor**

```bash
SERVER_PORT=8080         # Porta da API
GIN_MODE=release         # Modo do Gin (debug/release)
```

### **Sobrescrever no Docker Compose**

```yaml
services:
  api:
    environment:
      - DB_HOST=my-postgres-host
      - GIN_MODE=debug
```

## 📊 Monitoramento do Container

### **Logs**

```bash
# Logs em tempo real
docker logs -f alderaan-api

# Últimas 100 linhas
docker logs --tail 100 alderaan-api

# Logs com timestamp
docker logs -t alderaan-api
```

### **Recursos**

```bash
# Uso de CPU e memória
docker stats alderaan-api

# Limites de recursos (configurar no compose)
services:
  api:
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

### **Health Status**

```bash
# Status do health check
docker inspect alderaan-api \
  --format='{{.State.Health.Status}}'

# Últimas verificações
docker inspect alderaan-api \
  --format='{{json .State.Health}}' | jq
```

## 🎓 Comandos Úteis

```bash
# Rebuild completo
docker-compose up -d --build --force-recreate

# Ver imagens
docker images

# Remover imagens antigas
docker image prune -a

# Ver containers
docker ps -a

# Parar todos os containers
docker stop $(docker ps -q)

# Remover todos os containers parados
docker container prune

# Ver uso de disco do Docker
docker system df

# Limpar tudo (CUIDADO!)
docker system prune -a --volumes
```

## 📚 Recursos

- **Dockerfile Best Practices**: https://docs.docker.com/develop/dev-best-practices/
- **Multi-stage Builds**: https://docs.docker.com/build/building/multi-stage/
- **Docker Security**: https://docs.docker.com/engine/security/
- **Go + Docker**: https://docs.docker.com/language/golang/

---

**Anterior:** [Prometheus Monitoring](07-prometheus-monitoring.md) | **Voltar ao início:** [README](README.md)
