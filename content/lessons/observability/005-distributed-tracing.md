# Distributed Tracing (Jaeger / Tempo)

## 1. Learning Objectives
* **What you'll learn**: How to visualize the entire lifespan of a single HTTP request as it hops across 15 different Go microservices using Distributed Tracing tools like Jaeger or Grafana Tempo.
* **Why it matters**: In a monolith, you can just read the stack trace to find a bug. In microservices, a user clicks "Buy", and the request travels through the Gateway, Auth, Billing, Inventory, and Shipping services. If it takes 5 seconds, how do you know *which* of those 5 services caused the delay? Tracing solves this.
* **Where it's used**: Complex microservice architectures where standard logs are impossible to stitch together manually.

---

## 2. Real-world Story
Imagine mailing a package from New York to Tokyo.
Without tracking, you give the package to the Post Office, and a week later, it hasn't arrived. You have absolutely no idea where it got lost.
Distributed Tracing is the Barcode on the box. 
Every time a facility touches the box (A Microservice), they scan the barcode (`Trace ID`), record exactly how many hours it sat in their facility (`Span Duration`), and upload it. You can look at a beautiful Gantt chart showing exactly which warehouse delayed your package!

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[API Gateway] -->|TraceID: 123, SpanID: A| B(Auth Service)
    B -.->|Returns OK (10ms)| A
    
    A -->|TraceID: 123, SpanID: B| C(Billing Service)
    
    C -->|TraceID: 123, SpanID: C| D[(Postgres DB)]
    C -->|TraceID: 123, SpanID: D| E(Stripe API)
    
    style C fill:#ef4444,color:#fff
    style E fill:#f59e0b,color:#fff
```

---

## 4. Internal Working (Under the Hood)
A **Trace** is a tree of **Spans**.
1. **Trace ID**: A globally unique 16-byte UUID generated at the very edge of your architecture (usually the API Gateway). This ID remains exactly the same across every microservice.
2. **Span ID**: An 8-byte ID identifying a single unit of work (e.g., `BillingService.ChargeCard`).
3. **Parent Span ID**: If Billing calls Postgres, the Postgres Span declares the Billing Span as its Parent. This allows Jaeger to mathematically rebuild the exact parent-child hierarchy tree to render the visualization!

---

## 5. Compiler Behavior
* **HTTP Headers (W3C Trace Context)**: When the Go Auth Service makes a network call to the Go Billing Service, Go cannot magically pass memory across the network. The OTel HTTP Client intercepts the Go network call and injects a standard HTTP Header: `traceparent: 00-1234567890abcdef-11223344-01`. The receiving Go service reads this header and resumes the trace!

---

## 6. Memory Management
* **Jaeger Storage**: Traces are massive. A single HTTP request might generate 50 spans (JSON objects). If you have 10,000 Req/Sec, Jaeger will generate Terabytes of data per hour. Therefore, Jaeger/Tempo rarely uses PostgreSQL. They use hyper-scalable NoSQL stores (Cassandra) or massive Cloud Object Storage (AWS S3) to handle the insane write-throughput.

---

## 7. Code Examples

### 🔹 Example 1: Propagating Context over HTTP
```go
import "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

// When making an outbound HTTP request to another microservice, 
// you CANNOT use standard http.Get(). You must use the OTel wrapper!
func CallBillingService(ctx context.Context) {
    client := http.Client{
        Transport: otelhttp.NewTransport(http.DefaultTransport),
    }

    req, _ := http.NewRequestWithContext(ctx, "GET", "http://billing-api/charge", nil)
    
    // The OTel Transport automatically pulls the TraceID from 'ctx' 
    // and injects it into the HTTP 'traceparent' header!
    resp, err := client.Do(req)
}
```

### 🔹 Example 2: Receiving the Context
```go
// On the receiving microservice, wrap the HTTP handler!
func main() {
    handler := http.HandlerFunc(ChargeHandler)
    
    // This middleware looks for the 'traceparent' header. 
    // If found, it creates a Child Span. If not found, it generates a brand new Trace ID!
    wrappedHandler := otelhttp.NewHandler(handler, "Billing.Charge")
    
    http.Handle("/charge", wrappedHandler)
}
```

### 🔹 Example 3: Advanced (Custom Span Events)
```go
func ChargeHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    span := trace.SpanFromContext(ctx)
    
    // You can attach "Events" (similar to logs) directly onto the visual timeline!
    span.AddEvent("Connecting to Stripe API")
    
    err := Stripe.Charge()
    if err != nil {
        // If it fails, mark the Span as visually RED on the Jaeger dashboard!
        span.SetStatus(codes.Error, "Stripe charge failed")
        span.RecordError(err)
    }
}
```

### 🔹 Example 4: Production (Jaeger UI Flame Graph)
```text
Trace: 12345 (Total Time: 1.5s)
|--- API Gateway (1.5s)
      |--- Auth Service (0.1s)
      |--- Billing Service (1.4s) 💥 ERROR!
            |--- Postgres DB (0.2s)
            |--- Stripe API (1.2s) 💥 Timeout!
