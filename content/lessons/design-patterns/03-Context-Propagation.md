# Context Propagation

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Real-World Analogy
* Core Concepts
* Architecture Diagram
* Step-by-Step Implementation
* Syntax
* Beginner Example
* Intermediate Example
* Advanced Example
* Production Use Cases
* Performance Analysis
* Best Practices
* Common Mistakes
* Debugging Guide
* Exercises
* Quiz
* Interview Questions
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

In Go, the `context` package is arguably the most pervasive design pattern across the entire ecosystem. Almost every standard library package and third-party library that deals with I/O (network, database, file system) expects a `context.Context` as the first argument of its functions.

**Context Propagation** is the design pattern of passing a `Context` explicitly down the call stack, from the very beginning of a request (like an incoming HTTP call) all the way down to the deepest database queries, allowing you to seamlessly manage timeouts, cancellations, and request-scoped metadata across boundaries.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand why Context is explicitly passed in Go, rather than implicitly stored.
* Use `context.WithTimeout` and `context.WithCancel` to protect resources.
* Pass request-scoped values (like User IDs or Tracing IDs) safely through contexts.
* Avoid the common pitfalls of context abuse.

---

# Prerequisites

Before reading this chapter you should know:

* Basic Go functions.
* Goroutines (`08-Goroutines.md`).
* A basic understanding of what `context` is (`19-Context.md`).

---

# Why This Topic Exists

Imagine a user loads a webpage that fetches data from 5 different microservices. 
After 1 second, the user gets impatient, closes the browser tab, and leaves.

If your Go server doesn't use Context Propagation, those 5 microservices will continue executing their heavy database queries for the next 10 seconds, wasting CPU and memory on a response that no one is listening to anymore.

By passing a `Context` from the HTTP handler down to every database call, as soon as the user disconnects, the context is cancelled. This cancellation signal instantly propagates down the entire call stack, instantly aborting all 5 database queries and freeing up your server's resources.

---

# Real-World Analogy

### The Construction Project Manager

* **The Request**: A client calls the Project Manager and orders a skyscraper.
* **The Context**: The Project Manager gives a specific walkie-talkie channel to the foremen, the plumbers, and the electricians. 
* **Propagation**: The foremen hand out walkie-talkies (on the same channel) to their individual workers.
* **Cancellation**: If the client suddenly goes bankrupt and cancels the project, the Project Manager presses the button on the walkie-talkie. Instantly, the foremen, plumbers, electricians, and every single worker hears the "ABORT" signal and drops their tools. No further money (CPU) is wasted.

---

# Core Concepts

* **Explicit Passing**: Go does not have "Thread Local Storage" (unlike Java/Python). You *must* physically pass the context as the first argument to every function that needs it: `func doWork(ctx context.Context, data string)`.
* **Immutability**: A `Context` is immutable. When you add a timeout or a value to it, you don't mutate the original; you create a *new* wrapped context that inherits from the parent.
* **Done Channel**: A channel `<-ctx.Done()` that closes when the context is cancelled or times out.

---

# Architecture Diagram

```mermaid
flowchart TD
    HTTP[HTTP Handler<br/>Creates Context]
    Auth[Auth Middleware]
    Service[Business Logic Service]
    DB[Database Client]
    HTTPAPI[External API Client]
    
    HTTP -- Passes ctx --> Auth
    Auth -- Adds UserID to ctx --> Service
    Service -- Passes ctx --> DB
    Service -- Passes ctx --> HTTPAPI
    
    note right of HTTP: If HTTP Request is aborted,<br/>Context cancels.<br/>DB and HTTPAPI stop instantly.
```

---

# Step-by-Step Implementation

1. **Start the Context**: At the entry point of your app (e.g., inside an HTTP handler), extract the context from the request: `ctx := r.Context()`. (For background jobs, use `context.Background()`).
2. **Propagate**: Modify all functions in your call stack to accept `ctx context.Context` as their *first* parameter.
3. **Wrap (Optional)**: If a specific sub-task needs a stricter timeout, wrap the context: `ctx, cancel := context.WithTimeout(ctx, 2*time.Second)`.
4. **Consume**: At the very bottom of the stack, pass the context into the standard library (e.g., `http.NewRequestWithContext(ctx, ...)` or `sql.QueryContext(ctx, ...)`).

---

# Syntax

```go
// The standard pattern: ctx is ALWAYS the first argument
func FetchUser(ctx context.Context, id int) (*User, error) {
    // Pass it down to the next layer
    return db.QueryUser(ctx, id) 
}
```

