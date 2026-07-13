# Entruster

O **Entruster** é uma API REST desenvolvida em **Go (Golang)** projetada para registrar, gerenciar e auditar metadados de saúde estruturados (inspirados no padrão **HL7 FHIR**) em uma rede blockchain **Hyperledger Fabric**. O sistema foi concebido como um projeto de Mestrado com foco em interoperabilidade de dados de saúde baseada em blockchain (Interopchain).

Para mitigar a latência nativa de leitura da blockchain e fornecer recursos avançados de filtragem e busca textual, o Entruster utiliza o **Elasticsearch** como uma réplica de leitura indexada e otimizada (CQRS simplificado) alimentada assincronamente via **Observer Pattern**. Adicionalmente, possui uma camada de auditoria que permite rastrear acessos e ações de modificação de dados, combinando a rapidez do Elasticsearch com a imutabilidade e prova histórica nativa do ledger do Hyperledger Fabric.

---

## 🏗️ Arquitetura do Sistema

A arquitetura do Entruster divide-se em camadas bem definidas, promovendo a separação de responsabilidades (Separation of Concerns) e facilitando a manutenção e a extensibilidade do sistema.

```mermaid
graph TD
    Client[Cliente / Frontend] -->|API REST com JWT| Server[Gin API Router]
    Server --> Controllers[Controllers: Metadata CRUD & Search]
    
    Controllers --> DTO[MetadataDTO / Model]
    Controllers --> MetadataService[MetadataService]
    Controllers --> AuditSvc[AuditService]
    
    MetadataService --> ChaincodeQuery[ChaincodeQuery]
    ChaincodeQuery --> FabricGateway[Hyperledger Fabric Gateway]
    FabricGateway -->|Smart Contract| Chaincode[MetadataContract Go Chaincode]
    
    MetadataService --> Observer[MetadataObserver]
    Observer -->|Indexação de Leitura| ElasticSvc[Elasticsearch Service]
    Observer -->|Log de Evento| AuditSvc
    
    AuditSvc -->|Salvar logs| ElasticSvc
    AuditSvc -->|Histórico Nativo| ChaincodeQuery
    
    subgraph Armazenamento e Ledger
        Chaincode -->|Ledger State & History| StateDB[(CouchDB / Ledger State)]
        ElasticSvc -->|Buscas Rápidas & Auditoria| ESIndex[(Indices: metadata & audit)]
    end
```

### 1. Camada de Apresentação (HTTP & API)
*   **Gin Router (`routes/api_routes.go`)**: Expõe endpoints HTTP RESTful para gerenciamento de metadados, auditoria, testes de saúde e mock de usuário.
*   **Middleware de Autenticação (`app/core/middleware/auth.go`)**: Valida um JSON Web Token (JWT) presente no cabeçalho `Authorization: Bearer <token>`.
*   **Swagger (`app/core/swagger/`)**: Documentação interativa e auto-gerada da API por meio do Swagger UI em `/swagger/index.html`.

### 2. Camada de Domínio
*   **Modelo de Metadados (`app/domain/metadata/MetadataModel.go`)**: Estrutura de dados principal que trafega entre a aplicação, a blockchain e o Elasticsearch. Contém informações essenciais de metadados de saúde, como referências ao paciente (`patient_id`), ao recurso original (`asset_id`), provas criptográficas de integridade / zero-knowledge (`category`), metadados de ciclo de vida (`created_at`, `updated_at`, `deleted_at`) e autoria (`created_by`, `updated_by`).
*   **Controllers (`MetadataCrudController.go` & `MetadataSearchController.go`)**: Mapeiam requisições HTTP para os serviços adequados e estruturam as respostas de sucesso ou erro da API.

### 3. Camada de Observação (Decoupling)
*   **Metadata Observer (`app/domain/metadata/MetadataObserver.go`)**: Implementa o padrão de projeto *Observer*. Sempre que ocorre uma operação de mutação com sucesso na blockchain (inserção, atualização ou deleção), o Observer é notificado. Ele executa duas ações secundárias cruciais de forma independente:
    1.  Indexa ou atualiza a cópia do documento no **Elasticsearch** (otimizando buscas).
    2.  Registra uma entrada de log estruturada na camada de auditoria através do `AuditService`.

