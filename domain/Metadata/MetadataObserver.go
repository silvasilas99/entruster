package metadata

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/silvasilas99/entruster/audit"
	"github.com/silvasilas99/entruster/elasticsearch"
)

const entityName = "Metadata"
const esIndexName = "metadata"

// MetadataObserver listens to MetadataModel lifecycle events and delegates
// the creation of audit entries to the injected AuditService, as well as
// indexing in Elasticsearch.
type MetadataObserver struct {
	auditService *audit.AuditService
	elasticSvc   *elasticsearch.ElasticService
}

// NewMetadataObserver constructs a MetadataObserver.
func NewMetadataObserver(svc *audit.AuditService, es *elasticsearch.ElasticService) *MetadataObserver {
	if svc == nil {
		panic("metadata.NewMetadataObserver: auditService must not be nil")
	}
	if es == nil {
		panic("metadata.NewMetadataObserver: elasticSvc must not be nil")
	}
	return &MetadataObserver{auditService: svc, elasticSvc: es}
}

// OnCreate is fired after a Metadata record is successfully created.
func (o *MetadataObserver) OnCreate(id string, req MetadataModel) {
	details := fmt.Sprintf("patient_id=%s asset_id=%s name=%q", req.PatientID, req.AssetID, req.Name)
	o.auditService.Record(entityName, id, audit.ActionCreate, req.CreatedBy, details)

	req.ID, _ = strconv.ParseUint(id, 10, 64)
	err := o.elasticSvc.IndexDocument(context.Background(), esIndexName, id, req)
	if err != nil {
		log.Printf("Failed to index created metadata %s in Elasticsearch: %v", id, err)
	}
}

// OnUpdate is fired after a Metadata record is successfully updated.
func (o *MetadataObserver) OnUpdate(id string, req MetadataModel) {
	details := fmt.Sprintf("name=%q version=%q updated_by=%q", req.Name, req.Version, req.UpdatedBy)
	o.auditService.Record(entityName, id, audit.ActionUpdate, req.UpdatedBy, details)

	req.ID, _ = strconv.ParseUint(id, 10, 64)
	err := o.elasticSvc.IndexDocument(context.Background(), esIndexName, id, req)
	if err != nil {
		log.Printf("Failed to update metadata %s in Elasticsearch: %v", id, err)
	}
}

// OnDelete is fired after a Metadata record is successfully deleted.
func (o *MetadataObserver) OnDelete(id string) {
	o.auditService.Record(entityName, id, audit.ActionDelete, "system", "soft-delete committed to ledger")

	deletedAt := time.Now().UTC().Format(time.RFC3339)
	err := o.elasticSvc.UpdateDocument(context.Background(), esIndexName, id, map[string]interface{}{
		"deleted_at": deletedAt,
	})
	if err != nil {
		log.Printf("Failed to mark metadata %s as deleted in Elasticsearch: %v", id, err)
	}
}

// OnList is fired after a successful listing of all Metadata records.
func (o *MetadataObserver) OnList(count int) {
	details := fmt.Sprintf("returned %d record(s)", count)
	o.auditService.Record(entityName, "", audit.ActionList, "system", details)
}
