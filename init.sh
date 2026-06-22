#!/bin/bash
set -e

echo "Starting Fabric Network Setup..."

# Ensure we have docker and docker-compose installed
if ! command -v docker &> /dev/null; then
    echo "Docker is required but not installed."
    exit 1
fi

# Fix dubious ownership error for go list in git repos
git config --global --add safe.directory '*'

cd $PROJECT_DIR/fabric-samples/test-network

export CONTAINER_CLI_COMPOSE="docker compose"

# Remove the sentinel so the API waits for new certs (don't delete organizations/
# as network.sh needs the existing crypto to run 'down' cleanly).
rm -f "$PROJECT_DIR/.fabric_ready"

# Always run down first to clear stale volumes as per BRAIN.md
echo "Running network.sh down..."
./network.sh down

# Create channel
echo "Running network.sh up createChannel..."
./network.sh up createChannel -c metadatachannel -ca

# Deploy chaincode
echo "Deploying chaincode..."
./network.sh deployCC \
  -c metadatachannel -ccn basic \
  -ccp $PROJECT_DIR/chaincode \
  -ccl go

# Signal the API container that fresh certs are ready
touch "$PROJECT_DIR/.fabric_ready"

# Reconnect the API container to the newly-recreated fabric_test network
echo "Connecting API container to the Fabric network..."
docker network connect fabric_test entruster-api || echo "Warning: Could not connect api container (may already be connected)"

echo "Fabric Network Setup Completed Successfully."
echo "The API container will detect the new certs and start automatically."
