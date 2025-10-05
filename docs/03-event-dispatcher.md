# Event Dispatcher (Despachador de Eventos)

## 📚 O que é Event Dispatcher?

**Event Dispatcher** é um padrão de design que implementa o padrão **Observer** (Observador), permitindo que diferentes partes de um sistema se comuniquem de forma **desacoplada** através de eventos.

## 🎯 Problema que Resolve

### ❌ **Sem Event Dispatcher (Acoplado)**

```go
func CreateProduct(name string, sku int, price int) error {
    product := NewProduct(name, sku, price)

    // Múltiplas responsabilidades no mesmo lugar
    SaveToDatabase(product)
    SendEmailNotification(product)      // ❌ Acoplado
    UpdateSearchIndex(product)          // ❌ Acoplado
    LogAnalytics(product)               // ❌ Acoplado
    NotifyInventorySystem(product)      // ❌ Acoplado

    return nil
}
```

**Problemas:**
- Alto acoplamento
- Difícil de testar
- Difícil adicionar novas funcionalidades
- Viola o Single Responsibility Principle

### ✅ **Com Event Dispatcher (Desacoplado)**

```go
func CreateProduct(name string, sku int, price int) error {
    product := NewProduct(name, sku, price)

    // Apenas dispara o evento
    event := NewProductCreatedEvent(product)
    dispatcher.Dispatch("product.created", event)  // ✅

    return nil
}

// Outros componentes se registram para ouvir o evento
dispatcher.Register("product.created", SendEmailHandler)
dispatcher.Register("product.created", UpdateSearchIndexHandler)
dispatcher.Register("product.created", LogAnalyticsHandler)
```

**Benefícios:**
- Baixo acoplamento
- Fácil adicionar/remover handlers
- Testável
- Segue o Open/Closed Principle

## 🏗️ Componentes

### 1. **Event (Evento)**

Representa algo que aconteceu no sistema.

```go
// internal/domain/product/events/product_created_event.go
type ProductCreatedEvent struct {
    Name       string
    Sku        int
    Categories []string
    Price      int
}

func NewProductCreatedEvent(name string, sku int, categories []string, price int) *ProductCreatedEvent {
    return &ProductCreatedEvent{
        Name:       name,
        Sku:        sku,
        Categories: categories,
        Price:      price,
    }
}

func (e *ProductCreatedEvent) EventName() string {
    return "product.created"
}
```

### 2. **Event Handler (Manipulador de Eventos)**

Função que é executada quando um evento ocorre.

```go
type EventHandler func(event Event)

// Exemplo de handler
func EmailNotificationHandler(event Event) {
    if productEvent, ok := event.(*ProductCreatedEvent); ok {
        fmt.Printf("Sending email for product: %s\n", productEvent.Name)
        // Lógica de envio de email
    }
}
```

### 3. **Event Dispatcher (Despachador)**

Gerencia o registro de handlers e despacho de eventos.

```go
// internal/shared/domain/events/events_handler.go
type EventDispatcher struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex
}

func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{
        handlers: make(map[string][]EventHandler),
    }
}

func (d *EventDispatcher) Register(eventName string, handler EventHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()

    d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *EventDispatcher) Dispatch(eventName string, event Event) {
    d.mu.RLock()
    defer d.mu.RUnlock()

    if handlers, ok := d.handlers[eventName]; ok {
        for _, handler := range handlers {
            go handler(event)  // Executa assincronamente
        }
    }
}
```

## 🔄 Fluxo de Execução

```
1. Ação acontece → Product é criado
2. Evento é criado → ProductCreatedEvent
3. Evento é despachado → dispatcher.Dispatch()
4. Dispatcher encontra handlers registrados
5. Handlers são executados (assincronamente)
   ├── EmailHandler
   ├── SearchIndexHandler
   └── AnalyticsHandler
```

## 💡 Implementação Completa

### **1. Interface de Evento**

```go
type Event interface {
    EventName() string
}
```

### **2. Evento Concreto**

```go
type ProductCreatedEvent struct {
    ProductID  string
    Name       string
    Price      int
    CreatedAt  time.Time
}

func (e *ProductCreatedEvent) EventName() string {
    return "product.created"
}
```

### **3. Registrar Handlers**

```go
func main() {
    dispatcher := NewEventDispatcher()

    // Registrar handlers
    dispatcher.Register("product.created", func(event Event) {
        fmt.Println("Handler 1: Product created!")
    })

    dispatcher.Register("product.created", func(event Event) {
        fmt.Println("Handler 2: Sending email...")
    })

    dispatcher.Register("product.created", func(event Event) {
        fmt.Println("Handler 3: Updating search index...")
    })

    // ... resto do código
}
```

### **4. Disparar Evento**

```go
func NewProduct(name string, sku int, categories []string, price int,
    dispatcher *EventDispatcher) (*Product, *ProductCreatedEvent, error) {

    // Validações
    ok, err := Validate(name, sku, categories, price)
    if !ok {
        return nil, nil, err
    }

    p := &Product{Name: name, Sku: sku, Categories: categories, Price: price}

    // Criar e despachar evento
    event := NewProductCreatedEvent(p.GetName(), p.GetSku(),
        p.GetCategories(), p.GetPrice())

    if dispatcher != nil {
        dispatcher.Dispatch(event.EventName(), event)  // 🚀
    }

    return p, event, nil
}
```

