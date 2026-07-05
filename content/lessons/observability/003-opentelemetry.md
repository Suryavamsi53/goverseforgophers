# OpenTelemetry (OTel)

## 1. Learning Objectives
* **What you'll learn**: The architecture of OpenTelemetry, the industry-standard framework for generating, collecting, and exporting observability data (Metrics, Logs, Traces) from Go applications.
* **Why it matters**: Before OTel, if you used Datadog, you had to write Go code using the Datadog SDK. If you switched to New Relic, you had to rewrite 100% of your code! OTel provides a vendor-agnostic standard. Write the code once, and route the data anywhere.
* **Where it's used**: Adopted by Google, AWS, Microsoft, and the CNCF. It is the second most active CNCF project behind Kubernetes.

---

## 2. Real-world Story
Imagine standardizing electrical plugs. 
In the past, every TV required a proprietary wall socket. If you bought a Sony TV, you had to wire a Sony socket into your wall (Vendor Lock-in).
OpenTelemetry is the Universal Wall Socket. You write your Go application to plug into the OTel standard. Behind the wall, the OTel Collector can seamlessly route the electricity (data) to Datadog, Splunk, Prometheus, or Honeycomb without ever touching the TV (your Go code).

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Go App (OTel SDK)] -->|Generates Traces/Metrics| B[OTLP gRPC Protocol]
    
    B -->|Sends to| C{OTel Collector}
    
    C -->|Processes & Exports| D[(Jaeger / Tempo)]
    C -->|Processes & Exports| E[(Prometheus)]
    C -->|Processes & Exports| F[(Datadog SaaS)]
    
    style C fill:#8b5cf6,color:#fff
```

---

## 4. Internal Working (Under the Hood)
OpenTelemetry consists of three main pillars:
1. **The API & SDK**: The libraries you import into your Go code (`go.opentelemetry.io/otel`).
2. **OTLP (OpenTelemetry Protocol)**: A highly optimized, protobuf-based gRPC protocol used to transmit the telemetry data over the network rapidly.
3. **The Collector**: A standalone Go binary running in your infrastructure. It Receives data via OTLP, Processes it (scrubs PII, batches it), and Exports it to any vendor database.

---

## 5. Compiler Behavior
* **Context Propagation**: OTel in Go relies *entirely* on the standard library `context.Context`. The Trace ID and Span ID are invisibly embedded inside the `ctx`. If you fail to pass the `ctx` down the call chain of your Go functions, the trace is instantly broken, and the observability graph falls apart!

---

## 6. Memory Management
* **Batch Exporting**: The Go OTel SDK does NOT send an HTTP request every time a trace is generated, which would destroy network performance. It holds the traces in a Go slice in RAM, and a background Goroutine asynchronously flushes them to the Collector every 5 seconds (Batching), ensuring zero impact on your API latency.

---

## 7. Code Examples

### 🔹 Example 1: Simple (Initialization)
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

func initTracer() *trace.TracerProvider {
    // 1. Configure the Exporter (Where to send data via gRPC)
    exporter, _ := otlptracegrpc.New(context.Background())

    // 2. Configure the Provider (Batching logic)
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
    )
    
    // 3. Set as the Global standard
    otel.SetTracerProvider(tp)
    return tp
}
```

### 🔹 Example 2: Intermediate (Creating Spans)
```go
// The Context MUST be passed in!
func ProcessOrder(ctx context.Context, orderID string) {
    // 1. Create a "Span" representing this unit of work
    tracer := otel.Tracer("order-service")
    ctx, span := tracer.Start(ctx, "ProcessOrder")
    defer span.End() // Crucial: Ends the timer!

    // 2. Add custom business data (Attributes)
    span.SetAttributes(attribute.String("order.id", orderID))

    // 3. Pass the newly mutated context to the next function!
    ChargeCreditCard(ctx, orderID)
}
```

### 🔹 Example 3: Advanced (Auto-Instrumentation)
```go
// You don't have to manually write spans for database calls!
// OTel provides drop-in wrapper libraries for standard Go packages.
import "github.com/uptrace/opentelemetry-go-extra/otelsql"

// Just wrap standard sql.Open, and OTel will automatically trace EVERY SQL query!
db, err := otelsql.Open("postgres", "postgres://user:pass@localhost:5432/db")

// The trace will perfectly capture: "SELECT * FROM users", latency: 5ms
```

### 🔹 Example 4: Production (The Collector YAML)
```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      grpc:

processors:
  batch: # Batch data for performance
  attributes/pii:
    actions:
      # Automatically scrub sensitive data before it reaches Datadog!
      - action: hash
        key: user.email

exporters:
  prometheus:
    endpoint: "0.0.0.0:8889"
  datadog:
    api:
      key: "my-secret-key"

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [batch, attributes/pii]
      exporters: [datadog]
```

### 🔹 Example 5: Interview
```go
// Q: Why run an OTel Collector instead of sending data straight from the Go app to Datadog?
// A: Separation of concerns! The Collector handles API Key rotation, network retries, 
// data scrubbing, and vendor switching. Your Go code remains completely ignorant of infrastructure details.
```

