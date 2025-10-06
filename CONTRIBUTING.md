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
6. **Commit suas mudanças** usando [Conventional Commits](https://www.conventionalcommits.org/pt-br/):
    ```sh
    # Para novas funcionalidades
    git commit -m "feat: adiciona suporte a autenticação OAuth"
    
    # Para correções de bugs
    git commit -m "fix: corrige vazamento de memória no handler"
    
    # Para documentação
    git commit -m "docs: atualiza README com exemplos"
    
    # Para testes
    git commit -m "test: adiciona testes para validação de produtos"
    ```
    
    📖 [**Veja o guia completo de commits →**](docs/12-automated-releases.md#-tipos-de-commits)
7. **Envie a sua branch**:
    ```sh
    git push origin minha-feature
    ```
8. **Abra um Pull Request**: Vá até a página do repositório original e clique em "New Pull Request". Compare a sua branch com a branch `main` do repositório original e envie o Pull Request.
