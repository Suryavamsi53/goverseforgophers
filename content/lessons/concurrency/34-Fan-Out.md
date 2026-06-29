# Fan-Out Pattern

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

If **Fan-In** is the process of merging many channels into one, **Fan-Out** is the exact opposite. Fan-Out is the process of taking a single input channel and distributing its work across multiple downstream worker Goroutines. 

By having multiple Goroutines reading from the exact same channel, Go automatically load-balances the work. Whichever Goroutine finishes its current task first will simply grab the next item off the channel.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand how multiple Goroutines safely read from a single channel.
* Implement the Fan-Out pattern to parallelize CPU-intensive workloads.
* Combine Fan-Out with Fan-In to build high-performance data pipelines.
* Differentiate between Fan-Out and a basic Worker Pool.

---

# Prerequisites

Before reading this chapter you should know:

* Channels (`10-Channels.md`)
* Worker Pool (`32-Worker-Pool.md`)
* Fan-In (`33-Fan-In.md`)

---

# Why This Topic Exists

Suppose you have a channel that streams 1000 video files to be encoded. Video encoding is incredibly slow. If you only have one Goroutine reading from that channel and doing the encoding, the pipeline is extremely bottlenecked.

By "fanning out", you spin up 8 Goroutines that all read from that single video channel. Now you are encoding 8 videos simultaneously, utilizing 100% of your multi-core CPU. 

---

# Real-World Analogy

### The Call Center

* **The Phone Line**: A single unified phone number (the input channel) that customers call.
* **Fan-Out**: There are 20 customer service representatives (Goroutines) all listening to that single phone line. 
* When a call comes in, the phone system routes it to the *first available* representative. 
* If a rep takes 5 minutes to help a customer, they are busy. Other reps will take the subsequent calls. This perfectly balances the load across the entire room.

---

# Core Concepts

* **Competing Consumers**: When multiple Goroutines read from the same channel, they are "competing". Go guarantees that a single value sent to the channel will only be received by **exactly one** Goroutine.
* **Load Balancing**: Because Goroutines block when they are busy, fast workers will naturally process more items from the channel than slow workers.
* **Fan-Out vs Worker Pool**: They are mechanically identical. "Worker Pool" usually refers to the architectural design of limiting concurrency, while "Fan-Out" refers specifically to the data flow (one channel splitting to many).

---

# Architecture Diagram

```mermaid
flowchart LR
    Producer[Producer Goroutine]
    InQ[(Input Channel)]
    
    subgraph Fan-Out
    W1[Worker Goroutine 1]
    W2[Worker Goroutine 2]
    W3[Worker Goroutine 3]
    end

    Producer -- Sends Data --> InQ
    
    InQ -.-> W1
    InQ -.-> W2
    InQ -.-> W3
    
    note over InQ,W2: Each piece of data goes to exactly ONE worker
```

---

# Step-by-Step Implementation

1. Create a single `jobs` channel.
2. Launch `N` Goroutines (Fan-Out).
3. Inside each Goroutine, write a `for job := range jobs` loop.
4. From the main producer Goroutine, send data into the `jobs` channel.
5. Close the `jobs` channel when you are done producing.
6. The `N` Goroutines will all exit their `range` loops automatically when the channel is empty and closed.

---

# Syntax

```go
// 1. Create channel
jobs := make(chan int)

// 2. Fan-Out to 3 workers
for i := 0; i < 3; i++ {
    go func() {
        for j := range jobs { // 3. Competing Consumers
            fmt.Println("Processed", j)
        }
    }()
}
```

---

# Beginner Example

Fanning out to 3 workers to process 10 jobs.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	jobs := make(chan int, 10)
	var wg sync.WaitGroup

	// Fan-Out: Start 3 competing consumers
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", id, job)
				time.Sleep(500 * time.Millisecond) // Simulate work
			}
			fmt.Printf("Worker %d done\n", id)
		}(w)
	}

	// Producer
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs) // Signal that no more jobs are coming

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All work complete.")
}
```

---

# Intermediate Example

Combining Fan-Out and Fan-In to create a full Map-Reduce pipeline. We Fan-Out work to calculate squares, and Fan-In the results to a single slice.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// The Worker: Takes a job, does math, sends to results
func worker(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for n := range jobs {
		time.Sleep(100 * time.Millisecond)
		results <- n * n // Square the number
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	var wg sync.WaitGroup

	// FAN-OUT: Start 4 workers reading from the 'jobs' channel
	for w := 1; w <= 4; w++ {
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// PRODUCER: Send 10 numbers
	go func() {
		for i := 1; i <= 10; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// FAN-IN: Wait for all workers, then close results
	go func() {
		wg.Wait()
		close(results)
	}()

	// CONSUMER: Read the merged results
	var final []int
	for res := range results {
		final = append(final, res)
	}

	fmt.Println("Final squared results:", final)
}
```

---

# Advanced Example

