package audit

import (
	"fmt"
	"sync"
)

// AuditService persists and exposes audit entries produced by domain observers.
// The current implementation uses an in-memory store; replace the storage
// back-end (e.g. a database repository) without changing the public interface.
type AuditService struct {
	mu      sync.RWMutex
	entries []AuditModel
	counter uint64
}

// NewAuditService returns a ready-to-use AuditService.
func NewAuditService() *AuditService {
	return &AuditService{
		entries: make([]AuditModel, 0),
	}
}

// Record creates and persists a new audit entry for the given operation.
//
//   - entity   — name of the domain object (e.g. "Metadata")
//   - entityID — primary key of the affected record, or "" for collection ops
//   - action   — one of ActionCreate, ActionUpdate, ActionDelete, ActionList
//   - actor    — identity of the requester; falls back to "system" if blank
//   - details  — optional free-text context (payload snippet, description, …)
func (s *AuditService) Record(entity, entityID string, action AuditAction, actor, details string) {
	if actor == "" {
		actor = "system"
	}

	entry := newAuditModel(entity, entityID, action, actor, details)

	s.mu.Lock()
	s.counter++
	entry.ID = s.counter
	s.entries = append(s.entries, entry)
	s.mu.Unlock()

	fmt.Printf("[AUDIT] #%d | %s | entity=%s id=%q actor=%s | %s\n",
		entry.ID, entry.Action, entry.Entity, entry.EntityID, entry.Actor, entry.OccurredAt.Format("2006-01-02T15:04:05Z"))
}

// GetAll returns a snapshot of every audit entry, ordered by insertion time.
func (s *AuditService) GetAll() []AuditModel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := make([]AuditModel, len(s.entries))
	copy(snapshot, s.entries)
	return snapshot
}

// GetByEntity returns all audit entries for a given entity type (e.g. "Metadata").
func (s *AuditService) GetByEntity(entity string) []AuditModel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []AuditModel
	for _, e := range s.entries {
		if e.Entity == entity {
			result = append(result, e)
		}
	}
	return result
}

// GetByEntityID returns all audit entries for a specific entity instance.
func (s *AuditService) GetByEntityID(entity, entityID string) []AuditModel {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []AuditModel
	for _, e := range s.entries {
		if e.Entity == entity && e.EntityID == entityID {
			result = append(result, e)
		}
	}
	return result
}