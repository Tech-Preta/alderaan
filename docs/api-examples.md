# 🚀 Exemplos de Uso da API

Este arquivo contém exemplos práticos para testar a API de produtos.

## 📍 Base URL

```
http://localhost:8080
```

## 📚 Documentação Interativa

Acesse o Swagger UI para testar interativamente:

```
http://localhost:8080/swagger/index.html
```

---

## ✅ Health Check

Verifica se o servidor está funcionando:

```bash
curl http://localhost:8080/health
```

**Resposta:**
```json
{
  "status": "ok"
}
```

---

## 📦 Criar Produto

### Exemplo 1: Notebook

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notebook Dell Inspiron",
    "sku": 12345,
    "categories": ["Eletrônicos", "Computadores"],
    "price": 3500
  }'
```

### Exemplo 2: Smartphone

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro",
    "sku": 67890,
    "categories": ["Eletrônicos", "Smartphones"],
    "price": 7500
  }'
```

### Exemplo 3: Livro

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Domain-Driven Design",
    "sku": 11111,
    "categories": ["Livros", "Tecnologia"],
    "price": 120
  }'
```

**Resposta de Sucesso (201 Created):**
```json
{
  "name": "Notebook Dell Inspiron",
  "sku": 12345,
  "categories": ["Eletrônicos", "Computadores"],
  "price": 3500
}
```

**Resposta de Erro (400 Bad Request):**
```json
{
  "error": "name is required"
}
```

**Resposta de Conflito (409 Conflict):**
```json
{
  "error": "product already exists"
}
```

---

## 📋 Listar Todos os Produtos

```bash
curl http://localhost:8080/api/v1/products
```

**Resposta (200 OK):**
```json
[
  {
    "name": "Notebook Dell Inspiron",
    "sku": 12345,
    "categories": ["Eletrônicos", "Computadores"],
    "price": 3500
  },
  {
    "name": "iPhone 15 Pro",
    "sku": 67890,
    "categories": ["Eletrônicos", "Smartphones"],
    "price": 7500
  }
]
```

---

## 🔍 Buscar Produto por Nome

```bash
curl http://localhost:8080/api/v1/products/iPhone%2015%20Pro
```

**Nota:** Espaços devem ser codificados como `%20` na URL.

**Resposta (200 OK):**
```json
{
  "name": "iPhone 15 Pro",
  "sku": 67890,
  "categories": ["Eletrônicos", "Smartphones"],
  "price": 7500
}
```

**Resposta de Erro (404 Not Found):**
```json
{
  "error": "product not found"
}
```

---

## 🧪 Testando Validações

### ❌ Produto sem nome

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "sku": 12345,
    "categories": ["Eletrônicos"],
    "price": 3500
  }'
```

**Resposta (400 Bad Request):**
```json
{
  "error": "Key: 'CreateProductInput.Name' Error:Field validation for 'Name' failed on the 'required' tag"
}
```

### ❌ Produto com SKU inválido

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Teste",
    "sku": 0,
    "categories": ["Teste"],
    "price": 100
  }'
```

**Resposta (400 Bad Request):**
```json
{
  "error": "sku is required"
}
```

### ❌ Produto sem categorias

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Teste",
    "sku": 12345,
    "categories": [],
    "price": 100
  }'
```

**Resposta (400 Bad Request):**
```json
{
  "error": "categories is required"
}
```

### ❌ Produto com preço inválido

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Teste",
    "sku": 12345,
    "categories": ["Teste"],
    "price": -100
  }'
```

**Resposta (400 Bad Request):**
```json
{
  "error": "price is required"
}
```

---

## 🎯 Testando Produto Duplicado

### 1. Criar um produto

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Teste Duplicado",
    "sku": 99999,
    "categories": ["Teste"],
    "price": 100
  }'
```

### 2. Tentar criar novamente (mesmo nome)

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Teste Duplicado",
    "sku": 88888,
    "categories": ["Teste"],
    "price": 200
  }'
