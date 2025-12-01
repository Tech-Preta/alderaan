# 📝 Guia Rápido de Commits

Guia rápido para escrever commits que geram releases automáticos.

## 🎯 Formato Básico

```
<tipo>(<escopo>): <descrição curta>

<corpo opcional>

<rodapé opcional>
```

## 🏷️ Tipos e Versões

| Tipo | Versão | Release? | Uso |
|------|--------|----------|-----|
| `feat:` | **MINOR** 1.0.0 → 1.1.0 | ✅ Sim | Nova funcionalidade |
| `fix:` | **PATCH** 1.0.0 → 1.0.1 | ✅ Sim | Correção de bug |
| `perf:` | **PATCH** | ✅ Sim | Melhoria de performance |
| `docs:` | **PATCH** | ✅ Sim | Documentação |
| `style:` | **PATCH** | ✅ Sim | Formatação (sem mudar lógica) |
| `refactor:` | **PATCH** | ✅ Sim | Refatoração (sem mudar comportamento) |
| `test:` | **PATCH** | ✅ Sim | Testes |
| `build:` | **PATCH** | ✅ Sim | Sistema de build |
| `ci:` | **PATCH** | ✅ Sim | CI/CD |
| `chore:` | - | ❌ Não | Tarefas de manutenção |

## 💥 Breaking Changes

### Forma 1: Exclamação
```bash
feat!: remove suporte a API v1
```

### Forma 2: Rodapé
```bash
feat: redesenha API para v2

BREAKING CHANGE: API v1 removida. Veja guia de migração.
```

Ambas geram **MAJOR** version: 1.0.0 → 2.0.0

## 📋 Exemplos Práticos

### ✅ BOM - Gera release 1.1.0
```bash
feat(auth): adiciona autenticação JWT

Implementa sistema de autenticação usando JWT.
Inclui tokens de acesso e refresh.

Closes #42
```

### ✅ BOM - Gera release 1.0.1
```bash
fix(api): corrige timeout em requisições grandes

Aumenta timeout para 30s em uploads.
Adiciona retry automático em falhas de rede.
```

### ✅ BOM - Gera release 2.0.0
```bash
feat!: remove suporte a Go 1.19

BREAKING CHANGE: Versão mínima do Go agora é 1.21
```

### ❌ RUIM - Não gera release
```bash
chore: atualiza .gitignore
```

### ❌ RUIM - Não segue padrão
```bash
adiciona nova feature de login
```

## 🎨 Escopos Comuns

Use escopos para organizar melhor:

```bash
feat(api): novo endpoint
feat(auth): sistema de permissões
feat(cache): cache Redis
feat(database): nova tabela
feat(monitoring): novo dashboard

fix(api): corrige validação
fix(auth): corrige token expirado
fix(handlers): corrige memory leak

docs(readme): atualiza instalação
docs(api): adiciona exemplos

test(integration): testes E2E
test(unit): testes unitários

ci(actions): otimiza workflow
ci(docker): melhora build

build(deps): atualiza dependências
build(docker): otimiza Dockerfile
```

## 🔗 Referências a Issues

Sempre referencie issues no corpo ou rodapé:

```bash
feat(api): adiciona paginação

Implementa paginação em todos os endpoints de listagem.

Closes #123
Refs #124, #125
```

Ou use palavras-chave:
- `Closes #123` - Fecha a issue
- `Fixes #123` - Corrige a issue
- `Refs #123` - Referencia a issue
- `See #123` - Veja a issue

## 🚫 Não Fazer

### ❌ Commits vagos
```bash
fix: correção
feat: nova feature
chore: updates
```

### ❌ Sem tipo
```bash
adiciona nova funcionalidade
corrige bug
```

### ❌ Tipo errado
```bash
chore: adiciona autenticação  # Use feat:
```

### ❌ Múltiplas mudanças
```bash
feat: adiciona auth, corrige bugs, atualiza docs
# Divida em 3 commits separados!
```

## ✅ Fazer

### ✅ Commits descritivos
```bash
feat(auth): adiciona autenticação via OAuth2
fix(database): corrige deadlock em transações concorrentes
docs(api): adiciona exemplos de uso completos
```

### ✅ Commits atômicos
Um commit = uma mudança lógica

```bash
# Ao invés de:
feat: adiciona auth e corrige bugs

# Faça:
feat(auth): adiciona autenticação JWT
fix(handlers): corrige memory leak no handler
docs(auth): documenta sistema de autenticação
```

### ✅ Corpo detalhado (quando necessário)
```bash
feat(cache): implementa cache distribuído com Redis

- Adiciona cliente Redis configurável
- Implementa cache para produtos e categorias
- TTL configurável via variável de ambiente
- Fallback para banco quando cache falha

Melhora performance em 40% nos endpoints de listagem.

Closes #87
```

## 🎯 Dicas Rápidas

1. **Use o presente**: "adiciona" não "adicionado"
2. **Seja específico**: "corrige timeout na API" não "corrige bug"
3. **Primeira letra minúscula**: "feat: adiciona" não "feat: Adiciona"
4. **Sem ponto final**: "feat: nova feature" não "feat: nova feature."
5. **Máximo 50 caracteres** no título
6. **Use corpo para detalhes**, não título longo

## 🧪 Testar Antes de Commit

```bash
# 1. Verifique o que vai commitar
git diff

# 2. Rode os testes
make test

# 3. Verifique linter
make lint  # se disponível

# 4. Então faça commit
git commit -m "feat(api): adiciona novo endpoint"
```

## 📚 Recursos

- [Conventional Commits](https://www.conventionalcommits.org/pt-br/)
- [Semantic Versioning](https://semver.org/lang/pt-BR/)
- [Documentação Completa](../docs/12-automated-releases.md)

---

**💡 Dica**: Salve este arquivo como referência rápida!
