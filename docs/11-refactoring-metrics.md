# 🔄 Refatoração: Desacoplamento do Pacote Metrics

Documentação da refatoração que desacoplou o pacote `metrics` de `shared` para sua própria pasta em `internal/metrics`.

## 📝 Contexto

Originalmente, o pacote de métricas estava localizado em `internal/shared/metrics/`, o que implica que métricas são um componente "compartilhado" genérico. No entanto, métricas do Prometheus representam uma preocupação infraestrutural específica e bem definida, não necessariamente um utilitário compartilhado.

## 🎯 Objetivo

Melhorar a arquitetura do projeto seguindo o princípio de **Separation of Concerns** (Separação de Responsabilidades):

- ✅ **Antes**: `internal/shared/metrics/` - Métricas como utilitário compartilhado
- ✅ **Depois**: `internal/metrics/` - Métricas como módulo independente de infraestrutura

## 🏗️ Mudanças Realizadas

### **1. Estrutura de Diretórios**

```diff
internal/
├── domain/
├── infra/
+ ├── metrics/              ← NOVO: Pacote independente
+ │   ├── metrics.go
+ │   └── middleware.go
└── shared/
-   ├── metrics/            ← REMOVIDO
-   │   ├── metrics.go
-   │   └── middleware.go
    ├── config/
    ├── database/
    └── domain/
```

### **2. Imports Atualizados**

#### **cmd/main.go**

```diff
import (
    _ "github.com/williamkoller/golang-domain-driven-design/docs"
    product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
    product_handlers "github.com/williamkoller/golang-domain-driven-design/internal/infra/http/handlers"
    product_router "github.com/williamkoller/golang-domain-driven-design/internal/infra/http/router"
    "github.com/williamkoller/golang-domain-driven-design/internal/infra/persistence"
+   "github.com/williamkoller/golang-domain-driven-design/internal/metrics"
    "github.com/williamkoller/golang-domain-driven-design/internal/shared/config"
    "github.com/williamkoller/golang-domain-driven-design/internal/shared/database"
    shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
-   "github.com/williamkoller/golang-domain-driven-design/internal/shared/metrics"
)
```

#### **internal/infra/http/handlers/product_handler.go**

```diff
import (
    "net/http"

    "github.com/gin-gonic/gin"
    product_entity "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/entity"
    product_repository "github.com/williamkoller/golang-domain-driven-design/internal/domain/product/repository"
+   "github.com/williamkoller/golang-domain-driven-design/internal/metrics"
    shared_events "github.com/williamkoller/golang-domain-driven-design/internal/shared/domain/events"
-   "github.com/williamkoller/golang-domain-driven-design/internal/shared/metrics"
)
```

#### **internal/infra/http/router/product_router.go**

```diff
import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    product_handlers "github.com/williamkoller/golang-domain-driven-design/internal/infra/http/handlers"
+   "github.com/williamkoller/golang-domain-driven-design/internal/metrics"
-   "github.com/williamkoller/golang-domain-driven-design/internal/shared/metrics"
)
```

### **3. Comandos Executados**

```bash
# 1. Criar novo diretório
mkdir -p internal/metrics

# 2. Mover arquivos
mv internal/shared/metrics/*.go internal/metrics/

# 3. Remover diretório vazio
rmdir internal/shared/metrics

# 4. Atualizar imports (3 arquivos modificados)

# 5. Atualizar dependências
go mod tidy

# 6. Regenerar Swagger
make swagger

# 7. Compilar
make build
```

## ✅ Validação

### **Compilação**

```bash
$ make build
✅ Documentação gerada em docs/
✅ Binário compilado em bin/server
```

### **Testes**

```bash
# Iniciar servidor
$ ./bin/server
🚀 Server running on http://localhost:8080
📊 Metrics available at http://localhost:8080/metrics
📚 Swagger UI at http://localhost:8080/swagger/index.html

# Verificar health
$ curl http://localhost:8080/health
{"status":"ok"}

# Verificar métricas
$ curl http://localhost:8080/metrics | grep -E "^(http_|products_)"
http_request_duration_seconds_bucket{...} 0
http_requests_total{...} 1
products_total 9
products_created_total 0
...
```

## 🎨 Benefícios da Refatoração

### **1. Separação de Responsabilidades**

- **Antes**: `shared/metrics` mistura infraestrutura com utilitários
- **Depois**: `internal/metrics` é claramente um módulo de infraestrutura

### **2. Clareza Arquitetural**

```
internal/
├── domain/          → Lógica de negócio pura
├── infra/           → Infraestrutura (HTTP, persistência)
├── metrics/         → Infraestrutura de observabilidade ✨
└── shared/          → Utilitários verdadeiramente compartilhados
    ├── config/      → Configuração
    ├── database/    → Cliente de banco
    └── domain/      → Eventos de domínio
```

