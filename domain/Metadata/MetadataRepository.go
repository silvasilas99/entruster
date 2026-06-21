package metadata

import (
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// TODO: Create a separate package for the chaincode SDK connection logic. This will help to keep
// the code organized and maintainable. The chaincode SDK connection logic is currently mixed with
// the controller logic, which can make it difficult to understand and maintain.

func RegisterMetadataOnNetwork(contract *client.Contract, req MetadataModel) error {
	fmt.Printf("--> Submit Transaction: RegisterMetadataOnNetwork | PatientID: %s AssetID: %s\n", req.PatientID, req.AssetID)
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
		return fmt.Errorf("metadata.RegisterMetadataOnNetwork: Internal error. Failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")
	return nil
}

func GetAllMetadataFromNetwork(contract *client.Contract) ([]byte, error) {
	fmt.Printf("--> Evaluate Transaction: GetAllMetadataFromNetwork\n")
	result, err := contract.EvaluateTransaction("GetAllMetadataFromNetwork")
	if err != nil {
		return nil, fmt.Errorf("metadata.GetAllMetadataFromNetwork: Internal error. Failed to evaluate transaction: %w", err)
	}
	return result, nil
}

