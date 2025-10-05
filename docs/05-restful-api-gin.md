# RESTful API com Gin

## 📚 O que é Gin?

**Gin** é um framework web HTTP de alta performance para Go, inspirado no Martini, mas até **40x mais rápido** devido ao uso de httprouter.

**Características principais:**
- ⚡ Extremamente rápido
- 🛣️ Roteamento baseado em radix tree
- 📝 Validação de JSON
- 🔒 Middleware robusto
- 🎯 API simples e intuitiva

## 🎯 O que é REST?

**REST (Representational State Transfer)** é um estilo arquitetural para APIs que usa HTTP de forma semântica.

### Princípios REST:

1. **Client-Server**: Separação entre cliente e servidor
2. **Stateless**: Cada requisição contém toda informação necessária
3. **Cacheable**: Respostas podem ser cacheadas
4. **Uniform Interface**: Interface consistente
5. **Layered System**: Arquitetura em camadas

## 🛣️ Métodos HTTP

| Método | Operação | Idempotente | Seguro |
|--------|----------|-------------|---------|
| GET    | Ler recursos | ✅ | ✅ |
| POST   | Criar recurso | ❌ | ❌ |
| PUT    | Atualizar (completo) | ✅ | ❌ |
| PATCH  | Atualizar (parcial) | ❌ | ❌ |
| DELETE | Remover recurso | ✅ | ❌ |

## 📦 Instalação

```bash
go get -u github.com/gin-gonic/gin
```

**go.mod:**
```go
module github.com/seu-usuario/seu-projeto

go 1.23

require github.com/gin-gonic/gin v1.10.0
```

## 🚀 Exemplo Básico

```go
package main

import "github.com/gin-gonic/gin"

func main() {
    // Criar router com middleware padrão (logger e recovery)
    r := gin.Default()

    // Definir rota
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    // Iniciar servidor
    r.Run(":8080")  // localhost:8080
}
```

## 🛣️ Roteamento

### **Rotas Básicas**

```go
r := gin.Default()

// GET - Listar/Buscar
r.GET("/products", listProducts)

// POST - Criar
r.POST("/products", createProduct)

// PUT - Atualizar completo
r.PUT("/products/:id", updateProduct)

// PATCH - Atualizar parcial
r.PATCH("/products/:id", patchProduct)

// DELETE - Remover
r.DELETE("/products/:id", deleteProduct)
```

### **Parâmetros de Rota**

```go
// URL: /products/123
r.GET("/products/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{"product_id": id})
})

// URL: /users/john/books/golang
r.GET("/users/:name/books/:title", func(c *gin.Context) {
    name := c.Param("name")
    title := c.Param("title")
    c.JSON(200, gin.H{
        "user": name,
        "book": title,
    })
})
```

### **Query Parameters**

```go
// URL: /search?q=golang&page=1
r.GET("/search", func(c *gin.Context) {
    query := c.Query("q")           // "golang"
    page := c.DefaultQuery("page", "1")  // "1" (valor padrão)
    
    c.JSON(200, gin.H{
        "query": query,
        "page":  page,
    })
})
```

### **Grupos de Rotas**

```go
r := gin.Default()

// API v1
v1 := r.Group("/api/v1")
{
    v1.POST("/products", createProduct)
    v1.GET("/products", listProducts)
    v1.GET("/products/:name", getProduct)
}

// API v2
v2 := r.Group("/api/v2")
{
    v2.POST("/products", createProductV2)
    v2.GET("/products", listProductsV2)
}

// Admin
admin := r.Group("/admin")
admin.Use(AuthMiddleware())  // Middleware para autenticação
{
    admin.GET("/users", listUsers)
    admin.DELETE("/users/:id", deleteUser)
}
```

## 📝 Binding e Validação

### **JSON Binding**

```go
type Product struct {
    Name       string   `json:"name" binding:"required"`
    Sku        int      `json:"sku" binding:"required,gt=0"`
    Categories []string `json:"categories" binding:"required,min=1"`
    Price      int      `json:"price" binding:"required,gt=0"`
}

func createProduct(c *gin.Context) {
    var product Product

    // ShouldBindJSON valida automaticamente
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Produto validado e pronto para uso
    c.JSON(http.StatusCreated, product)
}
```

