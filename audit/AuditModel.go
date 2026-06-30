package audit

import "time"

// AuditAction enumerates the operations that can be audited.
type AuditAction string

const (
	ActionCreate AuditAction = "CREATE"
	ActionUpdate AuditAction = "UPDATE"
	ActionDelete AuditAction = "DELETE"
	ActionList   AuditAction = "LIST"
)

// AuditModel represents a single audit entry that records an operation
// performed on a domain entity (e.g. Metadata).
type AuditModel struct {
	// ID is an auto-incremented sequence number assigned at insertion time.
	ID uint64 `json:"id"`

	// Entity identifies the domain object being audited (e.g. "Metadata").
	Entity string `json:"entity"`

	// EntityID is the primary-key value of the affected record.
	// Empty string means the operation targeted the entire collection (e.g. LIST).
	EntityID string `json:"entity_id,omitempty"`

	// Action is the operation that was performed (CREATE, UPDATE, DELETE, LIST).
	Action AuditAction `json:"action"`

	// Actor is the identity that triggered the operation (e.g. a username or
	// certificate CN). Defaults to "system" when no actor is available.
	Actor string `json:"actor,omitempty"`

	// OccurredAt is the UTC timestamp of when the operation was recorded.
	OccurredAt time.Time `json:"occurred_at"`

	// Details holds any additional context about the operation
	// (e.g. the payload snapshot or a short description).
	Details string `json:"details,omitempty"`
}

// newAuditModel is a constructor that pre-populates OccurredAt with the
// current UTC time.
func newAuditModel(entity, entityID string, action AuditAction, actor, details string) AuditModel {
	return AuditModel{
		Entity:     entity,
		EntityID:   entityID,
		Action:     action,
		Actor:      actor,
		OccurredAt: time.Now().UTC(),
		Details:    details,
	}
}