package thirdpartyappevent

import (
	"github.com/silvasilas99/entruster/app/core/chaincode"
)

type ThirdPartyAppEventDTO struct {
	chaincode.TransferableDTO
	AppID       string `json:"app_id"`
	EventType   string `json:"event_type"`
	EventData   string `json:"event_data"`
	Description string `json:"description"`
}

func NewThirdPartyAppEventDTO(
	id uint64,
	appID string,
	eventType string,
	eventData string,
	description string,
	createdAt string,
	createdBy string,
	updatedAt string,
	updatedBy string,
	deletedAt string,
	deletedBy string,
) ThirdPartyAppEventDTO {
	return ThirdPartyAppEventDTO{
		TransferableDTO: chaincode.TransferableDTO{
			ID:        id,
			CreatedAt: createdAt,
			CreatedBy: createdBy,
			UpdatedAt: updatedAt,
			UpdatedBy: updatedBy,
			DeletedAt: deletedAt,
			DeletedBy: deletedBy,
		},
		AppID:       appID,
		EventType:   eventType,
		EventData:   eventData,
		Description: description,
	}
}

func ToModel(dto ThirdPartyAppEventDTO) ThirdPartyAppEventModel {
	return ThirdPartyAppEventModel{
		ID:          dto.TransferableDTO.ID,
		AppID:       dto.AppID,
		EventType:   dto.EventType,
		EventData:   dto.EventData,
		Description: dto.Description,

		CreatedAt: dto.TransferableDTO.CreatedAt,
		CreatedBy: dto.TransferableDTO.CreatedBy,
		UpdatedAt: dto.TransferableDTO.UpdatedAt,
		UpdatedBy: dto.TransferableDTO.UpdatedBy,
		DeletedAt: dto.TransferableDTO.DeletedAt,
		DeletedBy: dto.TransferableDTO.DeletedBy,
	}
}
