# Brain - Project Knowledge & Setup Log

## Context & Chaincode Deployment Error
When attempting to package and deploy the basic chaincode onto the Hyperledger Fabric test-network using:
```bash
./network.sh deployCC \
  -c metadatachannel -ccn basic \
  -ccp ../asset-transfer-basic/chaincode-go -ccl go
```
the deployment failed with the following error:
```
Error: failed to normalize chaincode path: failed to determine module root: exec: "go": executable file not found in $PATH
Error: failed to read chaincode package at 'basic.tar.gz': open basic.tar.gz: no such file or directory
Chaincode packaging has failed
```

### Analysis & Findings
- **WSL vs. Host Environment:** The project is running within a WSL Ubuntu environment. While Go was installed on the Windows host and located in the shared PATH (`/mnt/d/Program Files/Go/bin`), it was a Windows executable (`go.exe`). The Fabric `peer` CLI tool executing inside the WSL environment requires a native Linux Go binary to compile and package the chaincode.
- **Permission constraints:** As `sudo` required a password, any global native package installation via apt was not feasible without interactive credentials.

---

## Solutions & Actions Taken

### 1. Local Native Go Installation
We downloaded and installed a local version of Go `1.22.10` directly inside WSL without needing `sudo` privileges:
```bash
curl -Lo /tmp/go1.22.10.linux-amd64.tar.gz https://dl.google.com/go/go1.22.10.linux-amd64.tar.gz
mkdir -p /home/silas/go_dist
tar -C /home/silas/go_dist -xzf /tmp/go1.22.10.linux-amd64.tar.gz
```

### 2. Environment Shell Configuration
To ensure `go` remains in the PATH for all interactive WSL sessions, we appended the required environment variables to both `~/.bashrc` and `~/.profile`:
```bash
# Go environment variables
export GOROOT=/home/silas/go_dist/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/home/silas/go
export PATH=$PATH:$GOPATH/bin
```

### 3. Deploy Validation
After adding the native Linux Go binaries to the active session's `PATH`, we successfully ran the deployment script. The chaincode definition has been successfully:
1. Compiled and packaged as `basic_1.0`.
2. Installed on peer `peer0.org1` and `peer0.org2`.
3. Approved by both Org1 and Org2 MSPs.
4. Committed onto `metadatachannel`.

---

## Fabric Gateway Dependency Import Error
When running `go mod tidy`, the command failed with the following error:
```
go: finding module for package github.com/hyperledger/fabric-gateway/pkg/networkClient
...
github.com/hyperledger/fabric-gateway/pkg/networkClient: module github.com/hyperledger/fabric-gateway@latest found (v1.11.0), but does not contain package github.com/hyperledger/fabric-gateway/pkg/networkClient
```

