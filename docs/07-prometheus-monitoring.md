# Monitoramento com Prometheus - Golden Signals

## 📚 O que é Prometheus?

**Prometheus** é um sistema de monitoramento e alerta open-source criado pela SoundCloud, agora parte da Cloud Native Computing Foundation (CNCF).

**Características principais:**
- 📊 Modelo de dados multi-dimensional
- 🔍 Consultas poderosas com PromQL
- 🎯 Pull-based (Prometheus faz scraping)
- 📈 Armazenamento time-series eficiente
- 🚨 Sistema de alertas integrado

## 🎯 Golden Signals (Google SRE)

Os **Golden Signals** são as 4 métricas fundamentais propostas pelo Google SRE (Site Reliability Engineering) para monitorar qualquer sistema:

### **1. Latency (Latência)**
> Quanto tempo leva para processar uma requisição?

**Métricas implementadas:**
```go
http_request_duration_seconds{method="POST", endpoint="/api/v1/products", status="201"}
```

**Por que é importante:**
- Usuários percebem latência alta como "sistema lento"
- Indica problemas de performance
- Pode sinalizar sobrecarga ou bugs

**O que observar:**
- Percentil 50 (mediana): experiência típica
- Percentil 95: experiência da maioria
- Percentil 99: experiência dos casos extremos
- Percentil 99.9: detectar outliers

**Queries PromQL úteis:**
```promql
# Latência média por endpoint (últimos 5 minutos)
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# Percentil 95 de latência
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Percentil 99 de latência
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### **2. Traffic (Tráfego)**
> Quanta demanda o sistema está recebendo?

**Métricas implementadas:**
```go
http_requests_total{method="GET", endpoint="/api/v1/products", status="200"}
```

**Por que é importante:**
- Indica popularidade do serviço
- Ajuda no planejamento de capacidade
- Detecta picos ou quedas anormais

**O que observar:**
- Taxa de requisições por segundo (RPS)
- Distribuição por endpoint
- Padrões de uso ao longo do dia

**Queries PromQL úteis:**
```promql
# Requisições por segundo (últimos 5 minutos)
rate(http_requests_total[5m])

# Requisições por segundo por endpoint
sum(rate(http_requests_total[5m])) by (endpoint)

# Total de requisições nas últimas 24 horas
increase(http_requests_total[24h])
```

### **3. Errors (Erros)**
> Qual é a taxa de falhas das requisições?

**Métricas implementadas:**
```go
http_request_errors_total{method="POST", endpoint="/api/v1/products", status="400", error_type="client_error"}
```

**Por que é importante:**
- Indica problemas na aplicação
- Afeta diretamente a experiência do usuário
- SLA/SLO geralmente baseados em taxa de erro

**O que observar:**
- Taxa de erro (error rate)
- Erros 4xx (cliente) vs 5xx (servidor)
- Endpoints com mais erros

**Queries PromQL úteis:**
```promql
# Taxa de erro geral (%)
(sum(rate(http_request_errors_total[5m])) / sum(rate(http_requests_total[5m]))) * 100

# Taxa de erro por endpoint
sum(rate(http_request_errors_total[5m])) by (endpoint) / sum(rate(http_requests_total[5m])) by (endpoint)

# Erros 5xx (servidor)
sum(rate(http_request_errors_total{status=~"5.."}[5m]))
```

### **4. Saturation (Saturação)**
> Quão "cheio" está o sistema?

**Métricas implementadas:**
```go
http_in_flight_requests
```

**Por que é importante:**
- Indica proximidade do limite de capacidade
- Previne sobrecarga antes que aconteça
- Ajuda no dimensionamento de recursos

**O que observar:**
- Requisições simultâneas
- Uso de memória
- Uso de CPU
- Conexões do banco de dados

**Queries PromQL úteis:**
```promql
# Requisições em andamento (média últimos 5 min)
avg(http_in_flight_requests)

# Máximo de requisições simultâneas
max_over_time(http_in_flight_requests[5m])

# Saturação em % (assumindo limite de 1000 requisições)
(http_in_flight_requests / 1000) * 100
```

## 💼 Métricas de Negócio

Além dos Golden Signals, implementamos métricas específicas do domínio:

### **Produtos Criados**
```go
products_created_total
```
Contador que nunca diminui. Mostra total de produtos criados desde o início.

**Queries úteis:**
```promql
# Taxa de criação de produtos por minuto
rate(products_created_total[5m]) * 60

# Total criado nas últimas 24h
increase(products_created_total[24h])
```

### **Total de Produtos**
```go
products_total
```
Gauge que mostra quantidade atual de produtos.

**Queries úteis:**
```promql
# Número atual de produtos
products_total

