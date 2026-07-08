# Plano de Refatoração e Solução de Problemas

## 1. Problemas Identificados (Erros e Bugs)

Durante a varredura do projeto, identificamos diversos erros de compilação, problemas sintáticos e imports incorretos que quebraram a aplicação. A lista de problemas é detalhada abaixo:

### 1.1 Erros Sintáticos e de Compilação
- **Vírgulas em blocos de importação:** As importações múltiplas em Go não utilizam vírgulas (`,`) para separar os pacotes, mas apenas quebras de linha. Este erro ocorre em vários arquivos:
  - `app/core/chaincode/ChaincodeQuery.go`
  - `app/core/chaincode/ChaincodeService.go`
  - `app/domain/metadata/MetadataDTO.go`
- **Sintaxe de inicialização e chamadas de métodos:**
  - `app/domain/metadata/MetadataService.go`: Existem vírgulas sobrando, faltam chaves ou vírgulas em "composite literals", e existe código (como condicionais `if`) escrito fora do escopo de funções (no nível do pacote).
- **Incompatibilidade de structs e métodos:**
  - `app/core/chaincode/ChaincodeQuery.go`: Métodos recebem o "receiver" `(c *Contract)` ou `(chaincodeConnection *ChaincodeConnection)`, mas deveriam estar associados à struct definida `ChaincodeQuery`. Dentro do método, variáveis incompatíveis são chamadas (ex: uso de `c.ChaincodeService` dentro de um método que recebe `chaincodeConnection`).

### 1.2 Pacotes e Importações Quebradas
- Diversas importações assumem diretórios raiz, mas na verdade eles estão dentro da pasta `app/` ou não existem:
  - `github.com/silvasilas99/entruster/domain/metadata` deveria ser `github.com/silvasilas99/entruster/app/domain/metadata`.
  - `github.com/silvasilas99/entruster/elasticsearch` deveria ser `github.com/silvasilas99/entruster/app/core/elasticsearch`.
  - `github.com/silvasilas99/entruster/audit`: O pacote de auditoria foi referenciado na inicialização, mas não existe e precisará ser implementado.
- Múltiplos nomes de pacote na mesma pasta: Em `app/core/elasticsearch/`, existem arquivos com declaração `package elastic_query` e `package elastic_service`. O Go exige que todos os arquivos da mesma pasta tenham o mesmo nome de pacote (ex: `package elasticsearch`).

## 2. Tarefas e Plano de Execução Passo a Passo

Para restaurar o funcionamento do projeto, adicionaremos todas as dependências e arquivos de rede necessários, seguindo o passo a passo de forma estrita.

### Passo 1: Correção de Erros de Sintaxe e Pacotes
- Remover as vírgulas após as strings nos blocos `import (...)` em `ChaincodeQuery.go`, `ChaincodeService.go` e `MetadataDTO.go`.
- Corrigir `MetadataService.go`, envolvendo trechos soltos em suas respectivas funções e resolvendo as vírgulas/pontos e vírgulas (`;`) esperados.
- Em `app/core/elasticsearch/`, padronizar o nome do pacote nos arquivos para `package elasticsearch`.
- Atualizar todos os `import` no projeto (como em `routes/api_routes.go`) apontando para as pastas corretas de `app/domain/metadata` e `app/core/elasticsearch`.

### Passo 2: Implementação da Camada de Auditoria (Audit)
- Criar a pasta `app/core/audit` e os arquivos com o mock ou serviço para `AuditService` e os structs (modelos de auditoria).
- Corrigir qualquer erro de lógica relacionada ao `auditSvc` no arquivo de rotas e injeção de dependências.

### Passo 3: Criação do Script de Setup da Rede (network-setup.sh)
- Criar arquivo `./scripts/network-setup.sh`.
- Escrever os comandos shell para baixar o Fabric test-network (caso não exista na pasta `fabric-samples`), subir os containers via `./network.sh up`, criar o canal `metadatachannel` e implantar o chaincode escrito em Go em `/app/core/chaincode`.

### Passo 4: Configuração da Infraestrutura (Docker e Compose)
- Criar o manifesto de rede `./docker/docker-compose.yml` que provisionará Elasticsearch, MongoDB (se necessário), além do próprio container da API Entruster, conectando-os devidamente.
- Criar o `./docker/Dockerfile` base para compilar e iniciar a API Golang. Este arquivo compilará a aplicação via `go build` no `app/core/bootstrap/Server.go` ou em um `main.go` apropriado, e definirá o ponto de entrada.

### Passo 5: Validação Final
- Rodar `go mod tidy` para baixar/atualizar as definições de módulos perfeitamente, garantindo que não restem dependências externas desconhecidas.
- Rodar `go build ./...` e garantir 0 erros de compilação.
- Instruir sobre a execução de `docker compose up -d` para provar a estabilidade final.
