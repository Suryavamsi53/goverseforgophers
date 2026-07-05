# The Context Package

The `context` package is arguably the most important package in modern Go enterprise engineering. 

If you look at the standard library, almost every single method has a variant that accepts a Context (e.g., `db.QueryContext`, `http.NewRequestWithContext`, `exec.CommandContext`).

## 1. Why Context?

In modern microservices, a single HTTP request can trigger a massive chain reaction. 
1. The User hits the API Gateway.
2. The Gateway hits the Order Service.
3. The Order Service queries the Database AND hits the Auth Service simultaneously.

**The Ghost Goroutine Problem:**
If the User loses their 5G cell signal and the connection to the API Gateway drops, the API Gateway immediately stops working. But what about the Order Service? It has no idea the user disconnected! It will continue processing the request, spinning up CPU cycles, and querying the database for a user that is no longer there. 

If 1,000 users drop their connection, you will have 1,000 "Ghost Goroutines" running in the background, consuming RAM and Database Connections. Your server will crash.

The `context.Context` is a standard mechanism to pass **cancellation signals** down this massive chain.

## 2. The Context Tree

Contexts form an immutable tree hierarchy.

You always start with the root context:
`ctx := context.Background()`

You can then derive "Child Contexts" from the root. If you cancel a Parent context, all of its Children are instantly cancelled as well. But if you cancel a Child context, the Parent is unaffected.

## 3. Passing Context

The absolute strict rule of Go is that `context.Context` should always be the **first** parameter of any function that performs I/O or blocks.

```go
// CORRECT: ctx is the first parameter
func FetchData(ctx context.Context, id int) error {
    // ...
}

// INCORRECT: Do not store Context inside a struct!
type Worker struct {
    ctx context.Context // Anti-pattern!
}
```

## 4. Context Metadata (Values)

Aside from cancellation, Contexts can also carry Request-Scoped data. 

If the API Gateway validates a JWT and extracts the `UserID` and a `TraceID` (for observability), it can inject them into the Context and pass them down the chain.

```go
// 1. Injecting data
ctx = context.WithValue(ctx, "userID", 42)

// 2. Extracting data deep in the call stack
if id, ok := ctx.Value("userID").(int); ok {
    fmt.Printf("Processing for User: %d\n", id)
}
```

**Warning:** Do NOT use `context.WithValue` to pass optional function arguments (like database connections or logger instances). It bypasses Go's strict type-checking. It should only be used for request-scoped metadata (TraceIDs, Auth Tokens).
