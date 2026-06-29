# Context

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Real-World Analogy
* Core Concepts
* Internal Runtime Explanation
* Memory Layout
* Architecture Diagram
* Step-by-Step Execution
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
* Mini Project
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

We have spent the last 18 chapters learning how to build massive, highly concurrent systems using Goroutines and Channels. But as systems grow, a critical problem emerges: **How do you cancel a massive tree of Goroutines?**

If an HTTP request spawns a handler Goroutine, which spawns 3 database query Goroutines, which spawn 5 API request Goroutines... and the user suddenly clicks "Cancel" on their browser, how do you instantly kill all 9 of those Goroutines to save CPU?

The `context` package is Go's elegant, standardized solution to this exact problem.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand the `context.Context` interface.
* Pass context down the call stack conventionally.
* Use `context.Background()` and `context.TODO()`.
* Extract request-scoped values from a Context.
* Prepare for Context Cancellation (Chapter 20).

---

# Prerequisites

Before reading this chapter you should know:

* `select` statements (`16-Select.md`)
* Interfaces in Go
* Channel Closing (Broadcast pattern) (`14-Channel-Closing.md`)

---

# Why This Topic Exists

Before Go 1.7 (when `context` was added to the standard library), developers had to pass multiple parameters down through every function: a `done` channel for cancellation, a `deadline` for timeouts, and a `map` for request-scoped data like User IDs.

Function signatures became massive and bloated:
`func fetchUser(db *DB, done chan bool, timeout time.Time, userID int, traceID string)`

The `context` package unified all of these concepts into a single, standardized interface: `ctx context.Context`. It is now the most pervasive pattern in all of Go backend engineering.

---

# Real-World Analogy

### The Corporate Badge

When you join a large corporation, you are given a Badge (the Context). 
* As you walk through the building, you must show your badge at every door (passing context to every function).
* The badge contains metadata about you (Request-scoped values like `Role=Admin`).
* If HR suddenly invalidates your badge (Cancellation), every door instantly stops letting you through, and security escorts you out (Goroutines terminate).

---

# Core Concepts

* **`context.Context`**: An immutable interface. It is never modified, only wrapped.
* **`context.Background()`**: The empty root context. It is never canceled and has no values. Used at the very top level of main or tests.
* **`context.TODO()`**: Identical to Background, but signals to other developers: "I don't know what context to use here yet, I'll fix this later."
* **Request-Scoped Values**: Contexts can carry small pieces of metadata (like Trace IDs or Auth Tokens) down the call stack.

---

# Internal Runtime Explanation

The `Context` is just an interface with 4 methods:
1. `Deadline() (deadline time.Time, ok bool)`
2. `Done() <-chan struct{}`
3. `Err() error`
4. `Value(key any) any`

When you wrap a context (e.g., using `context.WithValue`), it creates a new struct that points to the parent context, forming a linked list. When you call `Value()`, it traverses up the linked list to find the key, much like prototype inheritance in JavaScript.

---

# Memory Layout

```text
Heap Memory (Linked List of Contexts)

+-----------------------+
| *emptyCtx (Background)|
+-----------------------+
           ^
           | (parent)
+-----------------------+
| *valueCtx             |  -> Key: "TraceID", Value: "12345"
+-----------------------+
           ^
           | (parent)
+-----------------------+
| *valueCtx             |  -> Key: "UserID", Value: "99"
+-----------------------+
```

---

# Architecture Diagram

```mermaid
flowchart TD
    Main[Main HTTP Handler] -->|ctx| DB[Database Fetcher]
    Main -->|ctx| API[External API Caller]
    Main -->|ctx| Cache[Redis Cache Checker]
    
    API -->|ctx| Logger[Log Writer]
    
    style Main fill:#f9f,stroke:#333,stroke-width:2px
    Note right of Main: Context originates here
```

---

# Step-by-Step Execution

