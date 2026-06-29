# Rate Limiter

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
* Mini Project
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

In previous chapters, we learned how to limit the *number of concurrent executions* using Worker Pools and Semaphores. However, sometimes you need to limit the *frequency of executions over time*, regardless of how many Goroutines are running. This is called **Rate Limiting**.

Go's extended standard library provides a powerful Token Bucket rate limiter in `golang.org/x/time/rate`. It allows you to enforce rules like "allow exactly 5 requests per second, with a maximum burst of 10 requests at once."

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand the Token Bucket algorithm.
* Differentiate between concurrency limits (Semaphores) and rate limits.
* Use `time/rate.Limiter` to throttle Goroutines.
* Implement a burst-capable rate limiter for an HTTP API client.

---

# Prerequisites

Before reading this chapter you should know:

* `time.Ticker` (`18-Tickers.md`)
* Context (`19-Context.md`)
* Semaphores (`36-Semaphore.md`)

---

# Why This Topic Exists

Imagine you are building a tool that scrapes prices from Amazon. Amazon's servers have a strict rule: if you make more than 10 requests per second, they ban your IP address. 

If you use a Worker Pool of 10 Goroutines, they might execute 10 requests in 1 millisecond. The pool limits the *concurrent* connections, but the *rate* is 10,000 requests per second! You will instantly get banned. You need a Rate Limiter to guarantee your app pauses exactly 100 milliseconds between every request, perfectly pacing the traffic to exactly 10 req/sec.

---

# Real-World Analogy

### The Subway Turnstile (Token Bucket)

* **The Bucket**: The subway station manager has a bucket that holds up to 5 physical tokens (Burst limit).
* **The Refill Rate**: Every 10 seconds, the manager drops exactly 1 new token into the bucket. If the bucket is full, they throw the token away.
* **The Customers (Requests)**: To enter the subway, a customer MUST take a token from the bucket.
* **The Behavior**: 
  - If 5 people arrive at once, they can all grab a token and enter instantly (Burst).
  - If a 6th person arrive immediately, they must wait 10 seconds for the manager to drop the next token.
  - The long-term entry rate will strictly average exactly 1 person every 10 seconds (Rate Limit), but allows for sudden, small crowds (Burst).

---

# Core Concepts

* **Rate (r)**: How frequently new tokens are added to the bucket (e.g., 5 tokens per second).
* **Burst (b)**: The maximum size of the bucket. If the bucket is full, new tokens are discarded. This allows sudden spikes in traffic to be processed instantly, up to the burst limit.
* **Wait(ctx)**: A method that blocks the Goroutine until a token is available.
* **Allow()**: A non-blocking method that returns `true` if a token is available, or `false` instantly if it is not.

---

# Architecture Diagram

```mermaid
flowchart TD
    Clock[System Clock]
    Bucket[Token Bucket (Max 3)]
    Req1[Request 1]
    Req2[Request 2]
    Req3[Request 3]
    Req4[Request 4]
    
    Clock -- Adds 1 token/sec --> Bucket
    
    Req1 -- Grabs Token (Success) --> Bucket
    Req2 -- Grabs Token (Success) --> Bucket
    Req3 -- Grabs Token (Success) --> Bucket
    
    Req4 -- Grabs Token (Fails/Waits) -.-> Bucket
```

---

# Step-by-Step Implementation

1. Install the package if needed: `go get golang.org/x/time/rate`.
2. Create a limiter: `limiter := rate.NewLimiter(rate.Every(1 * time.Second), 3)`. (1 req/sec, burst of 3).
3. Inside your Goroutine, before doing the expensive work, call `limiter.Wait(context.Background())`.
4. The Goroutine will automatically sleep and wake up exactly when it is legally allowed to proceed based on the rate.

---

# Syntax

```go
import "golang.org/x/time/rate"

// rate.Limit is a float64 representing events per second.
// Burst is an int representing the max bucket size.
r := rate.Limit(5.0) // 5 per second
b := 10              // Burst of 10
limiter := rate.NewLimiter(r, b)

// Blocking wait
err := limiter.Wait(ctx)

// Non-blocking check
allowed := limiter.Allow()
```

---

# Beginner Example

A simple rate limiter pacing an infinite loop to exactly 2 executions per second.

