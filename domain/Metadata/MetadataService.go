package metadata

import (
    "fmt"
    "github.com/hyperledger/fabric-gateway/pkg/client"
    "github.com/silvasilas99/entruster/fabric"
    "github.com/silvasilas99/entruster/domain/metadata"
)

// TODO:    Bring the function RegisterMetadata from the fabric package to this package, and make it call the fabric
//          function. This way, the controller only calls the service layer, and the service layer is responsible for
//          calling the fabric. This will improve the separation of concerns and make the code more maintainable. The
//          same should be done for all other functions that interact with the fabric.

func RegisterMetadata(contract *client.Contract, req metadata.MetadataModel) error {
    // Validate required fields first — if invalid, Fabric never gets called
    if req.id == "" {
        return fmt.Errorf("ID is required")
    }
    // Validation passed — call fabric
    return fabric.RegisterMetadataOnNetwork(contract, req)
}

func GetAllMetadata(contract *client.Contract) ([]metadata.MetadataModel, error) {
    return fabric.GetAllMetadataFromNetwork(contract)
}