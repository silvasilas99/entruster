package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// CreateMetadata submits a RegisterMetadataOnNetwork transaction to the ledger.
// The chaincode auto-generates the asset ID using an internal counter.
// CreatedAt and UpdatedAt are set to the current UTC time if not already set.
// If observer is non-nil, OnCreate is fired after a successful commit.
func CreateMetadata(contract *client.Contract, req MetadataModel, observer *MetadataObserver) error {
	fmt.Printf("--> Submit Transaction: CreateMetadata | PatientID: %d AssetID: %d\n", req.PatientID, req.AssetID)

	now := nowRFC3339()
	if req.CreatedAt == "" {
		req.CreatedAt = now
	}
	if req.UpdatedAt == "" {
		req.UpdatedAt = now
	}

	result, err := contract.SubmitTransaction(
		"RegisterMetadataOnNetwork",
		fmt.Sprintf("%d", req.PatientID), // uint64 — chaincode converts from string
		fmt.Sprintf("%d", req.AssetID),   // uint64 — chaincode converts from string
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
	id := string(result)
	fmt.Printf("*** Transaction committed successfully. Generated ID: %s\n", id)

	if observer != nil {
		observer.OnCreate(id, req)
	}
	return nil
}

// GetAllMetadata evaluates GetAllMetadataFromNetwork and returns the parsed slice.
// If observer is non-nil, OnList is fired with the number of records returned.
func GetAllMetadata(contract *client.Contract, observer *MetadataObserver) ([]MetadataModel, error) {
	fmt.Println("--> Evaluate Transaction: GetAllMetadata")
	result, err := contract.EvaluateTransaction("GetAllMetadataFromNetwork")
	if err != nil {
		return nil, fmt.Errorf("metadata.GetAllMetadata: failed to evaluate transaction: %w", err)
	}
	var list []MetadataModel
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, fmt.Errorf("metadata.GetAllMetadata: failed to unmarshal response: %w", err)
	}

	if observer != nil {
		observer.OnList(len(list))
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
// UpdatedAt is set to the current UTC time if not already set.
// If observer is non-nil, OnUpdate is fired after a successful commit.
func UpdateMetadataByID(contract *client.Contract, id string, req MetadataModel, observer *MetadataObserver) error {
	fmt.Printf("--> Submit Transaction: UpdateMetadataByID | ID: %s\n", id)

	if req.UpdatedAt == "" {
		req.UpdatedAt = nowRFC3339()
	}

	_, err := contract.SubmitTransaction(
		"UpdateMetadataById",
		id,
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

	if observer != nil {
		observer.OnUpdate(id, req)
	}
	return nil
}

// DeleteMetadataByID submits a DeleteMetadataById transaction.
// deletedAt is set to the current UTC time and recorded on the ledger (soft-delete marker).
// If observer is non-nil, OnDelete is fired after a successful commit.
func DeleteMetadataByID(contract *client.Contract, id string, observer *MetadataObserver) error {
	fmt.Printf("--> Submit Transaction: DeleteMetadataByID | ID: %s\n", id)
	deletedAt := nowRFC3339()
	_, err := contract.SubmitTransaction("DeleteMetadataById", id, deletedAt)
	if err != nil {
		return fmt.Errorf("metadata.DeleteMetadataByID: failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")

	if observer != nil {
		observer.OnDelete(id)
	}
	return nil
}


