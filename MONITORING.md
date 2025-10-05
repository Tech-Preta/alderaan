# 🎯 Configuração de Monitoramento Completa

Resumo da configuração da stack completa de observabilidade para o projeto Alderaan.

## ✅ O que foi Criado

### **1. Dashboard do Grafana** 📈

**Arquivo**: `monitoring/grafana/dashboards/alderaan-overview.json` (829 linhas)

**Painéis Criados** (13 no total):

#### **Golden Signals**
1. **Requests/sec (Traffic)** - Stat panel
   - Taxa de requisições por segundo
   
2. **P95 Latency** - Stat panel
   - Latência P95 com thresholds (amarelo > 0.5s, vermelho > 1s)
   
3. **Error Rate** - Stat panel
   - Taxa de erro percentual
   
4. **In-Flight Requests (Saturation)** - Stat panel
   - Requisições concorrentes

5. **Request Rate by Endpoint** - Time series
   - Gráfico de requisições por endpoint (GET/POST)
   
6. **Latency Percentiles by Endpoint** - Time series
   - P50, P95, P99 por endpoint

7. **HTTP Status Codes** - Time series (stacked)
   - Distribuição de status codes ao longo do tempo

#### **Métricas de Negócio**
8. **Total Products** - Stat panel
   - Total de produtos cadastrados
   
9. **Total Inventory Value** - Stat panel
   - Valor total em USD
   
10. **Average Product Price** - Stat panel
    - Preço médio em USD
    
11. **Products Created (Total)** - Stat panel
    - Total de produtos criados
    
12. **Products by Category** - Pie chart
    - Distribuição visual por categoria
    
13. **Product Creation Rate** - Time series
    - Taxa de criação de produtos/segundo

**Configurações**:
- Auto-refresh: 10 segundos
- Time range: Última 1 hora
- Tags: alderaan, api, golang
- UID: `alderaan-api-overview`

### **2. Alertmanager** 🚨

**Arquivo**: `monitoring/alertmanager.yml` (131 linhas)

**Configuração de Rotas**:
- **Critical**: Espera 5s, repete a cada 5 minutos
- **Warning**: Espera 30s, repete a cada 1 hora
- **Info**: Espera 1 minuto, repete a cada 24 horas

**Receivers Configurados**:
1. `default` - Webhook padrão
2. `critical-alerts` - Alertas críticos (pronto para Slack/Email)
3. `warning-alerts` - Warnings (pronto para Slack)
4. `info-alerts` - Informativos

**Regras de Inibição**:
- Se API está down → não alerta sobre latência
- Se há erro crítico → silencia warnings relacionados

**Integrações Prontas** (comentadas, é só descomentar):
- ✅ Email (SMTP)
- ✅ Slack
- ✅ Webhook customizado
- ✅ Templates customizados

### **3. Regras de Alerta** 📋

**Arquivo**: `monitoring/alerts.yml` (186 linhas - já existia, agora ativo)

**15+ Regras Configuradas**:

#### **Golden Signals**
- `HighLatencyP99` - P99 > 1s (warning)
- `CriticalHighLatencyP99` - P99 > 3s (critical)
- `HighErrorRate` - Erros > 5% (warning)
- `CriticalErrorRate` - Erros > 10% (critical)
- `HighConcurrentRequests` - > 50 req in-flight (warning)
- `NoTrafficDetected` - Zero tráfego por 5 minutos (warning)
- `APIDown` - API não responde (critical)

#### **Métricas de Negócio**
- `HighProductCreationRate` - > 10 produtos/min
- `NoProductsCreated` - Zero produtos em 1 hora
- `InventoryValueHigh` - Valor > $100k
- `LowAveragePrice` - Preço médio < $5
- `CategoryImbalance` - Categoria > 50% total

#### **Sistema**
- `PrometheusDown` - Prometheus down
- `TargetDown` - Target não responde

Todas incluem:
- `description`: Descrição clara do problema
- `summary`: Resumo executivo
- `runbook_url`: Link para documentação
- `severity`: critical/warning/info

### **4. Docker Compose Atualizado** 🐳

**Adicionado serviço Alertmanager**:
```yaml
alertmanager:
  image: prom/alertmanager:latest
  container_name: alderaan-alertmanager
  ports:
    - "9093:9093"
  volumes:
    - ./monitoring/alertmanager.yml:/etc/alertmanager/alertmanager.yml:ro
    - alertmanager_data:/alertmanager
```

**Volume adicionado**:
- `alertmanager_data` - Persistência de dados do Alertmanager

### **5. Prometheus Configurado** 📊

**Arquivo**: `monitoring/prometheus.yml`

**Mudanças**:
- ✅ Alertmanager habilitado: `alertmanager:9093`
- ✅ Regras de alerta ativas: `/etc/prometheus/alerts.yml`

### **6. Makefile Atualizado** 🔧

**Novos comandos incluem Alertmanager**:
```bash
make monitoring-up    # Inicia Prometheus + Alertmanager + Grafana
make monitoring-down  # Para todos
make monitoring-logs  # Logs de todos
```

**Output do `platform-up` mostra**:
```
🗄️  PostgreSQL:   http://localhost:5432
🚀 API:           http://localhost:8080
📚 Swagger:       http://localhost:8080/swagger/index.html
📊 Prometheus:    http://localhost:9090
🚨 Alertmanager:  http://localhost:9093  ← NOVO!
📈 Grafana:       http://localhost:3000
```

### **7. Documentação** 📚

