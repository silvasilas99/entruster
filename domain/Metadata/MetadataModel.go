package metadata

type MetadataModel struct {
    id           	string `json:"id"`
    name         	string `json:"name"`
    value        	string `json:"value"`
    version        	string `json:"version"`
    owner 			string `json:"owner"`
    rights 			string `json:"rights"`
    terms_of_access string `json:"terms_of_access"`,
    created_at 		string `json:"created_at"`
    updated_at 		string `json:"updated_at"`
    created_by      string `json:"created_by"`
    updated_by 		string `json:"updated_by"`
}