### 4. Camada de Integração (Blockchain e Elasticsearch)
*   **Chaincode Service (`app/core/chaincode/ChaincodeService.go`)**: Estabelece conexão com o peer da rede Hyperledger Fabric usando o driver oficial `fabric-gateway` de Go. Configura a criptografia TLS de forma robusta importando os certificados da organização e do orderer.
*   **Chaincode Query (`app/core/chaincode/ChaincodeQuery.go`)**: Wrapper de conveniência que submete transações (`SubmitTransaction`) ou avalia consultas (`EvaluateTransaction`) no contrato inteligente.
*   **Elasticsearch Service (`app/core/elasticsearch/`)**: Controla a indexação de dados, tratamento de filtros complexos e consultas de texto completo utilizando a biblioteca Go oficial da Elastic.

### 5. Camada de Auditoria (`app/core/audit/`)
*   **Audit Service (`AuditService.go`)**: Fornece recursos para logar ações na aplicação (criação, edição, exclusão e leitura de dados) salvando-as no Elasticsearch para auditoria rápida. Fornece também o endpoint de histórico nativo da blockchain (`GetNativeHistory`), que executa a leitura da árvore de estados de transação da rede para garantir que os dados históricos não sofreram violações.

### 6. Camada de Smart Contract (Chaincode)
*   **Metadata Contract (`chaincode/main.go`)**: Contrato inteligente construído com o `fabric-contract-api-go`. Executa na própria rede descentralizada e gerencia o estado global (World State) dos metadados através das seguintes funções nativas:
    *   `InitLedger`: Inicializa o contador sequencial de ID de metadados.
    *   `RegisterMetadataOnNetwork`: Valida e armazena um novo metadado na blockchain gerando um ID incremental automático.
    *   `GetAllMetadataFromNetwork`: Retorna todos os metadados ativos (ignorando os deletados logicamente).
    *   `GetMetadataById`: Recupera um registro individual.
    *   `UpdateMetadataById`: Atualiza os campos editáveis.
    *   `DeleteMetadataById`: Efetua o **Soft Delete** do registro marcando a propriedade `deleted_at`.
    *   `GetMetadataHistoryById`: Acessa a API nativa da blockchain `GetHistoryForKey` para extrair todas as versões passadas do registro e as IDs de transação que o modificaram.

---

## 📁 Estrutura de Arquivos

```
entruster/
├── app/
│   ├── cmd/
│   │   ├── heater/                 # Lógica de aquecimento/warmup de infraestrutura
│   │   ├── rollback/               # Utilitários de rollback de transações
│   │   └── seeding/                # Script de seeding FHIR e Benchmark (MetadataSeeding.go)
│   ├── core/
│   │   ├── audit/                  # Camada de auditoria baseada em blockchain e ES
│   │   │   ├── AuditModel.go
│   │   │   └── AuditService.go
│   │   ├── bootstrap/
│   │   │   └── Server.go           # Ponto de entrada da API REST (inicialização gRPC e Gin)
│   │   ├── chaincode/              # Integração gRPC/Gateway com Hyperledger Fabric
│   │   │   ├── ChaincodeQuery.go
│   │   │   ├── ChaincodeService.go
│   │   │   └── TransferableDTO.go
│   │   ├── elasticsearch/          # Módulo de conexão e consultas ao Elasticsearch
│   │   │   ├── ElasticQuery.go
│   │   │   └── ElasticService.go
│   │   ├── middleware/             # Middlewares Gin (Autenticação JWT)
│   │   │   ├── auth.go
│   │   │   └── auth_test.go
│   │   └── swagger/                # Arquivos estáticos Swagger Docs
│   └── domain/
│       └── metadata/               # Domínio de Negócio: CRUD, Observador e Serviços
│           ├── MetadataCrudController.go
│           ├── MetadataDTO.go
│           ├── MetadataModel.go
│           ├── MetadataObserver.go
│           ├── MetadataSearchController.go
│           └── MetadataService.go
├── chaincode/                      # Código do Smart Contract (Go Chaincode)
│   ├── main.go                     # Definição e lógica de transações do contrato
│   ├── go.mod
│   └── go.sum
├── config/
│   └── config.go                   # Constantes de rede, caminhos de TLS e chaves estáticas
├── docker/
│   ├── Dockerfile                  # Build em múltiplos estágios do container da API
│   └── docker-compose.yml          # Manifesto Docker local para ES, MongoDB e Fabric Network
├── routes/
│   └── api_routes.go               # Definição de rotas do Gin e Handlers de mock
├── scripts/
│   └── network-setup.sh            # Script de automação do Test-Network do Hyperledger Fabric
├── utils/
│   └── response.go                 # Padronizador de respostas HTTP
├── docker-compose.yml              # Orquestrador principal na raiz do projeto
├── go.mod                          # Módulo Go principal
├── go.sum
├── TODO.md                         # Registro de desenvolvimento
└── BENCHMARK.md                    # Relatórios de desempenho concorrente e speedup
```