### **Tags de Validação**

```go
type User struct {
    Name     string `binding:"required,min=3,max=50"`
    Email    string `binding:"required,email"`
    Age      int    `binding:"required,gte=18,lte=120"`
    Password string `binding:"required,min=8"`
    Website  string `binding:"omitempty,url"`
}
```

**Tags disponíveis:**
- `required` - Campo obrigatório
- `min`, `max` - Tamanho mínimo/máximo
- `gte`, `lte`, `gt`, `lt` - Comparações numéricas
- `email` - Validar email
- `url` - Validar URL
- `omitempty` - Campo opcional

### **Form Binding**

```go
type LoginForm struct {
    Username string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required"`
}

r.POST("/login", func(c *gin.Context) {
    var form LoginForm
    
    if err := c.ShouldBind(&form); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, gin.H{"username": form.Username})
})
```

## 🎨 Responses

### **JSON Response**

```go
// Objeto
c.JSON(200, gin.H{
    "message": "success",
    "data": product,
})

// Array
c.JSON(200, []Product{product1, product2})

// Struct
c.JSON(200, product)
```

### **Status Codes**

```go
import "net/http"

// Sucesso
c.JSON(http.StatusOK, data)              // 200
c.JSON(http.StatusCreated, data)         // 201
c.JSON(http.StatusNoContent, nil)        // 204

// Erro do Cliente
c.JSON(http.StatusBadRequest, error)     // 400
c.JSON(http.StatusUnauthorized, error)   // 401
c.JSON(http.StatusForbidden, error)      // 403
c.JSON(http.StatusNotFound, error)       // 404
c.JSON(http.StatusConflict, error)       // 409

// Erro do Servidor
c.JSON(http.StatusInternalServerError, error)  // 500
```

### **Outros Formatos**

```go
// XML
c.XML(200, gin.H{"message": "hello"})

// String
c.String(200, "Hello %s", name)

// HTML
c.HTML(200, "index.html", gin.H{"title": "Home"})

// File
c.File("./public/image.png")

// Redirect
c.Redirect(http.StatusMovedPermanently, "https://google.com")
```

## 🔧 Middleware

Middleware são funções executadas antes/depois dos handlers.

### **Middleware Básico**

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        // Processar requisição
        c.Next()

        // Após handler
        latency := time.Since(start)
        status := c.Writer.Status()
        
        log.Printf("[%d] %s %s - %v", status, c.Request.Method, 
            c.Request.URL.Path, latency)
    }
}

// Usar middleware
r := gin.New()  // Sem middleware padrão
r.Use(Logger())
r.Use(gin.Recovery())  // Recovery de panics
```

### **Middleware de Autenticação**

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()  // Para execução
            return
        }

        // Validar token
        if !isValidToken(token) {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }

        // Passar dados para próximo handler
        c.Set("user_id", getUserIDFromToken(token))
        
        c.Next()  // Continuar para próximo handler
    }
}

// Usar em rotas específicas
admin := r.Group("/admin")
admin.Use(AuthMiddleware())
{
    admin.GET("/users", listUsers)
}
```

### **Middleware de CORS**

```go
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", 
            "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", 
            "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

r.Use(CORSMiddleware())
```

## 💡 Implementação Completa - CRUD de Produtos

### **Router**

```go
// internal/infra/http/router/product_router.go
package product_router

import (
    "github.com/gin-gonic/gin"
    product_handlers "seu-projeto/internal/infra/http/handlers"
)

func SetupProductRouter(productHandler *product_handlers.ProductHandler) *gin.Engine {
    r := gin.Default()

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API v1
    v1 := r.Group("/api/v1")
    {
        // Produtos
        products := v1.Group("/products")
        {
            products.POST("", productHandler.Create)
            products.GET("", productHandler.FindAll)
            products.GET("/:name", productHandler.FindOne)
            products.PUT("/:name", productHandler.Update)
            products.DELETE("/:name", productHandler.Delete)
        }
    }

    return r
}
```

### **Handler**

```go
// internal/infra/http/handlers/product_handler.go
package product_handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    product_entity "seu-projeto/internal/domain/product/entity"
    product_repository "seu-projeto/internal/domain/product/repository"
    shared_events "seu-projeto/internal/shared/domain/events"
)

