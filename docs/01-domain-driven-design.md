# Domain-Driven Design (DDD)

## 📚 O que é DDD?

**Domain-Driven Design (DDD)** é uma abordagem de desenvolvimento de software proposta por Eric Evans em 2003, que coloca o **domínio do negócio** no centro da aplicação.

O DDD não é sobre tecnologia, mas sobre **entender profundamente o problema** que você está resolvendo e **modelar o software de acordo com a realidade do negócio**.

## 🎯 Conceitos Fundamentais

### 1. **Domínio (Domain)**

O domínio é a área de conhecimento ou atividade para a qual o software está sendo construído.

**Exemplo:** Em um e-commerce, o domínio inclui conceitos como:
- Produtos
- Pedidos
- Carrinho de compras
- Pagamentos
- Clientes

### 2. **Entidades (Entities)**

Entidades são objetos que têm **identidade única** e que podem mudar ao longo do tempo.

**Características:**
- Possuem identidade única (ID)
- Podem ser mutáveis
- São comparadas pelo ID, não pelos atributos

```go
type Product struct {
    ID         string    // Identidade única
    Name       string
    Sku        int
    Categories []string
    Price      int
}
```

### 3. **Value Objects (Objetos de Valor)**

São objetos **imutáveis** que não têm identidade própria, sendo definidos apenas por seus atributos.

**Características:**
- Imutáveis
- Comparados pelos valores de seus atributos
- Podem ser compartilhados entre entidades

```go
type Address struct {
    Street     string
    Number     string
    City       string
    PostalCode string
}
```

### 4. **Agregados (Aggregates)**

Um agregado é um grupo de objetos relacionados que são tratados como uma **unidade única** para alterações de dados.

**Exemplo:**
```
Pedido (Aggregate Root)
  ├── Itens do Pedido
  ├── Endereço de Entrega
  └── Informações de Pagamento
```

### 5. **Repositórios (Repositories)**

Repositórios abstraem a lógica de acesso aos dados, permitindo que o domínio não se preocupe com detalhes de persistência.

```go
type IProductRepository interface {
    Add(product Product) error
    Find() ([]Product, error)
    FindOne(name string) (Product, error)
}
```

### 6. **Serviços de Domínio (Domain Services)**

Operações que não pertencem naturalmente a nenhuma entidade ou value object.

```go
type PriceCalculatorService struct {}

func (s *PriceCalculatorService) CalculateDiscount(
    product Product,
    coupon Coupon,
) float64 {
    // Lógica de cálculo de desconto
}
```

### 7. **Eventos de Domínio (Domain Events)**

Representam algo que aconteceu no domínio e que outros componentes podem estar interessados.

```go
type ProductCreatedEvent struct {
    Name       string
    Sku        int
    Categories []string
    Price      int
}
```

## 🏗️ Camadas no DDD

### **Domain Layer (Camada de Domínio)**
- Entidades
- Value Objects
- Agregados
- Eventos de Domínio
- Interfaces de Repositórios

### **Application Layer (Camada de Aplicação)**
- Casos de uso
- Serviços de aplicação
- Comandos e Queries

### **Infrastructure Layer (Camada de Infraestrutura)**
- Implementação de repositórios
- Adaptadores externos (APIs, banco de dados)
- Frameworks e bibliotecas

### **Presentation Layer (Camada de Apresentação)**
- Controllers HTTP
- Serialização/Deserialização
- Validação de entrada

## ✅ Benefícios do DDD

1. **Código alinhado com o negócio**: O código reflete a linguagem e as regras do domínio
2. **Manutenibilidade**: Facilita mudanças e evoluções
3. **Testabilidade**: Domínio isolado é mais fácil de testar
4. **Comunicação**: Linguagem ubíqua facilita comunicação entre desenvolvedores e especialistas de domínio
5. **Escalabilidade**: Código bem organizado escala melhor

## 🚫 Quando NÃO usar DDD

- Projetos muito simples (CRUD básico)
- Aplicações com domínio trivial
- Protótipos rápidos
- Quando a equipe não tem experiência com DDD
- Quando não há acesso a especialistas do domínio

## 💡 Exemplo Prático no Projeto

No nosso projeto, aplicamos DDD da seguinte forma:

```
internal/domain/product/
├── entity/
│   └── product.go          # Entidade Product
├── events/
│   └── product_created_event.go  # Evento de Domínio
└── repository/
    └── product_repository.go     # Interface do Repositório
```

### Entidade com Validação de Negócio

```go
func NewProduct(name string, sku int, categories []string, price int,
    dispatcher *EventDispatcher) (*Product, *ProductCreatedEvent, error) {

    ok, err := Validate(name, sku, categories, price)
    if !ok {
        return nil, nil, err
    }

    p := &Product{Name: name, Sku: sku, Categories: categories, Price: price}

    event := NewProductCreatedEvent(p.GetName(), p.GetSku(),
        p.GetCategories(), p.GetPrice())

    if dispatcher != nil {
        dispatcher.Dispatch(event.EventName(), event)
    }

    return p, event, nil
}
```

### Validações de Domínio

```go
func Validate(name string, sku int, categories []string, price int) (bool, error) {
    if name == "" {
        return false, errors.New("name is required")
    }
    if sku <= 0 {
        return false, errors.New("sku is required")
    }
    if len(categories) == 0 {
        return false, errors.New("categories is required")
    }
    if price <= 0 {
        return false, errors.New("price is required")
    }
    return true, nil
}
```

## 📚 Recursos para Aprofundamento

- **Livro:** "Domain-Driven Design" - Eric Evans
- **Livro:** "Implementing Domain-Driven Design" - Vaughn Vernon
- **Livro:** "Domain-Driven Design Distilled" - Vaughn Vernon

## 🎓 Princípios Chave

1. **Linguagem Ubíqua**: Use a mesma linguagem do negócio no código
2. **Bounded Contexts**: Separe domínios complexos em contextos menores
3. **Foco no Core Domain**: Invista mais esforço no que diferencia seu negócio
4. **Iteração Contínua**: O modelo de domínio evolui com o entendimento do negócio

---

**Próximo:** [Arquitetura Limpa](02-clean-architecture.md)
