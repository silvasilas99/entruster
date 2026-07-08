package metadata

import (
	"github.com/silvasilas99/entruster/app/core/chaincode"
)

type MetadataDTO struct {
	chaincode.TransferableDTO
	PatientID     uint64 `json:"patient_id"`
	AssetID       uint64 `json:"asset_id"`
	ZKPProof      string `json:"zkp_proof"`
	Name          string `json:"name"`
	Value         string `json:"value"`
	Version       string `json:"version"`
	Owner         string `json:"owner"`
	Rights        string `json:"rights"`
	TermsOfAccess string `json:"terms_of_access"`
}

func NewMetadataDTO(
	id uint64,
	patientID uint64,
	assetID uint64,
	zkpProof string,
	name string,
	value string,
	version string,
	owner string,
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
		PatientID:     patientID,
		AssetID:       assetID,
		ZKPProof:      zkpProof,
		Name:          name,
		Value:         value,
		Version:       version,
		Owner:         owner,
		Rights:        rights,
		TermsOfAccess: termsOfAccess,
	}
}

func ToModel(dto MetadataDTO) MetadataModel{
	return MetadataModel{
		ID: dto.TransferableDTO.ID,
		PatientID:     dto.PatientID,
		AssetID:       dto.AssetID,
		ZKPProof:      dto.ZKPProof,
		Name:          dto.Name,
		Value:         dto.Value,
		Version:       dto.Version,
		Owner:         dto.Owner,
		Rights:        dto.Rights,
		TermsOfAccess: dto.TermsOfAccess,
		CreatedAt:     dto.TransferableDTO.CreatedAt,
		CreatedBy:     dto.TransferableDTO.CreatedBy,
		UpdatedAt:     dto.TransferableDTO.UpdatedAt,
		UpdatedBy:     dto.TransferableDTO.UpdatedBy,
		DeletedAt:     dto.TransferableDTO.DeletedAt,
		DeletedBy:     dto.TransferableDTO.DeletedBy,
	}
}

