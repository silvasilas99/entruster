# TOSEARCH:
- SmartContract + DIDContract
- gRPC
- ElasticSearch + Go

# TODO:
- Create a command to seed data

- Transform ./fabric/ChaincodeSDK.go into ./service/ChaincodeService.go
- Create .env and use it on ./config/config.go
- Improve semantic of vars and functions
- Create a MetadataValidador that will be used to validate the sintaxe and semanthic of the data. All property is required.
- Create MetadataDTO that will keep the metadata entry structure versionated, and provide getter and setters
- Create a README.md to the project
- Add ElasticSearch
- Apply repository pattern
- Apply observer

As a developer specialized on search engines like elasticSearch, Golang and blockchain with Hyperledger Fabric, add
  ElasticSearch to this projet, create the indexes on while make a transaction to create metadata; heat when the
  metadata is updated, and list the metadata with the data in elasticsearch in the searchs of GetAllMetadataHandler
  of domain/metadata/MetadataController

Add what was learned and done in brain.md