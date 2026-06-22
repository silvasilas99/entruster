# Setting things up

## Install the project dependencies

```bash
# Install HyperledgerFabric binaries
curl -sSL https://bit.ly/2ysbOFE | bash -s

# Set required vars (persist by adding to ~/.bashrc)
export GOROOT=/home/silas/go_dist/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/home/silas/go
export PATH=$PATH:$GOPATH/bin
export CONTAINER_CLI_COMPOSE="docker compose"
```

Install Go dependencies:
```bash
# API server
go mod tidy

# Chaincode (separate module)
cd chaincode && go mod tidy && cd ..
```

-----

## Up the network

> ⚠️ Always run `./network.sh down` first to clear stale Docker volumes.

```bash
cd fabric-samples/test-network

# Purge previous state (containers + volumes)
./network.sh down

# Create the test network with CouchDB and the metadatachannel
./network.sh up createChannel -c metadatachannel -ca

# Deploy the chaincode (absolute path required)
./network.sh deployCC \
  -c metadatachannel -ccn basic \
  -ccp /mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/chaincode \
  -ccl go

cd ../..
```

-----

## Start the server

```bash
# Set the Hyperledger Fabric test network path environment variable
export TEST_NETWORK_PATH='/mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric-samples/test-network'

# Compile and run the server entrypoint
go run cmd/server/main.go
```

Server: http://localhost:8080
Swagger: http://localhost:8080/swagger/index.html
