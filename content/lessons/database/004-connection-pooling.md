# Connection Pooling & PgBouncer

## 1. Learning Objectives
* **What you'll learn**: The mechanics of database connection pooling in Go, how `database/sql` safely manages concurrent Goroutines, and when to introduce PgBouncer for massive horizontal scaling.
* **Why it matters**: A single PostgreSQL backend process takes up to 10MB of RAM. If 10,000 Go Goroutines open 10,000 direct connections, Postgres will run out of memory (OOM) and crash instantly. Pooling is mandatory for stability.
* **Where it's used**: Every high-throughput Go backend, specifically when deployed on autoscaling Kubernetes clusters.

---

## 2. Real-world Story
Imagine a busy restaurant (PostgreSQL) with 50 tables (Connections). 
If 500 customers (Go Goroutines) walk in at the same time and demand their own private table, the restaurant cannot seat them all. 
Instead, the Hostess (The Connection Pool) manages the tables. When a customer walks in, the Hostess seats them at an empty table. They eat, leave, and the Hostess immediately gives that exact same table to the next waiting customer. 50 tables can efficiently serve 500 customers by *sharing and reusing* the physical space!

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Go Goroutine 1] -->|Needs DB| B{Go Connection Pool}
    C[Go Goroutine 2] -->|Needs DB| B
    D[Go Goroutine 3] -->|Needs DB| B
    
    B -->|Check out| E((Physical TCP Connection 1))
    B -->|Check out| F((Physical TCP Connection 2))
    
    E -->|Wait for query...| G[(PostgreSQL)]
    F -->|Return to pool| B
    
    style B fill:#3b82f6,color:#fff
```

---

## 4. Internal Working (Under the Hood)
When you call `db.Query()`, Go does not open a new TCP connection to Postgres. 
1. Go locks a Mutex and checks a slice of `idleConnections`.
2. If an idle connection exists, Go gives it to your Goroutine.
3. If no idle connections exist, Go checks if the `MaxOpenConns` limit is reached.
4. If not reached, it opens a *new* physical TCP connection.
5. If the limit IS reached, your Goroutine physically goes to sleep (blocks) and waits in a queue until another Goroutine calls `rows.Close()`, which returns its connection to the pool!

---

## 5. Compiler Behavior
* **Thread-Safety via Channels/Mutexes**: The `database/sql` package is one of the most brilliant examples of concurrent Go engineering. It uses a combination of `sync.Mutex` for fast state checks, and `chan connRequest` to wake up sleeping Goroutines exactly when a connection becomes available.

---

## 6. Memory Management
* **Idle Connections**: Holding 100 idle connections open means Postgres is wasting 1GB of RAM for absolutely no reason. By setting `SetMaxIdleConns(10)`, Go will intelligently close physical TCP sockets when traffic drops, freeing up RAM on your Postgres server.

---

## 7. Code Examples

### 🔹 Example 1: Simple
```go
// The standard database/sql Connection Pool setup in main.go
db, err := sql.Open("pgx", "postgres://user:pass@localhost/db")
if err != nil { log.Fatal(err) }

// VERY IMPORTANT: Configure the pool!
db.SetMaxOpenConns(50)    // Max total connections
db.SetMaxIdleConns(10)    // Max idle connections left open
db.SetConnMaxLifetime(time.Hour) // Close connections after 1 hour safely
```

### 🔹 Example 2: Intermediate
```go
// The pgxpool (Native PostgreSQL) Configuration
config, _ := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
config.MaxConns = 50
config.MinConns = 10
config.MaxConnLifetime = time.Hour

pool, _ := pgxpool.NewWithConfig(context.Background(), config)
```

### 🔹 Example 3: Advanced
```go
// Context Cancellation handling in the Pool queue
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// If the pool is full (MaxOpenConns=50) and all 50 are busy for 3 seconds,
// THIS Query will automatically abort after 2 seconds because of the context timeout,
// freeing up the Goroutine instead of deadlocking!
err := pool.QueryRow(ctx, "SELECT * FROM users").Scan(&u)
```

### 🔹 Example 4: Production
```go
// Emitting Pool Statistics to Prometheus!
stats := db.Stats()
fmt.Printf("In Use: %d, Idle: %d, Wait Count: %d\n", 
    stats.InUse, stats.Idle, stats.WaitCount)
