# Centralized Logging

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* The Problem with `fmt.Println`
* Structured Logging
* The Logging Pipeline (ELK / EFK Stack)
* Architecture Diagram
* Step-by-Step Implementation (`slog`)
* Beginner Example
* Intermediate Example (Correlating Traces and Logs)
* Production Use Cases
* Best Practices
* Common Mistakes
* Exercises
* Quiz
* Interview Questions
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

In a local environment, when something goes wrong, you look at your terminal and read the logs. 

In a distributed system, you might have 50 microservices, each running 3 instances across Kubernetes pods. That's 150 independent servers generating text simultaneously. You cannot SSH into 150 servers to `tail -f` their log files. 

**Centralized Logging** solves this by shipping all logs from every server to a single, searchable database (like Elasticsearch). Furthermore, to make these logs searchable by machines, we abandon plain text and embrace **Structured Logging**.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand the architecture of a centralized logging pipeline.
* Explain why unstructured plain-text logging is useless at scale.
* Use Go's built-in `log/slog` package to emit structured JSON logs.
* Inject Trace IDs into logs to seamlessly correlate Logging with Distributed Tracing.

---

# Prerequisites

Before reading this chapter you should know:

* Distributed Tracing (`13-Distributed-Tracing.md`).

---

# Why This Topic Exists

Imagine you write this in your code:
`fmt.Printf("User %s failed to pay $%d\n", username, amount)`

It outputs: `User alice123 failed to pay $50`.
This log is shipped to your centralized logging database. Later, your CEO asks: "How many users failed to pay more than $40 today?"

To answer this, you have to write a complex Regex query to extract the number after the `$` sign from a giant wall of plain text. It is incredibly slow and fragile. If another developer changes the log to `"Payment failed for user alice123, amount: $50"`, your Regex breaks, and your dashboards fail.

You need **Structured Logging**. Instead of a string, you emit a JSON object:
`{"level": "error", "event": "payment_failed", "user": "alice123", "amount": 50}`
Now, answering the CEO's question is a simple SQL-like query: `SELECT count(*) WHERE event="payment_failed" AND amount > 40`.

---

# Structured Logging

Structured logging forces developers to treat logs as Data Data, not Human Text. 
Every log entry is a JSON object (or similar key-value format). 
* **The Message**: A static, hardcoded string describing the event (e.g., "payment_failed"). NEVER put dynamic variables inside the message.
* **The Attributes**: Key-value pairs containing the dynamic context (e.g., `user_id: 123`, `latency_ms: 45`).

As of Go 1.21, Go includes a highly optimized, standard library package for structured logging called **`log/slog`**.

---

# The Logging Pipeline (ELK Stack)

Shipping logs to a central database involves an architecture commonly known as the ELK or EFK stack.

1. **The Application**: Your Go binary outputs JSON logs to `stdout`. (It does NOT connect directly to the database!).
2. **The Collector / Shipper (Fluentd, Promtail, Filebeat)**: A background agent running on the server that reads the `stdout` stream, batches the JSON logs, and securely ships them across the network.
3. **The Storage (Elasticsearch / Loki)**: A massive, distributed NoSQL database optimized for indexing and searching billions of JSON documents.
4. **The UI (Kibana / Grafana)**: The web interface where engineers write queries to search the logs.

---

# Architecture Diagram

```mermaid
flowchart LR
    Go1[Go App 1<br/>(stdout)]
    Go2[Go App 2<br/>(stdout)]
    
    Agent[Log Shipper<br/>(Fluentd)]
    
    DB[(Elasticsearch / Loki)]
    UI[Kibana / Grafana]

    Go1 -- "JSON Logs" --> Agent
    Go2 -- "JSON Logs" --> Agent
    
    Agent -- "Batches over Network" --> DB
    
    UI -- "Searches" --> DB
```

---

# Step-by-Step Implementation (`slog`)

1. Import `log/slog` and `os`.
2. Create a JSON Handler: `handler := slog.NewJSONHandler(os.Stdout, nil)`.
3. Create the Logger: `logger := slog.New(handler)`.
4. (Optional) Set it as the global default: `slog.SetDefault(logger)`.
5. Emit logs using `logger.Info()`, `logger.Error()`, providing the static message and alternating Key/Value pairs.

---

# Beginner Example

Replacing `fmt.Printf` with `slog` for structured JSON logging.

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	// 1. Configure the logger to output JSON to standard out
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	userID := 1024
	amount := 50.0

	// THE BAD WAY (Unstructured)
	// You should avoid this in distributed systems!
	// fmt.Printf("User %d paid $%.2f\n", userID, amount)

	// THE GOOD WAY (Structured)
	// Notice the message is a static string. Variables are passed as Key-Value pairs.
	logger.Info("payment_processed", 
		slog.Int("user_id", userID),
		slog.Float64("amount", amount),
	)

	// Warning with attributes
	logger.Warn("retry_attempted",
		slog.String("service", "stripe"),
		slog.Int("attempt", 2),
	)
}
```
*Output:*
```json
{"time":"2026-08-15T10:00:00Z","level":"INFO","msg":"payment_processed","user_id":1024,"amount":50}
{"time":"2026-08-15T10:00:00Z","level":"WARN","msg":"retry_attempted","service":"stripe","attempt":2}
```

---

# Intermediate Example (Correlating Traces and Logs)

If you have Distributed Tracing (Chapter 13) and Centralized Logging, they must talk to each other. When you emit a log, you should inject the current `Trace ID` into the JSON log. 
This allows you to see a spike in a Grafana dashboard, click the specific error log, and instantly jump to the Jaeger trace showing the entire request journey!

```go
package main

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel/trace"
)

