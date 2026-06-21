package metadata

import "encoding/json"

type MetadataModel struct {
	PatientID     string `json:"patient_id"`
	AssetID       string `json:"asset_id"`
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

func (m *MetadataModel) UnmarshalJSON(data []byte) error {
	type metadataAlias MetadataModel
	aux := struct {
		metadataAlias
		TermsOfAccessCamel string `json:"termsOfAccess"`
		CreatedAtCamel     string `json:"createdAt"`
		UpdatedAtCamel     string `json:"updatedAt"`
		CreatedByCamel     string `json:"createdBy"`
		UpdatedByCamel     string `json:"updatedBy"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	*m = MetadataModel(aux.metadataAlias)
	if m.TermsOfAccess == "" {
		m.TermsOfAccess = aux.TermsOfAccessCamel
	}
	if m.CreatedAt == "" {
		m.CreatedAt = aux.CreatedAtCamel
	}
	if m.UpdatedAt == "" {
		m.UpdatedAt = aux.UpdatedAtCamel
	}
	if m.CreatedBy == "" {
		m.CreatedBy = aux.CreatedByCamel
	}
	if m.UpdatedBy == "" {
		m.UpdatedBy = aux.UpdatedByCamel
	}

	return nil
}
