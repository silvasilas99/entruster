package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// MetadataContract provides functions for managing metadata assets on the ledger.
type MetadataContract struct {
	contractapi.Contract
}

// counterKey is the ledger key used to persist the auto-increment ID counter.
const counterKey = "_metadata_id_counter"

// MetadataAsset is the on-chain representation of a metadata record.
type MetadataAsset struct {
	ID            uint64 `json:"id"`
	PatientID     uint64 `json:"patient_id"`
	AssetID       uint64 `json:"asset_id"`
	ZKPProof      string `json:"zkp_proof"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Version       string `json:"version"`
	Owner         string `json:"owner"`
	Rights        string `json:"rights"`
	TermsOfAccess string `json:"terms_of_access"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	CreatedBy     string `json:"created_by"`
	UpdatedBy     string `json:"updated_by"`
}

// HistoryEntry wraps a single record from the asset's audit history.
type HistoryEntry struct {
	TxID      string        `json:"tx_id"`
	Timestamp string        `json:"timestamp"`
	IsDelete  bool          `json:"is_delete"`
	Value     MetadataAsset `json:"value"`
}

// ────────────────────────────────────────────────────────────
//  Initialisation
// ────────────────────────────────────────────────────────────

// InitLedger seeds the counter at 0 so the first asset receives ID 1.
func (c *MetadataContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return ctx.GetStub().PutState(counterKey, []byte("0"))
}

// ────────────────────────────────────────────────────────────
//  Helpers
// ────────────────────────────────────────────────────────────

func nextID(ctx contractapi.TransactionContextInterface) (uint64, error) {
	raw, err := ctx.GetStub().GetState(counterKey)
	if err != nil {
		return 0, fmt.Errorf("failed to read counter: %w", err)
	}
	var current uint64
	if raw != nil {
		current, err = strconv.ParseUint(string(raw), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse counter: %w", err)
		}
	}
	next := current + 1
	if err := ctx.GetStub().PutState(counterKey, []byte(strconv.FormatUint(next, 10))); err != nil {
		return 0, fmt.Errorf("failed to update counter: %w", err)
	}
	return next, nil
}

func assetKey(id uint64) string {
	return fmt.Sprintf("metadata_%d", id)
}

func putAsset(ctx contractapi.TransactionContextInterface, asset MetadataAsset) error {
	data, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("failed to marshal asset: %w", err)
	}
	return ctx.GetStub().PutState(assetKey(asset.ID), data)
}

func getAsset(ctx contractapi.TransactionContextInterface, id uint64) (*MetadataAsset, error) {
	data, err := ctx.GetStub().GetState(assetKey(id))
	if err != nil {
		return nil, fmt.Errorf("failed to read asset %d: %w", id, err)
	}
	if data == nil {
		return nil, fmt.Errorf("asset %d does not exist", id)
	}
	var asset MetadataAsset
	if err := json.Unmarshal(data, &asset); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
	}
	return &asset, nil
}

func mustParseUint64(s string, field string) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("field %q must be a non-negative integer, got %q", field, s)
	}
	return v, nil
}

// ────────────────────────────────────────────────────────────
//  Chaincode transactions
// ────────────────────────────────────────────────────────────

// RegisterMetadataOnNetwork creates a new metadata asset on the ledger.
// The ID is auto-generated; the caller must NOT supply it.
// Arguments order (all strings):
//
//	patientID, assetID, zkpProof, name, value, version,
//	owner, rights, termsOfAccess, createdAt, updatedAt, createdBy, updatedBy
func (c *MetadataContract) RegisterMetadataOnNetwork(
	ctx contractapi.TransactionContextInterface,
	patientIDStr string,
	assetIDStr string,
	zkpProof string,
	name string,
	value string,
	version string,
	owner string,
	rights string,
	termsOfAccess string,
	createdAt string,
	updatedAt string,
	createdBy string,
	updatedBy string,
) error {
	patientID, err := mustParseUint64(patientIDStr, "patientID")
	if err != nil {
		return err
	}
	assetID, err := mustParseUint64(assetIDStr, "assetID")
	if err != nil {
		return err
	}

	id, err := nextID(ctx)
	if err != nil {
		return err
	}

	asset := MetadataAsset{
		ID:            id,
		PatientID:     patientID,
		AssetID:       assetID,
		ZKPProof:      zkpProof,
		Name:          name,
		Value:         value,
		Version:       version,
		Owner:         owner,
		Rights:        rights,
		TermsOfAccess: termsOfAccess,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		CreatedBy:     createdBy,
		UpdatedBy:     updatedBy,
	}
	return putAsset(ctx, asset)
}