### **3. Coesão de Módulos**

- **Métricas** agora tem seu próprio namespace claro
- **Shared** contém apenas utilitários genéricos (config, database, domain events)
- Facilita encontrar e manter código relacionado

### **4. Preparação para Crescimento**

Com métricas isoladas, fica mais fácil:
- Adicionar novos tipos de métricas (tracing, logs)
- Criar exporters customizados
- Implementar métricas específicas por contexto
- Migrar para outras soluções de observabilidade

### **5. Melhora na Testabilidade**

```go
// Antes: import longo e confuso
import "github.com/.../internal/shared/metrics"

// Depois: import claro e direto
import "github.com/.../internal/metrics"
```

## 🔍 Impacto da Mudança

### **Arquivos Modificados**: 3

1. `cmd/main.go` - Import path atualizado
2. `internal/infra/http/handlers/product_handler.go` - Import path atualizado
3. `internal/infra/http/router/product_router.go` - Import path atualizado

### **Arquivos Movidos**: 2

1. `internal/shared/metrics/metrics.go` → `internal/metrics/metrics.go`
2. `internal/shared/metrics/middleware.go` → `internal/metrics/middleware.go`

### **Diretórios Removidos**: 1

- `internal/shared/metrics/` (vazio após mover arquivos)

### **Compatibilidade**

- ✅ **API pública**: Nenhuma mudança (endpoints iguais)
- ✅ **Métricas**: Nenhuma mudança (nomes iguais)
- ✅ **Swagger**: Regenerado automaticamente
- ✅ **Docker**: Funciona sem modificações
- ✅ **Tests**: Nenhum impacto (se existissem)

## 📚 Princípios Aplicados

### **1. Single Responsibility Principle (SRP)**

Cada módulo tem uma responsabilidade clara:
- `metrics/` → Observabilidade com Prometheus
- `shared/config/` → Configuração da aplicação
- `shared/database/` → Conexão com banco de dados

### **2. Package by Feature/Layer**

Seguindo a estrutura de Clean Architecture:
```
internal/
├── domain/          → Camada de Domínio
├── infra/           → Camada de Infraestrutura
│   ├── http/
│   └── persistence/
├── metrics/         → Camada de Observabilidade
└── shared/          → Componentes Cross-cutting
```

### **3. Dependency Inversion**

Métricas são passadas por injeção de dependência:
```go
// main.go
m := metrics.NewMetrics()
productHandler := product_handlers.NewProductHandler(repo, dispatcher, m)
router := product_router.SetupProductRouter(productHandler, m)
```

## 🚀 Próximos Passos

Com essa refatoração, o projeto está preparado para:

1. **Adicionar Tracing** (ex: Jaeger, Tempo)
   ```
   internal/
   ├── metrics/     → Prometheus
   ├── tracing/     → OpenTelemetry/Jaeger
   └── logging/     → Structured logging
   ```

2. **Métricas por Contexto**
   ```go
   internal/metrics/
   ├── http/        → Métricas HTTP
   ├── database/    → Métricas de DB
   └── business/    → Métricas de negócio
   ```

3. **Exporters Customizados**
   ```go
   internal/metrics/
   ├── prometheus/  → Exporter Prometheus
   ├── datadog/     → Exporter DataDog
   └── cloudwatch/  → Exporter CloudWatch
   ```

## 📖 Referências

- **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- **Package Organization**: https://go.dev/doc/modules/layout
- **Domain-Driven Design**: https://martinfowler.com/bliki/DomainDrivenDesign.html
- **Separation of Concerns**: https://en.wikipedia.org/wiki/Separation_of_concerns

## 🎓 Lições Aprendidas

### **O que "shared" realmente significa?**

❌ **Não é shared**:
- Código específico de infraestrutura (métricas, tracing)
- Código de domínio específico
- Implementações concretas de serviços

✅ **É shared**:
- Utilitários genéricos (helpers, validators)
- Configuração da aplicação
- Interfaces compartilhadas entre camadas
- Eventos de domínio (DDD)

### **Quando refatorar?**

Considere mover para fora de `shared/` quando:
- O módulo tem responsabilidade específica clara
- O módulo pode crescer independentemente
- O módulo representa uma camada arquitetural
- O nome do módulo descreve melhor sua função que "shared"

## 🔗 Commits Relacionados

```bash
# Ver histórico da mudança
git log --oneline --all --graph -- internal/metrics/
git log --oneline --all --graph -- internal/shared/metrics/

# Ver diff completo
git diff HEAD~1 HEAD -- internal/
```

---

**Data da Refatoração**: 2024-10-04
**Tipo**: Structural Refactoring
**Impacto**: Baixo (apenas imports)
**Status**: ✅ Completo e Testado
