# Semaphore

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

A **Mutex** allows exactly *one* Goroutine to access a resource at a time. 
A **Semaphore** allows exactly *N* Goroutines to access a resource at a time. 

While Go does not have a built-in `sync.Semaphore` type in the standard library (like Java or C# do), the community standard is to implement Semaphores using **Buffered Channels**. Alternatively, the extended standard library provides a highly optimized semaphore in `golang.org/x/sync/semaphore`.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand the difference between a Mutex and a Semaphore.
* Implement a basic Semaphore using a Buffered Channel.
* Use `golang.org/x/sync/semaphore` for advanced, weighted concurrency limits.
* Apply Semaphores to rate-limit external API calls.

---

# Prerequisites

Before reading this chapter you should know:

* Mutexes (`21-Mutex.md`)
* Buffered Channels (`11-Buffered-Channels.md`)

---

# Why This Topic Exists

Imagine you have a microservice that generates PDF reports. PDF generation takes 2 seconds and uses 500MB of RAM. If 10 users request a PDF, that's 5GB of RAM. If 100 users request a PDF simultaneously, your server tries to allocate 50GB of RAM and crashes instantly (OOM Kill).

You can't use a Mutex, because that would force 1 PDF to generate at a time (too slow). You need a way to say: "Allow up to exactly 4 PDFs to generate simultaneously, and make anyone else wait in line." This is exactly what a Semaphore does.

---

# Real-World Analogy

### The Nightclub Bouncer

* **Mutex**: A single-occupancy bathroom. The lock (Mutex) only allows 1 person inside. Everyone else waits in line.
* **Semaphore**: The Nightclub Bouncer. The fire code says the club can hold exactly 100 people. 
  - The Bouncer starts with 100 tickets.
  - As a person enters, the Bouncer takes a ticket (Acquire).
  - When the club has 100 people, the Bouncer has 0 tickets. He puts up his hand and forces the line to wait.
  - When someone leaves the club, they hand their ticket back to the Bouncer (Release).
  - The Bouncer immediately lets the next person in line enter.

---

# Core Concepts

* **Capacity (Weight)**: The maximum number of concurrent operations allowed.
* **Acquire (Wait)**: A Goroutine requests permission to proceed. If the semaphore is full, the Goroutine blocks (sleeps).
* **Release (Signal)**: A Goroutine finishes its work and returns its permission to the semaphore, waking up a waiting Goroutine.

---

# Architecture Diagram

```mermaid
flowchart TD
    Req1[Goroutine 1]
    Req2[Goroutine 2]
    Req3[Goroutine 3]
    Req4[Goroutine 4 (Waiting)]
    
    Sem{Semaphore<br/>Limit: 3}
    
    DB[(Database)]
    
    Req1 -- Acquire --> Sem
    Req2 -- Acquire --> Sem
    Req3 -- Acquire --> Sem
    Req4 -- Acquire (Blocks) --> Sem
    
    Sem --> DB
```

---

# Step-by-Step Implementation (Using Channels)

1. Create a buffered channel of empty structs: `sem := make(chan struct{}, 3)`.
2. To **Acquire**: Push a value into the channel `sem <- struct{}{}`. If the channel is full (already has 3 items), this operation will block the Goroutine until there is space.
3. To **Release**: Pull a value out of the channel `<-sem`. This instantly creates space in the channel, allowing one blocked Goroutine to proceed.

---

# Syntax

```go
// Create a Semaphore with capacity 3
sem := make(chan struct{}, 3)

// Acquire
sem <- struct{}{}

// Do work...

// Release
<-sem
```
*(Note: We use `struct{}{}` because an empty struct takes exactly 0 bytes of memory in Go).*

---

# Beginner Example

Limiting 10 Goroutines to only execute 2 at a time using a Buffered Channel.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// A Semaphore that allows 2 concurrent operations
	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			fmt.Printf("Worker %d waiting for semaphore...\n", id)
			
			// ACQUIRE: Blocks if 2 workers are already running
			sem <- struct{}{} 
			
			fmt.Printf("--> Worker %d ACQUIRED semaphore and is running!\n", id)
			time.Sleep(1 * time.Second) // Simulate work
			
			// RELEASE: Frees up a slot
			fmt.Printf("<-- Worker %d RELEASING semaphore\n", id)
			<-sem 
			
		}(i)
	}

	wg.Wait()
	fmt.Println("All done.")
}
```

---

# Intermediate Example

Using the official `golang.org/x/sync/semaphore` package. This package is better than the channel approach because it provides Context cancellation and "Weighted" acquires.

```go
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

