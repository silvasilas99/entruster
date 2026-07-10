package thirdpartyappevent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/silvasilas99/entruster/app/core/chaincode"
)

type ThirdPartyAppEventService struct {
	chaincodeQuery             *chaincode.ChaincodeQuery
	thirdPartyAppEventObserver *ThirdPartyAppEventObserver
}

func NewThirdPartyAppEventService(contract *client.Contract, observer *ThirdPartyAppEventObserver) *ThirdPartyAppEventService {
	query := chaincode.NewChaincodeQuery(contract)

	return &ThirdPartyAppEventService{
		chaincodeQuery:             query,
		thirdPartyAppEventObserver: observer,
	}
}

func (m *ThirdPartyAppEventService) RegisterThirdPartyAppEvent(thirdPartyAppEventDTO ThirdPartyAppEventDTO) error {
	fmt.Printf("[ThirdPartyAppEventService][RegisterThirdPartyAppEvent]: Storing thirdpartyappevent on blockchain: %+v\n", thirdPartyAppEventDTO)

	if _, err := time.Parse(time.RFC3339, thirdPartyAppEventDTO.CreatedAt); err != nil {
		thirdPartyAppEventDTO.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	payload, err := json.Marshal(thirdPartyAppEventDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal thirdpartyappevent: %w", err)
	}

	result, err := m.chaincodeQuery.StoreThirdPartyAppEventOnChain(payload)
	if err != nil {
		return fmt.Errorf("[ThirdPartyAppEventService][RegisterThirdPartyAppEvent]: Failed to submit transaction: %w", err)
	}

	id := string(result)
	fmt.Printf("*** Transaction committed successfully. Generated ID: %s\n", id)

	if m.thirdPartyAppEventObserver != nil {
		m.thirdPartyAppEventObserver.OnCreate(id, ToModel(thirdPartyAppEventDTO))
	}

	return nil
}

func (m *ThirdPartyAppEventService) GetAllThirdPartyAppEvent(filters []string) ([]ThirdPartyAppEventModel, error) {
	fmt.Println("--> Evaluate Transaction: GetAllThirdPartyAppEvent")
	result, err := m.chaincodeQuery.GetAllThirdPartyAppEventsOnChain()
	if err != nil {
		return nil, fmt.Errorf("thirdpartyappevent.GetAllThirdPartyAppEvent: failed to evaluate transaction: %w", err)
	}
	var list []ThirdPartyAppEventModel
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, fmt.Errorf("thirdpartyappevent.GetAllThirdPartyAppEvent: failed to unmarshal response: %w", err)
	}

	if m.thirdPartyAppEventObserver != nil {
		// m.thirdPartyAppEventObserver.OnList(len(list), "system")
	}
	return list, nil
}

func (m *ThirdPartyAppEventService) GetThirdPartyAppEventByID(id string) (*ThirdPartyAppEventModel, error) {
	fmt.Printf("--> Evaluate Transaction: GetThirdPartyAppEventByID | ID: %s\n", id)
	result, err := m.chaincodeQuery.GetThirdPartyAppEventById(id)
	if err != nil {
		return nil, fmt.Errorf("thirdpartyappevent.GetThirdPartyAppEventByID: failed to evaluate transaction: %w", err)
	}
	var model ThirdPartyAppEventModel
	if err := json.Unmarshal(result, &model); err != nil {
		return nil, fmt.Errorf("thirdpartyappevent.GetThirdPartyAppEventByID: failed to unmarshal response: %w", err)
	}
	return &model, nil
}

func (m *ThirdPartyAppEventService) UpdateThirdPartyAppEventByID(id string, req ThirdPartyAppEventModel) error {
	fmt.Printf("--> Submit Transaction: UpdateThirdPartyAppEventByID | ID: %s\n", id)

	if _, err := time.Parse(time.RFC3339, req.UpdatedAt); err != nil {
		req.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	_, err := m.chaincodeQuery.UpdateOnChain(
		"UpdateThirdPartyAppEventById",
		id,
		req.EventType,
		req.EventData,
		req.Description,
		req.UpdatedAt,
		req.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("thirdpartyappevent.UpdateThirdPartyAppEventByID: failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")

	if m.thirdPartyAppEventObserver != nil {
		m.thirdPartyAppEventObserver.OnUpdate(id, req)
	}

	return nil
}

func (m *ThirdPartyAppEventService) DeleteThirdPartyAppEventByID(id string, deletedBy string) error {
	fmt.Printf("--> Submit Transaction: DeleteThirdPartyAppEventByID | ID: %s | DeletedBy: %s\n", id, deletedBy)
	deletedAt := time.Now().UTC().Format(time.RFC3339)
	_, err := m.chaincodeQuery.DeleteOnChain("DeleteThirdPartyAppEventById", id, deletedAt)
	if err != nil {
		return fmt.Errorf("thirdpartyappevent.DeleteThirdPartyAppEventByID: failed to submit transaction: %w", err)
	}
	fmt.Println("*** Transaction committed successfully")

	if m.thirdPartyAppEventObserver != nil {
		m.thirdPartyAppEventObserver.OnDelete(id, deletedBy)
	}

	return nil
}
