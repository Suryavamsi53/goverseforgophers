# Metrics: RED and USE Methods

When instrumenting a Go application with metrics, you shouldn't just guess what to measure. If you measure too much, you bloat your Prometheus database and create "alert fatigue". If you measure too little, you miss critical outages.

Site Reliability Engineers (SREs) rely on two universal frameworks: The **RED Method** for services, and the **USE Method** for infrastructure.

## 1. The RED Method (Services)

The RED Method dictates exactly what you must measure for every single microservice and API endpoint.

* **Rate**: The number of requests your service is receiving per second.
  * *Metric Type*: Prometheus Counter (`http_requests_total`)
* **Errors**: The number of those requests that are failing (e.g., HTTP 5xx codes).
  * *Metric Type*: Prometheus Counter (`http_requests_total{status="500"}`)
* **Duration**: The time it takes to process a request (Latency).
  * *Metric Type*: Prometheus Histogram (`http_request_duration_seconds`)

If you have a Grafana dashboard showing these three metrics, you can instantly answer the most important question during an incident: *"Are users experiencing pain right now?"*

## 2. The USE Method (Infrastructure)

While RED monitors the *user experience*, the USE method monitors the *underlying hardware and resources* (Servers, Databases, Disks).

* **Utilization**: The average time the resource was busy servicing work (e.g., CPU is at 90% capacity).
* **Saturation**: The degree to which the resource has extra work queued up that it cannot process (e.g., Linux Load Average, or a Go Channel that is completely full).
* **Errors**: The count of error events (e.g., Network interface dropping packets, or Disk read failures).

## 3. High-Cardinality Explosions (The Trap)

When defining Prometheus metrics in Go, you attach "Labels" (Tags) to provide context.

```go
var requestCount = promauto.NewCounterVec(
    prometheus.CounterOpts{Name: "http_requests"},
    []string{"method", "status", "endpoint"}, // Labels
)

// usage: requestCount.WithLabelValues("GET", "200", "/users").Inc()
```
This allows you to query the error rate for specific endpoints!

**⚠️ THE TRAP:**
Never put unbounded, highly dynamic data (like User IDs or IP Addresses) into a metric label!

```go
// ❌ CATASTROPHIC BUG
[]string{"method", "status", "user_id"}
```
Prometheus creates a brand new time-series database in memory for *every unique combination of labels*. If you have 1,000,000 users, Prometheus will instantly try to create 1,000,000 databases in RAM, immediately crashing the Prometheus server with an Out-Of-Memory (OOM) panic!

**Metrics are for aggregate trends. Logs are for high-cardinality specifics.**
