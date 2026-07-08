package metadata

import "time"

type MetadataModel struct {
	ID            	uint64 `json:"id,omitempty"`
	PatientID     	uint64 `json:"patient_id"`
	AssetID       	uint64 `json:"asset_id"`
	Name          	string `json:"name"`
	Value         	string `json:"value"`
	Version       	string `json:"version"`
	Owner         	string `json:"owner"`
	Rights        	string `json:"rights"`
	TermsOfAccess 	string `json:"terms_of_access"`
	CreatedAt 		string  `json:"created_at,omitempty"`
	UpdatedAt 		string  `json:"updated_at,omitempty"`
	DeletedAt 		string `json:"deleted_at,omitempty"`
	CreatedBy 		string  `json:"created_by"`
	UpdatedBy 		string  `json:"updated_by"`
	DeletedBy 		string  `json:"deleted_by"`
}

func (m *MetadataModel) BeforeCreate() {
	m.CreatedAt = nowRFC3339()
}

func (m *MetadataModel) BeforeUpdate() {
	m.UpdatedAt = nowRFC3339()
}

func (m *MetadataModel) BeforeDelete() {
	m.DeletedAt = nowRFC3339()
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}