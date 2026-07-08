package chaincode

type TransferableDTO struct {
	ID       	  	uint64 `json:"id"`
	CreatedAt    	string `json:"created_at"`
	UpdatedAt     	string `json:"updated_at,omitempty"`
	DeletedAt     	string `json:"deleted_at,omitempty"`
	CreatedBy     	string `json:"created_by"`
	UpdatedBy     	string `json:"updated_by,omitempty"`
	DeletedBy     	string `json:"deleted_by,omitempty"`
}
