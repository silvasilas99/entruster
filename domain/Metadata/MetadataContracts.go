package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// CreateMetadata submits a RegisterMetadataOnNetwork transaction to the ledger.
// The chaincode auto-generates the asset ID using an internal counter.
func CreateMetadata(contract *client.Contract, req MetadataModel) error {
	fmt.Printf("--> Submit Transaction: CreateMetadata | PatientID: %s AssetID: %s\n", req.PatientID, req.AssetID)
	_, err := contract.SubmitTransaction(
		"RegisterMetadataOnNetwork",
		req.PatientID,     // uint64 — chaincode converts from string
		req.AssetID,       // uint64 — chaincode converts from string
		req.ZKPProof,
		req.Name,
		req.Value,
		req.Version,
		req.Owner,
		req.Rights,
		req.TermsOfAccess,
		req.CreatedAt,
		req.UpdatedAt,
		req.CreatedBy,
		req.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("metadata.CreateMetadata: failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")
	return nil
}

// GetAllMetadata evaluates GetAllMetadataFromNetwork and returns the parsed slice.
func GetAllMetadata(contract *client.Contract) ([]MetadataModel, error) {
	fmt.Println("--> Evaluate Transaction: GetAllMetadata")
	result, err := contract.EvaluateTransaction("GetAllMetadataFromNetwork")
	if err != nil {
		return nil, fmt.Errorf("metadata.GetAllMetadata: failed to evaluate transaction: %w", err)
	}
	var list []MetadataModel
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, fmt.Errorf("metadata.GetAllMetadata: failed to unmarshal response: %w", err)
	}
	return list, nil
}

// GetMetadataByID evaluates GetMetadataById and returns the matching asset.
func GetMetadataByID(contract *client.Contract, id string) (*MetadataModel, error) {
	fmt.Printf("--> Evaluate Transaction: GetMetadataByID | ID: %s\n", id)
	result, err := contract.EvaluateTransaction("GetMetadataById", id)
	if err != nil {
		return nil, fmt.Errorf("metadata.GetMetadataByID: failed to evaluate transaction: %w", err)
	}
	var m MetadataModel
	if err := json.Unmarshal(result, &m); err != nil {
		return nil, fmt.Errorf("metadata.GetMetadataByID: failed to unmarshal response: %w", err)
	}
	return &m, nil
}

// UpdateMetadataByID submits an UpdateMetadataById transaction.
// id is the asset ID (uint64 as string); req carries the new field values.
func UpdateMetadataByID(contract *client.Contract, id string, req MetadataModel) error {
	fmt.Printf("--> Submit Transaction: UpdateMetadataByID | ID: %s\n", id)
	_, err := contract.SubmitTransaction(
		"UpdateMetadataById",
		id,
		req.ZKPProof,
		req.Name,
		req.Value,
		req.Version,
		req.Owner,
		req.Rights,
		req.TermsOfAccess,
		req.UpdatedAt,
		req.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("metadata.UpdateMetadataByID: failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")
	return nil
}

// DeleteMetadataByID submits a DeleteMetadataById transaction.
func DeleteMetadataByID(contract *client.Contract, id string) error {
	fmt.Printf("--> Submit Transaction: DeleteMetadataByID | ID: %s\n", id)
	_, err := contract.SubmitTransaction("DeleteMetadataById", id)
	if err != nil {
		return fmt.Errorf("metadata.DeleteMetadataByID: failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")
	return nil
}

// GetMetadataAuditoryByID evaluates GetMetadataAuditoryById and returns
// the full history (audit trail) for the given asset ID.
func GetMetadataAuditoryByID(contract *client.Contract, id string) ([]MetadataHistoryEntry, error) {
	fmt.Printf("--> Evaluate Transaction: GetMetadataAuditoryByID | ID: %s\n", id)
	result, err := contract.EvaluateTransaction("GetMetadataAuditoryById", id)
	if err != nil {
		return nil, fmt.Errorf("metadata.GetMetadataAuditoryByID: failed to evaluate transaction: %w", err)
	}
	var history []MetadataHistoryEntry
	if err := json.Unmarshal(result, &history); err != nil {
		return nil, fmt.Errorf("metadata.GetMetadataAuditoryByID: failed to unmarshal response: %w", err)
	}
	return history, nil
}
