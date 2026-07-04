# Structured Logging

Logging `fmt.Println("Error occurred")` is perfectly fine for a local script. But in a production enterprise environment, logs are ingested by massive aggregating systems (like Datadog, ELK, or Splunk). 

If you use plain text logs, those systems cannot easily search, filter, or alert on your data. You must use **Structured Logging**.

## 1. The Standard `log` Package

The built-in `log` package is simple and thread-safe.

```go
import "log"

func main() {
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
    log.Println("Starting server...") 
    // Output: 2026/07/04 12:00:00 main.go:5: Starting server...
}
```
While useful, this is unstructured text. If you want to filter all logs for `userID=123`, you have to use complex, fragile Regex parsers.

## 2. Enter `slog` (Go 1.21+)

For years, the community relied on third-party libraries like Uber's `zap` or `logrus` for structured logging. In Go 1.21, the Go team officially introduced the `log/slog` package to the standard library.

`slog` outputs logs as structured key-value pairs (usually JSON).

```go
import (
    "log/slog"
    "os"
)

func main() {
    // 1. Create a JSON Handler pointing to standard output
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    
    // 2. Set it as the global default logger
    slog.SetDefault(logger)

    // 3. Log with structured key-value pairs!
    slog.Info("user logged in", 
        slog.String("userID", "user_123"), 
        slog.Int("attempt", 1),
    )
}
```

**The Output:**
```json
{"time":"2026-07-04T12:00:00Z","level":"INFO","msg":"user logged in","userID":"user_123","attempt":1}
```

Now, your aggregation tools can instantly query `level == "ERROR"` or `userID == "user_123"` with zero parsing overhead!

## 3. Contextual Logging

In large web servers, you want every log line generated during a request to automatically include the `TraceID` so you can group them together.

You achieve this by extracting the TraceID from the `context.Context` and injecting it into the logger.

```go
func processOrder(ctx context.Context, orderID int) {
    // Extract the trace ID (injected earlier by middleware)
    traceID := ctx.Value("traceID").(string)
    
    // Create a sub-logger that ALWAYS includes this traceID
    reqLogger := slog.Default().With(slog.String("traceID", traceID))
    
    reqLogger.Info("processing order", slog.Int("orderID", orderID))
    reqLogger.Error("database timeout") 
    // Both logs will have the exact same traceID attached!
}
```
