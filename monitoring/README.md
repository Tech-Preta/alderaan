# 📊 Monitoramento - Prometheus e Grafana

Este diretório contém toda a configuração necessária para monitorar a API Alderaan com Prometheus e Grafana.

## 🚀 Início Rápido

### **Opção 1: Stack Completa (Recomendado)**

Inicie todos os serviços de uma vez (PostgreSQL + API + Prometheus + Grafana):

```bash
docker-compose up -d
```

Aguarde ~30 segundos para todos os serviços iniciarem. Você terá acesso a:

- **API**: `http://localhost:8080`
- **Swagger**: `http://localhost:8080/swagger/index.html`
- **Métricas**: `http://localhost:8080/metrics`
- **Prometheus**: `http://localhost:9090`
- **Grafana**: `http://localhost:3000` (admin/admin)
- **PostgreSQL**: `localhost:5432`

### **Opção 2: Serviços Individuais**

```bash
# Iniciar apenas banco + API
docker-compose up -d postgres api

# Adicionar monitoramento depois
docker-compose up -d prometheus grafana

# Ou usar Makefile
make db-up              # Apenas PostgreSQL
make monitoring-up      # Prometheus + Grafana via compose principal
```

### **3. Gerar Dados de Teste**

Crie alguns produtos para gerar métricas:

```bash
# Criar produto 1
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Notebook",
    "sku": 12345,
    "categories": ["Eletrônicos", "Computadores"],
    "price": 3500
  }'

# Criar produto 2
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mouse Gamer",
    "sku": 67890,
    "categories": ["Periféricos", "Gaming"],
    "price": 250
  }'

# Listar produtos (gera métricas de leitura)
for i in {1..10}; do
  curl -s http://localhost:8080/api/v1/products > /dev/null
  sleep 1
done

# Buscar produto inexistente (gera erro 404)
curl http://localhost:8080/api/v1/products/ProdutoInexistente
```

### **4. Visualizar Métricas**

#### **No Prometheus** (`http://localhost:9090`)

Execute estas queries na aba "Graph":

```promql
# Taxa de requisições por segundo
rate(http_requests_total[1m])

# Latência P95
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Taxa de erro
rate(http_request_errors_total[5m])

# Requisições em andamento
http_in_flight_requests

# Produtos criados
products_created_total

# Total de produtos
products_total
```

#### **No Grafana** (`http://localhost:3000`)

1. Faça login (admin/admin)
2. O datasource Prometheus já está configurado automaticamente
3. Crie um dashboard ou use as queries acima

## 📁 Estrutura de Arquivos

```
monitoring/
├── prometheus.yml              # Configuração do Prometheus
├── alerts.yml                  # Regras de alerta (15+ alertas)
├── grafana/
│   ├── provisioning/
│   │   ├── datasources/
│   │   │   └── prometheus.yml  # Auto-config do datasource
│   │   └── dashboards/
│   │       └── default.yml     # Auto-load de dashboards
│   └── dashboards/             # Dashboards JSON (você pode adicionar aqui)
└── README.md                   # Documentação de monitoramento

Nota: O docker-compose.yml está na raiz do projeto e inclui todos os serviços:
- PostgreSQL
- API
- Prometheus
- Grafana
```

## 📊 Métricas Disponíveis

### **Golden Signals**

| Métrica | Tipo | Descrição |
|---------|------|-----------|
| `http_request_duration_seconds` | Histogram | Latência das requisições |
| `http_requests_total` | Counter | Total de requisições |
| `http_request_errors_total` | Counter | Total de erros |
| `http_in_flight_requests` | Gauge | Requisições simultâneas |

### **Métricas de Negócio**

| Métrica | Tipo | Descrição |
|---------|------|-----------|
| `products_created_total` | Counter | Produtos criados (total) |
| `products_total` | Gauge | Produtos atuais |
| `products_by_category` | Gauge | Produtos por categoria |
| `products_total_value` | Gauge | Valor total do inventário |
| `products_average_price` | Gauge | Preço médio |

## 🎯 Queries PromQL Úteis

### **Latência**

```promql
# Latência média
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# P50 (mediana)
histogram_quantile(0.50, rate(http_request_duration_seconds_bucket[5m]))

# P95
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# P99
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))

# Latência por endpoint
histogram_quantile(0.95,
  sum(rate(http_request_duration_seconds_bucket[5m])) by (endpoint, le)
)
```

### **Tráfego**

```promql
# Requisições por segundo (RPS)
rate(http_requests_total[5m])

# RPS por endpoint
sum(rate(http_requests_total[5m])) by (endpoint)

# RPS por método HTTP
sum(rate(http_requests_total[5m])) by (method)

# Total de requisições nas últimas 24h
increase(http_requests_total[24h])
```

### **Erros**

```promql
# Taxa de erro (%)
(sum(rate(http_request_errors_total[5m])) / sum(rate(http_requests_total[5m]))) * 100

# Erros por endpoint
sum(rate(http_request_errors_total[5m])) by (endpoint)

# Erros 4xx vs 5xx
sum(rate(http_request_errors_total[5m])) by (error_type)

# Erros por status code
sum(rate(http_request_errors_total[5m])) by (status)
```

### **Saturação**

