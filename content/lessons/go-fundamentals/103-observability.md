# Observability (The Three Pillars)

Once an application is in production, you are completely blind unless you implement Observability. If a user complains "the checkout page is slow", you need data to prove *why*.

Observability is built on three pillars: **Logs, Metrics, and Traces**. 

Go's ecosystem is heavily integrated with the OpenTelemetry (OTel) standard, which provides a unified way to collect all three.

## 1. Logs (What happened?)

As we discussed in the Logging chapter, `log/slog` is used for discrete, timestamped events. 

* **Use Case**: "User 123 failed authentication." 
* **Cost**: High. Logs consume massive disk space, so they should be reserved for actionable events and errors, not every minor function call.

## 2. Metrics (How often is it happening?)

Metrics are aggregated, numerical data points collected over time. Because they are just counters in memory, they are incredibly cheap to track. The industry standard is **Prometheus**.

In Go, you expose a `/metrics` HTTP endpoint. Prometheus scrapes this endpoint every 10 seconds.

* **Counters**: Total number of HTTP requests processed.
* **Gauges**: Current number of active Goroutines or RAM usage.
* **Histograms**: The latency of database queries (e.g., 99th percentile response times).

```go
// Example using prometheus client
var reqCounter = promauto.NewCounter(prometheus.CounterOpts{
    Name: "myapp_http_requests_total",
    Help: "The total number of HTTP requests",
})

func handler(w http.ResponseWriter, r *http.Request) {
    reqCounter.Inc() // Increment the metric (Virtually zero CPU cost!)
    w.Write([]byte("Hello"))
}
```

## 3. Traces (Where is the bottleneck?)

If a request hits your API gateway, flows through three microservices, hits a Postgres database, and returns, how do you know which step took the longest?

**Distributed Tracing** solves this. 
When a request enters the system, a unique `TraceID` is generated. As the request moves through your Go application, that ID is passed through the `context.Context` to every single function.

Each function creates a "Span" (a start and end timestamp).

```mermaid
gantt
    title Distributed Trace (TraceID: a1b2c3d4)
    dateFormat  s
    axisFormat %S
    
    section API Gateway
    HTTP /checkout : 0, 5s
    
    section Order Service (Go)
    Validate Cart : 1s, 2s
    Call Payment Service : 3s, 2s
    
    section Payment Service (Go)
    Stripe API Call : 3s, 2s
```
By passing the `Context` down your entire call stack, OpenTelemetry can stitch these spans together into a single waterfall graph, allowing you to instantly see that the Stripe API Call was the bottleneck causing the 5-second delay!
