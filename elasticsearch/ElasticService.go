package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

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

// IndexDocument indexes a generic document to a specified index
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

// UpdateDocument performs a partial update on a document in a specified index
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

// SearchDocuments performs an advanced search query on a specified index
// It returns a slice of unmarshaled maps representing the generic entities found
func (s *ElasticService) SearchDocuments(ctx context.Context, indexName string, query map[string]interface{}) ([]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(indexName),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch search error: %s", res.String())
	}

	return s.parseSearchResponse(res)
}

// parseSearchResponse parses the generic Search response
func (s *ElasticService) parseSearchResponse(res *esapi.Response) ([]map[string]interface{}, error) {
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &raw); err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %w", err)
	}

	hitsMap, ok := raw["hits"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits format")
	}

	hitsList, ok := hitsMap["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid hits array format")
	}

	var results []map[string]interface{}
	for _, hit := range hitsList {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			continue
		}

		sourceMap, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			continue
		}

		// Optionally include the _id in the source map
		if id, ok := hitMap["_id"].(string); ok {
			sourceMap["_id"] = id
		}

		results = append(results, sourceMap)
	}

	return results, nil
}

// AdvancedFilter constructs a query based on multiple field conditions (term matches)
func (s *ElasticService) AdvancedFilter(ctx context.Context, indexName string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	var mustClauses []map[string]interface{}

	for key, value := range filters {
		mustClauses = append(mustClauses, map[string]interface{}{
			"match": map[string]interface{}{
				key: value,
			},
		})
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": mustClauses,
			},
		},
	}

	return s.SearchDocuments(ctx, indexName, query)
}

// FullTextSearch performs a multi-match full text search across specified fields
func (s *ElasticService) FullTextSearch(ctx context.Context, indexName string, searchText string, fields []string) ([]map[string]interface{}, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  searchText,
				"fields": fields,
			},
		},
	}

	return s.SearchDocuments(ctx, indexName, query)
}

// DateRangeFilter performs a range query on a specified date field.
// fromDate and toDate can be formatted date strings (e.g., RFC3339). Pass empty strings to omit bounds.
func (s *ElasticService) DateRangeFilter(ctx context.Context, indexName string, dateField string, fromDate string, toDate string) ([]map[string]interface{}, error) {
	rangeQuery := map[string]interface{}{}

	if fromDate != "" {
		rangeQuery["gte"] = fromDate
	}
	if toDate != "" {
		rangeQuery["lte"] = toDate
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				dateField: rangeQuery,
			},
		},
	}

	return s.SearchDocuments(ctx, indexName, query)
}