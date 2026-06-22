# Local Development Setup

This project uses Docker Compose to fully automate the provisioning of the Hyperledger Fabric network, the deployment of the chaincode, and the booting of the Go API server. 

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed and running.
- [Docker Compose](https://docs.docker.com/compose/install/) plugin installed.
- Ensure your user has permissions to run Docker commands (or use `sudo`).

---

## 1. Start the Environment

To spin up the entire architecture (Fabric Network + Chaincode + API), simply run:

```bash
docker compose up -d
```

### What happens under the hood?
1. The **`setup` container** boots up and communicates with your host's Docker daemon.
2. It automatically purges any stale network or volumes.
3. It creates the `metadatachannel` and brings up the Fabric Certificate Authorities (CAs), Orderers, and Peers.
4. It compiles, packages, and deploys the Go chaincode (`basic_1.0`) to the peers.
5. It bridges the API container directly into the `fabric_test` network.
6. The **`api` container** continuously restarts until the network is ready, eventually establishing a gRPC connection to the peer and serving HTTP traffic.

You can monitor the setup progress by checking the logs:
```bash
# Watch the Fabric network provisioning and chaincode deployment
docker compose logs -f setup

# Watch the API server boot
docker compose logs -f api
```

---

## 2. Accessing the Application

Once the `setup` container finishes and the `api` container reports `✅ Connected to Fabric — metadatachannel`, the server is ready to use!

- **Base URL:** http://localhost:8080
- **Swagger Documentation:** http://localhost:8080/swagger/index.html

---

## 3. Shutting Down & Cleanup

To stop the API and completely tear down the Fabric network, destroying all containers and clearing ledger volumes:

```bash
docker compose down -v
```

> **Note:** The `-v` flag is highly recommended to ensure Fabric ledger volumes and crypto material are wiped clean, preventing `ledger already exists` errors on the next boot.

---

## (Optional) Local Native Development

If you prefer to run the API natively on your host machine against the Dockerized Fabric network (instead of running the API inside Docker):

1. Bring up the network using the setup container:
   ```bash
   docker compose up setup -d
   ```
2. Make sure you have Go installed on your machine (`go1.22+`).
3. Set the required environment variables to point to the local crypto material:
   ```bash
   export TEST_NETWORK_PATH="${PWD}/fabric-samples/test-network"
   ```
4. Run the server natively:
   ```bash
   go run cmd/server/main.go
   ```
