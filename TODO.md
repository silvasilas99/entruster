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

## 1: Mover ./chaincode/ ./app/core/chaincode/
## 2: Usar ChaincodeQuery.go invés de main.go, para remover redundância
## 3: Remover inicialização dos serviços como auditória, elasticsearch, observer, etc de dentro do api_routes.go
## 4: Criar mecanismo para exportação de metadados de forma segura
## 10: Criar testes de unidade para principais métodos (que façam sentido ser testados estaticamente)
## 11: Criar testes de integração para operações com ES e com o HF

# DICA
Seja extremamente lúcido e ative as skills necesssárias
