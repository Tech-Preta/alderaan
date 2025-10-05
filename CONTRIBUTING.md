# Contribuindo para o Projeto

Obrigado por considerar contribuir para este projeto! Abaixo estão algumas diretrizes para ajudar você a começar.

## Como Contribuir

### Reportando Problemas

Se você encontrar um bug ou tiver uma sugestão de melhoria, por favor, abra uma issue no repositório. Certifique-se de incluir o máximo de detalhes possível, incluindo passos para reproduzir o problema, se aplicável.

### Enviando Pull Requests

1. **Fork o repositório**: Clique no botão "Fork" no topo da página do repositório.
2. **Clone o seu fork**:
    ```sh
    git clone https://github.com/Tech-Preta/repository_sample
    cd repository_sample
    ```
3. **Crie uma branch para a sua feature ou correção**:
    ```sh
    git checkout -b minha-feature
    ```
4. **Faça as suas mudanças**: Adicione ou modifique o código conforme necessário.
5. **Adicione testes**: Certifique-se de que suas mudanças estão cobertas por testes.
6. **Commit suas mudanças**:
    ```sh
    git commit -m "Descrição das minhas mudanças"
    ```
7. **Envie a sua branch**:
    ```sh
    git push origin minha-feature
    ```
8. **Abra um Pull Request**: Vá até a página do repositório original e clique em "New Pull Request". Compare a sua branch com a branch `main` do repositório original e envie o Pull Request.

## CI/CD Pipeline

O projeto utiliza GitHub Actions para CI/CD automático. Quando você abrir um Pull Request:

### O que acontece automaticamente:
- ✅ **Build da imagem Docker** é executado (mas não publica)
- ✅ **Validação do Helm Chart** (se modificado)
- ✅ **Testes de segurança** (Trivy, OSV-Scanner)
- ✅ **Linting de código**

### Workflows disponíveis:

1. **docker-build.yml**: Build e push de imagens Docker
   - Triggers: push em `main`/`develop`, tags `v*`, PRs
   - Multi-arch: linux/amd64, linux/arm64
   - Publica em: ghcr.io/tech-preta/alderaan-api

2. **helm-publish.yml**: Empacotamento e publicação do Helm chart
   - Triggers: push em `main`, tags `v*`, manual
   - Publica em: ghcr.io/tech-preta/helm-charts

3. **release.yml**: Criação automática de releases
   - Triggers: tags `v*`
   - Gera changelog automático
   - Cria release no GitHub

### Criando uma Release:

Para criar uma nova release:

```bash
# 1. Atualize a versão no Chart.yaml
sed -i 's/version: .*/version: 1.0.1/' charts/alderaan/Chart.yaml

# 2. Commit e push
git add charts/alderaan/Chart.yaml
git commit -m "chore: bump version to 1.0.1"
git push

# 3. Crie e publique a tag
git tag v1.0.1
git push origin v1.0.1
```

O pipeline criará automaticamente:
- Imagem Docker com a tag correspondente
- Helm chart publicado
- Release no GitHub com changelog

### Mais informações:
Consulte a [documentação completa do CI/CD](docs/11-cicd-pipeline.md).
