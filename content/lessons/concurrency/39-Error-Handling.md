# Concurrent Error Handling

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Core Concepts
* Architecture Diagram
* Step-by-Step Implementation
* Syntax
* Beginner Example
* Intermediate Example
* Advanced Example
* Production Use Cases
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

In a normal, synchronous Go program, error handling is straightforward: you call a function, it returns `(value, err)`, and you check `if err != nil`.

In a highly concurrent program, functions are executed inside Goroutines running in the background. They cannot simply `return` an error to the main thread. If a background Goroutine panics or encounters a fatal error, you need a robust mechanism to capture that error, bring it back to the main thread, and potentially cancel all other running Goroutines.

---

# Learning Objectives

After completing this chapter you will be able to:

* Safely transmit errors from Goroutines back to the main thread using channels.
* Use `errgroup.Group` to dramatically simplify concurrent error handling.
* Handle Goroutine panics using `defer` and `recover`.
* Ensure that one failing Goroutine instantly cancels all sibling Goroutines.

---

# Prerequisites

Before reading this chapter you should know:

* Goroutines (`08-Goroutines.md`)
* Channels (`10-Channels.md`)
* Context (`19-Context.md`)
* Errgroup (`27-errgroup.md`)

---

# Why This Topic Exists

Imagine you have a Worker Pool downloading 100 images. Worker 5 encounters an "HTTP 401 Unauthorized" error. 
If Worker 5 just prints the error and exits, the other 99 workers will uselessly continue trying to download images (and failing) because the API key is expired. 
If Worker 5 tries to `panic`, the entire Go application crashes immediately, bringing down the whole server.

You need a way for Worker 5 to safely tell the main thread: "I hit a fatal error. Stop everyone else, and report this error back to the user."

---

# Core Concepts

* **Error Channels**: A channel specifically created to pass `error` types (`chan error`).
* **First-Error-Wins**: In concurrent operations, usually the *first* error encountered by any Goroutine is the one that should trigger a global cancellation.
* **Panic Propagation**: By default, a panic inside a Goroutine crashes the whole program. You must explicitly recover it and convert it to a standard error.

---

# Architecture Diagram

```mermaid
flowchart TD
    Main[Main Thread]
    ErrQ[(Error Channel)]
    Ctx[(Context Cancel)]
    
    W1[Worker 1]
    W2[Worker 2 (Fails!)]
    W3[Worker 3]
    
    W2 -- Sends Error --> ErrQ
    ErrQ -- Main reads first error --> Main
    Main -- Triggers Cancel --> Ctx
    
    Ctx -.-> W1
    Ctx -.-> W3
    note right of W1: Remaining workers stop<br/>due to Context cancellation
```

---

# Step-by-Step Implementation (Raw Channels)