---

## 🔄 Fluxo de Dados (Exemplo: Criação de Metadado)

Abaixo, descrevemos o ciclo de vida de uma requisição de registro de dados para exemplificar a interação de todas as camadas:

1.  **Requisição HTTP POST**: O cliente envia uma requisição `POST /api/metadata/` contendo o payload do metadado e o header `Authorization: Bearer <token>`.
2.  **Validação JWT**: O middleware `JWTAuth()` intercepta a chamada, valida a assinatura e autenticação do token com base na chave simétrica em `config.go`, chama o `/api/mock/user` para recuperar detalhes do ator (Ex: Dr. Silas Silva) e injeta esse contexto na requisição.
3.  **Controller**: O `MetadataCrudController` decodifica o corpo JSON para a estrutura `MetadataModel`, acionando as rotinas `BeforeCreate` para carregar o timestamp no padrão RFC3339.
4.  **Service & Blockchain Gateway**: O controller delega a operação para o `ChaincodeQuery` (através do serviço de metadados), que desmembra os dados estruturados em strings individuais e executa o método `RegisterMetadataOnNetwork` do Smart Contract na rede Hyperledger Fabric por meio de gRPC cifrado (TLS).
5.  **Smart Contract**: O chaincode incrementa o ID sequencial na base, escreve o estado serializado na blockchain (CouchDB) e responde com o ID gerado de sucesso.
6.  **Observer Trigger**: Ao receber o ID da blockchain, o `MetadataObserver` é invocado de forma transparente:
    *   Registra um log de auditoria no índice `audit` do Elasticsearch apontando o ID do recurso criado, quem o criou e informações adicionais do payload.
    *   Indexa o documento completo no índice `metadata` do Elasticsearch, permitindo busca de texto livre ou filtragem imediata.
7.  **Retorno da API**: O controller encapsula o status em formato JSON estruturado e responde ao cliente HTTP.

---

## 🚀 Guia de Configuração Local do Ambiente

Siga as instruções abaixo para subir a rede blockchain, o Elasticsearch e a API localmente a partir do zero até a realização da primeira requisição de API.

### Pré-requisitos
*   **Docker** e **Docker Compose** instalados e rodando.
*   **Go** (versão 1.21 ou superior) instalado localmente (caso queira rodar ou compilar testes localmente fora do container).
*   **Curl** e **Bash** instalados.

---

### Passo 1: Subir a Infraestrutura pelo Docker Compose

Na raiz do projeto, execute o comando a seguir. Este comando é projetado para subir todo o ambiente de forma automática:

```bash
docker compose up --build -d
```

#### O que este comando faz nos bastidores?
1.  **Serviço `setup` (fabric-setup)**: Executa uma imagem leve do Docker CLI contendo bash, curl e Go. Ela checa a existência da pasta `fabric-samples`. Se não existir, ela baixa automaticamente o script de instalação oficial da Hyperledger Fabric, clona os exemplos (`fabric-samples`) e baixa as imagens e binários específicos da versão **2.5.4** do Fabric e **1.5.7** do Fabric CA.
2.  **Inicialização da Rede**: O container de setup chama o script `./scripts/network-setup.sh`. Esse script:
    *   Desliga qualquer estado anterior da test-network (`network.sh down`).
    *   Inicia a test-network contendo a organização `Org1` e uma CA (Certificate Authority), configurando o banco de dados de estado CouchDB (`network.sh up createChannel -c metadatachannel -ca -s couchdb`).
    *   Empacota e realiza o deploy do chaincode localizado na pasta `/chaincode` sob o nome `basic` e canal `metadatachannel`.
    *   Ao finalizar, gera um arquivo sentinela `.fabric_ready` no diretório raiz compartilhado.
3.  **Serviço `elasticsearch`**: Sobe o Elasticsearch 8.13.0 em modo de nó único (`single-node`) com a segurança X-Pack desativada para facilidade de desenvolvimento. Possui um healthcheck HTTP na porta `9200`.
4.  **Serviço `api` (entruster-api)**: Compila e executa o binário Go da API do Entruster. O container aguarda em loop até que o arquivo sentinela `/app/shared/.fabric_ready` seja criado pelo container de setup. Quando detectado, inicializa a conexão gRPC com a rede Fabric e sobe o servidor HTTP na porta `8080`.

---

### Passo 2: Verificar o Status do Ambiente

Você pode monitorar a inicialização através dos logs:

```bash
docker compose logs -f api
```

Ao terminar, você verá mensagens de sucesso indicando a conexão bem-sucedida ao Fabric e o servidor HTTP rodando:

