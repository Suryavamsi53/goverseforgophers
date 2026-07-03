package repository

import (
	"context"
	"errors"

	"github.com/suryavamsivaggu/goverse/internal/domain"
)

type MockProjectRepository struct {
	projects []domain.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{
		projects: []domain.Project{
			{
				ID:          "proj-1",
				Slug:        "url-shortener",
				Title:       "URL Shortener API",
				Description: "Design and build a production-ready REST API for shortening URLs. Includes PostgreSQL integration, Redis caching, and rate limiting.",
				Difficulty:  domain.DifficultyBeginner,
				Tags:        []string{"REST API", "PostgreSQL", "Redis"},
				Icon:        "🔗",
				Color:       "emerald",
				Requirements: []string{
					"Create an endpoint to submit a long URL and receive a short code.",
					"Create an endpoint to redirect from a short code to the original URL.",
					"Use PostgreSQL to persist the URL mappings.",
					"Implement a Redis cache layer to speed up redirections.",
					"Add basic rate limiting middleware to prevent abuse.",
				},
				StarterCode: `package main

import (
	"encoding/json"
	"net/http"
)

type ShortenRequest struct {
	URL string ` + "`json:\"url\"`" + `
}

type ShortenResponse struct {
	ShortCode string ` + "`json:\"short_code\"`" + `
}

var urlStore = make(map[string]string)

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse JSON request
	// 2. Generate short code (can be mock or simple hash)
	// 3. Store in map
	// 4. Return JSON response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ShortenResponse{ShortCode: "mock_code"})
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get short code from URL path
	// 2. Look up original URL
	// 3. Redirect (HTTP 302) or return 404
	http.Redirect(w, r, "https://example.com", http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", ShortenHandler)
	http.HandleFunc("/", RedirectHandler)
	http.ListenAndServe(":8080", nil)
}
`,
				TestFile: `package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestURLShortener(t *testing.T) {
	// Test Shorten
	reqBody := ` + "`{\"url\":\"https://example.com\"}`" + `
	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
	rr := httptest.NewRecorder()
	
	handler := http.HandlerFunc(ShortenHandler)
	handler.ServeHTTP(rr, req)
	
	if rr.Code != http.StatusOK && rr.Code != http.StatusCreated {
		t.Errorf("expected 200 or 201, got %v", rr.Code)
	}
	
	var resp struct {
		ShortCode string ` + "`json:\"short_code\"`" + `
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ShortCode == "" {
		t.Error("expected non-empty short_code")
	}
	
	// Test Redirect
	req2, _ := http.NewRequest("GET", "/"+resp.ShortCode, nil)
	rr2 := httptest.NewRecorder()
	
	redirectHandler := http.HandlerFunc(RedirectHandler)
	redirectHandler.ServeHTTP(rr2, req2)
	
	if rr2.Code != http.StatusFound {
		t.Errorf("expected 302 redirect, got %v", rr2.Code)
	}
	
	loc := rr2.Header().Get("Location")
	if loc == "" {
		t.Errorf("expected Location header to be present")
	}
}
`,
			},
			{
				ID:          "proj-2",
				Slug:        "concurrent-scraper",
				Title:       "Concurrent Web Scraper",
				Description: "Build a high-performance web scraper using goroutines, channels, and worker pools. Learn to manage rate limits and graceful shutdown.",
				Difficulty:  domain.DifficultyIntermediate,
				Tags:        []string{"Goroutines", "Channels", "HTTP Client"},
				Icon:        "🕷️",
				Color:       "blue",
				Requirements: []string{
					"Implement a worker pool pattern to fetch multiple URLs concurrently.",
					"Use channels to safely pass URLs to workers and collect results.",
					"Implement rate limiting to avoid overwhelming target servers.",
					"Add graceful shutdown using context.Context.",
					"Extract title and meta tags from the fetched HTML.",
				},
				StarterCode: `package main

import (
	"context"
	"sync"
)

type ScrapeResult struct {
	URL    string
	Length int
}

func ScrapeWorker(ctx context.Context, id int, jobs <-chan string, results chan<- ScrapeResult, wg *sync.WaitGroup) {
	defer wg.Done()
	// TODO: implement worker loop
	// Read from jobs channel, process (mock by setting Length to len(URL)), and send to results
}

func RunScraper(urls []string, numWorkers int) []ScrapeResult {
	// TODO: setup channels, waitgroups, and start workers
	return nil
}

func main() {
	urls := []string{"https://go.dev", "https://google.com", "https://github.com"}
	results := RunScraper(urls, 2)
	for _, r := range results {
		println(r.URL, r.Length)
	}
}
`,
				TestFile: `package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestScrapeWorker(t *testing.T) {
	jobs := make(chan string, 2)
	results := make(chan ScrapeResult, 2)
	var wg sync.WaitGroup
	
	wg.Add(1)
	go ScrapeWorker(context.Background(), 1, jobs, results, &wg)
	
	jobs <- "http://test.com"
	jobs <- "http://example.com"
	close(jobs)
	
	// Wait with timeout
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	
	select {
	case <-c:
	case <-time.After(time.Second):
		t.Fatal("worker did not finish in time (WaitGroup issue)")
	}
	
	close(results)
	var count int
	for range results {
		count++
	}
	
	if count != 2 {
		t.Errorf("expected 2 results, got %d", count)
	}
}

func TestRunScraper(t *testing.T) {
	urls := []string{"a", "b", "c", "d"}
	res := RunScraper(urls, 2)
	if len(res) != len(urls) {
		t.Errorf("expected %d results, got %d", len(urls), len(res))
	}
}
`,
			},
			{
				ID:          "proj-3",
				Slug:        "distributed-kv",
				Title:       "Distributed KV Store",
				Description: "Implement a highly available distributed Key-Value store. Master the Raft consensus algorithm, gRPC communication, and WAL logging.",
				Difficulty:  domain.DifficultyAdvanced,
				Tags:        []string{"gRPC", "Raft Consensus", "Mutexes"},
				Icon:        "🗄️",
				Color:       "purple",
				Requirements: []string{
					"Define a gRPC service for Put, Get, and Delete operations.",
					"Implement a thread-safe in-memory map using sync.RWMutex.",
					"Write an append-only Write-Ahead Log (WAL) to disk for durability.",
					"Implement basic leader election using a simplified Raft concept.",
					"Replicate write operations to follower nodes.",
				},
				StarterCode: `package main

import (
	"sync"
	"errors"
)

type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

func (kv *KVStore) Put(key, value string) error {
	// TODO: Thread-safe write
	return nil
}

func (kv *KVStore) Get(key string) (string, error) {
	// TODO: Thread-safe read
	return "", errors.New("not implemented")
}

func (kv *KVStore) Delete(key string) error {
	// TODO: Thread-safe delete
	return nil
}

func main() {
	store := NewKVStore()
	store.Put("hello", "world")
	val, _ := store.Get("hello")
	println(val)
}
`,
				TestFile: `package main

import (
	"sync"
	"testing"
)

func TestKVStore(t *testing.T) {
	store := NewKVStore()
	
	err := store.Put("key1", "val1")
	if err != nil {
		t.Errorf("unexpected error on Put: %v", err)
	}
	
	val, err := store.Get("key1")
	if err != nil {
		t.Errorf("unexpected error on Get: %v", err)
	}
	if val != "val1" {
		t.Errorf("expected val1, got %v", val)
	}
	
	err = store.Delete("key1")
	if err != nil {
		t.Errorf("unexpected error on Delete: %v", err)
	}
	
	_, err = store.Get("key1")
	if err == nil {
		t.Errorf("expected error on Get for deleted key")
	}
}

func TestKVStoreConcurrency(t *testing.T) {
	store := NewKVStore()
	var wg sync.WaitGroup
	
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			store.Put("key", "val")
			store.Get("key")
		}(i)
	}
	wg.Wait()
}
`,
			},
		},
	}
}

func (r *MockProjectRepository) GetAll(ctx context.Context) ([]domain.Project, error) {
	return r.projects, nil
}

func (r *MockProjectRepository) GetBySlug(ctx context.Context, slug string) (*domain.Project, error) {
	for _, p := range r.projects {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, errors.New("project not found")
}