---

## 8. Production Examples
1. **Vendor Migration**: A company's Datadog bill reaches $1,000,000/year. Because they used OTel, they change 3 lines in the Collector YAML to route to Prometheus/Jaeger instead. They migrate in 5 minutes with zero code changes, saving a million dollars.
2. **Polyglot Microservices**: A Go service calls a Python service via gRPC. OTel injects the `trace_parent` ID into the gRPC metadata headers. The Python OTel SDK extracts it, perfectly linking the Go and Python traces together seamlessly.

---

## 9. Performance & Benchmarking
* **Sampling**: Generating a trace for every single HTTP request on a 10,000 Req/Sec Go API will melt your database. OTel allows you to configure "Head-Based Sampling" (e.g., only trace 1% of requests).

---

## 10. Best Practices
* ✅ **Do**: Name your spans dynamically but concisely (e.g., `UpdateUser`, `QueryPostgres`).
* ❌ **Don't**: Include high-cardinality data in the span name (e.g., `UpdateUser_12345`). The user ID belongs in the Span Attributes, not the Name!
* 🏢 **Google / Uber / Netflix Style**: Use "Tail-Based Sampling" in the Collector. The Collector stores 100% of traces in RAM for 10 seconds. If a trace completes normally, it drops it. If a trace contains an `Error` or takes `> 5 seconds`, it keeps it and exports it!

---

## 11. Common Mistakes
1. **Forgetting `defer span.End()`**: If you forget to end the span, the memory leaks in the Go SDK, the trace is never dispatched, and it looks like the function hung infinitely on the dashboard.
2. **Context Bleeding**: Using `context.Background()` deep inside a nested function instead of passing the parent `ctx`. The trace instantly snaps in half, and you get two disconnected fragments on the dashboard instead of one unified flow.

---

## 12. Debugging
How to troubleshoot OTel in production:
* **Logging Exporter**: If traces aren't appearing in Datadog, change the Collector config to export to `logging`. The Collector will simply print the raw JSON traces to `stdout`, proving whether the Go app is actually sending data.

---

## 13. Exercises
1. **Easy**: Initialize the OTel SDK in a Go `main.go` file.
2. **Medium**: Create a span inside a function and add an Attribute (`http.status_code = 200`).
3. **Hard**: Use the `otelhttp` middleware wrapper to automatically generate spans for every incoming HTTP request to your Go server.
4. **Expert**: Spin up an OTel Collector using Docker Compose, point your Go app to it via gRPC, and configure the Collector to output the traces to a local Jaeger container!

---

## 14. Quiz
1. **MCQ**: What is the primary benefit of OpenTelemetry?
   * (A) It is a database (B) It prevents vendor lock-in (C) It replaces Kubernetes. *(Answer: B)*
2. **Code Review**: `span.SetAttributes(attribute.String("password", user.Pass))`. What is wrong with this? *(NEVER put passwords, PII, or auth tokens in telemetry attributes! They will be sent in plaintext to your observability vendor, violating GDPR and compliance).*

---

## 15. FAANG Interview Questions
* **Beginner**: Explain the three pillars of observability (Metrics, Logs, Traces).
* **Intermediate**: How does Distributed Tracing propagate context across HTTP boundaries? (Hint: W3C Trace Context HTTP Headers).
* **Senior (Google/Meta)**: Architect a telemetry pipeline for a global scale system processing 5 million spans per second. Explain how you would implement Tail-Based sampling using an OTel Collector cluster backed by Kafka.

---

## 16. Mini Project
**The Polyglot Tracer**
* Write a Go API (Service A).
* Write a simple Node.js API (Service B).
* Service A makes an HTTP request to Service B.
* Use OTel in both. Configure the HTTP clients to inject the W3C Trace headers.
* Open Jaeger and view the beautiful flame graph seamlessly spanning across two completely different programming languages!

---

## 17. Enterprise Features & Observability
* **Baggage**: OTel supports `Baggage`, a mechanism to pass key-value pairs down the entire trace tree. If the API Gateway sets `tenant_id=xyz`, every single microservice downstream can instantly access that `tenant_id` without needing to add it to their function signatures!

---

## 18. Source Code Reading
Walkthrough of `go.opentelemetry.io/otel`.
* **The Global State**: Study how OTel uses global variables cautiously (`otel.GetTracerProvider()`). This allows any package in your entire codebase to generate spans without needing to pass a `Tracer` object through 50 layers of function arguments.

---

## 19. Architecture
* **The eBPF Revolution**: The future of OTel is eBPF. Instead of manually writing `span.Start()` in your Go code, an eBPF program running in the Linux Kernel intercepts the Go network syscalls and automatically generates traces with literally zero code changes to your application!

---

## 20. Summary & Cheat Sheet
* **OTel**: Vendor-agnostic standard.
* **Trace**: A complete journey of a request.
* **Span**: A single unit of work in that journey.
* **Collector**: The router/scrubber in the middle.
