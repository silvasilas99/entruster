package audit

import "time"

type AuditModel struct {
	ID 						uint64 			`json:"id"`
	EntityName 				string 			`json:"entity_name"`
	EntityId 				string 			`json:"entity_id"`
	EntityDataBeforeAction 	string 			`json:"entity_data_before_action,omitempty"`
	EntityDataAfterAction 	string 			`json:"entity_data_after_action,omitempty"`
	Actor 					string 			`json:"actor"`
	RequestData 			string 			`json:"request_data"`
	Action 					string 			`json:"action"`
	Details 				string 			`json:"details"`
	OccurredAt 				time.Time 		`json:"occurred_at"`
}

/**
 * 	NewAuditModel creates a new instance of AuditModel with the provided parameters, like a constructor.
 *	It initializes the fields of the AuditModel struct with the given values.
 */
func NewAuditModel(
	entityID string,
	entityName string,
	entityDataBeforeAction string,
	entityDataAfterAction string,
	actor string,
	requestData string,
	action string,
	details string,
) AuditModel {
	return AuditModel{
		EntityId: entityID,
		EntityName: entityName,
		EntityDataBeforeAction: entityDataBeforeAction,
		EntityDataAfterAction: entityDataAfterAction,
		Actor: actor,
		RequestData: requestData,
		Details: details,
		OccurredAt: time.Now().UTC(),
	}
}