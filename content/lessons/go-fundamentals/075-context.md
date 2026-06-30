# Context in Go

## 1️⃣ Learning Objectives
* **What you'll learn**: Master `context.Context` for request scoped data, cancellation signals, and deadlines.
* **Why it matters**: Without `Context`, API requests would hang forever if a downstream database or microservice failed to respond. It is the definitive tool to prevent goroutine leaks.
* **Where it's used**: Passed universally as the very first argument to almost every function in production Go code (HTTP handlers, gRPC, database queries).

---

## 2️⃣ Real-world Story
Imagine calling customer support (an API request) to cancel a subscription. The representative puts you on hold to call the billing department (a database query). 

If you get tired of waiting and **hang up the phone** (Client Disconnect), the representative shouldn't stay on hold with the billing department forever! They should immediately hang up their end too. 

`Context` is the mechanism that tells the entire chain of backend workers: *"The original user hung up. Stop what you are doing, drop the database connection, and kill the goroutines right now."*

---

## 3️⃣ Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[HTTP Request] -->|Creates ctx| B(Handler)
    B -->|Passes ctx| C(Database Query)
    B -->|Passes ctx| D(External API Call)
    
    E[Client Disconnects] -.->|Triggers Cancel()| B
    B -.->|Propagates Cancel()| C
    B -.->|Propagates Cancel()| D
    
    C -->|Closes DB Conn| F[Goroutine Exits]
    D -->|Drops HTTP Conn| G[Goroutine Exits]
```

---

## 4️⃣ Internal Working (Under the Hood)
Deep dive into `src/context/context.go`.
`Context` is an interface with four methods:
```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```
* **Done()**: The magic mechanism. It returns an unbuffered channel `<-chan struct{}`. When a context is cancelled, this channel is closed. Any goroutine `select`-ing on this channel will instantly wake up!
* **Context Tree**: Contexts form an immutable linked list. A parent `cancel()` call iterates through all children contexts and closes their `Done` channels sequentially.

---

## 5️⃣ Compiler Behavior
* **Escape Analysis**: Since Contexts are heavily nested linked lists of interfaces (like `cancelCtx` wrapping `emptyCtx`), passing `context.WithValue` almost universally forces heap allocations. Avoid storing large objects in it!

---

## 6️⃣ Memory Management
* **Memory Leaks**: If you create a `ctx, cancel := context.WithCancel(parent)`, and you FORGET to call `cancel()`, the child context stays permanently attached to the parent's linked list. This is a severe memory leak. Always `defer cancel()` immediately!

---

## 7️⃣ Code Examples

### 🔹 Example 1: Simple (Timeout)
```go
func fetchWithTimeout() {
    // Fails automatically if it takes > 2 seconds
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel() // MUST DEFER!

    req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.github.com", nil)
    http.DefaultClient.Do(req)
}
```

### 🔹 Example 2: Intermediate (Cancellation from Client)
```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    // r.Context() automatically cancels if the user closes their browser
    err := slowDatabaseQuery(r.Context())
    if err != nil {
        fmt.Println("User left early, DB query aborted!")
    }
}
```

### 🔹 Example 3: Advanced (Context Values)
```go
type key string
const userIDKey key = "userID"

