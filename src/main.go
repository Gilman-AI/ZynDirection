// Copyright (c) 2024 Herbert F. Gilman _a.k.a._ HFG Ventures LLC.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/storage"
)

var (
	mappings      map[string]string
	mappingsMutex sync.RWMutex
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	err = loadMappings(ctx, client)
	if err != nil {
		log.Fatalf("Failed to load mappings: %v", err)
	}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for {
			<-ticker.C
			if err := loadMappings(ctx, client); err != nil {
				log.Printf("Failed to refresh mappings: %v", err)
			} else {
				log.Println("Mappings successfully refreshed")
			}
		}
	}()
	http.HandleFunc("/", redirectHandler)
	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func loadMappings(ctx context.Context, client *storage.Client) error {
	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		return fmt.Errorf("BUCKET_NAME environment variable not set")
	}
	objectName := os.Getenv("MAPPINGS_OBJECT_NAME")
	if objectName == "" {
		objectName = "mappings.json"
	}
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(objectName)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to create object reader: %v", err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read object data: %v", err)
	}
	var newMappings map[string]string
	if err := json.Unmarshal(data, &newMappings); err != nil {
		return fmt.Errorf("failed to parse JSON: %v", err)
	}
	mappingsMutex.Lock()
	mappings = newMappings
	mappingsMutex.Unlock()
	return nil
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	mappingsMutex.RLock()
	originalURL, ok := mappings[path]
	mappingsMutex.RUnlock()
	if ok {
		http.Redirect(w, r, originalURL, http.StatusFound)
	} else {
		http.NotFound(w, r)
	}
}
