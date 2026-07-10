#!/bin/bash
set -e

# Remove sentinel file if it exists to signal that the network is not ready during setup
rm -f "$(dirname "$0")/../.fabric_ready"

# Define variables
FABRIC_VERSION="2.5.4"
CA_VERSION="1.5.7"
FABRIC_DIR="./fabric-samples"
CHANNEL_NAME="metadatachannel"
CHAINCODE_NAME="basic"
CHAINCODE_PATH="${PWD}/chaincode" 

echo "=================================================="
echo "    Setting up Hyperledger Fabric Test Network    "
echo "=================================================="

# Download fabric-samples if not exists
if [ ! -d "$FABRIC_DIR" ]; then
    echo "Downloading Fabric binaries and docker images..."
    curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
    ./install-fabric.sh docker samples binary --fabric-version $FABRIC_VERSION --ca-version $CA_VERSION
    rm install-fabric.sh
fi

cd "$FABRIC_DIR/test-network"

echo "Tearing down any previous state..."
./network.sh down

echo "Starting test network with CouchDB..."
./network.sh up createChannel -c $CHANNEL_NAME -ca -s couchdb

echo "Deploying the chaincode..."
# Note: Ensure that the path contains a valid go.mod and main package if it's a chaincode. 
# Depending on the setup, we deploy it using the absolute path.
./network.sh deployCC \
  -c $CHANNEL_NAME \
  -ccn $CHAINCODE_NAME \
  -ccp $CHAINCODE_PATH \
  -ccl go

echo "=================================================="
echo "          Network Setup Completed!                "
echo "=================================================="