// WaitCount > 0 means Goroutines are queuing up! Your pool is too small!
```

### 🔹 Example 5: Interview
```go
// Q: Why do you need `SetConnMaxLifetime`?
// A: Corporate firewalls, AWS Load Balancers, or HAProxy will often silently terminate 
// idle TCP connections after 5 minutes. If Go doesn't know the connection is dead, 
// the next query will fail with a "broken pipe" error! Setting a lifetime ensures Go gracefully cycles connections.
```

---

## 8. Production Examples
1. **Kubernetes Autoscaling**: If you have 100 Go pods, and each pod has a `MaxOpenConns` of 50, your Postgres DB might suddenly receive 5,000 connections! This will crush Postgres. This is where PgBouncer is introduced.
2. **Serverless (AWS Lambda)**: Lambda functions spin up and die in seconds. They cannot maintain a persistent connection pool. You must use a database proxy (like AWS RDS Proxy).

---

## 9. Performance & Benchmarking
* **The PgBouncer Multiplier**: PgBouncer sits *between* Go and Postgres. 100 Go pods can open 5,000 connections to PgBouncer, but PgBouncer multiplexes those requests down to just 100 physical connections to Postgres. It prevents Postgres from OOMing under massive scale.

---

## 10. Best Practices
* ✅ **Do**: Tune `MaxOpenConns`. A value between `20` and `50` is usually optimal per Go server instance, regardless of how many thousands of requests you get.
* ❌ **Don't**: Let `SetMaxIdleConns` be larger than `SetMaxOpenConns`. Go will silently truncate it.
* 🏢 **Google / Uber / Netflix Style**: For hyper-scale, use `pgbouncer` in "Transaction Pooling" mode.

---

## 11. Common Mistakes
1. **Leaking Connections**: `rows, _ := db.Query()`. If you do not call `rows.Close()`, the connection is permanently checked out of the pool. Do this 50 times, and your entire Go server freezes.
2. **Infinite Lifetimes**: Relying on the default `SetConnMaxLifetime(0)`. This causes "connection reset by peer" errors during network blips.

---

## 12. Debugging
How to troubleshoot Connection Pooling in production:
* **The "FATAL: remaining connection slots are reserved" error**: This means your Go apps have opened too many connections. You must lower `MaxOpenConns` or install PgBouncer.
* **Go App Freezes without crashing**: Check `db.Stats().WaitCount`. If it's soaring, it means you leaked a connection somewhere (forgot a `rows.Close()`), and all Goroutines are asleep in the queue.

---

## 13. Exercises
1. **Easy**: Configure a standard `database/sql` pool with 25 Max and 5 Idle conns.
2. **Medium**: Spin up 500 Goroutines that execute a query with a 50ms `pg_sleep`. Use `db.Stats()` to prove that Go only opened 25 physical connections and queued the rest.
3. **Hard**: Install PgBouncer locally using Docker, and configure your Go app to connect to PgBouncer instead of directly to Postgres.
4. **Expert**: Intentionally leak a connection by omitting `rows.Close()`, and write a diagnostic tool using `db.Stats()` to detect the leak.

---

## 14. Quiz
1. **MCQ**: What happens when a Goroutine requests a connection, but `MaxOpenConns` is reached and all connections are busy?
   * (A) `db.Query` panics (B) Go opens a new one anyway (C) The Goroutine blocks and waits. *(Answer: C)*
2. **System Design Follow-up**: Why does "Transaction Pooling" in PgBouncer break Go's Prepared Statements? *(Because Prepared Statements are tied to a physical TCP Session. Transaction pooling swaps the TCP connection out from under you between queries!)*

---

## 15. FAANG Interview Questions
* **Beginner**: Why not just open a new connection for every HTTP request?
* **Intermediate**: Explain the difference between Session Pooling and Transaction Pooling in PgBouncer.
* **Senior (Google/Meta)**: Architect the database connectivity for 10,000 microservices connecting to a single monolithic PostgreSQL database cluster.

---

## 16. Mini Project
**The Load Tester**
* Build a Go server with a `/slow` endpoint that queries `SELECT pg_sleep(1)`.
* Set `MaxOpenConns` to 10.
* Hit the endpoint with 100 concurrent requests using `hey` or `wrk`.
* Verify via logs that the server flawlessly processes 10 at a time, queueing the rest without a single crash or dropped request.

---

## 17. Enterprise Features & Observability
* **Liveness Probes**: Kubernetes should periodically ping `db.PingContext()`. If the pool is deadlocked, the ping fails, and Kubernetes restarts the Pod.

---

## 18. Source Code Reading
Walkthrough of `database/sql/sql.go`.
* **The `connRequest` Channel**: Study the `freeConn` function. It looks at the `connRequests` map. If there's a sleeping Goroutine waiting, it sends the connection pointer directly into that Goroutine's channel, instantly waking it up!

---

## 19. Architecture
* **Decoupling Configurations**: Your `MaxOpenConns` must not be hardcoded. It should be injected via Environment Variables (`DB_MAX_OPEN_CONNS=50`) so DevOps can scale the pool up/down without recompiling the Go binary.

---

## 20. Summary & Cheat Sheet
* **Limit**: `SetMaxOpenConns(50)`
* **Idle**: `SetMaxIdleConns(10)`
* **Lifetime**: `SetConnMaxLifetime(1 * time.Hour)`
* **PgBouncer**: Use when scaling to 100+ Pods.