```

**Resposta (409 Conflict):**
```json
{
  "error": "product already exists"
}
```

---

## 📊 Testando com jq (Pretty Print)

Se você tem `jq` instalado, pode formatar as respostas JSON:

```bash
# Listar produtos formatados
curl -s http://localhost:8080/api/v1/products | jq '.'

# Buscar produto formatado
curl -s http://localhost:8080/api/v1/products/iPhone%2015%20Pro | jq '.'

# Criar e formatar resposta
curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mouse Gamer",
    "sku": 55555,
    "categories": ["Periféricos"],
    "price": 150
  }' | jq '.'
```

---

## 🔧 Usando HTTPie (Alternativa ao curl)

Se preferir HTTPie para uma sintaxe mais amigável:

```bash
# Health check
http GET localhost:8080/health

# Criar produto
http POST localhost:8080/api/v1/products \
  name="Teclado Mecânico" \
  sku:=44444 \
  categories:='["Periféricos", "Teclados"]' \
  price:=350

# Listar produtos
http GET localhost:8080/api/v1/products

# Buscar produto
http GET localhost:8080/api/v1/products/Teclado%20Mecânico
```

---

## 🐳 Testando com Docker

Se você executar a API em um container Docker, use:

```bash
# Criar produto
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto via Docker",
    "sku": 77777,
    "categories": ["Docker"],
    "price": 999
  }'
```

---

## 📱 Testando com Postman

1. Importe a especificação OpenAPI:
   - Abra o Postman
   - Import → Link → `http://localhost:8080/swagger/doc.json`
   - Todos os endpoints serão importados automaticamente!

2. Ou use a collection manual:
   - **POST** `http://localhost:8080/api/v1/products`
   - Headers: `Content-Type: application/json`
   - Body (raw JSON):
     ```json
     {
       "name": "Teste Postman",
       "sku": 33333,
       "categories": ["Teste"],
       "price": 250
     }
     ```

---

## 🎭 Scripts de Teste Automatizado

### Script Bash para Testes

Crie um arquivo `test-api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "🧪 Testando API..."

# Health check
echo "\n1. Health Check"
curl -s $BASE_URL/health | jq '.'

# Criar produtos
echo "\n2. Criando produtos..."
curl -s -X POST $BASE_URL/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Teste 1",
    "sku": 11111,
    "categories": ["Teste"],
    "price": 100
  }' | jq '.'

curl -s -X POST $BASE_URL/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Produto Teste 2",
    "sku": 22222,
    "categories": ["Teste"],
    "price": 200
  }' | jq '.'

# Listar produtos
echo "\n3. Listando todos os produtos..."
curl -s $BASE_URL/api/v1/products | jq '.'

# Buscar produto específico
echo "\n4. Buscando produto específico..."
curl -s $BASE_URL/api/v1/products/Produto%20Teste%201 | jq '.'

echo "\n✅ Testes concluídos!"
```

Execute:
```bash
chmod +x test-api.sh
./test-api.sh
```

---

## 📈 Monitoramento e Logs

Observe os logs do servidor enquanto faz requisições:

```bash
# Terminal 1: Executar servidor
make run

# Terminal 2: Fazer requisições
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Monitor 4K", "sku": 66666, "categories": ["Monitores"], "price": 2500}'
```

Você verá logs como:
```
[GIN] 2024/10/04 - 21:53:15 | 201 |    1.234567ms |       127.0.0.1 | POST     "/api/v1/products"
```

---

## 🎓 Recursos Adicionais

- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **OpenAPI Spec (JSON)**: `http://localhost:8080/swagger/doc.json`
- **OpenAPI Spec (YAML)**: Ver arquivo `docs/swagger.yaml`
- **Documentação Técnica**: Ver pasta `docs/`

---

**💡 Dica:** Use o Swagger UI para testar interativamente - é muito mais fácil que curl! 🚀