```go
package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	// Rate: 2 tokens per second. Burst: 1.
	limiter := rate.NewLimiter(rate.Every(500*time.Millisecond), 1)
	ctx := context.Background()

	fmt.Println("Starting process...")
	start := time.Now()

	for i := 1; i <= 5; i++ {
		// Wait blocks until a token is available
		if err := limiter.Wait(ctx); err != nil {
			fmt.Println("Error waiting:", err)
			return
		}
		
		fmt.Printf("Executed task %d at %v\n", i, time.Since(start))
	}
}
```
*Output:*
```text
Executed task 1 at 0s
Executed task 2 at 500ms
Executed task 3 at 1s
Executed task 4 at 1.5s
Executed task 5 at 2s
```

---

# Intermediate Example

Demonstrating **Burst** capability. We allow 1 request per second, but a burst of 3.

```go
package main

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"time"
)

func main() {
	// Rate: 1 token per second. Burst: 3.
	// Since the bucket starts full, the first 3 requests are INSTANT.
	limiter := rate.NewLimiter(rate.Limit(1.0), 3)
	ctx := context.Background()

	start := time.Now()

	for i := 1; i <= 5; i++ {
		limiter.Wait(ctx)
		fmt.Printf("Task %d executed at %v\n", i, time.Since(start).Round(time.Millisecond))
	}
}
```
*Output:*
```text
Task 1 executed at 0ms (Burst)
Task 2 executed at 0ms (Burst)
Task 3 executed at 0ms (Burst)
Task 4 executed at 1s  (Had to wait for refill)
Task 5 executed at 2s  (Had to wait for refill)
```

---

# Advanced Example

Using `Allow()` to build an HTTP Rate Limiting Middleware. Instead of making the user wait 30 seconds for an API response, we instantly reject them if they exceed the limit.

```go
package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
)

// In a real app, you would have a map[string]*rate.Limiter (one for each IP address)
var globalLimiter = rate.NewLimiter(rate.Limit(2.0), 5) // 2 req/sec, burst 5

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow() instantly returns false if the bucket is empty
		if !globalLimiter.Allow() {
			http.Error(w, "429 Too Many Requests - Slow Down!", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, you successfully bypassed the rate limiter!")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	// Wrap our router with the rate limiter
	wrappedMux := rateLimitMiddleware(mux)

	fmt.Println("Server running on :8080...")
	http.ListenAndServe(":8080", wrappedMux)
}
```

---

# Production Use Cases

### 1. API Gateways
Every commercial API gateway (AWS API Gateway, Kong, Nginx) uses Token Bucket rate limiting. In Go, you can build a custom gateway that limits requests on a per-User or per-IP basis by keeping a `sync.Map` of `*rate.Limiter` objects, preventing any single user from DDoSing your service.

### 2. Background Sync Jobs
If you need to sync 1 million rows from a PostgreSQL database to Elasticsearch, a simple Goroutine loop will crush Elasticsearch with CPU usage. By putting a `limiter.Wait(ctx)` at 500 req/sec inside the loop, you guarantee a smooth, flat CPU utilization curve on Elasticsearch for the entire duration of the sync.

---

# Performance Analysis

* `rate.Limiter` is heavily optimized. It does not use a background Goroutine to "refill" the bucket on a timer. Instead, it uses simple math: when `Wait()` is called, it checks `time.Now()`, calculates how much time has passed since the last call, adds the mathematical number of tokens, and then sleeps the required amount.
* Because of this mathematical approach, evaluating `Allow()` or `Wait()` takes mere nanoseconds. It is safe to use in the most high-performance critical paths of your application.

---

# Best Practices

* **Understand Burst**: A Burst of 1 means requests are strictly paced (e.g., exactly 1 every second). A Burst of 100 means if the app was idle for a while, it can instantly fire 100 requests at the same exact millisecond. Ensure your downstream service can handle your Burst limit.
* **Per-Tenant Limiters**: Do not use a single global limiter for a multi-tenant web application. If User A spams the API, the global bucket empties, and User B is unfairly rate-limited. You must map limiters to IP addresses or API Keys.

---

# Common Mistakes

