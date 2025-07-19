package elasticlog

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

// LogToElastic logs a semantic event to Elasticsearch.
// If username, password, or index are empty, falls back to environment variables or "contoso-logs" for index.
// If index ends with "-", appends the current date (YYYY-MM-DD) for rolling indices.
func LogToElastic(level, msg string, fields map[string]interface{}, index, username, password string) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if username == "" {
		username = os.Getenv("ELASTICSEARCH_USERNAME")
	}
	if password == "" {
		password = os.Getenv("ELASTICSEARCH_PASSWORD")
	}
	if index == "" {
		index = "contoso-"
	}
	// Rolling index: if index ends with "-", append date
	if strings.HasSuffix(index, "-") {
		index = index + time.Now().Format("2006-01-02")
	}
	if esURL == "" || username == "" || password == "" {
		log.Println("Elasticsearch credentials not set, skipping log")
		return
	}
	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
		Username:  username,
		Password:  password,
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Printf("Failed to create Elasticsearch client: %v", err)
		return
	}
	doc := map[string]interface{}{
		"@timestamp": time.Now().Format(time.RFC3339),
		"level":      level,
		"service":    "contoso-backend",
		"message":    msg,
	}
	for k, v := range fields {
		doc[k] = v
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		log.Printf("Failed to encode log document: %v", err)
		return
	}
	_, err = es.Index(
		index,
		&buf,
		es.Index.WithContext(context.Background()),
	)
	if err != nil {
		log.Printf("Failed to log to Elasticsearch: %v", err)
	}
}
