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
			{
				ID:          "proj-4",
				Slug:        "log-analyzer",
				Title:       "Concurrent Log Analyzer",
				Description: "Build a blazing fast, concurrent log analysis pipeline using worker pools and regex. Analyze production logs and extract critical P99 latency metrics.",
				Scenario: `Your company's API Gateway generates massive amount of logs. 
During a production incident, you need to quickly identify:
1. P95 and P99 Latency spikes.
2. Malicious IP addresses hitting rate limits (429) or auth failures (401).

You must build a high-performance log analyzer in Go that parses 
these log lines concurrently using a Worker Pool.

**Input Format:**
` + "`" + `2026/07/04 10:41:31 [req-id] "GET /api/users HTTP/1.1" from 192.168.1.10:35196 - 200 15624B in 3.845ms` + "`" + `

**Output Requirements:**
- Aggregate Total Requests, Max Latency
- Calculate P95 and P99 Latency
- Track counts of Suspicious IPs (status 401 or 429)`,
				Difficulty:  domain.DifficultyIntermediate,
				Tags:        []string{"Concurrency", "Worker Pool", "Regex"},
				Icon:        "📊",
				Color:       "orange",
				Requirements: []string{
					"Parse Nginx or GoVerse formatted log lines using Regular Expressions.",
					"Implement a concurrent worker pool pattern to process thousands of lines per second.",
					"Safely aggregate metrics (IP counts, endpoint latencies) using sync.Mutex.",
					"Calculate P95 and P99 latency percentiles.",
					"Identify malicious IPs hitting 401 or 429 status codes.",
				},
				Tips: []string{
					"Use `regexp` to extract Method, Endpoint, IP, Status, and Duration.",
					"Use a `chan string` to send log lines to a pool of worker goroutines.",
					"Use `sync.Mutex` inside LogStats to prevent race conditions during aggregation.",
					"Store durations in a slice, sort them, and pick the index at len*0.95 and len*0.99.",
					"Strip the port number (e.g., `:35196`) from the IP address for aggregation.",
				},
				StarterCode: `package main

import (
	"fmt"
	"regexp"
	"sync"
	"time"
)

type LogEntry struct {
	IP       string
	Endpoint string
	Status   string
	Duration time.Duration
}

type LogStats struct {
	mu            sync.Mutex
	TotalReqs     int
	SuspiciousIPs map[string]int
	Durations     []time.Duration
	// TODO: Add other fields as needed (e.g., MaxLatency)
}

func NewLogStats() *LogStats {
	return &LogStats{
		SuspiciousIPs: make(map[string]int),
		Durations:     make([]time.Duration, 0),
	}
}

// ParseLine extracts data from a single log string using Regex
func ParseLine(line string) (LogEntry, error) {
	// Task 1: Parse Nginx or GoVerse formatted log lines
	// HINT: Use regexp.MustCompile
	return LogEntry{}, nil
}

// ProcessLogs orchestrates the worker pool pattern
func ProcessLogs(lines []string, numWorkers int) *LogStats {
	stats := NewLogStats()
	
	// Task 2: Implement concurrent worker pool
	// Task 3: Safely aggregate metrics (IP counts, endpoint latencies)
	// Task 5: Identify malicious IPs hitting 401 or 429 status codes
	
	return stats
}

// P95 returns the 95th percentile latency
func (s *LogStats) P95() time.Duration {
	// Task 4: Calculate P95 latency percentile
	return 0
}

// P99 returns the 99th percentile latency
func (s *LogStats) P99() time.Duration {
	// Task 4: Calculate P99 latency percentile
	return 0
}

func main() {
	fmt.Println("Log Analyzer Starter Code")
}
`,
				TestFile: `package main

import (
	"fmt"
	"testing"
	"time"
)

// Checkpoint 1: Parse Nginx or GoVerse formatted log lines (25 Marks)
func TestParseLine(t *testing.T) {
	line := ` + "`" + `2026/07/04 10:41:31 [id] "GET /api/users HTTP/1.1" from 192.168.1.10:35196 - 401 15B in 4.5ms` + "`" + `
	entry, err := ParseLine(line)
	if err != nil {
		t.Fatalf("Checkpoint 1 Failed: unexpected error: %v", err)
	}
	if entry.Status != "401" {
		t.Errorf("Checkpoint 1 Failed: expected status 401, got %v", entry.Status)
	}
	if entry.Endpoint != "/api/users" {
		t.Errorf("Checkpoint 1 Failed: expected endpoint /api/users, got %v", entry.Endpoint)
	}
	
	if entry.Duration != 4500*time.Microsecond {
		t.Errorf("Checkpoint 1 Failed: expected 4.5ms duration, got %v", entry.Duration)
	}
}

// Checkpoint 2: Safely aggregate metrics using sync.Mutex (25 Marks)
func TestLogStatsAggregation(t *testing.T) {
	lines := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		lines[i] = ` + "`" + `2026/07/04 10:41:31 [id] "GET / HTTP/1.1" from 10.0.0.1:1111 - 200 15B in 1ms` + "`" + `
	}
	stats := ProcessLogs(lines, 10)
	if stats == nil {
		t.Fatalf("Checkpoint 2 Failed: ProcessLogs returned nil")
	}
	if stats.TotalReqs != 1000 {
		t.Errorf("Checkpoint 2 Failed: expected 1000 TotalReqs, got %d. Did you use sync.Mutex and a Worker Pool?", stats.TotalReqs)
	}
}

// Checkpoint 3: Identify malicious IPs hitting 401 or 429 (25 Marks)
func TestSuspiciousIPs(t *testing.T) {
	lines := []string{
		` + "`" + `2026/07/04 [id] "GET / HTTP/1.1" from 10.0.0.1:1111 - 200 10B in 1ms` + "`" + `,
		` + "`" + `2026/07/04 [id] "POST /login HTTP/1.1" from 10.0.0.2:2222 - 401 10B in 2ms` + "`" + `,
		` + "`" + `2026/07/04 [id] "GET /api HTTP/1.1" from 10.0.0.2:3333 - 429 10B in 3ms` + "`" + `,
	}
	stats := ProcessLogs(lines, 2)
	if stats == nil {
		t.Fatalf("Checkpoint 3 Failed: ProcessLogs returned nil")
	}
	
	if count, ok := stats.SuspiciousIPs["10.0.0.2"]; !ok || count != 2 {
		t.Errorf("Checkpoint 3 Failed: expected IP 10.0.0.2 to have 2 suspicious requests, got %v", count)
	}
}

// Checkpoint 4: Calculate P95 and P99 latency percentiles (25 Marks)
func TestPercentiles(t *testing.T) {
	lines := make([]string, 100)
	for i := 0; i < 100; i++ {
		lines[i] = fmt.Sprintf(` + "`" + `2026/07/04 [id] "GET / HTTP/1.1" from 1.1.1.1:80 - 200 10B in %dms` + "`" + `, i+1)
	}
	stats := ProcessLogs(lines, 4)
	if stats == nil {
		t.Fatalf("Checkpoint 4 Failed: ProcessLogs returned nil")
	}
	
	p95 := stats.P95()
	if p95 < 94*time.Millisecond || p95 > 96*time.Millisecond {
		t.Errorf("Checkpoint 4 Failed: expected P95 around 95ms, got %v", p95)
	}
	
	p99 := stats.P99()
	if p99 < 98*time.Millisecond || p99 > 100*time.Millisecond {
		t.Errorf("Checkpoint 4 Failed: expected P99 around 99ms, got %v", p99)
	}
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
