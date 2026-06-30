package elasticsearch

import (
	"log"
	"os"

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
	return Client
}

// GetClient returns the initialized Elasticsearch client
func GetClient() *elasticsearch.Client {
	if Client == nil {
		return InitElasticClient()
	}
	return Client
}
