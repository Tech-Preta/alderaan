# 🤖 GitHub Automation & Workflows

Este diretório contém toda a automação e configuração do GitHub para o projeto Alderaan.

## 📋 Índice

- [Workflows](#workflows)
- [Documentação de Releases](#documentação-de-releases)
- [Templates](#templates)
- [Configurações](#configurações)

## 🔄 Workflows

### Release Automático

**Arquivo**: `workflows/release.yml`

Workflow principal que automatiza todo o processo de release:

```yaml
Trigger: Push para main (ou manual)
├── 1. Harden runner (security)
├── 2. Checkout code
├── 3. Setup Node.js 20
├── 4. Install semantic-release
├── 5. Setup Go 1.24
├── 6. Cache Go modules
├── 7. Install swag (Swagger)
├── 8. Generate Swagger docs
├── 9. Run tests
├── 10. Build application
└── 11. Execute semantic-release
     ├── Analisa commits
     ├── Determina versão
     ├── Atualiza CHANGELOG.md
     ├── Atualiza VERSION
     ├── Cria tag Git
     ├── Publica release no GitHub
     └── Anexa binário compilado
```

**Quando executa**: 
- Push para `main`
- Manualmente via workflow_dispatch

**O que faz**:
- ✅ Roda testes completos
- ✅ Gera documentação Swagger
- ✅ Compila binário
- ✅ Cria release automático
- ✅ Atualiza CHANGELOG
- ✅ Cria tags Git
- ✅ Publica no GitHub

### Outros Workflows

- **`dependency-review.yml`** - Revisa dependências em PRs
- **`issue.yml`** - Gerencia issues
- **`osv-scanner.yml`** - Scanner de vulnerabilidades OSV
- **`scorecard.yml`** - Supply-chain security scorecard
- **`trivy.yml`** - Gera SBOM com Trivy
- **`dependabot.yml`** - Configuração do Dependabot

## 📚 Documentação de Releases

### 1. Setup Guide (`RELEASE_SETUP.md`)

**Para quem**: Mantenedores do projeto
**Conteúdo**: 
- ✅ Checklist de configuração
- ⚙️ Permissões necessárias no GitHub
- 🛡️ Branch protection recomendada
- 🧪 Como testar o workflow
- 🐛 Troubleshooting específico

**Quando usar**: Ao configurar o projeto pela primeira vez

### 2. Commit Guide (`COMMIT_GUIDE.md`)

**Para quem**: Todos os desenvolvedores
**Conteúdo**:
- 📝 Formato de commits
- 🏷️ Tabela de tipos e versões
- ✅ Exemplos do que fazer
- ❌ Exemplos do que NÃO fazer
- 💡 Dicas rápidas

**Quando usar**: Antes de fazer qualquer commit

### 3. CHANGELOG Example (`CHANGELOG_EXAMPLE.md`)

**Para quem**: Curiosos e novos contribuidores
**Conteúdo**:
- 📖 Exemplo completo de CHANGELOG
- 🏗️ Estrutura das seções
- 🔗 Como links são gerados
- 📊 Organização por tipo

**Quando usar**: Para entender como ficará o CHANGELOG

### 4. Documentação Completa (`../docs/12-automated-releases.md`)

**Para quem**: Qualquer pessoa interessada
**Conteúdo**: 10KB+ de documentação completa
- 🎯 Visão geral do sistema
- 📝 Conventional Commits detalhado
- 🔢 Semantic Versioning explicado
- ⚙️ Como funciona internamente
- 💡 Exemplos práticos extensivos
- 🐛 Troubleshooting completo
- ✅ Boas práticas

**Quando usar**: Para entender profundamente o sistema

## 🎯 Guia Rápido por Persona

### 👨‍💻 Desenvolvedor Novo

1. Leia: [`COMMIT_GUIDE.md`](COMMIT_GUIDE.md) ⚡
2. Faça commits usando o formato
3. Pronto! 🎉

### 🔧 Mantenedor Configurando Projeto

1. Leia: [`RELEASE_SETUP.md`](RELEASE_SETUP.md)
2. Configure permissões no GitHub
3. Opcionalmente configure branch protection
4. Teste o workflow
5. Pronto! ✅

### 🤔 Curioso sobre CHANGELOG

1. Veja: [`CHANGELOG_EXAMPLE.md`](CHANGELOG_EXAMPLE.md)
2. Compare com o CHANGELOG real depois de alguns releases
3. Aproveite! 📖

### 🎓 Quer Entender Tudo

1. Leia: [`../docs/12-automated-releases.md`](../docs/12-automated-releases.md)
2. Estude os exemplos
3. Experimente com commits de teste
4. Domine o sistema! 🚀

## 📁 Estrutura de Arquivos

```
.github/
├── workflows/
│   ├── release.yml                    # ⭐ Workflow de release
│   ├── dependency-review.yml          # Revisão de dependências
│   ├── issue.yml                      # Gerenciamento de issues
│   ├── osv-scanner.yml                # Scanner OSV
│   ├── scorecard.yml                  # Security scorecard
│   ├── trivy.yml                      # SBOM generator
│   └── dependabot.yml                 # Config Dependabot
│
├── PULL_REQUEST_TEMPLATE/             # Templates de PR
├── ISSUE_TEMPLATE/                    # Templates de issues
│
├── RELEASE_SETUP.md                   # 📋 Setup guide
├── COMMIT_GUIDE.md                    # 📝 Guia rápido de commits
├── CHANGELOG_EXAMPLE.md               # 📖 Exemplo de CHANGELOG
├── README.md                          # 📚 Este arquivo
│
├── CODEOWNERS                         # Owners do código
└── FUNDING.yml                        # Funding/sponsorship
```

## 🔐 Configuração de Segurança

Todos os workflows seguem boas práticas de segurança:

✅ **Harden Runner** - Todos os workflows usam `step-security/harden-runner`
✅ **Versões Pinadas** - Actions pinadas por hash SHA
✅ **Egress Policy** - Política de auditoria de saída
✅ **Minimal Permissions** - Permissões mínimas necessárias
✅ **Tokens Scoped** - GITHUB_TOKEN com escopo limitado

## 🎯 Semantic Versioning

O projeto usa Semantic Versioning (SemVer) automaticamente:

```
MAJOR.MINOR.PATCH
  2  .  1  .  3

MAJOR: Breaking changes (incompatível)
MINOR: Novas funcionalidades (compatível)
PATCH: Correções de bugs (compatível)
```

### Como Commits Afetam Versão

| Commit | Antes | Depois | Tipo |
|--------|-------|--------|------|
| `feat: nova feature` | 1.0.0 | 1.1.0 | MINOR |
| `fix: corrige bug` | 1.0.0 | 1.0.1 | PATCH |
| `feat!: breaking` | 1.0.0 | 2.0.0 | MAJOR |
| `perf: otimiza` | 1.0.0 | 1.0.1 | PATCH |
| `docs: atualiza` | 1.0.0 | 1.0.1 | PATCH |
| `chore: mantém` | 1.0.0 | 1.0.0 | Nenhum |

## 📊 Estatísticas

- **Total de workflows**: 7
- **Documentação**: 1300+ linhas
- **Guias criados**: 4
- **Exemplos**: 30+
- **Tipos de commit**: 11

## 🚀 Releases no GitHub

Cada release publicado contém:

1. **Tag Git** - Versão semântica (ex: `v1.2.0`)
2. **Release Notes** - Geradas automaticamente em português
3. **Binário** - `alderaan-{version}-linux-amd64`
4. **CHANGELOG** - Atualizado automaticamente
5. **Links** - Para commits e comparações
6. **Seções** - Organizadas por tipo de mudança

## 🔗 Links Úteis

- [Conventional Commits](https://www.conventionalcommits.org/pt-br/)
- [Semantic Versioning](https://semver.org/lang/pt-BR/)
- [semantic-release](https://semantic-release.gitbook.io/)
- [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/)
- [GitHub Actions](https://docs.github.com/en/actions)

## ❓ FAQ

### Por que semantic-release e não GoReleaser?

**semantic-release** oferece:
- ✅ Análise automática de commits
- ✅ Versionamento 100% automático
- ✅ CHANGELOG automático em português
- ✅ Configuração mais simples
- ✅ Melhor para fluxo baseado em commits

**GoReleaser** é ótimo para:
- 📦 Builds multi-plataforma
- 🐳 Publicação em múltiplos registries
- 📝 Release notes manuais
- 🔧 Configuração complexa de build

Para este projeto, semantic-release + simple build é mais adequado.

### Por que não usar tags manualmente?

Tags manuais são propensas a erros:
- ❌ Esquecer de criar a tag
- ❌ Usar versão errada
- ❌ Esquecer de atualizar CHANGELOG
- ❌ Inconsistência no formato

Automação garante:
- ✅ Sempre consistente
- ✅ Sempre correto
- ✅ Sempre completo
- ✅ Economiza tempo

### Posso fazer release manual?

Sim! Use o workflow_dispatch:

```bash
# Via GitHub UI
Actions → Release → Run workflow → Run

# Ou via CLI
gh workflow run release.yml
```

### Como reverter um release?

```bash
# 1. Deletar tag local e remota
git tag -d v1.2.3
git push origin :refs/tags/v1.2.3

# 2. Deletar release no GitHub (UI)

# 3. Reverter commit se necessário
git revert <commit-hash>
git push origin main
```

## 💬 Suporte

- 📖 Leia a documentação
- 🐛 Abra uma issue
- 💬 Discuta no PR
- 📧 Contate os mantenedores

---

**Última atualização**: 2024-10-06
**Versão do sistema**: 1.0.0
**Status**: ✅ Ativo e funcional