```promql
# Requisições em andamento (atual)
http_in_flight_requests

# Média de requisições simultâneas (5 min)
avg_over_time(http_in_flight_requests[5m])

# Máximo de requisições simultâneas (5 min)
max_over_time(http_in_flight_requests[5m])
```

### **Negócio**

```promql
# Taxa de criação de produtos (por hora)
rate(products_created_total[1h]) * 3600

# Total de produtos
products_total

# Crescimento de produtos (últimas 24h)
products_total - (products_total offset 24h)

# Top 5 categorias
topk(5, products_by_category)

# Valor do inventário
products_total_value

# Preço médio
products_average_price
```

## 🚨 Alertas Configurados

Os alertas estão definidos em `alerts.yml`:

### **Alertas de Golden Signals**

- ✅ **HighLatencyP99**: Latência P99 > 1s por 5 min
- ✅ **HighErrorRate**: Taxa de erro > 5% por 5 min
- ✅ **HighConcurrentRequests**: > 100 requisições simultâneas por 5 min
- ✅ **TrafficDropped**: Queda súbita de tráfego

### **Alertas de Negócio**

- ✅ **NoProductsCreated**: Nenhum produto criado em 2h
- ✅ **EmptyInventory**: Inventário vazio
- ✅ **LowInventoryValue**: Valor total < $1000

### **Alertas de Sistema**

- ✅ **APIDown**: API não responde
- ✅ **APIFlapping**: API reiniciando repetidamente
- ✅ **SlowScrape**: Coleta de métricas lenta

## 🛠️ Comandos Úteis

### **Docker Compose**

```bash
# Iniciar todos os serviços (da raiz do projeto)
docker-compose up -d

# Iniciar apenas monitoramento
docker-compose up -d prometheus grafana

# Ver logs
docker-compose logs -f prometheus
docker-compose logs -f grafana

# Parar serviços
docker-compose stop prometheus grafana

# Parar tudo
docker-compose down

# Parar e remover volumes (limpa dados)
docker-compose down -v

# Recarregar configuração do Prometheus (sem restart)
curl -X POST http://localhost:9090/-/reload

# Ou usar Makefile (da raiz do projeto)
make monitoring-up      # Iniciar prometheus + grafana
make monitoring-down    # Parar monitoramento
make monitoring-logs    # Ver logs
```

### **Verificar Saúde**

```bash
# Prometheus
curl http://localhost:9090/-/healthy

# Grafana
curl http://localhost:3000/api/health

# API (métricas)
curl http://localhost:8080/metrics
```

## 📈 Criando Dashboard no Grafana

1. **Acesse Grafana**: `http://localhost:3000`
2. **Login**: admin/admin
3. **Criar Dashboard**: + → Dashboard → Add visualization
4. **Selecionar Datasource**: Prometheus
5. **Adicionar Query**: Use queries da seção anterior
6. **Configurar Visualização**: Escolha tipo de gráfico
7. **Salvar Dashboard**

### **Painéis Recomendados**

#### **Painel 1: Golden Signals**

- **Latência P95**: Line graph
- **RPS**: Line graph
- **Taxa de Erro**: Line graph com threshold
- **Requisições Simultâneas**: Gauge

#### **Painel 2: Negócio**

- **Produtos Criados**: Counter
- **Total de Produtos**: Stat
- **Produtos por Categoria**: Bar chart
- **Valor do Inventário**: Stat
- **Preço Médio**: Gauge

## 🔧 Troubleshooting

### **Prometheus não coleta métricas da API**

1. Verifique se a API está rodando:
   ```bash
   curl http://localhost:8080/health
   ```

2. Verifique se `/metrics` está acessível:
   ```bash
   curl http://localhost:8080/metrics
   ```

3. Verifique targets no Prometheus:
   - Acesse `http://localhost:9090/targets`
   - Verifique se `alderaan-api` está UP

4. Se estiver usando Docker no Mac/Windows:
   - Use `host.docker.internal:8080` no `prometheus.yml`

5. Se estiver usando Docker no Linux:
   - Use `172.17.0.1:8080` no `prometheus.yml`

### **Grafana não conecta no Prometheus**

1. Verifique se ambos estão na mesma rede Docker
2. Use `http://prometheus:9090` (nome do serviço, não localhost)
3. Teste conexão:
   ```bash
   docker-compose exec grafana curl http://prometheus:9090/-/healthy
   ```

### **Sem dados nas queries**

1. Aguarde 15-30 segundos após iniciar (scrape interval)
2. Gere tráfego na API
3. Verifique range de tempo no Grafana (últimos 5-15 minutos)

## 📚 Documentação Completa

Para mais detalhes sobre monitoramento, consulte:

📖 [docs/07-prometheus-monitoring.md](../docs/07-prometheus-monitoring.md)

## 🎓 Próximos Passos

1. ✅ Crie dashboards personalizados no Grafana
2. ✅ Configure AlertManager para enviar notificações
3. ✅ Adicione mais métricas de negócio conforme necessário
4. ✅ Configure retention apropriado para produção
5. ✅ Considere Thanos/Cortex para long-term storage

---

**🚀 Pronto! Agora você tem monitoramento completo com Golden Signals e métricas de negócio!**