### Confusing Rate Limiting with Concurrency Limiting
```go
// BAD IDEA: Trying to limit API concurrency using time.Sleep
for i := 0; i < 100; i++ {
    go func() {
        time.Sleep(1 * time.Second) // DOES NOT RATE LIMIT! 100 requests will fire instantly after 1 second.
        callAPI()
    }()
}

// GOOD: Use a rate limiter
limiter := rate.NewLimiter(10, 1)
for i := 0; i < 100; i++ {
    go func() {
        limiter.Wait(ctx)
        callAPI() // Guaranteed 10 per second
    }()
}
```

---

# Debugging Guide

* **"Context Deadline Exceeded"**: If you pass a `context.WithTimeout` to `limiter.Wait(ctx)`, and the bucket is empty, and the time it would take to refill the bucket is *longer* than your context timeout, `Wait()` will instantly return an error rather than sleeping. This is a feature, not a bug!

---

# Exercises

## Beginner
Create a `rate.Limiter` that allows exactly 10 requests per second with a burst of 1. Write a `for` loop that runs 20 times, calling `Wait()`, and printing the iteration number. Time the total execution; it should take exactly 2 seconds.

## Intermediate
Write a script that launches 50 Goroutines. Each Goroutine wants to hit an API, but must wait on a shared `rate.Limiter` configured for 5 req/sec (Burst 5). Verify that the first 5 print instantly, and the rest trickle out at 5 per second.

---

# Quiz

## Multiple Choice Questions
**1. What happens if you call `limiter.Allow()` and the token bucket is empty?**
A) The Goroutine goes to sleep until a token is added.
B) The application panics.
C) It instantly returns `false`.
*Answer*: C. (`Wait()` blocks; `Allow()` is non-blocking).

## True or False
**`rate.Limiter` creates a hidden background Goroutine that ticks every millisecond to refill the bucket.**
*Answer*: False. It uses mathematical time deltas on every call, avoiding the massive CPU overhead of background tickers.

---

# Interview Questions

## Beginner
**Q**: Explain the difference between `Wait()` and `Allow()` in the Go rate limiter.
*Answer*: `Wait()` is blocking; if no tokens are available, the Goroutine sleeps until one is generated. `Allow()` is non-blocking; it instantly returns true or false, making it perfect for HTTP middleware that needs to return a 429 response immediately.

## Intermediate
**Q**: What is the purpose of the "Burst" parameter in the Token Bucket algorithm?
*Answer*: Burst defines the maximum capacity of the token bucket. It allows a system to process sudden, brief spikes in traffic instantaneously, rather than forcing strict intervals. Once the burst capacity is exhausted, the system returns to the strict refill rate.

## Advanced
**Q**: How would you implement a distributed rate limiter across 10 Kubernetes pods?
*Answer*: The `x/time/rate` package only works in a single memory space. For distributed rate limiting, you must use a centralized data store like Redis. You can implement a sliding window or token bucket in Redis using Lua scripts to ensure atomic updates across all 10 pods simultaneously.

---

# Mini Project

**Requirement**: The Smart Scraper
1. You have a slice of 30 URLs to scrape.
2. The target server allows 5 requests per second (Burst of 10).
3. Use a Worker Pool of 50 Goroutines to fetch the URLs.
4. Pass a `rate.Limiter` to every worker.
5. Inside the worker, call `limiter.Wait(ctx)` before doing the "scrape" (`fmt.Printf`).
6. Observe how the worker pool provides the raw concurrency, but the Rate Limiter paces the actual execution perfectly.

---

# Cheat Sheet

* **Import**: `golang.org/x/time/rate`
* **Init**: `limiter := rate.NewLimiter(rate.Every(time.Second), 5)` (1 req/sec, burst 5)
* **Blocking (Clients/Workers)**: `limiter.Wait(ctx)`
* **Non-Blocking (Servers/Routers)**: `if !limiter.Allow() { return HTTP 429 }`

---

# Summary

Rate limiting is the crucial layer that protects the internet from itself. By understanding the Token Bucket algorithm and using Go's mathematically optimized `rate.Limiter`, you can easily tame wildly concurrent Goroutines, ensuring your applications interact safely with strictly paced external systems.

---

# Key Takeaways

* ✔ Rate limiting restricts frequency over time, not absolute concurrency.
* ✔ Token Bucket allows strict averages with sudden bursts.
* ✔ `Wait()` blocks until ready (for background workers).
* ✔ `Allow()` rejects instantly (for HTTP handlers).

---

# Further Reading
* [Go documentation for x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)

---

# Next Chapter
➡️ **Next:** `38-Performance.md`
