# Arquitetura Limpa (Clean Architecture)

## 📚 O que é Clean Architecture?

**Clean Architecture (Arquitetura Limpa)** é um padrão arquitetural proposto por Robert C. Martin (Uncle Bob) que visa criar sistemas:

- **Independentes de frameworks**
- **Testáveis**
- **Independentes da UI**
- **Independentes do banco de dados**
- **Independentes de qualquer agente externo**

## 🎯 Princípios Fundamentais

### 1. **Regra de Dependência**

> As dependências sempre apontam para dentro, em direção às regras de negócio de alto nível.

```
Infraestrutura → Aplicação → Domínio
      ↓              ↓           ↓
   (Externo)    (Casos de Uso) (Entidades)
```

**Camadas externas conhecem as internas, mas NUNCA o contrário.**

### 2. **Inversão de Dependência**

Use interfaces para inverter dependências, permitindo que o domínio defina contratos que a infraestrutura implementa.

```go
// Domínio define a interface
type IProductRepository interface {
    Add(product Product) error
    Find() ([]Product, error)
}

// Infraestrutura implementa
type ProductRepository struct {
    db *sql.DB
}

func (r *ProductRepository) Add(product Product) error {
    // Implementação com banco de dados
}
```

## 🏗️ Camadas da Clean Architecture

### **1. Entities (Entidades) - Centro**

Contém as **regras de negócio críticas** da empresa.

**Características:**
- Puras, sem dependências externas
- Contêm lógica de negócio essencial
- Podem ser usadas por múltiplas aplicações

```go
// internal/domain/product/entity/product.go
type Product struct {
    Name       string
    Sku        int
    Categories []string
    Price      int
}

func (p *Product) ApplyDiscount(percentage float64) error {
    if percentage < 0 || percentage > 100 {
        return errors.New("invalid discount percentage")
    }
    p.Price = int(float64(p.Price) * (1 - percentage/100))
    return nil
}
```

### **2. Use Cases (Casos de Uso) - Camada de Aplicação**

Contém as **regras de negócio específicas da aplicação**.

**Características:**
- Orquestra o fluxo de dados
- Coordena entidades
- Não depende de detalhes de implementação

```go
type CreateProductUseCase struct {
    repo       IProductRepository
    dispatcher EventDispatcher
}

func (uc *CreateProductUseCase) Execute(input CreateProductInput) error {
    product, event, err := NewProduct(
        input.Name, 
        input.Sku, 
        input.Categories, 
        input.Price,
        uc.dispatcher,
    )
    if err != nil {
        return err
    }
    
    return uc.repo.Add(*product)
}
```

### **3. Interface Adapters (Adaptadores de Interface)**

Convertem dados entre casos de uso e agentes externos.

**Exemplos:**
- Controllers HTTP
- Presenters
- Gateways
- Repositórios

```go
// internal/infra/http/handlers/product_handler.go
type ProductHandler struct {
    repo       *ProductRepository
    dispatcher *EventDispatcher
}

func (h *ProductHandler) Create(c *gin.Context) {
    var input struct {
        Name       string   `json:"name"`
        Sku        int      `json:"sku"`
        Categories []string `json:"categories"`
        Price      int      `json:"price"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    product, _, err := product_entity.NewProduct(
        input.Name, 
        input.Sku, 
        input.Categories, 
        input.Price, 
        h.dispatcher,
    )
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.repo.Add(*product); err != nil {
        c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, product)
}
```

### **4. Frameworks & Drivers (Frameworks e Drivers) - Camada Externa**

Contém frameworks, ferramentas e detalhes de implementação.

**Exemplos:**
- Gin (framework HTTP)
- Banco de dados
- APIs externas
- Bibliotecas de terceiros

## 📁 Estrutura de Pastas no Projeto

```
internal/
├── domain/              # Camada de Domínio (Entidades)
│   └── product/
│       ├── entity/
│       ├── events/
│       └── repository/  # Interfaces
├── application/         # Casos de Uso (não implementado neste exemplo)
│   └── usecases/
├── infra/              # Infraestrutura (Frameworks & Drivers)
│   ├── http/
│   │   ├── handlers/   # Adaptadores HTTP
│   │   └── router/
│   └── persistence/    # Implementação de repositórios
└── shared/             # Código compartilhado
    └── domain/
        └── events/
