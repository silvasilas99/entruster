package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/silvasilas99/entruster/app/core/chaincode"
	"github.com/silvasilas99/entruster/app/core/elasticsearch"
)

type ActionType string

const (
	ActionCreate ActionType = "CREATE"
	ActionUpdate ActionType = "UPDATE"
	ActionDelete ActionType = "DELETE"
	ActionList   ActionType = "LIST"
)

type AuditService struct {
	elasticSvc *elasticsearch.ElasticService
	cq         *chaincode.ChaincodeQuery
}

func NewAuditService(elasticSvc *elasticsearch.ElasticService, cq *chaincode.ChaincodeQuery) *AuditService {
	// Ensure audit index exists
	if elasticSvc != nil && elasticsearch.Client != nil {
		err := elasticsearch.InitializeIndex(elasticsearch.Client, "audit")
		if err != nil {
			log.Printf("Failed to initialize audit index: %v", err)
		}
	}
	return &AuditService{
		elasticSvc: elasticSvc,
		cq:         cq,
	}
}

// Record indexes the audit log in ElasticSearch. 
// For CREATE/UPDATE/DELETE, the data is inherently present in the ledger's transaction history.
func (s *AuditService) Record(entityName, entityID string, action ActionType, user, details string) {
	record := AuditModel{
		EntityName: entityName,
		EntityId:   entityID,
		Action:     string(action),
		Actor:      user,
		Details:    details,
		OccurredAt: time.Now().UTC(),
	}

	if s.elasticSvc == nil {
		return
	}

	docID := uuid.New().String()
	err := s.elasticSvc.IndexDocument(context.Background(), "audit", docID, record)
	if err != nil {
		log.Printf("Failed to index audit record in Elasticsearch: %v", err)
	}
}

// GetByEntityID retrieves the audit replicas from ElasticSearch using its search engine.
func (s *AuditService) GetByEntityID(entityName, entityID string) []AuditModel {
	if s.elasticSvc == nil {
		return nil
	}

	// Build a custom bool query to search by entityName and entityId in the "audit" index.
	filter := &elasticsearch.MetadataFilter{
		// Since we don't have entityName/entityId on MetadataFilter directly, 
		// we use the query string for exact matches or we can implement a custom query body.
		// For simplicity, we use the Query string field.
		Query: fmt.Sprintf("entity_name:%s AND entity_id:%s", entityName, entityID),
		Limit: 100,
	}
	
	results, _, err := s.elasticSvc.Search(context.Background(), "audit", filter)
	if err != nil {
		log.Printf("Error searching audit logs in ElasticSearch: %v", err)
		return nil
	}

	var audits []AuditModel
	for _, res := range results {
		b, _ := json.Marshal(res)
		var a AuditModel
		json.Unmarshal(b, &a)
		audits = append(audits, a)
	}
	return audits
}

// GetNativeHistory retrieves the native transaction history for an entity directly from the Hyperledger ledger.
func (s *AuditService) GetNativeHistory(entityID string) ([]byte, error) {
	if s.cq == nil {
		return nil, fmt.Errorf("chaincode query is not initialized")
	}
	return s.cq.GetMetadataHistoryById(entityID)
}