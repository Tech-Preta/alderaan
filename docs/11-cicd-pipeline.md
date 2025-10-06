# 🔄 CI/CD Pipeline - Publicação Automática

Este documento descreve o pipeline CI/CD implementado para publicação automática de Docker images e Helm charts no GitHub Container Registry.

## 📋 Visão Geral

O projeto utiliza GitHub Actions para automatizar:
- ✅ Build e publicação de imagens Docker multi-arquitetura
- ✅ Empacotamento e publicação de Helm charts
- ✅ Criação automática de releases com changelog
- ✅ Versionamento semântico
- ✅ Tags automáticas

## 🐳 Docker Build Pipeline

### Workflow: `.github/workflows/docker-build.yml`

**Triggers:**
- Push em branches `main` e `develop`
- Push de tags `v*`
- Pull requests para `main`

**Funcionalidades:**
- Geração automática da documentação Swagger antes do build
- Build multi-arquitetura (linux/amd64, linux/arm64)
- Cache otimizado do Docker Buildx
- Publicação automática no ghcr.io
- Tags automáticas baseadas em eventos Git

**Tags Geradas:**
```bash
# Branch main
ghcr.io/tech-preta/alderaan-api:main
ghcr.io/tech-preta/alderaan-api:latest

# Branch develop
ghcr.io/tech-preta/alderaan-api:develop

# Tags de versão (ex: v1.0.0)
ghcr.io/tech-preta/alderaan-api:v1.0.0
ghcr.io/tech-preta/alderaan-api:1.0

# Pull requests
ghcr.io/tech-preta/alderaan-api:pr-123
```

**Uso:**
```bash
# Baixar última versão
docker pull ghcr.io/tech-preta/alderaan-api:latest

# Baixar versão específica
docker pull ghcr.io/tech-preta/alderaan-api:v1.0.0

# Executar
docker run -p 8080:8080 ghcr.io/tech-preta/alderaan-api:latest
```

## 📦 Helm Chart Pipeline

### Workflow: `.github/workflows/helm-publish.yml`

**Triggers:**
- Push em branch `main`
- Push de tags `v*`
- Trigger manual (workflow_dispatch)

**Funcionalidades:**
- Empacotamento do Helm chart
- Publicação no OCI registry (ghcr.io)
- Versionamento automático baseado no Chart.yaml

**Instalação:**
```bash
# Instalar diretamente do registry
helm install alderaan oci://ghcr.io/tech-preta/helm-charts/alderaan \
  --version 1.0.0 \
  --namespace alderaan \
  --create-namespace

# Listar versões disponíveis
helm search repo alderaan --versions
```

## 🚀 Release Pipeline

### Workflow: `.github/workflows/release.yml`

**Triggers:**
- Push de tags `v*` (ex: v1.0.0, v2.1.0-beta)

**Funcionalidades:**
- Geração automática de changelog
- Criação de release no GitHub
- Instruções de instalação automáticas
- Detecção de pre-releases (alpha, beta, rc)

**Como Criar uma Release:**

```bash
# 1. Atualizar versão no Chart.yaml
sed -i 's/version: .*/version: 1.0.1/' charts/alderaan/Chart.yaml

# 2. Commit e push
git add charts/alderaan/Chart.yaml
git commit -m "chore: bump version to 1.0.1"
git push

# 3. Criar e push tag
git tag v1.0.1
git push origin v1.0.1

# O pipeline criará automaticamente:
# - Release no GitHub
# - Imagem Docker
# - Helm chart
```

## 🔐 Secrets Necessários

O pipeline usa o seguinte secret:

### USER_TOKEN
- **Tipo:** Personal Access Token (PAT)
- **Permissões Necessárias:**
  - `write:packages` - Para publicar no GitHub Container Registry
  - `read:packages` - Para ler pacotes
  - `repo` - Para criar releases

**Como Configurar:**

1. Vá em GitHub → Settings → Developer settings → Personal access tokens
2. Crie um novo token com as permissões acima
3. Adicione o token nos Secrets do repositório:
   - Repository → Settings → Secrets and variables → Actions
   - New repository secret: `USER_TOKEN`

## 📊 Monitoramento dos Workflows

### Status dos Workflows

Acesse: `https://github.com/Tech-Preta/alderaan/actions`

### Badges para README

```markdown
![Docker Build](https://github.com/Tech-Preta/alderaan/actions/workflows/docker-build.yml/badge.svg)
![Helm Publish](https://github.com/Tech-Preta/alderaan/actions/workflows/helm-publish.yml/badge.svg)
![Release](https://github.com/Tech-Preta/alderaan/actions/workflows/release.yml/badge.svg)
```

## 🔍 Verificação de Imagens Publicadas

### Docker Images

```bash
# Listar tags disponíveis
gh api /user/packages/container/alderaan-api/versions | jq '.[].metadata.container.tags[]'

# Inspecionar imagem
docker pull ghcr.io/tech-preta/alderaan-api:latest
docker inspect ghcr.io/tech-preta/alderaan-api:latest
```

### Helm Charts

```bash
# Mostrar informações do chart
helm show chart oci://ghcr.io/tech-preta/helm-charts/alderaan --version 1.0.0

# Mostrar valores padrão
helm show values oci://ghcr.io/tech-preta/helm-charts/alderaan --version 1.0.0
```

## 🐛 Troubleshooting

### Erro: Permission denied

```
Error: failed to authorize: failed to fetch anonymous token: 
unexpected status: 401 Unauthorized
```

**Solução:** Verificar se o token `USER_TOKEN` tem as permissões corretas.

### Erro: Package already exists

```
Error: package already exists
```

**Solução:** Incrementar a versão no `Chart.yaml` ou deletar a versão existente no GitHub Packages.

### Build Multi-arch Falha

```
Error: failed to solve: multiple platforms feature is currently not supported
```

**Solução:** O workflow usa Docker Buildx que suporta multi-platform. Verificar se o runner tem Buildx habilitado.

## 📚 Recursos Adicionais

- [GitHub Container Registry Docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
- [GitHub Packages Helm Docs](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-helm-package-registry)
- [Docker Buildx](https://docs.docker.com/buildx/working-with-buildx/)
- [Helm OCI Registry](https://helm.sh/docs/topics/registries/)
- [GitHub Actions](https://docs.github.com/en/actions)

## 🎯 Boas Práticas

1. **Versionamento Semântico**: Use [Semantic Versioning](https://semver.org/)
   - `MAJOR.MINOR.PATCH` (ex: 1.0.0)
   - Incremente MAJOR para breaking changes
   - Incremente MINOR para novas features
   - Incremente PATCH para bug fixes

2. **Tags de Pre-release**: Use sufixos para pre-releases
   - `v1.0.0-alpha.1`
   - `v1.0.0-beta.2`
   - `v1.0.0-rc.1`

3. **Changelog**: Mantenha mensagens de commit claras
   - Use conventional commits: `feat:`, `fix:`, `chore:`
   - O changelog é gerado automaticamente dos commits

4. **Testing**: Teste localmente antes de criar release
   ```bash
   # Build local
   docker build -t alderaan-api:test .
   docker run --rm alderaan-api:test
   
   # Lint helm chart
   helm lint charts/alderaan
   ```

5. **Security**: 
   - Nunca commite secrets
   - Use o scanner Trivy integrado
   - Revise dependências regularmente

---

**Última atualização:** 2024-10-05
