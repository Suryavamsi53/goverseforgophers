# Centralized Logging

If you have 50 instances (Pods) of a Go Microservice running in Kubernetes, and a user complains that they got a `500 Internal Server Error`, how do you find the log? 

You cannot SSH into 50 different servers to run `cat /var/log/app.log`. 

In a distributed system, logs must be shipped to a centralized, searchable database.

## 1. The 12-Factor App (Stdout)

The first rule of cloud-native logging is: **Your Go application should never write to a log file.**

Writing to files creates massive problems in Kubernetes (managing disk space, log rotation, permissions). Instead, the "12-Factor App" methodology dictates that your application should ONLY write to standard output (`os.Stdout` / `fmt.Println`).

In Kubernetes, the container runtime (Docker/containerd) automatically captures `Stdout` from your Go app. A background "DaemonSet" (like Fluent Bit or Promtail) tails those container logs and streams them over the network to a central database.

## 2. Structured Logging (JSON)

If you log:
`log.Printf("User %d failed to login from IP %s", userID, ip)`

The central database receives a raw string: `"User 42 failed to login from IP 192.168.1.1"`.
If you want to search the database for "Show me all failed logins from 192.168.1.1", the database has to do complex, slow Regex text parsing.

To fix this, Enterprise Go applications use **Structured Logging**. Every log is output as a JSON object.

```go
import "log/slog" // Built into Go 1.21+

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    logger.Error("failed to login", 
        slog.Int("user_id", 42),
        slog.String("ip", "192.168.1.1"),
    )
}
```
**Output:**
`{"time":"2026-07-01T12:00:00Z","level":"ERROR","msg":"failed to login","user_id":42,"ip":"192.168.1.1"}`

Now, the centralized database instantly indexes the `user_id` and `ip` fields, making searches lightning fast.

## 3. Correlation IDs

If you have millions of JSON logs flowing into your database, how do you group logs that belong to a single user's request?

You inject the **Trace ID** (from OpenTelemetry) into every single log statement! This is often called a Correlation ID.

```go
logger.Info("charging card", 
    slog.String("trace_id", ctx.Value("TraceID").(string)),
)
```

Now, when a user reports a bug, you take their `TraceID` and query your logging dashboard (e.g., Kibana): `trace_id: "abc-123"`. 
You instantly see all 15 logs across 4 different microservices that were generated during that exact millisecond, printed in perfect chronological order.

## 4. The ELK / PLG Stacks

The two most popular centralized logging stacks in the industry are:
1. **ELK Stack**: Elasticsearch (Database), Logstash (Collector), Kibana (Dashboard). Extremely powerful, but very heavy and expensive.
2. **PLG Stack**: Promtail (Collector), Loki (Database), Grafana (Dashboard). Created by Grafana, it is extremely lightweight because it does not index the log text (it only indexes the JSON labels). This is the modern standard for Go/Kubernetes environments.
