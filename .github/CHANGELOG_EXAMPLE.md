# Exemplo de CHANGELOG Gerado Automaticamente

Este arquivo mostra como o `CHANGELOG.md` será atualizado automaticamente pelo sistema de releases.

## Como o CHANGELOG é Gerado

Quando você faz commits seguindo [Conventional Commits](https://www.conventionalcommits.org/pt-br/) e faz push para `main`, o sistema:

1. Analisa todos os commits desde a última release
2. Agrupa os commits por tipo
3. Gera uma nova seção no CHANGELOG.md
4. Comita as mudanças automaticamente

## Exemplo: CHANGELOG Após Vários Releases

```markdown
# Changelog

Todas as mudanças notáveis neste projeto serão documentadas neste arquivo.

O formato é baseado em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/),
e este projeto adere ao [Semantic Versioning](https://semver.org/lang/pt-BR/).

## [2.0.0](https://github.com/Tech-Preta/alderaan/compare/v1.2.0...v2.0.0) (2024-10-10)

### ⚠ BREAKING CHANGES

* API v1 removida. Clientes devem migrar para v2.
* Estrutura de endpoints alterada

### ✨ Funcionalidades

* **api:** redesenha estrutura de endpoints para v2 ([abc1234](https://github.com/Tech-Preta/alderaan/commit/abc1234))
* **auth:** adiciona suporte a OAuth2 ([def5678](https://github.com/Tech-Preta/alderaan/commit/def5678))
* **cache:** implementa cache distribuído com Redis ([ghi9012](https://github.com/Tech-Preta/alderaan/commit/ghi9012))

### 🐛 Correções

* **database:** corrige deadlock em transações concorrentes ([jkl3456](https://github.com/Tech-Preta/alderaan/commit/jkl3456))
* **api:** corrige validação de campos obrigatórios ([mno7890](https://github.com/Tech-Preta/alderaan/commit/mno7890))

### ⚡ Performance

* **queries:** otimiza consultas PostgreSQL com índices ([pqr1234](https://github.com/Tech-Preta/alderaan/commit/pqr1234))

---

## [1.2.0](https://github.com/Tech-Preta/alderaan/compare/v1.1.0...v1.2.0) (2024-10-08)

### ✨ Funcionalidades

* **monitoring:** adiciona dashboard de custos no Grafana ([stu5678](https://github.com/Tech-Preta/alderaan/commit/stu5678))
* **api:** implementa paginação em endpoints de listagem ([vwx9012](https://github.com/Tech-Preta/alderaan/commit/vwx9012))
* **notifications:** adiciona sistema de notificações por email ([yza3456](https://github.com/Tech-Preta/alderaan/commit/yza3456))

### 🐛 Correções

* **handlers:** corrige vazamento de memória no handler de produtos ([bcd7890](https://github.com/Tech-Preta/alderaan/commit/bcd7890))
* **middleware:** corrige CORS para aceitar novos domínios ([efg1234](https://github.com/Tech-Preta/alderaan/commit/efg1234))

### 📚 Documentação

* **api:** adiciona exemplos de uso da API REST ([hij5678](https://github.com/Tech-Preta/alderaan/commit/hij5678))
* atualiza README com instruções de instalação ([klm9012](https://github.com/Tech-Preta/alderaan/commit/klm9012))

---

## [1.1.0](https://github.com/Tech-Preta/alderaan/compare/v1.0.0...v1.1.0) (2024-10-06)

### ✨ Funcionalidades

* **auth:** adiciona autenticação JWT ([nop3456](https://github.com/Tech-Preta/alderaan/commit/nop3456))
* **api:** implementa rate limiting ([qrs7890](https://github.com/Tech-Preta/alderaan/commit/qrs7890))

### 🐛 Correções

* **api:** corrige timeout em requisições grandes ([tuv1234](https://github.com/Tech-Preta/alderaan/commit/tuv1234))
* **database:** corrige pool de conexões ([wxy5678](https://github.com/Tech-Preta/alderaan/commit/wxy5678))

### ✅ Testes

* adiciona testes de integração E2E ([zab9012](https://github.com/Tech-Preta/alderaan/commit/zab9012))
* **handlers:** adiciona testes para validação de produtos ([cde3456](https://github.com/Tech-Preta/alderaan/commit/cde3456))

### 👷 CI/CD

* adiciona workflow de security scan ([fgh7890](https://github.com/Tech-Preta/alderaan/commit/fgh7890))
* otimiza cache de dependências ([ijk1234](https://github.com/Tech-Preta/alderaan/commit/ijk1234))

---

## [1.0.1](https://github.com/Tech-Preta/alderaan/compare/v1.0.0...v1.0.1) (2024-10-05)

### 🐛 Correções

* **docker:** corrige healthcheck no container ([lmn5678](https://github.com/Tech-Preta/alderaan/commit/lmn5678))
* **api:** corrige resposta de erro 500 ([opq9012](https://github.com/Tech-Preta/alderaan/commit/opq9012))

### 📚 Documentação

* corrige typos no README ([rst3456](https://github.com/Tech-Preta/alderaan/commit/rst3456))

---

## [1.0.0] - 2024-10-04

### 🎉 Lançamento Inicial

Primeira versão completa da API Alderaan com Domain-Driven Design, Clean Architecture, e observabilidade completa.

### ✨ Adicionado

#### **Arquitetura e Design**
- Implementação completa de Domain-Driven Design (DDD)
- Clean Architecture com separação clara de camadas
- Event Dispatcher para comunicação desacoplada
- Repository Pattern com thread-safety

#### **API RESTful**
- Framework Gin para alta performance HTTP
- Endpoints RESTful completos
- Health check endpoint
- Validação de entrada de dados

#### **Observabilidade e Monitoramento**
- Prometheus para coleta de métricas
- Grafana para visualização
- Alertmanager para alertas
- Golden Signals implementados

...

[2.0.0]: https://github.com/Tech-Preta/alderaan/compare/v1.2.0...v2.0.0
[1.2.0]: https://github.com/Tech-Preta/alderaan/compare/v1.1.0...v1.2.0
[1.1.0]: https://github.com/Tech-Preta/alderaan/compare/v1.0.0...v1.1.0
[1.0.1]: https://github.com/Tech-Preta/alderaan/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/Tech-Preta/alderaan/releases/tag/v1.0.0
```

## Estrutura de Cada Release

Cada nova release no CHANGELOG contém:

### 1. **Cabeçalho**
```markdown
## [1.1.0](link-comparação) (2024-10-06)
```
- Número da versão semântica
- Link para comparação de mudanças no GitHub
- Data do release

### 2. **Breaking Changes** (se houver)
```markdown
### ⚠ BREAKING CHANGES

* Descrição da mudança incompatível
```

### 3. **Seções Organizadas por Tipo**

Cada tipo de commit aparece em sua própria seção:

- **✨ Funcionalidades** (`feat:`) - Novas funcionalidades
- **🐛 Correções** (`fix:`) - Correções de bugs
- **⚡ Performance** (`perf:`) - Melhorias de performance
- **⏪ Reversões** (`revert:`) - Reversões de commits
- **📚 Documentação** (`docs:`) - Mudanças em documentação
- **💄 Estilo** (`style:`) - Mudanças de formatação
- **♻️ Refatoração** (`refactor:`) - Refatoração de código
- **✅ Testes** (`test:`) - Adição/modificação de testes
- **🏗️ Build** (`build:`) - Mudanças no sistema de build
- **👷 CI/CD** (`ci:`) - Mudanças em CI/CD

### 4. **Detalhes dos Commits**

Cada commit é listado com:
```markdown
* **escopo:** descrição ([hash](link-commit))
```
- Escopo opcional (ex: `api`, `auth`, `database`)
- Descrição do commit
- Hash e link para o commit no GitHub

### 5. **Links de Comparação**

No final do CHANGELOG:
```markdown
[1.1.0]: https://github.com/Tech-Preta/alderaan/compare/v1.0.0...v1.1.0
```

## Benefícios

1. **📖 Histórico Claro**: Qualquer pessoa pode ver exatamente o que mudou em cada versão
2. **🔗 Links Diretos**: Cada commit tem link para ver o código no GitHub
3. **🏷️ Categorização**: Mudanças organizadas por tipo
4. **⚠️ Breaking Changes**: Mudanças incompatíveis destacadas claramente
5. **📅 Cronologia**: Data de cada release registrada
6. **🔄 Comparações**: Links para diff entre versões

## Notas

- O CHANGELOG é **gerado automaticamente** - não edite manualmente
- Commits do tipo `chore:` **não aparecem** no CHANGELOG
- Use escopos para melhor organização (ex: `feat(api):`, `fix(database):`)
- Breaking changes são sempre destacados no topo
- Links são gerados automaticamente para commits e comparações
