# Metrics and Monitoring (Prometheus)

Logs tell you *what* went wrong (e.g., "Database connection refused").
Traces tell you *where* it went wrong (e.g., "The Order Service took 4 seconds").
**Metrics** tell you the overarching health of the system (e.g., "The CPU is at 99%", or "The 99th percentile latency is 500ms").

Metrics are numerical data aggregated over time. 

## 1. Prometheus Architecture (Pull vs Push)

The industry standard for metrics in Kubernetes is **Prometheus**.

Most monitoring systems (like Datadog) use a **Push** model: your Go application installs an SDK and constantly pushes network requests containing metrics to the Datadog servers.

Prometheus uses a **Pull (Scrape)** model. 
1. Your Go application exposes a simple HTTP endpoint: `GET /metrics`.
2. When you hit this endpoint, it returns a plain text list of numbers (e.g., `http_requests_total 425`).
3. The Prometheus Server is a background database that wakes up every 15 seconds, issues an HTTP `GET /metrics` to every single Pod in your Kubernetes cluster, downloads the numbers, and saves them to a Time-Series Database.

**Why Pull?** 
If your Go app is crashing from CPU exhaustion, it doesn't have the compute power to Push metrics to a server. With a Pull model, Prometheus does the heavy lifting, ensuring metrics are gathered even when the application is dying.

## 2. The 4 Golden Signals

Google's SRE (Site Reliability Engineering) handbook dictates that every microservice must export the "4 Golden Signals":

1. **Latency**: How long it takes to service a request (measured in percentiles: P50, P90, P99).
2. **Traffic**: The total demand placed on the system (Requests Per Second).
3. **Errors**: The rate of requests that fail (e.g., HTTP 5xx codes).
4. **Saturation**: How "full" the system is (CPU usage, RAM usage, Database Connection Pool limits).

## 3. Instrumenting Go for Prometheus

```go
import "github.com/prometheus/client_golang/prometheus"

// 1. Define a Metric (A Counter for total requests)
var requestsTotal = prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    },
    []string{"method", "status"}, // Labels for filtering
)

func init() {
    // Register it with the global Prometheus registry
    prometheus.MustRegister(requestsTotal)
}

func HandleCheckout(w http.ResponseWriter, r *http.Request) {
    // 2. Increment the metric!
    requestsTotal.WithLabelValues("POST", "200").Inc()
    w.Write([]byte("Success"))
}

func main() {
    // 3. Expose the /metrics endpoint for Prometheus to scrape
    http.Handle("/metrics", promhttp.Handler())
    http.ListenAndServe(":8080", nil)
}
```

## 4. Grafana Dashboards and Alerts

Once Prometheus has scraped the numbers into its database, you connect **Grafana** to visualize them.

You write **PromQL** queries in Grafana to render graphs (e.g., `rate(http_requests_total[5m])` shows the Requests-Per-Second over the last 5 minutes).

Finally, you configure **AlertManager**. If the PromQL query detects that the Error Rate exceeds 5% for more than 2 minutes, it automatically sends a Slack message or triggers a PagerDuty call to wake up the On-Call Engineer.
