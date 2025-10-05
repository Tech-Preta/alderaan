# Swagger / OpenAPI Documentation

## 📚 O que é Swagger?

**Swagger** (agora parte do **OpenAPI Initiative**) é um conjunto de ferramentas para:
- **Documentar APIs RESTful** de forma interativa
- **Testar endpoints** diretamente no navegador
- **Gerar código cliente** automaticamente
- **Validar contratos de API**

## 🎯 Por que usar Swagger?

### ✅ Benefícios

1. **Documentação Sempre Atualizada**: Gerada a partir do código
2. **Interface Interativa**: Teste sua API sem ferramentas externas
3. **Contratos Claros**: Especificação OpenAPI padronizada
4. **Facilita Integração**: Desenvolvedores entendem a API rapidamente
5. **Validação Automática**: Garante que requests/responses seguem o contrato

## 🚀 Configuração no Projeto

### **1. Dependências Instaladas**

```bash
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

### **2. Estrutura de Arquivos**

```
docs/
├── docs.go          # Código gerado automaticamente
├── swagger.json     # Especificação em JSON
└── swagger.yaml     # Especificação em YAML
```

**⚠️ IMPORTANTE**: Não edite manualmente os arquivos `docs.go`, `swagger.json` ou `swagger.yaml`. Eles são regenerados pelo comando `swag init`.

## 📝 Anotações Swagger

### **Anotações Gerais (main.go)**

```go
// @title           Servidor HTTP com Domain Driven Design
// @version         1.0
// @description     API RESTful com DDD, Gin e Graceful Shutdown
// @contact.name    API Support
// @contact.url     https://github.com/williamkoller
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /api/v1
```

**Campos disponíveis:**
- `@title` - Nome da API
- `@version` - Versão da API
- `@description` - Descrição geral
- `@contact.name` / `@contact.email` / `@contact.url` - Informações de contato
- `@license.name` / `@license.url` - Licença
- `@host` - Host da API
- `@BasePath` - Caminho base

### **Anotações de Endpoint**

```go
// Create godoc
// @Summary      Criar um novo produto
// @Description  Cria um novo produto com nome, SKU, categorias e preço
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      CreateProductInput  true  "Dados do produto"
// @Success      201      {object}  product_entity.Product
// @Failure      400      {object}  ErrorResponse
// @Failure      409      {object}  ErrorResponse
// @Router       /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
    // Implementação
}
```

**Campos disponíveis:**
- `@Summary` - Resumo curto do endpoint
- `@Description` - Descrição detalhada
- `@Tags` - Categoria/tag para agrupar endpoints
- `@Accept` - Formato de entrada aceito (json, xml, etc)
- `@Produce` - Formato de saída (json, xml, etc)
- `@Param` - Parâmetros do endpoint
- `@Success` - Respostas de sucesso
- `@Failure` - Respostas de erro
- `@Router` - Rota e método HTTP

### **Tipos de Parâmetros**

#### **1. Body Parameter**
```go
// @Param product body CreateProductInput true "Dados do produto"
```

#### **2. Path Parameter**
```go
// @Param name path string true "Nome do produto"
```

#### **3. Query Parameter**
```go
// @Param page query int false "Número da página" default(1)
// @Param limit query int false "Itens por página" default(10)
```

#### **4. Header Parameter**
```go
// @Param Authorization header string true "Bearer token"
```

### **Estruturas de Dados**

```go
// CreateProductInput representa os dados de entrada para criar um produto
type CreateProductInput struct {
    Name       string   `json:"name" binding:"required" example:"Notebook"`
    Sku        int      `json:"sku" binding:"required" example:"12345"`
    Categories []string `json:"categories" binding:"required" example:"Eletrônicos,Computadores"`
    Price      int      `json:"price" binding:"required" example:"3500"`
}
```

**Tags importantes:**
- `json` - Nome do campo no JSON
- `binding` - Validações (required, min, max, etc)
- `example` - Valor de exemplo na documentação

## 🔧 Gerando a Documentação

### **Comando Principal**

```bash
swag init -g cmd/main.go -o docs
```

**Parâmetros:**
- `-g` - Arquivo principal com anotações gerais
- `-o` - Diretório de saída (padrão: docs)

### **Quando Regenerar**

Execute `swag init` sempre que:
- ✅ Adicionar novos endpoints
- ✅ Modificar anotações Swagger
- ✅ Alterar estruturas de dados
- ✅ Mudar informações gerais da API

### **Script Helper**

Adicione no `Makefile`:

```makefile
.PHONY: swagger
swagger:
	@echo "Gerando documentação Swagger..."
	swag init -g cmd/main.go -o docs
	@echo "✅ Documentação gerada em docs/"

.PHONY: run
run: swagger
	go run cmd/main.go
```

Uso:
```bash
make swagger  # Apenas gerar documentação
make run      # Gerar e executar
```

## 🌐 Acessando a Documentação

### **1. Iniciar o Servidor**

```bash
go run cmd/main.go
```

### **2. Abrir no Navegador**

```
http://localhost:8080/swagger/index.html
```

### **3. Interface Swagger UI**

A interface oferece:
- 📋 Lista de todos os endpoints
- 📝 Documentação detalhada de cada endpoint
- 🧪 Botão "Try it out" para testar
- 📊 Schemas dos modelos de dados
- 💾 Download da especificação OpenAPI

## 🧪 Testando via Swagger UI

### **Exemplo: Criar Produto**

1. **Acesse** `http://localhost:8080/swagger/index.html`
2. **Expanda** o endpoint `POST /api/v1/products`
3. **Clique** em "Try it out"
4. **Edite** o JSON de exemplo:
   ```json
   {
     "name": "Notebook Dell",
     "sku": 54321,
     "categories": ["Eletrônicos", "Notebooks"],
     "price": 4500
   }
   ```