```
*Just by looking at the graph, a Junior Engineer can instantly see that the 1.5s delay was purely caused by the Stripe API timing out!*

### 🔹 Example 5: Interview
```go
// Q: Why must you pass `context.Context` as the very first parameter to almost every function in Go?
// A: Because Go does not have "Thread Local Storage" like Java or Python! 
// The ONLY way to pass the TraceID down a chain of 10 nested functions is manually via the Context parameter.
```

---

## 8. Production Examples
1. **Performance Bottlenecking**: You notice a request takes 2 seconds. You open Jaeger. You see the `Postgres.Select` span is exactly 1.8 seconds. You immediately know you are missing a database Index.
2. **The N+1 Query Problem**: You open a trace and see `Postgres.Select` visually repeated 50 times in a cascading staircase pattern on the Flame Graph. You instantly realize you have an N+1 ORM bug!

---

## 9. Performance & Benchmarking
* **Grafana Tempo vs Jaeger**: Jaeger requires an indexing database (Elasticsearch/Cassandra) which is expensive to run. Grafana Tempo was built to solve this: it doesn't index anything except the Trace ID, allowing it to store petabytes of traces directly in cheap AWS S3 buckets!

---

## 10. Best Practices
* ✅ **Do**: Trace your Database calls. If you use GORM or `database/sql`, use the OTel wrappers (`otelsql`). DB queries are the #1 cause of latency.
* ❌ **Don't**: Create a new Span for every tiny Go function (e.g., `AddTwoNumbers()`). This generates massive trace payloads and network overhead. Only create Spans for network calls, database queries, and significant business logic blocks (e.g., `ProcessInvoice`).
* 🏢 **Google / Uber / Netflix Style**: Use **Exemplars**. When a Prometheus metric spikes (e.g., `API Latency > 2s`), Prometheus stores the exact Trace ID of the slow request. In Grafana, you just click the metric spike, and it instantly opens the exact trace that caused it!

---

## 11. Common Mistakes
1. **Dropping the Context**: You receive the `ctx` in your HTTP handler, but you spawn a background Goroutine `go ProcessEmail(context.Background())`. You just wiped the Trace ID! The background email worker will not show up on the Jaeger trace. You must use `trace.SpanFromContext` or pass the parent context safely.
2. **Missing W3C Headers**: Testing your API Gateway using Postman without generating a Trace ID. The Gateway generates a new one, but if you have a complex mesh, Postman won't link to the frontend React app's trace. Use OTel in your React frontend as well!

---

## 12. Debugging
How to troubleshoot Tracing in production:
* **The Missing Link**: If you look at Jaeger and see the API Gateway trace ends abruptly, but you know it called the Auth Service, it means the Auth Service is not extracting the HTTP Headers correctly. Check the `otelhttp` middleware setup on the Auth Service.

---

## 13. Exercises
1. **Easy**: Run Jaeger all-in-one locally using Docker. Access the UI on port 16686.
2. **Medium**: Write a Go HTTP server that initializes the OTel Jaeger Exporter and creates a simple span. View it in the Jaeger UI.
3. **Hard**: Create two Go HTTP servers. Server A calls Server B. Use `otelhttp.NewTransport` to successfully propagate the trace from A to B.
4. **Expert**: Create a nested span in Server B representing a fake database call, and mark it with a simulated error status.

---

## 14. Quiz
1. **MCQ**: How does a downstream microservice know it is part of an ongoing Trace?
   * (A) It asks the database (B) It generates a new ID (C) It reads the TraceID from the incoming HTTP `traceparent` header. *(Answer: C)*
2. **System Design Follow-up**: How do you trace Asynchronous Kafka messages? *(You inject the W3C Trace headers into the Kafka Message Payload / Headers. The Go Consumer reads the Kafka headers, extracts the TraceID, and continues the trace!)*

---

## 15. FAANG Interview Questions
* **Beginner**: What is the difference between a Trace and a Span?
* **Intermediate**: Explain W3C Trace Context propagation.
* **Senior (Google/Meta)**: Architect a tracing pipeline for 10 million Req/Sec. Explain the mathematical difference between Head-Based sampling (random coin flip at the Gateway) and Tail-Based sampling (analyzing the full trace before deciding to keep it).

---

## 16. Mini Project
**The Flame Graph Visualizer**
* Build a 3-tier Go architecture: Gateway -> Service A -> Service B.
* Instrument all 3 with OTel.
* Add an artificial `time.Sleep(500 * time.Millisecond)` to Service B.
* Trigger a request at the Gateway.
* Open Jaeger and take a screenshot of the beautiful 3-tier Flame Graph showing the exact 500ms delay at the bottom level!

---

## 17. Enterprise Features & Observability
* **Service Dependency Maps**: Because Jaeger knows exactly who calls who, it can dynamically auto-generate an architectural diagram of your entire company! If a junior dev wants to know what services depend on the Billing API, they just look at the auto-generated Jaeger Node Graph.

---

## 18. Source Code Reading
Walkthrough of `go.opentelemetry.io/otel/propagation`.
* **The Propagator**: Study how the `TraceContext` propagator literally parses the W3C HTTP header string `00-{trace_id}-{span_id}-01` and injects it into the immutable Go `context.Context` struct!

---

## 19. Architecture
* **Trace Analytics**: Companies like Honeycomb ingest all your traces and use AI to find correlations. "We noticed that 99% of all traces taking > 5 seconds have the attribute `tenant_id = 42`. You have a noisy neighbor problem!"

---

## 20. Summary & Cheat Sheet
* **Goal**: Visualize distributed requests.
* **TraceID**: The unique journey.
* **SpanID**: The single hop.
* **Go Requirement**: `context.Context` MUST be passed everywhere.