// Middleware injects value
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), userIDKey, "user-123")
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Handler extracts value
func HandleProfile(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value(userIDKey).(string)
    fmt.Println(userID)
}
```

### 🔹 Example 4: Production (Graceful Shutdown)
```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    // Pass ctx to your HTTP server or worker pool.
    // When the user hits Ctrl+C, ctx.Done() closes!
    <-ctx.Done()
    fmt.Println("Shutting down cleanly...")
}
```

---

## 8️⃣ Production Examples
1. **Database Connection Pools**: `sql.DB` uses context to terminate long-running `SELECT` queries directly at the Postgres network driver level.
2. **Distributed Tracing**: OpenTelemetry heavily relies on `context.WithValue` to invisibly pass `TraceID` and `SpanID` down the call stack.

---

## 9️⃣ Performance & Benchmarking
* **CPU vs Memory Trade-offs**: Every time you call `context.WithValue`, it allocates a new `valueCtx` struct on the heap and creates a linked list node. Lookups are `O(N)` linear scans up the tree! 
**Never use Context as a generic dictionary.** Only use it for truly request-scoped variables (like Trace IDs).

---

## 🔟 Best Practices
* ✅ **Do**: Pass `ctx context.Context` as the very first parameter of a function.
* ✅ **Do**: Call it `ctx` (don't name it `c` or `context`).
* ❌ **Don't**: Store `Context` inside a struct (e.g., `type Service struct { ctx Context }`). Context must flow through function arguments!
* 🏢 **Google Style**: Do not use `context.Background()` deep inside business logic. Always pass the parent context down. Use `context.TODO()` if you are unsure what context to pass while refactoring.

---

## 11️⃣ Common Mistakes
1. **Forgetting `defer cancel()`**: Causes the parent context to hold onto the child forever in memory.
2. **Key Collisions in `WithValue`**: Using a basic `string` as a key.
```go
// BAD
ctx = context.WithValue(ctx, "user", 123) 

// GOOD: Define an unexported custom type to prevent collision across packages!
type contextKey int
const userKey contextKey = 0
ctx = context.WithValue(ctx, userKey, 123)
```

---

## 12️⃣ Debugging
* **Hanging Goroutines**: If you notice goroutines aren't dying (via `pprof`), look for `select` statements that forgot to include a `case <-ctx.Done():` channel read.

---

## 13️⃣ Exercises
1. **Easy**: Write a function that loops infinitely. Make it stop gracefully when a 3-second `WithTimeout` context expires.
2. **Medium**: Set up an HTTP handler, run `time.Sleep(10s)`, but detect if the client cancels the request mid-sleep using `r.Context().Done()`.
3. **Hard**: Build a custom nested Context hierarchy and prove that cancelling the parent cancels all children, but cancelling a child does NOT cancel the parent.

---

## 14️⃣ Quiz
1. **MCQ**: What is the Time Complexity of looking up a value in `context.WithValue` if the context has been nested 50 times?
   - A) O(1)
   - B) O(log N)
   - C) O(N)
*(Answer: C. It recursively walks up the linked list one node at a time!)*

---

## 15️⃣ FAANG Interview Questions
* **Beginner**: Why shouldn't you put database connections or large JSON payloads inside a Context?
* **Intermediate**: Explain the difference between `context.Background()` and `context.TODO()`.
* **Senior (Google/Meta)**: Design a circuit breaker pattern using Context that automatically cancels all pending downstream microservice requests if the error rate exceeds 50%.

---

## 16️⃣ Mini Project
**Context-Aware Web Scraper**
Build a concurrent web scraper. It receives a list of 100 URLs to scrape.
* Use `context.WithTimeout` to give the entire scraping job a hard limit of 5 seconds.
* Pass this context to every `http.NewRequestWithContext`.
* If 5 seconds hit, ensure all 100 goroutines instantly abort their active HTTP connections and exit without leaking memory.

---

## 17️⃣ Enterprise Features & Observability
* **Logging**: Inject the `TraceID` from the context into every single log line your application produces using `slog` or `logrus` hooks.

---

## 18️⃣ Source Code Reading
Read `src/context/context.go`.
* Observe the `cancelCtx` struct. See how it protects its children map using a `sync.Mutex` because `cancel()` can be called asynchronously by multiple goroutines at once!

---

## 19️⃣ Architecture
Context is the singular exception to the "Clean Architecture" rule of hiding framework details. It is universally accepted to pass `context.Context` deep into your innermost Domain and Repository layers.

---

## 20️⃣ Summary & Cheat Sheet
* **Timeout**: `ctx, cancel := context.WithTimeout(parent, time.Second)`
* **Cancel**: `ctx, cancel := context.WithCancel(parent)`
* **Values**: `ctx = context.WithValue(ctx, key, val)`
* **Check Cancellation**: 
```go
select {
case <-ctx.Done():
    return ctx.Err() // returns context.Canceled or context.DeadlineExceeded
}
```
