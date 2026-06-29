# Timeouts and Deadlines

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Real-World Analogy
* Core Concepts
* Timeouts vs Deadlines
* Architecture Diagram
* Step-by-Step Implementation (Go Context)
* Beginner Example
* Intermediate Example (Global Deadlines)
* Production Use Cases
* Best Practices
* Common Mistakes
* Debugging Guide
* Exercises
* Quiz
* Interview Questions
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

In a distributed system, a slow response is often worse than no response at all. If a service takes 30 seconds to reply, it holds database locks, consumes memory, and occupies network threads for the entire duration.

**Timeouts and Deadlines** act as a critical safety valve. They guarantee that an operation will either finish successfully or be violently aborted within a specific timeframe. In Go, this entire mechanism is elegantly managed and propagated across microservices using the `context.Context` package.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand why slow responses cause cascading failures.
* Distinguish between a Local Timeout and a Global Distributed Deadline.
* Use `context.WithTimeout` to protect Goroutines and Network I/O.
* Propagate deadlines across API boundaries using HTTP/gRPC headers.

---

# Prerequisites

Before reading this chapter you should know:

* The Go `context` package (`03-Context-Propagation.md` in Design Patterns).
* The 8 Fallacies of Distributed Computing (`02-Fallacies-of-Distributed-Computing.md`).

---

# Why This Topic Exists

Imagine your Go server can handle 1,000 concurrent requests. Normally, a request takes 0.1 seconds, so 1,000 Goroutines process very quickly and free up memory.

Suddenly, a third-party API you depend on starts taking 60 seconds to respond. 
Because you didn't configure a Timeout, your Goroutines sit there waiting. Within 1 second, all 1,000 available Goroutine slots are filled with requests waiting for the third-party API. New users trying to access completely unrelated parts of your website are blocked because the server is out of resources. 
**A slow dependency just took down your entire application.**

If you had set a 2-second timeout, the Goroutines would abort, return an error, and free up resources for other healthy parts of the system.

---

# Real-World Analogy

### The Fast Food Drive-Thru

* **No Timeout**: A customer pulls up to the window. The ice cream machine is broken. The cashier tells the customer to wait while they fix it. The customer waits 45 minutes. A line of 100 cars forms behind them, honking. The entire restaurant is paralyzed.
* **The Timeout**: The restaurant enforces a strict 2-minute "Timeout" rule. After 2 minutes, the cashier tells the customer: "I'm sorry, we cannot fulfill this right now. Please drive forward so we can serve the next car." The single order fails, but the restaurant as a whole remains healthy and operational.

---

# Core Concepts

* **Timeout**: A relative duration applied to a specific operation (e.g., "This HTTP request must finish within 2 seconds").
* **Deadline**: An absolute point in time (e.g., "This entire workflow must be completed by 14:05:00 UTC").
* **Cancellation Signal**: When a timeout expires, a signal must be broadcast to all running processes to immediately abort their work and clean up resources.

---

# Architecture Diagram

```mermaid
flowchart TD
    Client[Mobile App]
    API[API Gateway<br/>Deadline: 5s]
    Auth[Auth Service]
    DB[Database]

    Client -- "Request" --> API
    
    API -- "Calls Auth (Takes 2s)" --> Auth
    Auth -- "Returns OK" --> API
    
    note right of API: 3 Seconds Remaining on Deadline
    
    API -- "Calls DB with 3s Timeout" --> DB
    
    DB -.-x |"DB is Locked! Hangs..."| DB
    
    note left of API: 5s Expires! Context Cancelled.
    
    API -- "Interrupts DB Connection" --> DB
    API -- "Returns 504 Gateway Timeout" --> Client
```

---

# Step-by-Step Implementation (Go Context)

1. At the very top of your application (e.g., the HTTP Handler), create a context with a timeout: `ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)`.
2. **Crucial**: Always `defer cancel()` to prevent memory leaks in case the function finishes *before* the timeout.
3. Pass `ctx` as the first argument to every function, database call, and external HTTP request.
4. If the 5 seconds pass, Go automatically closes the `ctx.Done()` channel.
5. Standard libraries (like `database/sql` or `net/http`) listen to this channel and will instantly sever the TCP connection and return a `context deadline exceeded` error.

---

# Beginner Example

Protecting a slow, blocking function using `context`.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// A function that simulates heavy, slow work
func DoHeavyWork(ctx context.Context) error {
	fmt.Println("Work started...")

	// Listen for either the work finishing, or the context timing out
	select {
	case <-time.After(5 * time.Second): // The actual work takes 5 seconds
		fmt.Println("Work finished successfully!")
		return nil
	case <-ctx.Done(): // The safety valve!
		fmt.Println("ABORT: Context cancelled!")
		return ctx.Err()
	}
}

