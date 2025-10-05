# 📊 Prometheus Queries (PromQL)

Guia completo de queries PromQL disponíveis para consultar métricas da API Alderaan.

## 📋 Índice

- [Métricas Disponíveis](#-métricas-disponíveis)
- [Golden Signals](#-golden-signals)
- [Métricas de Negócio](#-métricas-de-negócio)
- [Queries Básicas](#-queries-básicas)
- [Queries Avançadas](#-queries-avançadas)
- [Funções PromQL Úteis](#-funções-promql-úteis)
- [Exemplos Práticos](#-exemplos-práticos)

## 📈 Métricas Disponíveis

### **HTTP Metrics (Golden Signals)**

```promql
# Histogram - Duração das requisições HTTP
http_request_duration_seconds_bucket{method="GET",path="/api/v1/products",le="0.1"}
http_request_duration_seconds_sum{method="GET",path="/api/v1/products"}
http_request_duration_seconds_count{method="GET",path="/api/v1/products"}

# Counter - Total de requisições HTTP
http_requests_total{method="GET",path="/api/v1/products",status="200"}

# Counter - Total de erros HTTP
http_request_errors_total{path="/api/v1/products",error_type="validation"}

# Gauge - Requisições em voo (concorrentes)
http_in_flight_requests
```

**Labels disponíveis**:
- `method`: Método HTTP (GET, POST, PUT, DELETE)
- `path`: Endpoint acessado
- `status`: Status code HTTP (200, 400, 500, etc)
- `error_type`: Tipo de erro (validation, database, etc)
- `le`: Less than or equal (para histograms)

### **Business Metrics (Métricas de Negócio)**

```promql
# Counter - Total de produtos criados (acumulado)
products_created_total

# Gauge - Total de produtos no banco
products_total

# Gauge - Valor total do inventário (em centavos)
products_total_value

# Gauge - Preço médio dos produtos (em centavos)
products_average_price

# Gauge - Produtos por categoria
products_by_category{category="Eletrônicos"}
products_by_category{category="Livros"}
```

**Labels disponíveis**:
- `category`: Nome da categoria do produto

## 🎯 Golden Signals

### **1. Latency (Latência)**

#### **Percentis de Latência**

```promql
# P50 (Mediana) - 50% das requisições são mais rápidas que este valor
histogram_quantile(0.50, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# P95 - 95% das requisições são mais rápidas que este valor
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# P99 - 99% das requisições são mais rápidas que este valor
histogram_quantile(0.99, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# P99.9 - 99.9% das requisições são mais rápidas que este valor
histogram_quantile(0.999, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```

#### **Latência por Endpoint**

```promql
# P95 por endpoint
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le, method, path)
)

# P50 apenas para GET /products
histogram_quantile(0.50,
  sum(rate(http_request_duration_seconds_bucket{method="GET",path="/api/v1/products"}[5m])) by (le)
)
```

#### **Latência Média**

```promql
# Latência média geral
rate(http_request_duration_seconds_sum[5m]) /
rate(http_request_duration_seconds_count[5m])

# Latência média por endpoint
rate(http_request_duration_seconds_sum[5m]) /
rate(http_request_duration_seconds_count[5m])
by (method, path)
```

#### **Latência Máxima**

```promql
# Latência máxima nos últimos 5 minutos
max_over_time(http_request_duration_seconds_sum[5m])

# Por endpoint
max_over_time(
  http_request_duration_seconds_sum{method="POST",path="/api/v1/products"}[5m]
)
```

### **2. Traffic (Tráfego)**

#### **Requisições por Segundo (RPS)**

```promql
# RPS total
sum(rate(http_requests_total[5m]))

# RPS nos últimos 1 minuto (mais granular)
sum(rate(http_requests_total[1m]))

# RPS por endpoint
sum(rate(http_requests_total[5m])) by (method, path)

# RPS apenas para POST (criação)
sum(rate(http_requests_total{method="POST"}[5m]))

# RPS apenas para endpoints de produtos
sum(rate(http_requests_total{path=~"/api/v1/products.*"}[5m]))
```

#### **Total de Requisições**

```promql
# Total acumulado de requisições
sum(http_requests_total)

# Total nos últimos 5 minutos
sum(increase(http_requests_total[5m]))

# Total nas últimas 24 horas
sum(increase(http_requests_total[24h]))

# Por status code
sum(http_requests_total) by (status)
```

#### **Requisições por Status Code**

```promql
# Distribuição por status
sum(rate(http_requests_total[5m])) by (status)

# Apenas sucesso (2xx)
sum(rate(http_requests_total{status=~"2.."}[5m]))

# Apenas erros de cliente (4xx)
sum(rate(http_requests_total{status=~"4.."}[5m]))

# Apenas erros de servidor (5xx)
sum(rate(http_requests_total{status=~"5.."}[5m]))
```

#### **Top Endpoints Mais Acessados**

```promql
# Top 5 endpoints por tráfego
topk(5, sum(rate(http_requests_total[5m])) by (method, path))

# Bottom 5 (menos acessados)
bottomk(5, sum(rate(http_requests_total[5m])) by (method, path))
```

### **3. Errors (Erros)**

#### **Taxa de Erro**

```promql
# Taxa de erro geral (porcentagem)
sum(rate(http_request_errors_total[5m])) /
sum(rate(http_requests_total[5m])) * 100

# Taxa de erro simplificada (0-1)
sum(rate(http_request_errors_total[5m])) /
sum(rate(http_requests_total[5m]))

# Taxa de erro por endpoint
sum(rate(http_request_errors_total[5m])) by (path) /
sum(rate(http_requests_total[5m])) by (path)
```

#### **Total de Erros**

```promql
# Total de erros
sum(http_request_errors_total)

# Erros nos últimos 5 minutos
sum(increase(http_request_errors_total[5m]))

# Erros por tipo
sum(http_request_errors_total) by (error_type)

# Erros por endpoint
sum(rate(http_request_errors_total[5m])) by (path)
```

#### **Taxa de Erros HTTP 5xx**

```promql
# Taxa de erros 5xx
sum(rate(http_requests_total{status=~"5.."}[5m])) /
sum(rate(http_requests_total[5m]))

# Taxa de erros 4xx
sum(rate(http_requests_total{status=~"4.."}[5m])) /
sum(rate(http_requests_total[5m]))
```

#### **Endpoints com Mais Erros**

```promql
# Top 5 endpoints com mais erros
topk(5, sum(rate(http_request_errors_total[5m])) by (path))

# Endpoints com taxa de erro > 5%
(
  sum(rate(http_request_errors_total[5m])) by (path) /
  sum(rate(http_requests_total[5m])) by (path)
) > 0.05
```

### **4. Saturation (Saturação)**

#### **Requisições Concorrentes**

```promql
# Requisições em voo (atualmente)
http_in_flight_requests

# Média de requisições em voo nos últimos 5 minutos
avg_over_time(http_in_flight_requests[5m])

# Máximo de requisições em voo
max_over_time(http_in_flight_requests[5m])

# Mínimo de requisições em voo
min_over_time(http_in_flight_requests[5m])
```

#### **Saturação de Capacidade**

```promql
# Porcentagem de capacidade usada (assumindo limite de 100)
http_in_flight_requests / 100 * 100

# Alerta se saturação > 50%
http_in_flight_requests > 50

# Alerta se saturação > 80%
http_in_flight_requests > 80
```

## 💼 Métricas de Negócio

### **Produtos**

#### **Total de Produtos**

```promql
# Total atual de produtos no banco
products_total

# Variação nas últimas 24 horas
products_total - products_total offset 24h

# Taxa de crescimento (%)
((products_total - products_total offset 1h) / products_total offset 1h) * 100
```

#### **Taxa de Criação de Produtos**

```promql
# Produtos criados por segundo
rate(products_created_total[5m])

# Produtos criados por minuto
rate(products_created_total[5m]) * 60

# Produtos criados por hora
rate(products_created_total[5m]) * 3600

# Total criado nas últimas 24 horas
increase(products_created_total[24h])

# Total criado hoje (desde meia-noite)
increase(products_created_total[1d:] @ start())
```

#### **Produtos por Categoria**

```promql
# Todas as categorias
products_by_category

# Categoria específica
products_by_category{category="Eletrônicos"}

# Top 5 categorias
topk(5, products_by_category)

# Porcentagem de cada categoria
products_by_category / sum(products_by_category) * 100

# Categorias com mais de 10 produtos
products_by_category > 10

# Categorias com menos produtos
bottomk(3, products_by_category)
```

### **Valores e Preços**

#### **Valor do Inventário**

```promql
# Valor total em centavos
products_total_value

# Valor total em dólares
products_total_value / 100

# Valor total em reais (assumindo conversão 1:5)
products_total_value / 100 * 5

# Variação de valor nas últimas 24h
products_total_value - products_total_value offset 24h

# Taxa de crescimento do valor
((products_total_value - products_total_value offset 1h) /
 products_total_value offset 1h) * 100
```

#### **Preço Médio**

```promql
# Preço médio em centavos
products_average_price

# Preço médio em dólares
products_average_price / 100

# Variação do preço médio
products_average_price - products_average_price offset 1h

# Preço médio por categoria (requer agregação manual)
products_total_value / products_total
```

#### **Valor por Categoria**

```promql
# Valor estimado por categoria (preço médio * quantidade)
products_by_category * products_average_price

# Top 3 categorias por valor
topk(3, products_by_category * products_average_price)
```

## 📝 Queries Básicas

### **Verificar se Métrica Existe**

```promql
# Ver últimos valores
http_requests_total

# Ver valores específicos
http_requests_total{method="GET"}

# Ver todas as labels disponíveis
count(http_requests_total) by (__name__)
```

### **Operações Matemáticas**

```promql
# Soma
sum(http_requests_total)

# Média
avg(http_request_duration_seconds_sum)

# Mínimo
min(products_by_category)

# Máximo
max(http_in_flight_requests)

# Contagem
count(http_requests_total)
```

### **Filtros por Label**

```promql
# Um label específico
http_requests_total{method="GET"}

# Múltiplos labels (AND)
http_requests_total{method="GET", status="200"}

# Regex match
http_requests_total{path=~"/api/v1/.*"}

# Regex não-match
http_requests_total{status!~"5.."}

# Múltiplos valores (OR)
http_requests_total{status=~"200|201|204"}
```

### **Range Vectors**

```promql
# Últimos 5 minutos
http_requests_total[5m]

# Última 1 hora
http_requests_total[1h]

# Último 1 dia
http_requests_total[1d]

# Com offset (1 hora atrás)
http_requests_total[5m] offset 1h
```

## 🚀 Queries Avançadas

### **SLO (Service Level Objective)**

#### **Disponibilidade**

```promql
# SLO de 99.9% de disponibilidade (0.999)
sum(rate(http_requests_total{status!~"5.."}[5m])) /
sum(rate(http_requests_total[5m])) > 0.999

# Porcentagem de disponibilidade
sum(rate(http_requests_total{status!~"5.."}[5m])) /
sum(rate(http_requests_total[5m])) * 100
```

#### **Latência (P95 < 1s)**

```promql
# Verificar se P95 está abaixo de 1 segundo
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (le)
) < 1
```

### **Correlações**

#### **Latência vs Tráfego**

```promql
# Latência aumenta com tráfego?
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
and
sum(rate(http_requests_total[5m]))
```

#### **Erros vs Latência**

```promql
# Ver se erros correlacionam com alta latência
(sum(rate(http_request_errors_total[5m])) / sum(rate(http_requests_total[5m]))) > 0.05
and
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 1
```

### **Tendências**

#### **Crescimento de Tráfego**

```promql
# Comparar tráfego atual vs 1 hora atrás
(sum(rate(http_requests_total[5m])) -
 sum(rate(http_requests_total[5m] offset 1h))) /
sum(rate(http_requests_total[5m] offset 1h)) * 100
```

#### **Crescimento de Produtos**

```promql
# Taxa de crescimento por hora
deriv(products_total[1h])

# Produtos criados por hora (média)
avg_over_time(rate(products_created_total[5m])[1h:]) * 3600
```

### **Previsões (Predictions)**

```promql
# Prever quando atingiremos 1000 produtos (linear regression)
predict_linear(products_total[1h], 3600)

# Prever tráfego nas próximas 2 horas
predict_linear(sum(rate(http_requests_total[5m]))[1h], 7200)
```

## 🔧 Funções PromQL Úteis

### **Funções de Agregação**

```promql
sum()        # Soma todos os valores
avg()        # Média
min()        # Mínimo
max()        # Máximo
count()      # Contagem
stddev()     # Desvio padrão
stdvar()     # Variância
```

### **Funções de Range**

```promql
rate()              # Taxa de mudança por segundo
irate()             # Taxa instantânea (para spikes)
increase()          # Aumento absoluto no período
delta()             # Diferença entre primeiro e último valor
idelta()            # Diferença instantânea
```

### **Funções over Time**

```promql
avg_over_time()     # Média durante o período
min_over_time()     # Mínimo
max_over_time()     # Máximo
sum_over_time()     # Soma
count_over_time()   # Contagem
quantile_over_time()# Quantil
stddev_over_time()  # Desvio padrão
```

### **Funções de Ordenação**

```promql
sort()              # Ordem crescente
sort_desc()         # Ordem decrescente
topk(n, ...)        # Top N valores
bottomk(n, ...)     # Bottom N valores
```

### **Funções de Tempo**

```promql
time()              # Timestamp atual
timestamp()         # Timestamp do sample
year()              # Ano
month()             # Mês
day_of_month()      # Dia do mês
hour()              # Hora
minute()            # Minuto
```

## 💡 Exemplos Práticos

### **Dashboard de Overview**

```promql
# Painel 1: RPS Total
sum(rate(http_requests_total[1m]))

# Painel 2: Latência P95
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# Painel 3: Taxa de Erro
sum(rate(http_request_errors_total[5m])) / sum(rate(http_requests_total[5m])) * 100

# Painel 4: Requisições em Voo
http_in_flight_requests
```

### **Alertas Importantes**

```promql
# Alerta: Alta Latência (P95 > 1s)
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le)) > 1

# Alerta: Taxa de Erro Alta (> 5%)
sum(rate(http_request_errors_total[5m])) / sum(rate(http_requests_total[5m])) > 0.05

# Alerta: Sem Tráfego (zero requisições em 5 min)
sum(rate(http_requests_total[5m])) == 0

# Alerta: Alta Saturação (> 50 requisições concorrentes)
http_in_flight_requests > 50
```

### **Análise de Performance**

```promql
# Qual endpoint é mais lento? (Top 5)
topk(5,
  histogram_quantile(0.95,
    sum(rate(http_request_duration_seconds_bucket[5m])) by (le, path)
  )
)

# Qual endpoint tem mais erros? (Top 5)
topk(5, sum(rate(http_request_errors_total[5m])) by (path))

# Qual endpoint tem mais tráfego? (Top 5)
topk(5, sum(rate(http_requests_total[5m])) by (path))
```

### **Análise de Negócio**

```promql
# Quantos produtos foram criados hoje?
increase(products_created_total[1d:] @ start())

# Qual categoria está crescendo mais rápido?
topk(1, deriv(products_by_category[1h]))

# Valor médio de produto criado por hora
(products_total_value - products_total_value offset 1h) /
(products_total - products_total offset 1h)
```

## 🌐 Como Usar

### **No Prometheus UI**

1. Acesse: http://localhost:9090
2. Vá em **Graph**
3. Cole a query no campo de texto
4. Clique em **Execute**
5. Escolha **Graph** ou **Console** para visualização

### **No Grafana**

1. Acesse: http://localhost:3000
2. Crie ou edite um painel
3. Selecione **Prometheus** como datasource
4. Cole a query no campo **Metrics**
5. Ajuste visualização (Time series, Stat, Gauge, etc)

### **Via API**

```bash
# Query instantânea
curl "http://localhost:9090/api/v1/query?query=sum(rate(http_requests_total[5m]))"

# Query com range
curl "http://localhost:9090/api/v1/query_range?query=sum(rate(http_requests_total[5m]))&start=2024-10-04T00:00:00Z&end=2024-10-04T23:59:59Z&step=15s"
```

## 🎓 Dicas de Otimização

### **1. Use Range Adequado**

```promql
# ❌ Range muito longo (lento)
sum(rate(http_requests_total[1d]))

# ✅ Range apropriado (rápido)
sum(rate(http_requests_total[5m]))
```

### **2. Agregue Antes de Calcular**

```promql
# ❌ Cálculo em cada série (lento)
histogram_quantile(0.95, http_request_duration_seconds_bucket)

# ✅ Agregação primeiro (rápido)
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```

### **3. Use Labels com Sabedoria**

```promql
# ❌ Muitos labels (cardinalidade alta)
sum(rate(http_requests_total[5m])) by (method, path, status, user_id, session_id)

# ✅ Labels essenciais (cardinalidade baixa)
sum(rate(http_requests_total[5m])) by (method, path)
```

## 📚 Recursos

- **PromQL Basics**: https://prometheus.io/docs/prometheus/latest/querying/basics/
- **PromQL Functions**: https://prometheus.io/docs/prometheus/latest/querying/functions/
- **PromQL Examples**: https://prometheus.io/docs/prometheus/latest/querying/examples/
- **Best Practices**: https://prometheus.io/docs/practices/naming/
- **Query Performance**: https://prometheus.io/docs/prometheus/latest/querying/basics/#time-series-selectors

---

**Anterior:** [Flyway Migrations](09-flyway-migrations.md) | **Próximo:** [Documentação](README.md) | **Voltar:** [README Principal](../README.md)