**Arquivos atualizados**:
- `monitoring/README.md` - Guia completo de monitoramento
- `README.md` - Seção de monitoramento expandida
- `QUICKSTART.md` - Inclui Alertmanager
- `.gitignore` - Ignora dados de volumes

## 🎯 Status dos Serviços

```
✅ postgres       - Up 2 minutes (healthy)
✅ alertmanager   - Up 2 minutes
✅ api            - Up 1 minute (healthy)
✅ prometheus     - Up 1 minute
✅ grafana        - Up 1 minute
```

## 🔍 Verificações Realizadas

### **1. Dashboard do Grafana**
```bash
# Logs confirmam provisionamento
✅ "finished to provision dashboards"
✅ Dashboard JSON carregada de /var/lib/grafana/dashboards/
```

### **2. Alertmanager**
```bash
# API responde corretamente
✅ GET /api/v2/status - Config carregada
✅ GET /api/v2/alerts - Pronto para receber alertas
✅ 3 receivers configurados (critical, warning, info)
```

### **3. Prometheus**
```bash
# Regras carregadas
✅ 15+ regras de alerta detectadas
✅ Todas com estado "inactive" (normal)
✅ Integração com Alertmanager ativa
```

### **4. Métricas**
```bash
# Tráfego gerado para teste
✅ 10 produtos criados via API
✅ Métricas coletadas pelo Prometheus
✅ Golden Signals + Business Metrics funcionando
```

## 📊 Como Acessar

### **Grafana Dashboard**
1. Acesse: http://localhost:3000
2. Login: `admin` / `admin`
3. Menu: Dashboards → "Alderaan API - Overview"
4. Auto-refresh a cada 10s

### **Prometheus**
1. Acesse: http://localhost:9090
2. Menu: Status → Rules (ver todas as regras)
3. Menu: Alerts (ver alertas ativos)
4. Graph: Testar queries PromQL

### **Alertmanager**
1. Acesse: http://localhost:9093
2. Ver alertas disparados (atualmente: 0)
3. Criar silences para manutenção
4. Ver configuração de receivers

## 🚀 Próximos Passos

### **Configurar Notificações**

#### **Slack**
Edite `monitoring/alertmanager.yml`:
```yaml
receivers:
  - name: 'critical-alerts'
    slack_configs:
      - channel: '#alerts'
        api_url: 'https://hooks.slack.com/services/XXX/YYY/ZZZ'
```

#### **Email**
```yaml
global:
  smtp_smarthost: 'smtp.gmail.com:587'
  smtp_from: 'alerts@alderaan.com'
  smtp_auth_username: 'your-email@gmail.com'
  smtp_auth_password: 'app-password'
```

### **Testar Alertas**

```bash
# 1. Gerar latência alta (simular)
# Fazer várias requisições simultâneas

# 2. Gerar erros (tentar criar produto duplicado)
curl -X POST http://localhost:8080/api/v1/products \
  -d '{"name":"Duplicate","sku":9999,"categories":["test"],"price":100}'

# 3. Ver alertas no Prometheus
curl http://localhost:9090/api/v1/alerts | jq

# 4. Ver no Alertmanager
curl http://localhost:9093/api/v2/alerts | jq
```

### **Customizar Dashboard**

1. Acesse Grafana
2. Abra "Alderaan API - Overview"
3. Clique em "Edit" em qualquer painel
4. Adicione/remova/modifique queries
5. Salve com "Save Dashboard"

### **Criar Novas Dashboards**

```bash
# Exportar dashboard existente
curl -u admin:admin http://localhost:3000/api/dashboards/uid/alderaan-api-overview | \
  jq '.dashboard' > my-dashboard.json

# Editar e colocar em monitoring/grafana/dashboards/
# Reiniciar Grafana
docker-compose restart grafana
```

## 📁 Estrutura Final

```
monitoring/
├── alertmanager.yml                      # ← NOVO: Config Alertmanager
├── alerts.yml                            # ← ATIVO: 15+ regras
├── prometheus.yml                        # ← ATUALIZADO: Com Alertmanager
├── grafana/
│   ├── dashboards/
│   │   └── alderaan-overview.json        # ← NOVO: Dashboard completa
│   └── provisioning/
│       ├── datasources/
│       │   └── prometheus.yml
│       └── dashboards/
│           └── default.yml
└── README.md                             # ← ATUALIZADO: Guia completo
```

## ✨ Resumo

### **Antes**
- ❌ Sem dashboard no Grafana
- ❌ Sem Alertmanager
- ❌ Regras de alerta comentadas
- ❌ Sem notificações

### **Depois**
- ✅ Dashboard completa com 13 painéis
- ✅ Alertmanager rodando e configurado
- ✅ 15+ regras de alerta ativas
- ✅ Pronto para notificações (Slack/Email)
- ✅ Golden Signals + Business Metrics
- ✅ Auto-refresh e time ranges
- ✅ Rotas e inibições configuradas
- ✅ Documentação completa

## 🎓 Recursos

- **Grafana Dashboard Docs**: https://grafana.com/docs/grafana/latest/dashboards/
- **Alertmanager Config**: https://prometheus.io/docs/alerting/latest/configuration/
- **PromQL Queries**: https://prometheus.io/docs/prometheus/latest/querying/basics/
- **Golden Signals**: https://sre.google/sre-book/monitoring-distributed-systems/

---

**Criado em**: 2025-10-04  
**Stack**: Prometheus 2.x + Alertmanager 0.27+ + Grafana 11.x  
**Status**: ✅ Produção-ready