1. An HTTP request hits the server. The router creates a Root Context.
2. The router injects the `TraceID` into the Context (wrapping it).
3. The router passes the `ctx` as the first parameter to the `HandleRequest(ctx, w, r)` function.
4. `HandleRequest` passes `ctx` into `DB.FetchUser(ctx)`.
5. `FetchUser` extracts the `TraceID` from the `ctx` to include in the SQL logs.

---

# Syntax

```go
import "context"

// 1. Creating the root
ctx := context.Background()

// 2. Injecting a value (Returns a NEW context!)
// Note: Keys should ideally be custom unexported types, not raw strings.
type contextKey string
const userIDKey contextKey = "userID"

ctxWithValue := context.WithValue(ctx, userIDKey, 42)

// 3. Extracting a value
val := ctxWithValue.Value(userIDKey)
```

---

# Beginner Example

Passing a Trace ID down the call stack so all logs share the same ID.

```go
package main

import (
	"context"
	"fmt"
)

type key string
const traceIDKey key = "traceID"

func queryDatabase(ctx context.Context) {
	// Extract the value from the context
	traceID := ctx.Value(traceIDKey)
	fmt.Printf("[Trace: %s] Executing SQL Query...\n", traceID)
}

func handleRequest(ctx context.Context) {
	fmt.Println("Handling incoming request...")
	queryDatabase(ctx)
}

func main() {
	// 1. Create a root context
	rootCtx := context.Background()

	// 2. Wrap it with a trace ID
	ctx := context.WithValue(rootCtx, traceIDKey, "REQ-998877")

	// 3. Pass it down the chain
	handleRequest(ctx)
}
```

---

# Intermediate Example

Understanding why we use custom types for Context Keys to prevent collision.

```go
package main

import (
	"context"
	"fmt"
)

// By defining a custom unexported type, it is mathematically impossible
// for another package to accidentally overwrite our key!
type authKeyType string
const AuthKey authKeyType = "auth"

// Imagine this is in a 3rd party package
type thirdPartyKey string
const OtherAuthKey thirdPartyKey = "auth"

func main() {
	ctx := context.Background()
	
	// We set our auth token
	ctx = context.WithValue(ctx, AuthKey, "my-secret-jwt")
	
	// A 3rd party library sets their auth token using the EXACT same string "auth"
	ctx = context.WithValue(ctx, OtherAuthKey, "their-secret")

	// Because the TYPES are different, they do not collide!
	fmt.Println("My Auth:", ctx.Value(AuthKey)) 
	fmt.Println("Their Auth:", ctx.Value(OtherAuthKey))
}
```

---

# Advanced Example

Using Context Values is controversial. The Go authors recommend using them *only* for request-scoped data that transits processes and APIs, not for passing optional parameters to functions.

```go
// GOOD Use of Context Values:
// - Trace IDs / Request IDs
// - Authentication Tokens / User IDs
// - IP Addresses of the caller

// BAD Use of Context Values (Anti-Patterns):
// - Database Connections (Pass this as a struct field!)
// - Logger Instances (Pass this as a struct field!)
// - Configuration structs
```

---

# Production Use Cases

### 1. Request Tracing (OpenTelemetry)
In microservice architectures, when Service A calls Service B, it generates a `TraceID`. This `TraceID` is shoved into the `context.Context` and passed to every single function. When Service B calls PostgreSQL, the Postgres driver extracts the `TraceID` from the context and attaches it to the database logs.

### 2. HTTP Middlewares
In standard Go `net/http`, every `*http.Request` has a built-in `r.Context()`. Middlewares often extract JWTs from the Authorization header, validate them, and shove the UserID into `r.WithContext(ctx)` so the final handler knows who the user is.

---

# Performance Analysis

* **Linked List Traversal**: `ctx.Value()` is an O(N) operation based on how many times the context was wrapped. If you wrap a context 1,000 times and ask for the very first key, it must traverse 1,000 pointers. It is relatively slow compared to a standard Map lookup. Do not abuse it!

---

# Best Practices

* **Always the first parameter**: By strict Go convention, `ctx context.Context` should always be the very first parameter of any function that does I/O or takes a long time.
* **Do not store Contexts in Structs**: Contexts are ephemeral and request-scoped. They should flow through function arguments, not sit permanently inside a struct field (with very rare exceptions).
* **Use custom types for keys**: Never use `context.WithValue(ctx, "userID", 123)`. Always define a type `type key string` to prevent collisions.