```

## 🔄 Fluxo de Dados

```
1. HTTP Request → Handler (Interface Adapter)
2. Handler → Use Case (Application Layer)
3. Use Case → Entity (Domain Layer)
4. Entity → Valida e processa
5. Use Case → Repository Interface (Domain)
6. Repository Implementation → Banco de Dados (Infrastructure)
7. Resultado retorna camada por camada
8. Handler → HTTP Response
```

## ✅ Benefícios

### 1. **Testabilidade**
```go
// Fácil mockar repositório para testes
type MockProductRepository struct {}

func (m *MockProductRepository) Add(product Product) error {
    return nil
}

func TestCreateProduct(t *testing.T) {
    mockRepo := &MockProductRepository{}
    dispatcher := NewEventDispatcher()
    
    handler := NewProductHandler(mockRepo, dispatcher)
    
    // Teste sem dependências externas
}
```

### 2. **Independência de Framework**

Se decidir trocar Gin por Echo ou Fiber, apenas a camada de infraestrutura muda:

```go
// Antes: Gin
func (h *ProductHandler) Create(c *gin.Context) { ... }

// Depois: Echo
func (h *ProductHandler) Create(c echo.Context) error { ... }
```

O domínio permanece intacto! ✅

### 3. **Independência de Banco de Dados**

```go
// Hoje: In-Memory
type ProductRepository struct {
    data map[string]Product
}

// Amanhã: PostgreSQL
type ProductRepository struct {
    db *sql.DB
}

// Depois: MongoDB
type ProductRepository struct {
    collection *mongo.Collection
}
```

A interface `IProductRepository` não muda!

### 4. **Manutenibilidade**

Cada camada tem responsabilidade clara:
- **Domínio**: Regras de negócio
- **Aplicação**: Orquestração
- **Infraestrutura**: Detalhes técnicos

### 5. **Escalabilidade**

Fácil adicionar novos recursos sem quebrar código existente.

## 🎯 Princípios SOLID Aplicados

### **S - Single Responsibility Principle**
Cada classe/módulo tem uma única responsabilidade.

### **O - Open/Closed Principle**
Aberto para extensão, fechado para modificação.

### **L - Liskov Substitution Principle**
Implementações podem ser substituídas sem quebrar o sistema.

### **I - Interface Segregation Principle**
Interfaces específicas são melhores que uma interface geral.

### **D - Dependency Inversion Principle**
Dependa de abstrações, não de implementações concretas.

## 🚫 Anti-Padrões a Evitar

### ❌ **Acoplamento entre Camadas**
```go
// ERRADO: Domínio importando HTTP
package product_entity

import "github.com/gin-gonic/gin"  // ❌

type Product struct {
    Name string
}

func (p *Product) ToJSON(c *gin.Context) { ... }  // ❌
```

### ❌ **Lógica de Negócio no Handler**
```go
// ERRADO: Validação de negócio no handler HTTP
func (h *ProductHandler) Create(c *gin.Context) {
    if input.Price <= 0 {  // ❌ Isso deve estar no domínio!
        return errors.New("price must be positive")
    }
}
```

### ✅ **CORRETO: Validação no Domínio**
```go
// internal/domain/product/entity/product.go
func Validate(price int) error {
    if price <= 0 {
        return errors.New("price must be positive")
    }
    return nil
}
```

## 💡 Exemplo Completo no Projeto

### **Domínio (Centro)**
```go
// internal/domain/product/entity/product.go
type Product struct {
    Name       string
    Sku        int
    Categories []string
    Price      int
}
```

### **Interface de Repositório (Domínio)**
```go
// internal/domain/product/repository/product_repository.go
type IProductRepository interface {
    Add(product Product) error
    Find() ([]Product, error)
}
```

### **Implementação (Infraestrutura)**
```go
// internal/domain/product/repository/product_repository.go
type ProductRepository struct {
    data map[string]Product
    mu   sync.RWMutex
}

func (r *ProductRepository) Add(product Product) error {
    // Implementação específica
}
```

### **Handler HTTP (Interface Adapter)**
```go
// internal/infra/http/handlers/product_handler.go
type ProductHandler struct {
    repo       *ProductRepository
    dispatcher *EventDispatcher
}

func (h *ProductHandler) Create(c *gin.Context) {
    // Adapta HTTP para domínio
}
```

## 📚 Recursos

- **Livro:** "Clean Architecture" - Robert C. Martin
- **Livro:** "Clean Code" - Robert C. Martin
- **Blog:** https://blog.cleancoder.com/

---

**Anterior:** [Domain-Driven Design](01-domain-driven-design.md) | **Próximo:** [Event Dispatcher](03-event-dispatcher.md)