---

# Beginner Example

Propagating a timeout down to a simulated database call.

```go
package main

import (
	"context"
	"fmt"
	"time"
)

// The lowest layer (e.g., a database query)
func queryDatabase(ctx context.Context) error {
	fmt.Println("DB query started...")
	
	// Simulate a query that takes 3 seconds
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("DB query finished successfully!")
		return nil
	case <-ctx.Done(): // Listen for the cancellation signal
		fmt.Println("DB query aborted!")
		return ctx.Err()
	}
}

// The middle layer (e.g., a service)
func getUserService(ctx context.Context) error {
	// Simply pass the context down
	return queryDatabase(ctx)
}

func main() {
	// Top layer: Create a context with a 1-second timeout
	// (But the DB query takes 3 seconds!)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel() // Always defer cancel to prevent memory leaks

	fmt.Println("Calling service...")
	err := getUserService(ctx)
	
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```
*Output: The query aborts after 1 second with `context deadline exceeded`.*

---

# Intermediate Example

Propagating Request-Scoped Values (like a Trace ID or User ID).

```go
package main

import (
	"context"
	"fmt"
)

// Define a custom unexported type for context keys to prevent collisions
type contextKey string
const userIDKey contextKey = "userID"

// Middleware-like function that adds a value to the context
func authMiddleware(ctx context.Context, token string) context.Context {
	// In reality, you'd decode the token here
	userID := "user_123"
	
	// Create a NEW context wrapping the old one, containing the value
	return context.WithValue(ctx, userIDKey, userID)
}

// Deep service function that extracts the value
func updateProfile(ctx context.Context) {
	// Extract the value. Must type-assert it.
	userID := ctx.Value(userIDKey)
	
	if userID == nil {
		fmt.Println("Error: No user ID found in context!")
		return
	}
	
	fmt.Printf("Updating profile for user: %s\n", userID.(string))
}

func main() {
	ctx := context.Background()
	
	// Request starts, passes through auth
	ctxWithUser := authMiddleware(ctx, "fake-jwt-token")
	
	// Pass the context deeply into the app
	updateProfile(ctxWithUser)
}
```

---

# Advanced Example

