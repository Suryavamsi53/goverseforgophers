# Prometheus & Time-Series Data

## 1. Learning Objectives
* **What you'll learn**: The architecture of Prometheus, a pull-based time-series database (TSDB), and how to expose Go application metrics via HTTP.
* **Why it matters**: You cannot fix what you cannot measure. Without Prometheus, you have no idea if your Go API is returning 500s or if memory is leaking until a customer complains.
* **Where it's used**: Almost every modern cloud infrastructure uses Prometheus as the standard for metric collection.

---

## 2. Real-world Story
Imagine managing a large hotel. 
A **Push-based** system is like having 500 guests call the front desk every 10 seconds to say "My room is fine!" The front desk gets overwhelmed and crashes.
Prometheus is a **Pull-based** system. The front desk (Prometheus) has a clipboard. Every 15 seconds, the manager walks down the hall, knocks on every door, and asks "Is everything okay?" If a room is empty or on fire, the manager knows instantly, and the guests never have to do any work.

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Prometheus Server] -->|HTTP GET /metrics| B[Go Service 1]
    A -->|HTTP GET /metrics| C[Go Service 2]
    A -->|HTTP GET /metrics| D[Node Exporter (Linux CPU)]
    
    A -->|Stores Data in| E[(Time-Series DB)]
    E -->|Queried via PromQL by| F[Grafana]
    
    A -.->|Fires Rules to| G[Alertmanager]
    G -.->|Slack / PagerDuty| H[On-Call Engineer]
```

---

## 4. Internal Working (Under the Hood)
Prometheus is specifically designed for numerical **Time-Series Data**.
A time-series is a stream of timestamped values belonging to the exact same metric name and labels.
Format: `<metric_name>{<label_name>=<label_value>, ...} <timestamp> <value>`
Example: `http_requests_total{method="GET", status="200"} 1689345600 42`
Prometheus scrapes this text data from your Go app, highly compresses it in RAM, and flushes it to disk in sequential blocks.

---

## 5. Compiler Behavior
* **Lock-Free Counters**: The official Go client `github.com/prometheus/client_golang` is insanely optimized. Updating a metric (e.g., `counter.Inc()`) does NOT lock your Go application with a heavy Mutex. It uses Go's atomic hardware instructions (`atomic.AddUint64`), meaning you can increment metrics millions of times per second with zero CPU penalty.

---

## 6. Memory Management
* **Cardinality Explosion**: The #1 reason Prometheus crashes. If you add a label `user_id="suryavamsi53"` to an HTTP metric, Prometheus creates a brand new time-series in RAM for *every single user*. If you have 1 million users, you create 1 million time-series. The Go memory will instantly OOM. Never use unbounded/high-cardinality data in metric labels!

---

## 7. Code Examples

### 🔹 Example 1: Simple (The HTTP Server)
```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

func main() {
    // Expose the /metrics endpoint for Prometheus to scrape!
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8080", nil)
}
```

### 🔹 Example 2: Intermediate (Counters and Gauges)
```go
import "github.com/prometheus/client_golang/prometheus"

// Counter: Only goes UP (e.g., Total HTTP Requests)
var httpRequestsTotal = prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests.",
    },
)

// Gauge: Can go UP and DOWN (e.g., Current Memory Usage, Active WebSockets)
var activeWebSockets = prometheus.NewGauge(
    prometheus.GaugeOpts{
        Name: "active_websockets_current",
        Help: "Current number of active WebSocket connections.",
    },
)

func init() {
    // You MUST register the metrics with the global registry!
    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(activeWebSockets)
}
```

### 🔹 Example 3: Advanced (Histograms for Latency)
```go
// Histograms track the DISTRIBUTION of data (like API Latency)
var requestLatency = prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
        Name:    "http_request_duration_seconds",
        Help:    "Histogram of request latencies.",
        Buckets: []float64{0.1, 0.5, 1, 2, 5}, // Important! Defines the buckets.
    },
    []string{"method", "route"}, // The Labels
)

func APIMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start).Seconds()
        
        // Record the latency!
        requestLatency.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}
