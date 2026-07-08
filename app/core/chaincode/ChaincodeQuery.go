package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type ChaincodeQuery struct {
	ChaincodeService *ChaincodeService
	CounterKey string
	Contract *client.Contract
}

func NewChaincodeQuery(contract *client.Contract) *ChaincodeQuery {
	// Not instantiating the service here, assuming it's injected or using the contract directly.
	return &ChaincodeQuery{
		CounterKey: "_metadata_id",
		Contract: contract,
	}
}

// registerPayload mirrors the JSON fields produced by MetadataService.RegisterMetadata.
// It is used only to unpack the payload before forwarding individual args to the chaincode.
type registerPayload struct {
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

// StoreOnChain submits the RegisterMetadataOnNetwork chaincode transaction.
//
// The chaincode function expects 13 individual string arguments (positional), NOT a JSON blob:
//
//	patientID, assetID, zkpProof, name, value, version,
//	owner, rights, termsOfAccess, createdAt, updatedAt, createdBy, updatedBy
//
// This method unmarshals the payload produced by MetadataService and forwards each
// field as a separate argument so the parameter count matches the chaincode signature.
func (c *ChaincodeQuery) StoreOnChain(
	payload []byte,
) (string, error) {
	var p registerPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return "", fmt.Errorf("StoreOnChain: failed to unmarshal payload: %w", err)
	}

	result, err := c.Contract.SubmitTransaction(
		"RegisterMetadataOnNetwork",
		strconv.FormatUint(p.PatientID, 10),
		strconv.FormatUint(p.AssetID, 10),
		p.ZKPProof,
		p.Name,
		p.Value,
		p.Version,
		p.Owner,
		p.Rights,
		p.TermsOfAccess,
		p.CreatedAt,
		p.UpdatedAt,
		p.CreatedBy,
		p.UpdatedBy,
	)
	if err != nil {
		return "", err
	}
	return string(result), nil
}


func (c *ChaincodeQuery) UpdateOnChain(
	transactionName string,
	args ...string,
) (string, error) {
	result, err := c.Contract.SubmitTransaction(transactionName, args...)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (c *ChaincodeQuery) DeleteOnChain(
	transactionName string,
	args ...string,
) (string, error) {
	result, err := c.Contract.SubmitTransaction(transactionName, args...)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func (c *ChaincodeQuery) GetAllOnChain() ([]byte, error) {
	result, err := c.Contract.EvaluateTransaction("GetAllMetadataFromNetwork")
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *ChaincodeQuery) GetMetadataById(
	id string,
) ([]byte, error) {
	result, err := c.Contract.EvaluateTransaction("GetMetadataById", id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