## ⚡ Execução Assíncrona

Os handlers são executados em **goroutines** para não bloquear a operação principal:

```go
func (d *EventDispatcher) Dispatch(eventName string, event Event) {
    d.mu.RLock()
    defer d.mu.RUnlock()

    if handlers, ok := d.handlers[eventName]; ok {
        for _, handler := range handlers {
            go handler(event)  // ⚡ Assíncrono!
        }
    }
}
```

**Vantagens:**
- Não bloqueia a operação principal
- Handlers lentos não afetam a resposta HTTP
- Melhor performance

**Desvantagens:**
- Não há garantia de ordem de execução
- Erros em handlers não são facilmente tratados

## 🛡️ Thread-Safety

O dispatcher usa `sync.RWMutex` para ser **thread-safe**:

```go
type EventDispatcher struct {
    handlers map[string][]EventHandler
    mu       sync.RWMutex  // 🛡️ Proteção contra race conditions
}

func (d *EventDispatcher) Register(eventName string, handler EventHandler) {
    d.mu.Lock()         // Lock exclusivo para escrita
    defer d.mu.Unlock()

    d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *EventDispatcher) Dispatch(eventName string, event Event) {
    d.mu.RLock()        // Lock compartilhado para leitura
    defer d.mu.RUnlock()

    // ...
}
```

## 📝 Exemplos de Uso

### **Exemplo 1: Notificações**

```go
dispatcher.Register("product.created", func(event Event) {
    e := event.(*ProductCreatedEvent)
    SendEmail(e.Name, "Product created successfully!")
})

dispatcher.Register("product.created", func(event Event) {
    e := event.(*ProductCreatedEvent)
    SendSMS("+5511999999999", fmt.Sprintf("New product: %s", e.Name))
})
```

### **Exemplo 2: Analytics**

```go
dispatcher.Register("product.created", func(event Event) {
    e := event.(*ProductCreatedEvent)
    analytics.Track("Product Created", map[string]interface{}{
        "product_name": e.Name,
        "sku":          e.Sku,
        "price":        e.Price,
    })
})
```

### **Exemplo 3: Cache Invalidation**

```go
dispatcher.Register("product.updated", func(event Event) {
    e := event.(*ProductUpdatedEvent)
    cache.Delete(fmt.Sprintf("product:%s", e.ProductID))
})
```

### **Exemplo 4: Search Index**

```go
dispatcher.Register("product.created", func(event Event) {
    e := event.(*ProductCreatedEvent)
    searchIndex.AddDocument(e.ProductID, e.Name, e.Categories)
})
```

## ✅ Benefícios

1. **Desacoplamento**: Componentes não conhecem uns aos outros
2. **Extensibilidade**: Fácil adicionar novos handlers sem modificar código existente
3. **Testabilidade**: Handlers podem ser testados isoladamente
4. **Manutenibilidade**: Cada handler tem uma responsabilidade única
5. **Performance**: Execução assíncrona não bloqueia operação principal

## 🚫 Quando NÃO Usar

- Quando você precisa de garantias transacionais (use sagas ou outbox pattern)
- Quando a ordem de execução é crítica
- Quando você precisa de resposta síncrona de todos os handlers
- Para operações muito simples onde o overhead não compensa

## 🔍 Padrões Relacionados

### **Observer Pattern**
Event Dispatcher é uma implementação do padrão Observer.

### **Pub/Sub (Publish/Subscribe)**
Similar, mas geralmente usado para comunicação entre sistemas distribuídos (ex: RabbitMQ, Kafka).

### **Event Sourcing**
Padrão que armazena eventos como fonte da verdade em vez do estado atual.

## 🎯 Boas Práticas

### 1. **Nomeação Clara de Eventos**
```go
// ✅ Bom
"product.created"
"user.registered"
"order.placed"

// ❌ Ruim
"event1"
"stuff_happened"
```

### 2. **Eventos Imutáveis**
```go
type ProductCreatedEvent struct {
    Name       string    // ✅ Campos são lidos, nunca modificados
    Sku        int
    CreatedAt  time.Time
}
```

### 3. **Handlers Independentes**
Cada handler deve funcionar independentemente dos outros.

### 4. **Tratamento de Erros**
```go
dispatcher.Register("product.created", func(event Event) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Handler panic: %v", r)
        }
    }()

    // Lógica do handler
})
```

## 📚 Evolução Futura

Para sistemas mais complexos, considere:

1. **Event Store**: Persistir eventos para auditoria
2. **Event Replay**: Reproduzir eventos para testes ou recuperação
3. **Message Broker**: Usar RabbitMQ, Kafka, etc para eventos distribuídos
4. **CQRS**: Command Query Responsibility Segregation

---

**Anterior:** [Arquitetura Limpa](02-clean-architecture.md) | **Próximo:** [Graceful Shutdown](04-graceful-shutdown.md)
