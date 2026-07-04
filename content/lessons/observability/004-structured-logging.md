# Structured Logging Pipeline

We covered `log/slog` in the Go Fundamentals track, but how does structured logging fit into the broader observability pipeline?

## 1. The Problem with Text Logs

If your Go server logs plain text:
`2023/10/01 10:00:00 Failed to process order 9923 for user alice@gmail.com`

When this log arrives in Elasticsearch or Splunk, it is treated as a single, massive string. If you want to build a graph showing "Failed orders per user", the log aggregation system has to run complex Regex operations on millions of logs to extract the email address. This is incredibly slow and burns massive CPU.

## 2. JSON Ingestion

By using Go's `slog.NewJSONHandler`, your output looks like this:
```json
{
  "time": "2023-10-01T10:00:00Z",
  "level": "ERROR",
  "msg": "Failed to process order",
  "orderID": 9923,
  "userEmail": "alice@gmail.com"
}
```

When Elasticsearch receives this JSON payload, it doesn't just save a string. It dynamically parses the JSON and creates indexed, searchable database columns for `orderID` and `userEmail`. 

Querying `userEmail == "alice@gmail.com"` is now an indexed database lookup, returning results in milliseconds instead of minutes.

## 3. Standardizing Keys (The Taxonomy)

If Team A logs `"user_id": 123` and Team B logs `"userID": 123`, your log aggregator will create two different columns, making cross-team queries impossible.

Enterprise architectures require a strict **Logging Taxonomy**. You define constants in a shared Go package:

```go
// pkg/logkeys/keys.go
const (
    KeyUserID  = "user_id"
    KeyOrderID = "order_id"
    KeyLatency = "latency_ms"
)
```
Teams are forced to use these keys when writing logs:
```go
slog.Info("order saved", slog.Int(logkeys.KeyOrderID, 99))
```

## 4. Contextual Log Injection

In an HTTP server, every log line generated during a request should contain the `TraceID` and the `UserID`. Instead of manually adding these to every single `slog.Info()` call, you inject a pre-configured logger into the `context.Context` using Middleware.

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        traceID := generateTraceID()
        
        // Create a logger bound with the traceID
        logger := slog.Default().With(slog.String("trace_id", traceID))
        
        // Inject the logger into the request context
        ctx := context.WithValue(r.Context(), "logger", logger)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```
Now, downstream repository functions can extract the logger from the context and use it, guaranteeing that every database error log perfectly correlates with the original HTTP request!
