# 🔧 Configuração do Sistema de Releases

Este documento descreve as configurações necessárias para o funcionamento completo do sistema de releases automáticos.

## ✅ Checklist de Configuração

### 1. Configurações Básicas (✅ Já Configuradas)

- [x] Workflow `.github/workflows/release.yml` criado
- [x] Configuração `.releaserc.json` criada
- [x] Arquivo `VERSION` inicializado
- [x] Documentação em `docs/12-automated-releases.md`
- [x] Comandos Make adicionados

### 2. Configurações do Repositório GitHub (⚠️ Necessário Ação Manual)

Para que o workflow funcione corretamente, configure o seguinte no GitHub:

#### 2.1 Permissões do Workflow

1. Vá em `Settings` → `Actions` → `General`
2. Em **Workflow permissions**, selecione:
   - ✅ **Read and write permissions**
   - ✅ **Allow GitHub Actions to create and approve pull requests**
3. Clique em **Save**

#### 2.2 Branch Protection (Recomendado)

Para evitar pushes diretos na `main`:

1. Vá em `Settings` → `Branches`
2. Adicione uma **Branch protection rule** para `main`:
   - ✅ **Require a pull request before merging**
   - ✅ **Require approvals**: 1
   - ✅ **Require status checks to pass before merging**
   - ✅ **Require branches to be up to date before merging**
   - ✅ **Do not allow bypassing the above settings** (mesmo para admins)

#### 2.3 Configuração de Tags (Opcional, mas Recomendado)

Para proteger tags de serem deletadas:

1. Vá em `Settings` → `Tags`
2. Adicione uma **Tag protection rule**:
   - Pattern: `v*`
   - ✅ Isso previne deleção acidental de tags de versão

### 3. Configurações Opcionais

#### 3.1 Notificações de Release

Para notificar a equipe sobre novos releases:

1. Vá em `Settings` → `Webhooks`
2. Configure webhook para seu canal de comunicação (Slack, Discord, etc.)
3. Eventos: Selecione **Releases**

#### 3.2 Release Drafter (Alternativa)

Se preferir revisar releases antes de publicar, considere usar [Release Drafter](https://github.com/release-drafter/release-drafter) em vez de publicação automática.

## 🧪 Como Testar

### Teste Local (Dry-run)

```bash
# Verifica se está pronto para release
make release-check

# Mostra versão atual
make release-version

# Simula release (sem executar)
make release-dry-run
```

### Teste Real (Primeiro Release)

Para testar o workflow real:

1. **Crie um commit de teste**:
   ```bash
   git checkout -b test-release
   echo "# Test" >> test.md
   git add test.md
   git commit -m "feat: teste do sistema de releases"
   ```

2. **Crie um Pull Request**:
   ```bash
   gh pr create --title "feat: teste de releases" --body "PR para testar o workflow de releases"
   ```

3. **Após aprovação, faça merge para main**

4. **Verifique o workflow**:
   - Vá em `Actions` no GitHub
   - Procure pelo workflow `Release`
   - Aguarde ~2-3 minutos
   - Verifique se:
     - ✅ Workflow executou com sucesso
     - ✅ CHANGELOG.md foi atualizado
     - ✅ Tag foi criada (ex: `v1.1.0`)
     - ✅ Release foi publicado no GitHub
     - ✅ Binário foi anexado ao release

## 🐛 Troubleshooting

### "Permission denied" no workflow

**Solução**: Configure permissões conforme item 2.1 acima.

### "Resource not accessible by integration"

**Solução**: O token GITHUB_TOKEN precisa de permissão de escrita. Veja item 2.1.

### Release não foi criado

**Possíveis causas**:
1. Commit não segue Conventional Commits
2. Commit é do tipo `chore:` (não gera release)
3. Não há mudanças desde a última release
4. Workflow falhou (veja logs em Actions)

### Tag não foi criada

**Solução**: Verifique permissões do workflow e logs em Actions.

## 📋 Formato de Commits Aceito

O workflow analisa commits usando estes padrões:

| Tipo | Versão | Exemplo |
|------|--------|---------|
| `feat:` | MINOR (1.0.0 → 1.1.0) | `feat: adiciona autenticação` |
| `fix:` | PATCH (1.0.0 → 1.0.1) | `fix: corrige bug no login` |
| `perf:` | PATCH | `perf: otimiza queries` |
| `docs:` | PATCH | `docs: atualiza README` |
| `style:` | PATCH | `style: formata código` |
| `refactor:` | PATCH | `refactor: simplifica handler` |
| `test:` | PATCH | `test: adiciona testes` |
| `build:` | PATCH | `build: atualiza deps` |
| `ci:` | PATCH | `ci: otimiza workflow` |
| `chore:` | Nenhuma | `chore: atualiza .gitignore` |
| `feat!:` ou `BREAKING CHANGE:` | MAJOR (1.0.0 → 2.0.0) | `feat!: remove API v1` |

## 📚 Documentação Adicional

- [Conventional Commits](https://www.conventionalcommits.org/pt-br/)
- [Semantic Versioning](https://semver.org/lang/pt-BR/)
- [semantic-release Documentation](https://semantic-release.gitbook.io/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

## 🎯 Status da Configuração

- [x] Workflow criado
- [x] Configuração do semantic-release
- [x] Documentação completa
- [ ] Permissões do workflow configuradas (manual)
- [ ] Branch protection configurada (opcional, mas recomendado)
- [ ] Tag protection configurada (opcional)
- [ ] Teste do primeiro release executado

---

**Nota**: Após configurar as permissões no GitHub (item 2), o sistema estará totalmente funcional e criará releases automaticamente a cada push para `main` com commits válidos.