func main() {
	// Create a semaphore with a total weight of 3
	sem := semaphore.NewWeighted(3)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// Acquire a weight of 1. Blocks if weight limit (3) is reached.
			if err := sem.Acquire(ctx, 1); err != nil {
				fmt.Printf("Failed to acquire for worker %d\n", id)
				return
			}
			
			// Always use defer to ensure Release is called, even on panics!
			defer sem.Release(1)

			fmt.Printf("Worker %d running\n", id)
			time.Sleep(1 * time.Second)
			
		}(i)
	}

	wg.Wait()
}
```

---

# Advanced Example

Using Weighted Semaphores. Imagine your server has 8GB of RAM. You have two types of jobs: "Thumbnail" (uses 1GB) and "4k Video" (uses 4GB). You can use a Weighted Semaphore to track actual memory usage!

```go
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/semaphore"
	"sync"
	"time"
)

func main() {
	// Our server has "8" units of capacity (e.g., 8GB RAM)
	maxRAM := int64(8)
	sem := semaphore.NewWeighted(maxRAM)
	ctx := context.Background()
	var wg sync.WaitGroup

	runJob := func(name string, cost int64) {
		defer wg.Done()
		fmt.Printf("Job %s wants %dGB...\n", name, cost)
		
		// Blocks until enough total capacity is available!
		sem.Acquire(ctx, cost) 
		defer sem.Release(cost)

		fmt.Printf("+++ Job %s started! (Consuming %dGB)\n", name, cost)
		time.Sleep(2 * time.Second)
		fmt.Printf("--- Job %s finished! (Releasing %dGB)\n", name, cost)
	}

	wg.Add(3)
	go runJob("Video 1", 4) // Takes 4GB
	go runJob("Video 2", 4) // Takes 4GB (Server is now full!)
	time.Sleep(100 * time.Millisecond)
	
	// This Thumbnail wants 1GB, but 8/8GB is used. It MUST wait!
	go runJob("Thumbnail", 1) 

	wg.Wait()
}
```

---

# Production Use Cases

### 1. API Rate Limiting (Outbound)
If you are calling a third-party API (like Stripe or Twilio) that allows exactly 20 requests per second, wrapping the HTTP call in a Semaphore with a capacity of 20 ensures your microservice will never accidentally overwhelm their servers and get temporarily banned.

### 2. Protecting Heavy Endpoints (Inbound)
If a specific HTTP handler in your Chi or Gin router does a massive SQL `GROUP BY` query, you can put a Semaphore inside the handler. If 100 users hit it at once, only 5 will execute the query simultaneously, while the other 95 simply block and wait their turn. The SQL database survives the traffic spike.

---

# Performance Analysis

* **Channel Semaphore**: Extremely fast and requires no external dependencies. Perfect for simple limits (e.g., max 10 concurrent requests).
* **Weighted Semaphore**: Slightly slower (uses a Mutex internally), but offers incredible flexibility for resource-aware limits (e.g., tracking Megabytes of RAM used instead of just tracking the number of Goroutines).
* **Semaphore vs Worker Pool**: Both limit concurrency. A Worker Pool is usually better for background processing of queues. A Semaphore is usually better for inline HTTP request handlers where you want the caller Goroutine to block until it gets a turn.

---

# Best Practices

* **Always Defer Release**: Just like `mu.Unlock()`, you must always `defer sem.Release(weight)` or `defer func() { <-sem }()`. If a Goroutine panics or returns an error early without releasing, that semaphore slot is permanently lost (a capacity leak).
* **Use `TryAcquire`**: If a user hits a busy endpoint, you might not want them to wait for 30 seconds. The `x/sync/semaphore` package has `TryAcquire()`. If it returns false, you can instantly return an `HTTP 429 Too Many Requests` error to the user.

---

# Common Mistakes

### Forgetting to Release
```go
sem.Acquire(ctx, 1)
if err != nil {
    return // BAD: If an error happens here, the semaphore is never released!
}
sem.Release(1)