5. **Clique** em "Execute"
6. **Veja** a resposta abaixo

### **Exemplo: Listar Produtos**

1. **Expanda** `GET /api/v1/products`
2. **Clique** em "Try it out"
3. **Clique** em "Execute"
4. **Veja** a lista de produtos

## 📦 Estrutura do Projeto com Swagger

```
.
├── cmd/
│   └── main.go                     # Anotações gerais @title, @version, etc
├── docs/
│   ├── docs.go                     # Gerado automaticamente ⚠️
│   ├── swagger.json                # Especificação OpenAPI em JSON ⚠️
│   ├── swagger.yaml                # Especificação OpenAPI em YAML ⚠️
│   └── 06-swagger-documentation.md # Este guia
├── internal/
│   └── infra/
│       └── http/
│           ├── handlers/
│           │   └── product_handler.go  # Anotações de endpoint
│           └── router/
│               └── product_router.go   # Rota /swagger/*any
└── go.mod
```

## 🎨 Personalizando a Documentação

### **Adicionar Autenticação**

```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// No endpoint:
// @Security BearerAuth
func (h *Handler) ProtectedEndpoint(c *gin.Context) {
    // ...
}
```

### **Adicionar Tags Personalizadas**

```go
// @tag.name products
// @tag.description Operações relacionadas a produtos

// @tag.name users
// @tag.description Gerenciamento de usuários
```

### **Respostas com Múltiplos Schemas**

```go
// @Success 200 {object} SuccessResponse
// @Success 200 {object} AlternativeResponse
// @Failure 400 {object} BadRequestError
// @Failure 404 {object} NotFoundError
// @Failure 500 {object} InternalServerError
```

## 🔍 Exemplos Avançados

### **Endpoint Completo**

```go
// ListProducts godoc
// @Summary      Listar produtos com paginação
// @Description  Retorna lista paginada de produtos com filtros opcionais
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        page      query  int     false  "Número da página"      default(1)      minimum(1)
// @Param        limit     query  int     false  "Itens por página"      default(10)     minimum(1)  maximum(100)
// @Param        category  query  string  false  "Filtrar por categoria"
// @Param        minPrice  query  int     false  "Preço mínimo"          minimum(0)
// @Param        maxPrice  query  int     false  "Preço máximo"
// @Success      200       {object}  PaginatedProductResponse
// @Failure      400       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Router       /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
    // Implementação
}
```

### **Upload de Arquivo**

```go
// UploadImage godoc
// @Summary      Upload de imagem do produto
// @Description  Faz upload de uma imagem para o produto
// @Tags         products
// @Accept       multipart/form-data
// @Produce      json
// @Param        id    path    string  true  "ID do produto"
// @Param        file  formData  file  true  "Arquivo de imagem"
// @Success      200   {object}  UploadResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /products/{id}/image [post]
func (h *ProductHandler) UploadImage(c *gin.Context) {
    // Implementação
}
```

## 📊 Boas Práticas

### ✅ **DO (Faça)**

1. **Sempre use anotações completas**
   ```go
   // @Summary, @Description, @Tags, @Accept, @Produce
   ```

2. **Documente todos os códigos de status**
   ```go
   // @Success 200
   // @Failure 400
   // @Failure 404
   // @Failure 500
   ```

3. **Use exemplos realistas**
   ```go
   example:"usuario@email.com"
   ```

4. **Agrupe endpoints por tags**
   ```go
   // @Tags products
   // @Tags users
   // @Tags auth
   ```

5. **Valide com binding tags**
   ```go
   `binding:"required,email"`
   ```

### ❌ **DON'T (Não Faça)**

1. **Não edite arquivos gerados manualmente**
   ```
   ❌ docs/docs.go
   ❌ docs/swagger.json
   ❌ docs/swagger.yaml
   ```

2. **Não use descrições vagas**
   ```go
   // ❌ @Summary Get data
   // ✅ @Summary Listar todos os produtos ativos
   ```

3. **Não esqueça de regenerar após mudanças**
   ```bash
   # Sempre execute após modificar anotações
   swag init -g cmd/main.go -o docs
   ```

## 🐳 Swagger em Produção

### **Desabilitar em Produção (Opcional)**

```go
if os.Getenv("GIN_MODE") != "release" {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

### **Usar URL Externa**

```go
url := ginSwagger.URL("http://api.exemplo.com/swagger/doc.json")
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
```

## 🔗 Exportando a Especificação

### **Compartilhar com Time**

```bash
# JSON
curl http://localhost:8080/swagger/doc.json > api-spec.json

# YAML
cp docs/swagger.yaml api-spec.yaml
```

### **Importar em Outras Ferramentas**

- **Postman**: Import → OpenAPI 3.0
- **Insomnia**: Import → From File
- **API Gateway**: Carregar especificação OpenAPI

## 📚 Recursos

- **Swag Docs**: https://github.com/swaggo/swag
- **Gin-Swagger**: https://github.com/swaggo/gin-swagger
- **OpenAPI Spec**: https://swagger.io/specification/
- **Swagger UI**: https://swagger.io/tools/swagger-ui/

## 🎯 Checklist

- [x] Dependências instaladas (`swag`, `gin-swagger`, `files`)
- [x] Anotações gerais no `main.go`
- [x] Anotações nos handlers
- [x] Rota `/swagger/*any` configurada
- [x] Documentação gerada com `swag init`
- [x] Interface acessível em `http://localhost:8080/swagger/index.html`

---

**Anterior:** [RESTful API com Gin](05-restful-api-gin.md) | **Voltar ao início:** [README](README.md)
