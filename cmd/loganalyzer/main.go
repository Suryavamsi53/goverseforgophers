package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// LogEntry holds parsed data from a single log line
type LogEntry struct {
	IP       string
	Method   string
	Endpoint string
	Status   string
	Duration time.Duration
}

// LogStats aggregates metrics
type LogStats struct {
	mu             sync.Mutex
	TotalReqs      int
	StatusCounts   map[string]int
	IPCounts       map[string]int
	PathCounts     map[string]int
	TotalLatency   time.Duration
	MaxLatency     time.Duration
	Durations      []time.Duration
	SuspiciousIPs  map[string]int
	EndpointErrors map[string]int
}

func NewLogStats() *LogStats {
	return &LogStats{
		StatusCounts:   make(map[string]int),
		IPCounts:       make(map[string]int),
		PathCounts:     make(map[string]int),
		Durations:      make([]time.Duration, 0, 10000),
		SuspiciousIPs:  make(map[string]int),
		EndpointErrors: make(map[string]int),
	}
}

// Add safely updates statistics with a new LogEntry
func (s *LogStats) Add(entry LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.TotalReqs++
	s.StatusCounts[entry.Status]++
	
	// Clean up IP (remove port if present)
	cleanIP := entry.IP
	if idx := strings.LastIndex(cleanIP, ":"); idx != -1 && !strings.HasSuffix(cleanIP, "]") {
		// handle ipv4 with port
		cleanIP = cleanIP[:idx]
	} else if strings.HasPrefix(cleanIP, "[") && strings.Contains(cleanIP, "]:") {
		// handle ipv6 with port
		cleanIP = cleanIP[:strings.LastIndex(cleanIP, "]:")+1]
	}
	s.IPCounts[cleanIP]++

	// Clean up URL to just path
	u, err := url.Parse(entry.Endpoint)
	path := entry.Endpoint
	if err == nil && u.Path != "" {
		path = u.Path
	}
	s.PathCounts[path]++

	// Track errors per endpoint (5xx or 4xx)
	if strings.HasPrefix(entry.Status, "4") || strings.HasPrefix(entry.Status, "5") {
		s.EndpointErrors[path]++
	}

	// Track suspicious IPs (401 Unauthorized or 429 Rate Limited)
	if entry.Status == "401" || entry.Status == "429" {
		s.SuspiciousIPs[cleanIP]++
	}

	// Latency metrics
	s.TotalLatency += entry.Duration
	if entry.Duration > s.MaxLatency {
		s.MaxLatency = entry.Duration
	}
	if entry.Duration > 0 {
		s.Durations = append(s.Durations, entry.Duration)
	}
}

func main() {
	formatFlag := flag.String("format", "goverse", "Log format to parse: 'goverse' or 'nginx'")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: loganalyzer [-format=goverse|nginx] <logfile>")
		os.Exit(1)
	}
	filePath := flag.Args()[0]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	var logRegex *regexp.Regexp
	
	if *formatFlag == "goverse" {
		// Example: 2026/07/04 10:41:31 [id] "GET http://localhost:8080/ HTTP/1.1" from [::1]:35196 - 200 15624B in 3.84544ms
		logRegex = regexp.MustCompile(`"([A-Z]+)\s(.*?)\sHTTP/.*?"\sfrom\s(.*?)\s-\s(\d{3})\s.*?\sin\s(.*)`)
	} else {
		// Example: 192.168.1.1 - - [date] "GET /api/users HTTP/1.1" 200 1024 "-" "Mozilla/5.0"
		logRegex = regexp.MustCompile(`^(\S+).+?"([A-Z]+)\s(.*?)\sHTTP/.*?"\s(\d{3})`)
	}

	stats := NewLogStats()
	linesCh := make(chan string, 1000)
	var wg sync.WaitGroup
	numWorkers := 4

	start := time.Now()

	// Start worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range linesCh {
				matches := logRegex.FindStringSubmatch(line)
				if *formatFlag == "goverse" && len(matches) == 6 {
					duration, _ := time.ParseDuration(matches[5])
					stats.Add(LogEntry{
						Method:   matches[1],
						Endpoint: matches[2],
						IP:       matches[3],
						Status:   matches[4],
						Duration: duration,
					})
				} else if *formatFlag == "nginx" && len(matches) == 5 {
					stats.Add(LogEntry{
						IP:       matches[1],
						Method:   matches[2],
						Endpoint: matches[3],
						Status:   matches[4],
						Duration: 0,
					})
				}
			}
		}()
	}

	// Read file and send to workers
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linesCh <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	close(linesCh)
	wg.Wait()

	elapsed := time.Since(start)

	// Print Report
	fmt.Println("=====================================")
	fmt.Printf("   LOG ANALYSIS REPORT (%s)\n", strings.ToUpper(*formatFlag))
	fmt.Println("=====================================")
	fmt.Printf("Total Requests : %d\n", stats.TotalReqs)
	if stats.TotalReqs > 0 && *formatFlag == "goverse" {
		avgLatency := stats.TotalLatency / time.Duration(stats.TotalReqs)
		fmt.Printf("Avg Latency    : %v\n", avgLatency)
		fmt.Printf("Max Latency    : %v\n", stats.MaxLatency)
		
		// Calculate Percentiles
		if len(stats.Durations) > 0 {
			sort.Slice(stats.Durations, func(i, j int) bool {
				return stats.Durations[i] < stats.Durations[j]
			})
			p95Index := int(float64(len(stats.Durations)) * 0.95)
			p99Index := int(float64(len(stats.Durations)) * 0.99)
			fmt.Printf("P95 Latency    : %v\n", stats.Durations[p95Index])
			fmt.Printf("P99 Latency    : %v\n", stats.Durations[p99Index])
		}
	}
	fmt.Printf("Time Taken     : %v\n", elapsed)
	fmt.Println("-------------------------------------")
	
	fmt.Println("HTTP Status Codes:")
	printTop(stats.StatusCounts, 5)
	
	fmt.Println("\nTop 5 IP Addresses:")
	printTop(stats.IPCounts, 5)

	if len(stats.SuspiciousIPs) > 0 {
		fmt.Println("\nTop 5 Suspicious IPs (401/429):")
		printTop(stats.SuspiciousIPs, 5)
	}
	
	fmt.Println("\nTop 5 Endpoints:")
	printTop(stats.PathCounts, 5)

	if len(stats.EndpointErrors) > 0 {
		fmt.Println("\nTop 5 Endpoints by Errors (4xx/5xx):")
		printTop(stats.EndpointErrors, 5)
	}
	fmt.Println("=====================================")
}

func printTop(counts map[string]int, limit int) {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range counts {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	for i, kv := range ss {
		if i >= limit {
			break
		}
		padLen := 20 - len(kv.Key)
		padding := " "
		if padLen > 0 {
			padding = strings.Repeat(" ", padLen)
		}
		fmt.Printf("  %s%s : %d\n", kv.Key, padding, kv.Value)
	}
}