### Analysis & Findings
- **Incorrect Import Path:** [ChaincodeSdk.go](file:///mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric/ChaincodeSdk.go) was importing `"github.com/hyperledger/fabric-gateway/pkg/networkClient"`.
- **Correct Path:** The correct package containing the Fabric Gateway Go SDK client abstractions (`Contract`, `Gateway`, `Connect`, etc.) is `"github.com/hyperledger/fabric-gateway/pkg/client"`.

### Solutions & Actions Taken
- **Corrected Import:** Updated [ChaincodeSdk.go](file:///mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric/ChaincodeSdk.go) to import `"github.com/hyperledger/fabric-gateway/pkg/client"` instead of the invalid `"github.com/hyperledger/fabric-gateway/pkg/networkClient"`. 
- **Aliased Namespace:** Kept the package alias as `networkClient` (e.g. `networkClient "github.com/hyperledger/fabric-gateway/pkg/client"`) to preserve compatibility with existing code references without having to modify references like `networkClient.Contract` or `networkClient.Gateway` throughout the file.
- **Verification:** Verified that other source files (`MetadataController.go`, `MetadataService.go`, and `api_routes.go`) are already using the correct standard `"github.com/hyperledger/fabric-gateway/pkg/client"` import path.

---

## Fabric Chaincode Deployment Redefined Sequence Error (Status 500)
When attempting to deploy or upgrade the chaincode using `./network.sh deployCC`, the deployment failed with the following error:
```
Error: proposal failed with status: 500 - failed to invoke backing implementation of 'ApproveChaincodeDefinitionForMyOrg': attempted to redefine uncommitted sequence (3) for namespace basic with unchanged content
Chaincode definition approved on peer0.org1 on channel 'metadatachannel' failed
Deploying chaincode failed
```

### Analysis & Findings
- **State Mismatch:** The committed chaincode definition sequence on the channel was `2`. However, in a previous execution, Organization 1 had successfully approved sequence `3`, but the definition was never committed (either because Organization 2 hadn't approved it yet or the script stopped/failed before the commit phase).
- **Auto-detection Logic:** The deployment script checks the channel for the committed sequence (`2`) and queries the active peer organization (which was Org 2 from the previous step, having only approved sequence `2`). Seeing both committed and approved sequences match, it auto-detects that the next sequence to deploy is `3`.
- **Fabric Lifecycle Restriction:** Fabric does not allow an organization to re-approve the exact same definition on an uncommitted sequence (e.g., Org 1 trying to approve sequence `3` again with the same package ID and version). This results in the `attempted to redefine uncommitted sequence` error.

### Solutions & Actions Taken
1. **Manual Sequence Alignment:** We manually ran commands setting variables for Organization 2 to approve sequence `3` using the existing package ID:
   ```bash
   export PATH="/mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric-samples/bin:$PATH"
   export FABRIC_CFG_PATH="/mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/fabric-samples/config/"
   source scripts/envVar.sh
   setGlobals 2
   peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" --channelID metadatachannel --name basic --version 1.0 --package-id basic_1.0:1f66cc610fb8ac88d69e78edc4c1cb09bac7c839fb5399009310ac8dfa703799 --sequence 3
   ```
2. **Manual Commit:** Once both Org 1 and Org 2 approved sequence `3`, we manually committed sequence `3` on the channel:
   ```bash
   setGlobals 1
   peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$ORDERER_CA" --channelID metadatachannel --name basic --version 1.0 --sequence 3 --peerAddresses localhost:7051 --tlsRootCertFiles "$PEER0_ORG1_CA" --peerAddresses localhost:9051 --tlsRootCertFiles "$PEER0_ORG2_CA"
   ```
3. **Successful Execution:** With sequence `3` successfully committed to the ledger, running the `./network.sh deployCC` command again allowed the script to auto-detect sequence `4`, approve it for both organizations, and commit it cleanly.

---

## Go Import Cycle: `fabric` ↔ `domain/metadata`

When trying to build the project, the following error occurred:
```
package command-line-arguments
    imports github.com/silvasilas99/entruster/fabric
    imports github.com/silvasilas99/entruster/domain/metadata
    imports github.com/silvasilas99/entruster/fabric: import cycle not allowed
```

### Analysis & Findings
- `fabric/ChaincodeSdk.go` imported `domain/metadata` (to use `metadata.MetadataModel` as a parameter type in `RegisterMetadataOnNetwork` and `GetAllMetadataFromNetwork`).
- `domain/metadata/MetadataController.go` was calling `RegisterMetadataOnNetwork` as a local package function, expecting it to live inside the `metadata` package — but the function was defined in `fabric`.
- This created a circular dependency: `fabric → metadata → fabric`.

### Solutions & Actions Taken
1. **Moved `RegisterMetadataOnNetwork` and `GetAllMetadataFromNetwork`** from `fabric/ChaincodeSdk.go` into `domain/metadata/MetadataRepository.go`. Both functions only need `fabric-gateway/pkg/client` and local types, so no cycle is introduced.
2. **Removed the `domain/metadata` import** from `fabric/ChaincodeSdk.go`. The `fabric` package now only handles connection/identity setup and exposes `Connect()`.
3. The dependency graph became acyclic: `fabric → config` and `metadata → fabric-gateway/client`.

---

## Chaincode Argument Mismatch: `RegisterMetadataOnNetwork`

After fixing the import cycle, calling `POST /api/metadata` returned:
```json
{
  "message": "metadata.RegisterMetadataOnNetwork: Internal error. Failed to submit transaction: rpc error: code = Aborted desc = failed to endorse transaction",
  "success": false
}
```

### Analysis & Findings
The Go application was calling `SubmitTransaction` with **11 arguments** in the wrong order, while the chaincode's `MetadataContract.RegisterMetadataOnNetwork` expected **13 arguments**:

| Position | Application (wrong) | Chaincode (expected) |
|---|---|---|
| 1 | `ID` (string) | `patientID` (uint64) |
| 2 | `name` | `assetID` (uint64) |
| 3 | `value` | `zkpProof` (string) |
| 4–11 | `version … updatedBy` | `name … updatedBy` |
| — | *(missing)* | `assetID`, `zkpProof` |

Additionally, the `ID` field was being sent by the client, but the chaincode **auto-generates the ID** using an internal counter (`_metadata_id_counter`). Sending it as the first argument caused a type mismatch (`string` vs `uint64`).

The `MetadataModel` struct was also missing the `ZKPProof` field that the chaincode requires.

### Solutions & Actions Taken
1. **Updated `MetadataModel`**: Removed the `ID` field (auto-generated by chaincode) and added `ZKPProof string \`json:"zkp_proof"\``.
2. **Fixed `RegisterMetadataOnNetwork` in `MetadataRepository.go`**: Changed `SubmitTransaction` to pass arguments in the correct order: `PatientID, AssetID, ZKPProof, Name, Value, Version, Owner, Rights, TermsOfAccess, CreatedAt, UpdatedAt, CreatedBy, UpdatedBy`.
3. **Fixed `MetadataController.go`**: Updated the success response to return `patient_id` and `asset_id` instead of the now-removed `id` field.

### Correct Request Body
```json
{
  "patient_id": "1",
  "asset_id": "1",
  "zkp_proof": "",
  "name": "name",
  "value": "value",
  "version": "version",
  "owner": "owner",
  "rights": "rights",
  "terms_of_access": "termsOfAccess",
  "created_at": "createdAt",
  "updated_at": "updatedAt",
  "created_by": "createdBy",
  "updated_by": "updatedBy"
}
```
> `patient_id` and `asset_id` must be numeric strings (e.g. `"1"`, `"42"`) since the chaincode converts them to `uint64`.

---

## Full Chaincode Layer Implementation

### Context
`MetadataContracts.go` existed but was completely empty. `MetadataController.go` had all handlers stubbed with `501 Not Implemented`. `api_routes.go` only wired `POST /` and `GET /`, with the remaining routes commented out. `MetadataRepository.go` held duplicate versions of the two functions already moved from `fabric/`.

### Actions Taken

#### 1. `MetadataContracts.go` — implemented all 6 contract functions

| Function | SDK Method | Chaincode Transaction |
|---|---|---|
| `CreateMetadata` | `SubmitTransaction` | `RegisterMetadataOnNetwork` |
| `GetAllMetadata` | `EvaluateTransaction` | `GetAllMetadataFromNetwork` |
| `GetMetadataByID` | `EvaluateTransaction` | `GetMetadataById` |
| `UpdateMetadataByID` | `SubmitTransaction` | `UpdateMetadataById` |
| `DeleteMetadataByID` | `SubmitTransaction` | `DeleteMetadataById` |
| `GetMetadataAuditoryByID` | `EvaluateTransaction` | `GetMetadataAuditoryById` |

- `EvaluateTransaction` is used for read-only queries (no ledger write, no consensus).
- `SubmitTransaction` is used for state-changing operations (goes through endorsement + ordering).
- All responses from `Evaluate` are JSON-unmarshalled into typed structs.

#### 2. `MetadataModel.go` — added `MetadataHistoryEntry`

```go
type MetadataHistoryEntry struct {
    TxID      string        `json:"tx_id"`
    Timestamp string        `json:"timestamp"`
    IsDelete  bool          `json:"is_delete"`
    Value     MetadataModel `json:"value"`
}
```
Used as the return type of `GetMetadataAuditoryByID` to represent each entry in an asset's audit trail.

#### 3. `MetadataController.go` — fully implemented all handlers

Replaced all `501 Not Implemented` stubs with real implementations calling the corresponding contract function. Added `GetMetadataAuditoryByIDHandler`. Removed the old `ExportMetadataAsCsvHandler` stub.

#### 4. `api_routes.go` — wired all 6 routes

```
POST   /api/metadata/
GET    /api/metadata/
GET    /api/metadata/:id
PUT    /api/metadata/:id
DELETE /api/metadata/:id
GET    /api/metadata/:id/auditory
```

The `/auditory` sub-path avoids a gin router conflict with the plain `/:id` parameter route.

#### 5. `MetadataRepository.go` — cleaned up

Removed the duplicate `RegisterMetadataOnNetwork` and `GetAllMetadataFromNetwork` functions. The file now serves as a package-responsibility comment only, since all chaincode I/O lives in `MetadataContracts.go`.

### Responsibility Split (final)

```
MetadataModel.go      — data types & JSON unmarshalling
MetadataContracts.go  — Fabric ledger I/O (chaincode calls)
MetadataController.go — HTTP handlers (gin)
MetadataRepository.go — (empty, reserved for future DB/off-chain queries)
```

---

## Added Incremental `ID` Field to `MetadataModel`

### Context
The chaincode auto-generates asset IDs via an internal counter (`_metadata_id_counter`). Previously the Go model had no `ID` field, so read responses from the ledger would silently drop the `id` value.

### Actions Taken
Added `ID uint64` as the first field of `MetadataModel`:

```go
ID uint64 `json:"id,omitempty"`
```

**Design decisions:**
- **`uint64`** — matches the chaincode's counter type.
- **`omitempty`** — the field is absent from `POST` bodies when zero, so clients never need to send it on create.
- **No change to `CreateMetadata`** — `ID` is never passed as a transaction argument; the chaincode generates it.
- **Automatic on reads** — the existing `UnmarshalJSON` picks up `"id"` from chaincode JSON responses for `GetAll`, `GetById`, and `GetAuditory` without any additional code.

---

## Swagger Documentation Integration

### Context
The application had partial setup for Swagger (the Swagger annotations were defined on `main.go` and `MetadataController.go`, and the anonymous import of `_ "github.com/silvasilas99/entruster/docs"` was present in `main.go`). However:
1. The `docs/` package had not been generated yet, causing compile errors.
2. The `/swagger/*any` route handler was not registered in `routes/api_routes.go`.
3. Generating the documentation failed because the annotations in `MetadataController.go` referenced `utils.SuccessResponse` and `utils.ErrorResponse` structs which did not exist.

### Actions Taken
1. **Added Response Structs**: Created `SuccessResponse` and `ErrorResponse` structs inside [response.go](file:///mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/utils/response.go). Updated `SendSuccess` and `SendError` helpers to use these typed structs instead of untyped `gin.H` maps.
2. **Configured Router**: Added the `/swagger/*any` GET route inside [api_routes.go](file:///mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/routes/api_routes.go) to bind the Gin middleware using `github.com/swaggo/gin-swagger` and `github.com/swaggo/files`.
3. **Generated Documentation**: Executed `swag init -g cmd/server/main.go` to generate the `docs/` directory (`docs.go`, `swagger.json`, `swagger.yaml`).
4. **Resolved Module Dependencies**: Ran `go mod tidy` to add the Swagger middleware modules (`github.com/swaggo/gin-swagger` and `github.com/swaggo/files`) to the direct requirements in `go.mod`.
5. **Verified Build**: Verified that the server successfully compiles without any errors (`go build ./cmd/server`).

---

## Fixed JSON Unmarshal Error: Type Mismatch for `patient_id` and `asset_id`

### Context
When calling `GET /api/metadata`, the application failed with the error:
`metadata.GetAllMetadata: failed to unmarshal response: json: cannot unmarshal number into Go struct field .patient_id of type string`

This happens because the chaincode Go code treats `patientID` and `assetID` as `uint64`. When writing to/reading from the Fabric ledger state, the values are marshalled into JSON as raw numeric values (e.g. `1`), whereas the API's `MetadataModel` defined `PatientID` and `AssetID` as `string`.

### Actions Taken
1. **Created `flexibleString` Helper**: Introduced a custom type alias `flexibleString string` inside [MetadataModel.go](file:///mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/domain/metadata/MetadataModel.go). Implemented custom JSON unmarshalling logic on it to handle unmarshalling from both string values and numeric values (using `json.Number`).
2. **Updated `UnmarshalJSON`**: Configured `MetadataModel`'s custom JSON unmarshaller to parse `patient_id` and `asset_id` using this `flexibleString` helper. Outer fields override the inner embedded struct tags during unmarshalling, successfully avoiding the unmarshal error while preserving the `string` representation inside the Go code.
3. **Verified and Rebuilt**: Regenerated Swagger docs and successfully compiled the server with `go build ./cmd/server`.

---

## Missing Standalone Chaincode Module

### Context
Attempting `./network.sh deployCC -ccp ../../chaincode -ccl go` failed with:
```
Path to chaincode does not exist. Please provide different path.
```
Additionally, the peer was unreachable (`connection refused on 7051`) because the network was not running.

### Analysis & Findings
- **Architecture mismatch:** `domain/metadata/MetadataContracts.go` is the **API client** — it uses `fabric-gateway/pkg/client` to call chaincode via gRPC. It is NOT on-chain chaincode. No actual chaincode existed in the project.
- **Real chaincode requirements:** On-chain chaincode must use `github.com/hyperledger/fabric-contract-api-go`, implement `contractapi.Contract`, have its own `go.mod`, and live in a separate Go module.
- **Stale Docker state:** The network containers from a previous session were in `Exited` state but their ledger volumes were still present. Running `network.sh up` again caused `"ledger already exists"` errors on the peers, preventing them from starting correctly.

### Actions Taken
1. **Created standalone chaincode module** at `chaincode/` with its own `go.mod` (`module github.com/silvasilas99/entruster/chaincode`) and `chaincode/main.go` implementing all 6 transactions:
   - `RegisterMetadataOnNetwork` — creates asset with auto-incremented ID
   - `GetAllMetadataFromNetwork` — range query over all assets
   - `GetMetadataById` — single asset lookup
   - `UpdateMetadataById` — replaces mutable fields
   - `DeleteMetadataById` — removes asset from ledger
   - `GetMetadataAuditoryById` — returns full audit history via `GetHistoryForKey`
2. **Full teardown before restart:** Ran `./network.sh down` to remove stale containers AND volumes, then `./network.sh up createChannel -c metadatachannel -ca` for a clean boot.
3. **Deployed with correct absolute path:**
   ```bash
   export GOROOT=/home/silas/go_dist/go && export PATH=$PATH:$GOROOT/bin
   ./network.sh deployCC \
     -c metadatachannel -ccn basic \
     -ccp /mnt/d/@PROJETOS/Mestrado/Interopchain/entruster/chaincode \
     -ccl go
   ```
4. **Result:** Chaincode `basic_1.0` (sequence 1) committed on `metadatachannel`, approved by both Org1MSP and Org2MSP.
   - Package ID: `basic_1.0:da066418ccd0481272dd92c966567a19a7a2a32409ab15e2f03c28d6a8b64127`

### Key Rule Going Forward
**Always run `./network.sh down` before `./network.sh up`** when restarting after a previous session, to avoid stale volume conflicts. The correct deploy command uses the absolute path to `chaincode/`.

---

## Containerization & Environment Automation

### Context
The user wanted to act as a DevOps engineer and dockerize the application so that a single `docker compose up -d` command could spin up the entire Hyperledger Fabric network (following the steps in `SETUP.md`), deploy the chaincode, and boot up the API container correctly connected to the network.

### Analysis & Findings
1. **Network Mismatches**: The default `docker-compose.yml` network rules conflict with Fabric's `network.sh`, which automatically provisions its own network `fabric_test`. When the Go API container booted on the docker-compose network, it was unable to resolve peer nodes.
2. **Setup Script Containerization (`docker-in-docker`)**: Running `network.sh` requires access to the docker daemon. To automate this, we built a `setup` container that mounts `/var/run/docker.sock`. However, because the setup container was technically running on a separate docker bridge network, mapping ports and executing Docker CLI commands created labeling and working directory mismatches with the host daemon.
3. **Git Ownership during Chaincode Package**: Inside the setup container, packaging the chaincode failed with `fatal: detected dubious ownership in repository` because the container user (root) didn't match the host user that owned the mounted volume.
4. **Race Condition**: The `api` container booted before the Fabric network and chaincode were ready. Fabric's CA generation and chaincode deployment can take 1-2 minutes.

### Solutions & Actions Taken
1. **Docker-in-Docker Setup Container**: Created `Dockerfile.setup` and `init.sh` to encapsulate the deployment logic. We configured the `setup` container in `docker-compose.yml` with `network_mode: "host"`, mapped `${PWD}` identically, and set the `working_dir` explicitly. This allows the setup container to seamlessly control the host's docker daemon and communicate with Fabric CAs natively over localhost during network bring-up without volume mounting conflicts.
2. **Git Safe Directory Fix**: Modified `init.sh` to execute `git config --global --add safe.directory '*'` before calling `network.sh deployCC`, solving the `go list` packaging failures.
3. **Dynamic Network Join**: Modified `init.sh` to run `docker network connect fabric_test entruster-api` upon successful deployment, natively bridging the API to Fabric's internal network.
4. **API Container Resiliency**: Created `Dockerfile.api` using a multi-stage Go build. Set `restart: on-failure` for the `api` service in docker-compose. The API container continuously attempts to boot and connect; once `init.sh` completes the Fabric setup and joins the `api` container to `fabric_test`, the API successfully establishes its gRPC connection and serves traffic on port `8080`.
5. **Config & Environment Variables**: Refactored the Go API configuration (`config/config.go`) to read paths dynamically via environment variables so that paths resolve correctly whether running locally or inside the Docker container.

---

## API to Fabric Connection Errors (TLS & Reachability)

### Context
When testing the newly dockerized API container, calling the POST `/api/metadata` endpoint resulted in various gRPC connection errors:
1. `rpc error: code = Unavailable desc = name resolver error: produced zero addresses`
2. `rpc error: code = Unavailable desc = connection error: desc = "transport: authentication handshake failed: tls: failed to verify certificate: x509: certificate signed by unknown authority"`
3. `rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing: dial tcp 172.19.0.5:7051: connect: no route to host"`
4. `rpc error: code = FailedPrecondition desc = no combination of peers can be derived which satisfy the endorsement policy: No metadata was found for chaincode basic in channel metadatachannel`

### Analysis & Findings
1. **Network Disconnection (`no route to host`)**: Occurred because `init.sh` ran `network.sh down` to clean up old state, which wiped and recreated the `fabric_test` network. The API container was still attached to the old, deleted network interface and couldn't route traffic to the newly created peers.
2. **Missing `fabric_test` in Compose**: The `docker-compose.yml` file did not declare `fabric_test` as an external network, so the `docker compose up` command couldn't natively wire the API container to it from the start.
3. **TLS Handshake Failure**: The Fabric Gateway SDK loads the TLS certificate at initialization. Because the API container started and cached the certificates from an earlier failed network boot, it rejected connections when `init.sh` eventually booted a fresh network with new cryptographic material.
4. **Broken Crypto Material Generation**: A previous cleanup attempt using `rm -rf` accidentally deleted git-tracked files like `organizations/fabric-ca/registerEnroll.sh`. As a result, subsequent `./network.sh up` commands failed silently midway, leaving empty TLS folders for the orderer and causing Genesis block generation to fail.
5. **Chaincode Discovery Failure (`No metadata was found...`)**: Even after the network was fixed, the API container booted slightly before the chaincode was committed. The Gateway SDK cached the discovery response showing no chaincode, leading to `FailedPrecondition` errors.

### Solutions & Actions Taken
1. **Network Wiring**: Explicitly declared `fabric_test` as an external network in `docker-compose.yml` and assigned the `api` container to it, allowing Docker to handle the network lifecycle properly.
2. **Synchronized Boot via Sentinel**: Added a `.fabric_ready` file sentinel mechanism:
   - Modified `init.sh` to write `.fabric_ready` to the shared project directory *only after* the chaincode is fully deployed and committed.
   - Created a new `entrypoint.sh` for the `api` container that runs a `while` loop waiting for the `.fabric_ready` file before it starts the Go server.
   - This completely solved the TLS caching and chaincode discovery issues by ensuring the Go API only loads certificates and queries discovery after the chaincode is confirmed ready.
3. **Nuclear Reset & Git Restore**: Recovered the accidentally deleted `registerEnroll.sh` using `git restore test-network/organizations/`. Performed a full cleanup using `network.sh down` and removed stale `organizations/` directories to force a pristine regeneration of all crypto materials.
4. **Gateway Cache Clear**: With the `entrypoint.sh` wait-loop in place, the `api` container naturally starts with a clean Gateway cache holding the newly committed chaincode metadata, allowing successful transactions.

---

## Auditing Layer: AuditModel, AuditService, and MetadataObserver

### Context
The application needed an auditing mechanism to record every state-changing and listing operation performed on Metadata assets (CREATE, UPDATE, DELETE, LIST). The requirement was:
- Emit an audit entry every time a Metadata operation succeeds.
- Keep the audit logic decoupled from both the Fabric transport layer and the HTTP handlers.
- Use the **Observer pattern** so that `MetadataModel` lifecycle events drive `AuditService` calls without creating circular dependencies.

### Architecture

```
HTTP Request
    ↓
MetadataController  (gin handler)
    ↓
MetadataContracts   (Fabric ledger call)
    ↓  (on success only)
MetadataObserver.OnCreate / OnUpdate / OnDelete / OnList
    ↓
AuditService.Record(entity, entityID, action, actor, details)
    ↓
AuditModel stored in-memory + printed to stdout
```

### Files Created / Modified

| File | Role |
|---|---|
| `audit/AuditModel.go` | `AuditAction` enum (`CREATE`, `UPDATE`, `DELETE`, `LIST`) + `AuditModel` struct |
| `audit/AuditService.go` | Thread-safe in-memory store; `Record`, `GetAll`, `GetByEntity`, `GetByEntityID` |
| `domain/metadata/MetadataObserver.go` | Observer with typed hooks: `OnCreate`, `OnUpdate`, `OnDelete`, `OnList` |
| `domain/metadata/MetadataContracts.go` | Each mutating function now accepts `*MetadataObserver` and fires the hook after a successful ledger commit |
| `domain/metadata/MetadataController.go` | Handler constructors accept and forward `*MetadataObserver` |
| `routes/api_routes.go` | Bootstraps `AuditService → MetadataObserver` once at startup and injects into handlers |

### Audited Operations

| HTTP Operation | Observer Hook | Audit Action |
|---|---|---|
| `POST /api/metadata/` | `OnCreate(req)` | `CREATE` |
| `GET /api/metadata/` | `OnList(count)` | `LIST` |
| `PUT /api/metadata/:id` | `OnUpdate(id, req)` | `UPDATE` |
| `DELETE /api/metadata/:id` | `OnDelete(id)` | `DELETE` |
| `GET /api/metadata/:id` | *(not audited — read-only)* | — |
| `GET /api/metadata/:id/auditory` | *(not audited — read-only)* | — |

### Design Decisions

1. **Observer over direct coupling**: The observer is injected via the constructor (`NewMetadataObserver(svc)`), keeping `MetadataContracts` free of any direct audit dependency. Passing `nil` disables auditing silently, useful for tests.
2. **Hook fires only on success**: The observer methods are called *after* the Fabric transaction commits without error, so no phantom audit entries are created for failed operations.
3. **In-memory store**: `AuditService` currently uses a `sync.RWMutex`-protected slice. To persist audits across restarts, replace the slice/mutex with a database repository while keeping the same `Record`/`GetAll` public interface.
4. **Actor fallback**: If `actor` is blank when `Record` is called, it defaults to `"system"` to keep every audit row attributable.
5. **Single bootstrap point**: `AuditService` and `MetadataObserver` are instantiated once in `routes/api_routes.go` and shared across all handlers, ensuring a single consistent audit trail for the lifetime of the server.

### Key Rule Going Forward
Any new domain entity that requires auditing should follow the same pattern:
1. Create `<Entity>Observer.go` with typed hooks in the domain package.
2. Add `*<Entity>Observer` parameters to the contract functions that change state.
3. Wire `AuditService → <Entity>Observer` in `routes/api_routes.go`.

---

## Elasticsearch Integration & Lifecycle Sync

### Context
The user wanted to integrate a local containerized Elasticsearch engine for advanced metadata querying (e.g. by dates and matching properties) instead of relying solely on the Hyperledger Fabric ledger to read collections of data. They also requested that the newly created `AuditService` strictly handle the auditing endpoint (`/auditory`), fully decoupling this logic from the Fabric chaincode history functions.

### Analysis & Findings
1. **Chaincode Return Payload**: Our Go API's `CreateMetadata` calls `contract.SubmitTransaction`, which initially did not return the auto-generated `ID` assigned by the chaincode to the new asset. We couldn't properly sync Elasticsearch without this primary key.
2. **Coupled Audit Architecture**: The `/api/metadata/:id/auditory` endpoint was originally directly invoking a chaincode function (`GetMetadataAuditoryById`) which queried `GetHistoryForKey` and parsed `MetadataHistoryEntry` objects. This made the chaincode heavy and tightly coupled the API to ledger history capabilities, rendering the new `AuditService` partially redundant.
3. **Soft-deletion Support**: Records deleted from the ledger should still be searchable or traceable via Elasticsearch (or at least marked as soft-deleted) rather than permanently wiped.

### Actions Taken

#### 1. Elastic Search Engine (Docker)
- Added `elasticsearch` (v8.13.0) to `docker-compose.yml`, running on a single-node discovery mode and exposed on port `9200`.
- Passed the `ELASTICSEARCH_URL=http://elasticsearch:9200` environment variable to the `api` container so it automatically connects on boot.
- Developed `ElasticClient.go` to manage connection pooling via the official `go-elasticsearch/v8` client.
- Implemented `ElasticService.go` with foundational capabilities to `IndexDocument`, `UpdateDocument` (partial), `SearchDocuments`, and `DateRangeFilter` via native Elasticsearch Query DSL (`bool.must`, `range`).

#### 2. Synchronizing Data via MetadataObserver
We injected `ElasticService` into the `MetadataObserver` to automatically translate Fabric ledger lifecycle events into search engine sync operations:
- **`OnCreate`**: Maps to an `IndexDocument` command to insert the fresh metadata record.
- **`OnUpdate` (Heating)**: Re-indexes the updated `MetadataModel`, overwriting the older Elasticsearch document to keep the cache hot.
- **`OnDelete`**: Triggers a partial `UpdateDocument` payload setting a `deleted_at` UTC timestamp flag (soft-delete mapping).

#### 3. Chaincode ID Resolution
To ensure the `MetadataObserver.OnCreate` hook has access to the primary key:
- Edited `chaincode/main.go` -> `RegisterMetadataOnNetwork` to return `(string, error)` (the newly incremented ID).
- Updated `MetadataContracts.go` to parse the returned byte array into the ID and feed it into `observer.OnCreate(id, req)`.

#### 4. Audit Decoupling & Cleanup
Completely decoupled auditing logic from the blockchain ledger state:
- Removed `GetMetadataAuditoryById` and `HistoryEntry` from the smart contract (`chaincode/main.go`).
- Removed `MetadataHistoryEntry` from `MetadataModel.go`.
- Rerouted `GetMetadataAuditoryByIDHandler` in `MetadataController.go` to accept the `AuditService` directly and fetch immutable traces locally (`auditSvc.GetByEntityID("Metadata", id)`).
- Updated Swagger documentation to natively return the `audit.AuditModel` response structure and recompiled (`swag init`).

#### 5. Controller Search Overhaul
- Refactored `GetAllMetadataHandler` (`GET /api/metadata/`) to evaluate the Elasticsearch cluster via `ElasticService.SearchDocuments` instead of querying the Fabric network.
- Added native support for `patient_id` and `asset_id` matching, as well as `from` / `to` date range filtering (via `created_at`).
- Implemented a `must_not` exclusion clause enforcing that any record containing a `deleted_at` field is omitted from standard read queries.