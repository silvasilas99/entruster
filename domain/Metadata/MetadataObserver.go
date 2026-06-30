package metadata

import (
	"fmt"

	"github.com/silvasilas99/entruster/audit"
)

const entityName = "Metadata"

// MetadataObserver listens to MetadataModel lifecycle events and delegates
// the creation of audit entries to the injected AuditService.
//
// Wire it once at startup and pass it into every contract function that must
// produce an audit trail.
type MetadataObserver struct {
	auditService *audit.AuditService
}

// NewMetadataObserver constructs a MetadataObserver backed by the given
// AuditService.  Panics if auditService is nil.
func NewMetadataObserver(svc *audit.AuditService) *MetadataObserver {
	if svc == nil {
		panic("metadata.NewMetadataObserver: auditService must not be nil")
	}
	return &MetadataObserver{auditService: svc}
}

// OnCreate is fired after a Metadata record is successfully created on the
// ledger.
func (o *MetadataObserver) OnCreate(req MetadataModel) {
	details := fmt.Sprintf("patient_id=%d asset_id=%d name=%q", req.PatientID, req.AssetID, req.Name)
	o.auditService.Record(entityName, "", audit.ActionCreate, req.CreatedBy, details)
}

// OnUpdate is fired after a Metadata record is successfully updated on the
// ledger.  id is the string representation of the asset's numeric ID.
func (o *MetadataObserver) OnUpdate(id string, req MetadataModel) {
	details := fmt.Sprintf("name=%q version=%q updated_by=%q", req.Name, req.Version, req.UpdatedBy)
	o.auditService.Record(entityName, id, audit.ActionUpdate, req.UpdatedBy, details)
}

// OnDelete is fired after a Metadata record is successfully deleted from the
// ledger.  id is the string representation of the asset's numeric ID.
func (o *MetadataObserver) OnDelete(id string) {
	o.auditService.Record(entityName, id, audit.ActionDelete, "system", "soft-delete committed to ledger")
}

// OnList is fired after a successful listing of all Metadata records.
// count is the number of items returned by the query.
func (o *MetadataObserver) OnList(count int) {
	details := fmt.Sprintf("returned %d record(s)", count)
	o.auditService.Record(entityName, "", audit.ActionList, "system", details)
}