Combining Timeouts and Values in a real HTTP server. Notice how the HTTP request automatically provides a context that cancels if the client disconnects!

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func heavyQuery(ctx context.Context) error {
	// A fake query that takes 5 seconds
	select {
	case <-time.After(5 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err() // e.g., context canceled
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the Context from the incoming HTTP Request.
	// If the user closes their browser, this context cancels instantly!
	ctx := r.Context()
	
	// 2. We ALSO want to enforce a strict 2-second timeout on the DB.
	// We wrap the HTTP context. Now it cancels if the user disconnects OR 2 seconds pass.
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	fmt.Println("Starting heavy query...")
	err := heavyQuery(ctx)
	
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
		return
	}
	
	fmt.Fprintln(w, "Query successful!")
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server running on :8080 (Try hitting it and pressing Ctrl+C in your browser fast!)")
	http.ListenAndServe(":8080", nil)
}
```

---

# Production Use Cases

### 1. Distributed Tracing (OpenTelemetry)
In microservice architectures, when Service A calls Service B, it needs to pass a "Trace ID" so logs can be aggregated. The Trace ID is stored in the `context.Context` using `context.WithValue`. Service A reads the ID from the Context, injects it into HTTP Headers, sends it to Service B, which extracts it and puts it into a new Context.

### 2. Database Transaction Scoping
Some ORMs and database libraries allow you to start an SQL Transaction and store the transaction object inside a `context.Context`. When you call repository functions down the stack, they extract the transaction from the Context and use it, ensuring all queries run within the same ACID transaction.

---

# Performance Analysis

* Context creation and cancellation is extremely lightweight. `WithCancel` spawns a tiny amount of bookkeeping to link the parent to the child, but it is highly optimized.
* `context.WithValue` is slightly slower because it operates like a linked list. If you look up a key, it checks the current context, then the parent, then the grandparent. **Do not use Context as a general-purpose map for massive amounts of data.**

---

# Best Practices

* **Always the First Argument**: `ctx` should always be the very first parameter of a function: `func Do(ctx context.Context, arg1 string)`.
* **Never store Context in a Struct**: Contexts are request-scoped. Structs are usually long-lived. If you store a context in a struct field, it will eventually expire or be cancelled, permanently breaking the struct.
* **Custom Types for Keys**: When using `context.WithValue`, NEVER use basic types like `string` for the key. Always define a custom unexported type `type contextKey string` to prevent key collisions between different packages.

---

# Common Mistakes

### Passing `nil` Contexts
```go
// BAD: Will cause a panic if passed to standard library functions!
queryDatabase(nil)

// GOOD: If you don't have a context yet, or are writing a test, use TODO
queryDatabase(context.TODO())
// Or use Background for the top-level of a program
queryDatabase(context.Background())
```

### Forgetting to Cancel
```go
// BAD: Memory leak! The timer created by WithTimeout stays active in the background!
ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
doWork(ctx)

// GOOD: Always defer cancel immediately
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
doWork(ctx)
```

---

# Debugging Guide

* **"context deadline exceeded"**: This error means the operation took longer than the `WithTimeout` or `WithDeadline` allowed. Increase the timeout, or optimize the underlying query.
* **"context canceled"**: This usually means the parent operation was aborted. In HTTP servers, it almost always means the client disconnected (closed their browser/terminal) before the server could finish responding.

---

# Exercises

## Beginner
Write a function `doTask(ctx context.Context)` that simulates a 3-second task using a `select` and `time.After`. In `main`, call it with a 1-second timeout context. Print the resulting error.

## Intermediate
Create a context and add a value `userID = 55`. Create a child context from it with a 2-second timeout. Pass the child context into a function. Prove that the function can still read the `userID` from the wrapped child context!

---

# Quiz

## Multiple Choice Questions
**1. Why shouldn't you store a `context.Context` inside a long-lived struct field?**
A) Because it uses too much memory.
B) Because Contexts are tied to the lifecycle of a single request/operation. If the context is canceled, the struct becomes permanently unusable.
C) Because Go's garbage collector cannot clean up structs with contexts.
*Answer*: B

## True or False
**It is safe to use `string` types for keys in `context.WithValue` as long as you make the string unique, like `"my-app-user-id"`.**
*Answer*: False. While it might work, it is strongly discouraged by the Go team. Another package might accidentally use the exact same string, causing a catastrophic data collision. Always use an unexported custom type.

---

# Interview Questions

## Beginner
**Q**: What is the primary purpose of the `context` package in Go?
*Answer*: To pass cancellation signals, deadlines/timeouts, and request-scoped values across API boundaries and between Goroutines.

## Intermediate
**Q**: Explain why you must always call the `cancel` function returned by `context.WithTimeout`, even if the operation finishes successfully before the timeout.
*Answer*: `context.WithTimeout` starts a timer in the background. If the operation finishes quickly and you don't call `cancel`, that timer stays in memory until it expires, creating a memory leak. Calling `cancel` explicitly stops the timer and frees the resources immediately.

## Advanced
**Q**: If you have a parent context with a 5-second timeout, and you wrap it with a child context that has a 10-second timeout, what happens?
*Answer*: The child context will cancel after 5 seconds. Context deadlines can only be *shortened*, never lengthened. A child context is always canceled when its parent is canceled. The 5-second parent timeout acts as a hard upper bound.

---

# Cheat Sheet

* **Create Root Context**: `ctx := context.Background()`
* **Create Test/Placeholder Context**: `ctx := context.TODO()`
* **Timeout**: `ctx, cancel := context.WithTimeout(parent, 5*time.Second); defer cancel()`
* **Values (Write)**: `ctx = context.WithValue(ctx, myCustomKey("id"), 123)`
* **Values (Read)**: `val := ctx.Value(myCustomKey("id")).(int)`
* **Listen for Cancel**: `<-ctx.Done()`

---

# Summary

Context Propagation is the connective tissue of Go programs. It is the reason Go servers are incredibly resilient under load. By explicitly passing a Context down the stack, you ensure that no Goroutine ever works on stale or aborted requests, keeping your system lean, responsive, and robust.

---

# Key Takeaways

* ✔ Pass `ctx context.Context` as the first argument to functions.
* ✔ Use Context for Timeouts, Cancellations, and Request-Scoped Data.
* ✔ Never store Contexts in structs.
* ✔ Always `defer cancel()` when creating derived contexts.

---

# Further Reading
* [Go Blog: Go Concurrency Patterns: Context](https://go.dev/blog/context)

---

# Next Chapter
➡️ **Next:** `04-Error-Wrapping-and-Typing.md`