If you aren't using `errgroup`, here is how you manually handle concurrent errors:
1. Create a buffered error channel: `errCh := make(chan error, numWorkers)`. (It MUST be buffered so workers don't block if the main thread stops reading).
2. Pass the channel to the workers.
3. If a worker hits an error, send it: `errCh <- err`.
4. In the main thread, use a WaitGroup to close the error channel when all workers are done.
5. Range over the error channel. You can collect all errors, or just return the first one.

---

# Syntax

```go
errCh := make(chan error, 1) // Buffered to 1 for "first error wins"

go func() {
    if err := doWork(); err != nil {
        // Use select to send the error without blocking if the channel is full
        select {
        case errCh <- err:
        default:
        }
    }
}()
```

---

# Beginner Example

Collecting multiple errors from multiple Goroutines.

```go
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	// Buffer must be equal to the number of workers to prevent blocking
	errCh := make(chan error, 3) 

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			if id%2 != 0 {
				errCh <- fmt.Errorf("worker %d failed", id)
			}
		}(i)
	}

	// Wait in the background, then close the channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Range over the closed channel to collect all errors
	for err := range errCh {
		fmt.Println("Caught error:", err)
	}
	fmt.Println("Done")
}
```

---

# Intermediate Example

The modern, idiomatic way to handle concurrent errors is using the `golang.org/x/sync/errgroup` package. It handles the WaitGroup, the Error Channel, and the Context Cancellation for you.

```go
package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"time"
)

func main() {
	// Create a Group that returns a cancelable context
	g, ctx := errgroup.WithContext(context.Background())

	for i := 1; i <= 3; i++ {
		id := i // Capture loop variable
		
		// Launch a Goroutine via the group
		g.Go(func() error {
			fmt.Printf("Worker %d started\n", id)
			
			// Simulate long work, listening for cancellation
			select {
			case <-ctx.Done():
				fmt.Printf("Worker %d cancelled!\n", id)
				return ctx.Err()
			case <-time.After(500 * time.Millisecond):
				// Simulate a failure in Worker 2
				if id == 2 {
					return fmt.Errorf("fatal database error in worker 2")
				}
				fmt.Printf("Worker %d finished successfully\n", id)
				return nil
			}
		})
	}

	// g.Wait() blocks until all Goroutines finish.
	// It returns the VERY FIRST error it receives, and automatically 
	// cancels the context to abort the remaining Goroutines.
	if err := g.Wait(); err != nil {
		fmt.Println("Pipeline failed with:", err)
	} else {
		fmt.Println("All workers succeeded.")
	}
}
```

---

# Advanced Example

Handling **Panics**. A panic in a Goroutine will crash the entire server. To prevent this, every long-running background worker should have a `defer recover()` block that converts the panic into an error.

```go
package main

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// SafeWorker wraps a standard function with panic recovery
func SafeWorker(wg *sync.WaitGroup, errCh chan<- error, task func() error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		
		// The Recovery Block
		defer func() {
			if r := recover(); r != nil {
				// Convert panic into an error, attach the stack trace
				errCh <- fmt.Errorf("PANIC RECOVERED: %v\nStack: %s", r, debug.Stack())
			}
		}()

		// Execute actual work
		if err := task(); err != nil {
			errCh <- err
		}
	}()
}

func main() {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	// A terribly written function that panics
	badFunc := func() error {
		var pointer *int
		*pointer = 10 // Nil pointer panic!
		return nil
	}

	SafeWorker(&wg, errCh, badFunc)

	go func() {
		wg.Wait()
		close(errCh)
	}()

	// The main thread catches the panic as a normal error and stays alive!
	for err := range errCh {
		fmt.Println("Caught in main:", err)
	}
	
	fmt.Println("Server is still running...")
}
```

---

# Production Use Cases

### 1. Web Scraper Aborts
You fan-out to scrape 100 pages of a website. If the first page returns `403 Forbidden`, you don't want to hit the server 99 more times. By returning the error to an `errgroup`, the context is immediately cancelled, aborting the other 99 HTTP requests in-flight, saving massive amounts of bandwidth.

### 2. HTTP Server Request Scoping
When a user hits your API, you might spawn 3 Goroutines to fetch their User profile, their Billing status, and their Friends list concurrently. If the Billing database is down, you must catch that error, cancel the User and Friends queries immediately, and return a `500 Internal Server Error` to the user.

---

# Best Practices

* **Always Buffer Error Channels**: If you manually create an `errCh`, its buffer size must be greater than or equal to the number of Goroutines. If it is unbuffered, and 2 Goroutines error at the same time, the second Goroutine will deadlock trying to send to the channel.
* **Use errgroup**: For 99% of professional Go code, `golang.org/x/sync/errgroup` is the standard for handling concurrent errors. Don't reinvent the wheel with raw channels unless you have a highly specific streaming requirement.
* **Recover Panics at the Edge**: In HTTP servers, frameworks like Chi and Gin automatically recover panics in the request thread. But if *you* spawn a Goroutine `go func() {}` inside that handler, it is detached from the framework's recovery. You MUST add your own `defer recover()` inside that specific Goroutine.

---

# Common Mistakes

### The Unbuffered Error Deadlock
```go
// BAD: Unbuffered channel
errCh := make(chan error)

for i := 0; i < 10; i++ {
    go func() {
        // If two workers error, the first sends. 
        // The second blocks forever because no one is reading yet!
        errCh <- errors.New("boom") 
    }()
}
```

---

# Debugging Guide

* **Program exits unexpectedly without an error log**: A background Goroutine panicked. Check your logs for the standard Go panic stack trace. Always add a recovery block to long-running background tasks.
* **Goroutine Leak (pprof)**: If your workers are blocking on `errCh <- err`, it means your error channel buffer was too small, and the main thread stopped reading from it.

---

# Exercises

## Beginner
Write a script with 5 Goroutines. 3 of them return `nil`, 2 of them return an error using a buffered `errCh`. Have the main thread print exactly how many errors occurred.

## Intermediate
Refactor the beginner exercise to use `errgroup.Group`. Notice how `g.Wait()` only returns the *first* error it receives, not all of them.

---

# Quiz

## Multiple Choice Questions
**1. Why does an `errgroup` automatically cancel the Context when a worker returns an error?**
A) To save memory.
B) To trigger "fail-fast" behavior, aborting all other workers since the overall operation has already failed.
C) To restart the failed worker.
*Answer*: B

## True or False
**If a background Goroutine panics, the main thread will automatically catch it and keep the application running.**
*Answer*: False. A panic in *any* Goroutine will crash the entire Go program unless that specific Goroutine has a `defer recover()` block.

---

# Interview Questions

## Beginner
**Q**: How do you return an error from a Goroutine?
*Answer*: You cannot use the `return` keyword to pass an error back to the caller like a normal function. You must push the error into a shared `chan error`, or use a synchronization package like `errgroup`.

## Intermediate
**Q**: When using a raw `chan error` to collect errors from 10 workers, what size should the channel's buffer be?
*Answer*: The buffer should be at least 10. If it is smaller, and multiple workers fail simultaneously, the channel will fill up. The remaining failing workers will block forever trying to send their error to a full channel, creating a Goroutine leak.

## Advanced
**Q**: How does `errgroup` implement the "first-error-wins" pattern under the hood?
*Answer*: Internally, `errgroup` uses a `sync.Once` wrapper around the error assignment. When multiple workers return an error, the `sync.Once.Do()` ensures that only the very first error is saved to the group's state and triggers the context cancellation. All subsequent errors are simply ignored.

---

# Cheat Sheet

* **Errgroup Pattern**:
```go
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error {
    // Listen for ctx.Done()
    return fmt.Errorf("fail")
})
if err := g.Wait(); err != nil {
    // Handle the first error
}
```
* **Panic Recovery**:
```go
go func() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered:", r)
        }
    }()
    // risky code
}()
```

---

# Summary

Error handling in a concurrent environment requires a shift in mindset. You must assume that errors can happen simultaneously across multiple threads, and that panics are lethal to the entire application. By mastering buffered error channels and the `errgroup` package, you can build resilient systems that gracefully degrade and abort early when things go wrong.

---

# Key Takeaways

* ✔ Background Goroutines require `chan error` to report failures.
* ✔ Unbuffered error channels cause deadlocks; always use buffers.
* ✔ `errgroup` is the industry standard for concurrent error handling.
* ✔ A panic in a Goroutine crashes the whole server; always recover them.

---

# Further Reading
* [Go documentation for x/sync/errgroup](https://pkg.go.dev/golang.org/x/sync/errgroup)

---

# Next Chapter
➡️ **Next:** `40-Testing-Concurrency.md`
