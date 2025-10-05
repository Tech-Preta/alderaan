# Graceful Shutdown (Desligamento Controlado)

## 📚 O que é Graceful Shutdown?

**Graceful Shutdown** é o processo de encerrar uma aplicação de forma **controlada e segura**, garantindo que:

- Requisições em andamento sejam finalizadas
- Conexões sejam fechadas adequadamente
- Recursos sejam liberados corretamente
- Dados não sejam corrompidos
- Nenhuma operação seja interrompida abruptamente

## 🎯 Problema que Resolve

### ❌ **Sem Graceful Shutdown**

```go
func main() {
    http.ListenAndServe(":8080", handler)
    // Pressionar Ctrl+C mata o processo imediatamente
}
```

**Consequências:**
- ❌ Requisições HTTP em andamento são abortadas
- ❌ Conexões de banco de dados não são fechadas
- ❌ Arquivos podem ficar corrompidos
- ❌ Transações podem ficar incompletas
- ❌ Usuários recebem erros de conexão

### ✅ **Com Graceful Shutdown**

```go
func main() {
    server := &http.Server{Addr: ":8080", Handler: handler}
    
    go server.ListenAndServe()
    
    GracefulShutdown(server, 5*time.Second)
}
```

**Benefícios:**
- ✅ Requisições são finalizadas antes de encerrar
- ✅ Conexões são fechadas adequadamente
- ✅ Recursos são liberados
- ✅ Dados permanecem consistentes
- ✅ Zero downtime com deploys rolling

## 🏗️ Componentes

### 1. **Signal Handling (Captura de Sinais)**

O sistema operacional envia sinais para processos. Os principais são:

- **SIGINT** (Ctrl+C): Interrupção do teclado
- **SIGTERM**: Pedido de término (usado por orquestradores como Kubernetes)
- **SIGKILL**: Término forçado (não pode ser capturado)

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

<-quit  // Bloqueia até receber um sinal
fmt.Println("Shutdown signal received!")
```

### 2. **Context com Timeout**

Dá um tempo limite para finalizar operações:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

if err := server.Shutdown(ctx); err != nil {
    fmt.Printf("Forced shutdown: %v\n", err)
}
```

### 3. **Server Shutdown**

O servidor HTTP do Go tem método `Shutdown` que:
- Para de aceitar novas conexões
- Aguarda requisições ativas finalizarem
- Fecha o servidor

```go
if err := server.Shutdown(ctx); err != nil {
    // Timeout ou erro no shutdown
}
```

## 💡 Implementação Completa

```go
// cmd/main.go
package main

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    // Configurar servidor
    r := gin.Default()
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    server := &http.Server{
        Addr:    ":8080",
        Handler: r,
    }

    // Iniciar servidor em goroutine
    go func() {
        fmt.Println("Server running on http://localhost:8080")
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            fmt.Printf("Failed to start server: %v\n", err)
        }
    }()

    // Aguardar sinal de shutdown
    GracefulShutdown(server, 5*time.Second)
}

func GracefulShutdown(server *http.Server, timeout time.Duration) {
    // 1. Criar canal para receber sinais
    quit := make(chan os.Signal, 1)
    
    // 2. Registrar quais sinais queremos capturar
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // 3. Bloquear até receber um sinal
    <-quit
    fmt.Println("\nReceived termination signal. Shutting down server...")

    // 4. Criar context com timeout
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    // 5. Tentar desligar gracefully
    if err := server.Shutdown(ctx); err != nil {
        fmt.Printf("Error during server shutdown: %v\n", err)
    } else {
        fmt.Println("Server shut down gracefully.")
    }
}
```

## 🔄 Fluxo de Execução

```
1. Aplicação inicia
2. Servidor HTTP começa a aceitar conexões
3. Goroutine aguarda sinal de término
   │
   ├── Usuário pressiona Ctrl+C
   │   └── Sistema envia SIGINT
   │
4. Sinal é capturado
5. Servidor para de aceitar novas conexões
6. Aguarda requisições ativas finalizarem (até timeout)
   │
   ├── Cenário 1: Todas finalizaram antes do timeout ✅
   │   └── Shutdown graceful completo
   │
   └── Cenário 2: Timeout expirou com requisições ativas ⚠️
       └── Shutdown forçado
7. Recursos são liberados
8. Aplicação encerra
```

## ⏱️ Timeout

O timeout define quanto tempo esperar:

```go
// Timeout curto (5 segundos)
GracefulShutdown(server, 5*time.Second)

// Timeout longo (30 segundos) - para operações demoradas
GracefulShutdown(server, 30*time.Second)

// Sem timeout (aguarda indefinidamente)
ctx := context.Background()
server.Shutdown(ctx)
```

**Recomendação:** 
- **APIs rápidas**: 5-10 segundos
- **APIs com operações longas**: 30-60 segundos
- **Workers de background**: 2-5 minutos

## 🧪 Testando Graceful Shutdown

### **Teste Manual**

