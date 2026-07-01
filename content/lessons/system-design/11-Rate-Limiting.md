# Rate Limiting: Protecting Your Systems from Abuse

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Rate Limiting Algorithms
* Distributed Rate Limiting
* Code Examples & Good Principles
* Architecture Diagram
* Real-World Analogy
* Interview Questions
* Quiz
* Exercises
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

In a perfect world, clients would make well-behaved requests to your API. In reality, your systems will face malicious scrapers, brute-force attacks, buggy scripts, and sudden viral traffic spikes. 

Rate Limiting is a defensive mechanism that restricts the number of requests a client (identified by IP, User ID, or API Key) can make within a specified time window. It prevents resource starvation and ensures fair usage across all users.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand why rate limiting is essential for system reliability and cost control.
* Explain the pros and cons of different rate-limiting algorithms.
* Implement a Token Bucket algorithm in Go.
* Design a distributed rate limiter using Redis to handle traffic across multiple server instances.

---

# Prerequisites

Before reading this chapter you should know:

* Microservices & API Gateways (`10-Microservices.md`)
* Caching / Redis (`06-Caching.md`)

---

# Why This Topic Exists

Imagine an attacker tries to guess a user's password by sending 10,000 login requests per second. Or a premium API user runs a buggy infinite loop that hammers your database. Without rate limiting, these scenarios can cause a Denial of Service (DoS) for everyone else. Rate limiting is often one of the first lines of defense configured at the API Gateway or Load Balancer level.

---

# Rate Limiting Algorithms

There are several standard algorithms, each with different performance and strictness trade-offs.

### 1. Token Bucket
* **How it works**: Imagine a bucket that holds a maximum of *N* tokens. Tokens are added to the bucket at a fixed rate (e.g., 10 tokens per second). When a request arrives, it must take a token from the bucket. If the bucket is empty, the request is dropped (HTTP 429 Too Many Requests).
* **Pros**: Easy to implement. Allows for sudden bursts of traffic (up to the bucket capacity).
* **Cons**: If the bucket is full, a massive sudden burst might still briefly overwhelm downstream services.

### 2. Leaky Bucket
* **How it works**: Imagine a bucket with a hole in the bottom. Requests are poured into the top of the bucket. The bucket leaks (processes) requests at a strict, constant rate. If the bucket overflows, new requests are discarded.
* **Pros**: Smooths out traffic into a steady stream.
* **Cons**: Sudden bursts of valid traffic are dropped or delayed if the bucket is full.

### 3. Fixed Window Counter
* **How it works**: The timeline is divided into fixed windows (e.g., 12:00:00 to 12:01:00). A counter increments for each request in that window. If the limit is reached, requests are dropped until the next window starts.
* **Pros**: Very memory efficient.
* **Cons**: The "Boundary Problem". If the limit is 100/minute, a user could send 100 requests at 12:00:59, and another 100 requests at 12:01:01. The system just processed 200 requests in 2 seconds, effectively bypassing the intended limit.

### 4. Sliding Window Log / Counter
* **How it works**: Solves the boundary problem by keeping a precise timestamp (log) of every request, or mathematically smoothing fixed windows based on the exact current time.
* **Pros**: Highly accurate. No boundary problems.
* **Cons**: Can consume significant memory if tracking individual timestamps (Sliding Window Log).

---

# Distributed Rate Limiting

If you have one Go server, you can store the rate limit counters in local memory (RAM). But if you have 10 Go servers behind a load balancer, a user could theoretically make 10x the allowed requests by hitting different servers.

**Solution: Centralized Storage (Redis)**
To enforce a global rate limit, the API Gateway or Go servers must store their counters in a fast, centralized datastore like Redis. However, doing a `GET` and `SET` for every request introduces race conditions. We solve this using **Redis Lua Scripts** to ensure the check-and-decrement operation is atomic.

---

# Code Examples & Good Principles

### Principle: Token Bucket in Go (In-Memory for a single node)

The `golang.org/x/time/rate` package provides a robust, production-ready Token Bucket implementation.

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// Principle: Isolate limiters per client (e.g., per IP address).
// For simplicity in this example, we use a single global limiter.
// Allow 2 requests per second, with a maximum burst of 5.
var limiter = rate.NewLimiter(rate.Limit(2), 5)

func rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ask the limiter if we are allowed to proceed
		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1") // Be a good citizen, tell them when to retry
			http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Request Successful!")
	})

	// Wrap the mux with our middleware
	fmt.Println("Server running on :8080 (Try hitting it rapidly)")
	http.ListenAndServe(":8080", rateLimiterMiddleware(mux))
}
```

---

# Architecture Diagram

```mermaid
flowchart TD
    Client([Client IP: 192.168.1.1])
    
    subgraph API_Gateway_Cluster
        LB[Load Balancer]
        AG1[Gateway Node 1]
        AG2[Gateway Node 2]
    end
    
    subgraph Central_Cache
        Redis[(Redis Cache \nAtomic Lua Scripts)]
    end
    
    subgraph Backend
        Svc[Protected Microservice]
    end

    Client -- "1. Request" --> LB
    LB -- "2. Routes to" --> AG1
    
    AG1 -- "3. Check & Decrement Limit for IP" --> Redis
    
    alt Limit Exceeded
        Redis -. "4a. Return Denied" .-> AG1
        AG1 -. "HTTP 429" .-> Client
    else Limit Allowed
        Redis -. "4b. Return Allowed" .-> AG1
        AG1 -- "5. Forward Request" --> Svc
    end
```

---

# Real-World Analogy

* **Token Bucket**: A bouncer at a nightclub. He is given 10 VIP passes every hour. When a VIP arrives, he takes a pass and lets them in. If 5 VIPs arrive at once, they all get in instantly (burst). If 15 arrive, 10 get in, and 5 are rejected.
* **Leaky Bucket**: A toll booth on a highway. No matter how many cars arrive and wait in line, the toll booth operator processes exactly one car every 10 seconds.
* **Fixed Window**: A happy hour that runs strictly from 5:00 PM to 6:00 PM. At 6:01 PM, the limit resets and prices double, regardless of how long you've been sitting there.

---

# Interview Questions

## Beginner
**Q**: What HTTP status code should you return when a client exceeds their rate limit?
*Answer*: HTTP 429 (Too Many Requests). You should also include a `Retry-After` header indicating how many seconds the client must wait.

## Intermediate
**Q**: In a distributed system, why is using a simple Redis `GET` followed by a `SET` to track rate limits dangerous?
*Answer*: Race conditions. If the limit is 1, and two requests arrive simultaneously at two different gateway nodes, both might `GET` a counter of 0, both think they are allowed, and both `SET` the counter to 1. The limit of 1 was breached. You must use atomic operations (like `INCR` or Lua scripts).

## Advanced
**Q**: What happens if the centralized Redis rate limiting cluster goes down? Does your whole API go down?
*Answer*: It shouldn't. Rate limiting is a defensive mechanism, not a critical business function. You should "Fail Open". If the API Gateway cannot reach Redis, it should log a warning and allow the requests through to prevent a total system outage, perhaps falling back to a generous local, in-memory rate limit as a secondary defense.

---

# Quiz

## Multiple Choice Questions
**1. Which algorithm suffers from the "boundary problem", where traffic spikes at the edges of time windows can bypass the intended rate?**
A) Leaky Bucket
B) Sliding Window Log
C) Fixed Window Counter
*Answer*: C

## True or False
**The Token Bucket algorithm allows for sudden bursts of traffic, up to the maximum capacity of the bucket.**
*Answer*: True. If the bucket has capacity, multiple requests can consume tokens simultaneously, resulting in a burst.

---

# Exercises

## Beginner
Modify the Go code example to use `httputil.ReverseProxy` to rate limit a downstream microservice instead of just returning a string.

## Intermediate
Research Redis Lua scripts. Write a small Lua script that implements a Fixed Window counter atomically (using `INCR` and `EXPIRE`).

---

# Summary

Rate limiting is non-negotiable for public-facing APIs. It protects your databases from melting down and ensures fair resource distribution. The Token Bucket algorithm offers the best balance of flexibility and performance, while Redis provides the atomic, centralized state needed to enforce limits across a fleet of stateless microservices.

---

# Key Takeaways

* ✔ Always return HTTP 429 and a Retry-After header when limiting.
* ✔ Token Bucket allows bursts; Leaky Bucket strictly smooths traffic.
* ✔ Fixed Window counters are cheap but suffer from edge-case spikes.
* ✔ Distributed rate limiting requires atomic operations (e.g., Redis Lua scripts) to prevent race conditions.
* ✔ Always design rate limiters to "Fail Open" so a cache outage doesn't take down the primary application.

---

# Further Reading
* [Stripe Engineering: Scaling your API with rate limiters](https://stripe.com/blog/rate-limiters)
* [Envoy Proxy Rate Limiting](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/other_features/global_rate_limiting)

---

# Next Chapter
➡️ **Next:** `12-Consensus.md`
