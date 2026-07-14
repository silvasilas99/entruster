package thirdpartyappevent

import "time"

type EventTrigged

const (
	ActionInsert 	Action = "INSERT"
	ActionUpdate 	Action = "UPDATE"
	ActionDelete 	Action = "DELETE"
	ActionRead 		Action = "READ"
)

type ThirdPartyAppEventModel struct {
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
	CreatedAt     			string   			`json:"created_at"`
	UpdatedAt 				string 				`json:"updated_at,omitempty"`
	DeletedAt 				string 				`json:"deleted_at,omitempty"`
	CreatedBy 				string 				`json:"created_by"`
	UpdatedBy 				string 				`json:"updated_by,omitempty"`
	DeletedBy 				string 				`json:"deleted_by,omitempty"`
}

func (m *ThirdPartyAppEventModel) BeforeCreate() {
	m.CreatedAt = nowRFC3339()
}

func (m *ThirdPartyAppEventModel) BeforeUpdate() {
	m.UpdatedAt = nowRFC3339()
}

func (m *ThirdPartyAppEventModel) BeforeDelete() {
	m.DeletedAt = nowRFC3339()
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
