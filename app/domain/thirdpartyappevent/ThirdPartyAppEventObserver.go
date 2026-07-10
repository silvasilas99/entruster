package thirdpartyappevent

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/silvasilas99/entruster/app/core/audit"
	"github.com/silvasilas99/entruster/app/core/elasticsearch"
)

const entityName = "ThirdPartyAppEvent"
const esIndexName = "thirdpartyappevent"

// ThirdPartyAppEventObserver listens to ThirdPartyAppEventModel lifecycle events and delegates
// the creation of audit entries to the injected AuditService, as well as
// indexing in Elasticsearch.
type ThirdPartyAppEventObserver struct {
	auditService *audit.AuditService
	elasticSvc   *elasticsearch.ElasticService
}

// NewThirdPartyAppEventObserver constructs a ThirdPartyAppEventObserver.
func NewThirdPartyAppEventObserver(svc *audit.AuditService, es *elasticsearch.ElasticService) *ThirdPartyAppEventObserver {
	if svc == nil {
		panic("thirdpartyappevent.NewThirdPartyAppEventObserver: auditService must not be nil")
	}
	if es == nil {
		panic("thirdpartyappevent.NewThirdPartyAppEventObserver: elasticSvc must not be nil")
	}
	return &ThirdPartyAppEventObserver{auditService: svc, elasticSvc: es}
}

// OnCreate is fired after a ThirdPartyAppEvent record is successfully created.
func (o *ThirdPartyAppEventObserver) OnCreate(id string, req ThirdPartyAppEventModel) {
	details := fmt.Sprintf("app_id=%s event_type=%s description=%q", req.AppID, req.EventType, req.Description)
	o.auditService.Record(entityName, id, audit.ActionCreate, req.CreatedBy, details)

	req.ID, _ = strconv.ParseUint(id, 10, 64)
	err := o.elasticSvc.IndexDocument(context.Background(), esIndexName, id, req)
	if err != nil {
		log.Printf("Failed to index created thirdpartyappevent %s in Elasticsearch: %v", id, err)
	}
}

// OnUpdate is fired after a ThirdPartyAppEvent record is successfully updated.
func (o *ThirdPartyAppEventObserver) OnUpdate(id string, req ThirdPartyAppEventModel) {
	details := fmt.Sprintf("description=%q event_type=%q updated_by=%q", req.Description, req.EventType, req.UpdatedBy)
	o.auditService.Record(entityName, id, audit.ActionUpdate, req.UpdatedBy, details)

	req.ID, _ = strconv.ParseUint(id, 10, 64)
	err := o.elasticSvc.IndexDocument(context.Background(), esIndexName, id, req)
	if err != nil {
		log.Printf("Failed to update thirdpartyappevent %s in Elasticsearch: %v", id, err)
	}
}

// OnDelete is fired after a ThirdPartyAppEvent record is successfully deleted.
func (o *ThirdPartyAppEventObserver) OnDelete(id string, deletedBy string) {
	o.auditService.Record(entityName, id, audit.ActionDelete, deletedBy, "soft-delete committed to ledger")

	deletedAt := time.Now().UTC().Format(time.RFC3339)
	err := o.elasticSvc.UpdateDocument(context.Background(), esIndexName, id, map[string]interface{}{
		"deleted_at": deletedAt,
	})
	if err != nil {
		log.Printf("Failed to mark thirdpartyappevent %s as deleted in Elasticsearch: %v", id, err)
	}
}

// OnList is fired after a successful listing of all ThirdPartyAppEvent records.
func (o *ThirdPartyAppEventObserver) OnList(count int, actor string) {
	details := fmt.Sprintf("returned %d record(s)", count)
	o.auditService.Record(entityName, "", audit.ActionList, actor, details)
}
