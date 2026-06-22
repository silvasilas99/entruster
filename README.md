# Entruster

A REST API server that registers and queries **healthcare metadata** on a **Hyperledger Fabric** blockchain. Built with Go, Gin, and the Fabric Gateway SDK.

---

## Table of Contents

- [What is Entruster?](#what-is-entruster)
- [Architecture Overview](#architecture-overview)
  - [The Two-Layer Model](#the-two-layer-model)
  - [Request Flow](#request-flow)
  - [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Setup & Running](#setup--running)
  - [1. Install Dependencies](#1-install-dependencies)
  - [2. Start the Fabric Network](#2-start-the-fabric-network)
  - [3. Start the API Server](#3-start-the-api-server)
- [API Reference](#api-reference)
- [Chaincode Reference](#chaincode-reference)
- [Environment Variables](#environment-variables)
- [Restarting After a Previous Session](#restarting-after-a-previous-session)

---

## What is Entruster?

Entruster is a **middleware service** in a larger healthcare interoperability system. Its job is to act as a trusted registrar ("entruster") that:

1. Receives metadata about a healthcare asset (e.g. a patient record, a medical image) via a REST API.
2. Writes that metadata to a **Hyperledger Fabric ledger**, making it immutable and auditable.
3. Exposes read, update, delete, and audit-trail endpoints so other services can query the ledger through a familiar HTTP interface.

---

## Architecture Overview

### The Two-Layer Model

Entruster is composed of **two separate Go programs** that talk to each other over gRPC. This is fundamental to how Hyperledger Fabric works — you cannot merge them into one.

```
┌──────────────────────────────────────────────────────────────┐
│                    YOUR MACHINE (WSL)                         │
│                                                               │
│   ┌───────────────────────────────────────────────────────┐  │
│   │                  Gin REST API Server                  │  │
│   │              cmd/server/main.go  (:8080)              │  │
│   │                                                       │  │
│   │   routes/api_routes.go          ← HTTP routing        │  │
│   │   domain/metadata/                                    │  │
│   │     MetadataController.go       ← HTTP handlers       │  │
│   │     MetadataContracts.go        ← GATEWAY CLIENT      │  │
│   │     MetadataModel.go            ← shared data types   │  │
│   │   fabric/ChaincodeSdk.go        ← gRPC connection     │  │
│   │   config/config.go              ← channel / peer cfg  │  │
│   └───────────────────┬───────────────────────────────────┘  │
│                        │  gRPC over TLS (port 7051)           │
└────────────────────────┼─────────────────────────────────────┘
                         │
┌────────────────────────┼─────────────────────────────────────┐
│   DOCKER CONTAINERS  (Hyperledger Fabric Network)             │
│                        │                                      │
│   ┌────────────────────▼──────────────────────────────────┐  │
│   │          peer0.org1  /  peer0.org2  (channel)         │  │
│   │                                                        │  │
│   │   chaincode/main.go             ← ON-CHAIN CONTRACT   │  │
│   │   (fabric-contract-api-go)                             │  │
│   │   Runs INSIDE the peer, reads/writes the ledger        │  │
│   └───────────────────────────────────────────────────────┘  │
│                                                               │
│   orderer.example.com               ← ordering service        │
└───────────────────────────────────────────────────────────────┘
```

| Layer | File | Library | Runs on |
|---|---|---|---|
| **Gateway Client** | `domain/metadata/MetadataContracts.go` | `fabric-gateway/pkg/client` | Your machine |
| **gRPC Connection** | `fabric/ChaincodeSdk.go` | `fabric-gateway`, `google.golang.org/grpc` | Your machine |
| **On-chain Contract** | `chaincode/main.go` | `fabric-contract-api-go` | Inside the Docker peer |

> **Why two separate programs?**
> The API server is a **client** that sends requests to the peer over gRPC — analogous to an HTTP client. The chaincode is a **server** that receives those requests and writes to the immutable ledger. They use completely different SDKs and cannot be merged into one binary.

### Request Flow

A single `POST /api/metadata` traces through the system like this:

```
HTTP Client (Postman / frontend)
  → POST /api/metadata
    → routes/api_routes.go               (matches route)
      → MetadataController.go            (parses & validates JSON body)
        → MetadataContracts.go           (calls contract.SubmitTransaction)
          → fabric/ChaincodeSdk.go       (signs & sends gRPC to peer)
            → chaincode/main.go          (executes RegisterMetadataOnNetwork)
              → Fabric Ledger            (asset written, ID auto-generated ✅)
```

Read operations (`GET`) use `EvaluateTransaction` instead of `SubmitTransaction` — they query a single peer without going through the ordering service, making them faster and not creating a blockchain transaction.

### Project Structure

```
entruster/
│
├── cmd/server/main.go          # Entry point: connects to Fabric, starts Gin
│
├── config/config.go            # Peer endpoint, channel, chaincode name, TLS paths
│
├── fabric/
│   └── ChaincodeSdk.go         # Establishes gRPC + TLS connection to peer0.org1
│
├── domain/metadata/
│   ├── MetadataModel.go        # MetadataAsset struct + custom JSON unmarshalling
│   ├── MetadataContracts.go    # Gateway client: wraps all 6 chaincode calls
│   └── MetadataController.go  # Gin HTTP handlers for all 6 routes
│
├── routes/
│   └── api_routes.go           # Wires HTTP routes to controllers + Swagger
│
├── utils/
│   └── response.go             # SendSuccess / SendError helpers + typed structs
│
├── docs/                       # Auto-generated Swagger docs (do not edit manually)
│
├── chaincode/
│   ├── go.mod                  # Standalone Go module (separate from the API server)
│   └── main.go                 # On-chain Fabric contract: the 6 ledger transactions
│
├── fabric-samples/
│   └── test-network/           # Hyperledger Fabric test network scripts & Docker config
│
├── BRAIN.md                    # Developer knowledge log: all bugs encountered & fixed
├── SETUP.md                    # Quick setup reference
└── README.md                   # This file
```

---

## Prerequisites

- **Docker** and **Docker Compose**
- **Go** ≥ 1.21 (native Linux binary — see note below)
- **WSL Ubuntu** (project runs on Windows Subsystem for Linux)

> **WSL Go Note:** The Fabric peer CLI requires a native Linux Go binary to compile and package chaincode. Install it locally without `sudo`:
> ```bash
> curl -Lo /tmp/go1.22.10.linux-amd64.tar.gz https://dl.google.com/go/go1.22.10.linux-amd64.tar.gz
> mkdir -p ~/go_dist && tar -C ~/go_dist -xzf /tmp/go1.22.10.linux-amd64.tar.gz
> ```
> Then add to `~/.bashrc`:
> ```bash
> export GOROOT=/home/silas/go_dist/go
> export PATH=$PATH:$GOROOT/bin
> export GOPATH=/home/silas/go
> export PATH=$PATH:$GOPATH/bin
> ```

---

## Setup & Running

### 1. Install Dependencies

```bash
# Install Hyperledger Fabric binaries and Docker images
curl -sSL https://bit.ly/2ysbOFE | bash -s

# Set environment variables for this session (or persist in ~/.bashrc)
export GOROOT=/home/silas/go_dist/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/home/silas/go
export PATH=$PATH:$GOPATH/bin
export CONTAINER_CLI_COMPOSE="docker compose"
```

Install Go API server dependencies:
```bash
go mod tidy
```

Install chaincode dependencies:
```bash
cd chaincode && go mod tidy && cd ..
```

---

### 2. Start the Fabric Network

> ⚠️ **Always run `network.sh down` first** to avoid stale ledger volume conflicts from previous sessions.

```bash
cd fabric-samples/test-network

# Tear down any previous state (containers + volumes)
./network.sh down

# Start network with CouchDB and create the channel
./network.sh up createChannel -c metadatachannel -ca

# Deploy the chaincode (use absolute path)
./network.sh deployCC \
  -c metadatachannel \
  -ccn basic \
  -ccp /mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/chaincode \
  -ccl go

cd ../..
```

Verify the network is healthy:
```bash
docker ps --format "table {{.Names}}\t{{.Status}}"
# Should show: peer0.org1, peer0.org2, orderer, couchdb0, couchdb1 — all Up
```

---

### 3. Start the API Server

```bash
# Point the server at the test network's crypto material
export TEST_NETWORK_PATH='/mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric-samples/test-network'

# Run the server
go run cmd/server/main.go
```

The server will start on **http://localhost:8080**.
Interactive Swagger UI is available at **http://localhost:8080/swagger/index.html**.

---

## API Reference

All endpoints are prefixed with `/api/metadata`.

| Method | Path | Description | Chaincode |
|--------|------|-------------|-----------|
| `POST` | `/api/metadata/` | Register a new metadata asset | `RegisterMetadataOnNetwork` |
| `GET` | `/api/metadata/` | List all metadata assets | `GetAllMetadataFromNetwork` |
| `GET` | `/api/metadata/:id` | Get a single asset by ID | `GetMetadataById` |
| `PUT` | `/api/metadata/:id` | Update an existing asset | `UpdateMetadataById` |
| `DELETE` | `/api/metadata/:id` | Delete an asset | `DeleteMetadataById` |
| `GET` | `/api/metadata/:id/auditory` | Get full audit trail for an asset | `GetMetadataAuditoryById` |

### POST `/api/metadata/` — Example Request Body

```json
{
  "patient_id": "1",
  "asset_id": "1",
  "zkp_proof": "",
  "name": "Chest X-Ray",
  "value": "base64encodeddata",
  "version": "1.0",
  "owner": "Hospital A",
  "rights": "read-only",
  "terms_of_access": "research-only",
  "created_at": "2026-06-22T17:00:00Z",
  "updated_at": "2026-06-22T17:00:00Z",
  "created_by": "dr.silva",
  "updated_by": "dr.silva"
}
```

> `patient_id` and `asset_id` must be **numeric strings** (e.g. `"1"`, `"42"`). The asset `id` is **auto-generated** by the chaincode — never send it on create.

### Example Response

```json
{
  "success": true,
  "message": "Metadata registered successfully",
  "data": {
    "patient_id": "1",
    "asset_id": "1"
  }
}
```

---

## Chaincode Reference

The on-chain contract (`chaincode/main.go`) exposes the following Fabric transactions:

| Transaction | Type | Description |
|---|---|---|
| `InitLedger` | Submit | Seeds the ID counter to 0 (call once on deploy) |
| `RegisterMetadataOnNetwork` | Submit | Creates a new asset; auto-increments ID |
| `GetAllMetadataFromNetwork` | Evaluate | Returns all assets as a JSON array |
| `GetMetadataById` | Evaluate | Returns a single asset by its auto-generated ID |
| `UpdateMetadataById` | Submit | Replaces mutable fields of an existing asset |
| `DeleteMetadataById` | Submit | Removes an asset from the world state |
| `GetMetadataAuditoryById` | Evaluate | Returns the full history (every tx) for an asset |

> **Submit** = goes through endorsement + ordering → creates a blockchain transaction (state-changing).
> **Evaluate** = queries a single peer locally → no blockchain transaction (read-only).

---

## Environment Variables

| Variable | Required | Description |
|---|---|---|
| `TEST_NETWORK_PATH` | ✅ Yes | Absolute path to `fabric-samples/test-network/` |
| `CONTAINER_CLI_COMPOSE` | Optional | Set to `"docker compose"` if `docker-compose` binary is unavailable |

---

## Restarting After a Previous Session

Docker volumes persist ledger state between reboots. Always use the full teardown sequence:

```bash
# Set Go in PATH
export GOROOT=/home/silas/go_dist/go && export PATH=$PATH:$GOROOT/bin

cd fabric-samples/test-network

# 1. Wipe containers AND volumes
./network.sh down

# 2. Fresh start
./network.sh up createChannel -c metadatachannel -ca

# 3. Redeploy chaincode
./network.sh deployCC \
  -c metadatachannel -ccn basic \
  -ccp /mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/chaincode \
  -ccl go

cd ../..

# 4. Run the server
export TEST_NETWORK_PATH='/mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric-samples/test-network'
go run cmd/server/main.go
```

---

> For a detailed log of all bugs encountered and resolved during development, see [BRAIN.md](./BRAIN.md).
