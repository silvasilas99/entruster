package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/fabric"
)

// TODO:    Bring the function RegisterMetadata from the fabric package to this package, and make it call the fabric
//          function. This way, the controller only calls the service layer, and the service layer is responsible for
//          calling the fabric. This will improve the separation of concerns and make the code more maintainable. The
//          same should be done for all other functions that interact with the fabric.

func RegisterMetadata(contract *client.Contract, req MetadataModel) error {
	// Validate required fields first — if invalid, Fabric never gets called
	if req.ID == "" {
		return fmt.Errorf("ID is required")
	}
	// Validation passed — call fabric
	return fabric.RegisterMetadataOnNetwork(
		contract,
		req.ID,
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
}

func GetAllMetadata(contract *client.Contract) ([]MetadataModel, error) {
	result, err := fabric.GetAllMetadataFromNetwork(contract)
	if err != nil {
		return nil, err
	}
	var metadataList []MetadataModel
	err = json.Unmarshal(result, &metadataList)
	if err != nil {
		return nil, fmt.Errorf("metadataService.GetAllMetadata: Failed to unmarshal metadata list: %w", err)
	}
	return metadataList, nil
}