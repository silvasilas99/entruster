package metadata

type MetadataModel struct {
	ID            string `json:"id"`
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