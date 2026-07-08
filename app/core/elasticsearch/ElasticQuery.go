package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

// ElasticService provides generic operations for Elasticsearch
type ElasticService struct {
	client *elasticsearch.Client
}

// NewElasticService creates a new ElasticService
func NewElasticService() *ElasticService {
	return &ElasticService{
		client: GetClient(),
	}
}

// ---------------------------------------------------------------------------
// Search filter
// ---------------------------------------------------------------------------

// MetadataFilter defines advanced search criteria with support for full-text
// queries, term filters, date ranges, soft-delete exclusion, and pagination.
type MetadataFilter struct {
	// Full-text search (multi_match across name, value, tags)
	Query string `json:"query,omitempty"`

	// Term filters
	PatientID string   `json:"patient_id,omitempty"`
	AssetID   string   `json:"asset_id,omitempty"`
	Type      string   `json:"type,omitempty"`
	Status    string   `json:"status,omitempty"`
	Tags      []string `json:"tags,omitempty"`

	// Date-range filters
	CreatedFrom *time.Time `json:"created_from,omitempty"`
	CreatedTo   *time.Time `json:"created_to,omitempty"`
	UpdatedFrom *time.Time `json:"updated_from,omitempty"`
	UpdatedTo   *time.Time `json:"updated_to,omitempty"`

	// When true, records with a non-null deleted_at are excluded from results.
	ExcludeDeleted bool `json:"exclude_deleted,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// ---------------------------------------------------------------------------
// CRUD operations
// ---------------------------------------------------------------------------

// IndexDocument indexes (creates or fully replaces) a document in the specified index.
func (s *ElasticService) IndexDocument(ctx context.Context, indexName string, docID string, document interface{}) error {
	body, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("error marshaling document: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error indexing document: %s", res.String())
	}

	return nil
}

// UpdateDocument performs a partial update on a document in the specified index.
func (s *ElasticService) UpdateDocument(ctx context.Context, indexName string, docID string, doc map[string]interface{}) error {
	updateBody := map[string]interface{}{
		"doc": doc,
	}
	body, err := json.Marshal(updateBody)
	if err != nil {
		return fmt.Errorf("error marshaling update body: %w", err)
	}

	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: docID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elasticsearch error updating document: %s", res.String())
	}

	return nil
}

// DeleteDocument removes a document from the specified index by its ID.
func (s *ElasticService) DeleteDocument(ctx context.Context, indexName string, docID string) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: docID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return fmt.Errorf("document %q not found in index %q", docID, indexName)
	}

	if res.IsError() {
		return fmt.Errorf("elasticsearch error deleting document (status %d): %s", res.StatusCode, res.String())
	}

	return nil
}

// GetDocumentByID retrieves a single document from the specified index by its ID.
func (s *ElasticService) GetDocumentByID(ctx context.Context, indexName string, docID string) (map[string]interface{}, error) {
	req := esapi.GetRequest{
		Index:      indexName,
		DocumentID: docID,
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return nil, fmt.Errorf("error getting document: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, fmt.Errorf("document %q not found in index %q", docID, indexName)
	}

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error getting document (status %d): %s", res.StatusCode, res.String())
	}

	var result struct {
		Source map[string]interface{} `json:"_source"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding document: %w", err)
	}

	return result.Source, nil
}

// ---------------------------------------------------------------------------
// Search
// ---------------------------------------------------------------------------