func main() {
	// Create a strict 2-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Prevent leaks

	fmt.Println("Calling function...")
	err := DoHeavyWork(ctx)

	if err != nil {
		fmt.Println("Error:", err)
	}
}
```
*Output:*
```text
Calling function...
Work started...
ABORT: Context cancelled!
Error: context deadline exceeded
```

---

# Intermediate Example (Global Distributed Deadlines)

A **Timeout** is local (e.g., an HTTP client timeout). A **Deadline** is global.
If a user requests data, and the API sets a global deadline of 5 seconds, it passes that deadline down to Service A. If Service A takes 3 seconds, it must tell Service B: "You only have 2 seconds left to finish this."

Go's gRPC implementation handles this automatically!

**Client Code:**
```go
// Create a deadline of 5 seconds from NOW
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// gRPC automatically injects this deadline into the HTTP/2 headers!
response, err := grpcClient.GetUserData(ctx, request)
```

**Server Code (Service A):**
```go
func (s *server) GetUserData(ctx context.Context, req *Request) (*Response, error) {
    // The ctx here was automatically populated from the gRPC headers!
    // It knows exactly when the client will give up.
    
    deadline, ok := ctx.Deadline()
    if ok {
        fmt.Printf("I must finish this work before %v\n", deadline)
    }
    
    // We pass the SAME context to the database.
    // If the client's 5-second deadline passes, the DB query instantly aborts!
    rows, err := db.QueryContext(ctx, "SELECT * FROM users")
}
```

---

# Production Use Cases

### 1. Database Queries (`database/sql`)
Never use `db.Query()`. Always use `db.QueryContext(ctx)`. If a DBA accidentally locks a table, `db.Query()` will hang forever, locking up your web server. `QueryContext` ensures that after a few seconds, the connection is forcefully closed and the Goroutine is freed.

### 2. Tail Latency and Hedged Requests
Google found that in a system with 10,000 servers, 1 server is always having a random CPU spike, causing a slow response (Tail Latency). To fix this, they use a pattern called **Hedged Requests**:
Send a request to Server A. Set a very aggressive timeout (e.g., 50ms). If Server A doesn't reply in 50ms, don't return an error; instead, fire the exact same request to Server B! Accept whichever response comes back first.

---

# Best Practices

* **Top-Down Deadlines**: Define the deadline at the very top of the stack (the API Gateway) and pass it all the way down to the deepest database query.
* **Always Defer Cancel**: `context.WithTimeout` starts a background timer in memory. If your function finishes successfully in 1 second, the 5-second timer is still running in memory! `defer cancel()` stops the timer immediately, preventing memory leaks.
* **Network vs Application Timeouts**: Setting a timeout on the `http.Client` only protects the network connection. Using `context` protects the actual CPU execution and Goroutines. Always prefer `context`.

---

# Common Mistakes

### Catching Timeouts but not aborting work
```go
func BadCode(ctx context.Context) {
    // We check if the context timed out...
    if ctx.Err() != nil {
        return
    }
    
    // BAD! If the context times out exactly here, this heavy work STILL RUNS!
    time.Sleep(10 * time.Second) 
}
```
*Fix*: You must actively listen to the `<-ctx.Done()` channel using a `select` block around any heavy I/O or loops, as shown in the Beginner Example. Standard library functions (like `http` and `sql`) do this internally for you.

---

# Exercises

## Beginner
Create a background function that runs an infinite `for` loop, printing "Working..." every 500ms. In `main`, use `context.WithTimeout` to allow the function to run for exactly 2 seconds, and then gracefully exit the infinite loop when the context is cancelled.

## Intermediate
Write a function that makes an HTTP GET request to a public API using `http.NewRequestWithContext`. Set the context timeout to 1 millisecond. Verify that the network request forcefully aborts and returns an error containing "context deadline exceeded".

---

# Quiz

## Multiple Choice Questions
**1. Why must you always call `defer cancel()` when using `context.WithTimeout`?**
A) Because the Go compiler requires it.
B) To prevent memory leaks by immediately stopping the background timer if the function finishes before the timeout occurs.
C) To ensure the database transaction commits successfully.
*Answer*: B

## True or False
**If you pass a Context with a 5-second timeout to an `http.Client`, and the server takes 10 seconds to process the request, the server will automatically stop processing the request at 5 seconds.**
*Answer*: False. The *Client* will abort the TCP connection and stop waiting at 5 seconds. However, unless the *Server* was specifically written to detect client disconnection via context cancellation, the Server will blissfully continue processing the request for the full 10 seconds, wasting its own CPU!

---

# Interview Questions

## Beginner
**Q**: What is the difference between a Timeout and a Deadline?
*Answer*: A timeout is a relative duration (e.g., "Wait 3 seconds"). A deadline is an absolute point in time (e.g., "Wait until 12:00 PM"). In Go, `context.WithTimeout` converts the relative duration into an absolute deadline internally.

## Intermediate
**Q**: How do slow database queries cause cascading failures in web servers?
*Answer*: Web servers use finite resources (Goroutines, threads, connection pools) to handle concurrent users. If a DB query lacks a timeout and hangs, the thread handling that user is blocked. As more users request that endpoint, all available threads become blocked waiting for the DB. The server is now completely paralyzed and cannot serve even healthy endpoints.

## Advanced
**Q**: Explain how a Global Deadline propagates across three microservices (A -> B -> C) via gRPC.
*Answer*: Service A creates a context with a 10-second timeout. It calls Service B via gRPC. The gRPC client automatically serializes this absolute deadline time into the HTTP/2 headers (e.g., `grpc-timeout: 10S`). Service B parses this header, calculates how much time is left, and reconstructs a new `context.Context` locally with the remaining time. When Service B calls Service C, the process repeats, ensuring that if 10 seconds pass globally, all three services instantly abort their execution simultaneously.

---

# Summary

In the harsh reality of distributed systems, hope is not a strategy. You cannot hope that a network call will return quickly. Timeouts and Deadlines are the mandatory seatbelts of cloud architecture, ensuring that when systems inevitably slow down, they gracefully shed load and fail fast instead of paralyzing your entire infrastructure.

---

# Key Takeaways

* ✔ A slow response is often worse than an instant error.
* ✔ Use `context.WithTimeout` to protect every piece of Network/Disk I/O.
* ✔ Always `defer cancel()` to prevent timer memory leaks.
* ✔ Global Deadlines ensure all microservices abort simultaneously when the client gives up.

---

# Further Reading
* [Go Blog: Go Concurrency Patterns: Context](https://go.dev/blog/context)

---

# Next Chapter
➡️ **Next:** `10-Distributed-Transactions.md` (Beginning of Part 4: Consistency and State)
