package thirdpartyappevent

import (
	"github.com/silvasilas99/entruster/app/core/chaincode"
)

type ThirdPartyAppEventDTO struct {
	chaincode.TransferableDTO
	ID 						int64       		`json:"id"`
	EntityID      			int64      			`json:"entity_id"`
	EntityName    			string      		`json:"entity_name"`
	EntityBeforeAction   	json.RawMessage 	`json:"entity_before_action,omitempty"`
	EntityAfterAction    	json.RawMessage 	`json:"entity_after_action,omitempty"`
	PlatformName			string			    `json:"platform_name,omitempty"`		
	ActionType    			AuditAction 		`json:"action_type"`
	ActorID       			*string     		`json:"actor_id,omitempty"`
	IPAddress     			*string     		`json:"ip_address,omitempty"`
	UserAgent     			*string     		`json:"user_agent,omitempty"`
}

func NewThirdPartyAppEventDTO(
	id int64,
	entityID int64,
	entityName string,
	entityBeforeAction json.RawMessage,
	entityAfterAction json.RawMessage,
	actionType string,
	actorID int64,
	platformName string,
	ipAddress string,
	userAgent string
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
			DeletedBy: deletedBy
		},
		EntityID: 				entityId,
		EntityName: 			entityName,
		EntityBeforeAction: 	entityBeforeAction,
		EntityAfterAction: 		entityAfterAction,
		ActionType: 			actionType,
		ActorID: 				actorID,
		PlatformName: 			platformName,
		IPAddress: 				ipAddress,
		UserAgent: 				userAgent
	}
}

// TODO: Avaliar mudança desse getter para ThirdPartyAppEventModel.go
func ToModel(dto ThirdPartyAppEventDTO) ThirdPartyAppEventModel {
	return ThirdPartyAppEventModel{
		ID:          			dto.TransferableDTO.ID,
		EntityID: 				dto.EntityID,
		EntityName: 			dto.EntityName,
		EntityBeforeAction: 	dto.EntityBeforeAction,
		EntityAfterAction: 		dto.EntityAfterAction,
		PlatformName: 			dto.PlatformName,
		ActionType: 			dto.ActionType,
		ActorID: 				dto.ActorID,
		IPAddress: 				dto.IPAddress,
		UserAgent: 				dto.UserAgent,
		CreatedAt: 				dto.TransferableDTO.CreatedAt,
		CreatedBy: 				dto.TransferableDTO.CreatedBy,
		UpdatedAt: 				dto.TransferableDTO.UpdatedAt,
		UpdatedBy: 				dto.TransferableDTO.UpdatedBy,
		DeletedAt: 				dto.TransferableDTO.DeletedAt,
		DeletedBy: 				dto.TransferableDTO.DeletedBy
	}
}