```

### 🔹 Example 4: Production (PromQL Query)
```promql
# The magic of PromQL.
# Calculate the 99th percentile Latency over the last 5 minutes!
histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
```

### 🔹 Example 5: Interview
```go
// Q: What happens to the metrics if the Go server crashes?
// A: The metrics are lost! Because Prometheus is Pull-based (scrapes every 15s), 
// any metrics updated in the 14 seconds before the crash will never be scraped. 
// This is an acceptable trade-off for the massive performance gain.
```

---

## 8. Production Examples
1. **Node Exporter**: You run a tiny Go binary called `node_exporter` on every Linux VM. It exposes metrics about CPU, RAM, and Disk space for Prometheus to scrape.
2. **Alertmanager**: Prometheus doesn't send emails. It calculates rules (`if RAM > 90%`) and forwards an alert payload to `Alertmanager`, which handles deduplication and routing to Slack or PagerDuty.

---

## 9. Performance & Benchmarking
* **Rate over Absolute Values**: Never look at absolute counter values. A counter resetting to 0 when a Go pod restarts is completely normal. Always use the `rate()` function in PromQL, which calculates the per-second average and mathematically handles the resets transparently!

---

## 10. Best Practices
* ✅ **Do**: Suffix your metric names with the unit (e.g., `_total`, `_seconds`, `_bytes`) so engineers reading Grafana instantly know what the number means.
* ❌ **Don't**: Build an HTTP Client in Go to manually push metrics to a database. You will DDoS your own database. Just expose `/metrics` and let Prometheus do the work.
* 🏢 **Google / Uber / Netflix Style**: Use `ServiceMonitors` (Prometheus Operator) in Kubernetes. You don't configure Prometheus manually. You deploy a Go Pod, attach a `ServiceMonitor` YAML, and Prometheus automatically discovers and scrapes it!

---

## 11. Common Mistakes
1. **Too Many Histogram Buckets**: If you define a Histogram with 100 buckets, and it has 3 labels with 10 values each, you just created `100 * 10 * 10 * 10 = 100,000` time-series for a SINGLE metric. This will crash Prometheus.
2. **Scraping Too Often**: Setting a scrape interval of `1s`. Most dashboards don't need 1-second granularity. 15s or 30s is the industry standard and saves 95% of storage costs.

---

## 12. Debugging
How to troubleshoot Prometheus in production:
* **The `/targets` UI**: Open the Prometheus Web UI and go to Status -> Targets. If it says `DOWN`, it means Prometheus physically cannot reach your Go app (usually a Kubernetes networking/firewall issue).

---

## 13. Exercises
1. **Easy**: Write a Go HTTP server that imports `promhttp` and serves the default metrics. Hit `/metrics` in your browser.
2. **Medium**: Create a Custom Counter that increments every time someone visits the `/hello` route.
3. **Hard**: Spin up Prometheus locally using Docker Compose, point it at your Go app's IP, and view your custom metric in the Prometheus UI.
4. **Expert**: Implement an `http_request_duration_seconds` Histogram using middleware and write a PromQL query to calculate the `rate()` of requests per second.

---

## 14. Quiz
1. **MCQ**: What data structure should you use to track the memory usage of a Go process?
   * (A) Counter (B) Gauge (C) Histogram. *(Answer: B, because memory goes up AND down).*
2. **Code Review**: `metrics.WithLabelValues(r.URL.Path).Inc()`. Why is using the raw URL path highly dangerous? *(If the path contains an ID, like `/users/123` and `/users/456`, you will create infinite cardinality. You must normalize the route to `/users/{id}` BEFORE passing it to the label).*

---

## 15. FAANG Interview Questions
* **Beginner**: Explain Pull vs Push observability models.
* **Intermediate**: What is a Histogram and why is it preferred over calculating averages (Mean) for API latency?
* **Senior (Google/Meta)**: Explain Prometheus federation and Long-Term Storage (Thanos/Cortex). How do you query 3 years of metrics across 50 global Kubernetes clusters?

---

## 16. Mini Project
**The Go Load Tester**
* Write a Go API with a `/metrics` endpoint and a `/work` endpoint that sleeps for a random time between 10ms and 500ms.
* Run Prometheus in Docker to scrape it.
* Write a script that hits `/work` 1,000 times concurrently.
* View the resulting Histogram distribution in the Prometheus UI to visually prove the random latency curve.

---

## 17. Enterprise Features & Observability
* **Exemplars**: Modern Prometheus supports Exemplars. A single metric point (like a slow API call) can carry a `trace_id`. In Grafana, clicking on the spike in the graph instantly jumps you to the exact Distributed Trace that caused the spike!

---

## 18. Source Code Reading
Walkthrough of `github.com/prometheus/client_golang`.
* **The Registry**: Study the `prometheus.Registry`. It enforces strict thread-safety and validation, ensuring you don't accidentally register two metrics with the same name, which would corrupt the text output.

---

## 19. Architecture
* **The RED Method**: The 3 metrics you MUST monitor for every Go microservice: Rate (Requests/sec), Errors (500s/sec), and Duration (Latency Histograms).

---

## 20. Summary & Cheat Sheet
* **Model**: Pull-based TSDB.
* **Format**: Pure text, exposed on `/metrics`.
* **Types**: Counter (Up), Gauge (Up/Down), Histogram (Distribution).
* **Cardinality**: Keep label values limited (No UUIDs!).