```
Waiting for Fabric Test Network to be ready...
Fabric is ready. Starting API...
✅ Connected to Fabric — metadatachannel
🚀 Server running on http://localhost:8080
📄 Swagger UI at http://localhost:8080/swagger/index.html
```

---

### Passo 3: Obter o Token de Autenticação JWT

A API do Entruster é protegida por autenticação baseada em JWT. Para fins de desenvolvimento local, o sistema implementa uma validação de token estática definida no arquivo [config.go](file:///mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/config/config.go).

Você pode copiar a chave estática padrão diretamente do arquivo de configuração para realizar as chamadas HTTP.
O token JWT padrão pré-assinado é:
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c3JfNzc3IiwibmFtZSI6IkRyLiBTaWxhcyBTaWx2YSIsInJvbGUiOiJDaGllZiBDbGluaWNpYW4iLCJpYXQiOjE3ODIyOTc2MDB9.4M2R0_2nN0X4Y0tX6W7N8wYwS8uV2Fz3T9k_Y3Z1U_w
```

---

### Passo 4: Fazer a Primeira Requisição da API

Vamos realizar a primeira requisição criando um metadado de saúde (representando um recurso FHIR fictício de Paciente) por meio do `curl`:

#### 1. Inserir um Registro de Metadado (POST)
Submeta a chamada HTTP a seguir (substitua o token de autorização se necessário):

```bash
curl -X POST http://localhost:8080/api/metadata/ \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c3JfNzc3IiwibmFtZSI6IkRyLiBTaWxhcyBTaWx2YSIsInJvbGUiOiJDaGllZiBDbGluaWNpYW4iLCJpYXQiOjE3ODIyOTc2MDB9.4M2R0_2nN0X4Y0tX6W7N8wYwS8uV2Fz3T9k_Y3Z1U_w" \
  -H "Content-Type: application/json" \
  -d '{
    "patient_id": 1001,
    "asset_id": 2002,
    "category": "a8f5c2d6e9f8a3b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9",
    "name": "FHIR-Patient-Registration",
    "value": "Patient/1001-active-metadata",
    "version": "1.0",
    "owner": "Hospital Santa Casa",
    "rights": "Strictly Confidential",
    "terms_of_access": "Only authenticated medical staff can query"
  }'
```

**Resposta esperada (JSON):**
```json
{
  "code": 201,
  "status": "success",
  "message": "Metadata registered successfully",
  "data": {
    "id": 1
  }
}
```

#### 2. Consultar o Metadado Registrado por ID (GET)
Recupere as informações detalhadas gravadas na rede blockchain a partir do ID retornado no passo anterior:

```bash
curl -X GET http://localhost:8080/api/metadata/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c3JfNzc3IiwibmFtZSI6IkRyLiBTaWxhcyBTaWx2YSIsInJvbGUiOiJDaGllZiBDbGluaWNpYW4iLCJpYXQiOjE3ODIyOTc2MDB9.4M2R0_2nN0X4Y0tX6W7N8wYwS8uV2Fz3T9k_Y3Z1U_w"
```

---

### Passo 5: Executar o Seeding de Dados FHIR & Ferramenta de Benchmark

O projeto acompanha um script para carga de dados em massa (Seeding) contendo metadados de perfis padrão HL7 FHIR (como *Patient*, *Observation* e *Record*) que testa o desempenho da blockchain sob estresse.

Para rodar o processo de seeding e geração automática de métricas concorrentes, execute o comando Go nativo:

```bash
TEST_NETWORK_PATH=./fabric-samples/test-network go run ./app/cmd/seeding/MetadataSeeding.go -count 10 -workers 4
```

*   **`-count`**: Define o total de metadados a serem gerados para a simulação (default: 100).
*   **`-workers`**: Número de goroutines simultâneas que realizarão as submissões em lote (modo concorrente).

Este script calcula o tempo decorrido, o throughput médio em transações por segundo (TPS) para envios sequenciais e concorrentes, medindo o fator de aceleração (speedup) obtido através de paralelismo. Os resultados são exportados de forma formatada diretamente para o arquivo `BENCHMARK.md`.

---

## 🛠️ Comandos Úteis de Desenvolvimento

### Rodar os Testes Unitários
Para rodar os testes estáticos da API:
```bash
go test ./...
```

### Derrubar e Limpar a Infraestrutura Local
Para apagar todos os containers do Elasticsearch, da API e limpar completamente a rede blockchain (volumes e dados):
```bash
docker compose down -v
cd fabric-samples/test-network && ./network.sh down
```