# Crescimento nas últimas 24h
products_total - products_total offset 24h
```

### **Produtos por Categoria**
```go
products_by_category{category="Eletrônicos"}
```
Distribuição de produtos por categoria.

**Queries úteis:**
```promql
# Top 5 categorias
topk(5, products_by_category)

# Produtos por categoria (ordenado)
sort_desc(products_by_category)
```

### **Valor Total dos Produtos**
```go
products_total_value
```
Valor monetário total do inventário.

**Queries úteis:**
```promql
# Valor total atual
products_total_value

# Crescimento de valor (últimas 24h)
rate(products_total_value[24h])
```

### **Preço Médio**
```go
products_average_price
```
Preço médio dos produtos.

**Queries úteis:**
```promql
# Preço médio atual
products_average_price

# Variação do preço médio
deriv(products_average_price[1h])
```

## 📊 Estrutura das Métricas

### **Tipos de Métricas no Prometheus**

#### **1. Counter (Contador)**
Valor que só aumenta (ou reseta para zero).

**Exemplo:**
```go
prometheus.NewCounter(prometheus.CounterOpts{
    Name: "products_created_total",
    Help: "Total de produtos criados",
})
```

**Quando usar:**
- Requisições processadas
- Erros ocorridos
- Produtos criados

#### **2. Gauge (Medidor)**
Valor que pode subir ou descer.

**Exemplo:**
```go
prometheus.NewGauge(prometheus.GaugeOpts{
    Name: "products_total",
    Help: "Número atual de produtos",
})
```

**Quando usar:**
- Temperatura
- Memória usada
- Requisições em andamento

#### **3. Histogram (Histograma)**
Observa valores e os agrupa em buckets.

**Exemplo:**
```go
prometheus.NewHistogram(prometheus.HistogramOpts{
    Name: "http_request_duration_seconds",
    Help: "Duração das requisições",
    Buckets: []float64{0.001, 0.01, 0.1, 1.0},
})
```

**Quando usar:**
- Latência de requisições
- Tamanho de requisições
- Tempo de processamento

#### **4. Summary (Resumo)**
Similar ao histogram, mas calcula percentis no cliente.

**Quando usar:**
- Quando você precisa de percentis exatos
- Menor overhead no Prometheus

## 🚀 Acessando as Métricas

### **Endpoint /metrics**

```bash
curl http://localhost:8080/metrics
```

**Formato Prometheus:**
```
# HELP http_requests_total Total de requisições HTTP recebidas (Traffic)
# TYPE http_requests_total counter
http_requests_total{endpoint="/api/v1/products",method="POST",status="201"} 10

# HELP http_request_duration_seconds Duração das requisições HTTP em segundos (Latency)
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{endpoint="/api/v1/products",method="POST",status="201",le="0.001"} 5
http_request_duration_seconds_bucket{endpoint="/api/v1/products",method="POST",status="201",le="0.01"} 9
http_request_duration_seconds_bucket{endpoint="/api/v1/products",method="POST",status="201",le="+Inf"} 10
http_request_duration_seconds_sum{endpoint="/api/v1/products",method="POST",status="201"} 0.045
http_request_duration_seconds_count{endpoint="/api/v1/products",method="POST",status="201"} 10

# HELP products_created_total Total de produtos criados
# TYPE products_created_total counter
products_created_total 10

# HELP products_total Número atual de produtos cadastrados
# TYPE products_total gauge
products_total 10
```

## 🛠️ Configurando o Prometheus

### **1. Arquivo de Configuração**

Crie `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'alderaan-api'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
```

### **2. Executar Prometheus**

```bash
# Via Docker
docker run -p 9090:9090 \
  -v $(pwd)/prometheus.yml:/etc/prometheus/prometheus.yml \
  prom/prometheus

# Via binário
./prometheus --config.file=prometheus.yml
```

### **3. Acessar Interface**

```
http://localhost:9090
```

## 📈 Configurando Grafana

### **1. Executar Grafana**

```bash
docker run -d -p 3000:3000 grafana/grafana
```

### **2. Adicionar Data Source**

1. Acesse `http://localhost:3000` (admin/admin)
2. Configuration → Data Sources → Add data source
3. Escolha Prometheus
4. URL: `http://localhost:9090`
5. Save & Test

### **3. Importar Dashboard**

O projeto inclui dashboards prontos em `monitoring/dashboards/`.

## 🎨 Queries PromQL Essenciais

### **Dashboard de Golden Signals**

