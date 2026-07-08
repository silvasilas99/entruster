# CONTEXTO GERAL
Estou fazendo uma grande refatoração no meu projeto, e com isso, muitas
funções pararam de funcionar, referencias foram quebradas, importações
estão com caminho desatualizado, e muitos outros bugs estão espalhadas
pelos pacotes Go. O meu projeto basicamente visa construir uma API para
interação com uma rede blockchain baseada em Hyperledger Fabric.
Existe uma camada de auditoria, que deve construída, utilizando como base
o esquema de historico de transação do ledger nativo do Hyperledger Fabric.
As documentações de API também deve ser refatorado. Tudo isso roda sobre
docker, que deverá utilizar a rede test-network do Fabric (que tem os
binários, imagens docker e de containers em ./fabric-samples/test-network),
porem, nomes de canais, peers e organizações apropriados para o projeto.

# FLUXO DA APLICAÇÃO:
O fluxo de execução deve ser o seguinte: ./app/core/bootstrap/Server.go
--> ./routes/api_routes.go --> ./app/domain/metadata/MetadataCrudController.go
(ou ./app/domain/metadata/MetadataSearchController.go, dependendo da rota)
que utiliza o MetadataDTO para transferir filtros e propriedades de metadata
--> ./app/domain/metadata/MetadataService.go --> ./app/core/chaincode/ChaincodeQuery.go
(que usa ./app/core/chaincode/ChaincodeService.go para conectar com rede blockchain)
e/ou ./app/core/elasticsearch/ElasticQuery.go (que usa ./app/core/elasticsearch/ElasticService.go
para conectar com o ElasticSearch).

# PROMPT:
Aja como uma equipe de desenvolvimento especializada em desenvolvimento de
software, Golang e blockchain com Hyperledger Fabric. Primeiro, obtenha o
contexto da aplicação (ignore as pastas vendor, e binários); então planeje
e execute as tarefas listadas a seguir:

# TAREFAS:

## 1: Criar arquivo ./scripts/network-setup.sh, o qual baixará os binários do fabric em
./fabric-samples (se não existir); e então criará e configurará a rede
blockchain; para em seguida realizar o deploy do chaincode para a rede
blockchain. Sugestão de comandos (checar e adaptar conforme necessidades):
```
cd fabric-samples/test-network
# Tear down any previous state (containers + volumes)
./network.sh down
# Start network with CouchDB and create the channel
./network.sh up createChannel -c metadatachannel -ca
# Deploy the chaincode (use absolute path)
./network.sh deployCC \
  -c metadatachannel \
  -ccn basic \
  -ccp /mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/core/chaincode \
  -ccl go
```
## 2: Criar um manifesto Docker em ./docker/docker-compose.yml, que configurará os
serviços do elasticsearch, mongodb, etc e configurará a rede blockchain,
utilizando a rede blockchain criada com a execução do script ./scripts/network-setup.sh;
## 3: Criar um Dockerfile em ./docker/Dockerfile onde deve ser criado o mecanismo
para configurar a rede, os serviços e tudo necessário para o projeto rodar
localmente sobre o Docker;
## 4: Tornar todo o projeto compilavel, executavel e testavel do zero ao rodar
"docker compose up -d"
## 5: Implemente autenticação utilizando JWT. A principio o sistema deve usar
um token estático, que pode ser armazenado num arquivo de configuração em
./config/config.go
## 6: Criar um arquivo BRAIN.md contendo tudo que foi aprendido, planejado e
pensado durante a sessão;
## 7: Criar uma camada de auditoria no caminho "./app/core/audit". O requisito
principal é integrar o esquema de auditoria ao mecanismo nativo de historico
de transações do Hyperledger, e indexar replicas no ElasticSearch, para utilizar
a engine de pesquisa nas operações de leitura de auditoria, que devem ser efetuadas
no arquivo "./app/core/AuditSearchController.go", e ser invocado em ./routes/api_routes.go
## 8: Criar comandos de seeding de metadata dados com principais dados de saúde HL7 FHIR,
especificamente dos perfis tipo Patient, Observation e Record, em ./app/cmg/seeding/MetadataSeeding.go,
que seja executavel ao rodar um comando Go.
### 8.1: Para o seeding na blockchain, use Hyperledger Caliper
### 8.2: O serviço do Hyperledger Caliper deve ser configurado em ./app/core/caliper/CapilerService.go,
o qual deve conter a lógica para realizar todas as operações de interação com o Caliper;
### 8.3: O serviço do HyperledgerCaliper deve interagir com a rede blockchain através do
./app/core/chaincode/ChaincodeQuery.go;
### 8.4: Todos os dados relativos a benchmarch e speedup devem ser colhidos durante execução
do comando de seeding, e descritos no arquivo ./BENCHMARK.md que deve ser criado
## 9: Criar testes de unidade para principais métodos (que façam sentido ser testados estaticamente)
## 10: Criar testes de integração para operações com ES e com o HF;


# DICA
Seja extremamente lúcido e ative as skills necesssárias