// Search performs a filtered, paginated search on the specified index using
// a MetadataFilter. It builds a bool query with full-text, term, and date-range
// clauses, and returns the matching documents together with the total hit count.
func (s *ElasticService) Search(ctx context.Context, indexName string, filter *MetadataFilter) ([]map[string]interface{}, int64, error) {
	if filter == nil {
		filter = &MetadataFilter{}
	}

	// Pagination defaults / caps
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	// ---- bool query components ----
	must := []map[string]interface{}{}
	filters := []map[string]interface{}{}
	mustNot := []map[string]interface{}{}

	// Full-text search across name, value, tags
	if filter.Query != "" {
		must = append(must, map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  filter.Query,
				"fields": []string{"name^2", "value", "tags"},
			},
		})
	}

	// Term filter: patient_id
	if filter.PatientID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"patient_id": filter.PatientID,
			},
		})
	}

	// Term filter: asset_id
	if filter.AssetID != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"asset_id": filter.AssetID,
			},
		})
	}

	// Term filter: type
	if filter.Type != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"type": filter.Type,
			},
		})
	}

	// Term filter: status
	if filter.Status != "" {
		filters = append(filters, map[string]interface{}{
			"term": map[string]interface{}{
				"status": filter.Status,
			},
		})
	}

	// Terms filter: tags (match any)
	if len(filter.Tags) > 0 {
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{
				"tags": filter.Tags,
			},
		})
	}

	// Date-range filter: created_at
	createdRange := map[string]interface{}{}
	if filter.CreatedFrom != nil {
		createdRange["gte"] = filter.CreatedFrom.Format(time.RFC3339)
	}
	if filter.CreatedTo != nil {
		createdRange["lte"] = filter.CreatedTo.Format(time.RFC3339)
	}
	if len(createdRange) > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"created_at": createdRange,
			},
		})
	}

	// Date-range filter: updated_at
	updatedRange := map[string]interface{}{}
	if filter.UpdatedFrom != nil {
		updatedRange["gte"] = filter.UpdatedFrom.Format(time.RFC3339)
	}
	if filter.UpdatedTo != nil {
		updatedRange["lte"] = filter.UpdatedTo.Format(time.RFC3339)
	}
	if len(updatedRange) > 0 {
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{
				"updated_at": updatedRange,
			},
		})
	}

	// Exclude soft-deleted records
	if filter.ExcludeDeleted {
		mustNot = append(mustNot, map[string]interface{}{
			"exists": map[string]interface{}{
				"field": "deleted_at",
			},
		})
	}

	// ---- Assemble the bool query ----
	boolQuery := map[string]interface{}{}
	if len(must) > 0 {
		boolQuery["must"] = must
	}
	if len(filters) > 0 {
		boolQuery["filter"] = filters
	}
	if len(mustNot) > 0 {
		boolQuery["must_not"] = mustNot
	}

	var searchQuery map[string]interface{}
	if len(boolQuery) == 0 {
		searchQuery = map[string]interface{}{
			"match_all": map[string]interface{}{},
		}
	} else {
		searchQuery = map[string]interface{}{
			"bool": boolQuery,
		}
	}

	searchBody := map[string]interface{}{
		"query": searchQuery,
		"size":  filter.Limit,
		"from":  filter.Offset,
	}

	body, err := json.Marshal(searchBody)
	if err != nil {
		return nil, 0, fmt.Errorf("error encoding search body: %w", err)
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return nil, 0, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var errResp map[string]interface{}
		json.NewDecoder(res.Body).Decode(&errResp)
		return nil, 0, fmt.Errorf("elasticsearch search error (status %d): %v", res.StatusCode, errResp)
	}

	// Parse typed response to extract total + hits
	var searchResult struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID     string                 `json:"_id"`
				Source map[string]interface{} `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(res.Body).Decode(&searchResult); err != nil {
		return nil, 0, fmt.Errorf("error decoding search response: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(searchResult.Hits.Hits))
	for _, hit := range searchResult.Hits.Hits {
		doc := hit.Source
		doc["_id"] = hit.ID
		results = append(results, doc)
	}

	return results, searchResult.Hits.Total.Value, nil
}

// ---------------------------------------------------------------------------
// Cluster / index utilities
// ---------------------------------------------------------------------------

// GetHealth returns the Elasticsearch cluster health as a generic map.
func (s *ElasticService) GetHealth(ctx context.Context) (map[string]interface{}, error) {
	req := esapi.ClusterHealthRequest{}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return nil, fmt.Errorf("error getting cluster health: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch health check error (status %d): %s", res.StatusCode, res.String())
	}

	var health map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return nil, fmt.Errorf("error decoding health response: %w", err)
	}

	return health, nil
}

// CountDocuments returns the total number of documents in the specified index.
func (s *ElasticService) CountDocuments(ctx context.Context, indexName string) (int64, error) {
	req := esapi.CountRequest{
		Index: []string{indexName},
	}

	res, err := req.Do(ctx, s.client)
	if err != nil {
		return 0, fmt.Errorf("error counting documents: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return 0, fmt.Errorf("elasticsearch count error (status %d): %s", res.StatusCode, res.String())
	}

	var countResult struct {
		Count int64 `json:"count"`
	}

	if err := json.NewDecoder(res.Body).Decode(&countResult); err != nil {
		return 0, fmt.Errorf("error decoding count response: %w", err)
	}

	return countResult.Count, nil
}