// GetAllMetadataFromNetwork returns every metadata asset on the ledger as a JSON array.
func (c *MetadataContract) GetAllMetadataFromNetwork(ctx contractapi.TransactionContextInterface) ([]MetadataAsset, error) {
	iter, err := ctx.GetStub().GetStateByRange("metadata_", "metadata_~")
	if err != nil {
		return nil, fmt.Errorf("failed to get state range: %w", err)
	}
	defer iter.Close()

	var results []MetadataAsset
	for iter.HasNext() {
		kv, err := iter.Next()
		if err != nil {
			return nil, err
		}
		var asset MetadataAsset
		if err := json.Unmarshal(kv.Value, &asset); err != nil {
			return nil, fmt.Errorf("failed to unmarshal asset: %w", err)
		}
		results = append(results, asset)
	}
	return results, nil
}

// GetMetadataById returns the metadata asset with the given numeric ID.
func (c *MetadataContract) GetMetadataById(
	ctx contractapi.TransactionContextInterface,
	idStr string,
) (*MetadataAsset, error) {
	id, err := mustParseUint64(idStr, "id")
	if err != nil {
		return nil, err
	}
	return getAsset(ctx, id)
}

// UpdateMetadataById replaces the mutable fields of an existing asset.
// Arguments order (all strings):
//
//	id, zkpProof, name, value, version,
//	owner, rights, termsOfAccess, updatedAt, updatedBy
func (c *MetadataContract) UpdateMetadataById(
	ctx contractapi.TransactionContextInterface,
	idStr string,
	zkpProof string,
	name string,
	value string,
	version string,
	owner string,
	rights string,
	termsOfAccess string,
	updatedAt string,
	updatedBy string,
) error {
	id, err := mustParseUint64(idStr, "id")
	if err != nil {
		return err
	}
	asset, err := getAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.ZKPProof = zkpProof
	asset.Name = name
	asset.Value = value
	asset.Version = version
	asset.Owner = owner
	asset.Rights = rights
	asset.TermsOfAccess = termsOfAccess
	asset.UpdatedAt = updatedAt
	asset.UpdatedBy = updatedBy

	return putAsset(ctx, *asset)
}

// DeleteMetadataById removes the metadata asset with the given ID from the ledger.
func (c *MetadataContract) DeleteMetadataById(
	ctx contractapi.TransactionContextInterface,
	idStr string,
) error {
	id, err := mustParseUint64(idStr, "id")
	if err != nil {
		return err
	}
	// Verify it exists first so we return a meaningful error.
	if _, err := getAsset(ctx, id); err != nil {
		return err
	}
	return ctx.GetStub().DelState(assetKey(id))
}

// GetMetadataAuditoryById returns the full history (audit trail) for the given asset ID.
func (c *MetadataContract) GetMetadataAuditoryById(
	ctx contractapi.TransactionContextInterface,
	idStr string,
) ([]HistoryEntry, error) {
	id, err := mustParseUint64(idStr, "id")
	if err != nil {
		return nil, err
	}

	iter, err := ctx.GetStub().GetHistoryForKey(assetKey(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get history for asset %d: %w", id, err)
	}
	defer iter.Close()

	var history []HistoryEntry
	for iter.HasNext() {
		mod, err := iter.Next()
		if err != nil {
			return nil, err
		}
		entry := HistoryEntry{
			TxID:      mod.TxId,
			Timestamp: time.Unix(mod.Timestamp.Seconds, int64(mod.Timestamp.Nanos)).UTC().Format(time.RFC3339),
			IsDelete:  mod.IsDelete,
		}
		if !mod.IsDelete {
			if err := json.Unmarshal(mod.Value, &entry.Value); err != nil {
				return nil, fmt.Errorf("failed to unmarshal history entry: %w", err)
			}
		}
		history = append(history, entry)
	}
	return history, nil
}

// ────────────────────────────────────────────────────────────
//  Entry point
// ────────────────────────────────────────────────────────────

func main() {
	cc, err := contractapi.NewChaincode(&MetadataContract{})
	if err != nil {
		panic(fmt.Sprintf("error creating MetadataContract chaincode: %v", err))
	}
	if err := cc.Start(); err != nil {
		panic(fmt.Sprintf("error starting MetadataContract chaincode: %v", err))
	}
}
