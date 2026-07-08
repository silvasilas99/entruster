package audit

import "time"

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
	ActionList   ActionType = "LIST"
)

type AuditService struct {
	records []AuditModel
}

func NewAuditService() *AuditService {
	return &AuditService{
		records: make([]AuditModel, 0),
	}
}

func (s *AuditService) Record(entityName, entityID string, action ActionType, user, details string) {
	s.records = append(s.records, AuditModel{
		EntityName: entityName,
		EntityId:   entityID,
		Action:     string(action),
		Actor:      user,
		Details:    details,
		OccurredAt: time.Now().UTC(),
	})
}

func (s *AuditService) GetByEntityID(entityName, entityID string) []AuditModel {
	var result []AuditModel
	for _, r := range s.records {
		if r.EntityName == entityName && r.EntityId == entityID {
			result = append(result, r)
		}
	}
	return result
}