```bash
# Terminal 1: Iniciar servidor
go run cmd/main.go

# Terminal 2: Fazer requisição longa
curl http://localhost:8080/slow-endpoint

# Terminal 1: Pressionar Ctrl+C enquanto requisição está em andamento
# Observe que a requisição finaliza antes do servidor encerrar
```

### **Teste com Script**

```bash
#!/bin/bash

# Iniciar servidor em background
go run cmd/main.go &
SERVER_PID=$!

sleep 2

# Fazer requisição longa em background
curl http://localhost:8080/slow-endpoint &

sleep 1

# Enviar SIGTERM
kill -TERM $SERVER_PID

# Aguardar servidor encerrar
wait $SERVER_PID

echo "Graceful shutdown completed!"
```

## 🔧 Cenários Avançados

### **1. Múltiplos Recursos**

```go
func GracefulShutdown(server *http.Server, db *sql.DB, cache *redis.Client) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    fmt.Println("Shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Encerrar servidor HTTP
    if err := server.Shutdown(ctx); err != nil {
        log.Printf("HTTP server shutdown error: %v", err)
    }

    // Fechar conexão do banco de dados
    if err := db.Close(); err != nil {
        log.Printf("Database close error: %v", err)
    }

    // Fechar conexão do cache
    if err := cache.Close(); err != nil {
        log.Printf("Cache close error: %v", err)
    }

    fmt.Println("All resources closed successfully")
}
```

### **2. Workers em Background**

```go
type Worker struct {
    quit chan struct{}
}

func (w *Worker) Start() {
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                // Trabalho periódico
                fmt.Println("Working...")
            case <-w.quit:
                fmt.Println("Worker stopped")
                return
            }
        }
    }()
}

func (w *Worker) Stop() {
    close(w.quit)
}

func main() {
    worker := &Worker{quit: make(chan struct{})}
    worker.Start()

    server := &http.Server{Addr: ":8080"}
    go server.ListenAndServe()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // Shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    server.Shutdown(ctx)
    worker.Stop()  // Parar worker também
}
```

### **3. Graceful Shutdown com WaitGroup**

```go
func main() {
    var wg sync.WaitGroup

    // Iniciar múltiplos workers
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            // Trabalho do worker
        }(i)
    }

    // Aguardar sinal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    fmt.Println("Waiting for workers to finish...")
    wg.Wait()
    fmt.Println("All workers finished!")
}
```

## 🐳 Graceful Shutdown no Docker

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o server cmd/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

# Importante: usar exec para propagar sinais corretamente
CMD ["./server"]

# Ou com shell script:
# CMD exec ./server
```

**docker-compose.yml:**
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    stop_grace_period: 30s  # Tempo antes de SIGKILL
```

## ☸️ Graceful Shutdown no Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-server
spec:
  template:
    spec:
      containers:
      - name: api
        image: my-api:latest
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sh", "-c", "sleep 5"]
        # Tempo para graceful shutdown antes de SIGKILL
        terminationGracePeriodSeconds: 30
```

## ✅ Benefícios

1. **Sem Perda de Requisições**: Usuários não veem erros durante deploys
2. **Zero Downtime**: Com rolling updates no Kubernetes
3. **Integridade de Dados**: Transações são finalizadas
4. **Melhor UX**: Usuários não percebem o deploy
5. **Profissionalismo**: Sistemas de produção devem ter shutdown controlado

## 🚫 Erros Comuns

### ❌ **1. Não usar goroutine para ListenAndServe**
```go
// ERRADO
server.ListenAndServe()  // Bloqueia aqui
GracefulShutdown(server) // Nunca executado!

// CORRETO
go server.ListenAndServe()
GracefulShutdown(server)
```

### ❌ **2. Timeout muito curto**
```go
// Pode interromper requisições legítimas
GracefulShutdown(server, 1*time.Second)  // ❌ Muito curto!
```

### ❌ **3. Não capturar SIGTERM**
```go
// ERRADO - Kubernetes usa SIGTERM
signal.Notify(quit, syscall.SIGINT)  // ❌ Só SIGINT

// CORRETO
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)  // ✅
```

### ❌ **4. Não verificar erro de Shutdown**
```go
// ERRADO
server.Shutdown(ctx)  // ❌ Ignora erro

// CORRETO
if err := server.Shutdown(ctx); err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

## 📊 Monitoramento

```go
func GracefulShutdown(server *http.Server, timeout time.Duration) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    <-quit
    startTime := time.Now()
    log.Println("Shutdown initiated")

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Printf("Shutdown error after %v: %v", time.Since(startTime), err)
        // Enviar métrica de erro
        metrics.IncrementCounter("shutdown.errors")
    } else {
        log.Printf("Shutdown completed in %v", time.Since(startTime))
        // Enviar métrica de sucesso
        metrics.RecordDuration("shutdown.duration", time.Since(startTime))
    }
}
```

## 📚 Recursos

- **Documentação Go**: https://pkg.go.dev/net/http#Server.Shutdown
- **Kubernetes Lifecycle**: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/

---

**Anterior:** [Event Dispatcher](03-event-dispatcher.md) | **Próximo:** [RESTful API com Gin](05-restful-api-gin.md)