// GOOD:
sem.Acquire(ctx, 1)
defer sem.Release(1)
if err != nil {
    return // Safe, defer will run
}
```

---

# Debugging Guide

* **Server hangs under load**: You probably have a capacity leak. Check that every single `Acquire` is paired with a `defer Release`. If one Goroutine fails and skips the release, the capacity permanently drops. If it drops to 0, all future requests deadlock.
* **"Failed to acquire"**: If you request a weight *larger* than the total capacity of the semaphore (e.g., `max=8`, you ask for `10`), `Acquire` will block forever.

---

# Exercises

## Beginner
Use a buffered channel of size 3 to create a semaphore. Launch 10 Goroutines that print their ID, sleep for 1 second, and exit. Watch the terminal to verify only 3 print at a time.

## Intermediate
Create an HTTP server using `net/http`. Inside the handler for `/heavy`, use `semaphore.NewWeighted(2)`. Use `TryAcquire`. If the semaphore is full, return a 429 status code. Use an API testing tool (like `wrk` or `hey`) to bombard it and verify the 429s are working.

---

# Quiz

## Multiple Choice Questions
**1. How do you "Acquire" a slot in a channel-based semaphore?**
A) By reading from the channel `<-sem`
B) By pushing to the channel `sem <- struct{}{}`
C) By closing the channel
*Answer*: B. (Pushing into a buffered channel blocks if it is full, which is the exact behavior we want for acquiring a slot).

## True or False
**A Mutex is just a Semaphore with a capacity of exactly 1.**
*Answer*: True. In computer science theory, a binary semaphore (capacity 1) is functionally equivalent to a standard Mutex lock.

---

# Interview Questions

## Beginner
**Q**: What is the difference between a Mutex and a Semaphore?
*Answer*: A Mutex provides mutual exclusion (only 1 Goroutine can access the resource). A Semaphore provides bounded concurrency (up to N Goroutines can access the resource simultaneously).

## Intermediate
**Q**: When would you use a Weighted Semaphore instead of a simple Channel Semaphore?
*Answer*: You use a Weighted Semaphore when different tasks consume different amounts of resources. For example, if task A uses 1GB of RAM and task B uses 5GB of RAM, a channel semaphore just counts "2 tasks", but a weighted semaphore can accurately block if the total requested memory exceeds the server's limits.

## Advanced
**Q**: Compare using a Semaphore vs using a Worker Pool to protect a database from too many concurrent connections.
*Answer*: 
A **Semaphore** is inline blocking. If 100 HTTP requests come in, 100 Goroutines are spawned, 10 get the database semaphore, and 90 Goroutines sit idle in memory waiting their turn. This is simpler to write but consumes memory for the 90 sleeping Goroutines.
A **Worker Pool** is queue-based. 100 HTTP requests come in, they push an ID to a channel, and exactly 10 Worker Goroutines read from the channel and hit the database. This is more memory-efficient and provides better backpressure control.

---

# Mini Project

**Requirement**: The Museum Tour Guide.
1. A museum can only hold exactly 5 people at a time.
2. 20 Tourists (Goroutines) arrive at the door at the same time.
3. Use a Semaphore (Channel or `x/sync`) to strictly enforce the limit.
4. A Tourist should: Acquire, print "Tourist X entered", sleep 500ms, print "Tourist X left", Release.
5. Watch the output to confirm exactly 5 tourists are inside at any given time.

---

# Cheat Sheet

* **Channel Semaphore (Init)**: `sem := make(chan struct{}, limit)`
* **Channel Semaphore (Acquire)**: `sem <- struct{}{}`
* **Channel Semaphore (Release)**: `<-sem`
* **Weighted (Init)**: `sem := semaphore.NewWeighted(limit)`
* **Weighted (Acquire)**: `sem.Acquire(ctx, weight)`
* **Weighted (Release)**: `defer sem.Release(weight)`

---

# Summary

Semaphores are the ultimate "bouncer" for your Go applications. By placing them around expensive operations, you guarantee that sudden spikes in traffic will simply cause requests to wait nicely in line, rather than crashing your server in a fiery Out-Of-Memory explosion.

---

# Key Takeaways

* ✔ Limits concurrency to exactly `N` Goroutines.
* ✔ Idiomatic Go uses buffered channels of `struct{}{}`.
* ✔ `x/sync/semaphore` allows for complex weighted limits.
* ✔ Always use `defer` to release the semaphore.

---

# Further Reading
* [Go documentation for x/sync/semaphore](https://pkg.go.dev/golang.org/x/sync/semaphore)

---

# Next Chapter
➡️ **Next:** `37-Rate-Limiter.md`