```promql
# Latência (p95) por endpoint
histogram_quantile(0.95, 
  sum(rate(http_request_duration_seconds_bucket[5m])) by (endpoint, le)
)

# Taxa de requisições
sum(rate(http_requests_total[5m])) by (endpoint)

# Taxa de erro (%)
(
  sum(rate(http_request_errors_total[5m])) 
  / 
  sum(rate(http_requests_total[5m]))
) * 100

# Requisições simultâneas
http_in_flight_requests
```

### **Dashboard de Negócio**

```promql
# Produtos criados (taxa por hora)
rate(products_created_total[1h]) * 3600

# Total de produtos
products_total

# Top 10 categorias
topk(10, products_by_category)

# Valor total do inventário
products_total_value

# Preço médio
products_average_price
```

## 🚨 Alertas Recomendados

### **Latência Alta**

```yaml
- alert: HighLatency
  expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 1
  for: 5m
  annotations:
    summary: "Latência alta detectada"
    description: "P99 acima de 1 segundo"
```

### **Taxa de Erro Alta**

```yaml
- alert: HighErrorRate
  expr: (sum(rate(http_request_errors_total[5m])) / sum(rate(http_requests_total[5m]))) > 0.05
  for: 5m
  annotations:
    summary: "Taxa de erro alta"
    description: "Mais de 5% de erros"
```

### **Saturação Alta**

```yaml
- alert: HighSaturation
  expr: http_in_flight_requests > 100
  for: 5m
  annotations:
    summary: "Muitas requisições simultâneas"
    description: "Mais de 100 requisições em andamento"
```

## 📊 Boas Práticas

### ✅ **DO (Faça)**

1. **Use labels consistentes**
   ```go
   {method="POST", endpoint="/api/v1/products"}
   ```

2. **Evite cardinalidade alta**
   ```go
   // ❌ Ruim: user_id tem milhões de valores
   http_requests{user_id="12345"}
   
   // ✅ Bom: status tem poucos valores
   http_requests{status="200"}
   ```

3. **Nomeie métricas com sufixos apropriados**
   ```go
   requests_total  // counter
   requests_duration_seconds  // histogram
   memory_bytes  // gauge
   ```

4. **Adicione help text descritivo**
   ```go
   Help: "Total de requisições HTTP recebidas (Traffic)"
   ```

### ❌ **DON'T (Não Faça)**

1. **Não use labels com alta cardinalidade**
   ```go
   // ❌ Evite: IDs únicos, timestamps, emails
   {user_id="unique-123"}
   {timestamp="2024-10-04T12:00:00Z"}
   ```

2. **Não crie métricas demais**
   - Foque nas métricas que importam
   - Golden Signals + principais métricas de negócio

3. **Não use Counter para valores que diminuem**
   ```go
   // ❌ Errado: usar Counter
   activeUsers.Inc()  // Quando usuário conecta
   activeUsers.Dec()  // ❌ Não existe Dec() em Counter!
   
   // ✅ Correto: usar Gauge
   activeUsers.Set(count)
   ```

## 🔍 Troubleshooting

### **Métricas não aparecem**

1. Verifique se `/metrics` está acessível:
   ```bash
   curl http://localhost:8080/metrics
   ```

2. Verifique configuração do Prometheus:
   ```yaml
   targets: ['localhost:8080']  # Correto
   ```

3. Verifique logs do Prometheus

### **Latência das métricas**

- Prometheus faz scraping a cada `scrape_interval` (padrão: 15s)
- Métricas levam até 15s para aparecer
- Use `evaluation_interval` menor se precisar de mais agilidade

## 📚 Recursos

- **Prometheus Docs**: https://prometheus.io/docs/
- **Golden Signals**: https://sre.google/sre-book/monitoring-distributed-systems/
- **PromQL Tutorial**: https://prometheus.io/docs/prometheus/latest/querying/basics/
- **Grafana Docs**: https://grafana.com/docs/

## 🎓 Conceitos Importantes

### **Scraping vs Push**

**Prometheus (Pull-based):**
- ✅ Servidor controla quando coletar
- ✅ Fácil detectar se target está down
- ✅ Mais simples de debugar

**Push-based (ex: StatsD):**
- Aplicação envia métricas
- Usado quando pull não é possível

### **Retention (Retenção)**

Prometheus armazena dados localmente:
- Padrão: 15 dias
- Configurável via `--storage.tsdb.retention.time`
- Para long-term storage, use Thanos ou Cortex

### **High Availability**

Para produção:
- Execute múltiplas instâncias do Prometheus
- Use Thanos ou Cortex para HA
- Configure alertmanager em cluster

---

**Anterior:** [Swagger Documentation](06-swagger-documentation.md) | **Voltar ao início:** [README](README.md)

