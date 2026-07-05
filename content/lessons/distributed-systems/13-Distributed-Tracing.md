# Distributed Tracing

In a Monolith, if an HTTP request fails, you look at the logs for that single server. The Stack Trace tells you exactly what line of code triggered the error.

In a Microservices architecture, a single HTTP request from a user might travel through the API Gateway, hit the Auth Service, fan-out to the Order and Inventory services, and hit 4 different databases.
If the request takes 5 seconds to load, which of those 6 components was slow? 

You cannot use a simple Stack Trace. You need **Distributed Tracing**.

## 1. Trace IDs and Span IDs

To track a request across network boundaries, we inject a unique identifier into the HTTP headers.

1. **Trace ID**: When a user hits the API Gateway, the Gateway generates a UUID (e.g., `trace-1234`). This represents the *entire* lifecycle of the request.
2. **Span ID**: A Span represents a single unit of work (e.g., "Querying the Postgres DB"). Every Span has a unique ID, but it belongs to the parent `trace-1234`.

When the API Gateway calls the Order Service, it injects `Trace-ID: trace-1234` into the HTTP Headers. The Order Service extracts it, and uses it for all of its own database spans!

## 2. OpenTelemetry (The Industry Standard)

Historically, tracing was fragmented (Jaeger, Zipkin, New Relic, Datadog all had their own SDKs).

Today, the entire industry has standardized on **OpenTelemetry (OTel)**. You instrument your Go application using the official OTel Go SDK.

```go
import "go.opentelemetry.io/otel"

func HandleCheckout(ctx context.Context) {
    // 1. Create a Span from the incoming Context!
    // The Context automatically carries the Trace-ID from the HTTP middleware!
    ctx, span := otel.Tracer("order-service").Start(ctx, "HandleCheckout")
    defer span.End()

    // 2. Add metadata to the span
    span.SetAttributes(attribute.String("user.id", "bob"))

    // 3. Pass the Context down to the database query!
    db.QueryContext(ctx, "SELECT ...")
}
```

## 3. The Context Requirement

Notice that tracing is 100% dependent on Go's `context.Context`! 

If you fail to pass the `ctx` into a function, or you accidentally create a brand new `context.Background()` in the middle of a request, the Trace is broken. The downstream services will generate a brand new Trace ID, and you will lose the ability to visualize the full request path in Jaeger or Datadog.

## 4. Visualizing the Trace

Once your Go services generate these Spans, they batch them up and send them over gRPC to a centralized Tracing Backend (like **Jaeger** or **Grafana Tempo**).

The Tracing Backend stitches all the Spans together using the `Trace-ID` and renders a beautiful Waterfall Chart in your browser.

* `[API Gateway] ------------------ (500ms)`
  * `[Auth Service] -- (50ms)`
  * `[Order Service] -------------- (450ms)`
    * `[Postgres Query] - (10ms)`
    * `[Stripe API] ------------- (440ms)` **<-- FOUND THE BOTTLENECK!**

Without Distributed Tracing, debugging performance in Microservices is literally impossible.
