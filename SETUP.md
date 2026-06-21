# Setting things up

## Install the project dependencies

```bash
# Install HyperledgerFabric binaries
curl -sSL https://bit.ly/2ysbOFE | bash -s

# Set required vars (you can persist adding this in ~/.bashrc)
export CONTAINER_CLI_COMPOSE="docker compose"
export GOROOT=/home/silas/go_dist/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/home/silas/go
export PATH=$PATH:$GOPATH/bin
```

-----

## Up the network

```bash
# Purge the previous configs
cd fabric-samples/test-network
./network.sh down
cd ../../

# Create the test network, with CounchDB, Go and the metadatachannel and deploy the chaincode
cd fabric-samples/test-network
./network.sh up -s couchdb
./network.sh createChannel -c metadatachannel
./network.sh deployCC \
  -c metadatachannel -ccn basic \
  -ccp ../../chaincode -ccl go
cd ../../
```

-----

## Start the server

```bash
# Set the Hyperledger Fabric test network path environment variable
export TEST_NETWORK_PATH='/mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric-samples/test-network'

# Compile and run the server entrypoint
go run cmd/server/main.go
```
