# 🚀 Releases, Tags e CHANGELOG Automáticos

Este documento descreve o processo de release automatizado do projeto Alderaan, que utiliza **Semantic Versioning** e **Conventional Commits** para gerar versões, tags, CHANGELOG e releases do GitHub automaticamente.

## 📋 Índice

- [Visão Geral](#visão-geral)
- [Conventional Commits](#conventional-commits)
- [Semantic Versioning](#semantic-versioning)
- [Como Funciona](#como-funciona)
- [Tipos de Commits](#tipos-de-commits)
- [Exemplos de Uso](#exemplos-de-uso)
- [Comandos Make](#comandos-make)
- [Troubleshooting](#troubleshooting)

## 🎯 Visão Geral

O projeto utiliza **semantic-release** integrado ao GitHub Actions para:

- ✅ **Analisar commits** desde a última release
- ✅ **Determinar o tipo de versão** (major, minor, patch)
- ✅ **Gerar notas de release** automaticamente
- ✅ **Atualizar CHANGELOG.md** com as mudanças
- ✅ **Criar tag Git** com a nova versão
- ✅ **Publicar release no GitHub** com binários
- ✅ **Atualizar arquivo VERSION** com a versão atual

## 📝 Conventional Commits

O projeto segue a especificação [Conventional Commits](https://www.conventionalcommits.org/pt-br/).

### Formato

```
<tipo>[escopo opcional]: <descrição>

[corpo opcional]

[rodapé(s) opcional(is)]
```

### Exemplo

```bash
feat(api): adiciona endpoint de autenticação JWT

Implementa autenticação JWT com refresh tokens
para melhorar a segurança da API.

Closes #123
```

## 🔢 Semantic Versioning

O projeto adere ao [Semantic Versioning 2.0.0](https://semver.org/lang/pt-BR/).

### Formato: MAJOR.MINOR.PATCH

- **MAJOR (1.x.x)**: Mudanças incompatíveis na API
- **MINOR (x.1.x)**: Novas funcionalidades (compatível)
- **PATCH (x.x.1)**: Correções de bugs (compatível)

### Regras de Versionamento

| Tipo de Commit | Versão Incrementada | Exemplo |
|----------------|---------------------|---------|
| `fix:` | PATCH | 1.0.0 → 1.0.1 |
| `feat:` | MINOR | 1.0.0 → 1.1.0 |
| `BREAKING CHANGE:` | MAJOR | 1.0.0 → 2.0.0 |
| `perf:`, `docs:`, `style:`, `refactor:`, `test:`, `build:`, `ci:` | PATCH | 1.0.0 → 1.0.1 |
| `chore:` | Nenhuma | Não gera release |

## ⚙️ Como Funciona

### 1. Desenvolvimento

Desenvolva normalmente e faça commits seguindo o padrão Conventional Commits:

```bash
git add .
git commit -m "feat: adiciona cache Redis para produtos"
git push origin main
```

### 2. Trigger Automático

Ao fazer push para a branch `main`, o GitHub Actions:

1. Executa testes
2. Compila a aplicação
3. Analisa commits desde a última release
4. Determina a nova versão
5. Atualiza CHANGELOG.md
6. Cria tag Git
7. Publica release no GitHub

### 3. Release Criado

O sistema automaticamente:

- ✅ Cria release no GitHub
- ✅ Anexa binário compilado
- ✅ Gera notas de release
- ✅ Atualiza CHANGELOG.md
- ✅ Comita mudanças com `[skip ci]`

## 🏷️ Tipos de Commits

### 🎯 Principais (Geram Release)

#### `feat:` - Nova Funcionalidade (MINOR)

```bash
feat: adiciona suporte a autenticação OAuth2
feat(api): implementa paginação em endpoints de listagem
feat(monitoring): adiciona dashboard de custos no Grafana
```

#### `fix:` - Correção de Bug (PATCH)

```bash
fix: corrige vazamento de memória no handler de produtos
fix(database): resolve deadlock em transações concorrentes
fix(api): corrige validação de campos obrigatórios
```

#### `perf:` - Melhoria de Performance (PATCH)

```bash
perf: otimiza queries do PostgreSQL com índices
perf(cache): implementa cache em memória para produtos
```

### 🔧 Secundários (Geram Release)

#### `docs:` - Documentação (PATCH)

```bash
docs: atualiza README com instruções de instalação
docs(api): adiciona exemplos de uso da API REST
```

#### `style:` - Formatação (PATCH)

```bash
style: formata código com gofmt
style: aplica linter em todo o projeto
```

#### `refactor:` - Refatoração (PATCH)

```bash
refactor: simplifica lógica de validação de produtos
refactor(domain): extrai interface do repositório
```

#### `test:` - Testes (PATCH)

```bash
test: adiciona testes para handler de produtos
test(integration): implementa testes E2E da API
```

#### `build:` - Build (PATCH)

```bash
build: atualiza dependências do Go
build(docker): otimiza Dockerfile multi-stage
```

#### `ci:` - CI/CD (PATCH)

```bash
ci: adiciona workflow de security scan
ci(actions): otimiza cache de dependências
```

### 🚫 Não Geram Release

#### `chore:` - Manutenção

```bash
chore: atualiza .gitignore
chore: remove código comentado
```

### 💥 Breaking Changes (MAJOR)

Qualquer commit com `BREAKING CHANGE:` no rodapé:

```bash
feat(api): redesenha estrutura de endpoints

BREAKING CHANGE: endpoints de produtos movidos de /products para /api/v2/products
```

Ou com `!` após o tipo:

```bash
feat!: remove suporte a API v1
```

## 💡 Exemplos de Uso

### Exemplo 1: Correção de Bug

```bash
# 1. Faz a correção
git add .
git commit -m "fix(api): corrige timeout em requisições grandes"

# 2. Push para main
git push origin main

# 3. Aguarda ~2-3 minutos
# 4. Nova versão 1.0.1 é criada automaticamente
```

### Exemplo 2: Nova Funcionalidade

```bash
# 1. Implementa a funcionalidade
git add .
git commit -m "feat(auth): adiciona autenticação JWT com refresh tokens

Implementa sistema completo de autenticação usando JWT.
Inclui tokens de acesso e refresh para maior segurança.

Closes #42"

# 2. Push para main
git push origin main

# 3. Nova versão 1.1.0 é criada automaticamente
```

### Exemplo 3: Breaking Change

```bash
# 1. Implementa mudança incompatível
git add .
git commit -m "feat(api)!: redesenha API para versão 2.0

Reestrutura completamente os endpoints da API.
Remove suporte a API v1.

BREAKING CHANGE: API v1 removida. Clientes devem migrar para v2.
Veja guia de migração em docs/migration-v2.md"

# 2. Push para main
git push origin main

# 3. Nova versão 2.0.0 é criada automaticamente
```

### Exemplo 4: Múltiplos Commits

```bash
# Várias mudanças antes do release
git commit -m "feat: adiciona cache Redis"
git commit -m "feat: implementa rate limiting"
git commit -m "fix: corrige bug em validação"
git commit -m "docs: atualiza README"

git push origin main

# Release 1.2.0 com todas as mudanças
```

## 🛠️ Comandos Make

### `make release-check`

Verifica se o repositório está pronto para release:

```bash
make release-check
```

Verifica:
- ✅ Sem mudanças não commitadas
- ✅ Sem mudanças staged
- ✅ Branch limpo

### `make release-version`

Mostra a versão atual e última tag:

```bash
make release-version

# Saída:
# Versão atual:
# 1.2.3
# 
# Última tag:
# v1.2.3
```

### `make release-dry-run`

Simula um release (não executa):

```bash
make release-dry-run
```

## 📄 Arquivos Importantes

### `.releaserc.json`

Configuração do semantic-release:

- Define regras de versionamento
- Configura plugins
- Define formato do CHANGELOG
- Configura assets do release

### `.github/workflows/release.yml`

Workflow do GitHub Actions:

- Executa em push para `main`
- Roda testes
- Compila aplicação
- Executa semantic-release

### `VERSION`

Arquivo com a versão atual do projeto:

```
1.2.3
```

### `CHANGELOG.md`

Gerado automaticamente com todas as mudanças:

```markdown
# Changelog

## [1.2.0] - 2024-10-15

### ✨ Funcionalidades
- adiciona cache Redis para produtos
- implementa rate limiting na API

### 🐛 Correções
- corrige bug em validação de produtos
```

## 🔍 Troubleshooting

### Release não foi criado

**Problema**: Fiz push mas nenhum release foi criado.

**Soluções**:

1. Verifique se usou Conventional Commits:
   ```bash
   git log --oneline -5
   ```

2. Verifique se não é um commit `chore:`:
   ```bash
   # chore não gera release
   chore: atualiza .gitignore  ❌
   
   # Use outro tipo
   docs: atualiza .gitignore   ✅
   ```

3. Verifique os logs do GitHub Actions:
   - Vá em `Actions` no GitHub
   - Clique no workflow `Release`
   - Veja os logs de execução

### Versão errada foi gerada

**Problema**: Esperava MINOR mas gerou PATCH.

**Causa**: Tipo de commit incorreto.

**Solução**: Use o tipo correto:
- `feat:` → MINOR (1.0.0 → 1.1.0)
- `fix:` → PATCH (1.0.0 → 1.0.1)
- `BREAKING CHANGE:` → MAJOR (1.0.0 → 2.0.0)

### CHANGELOG não foi atualizado

**Problema**: Release criado mas CHANGELOG não mudou.

**Solução**:
1. Verifique se o commit está no formato correto
2. Veja logs do workflow
3. Verifique permissões do token GitHub

### Como reverter um release

```bash
# 1. Deletar tag local
git tag -d v1.2.3

# 2. Deletar tag remota
git push origin :refs/tags/v1.2.3

# 3. Deletar release no GitHub (interface web)
# Settings → Releases → Delete release

# 4. Reverter commit se necessário
git revert HEAD
git push origin main
```

## 📚 Referências

- [Conventional Commits](https://www.conventionalcommits.org/pt-br/)
- [Semantic Versioning](https://semver.org/lang/pt-BR/)
- [semantic-release](https://semantic-release.gitbook.io/)
- [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/)

## ✅ Boas Práticas

1. **Commits atômicos**: Um commit = uma mudança lógica
2. **Mensagens descritivas**: Seja claro e objetivo
3. **Use escopo**: Identifique a área afetada
4. **Breaking changes explícitos**: Sempre documente mudanças incompatíveis
5. **Revise antes do push**: Verifique se os commits estão corretos
6. **Referências a issues**: Use `Closes #123` no corpo do commit
7. **Testes antes do push**: Garanta que os testes passam
8. **Branch main protegida**: Configure branch protection rules

## 🎯 Workflow Recomendado

```bash
# 1. Criar branch de feature
git checkout -b feature/nova-funcionalidade

# 2. Desenvolver e fazer commits
git add .
git commit -m "feat: implementa funcionalidade X"

# 3. Fazer mais commits se necessário
git commit -m "test: adiciona testes para funcionalidade X"
git commit -m "docs: documenta funcionalidade X"

# 4. Criar Pull Request
gh pr create --title "feat: Nova Funcionalidade X" --body "Descrição detalhada"

# 5. Após aprovação, merge para main
gh pr merge --squash

# 6. Release automático é criado! 🎉
```

## 🔐 Segurança

O workflow de release:
- ✅ Usa tokens limitados ao escopo
- ✅ Hardening do runner com step-security
- ✅ Executa testes antes do release
- ✅ Valida binários compilados
- ✅ Usa versões pinadas de actions
- ✅ Auditoria de egress policy

---

**Nota**: Este processo foi configurado para simplificar o workflow de release e garantir versionamento consistente. Para sugestões ou problemas, abra uma issue no GitHub.
