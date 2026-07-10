package thirdpartyappevent

import "time"

type ThirdPartyAppEventModel struct {
	ID          uint64 `json:"id,omitempty"`
	AppID       string `json:"app_id"`
	EventType   string `json:"event_type"`
	EventData   string `json:"event_data"`
	Description string `json:"description"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	DeletedAt string `json:"deleted_at,omitempty"`
	CreatedBy string `json:"created_by"`
	UpdatedBy string `json:"updated_by,omitempty"`
	DeletedBy string `json:"deleted_by,omitempty"`
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
