package metadata

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/chaincode"
)

type MetadataService struct {
	chaincodeQuery *chaincode.ChaincodeQuery
	metadataObserver *MetadataObserver
}

func NewMetadataService(contract *client.Contract, observer *MetadataObserver) *MetadataService {
	query := chaincode.NewChaincodeQuery(contract)

	return &MetadataService{
		chaincodeQuery: query,
		metadataObserver: observer,
	}
}

func (m *MetadataService) RegisterMetadata(metadataDTO MetadataDTO) error {
	fmt.Printf("[MetadataService][RegisterMetadata]: Storing metadata on blockchain: %+v\n", metadataDTO)

	if metadataDTO.CreatedAt == "" {
		metadataDTO.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	payload, err := json.Marshal(metadataDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	result, err := m.chaincodeQuery.StoreOnChain(payload)
	if err != nil {
		return fmt.Errorf("[MetadataService][RegisterMetadata]: Failed to submit transaction: %w", err)
	}

	id := string(result)
	fmt.Printf("*** Transaction committed successfully. Generated ID: %s\n", id)

	if m.metadataObserver != nil {
		// Just a placeholder call, actual method depends on observer implementation
		// m.metadataObserver.OnCreate(id, metadataDTO)
	}

	return nil
}

func (m *MetadataService) GetAllMetadata(filters []string) ([]MetadataModel, error) {
	fmt.Println("--> Evaluate Transaction: GetAllMetadata")
	result, err := m.chaincodeQuery.GetAllOnChain()
	if err != nil {
		return nil, fmt.Errorf("metadata.GetAllMetadata: failed to evaluate transaction: %w", err)
	}
	var list []MetadataModel
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, fmt.Errorf("metadata.GetAllMetadata: failed to unmarshal response: %w", err)
	}

	if m.metadataObserver != nil {
		// m.metadataObserver.OnList(len(list))
	}
	return list, nil
}

func (m *MetadataService) GetMetadataByID(id string) (*MetadataModel, error) {
	fmt.Printf("--> Evaluate Transaction: GetMetadataByID | ID: %s\n", id)
	result, err := m.chaincodeQuery.GetMetadataById(id)
	if err != nil {
		return nil, fmt.Errorf("metadata.GetMetadataByID: failed to evaluate transaction: %w", err)
	}
	var model MetadataModel
	if err := json.Unmarshal(result, &model); err != nil {
		return nil, fmt.Errorf("metadata.GetMetadataByID: failed to unmarshal response: %w", err)
	}
	return &model, nil
}

func (m *MetadataService) UpdateMetadataByID(id string, req MetadataModel) error {
	fmt.Printf("--> Submit Transaction: UpdateMetadataByID | ID: %s\n", id)

	if req.UpdatedAt == "" {
		req.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := m.chaincodeQuery.UpdateOnChain(
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

	return nil
}

func (m *MetadataService) DeleteMetadataByID(id string) error {
	fmt.Printf("--> Submit Transaction: DeleteMetadataByID | ID: %s\n", id)
	deletedAt := time.Now().UTC().Format(time.RFC3339)
	_, err := m.chaincodeQuery.DeleteOnChain("DeleteMetadataById", id, deletedAt)
	if err != nil {
		return fmt.Errorf("metadata.DeleteMetadataByID: failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")

	return nil
}
