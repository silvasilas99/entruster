package elasticsearch

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

// Client is the global Elasticsearch client
var Client *elasticsearch.Client

// InitElasticClient establishes and exports a connection to the local Elasticsearch containerized server
func InitElasticClient() *elasticsearch.Client {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://localhost:9200" // Default for local testing
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			esURL,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	log.Printf("Successfully connected to Elasticsearch at %s", esURL)
	Client = es

	// Auto-create required indices on startup
	if err := InitializeIndex(es, "metadata"); err != nil {
		log.Printf("⚠️  Warning: failed to initialize 'metadata' index: %v", err)
	}

	return Client
}

// InitializeIndex checks whether the given index exists in Elasticsearch.
// If it does not exist, it creates the index with the mapping that matches
// the MetadataModel struct fields.
func InitializeIndex(es *elasticsearch.Client, indexName string) error {
	// 1. Check if the index already exists
	res, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return fmt.Errorf("error checking if index %q exists: %w", indexName, err)
	}
	defer res.Body.Close()

	// Status 200 means the index exists — nothing to do
	if res.StatusCode == 200 {
		log.Printf("✅ Elasticsearch index %q already exists", indexName)
		return nil
	}

	// 2. Create the index with explicit mappings
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0
		},
		"mappings": {
			"properties": {
				"id":              {"type": "long"},
				"patient_id":      {"type": "long"},
				"asset_id":        {"type": "long"},
				"name":            {"type": "text", "fields": {"keyword": {"type": "keyword"}}},
				"value":           {"type": "text"},
				"version":         {"type": "keyword"},
				"owner":           {"type": "keyword"},
				"rights":          {"type": "text"},
				"terms_of_access": {"type": "text"},
				"created_at":      {"type": "date"},
				"updated_at":      {"type": "date"},
				"deleted_at":      {"type": "date"},
				"created_by":      {"type": "keyword"},
				"updated_by":      {"type": "keyword"},
				"deleted_by":      {"type": "keyword"}
			}
		}
	}`

	createRes, err := es.Indices.Create(
		indexName,
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("error creating index %q: %w", indexName, err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		return fmt.Errorf("elasticsearch error creating index %q: %s", indexName, createRes.String())
	}

	log.Printf("✅ Elasticsearch index %q created successfully", indexName)
	return nil
}

// GetClient returns the initialized Elasticsearch client
func GetClient() *elasticsearch.Client {
	if Client == nil {
		return InitElasticClient()
	}
	return Client
}
