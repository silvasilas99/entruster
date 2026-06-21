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