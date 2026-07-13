package metadata

import (
	"github.com/silvasilas99/entruster/app/core/chaincode"
)

type MetadataDTO struct {
	chaincode.TransferableDTO
	PatientID     	uint64 		`json:"patient_id"`
	AssetID       	uint64 		`json:"asset_id"`
	Category      	string 		`json:"category"`
	ResourceType	string 		`json:"resource_type"`
	PrimitiveType	string 		`json:"primitive_type"`
	Name          	string 		`json:"name"`
	Value         	string 		`json:"value"`
	Version       	string 		`json:"version"`
	Source         	string 		`json:"source"`
	Rights        	string 		`json:"rights"`
	TermsOfAccess 	string 		`json:"terms_of_access"`
}

func NewMetadataDTO(
	id uint64,
	patientID uint64,
	assetID uint64,
	category string,
	resourceType string,
	primitiveType string,
	name string,
	value string,
	version string,
	source string,
	rights string,
	termsOfAccess string,
	createdAt string,
	createdBy string,
	updatedAt string,
	updatedBy string,
	deletedAt string,
	deletedBy string,
) MetadataDTO {
	return MetadataDTO{
		TransferableDTO: chaincode.TransferableDTO{
			ID:	id,
			CreatedAt: createdAt,
			CreatedBy: createdBy,
			UpdatedAt: updatedAt,
			UpdatedBy: updatedBy,
			DeletedAt: deletedAt,
			DeletedBy: deletedBy,
		},
		PatientID:     	patientID,
		AssetID:       	assetID,
		Category:      	category,
		ResourceType: 	resourceType,
		PrimitiveType: 	primitiveType,
		Name:          	name,
		Value:         	value,
		Version:       	version,
		Source:         source,
		Rights:        	rights,
		TermsOfAccess: 	termsOfAccess
	}
}

func ToModel(dto MetadataDTO) MetadataModel{
	return MetadataModel{
		ID: 				dto.TransferableDTO.ID,
		PatientID: 			dto.PatientID,
		AssetID: 			dto.AssetID,
		Category: 			dto.Category,
		ResourceType: 		dto.ResourceType,
		PrimitiveType: 		dto.PrimitiveType,
		Name: 				dto.Name,
		Value: 				dto.Value,
		Version: 			dto.Version,
		Source: 			dto.Source,
		Rights: 			dto.Rights,
		TermsOfAccess: 		dto.TermsOfAccess,
		CreatedAt:     		dto.TransferableDTO.CreatedAt,
		CreatedBy:     		dto.TransferableDTO.CreatedBy,
		UpdatedAt:     		dto.TransferableDTO.UpdatedAt,
		UpdatedBy:     		dto.TransferableDTO.UpdatedBy,
		DeletedAt:     		dto.TransferableDTO.DeletedAt,
		DeletedBy:     		dto.TransferableDTO.DeletedBy
	}
}