// A custom wrapper that extracts the TraceID from the Context
// and injects it into the slog JSON attributes.
func LogWithTrace(ctx context.Context, logger *slog.Logger, level slog.Level, msg string, args ...any) {
	// 1. Extract the active Span from the context
	spanContext := trace.SpanContextFromContext(ctx)
	
	// 2. If a valid trace exists, inject the trace_id into the log arguments
	if spanContext.HasTraceID() {
		traceID := spanContext.TraceID().String()
		args = append(args, slog.String("trace_id", traceID))
	}

	// 3. Emit the log
	logger.Log(ctx, level, msg, args...)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background() // In reality, this context holds the OTel trace

	// Simulate emitting an error log during a request
	LogWithTrace(ctx, logger, slog.LevelError, "database_connection_failed",
		slog.String("db_host", "localhost:5432"),
	)
}
```

---

# Production Use Cases

### 1. Security Auditing (SIEM)
Security Information and Event Management (SIEM) systems rely entirely on structured logs. If an application logs `{"event": "login_failed", "ip": "192.168.1.5"}`, the SIEM database can automatically run anomaly detection. If it sees 5,000 `login_failed` events from the same IP within 1 minute, it automatically updates firewall rules to block the IP.

### 2. Ephemeral Kubernetes Pods
In Kubernetes, when a Pod crashes, the container is completely deleted, along with its local hard drive. If you write your logs to a local file (e.g., `app.log`), those logs are permanently lost the moment you need them most (during a crash!). By writing logs to `stdout`, the Kubernetes daemonset (Fluentd) captures them in real-time and safely stores them in Elasticsearch *before* the pod is deleted.

---

# Best Practices

* **Do not log PII**: Never log Personally Identifiable Information (Passwords, Social Security Numbers, full Credit Card numbers). Logs are often widely accessible to many engineers in the company. Logging PII is a severe security violation.
* **Log to `stdout` (12-Factor App methodology)**: Never write logic in your Go application to connect to Elasticsearch directly, and never manage log file rotation on disk. Just print to standard out. Let the infrastructure (Docker/Kubernetes) handle the routing and shipping.
* **Use Static Messages**: `logger.Info(fmt.Sprintf("User %s logged in", name))` ruins the purpose of structured logging. The message should be exactly `"user_logged_in"`.

---

# Common Mistakes

### Logging Too Much (The Wallet Drainer)
Elasticsearch and Datadog charge you by the gigabyte for log ingestion and storage. If you log `"entering function"` and `"exiting function"` for every single HTTP request, a high-throughput microservice will generate Terabytes of logs a day, bankrupting your startup. 
Only log actionable events (Errors, Warnings, and major business milestones like "Order Placed").

---

# Quiz

## Multiple Choice Questions
**1. Why should you log to `stdout` instead of directly to a file like `var/log/app.log` in a modern distributed environment?**
A) Because writing to files is not allowed in Go.
B) Containers in the cloud are ephemeral. If the container crashes and is deleted, the local file is lost. By logging to `stdout`, infrastructure agents (like Promtail) can stream the logs safely to an external database in real-time.
C) `stdout` is faster.
*Answer*: B

## True or False
**Injecting the `trace_id` into your structured logs is considered a bad practice because it clutters the logs with unnecessary UUIDs.**
*Answer*: False. Injecting the `trace_id` is the industry standard best practice. It is the only way to correlate a specific log entry in Elasticsearch with the visual distributed trace in Jaeger, bridging the gap between logs and traces.

---

# Interview Questions

## Beginner
**Q**: What is Structured Logging and why is it superior to plain text logging?
*Answer*: Structured logging outputs logs as machine-readable data structures (usually JSON) consisting of a static event name and dynamic key-value attributes. It is superior because it allows centralized logging databases to index the keys, enabling fast, complex querying (like filtering by user ID or sorting by latency) without relying on fragile string parsing or regex.

## Intermediate
**Q**: According to the 12-Factor App principles, how should a Go application handle routing logs to an external database like Elasticsearch?
*Answer*: It shouldn't. The application should treat logs as an event stream and write them completely unbuffered to `stdout`. The execution environment (like a Kubernetes sidecar or daemonset) takes responsibility for capturing that stream, batching it, and shipping it to Elasticsearch. This decouples the application from the logging infrastructure.

## Advanced
**Q**: Explain how `log/slog` in Go 1.21 improves performance compared to passing a generic `map[string]interface{}` to a JSON encoder.
*Answer*: Passing a generic map requires the Go runtime to use Reflection to determine the types of the values at runtime, and forces heap allocations. `slog` uses strongly-typed attributes (e.g., `slog.Int("key", 5)`, `slog.String("key", "val")`). These structs do not require reflection, minimizing CPU overhead and preventing garbage collection pressure in high-throughput applications.

---

# Summary

Centralized, structured logging is the foundation of backend observability. By adopting `slog` and treating logs as JSON events rather than human-readable paragraphs, you transform your application's output from a chaotic wall of text into a powerful, queryable database that can alert you to production issues before your users even notice.

---

# Key Takeaways

* ✔ Plain text logs are unscalable. Use Structured JSON Logging.
* ✔ Emit a static message, use Key-Value pairs for variables.
* ✔ Use Go's built-in `log/slog` package.
* ✔ Always log to `stdout` in containerized environments.
* ✔ Inject Trace IDs into logs for perfect correlation.

---

# Further Reading
* [A Guide to the Go slog Package](https://betterstack.com/community/guides/logging/logging-in-go/#structured-logging-with-slog)
* [The 12-Factor App: Logs](https://12factor.net/logs)

---

# Next Chapter
➡️ **Next:** `15-Metrics-and-Monitoring.md`