Handling unfair workloads. If one job takes 10 seconds, and 9 jobs take 1 second, Fan-Out automatically ensures the other workers process the 9 fast jobs while the slow job occupies a single worker.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tasks := make(chan string, 10)
	var wg sync.WaitGroup

	// Fan-Out to 2 workers
	for w := 1; w <= 2; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for task := range tasks {
				fmt.Printf("Worker %d picked up %s\n", id, task)
				if task == "SLOW TASK" {
					time.Sleep(3 * time.Second)
				} else {
					time.Sleep(200 * time.Millisecond)
				}
				fmt.Printf("Worker %d finished %s\n", id, task)
			}
		}(w)
	}

	// Send 1 slow task, and 5 fast tasks
	tasks <- "SLOW TASK"
	for i := 1; i <= 5; i++ {
		tasks <- fmt.Sprintf("Fast Task %d", i)
	}
	close(tasks)

	wg.Wait()
}
// Notice in the output: Worker 1 gets stuck on SLOW TASK. 
// Worker 2 automatically picks up the slack and processes ALL the Fast Tasks!
```

---

# Production Use Cases

### 1. Data Processing Pipelines
Parsing massive JSON files. The main thread reads the file line by line and pushes raw JSON strings to a channel. Fan-Out is used to distribute these strings to 10 Goroutines that unmarshal the JSON into Go structs simultaneously.

### 2. Email Sending Service
When a marketing campaign is triggered, 100,000 emails need to be sent. A single producer queries the DB for email addresses and pushes them to a channel. A Fan-Out of 50 SMTP worker Goroutines read from the channel and send the emails over the network concurrently.

---

# Performance Analysis

* **Lock Contention**: When multiple Goroutines read from the exact same channel, Go internally uses a Mutex to ensure only one Goroutine gets the value. While highly optimized, if you have 10,000 Goroutines fanning out from a single channel, the Mutex lock contention on the channel itself can become a bottleneck.
* **The Sweet Spot**: Keep your Fan-Out worker count reasonable (e.g., `runtime.NumCPU()` for CPU tasks, or ~100 for I/O tasks). 

---

# Best Practices

* **Always use WaitGroups**: To know when the Fan-Out phase is complete, you must pass a `sync.WaitGroup` to every worker and wait for them to finish.
* **Buffer the Input Channel**: The channel you are fanning out from should usually be buffered. This allows the producer to rapidly enqueue work without waiting for a worker to become available for every single item.

---

# Common Mistakes

### Forgetting to close the channel
If the producer never calls `close(jobs)`, the Fan-Out workers will sit in their `for range` loops forever, waiting for more data. This is a Goroutine leak and will eventually cause a deadlock or OOM crash.

---

# Debugging Guide

* **Output order is mixed up**: This is expected! Fan-Out destroys ordering. Because workers run concurrently, Job 2 might finish before Job 1. If you need strictly ordered output, you must pass an index with the job and sort the results array at the very end.

---

# Exercises

## Beginner
Create a channel of strings containing 5 website URLs. Fan-Out to 3 workers that simply print "Downloading {URL}".

## Intermediate
Create a Fan-Out pipeline that processes integers. The workers should multiply the integer by 10. Fan-In the results to a slice. Prove that the final slice is not in the same order as the input (1, 2, 3, 4, 5).

---

# Quiz

## Multiple Choice Questions
**1. In a Fan-Out pattern, if 3 Goroutines are ranging over the same channel, and a single integer `5` is sent to the channel, which Goroutines receive it?**
A) All 3 Goroutines receive a copy of `5`.
B) Only one Goroutine receives `5`.
C) The channel panics because of multiple readers.
*Answer*: B

## True or False
**Fan-Out is useful for ensuring data is processed in the exact same chronological order it was received.**
*Answer*: False. Fan-Out inherently destroys ordering because Goroutines execute independently and finish at unpredictable times.

---

# Interview Questions

## Beginner
**Q**: What is the difference between Fan-In and Fan-Out?
*Answer*: Fan-In merges multiple channels into one. Fan-Out distributes one channel to multiple reader Goroutines.

## Intermediate
**Q**: How does Go handle the race condition of multiple Goroutines reading from the exact same channel at the same time?
*Answer*: Channels are completely thread-safe. Go's runtime uses internal locks (and lock-free fast paths) to ensure that a value in a channel is atomically handed off to exactly one waiting receiver.

## Advanced
**Q**: If you are fanning out to 1,000 workers to do HTTP requests, and the third-party API starts rate-limiting you, how do you gracefully slow down the Fan-Out?
*Answer*: You can use a Token Bucket rate limiter (like `golang.org/x/time/rate`). Before a worker pulls from the channel, or immediately after, it must call `limiter.Wait(ctx)`. This globally restricts the speed of all workers without needing to change the Fan-Out architecture.

---

# Mini Project

**Requirement**: The Prime Number Finder.
1. Generate numbers from 1 to 100 on a `numbers` channel.
2. Fan-Out to 4 Goroutines.
3. Each worker checks if the number is prime. (You can write a simple `isPrime(n)` function).
4. If it is prime, send it to a `primes` channel.
5. Fan-In the results, collect them into a slice, sort the slice, and print it.

---

# Cheat Sheet

* **Fan-Out Pattern**:
```go
jobs := make(chan int, 100)
for w := 0; w < runtime.NumCPU(); w++ {
    go func() {
        for j := range jobs {
            // Process j
        }
    }()
}
```

---

# Summary

Fan-Out is the engine that drives high-performance parallel processing in Go. By leveraging the fact that channels safely multiplex readers, you can build self-balancing workloads that automatically route around slow tasks and maximize CPU utilization.

---

# Key Takeaways

* ✔ Distributes work from one channel to many Goroutines.
* ✔ Automatically load-balances tasks.
* ✔ Destroys the original ordering of the data.
* ✔ Requires a closed input channel to cleanly shut down workers.

---

# Further Reading
* [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)

---

# Next Chapter
➡️ **Next:** `35-Pipeline.md`