type ProductHandler struct {
    repo       *product_repository.ProductRepository
    dispatcher *shared_events.EventDispatcher
}

func NewProductHandler(repo *product_repository.ProductRepository, 
    dispatcher *shared_events.EventDispatcher) *ProductHandler {
    return &ProductHandler{repo, dispatcher}
}

// POST /api/v1/products
func (h *ProductHandler) Create(c *gin.Context) {
    var input struct {
        Name       string   `json:"name" binding:"required"`
        Sku        int      `json:"sku" binding:"required,gt=0"`
        Categories []string `json:"categories" binding:"required,min=1"`
        Price      int      `json:"price" binding:"required,gt=0"`
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

// GET /api/v1/products
func (h *ProductHandler) FindAll(c *gin.Context) {
    products, err := h.repo.Find()
    if err != nil {
        c.JSON(http.StatusInternalServerError, 
            gin.H{"error": "failed to fetch products"})
        return
    }
    
    c.JSON(http.StatusOK, products)
}

// GET /api/v1/products/:name
func (h *ProductHandler) FindOne(c *gin.Context) {
    name := c.Param("name")
    
    product, err := h.repo.FindOne(name)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
        return
    }
    
    c.JSON(http.StatusOK, product)
}
```

## 🧪 Testando a API

### **Criar Produto**

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

### **Listar Produtos**

```bash
curl http://localhost:8080/api/v1/products
```

### **Buscar Produto**

```bash
curl http://localhost:8080/api/v1/products/Notebook
```

### **Health Check**

```bash
curl http://localhost:8080/health
```

## 🎯 Boas Práticas RESTful

### **1. Nomear Recursos com Substantivos no Plural**

```go
// ✅ Bom
GET /products
GET /users
GET /orders

// ❌ Ruim
GET /getProducts
GET /user
GET /order_list
```

### **2. Usar Hierarquia de Recursos**

```go
// ✅ Bom
GET /users/:userId/orders
GET /orders/:orderId/items

// ❌ Ruim
GET /getUserOrders/:userId
GET /getOrderItems/:orderId
```

### **3. Filtros via Query Parameters**

```go
// ✅ Bom
GET /products?category=electronics&minPrice=100&maxPrice=500
GET /users?page=1&limit=20&sort=name

// ❌ Ruim
GET /products/electronics/100-500
GET /users/page/1/limit/20
```

### **4. Versionamento de API**

```go
// ✅ Bom
GET /api/v1/products
GET /api/v2/products

// Alternativa no header
GET /products
Header: Accept: application/vnd.api+json; version=1
```

### **5. Status Codes Apropriados**

```go
// Criar recurso
c.JSON(http.StatusCreated, product)  // 201

// Recurso não encontrado
c.JSON(http.StatusNotFound, error)   // 404

// Validação falhou
c.JSON(http.StatusBadRequest, error) // 400

// Conflito (recurso já existe)
c.JSON(http.StatusConflict, error)   // 409
```

### **6. Respostas Consistentes**

```go
// Sucesso
{
  "data": { ... },
  "message": "Product created successfully"
}

// Erro
{
  "error": {
    "code": "PRODUCT_NOT_FOUND",
    "message": "Product with ID 123 not found"
  }
}
```

## ⚡ Performance

### **1. Usar gin.New() em vez de gin.Default()**

Para melhor controle sobre middleware:

```go
r := gin.New()
r.Use(gin.Recovery())  // Apenas recovery, sem logger
```

### **2. Limite de Request Body**

```go
r.MaxMultipartMemory = 8 << 20  // 8 MB
```

### **3. Timeout**

```go
server := &http.Server{
    Addr:           ":8080",
    Handler:        r,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
}
```

## 📚 Recursos

- **Documentação Oficial**: https://gin-gonic.com/docs/
- **GitHub**: https://github.com/gin-gonic/gin
- **REST API Best Practices**: https://restfulapi.net/

---

**Anterior:** [Graceful Shutdown](04-graceful-shutdown.md) | **Próximo:** [Swagger Documentation](06-swagger-documentation.md)