---

# Common Mistakes

### Using Context to pass dependencies
```go
// TERRIBLE ANTI-PATTERN:
func GetUser(ctx context.Context, id int) {
    // Extracting a database connection from a context is considered horrible practice.
    // You lose compile-time type safety!
    db := ctx.Value("DB").(*sql.DB) 
    db.Query(...)
}

// CORRECT:
type Repository struct {
    DB *sql.DB // Inject dependencies into structs!
}
func (r *Repository) GetUser(ctx context.Context, id int) { ... }
```

---

# Debugging Guide

* **Nil Context Panics**: Never pass a `nil` context to a function that expects a `context.Context`. It will likely panic when they call `ctx.Done()`. If you don't have a context, pass `context.TODO()`.

---

# Exercises

## Beginner
Write a function `printContext(ctx context.Context)`. In `main`, create a `Background` context, add a key `language` with value `"Go"`, and pass it to the function to print.

## Intermediate
Examine the standard library `net/http` package documentation. Find the method `r.Context()`. How would you use this in an HTTP handler to get a context?

---

# Quiz

## Multiple Choice Questions
**1. Which context should you use at the very top of your `main()` function?**
A) `context.TODO()`
B) `context.Background()`
C) `context.New()`
*Answer*: B

## True or False
**It is a best practice to pass your Database Connection Pool through `context.WithValue`.**
*Answer*: False! Dependencies should be injected via struct fields or direct parameters. Context values are strictly for request-scoped metadata like Trace IDs.

---

# Interview Questions

## Beginner
**Q**: What is the difference between `context.Background()` and `context.TODO()`?
*Answer*: Functionally, they are identical (both return an empty root context). Semantically, `Background` is used when you deliberately want a root context (like in `main`), while `TODO` is used when you are unsure or refactoring old code to support context.

## Intermediate
**Q**: Why must context keys be custom types rather than raw strings?
*Answer*: Because multiple packages might use the same string key (e.g., "auth"). If they do, they will collide in the context linked list, and one package might extract the other package's data, causing severe security bugs. Custom types prevent this collision at the compiler level.

## Google-Level Questions
**Q**: Explain the performance characteristics of `context.WithValue`. Is it suitable for storing a configuration map of 100 items?
*Answer*: No. `context.WithValue` creates a new Context node that points to its parent, forming a singly linked list. Calling `Value()` forces a linear O(N) traversal up the tree. Storing 100 items would create a deep tree that is incredibly slow to query. If you need 100 configuration items, you should store a *pointer* to a single Configuration struct as a single Context value, making it an O(1) lookup.

---

# Mini Project

**Requirement**: The Middleware Mock
Create an HTTP server with a single route `/dashboard`. Write a middleware function that checks if the request has an `X-User-Role` header. If it does, shove that role into the `r.Context()` using `context.WithValue`. In the final `/dashboard` handler, extract the role from the context. If the role is "Admin", print "Welcome Admin", otherwise print "Access Denied".

---

# Cheat Sheet

* **Root**: `ctx := context.Background()`
* **Placeholder**: `ctx := context.TODO()`
* **Set Value**: `ctx = context.WithValue(ctx, keyType("id"), 42)`
* **Get Value**: `val := ctx.Value(keyType("id")).(int)`

---

# Summary

The `context` package is the nervous system of modern Go applications. It standardizes how metadata flows down the call stack. While passing values is useful for tracing, the true superpower of Context is Cancellation, which we will explore in the very next chapter.

---

# Key Takeaways

* ✔ Always pass `ctx` as the first parameter.
* ✔ Use `context.Background()` at the root.
* ✔ Use custom types for keys to prevent collisions.
* ✔ Never store contexts in structs.

---

# Further Reading
* [Go Blog: Context](https://go.dev/blog/context)

---

# Next Chapter
➡️ **Next:** `20-Cancellation.md`
