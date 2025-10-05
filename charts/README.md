# Alderaan Helm Charts

Este diretório contém os Helm charts para deploy da aplicação Alderaan no Kubernetes.

## 📦 Charts Disponíveis

### Alderaan API
Chart principal que implanta a aplicação Alderaan com:
- ✅ API Go com DDD e Clean Architecture
- ✅ PostgreSQL (via subchart Bitnami)
- ✅ Horizontal Pod Autoscaler (HPA)
- ✅ Ingress com suporte a TLS
- ✅ Service Account e configurações de segurança
- ✅ ConfigMaps e Secrets para configuração
- ✅ Health checks e readiness probes
- ✅ Métricas Prometheus integradas
- ✅ Graceful shutdown
- Prometheus (via subchart community)
- Grafana (via subchart oficial)
- Migrations (Flyway)

## 🚀 Quick Start

### Instalação Rápida

```bash
# Instalar com valores padrão
helm install alderaan ./alderaan

# Instalar com valores customizados
helm install alderaan ./alderaan -f ./alderaan/examples/values-production.yaml

# Instalar em namespace específico
helm install alderaan ./alderaan --namespace alderaan --create-namespace
```

### Upgrade

```bash
helm upgrade alderaan ./alderaan -f custom-values.yaml
```

### Uninstall

```bash
helm uninstall alderaan
```

## 📋 Pré-requisitos

- **Kubernetes**: 1.23+
- **Helm**: 3.8+
- **Ingress Controller**: nginx-ingress, traefik, ou similar
- **Storage**: Provisioner para PersistentVolumes
- **cert-manager** (opcional): Para TLS automático com Let's Encrypt

### Instalando Pré-requisitos

```bash
# Instalar nginx-ingress
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install nginx-ingress ingress-nginx/ingress-nginx

# Instalar cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml
```

## 🎯 Exemplos de Uso

### Desenvolvimento Local (Minikube/Kind)

```bash
helm install alderaan ./alderaan -f ./alderaan/examples/values-development.yaml

# Port-forward para acesso local
kubectl port-forward svc/alderaan 8080:80
```

Acesse: http://localhost:8080

### Ambiente de Staging

```bash
helm install alderaan ./alderaan \
  --namespace staging \
  --create-namespace \
  -f ./alderaan/examples/values-staging.yaml \
  --set postgresql.auth.password="$(openssl rand -base64 32)"
```

### Ambiente de Produção

```bash
# Criar namespace
kubectl create namespace production

# Criar secret com senha do PostgreSQL
kubectl create secret generic alderaan-db-password \
  --from-literal=password="$(openssl rand -base64 32)" \
  -n production

# Instalar chart
helm install alderaan ./alderaan \
  --namespace production \
  -f ./alderaan/examples/values-production.yaml \
  --set postgresql.auth.existingSecret=alderaan-db-password
```

## 🔧 Customização

### Valores Principais

| Arquivo | Descrição |
|---------|-----------|
| `values.yaml` | Valores padrão |
| `examples/values-development.yaml` | Ambiente de desenvolvimento |
| `examples/values-staging.yaml` | Ambiente de staging |
| `examples/values-production.yaml` | Ambiente de produção |

### Override de Valores

```bash
# Via arquivo
helm install alderaan ./alderaan -f my-values.yaml

# Via command line
helm install alderaan ./alderaan \
  --set replicaCount=5 \
  --set image.tag=2.0.0 \
  --set postgresql.primary.persistence.size=200Gi
```

## 📊 Monitoramento

### Acessar Prometheus

```bash
kubectl port-forward svc/alderaan-prometheus-server 9090:80
# Acesse: http://localhost:9090
```

### Acessar Grafana

```bash
kubectl port-forward svc/alderaan-grafana 3000:80
# Acesse: http://localhost:3000
# User: admin / Password: (definido em values.yaml)
```

### Visualizar Métricas

```bash
# Métricas da aplicação
kubectl port-forward svc/alderaan 8080:80
curl http://localhost:8080/metrics
```

## 🗄️ Backup e Restore

### Backup do PostgreSQL

```bash
# Criar backup
kubectl exec -it alderaan-postgresql-0 -- \
  pg_dump -U alderaan alderaan_db > backup.sql

# Restaurar backup
kubectl exec -i alderaan-postgresql-0 -- \
  psql -U alderaan alderaan_db < backup.sql
```

### Backup de Volumes

```bash
# Usando Velero
velero backup create alderaan-backup --include-namespaces alderaan
```

## 🔍 Troubleshooting

### Verificar Status dos Pods

```bash
kubectl get pods -l app.kubernetes.io/name=alderaan
```

### Ver Logs

```bash
# Logs da aplicação
kubectl logs -l app.kubernetes.io/name=alderaan -f

# Logs do PostgreSQL
kubectl logs alderaan-postgresql-0 -f
```

### Testar Conectividade

```bash
# Test health endpoint
kubectl run test --rm -it --restart=Never --image=curlimages/curl -- \
  curl http://alderaan/health

# Test database connection
kubectl run psql --rm -it --restart=Never --image=postgres:15 -- \
  psql -h alderaan-postgresql -U alderaan -d alderaan_db
```

### Debug de Pods

```bash
# Describe pod
kubectl describe pod <pod-name>

# Shell interativo
kubectl exec -it <pod-name> -- /bin/sh

# Ver eventos
kubectl get events --sort-by=.metadata.creationTimestamp
```

## 📚 Estrutura do Chart

```
alderaan/
├── Chart.yaml              # Metadados do chart
├── values.yaml             # Valores padrão
├── .helmignore            # Arquivos a ignorar
├── README.md              # Documentação
├── templates/             # Templates Kubernetes
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── configmap.yaml
│   ├── secret.yaml
│   ├── serviceaccount.yaml
│   ├── hpa.yaml
│   ├── pdb.yaml
│   ├── servicemonitor.yaml
│   ├── networkpolicy.yaml
│   ├── migration-job.yaml
│   ├── _helpers.tpl
│   └── NOTES.txt
└── examples/              # Exemplos de configuração
    ├── values-development.yaml
    ├── values-staging.yaml
    └── values-production.yaml
```

## 🤝 Contribuindo

Para contribuir com melhorias no chart:

1. Fork o repositório
2. Crie uma branch (`git checkout -b feature/melhoria`)
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📄 Licença

MIT License - veja o arquivo LICENSE para detalhes

## 📞 Suporte

- **Issues**: https://github.com/williamkoller/golang-domain-driven-design/issues
- **Discussions**: https://github.com/williamkoller/golang-domain-driven-design/discussions
- **Documentação**: https://github.com/williamkoller/golang-domain-driven-design/tree/main/docs
