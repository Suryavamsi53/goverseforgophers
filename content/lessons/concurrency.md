# Go Concurrency: Goroutines and Channels

Concurrency in Go is a first-class citizen. Unlike threads in languages like Java or C++, Go uses **goroutines**, which are lightweight, user-space threads managed by the Go runtime.

## Goroutines

To start a new goroutine, you simply use the `go` keyword followed by a function call:

```go
package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from goroutine!")
}

func main() {
	go sayHello() // Starts a new goroutine
	
	// We need to wait a bit, otherwise main might exit before the goroutine runs
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Hello from main!")
}
```

## Channels

Goroutines communicate via **channels**. Channels provide a way for two goroutines to synchronize execution and communicate by passing a value of a specified element type.

```go
package main

import "fmt"

func main() {
	messages := make(chan string) // Create a new channel of strings

	go func() {
		messages <- "ping" // Send a value into the channel
	}()

	msg := <-messages // Receive a value from the channel
	fmt.Println(msg)
}
```

### Best Practices
- **Don't communicate by sharing memory; share memory by communicating.**
- Always close channels if the receiver needs to know that no more values will be sent.
- Be careful with unbuffered channels, they can lead to deadlocks if not handled properly.


# Go Concurrency Handbook

## Goal

By the end of this chapter, you should be able to confidently design, implement, debug, and optimize concurrent Go applications for production systems.

This chapter assumes no prior concurrency knowledge and gradually progresses to advanced concepts used in large-scale backend services.

---

# Table of Contents

## Part I — Foundations

1. What is Concurrency?
2. Why Concurrency Matters
3. Concurrency vs Parallelism
4. Processes vs Threads vs Goroutines
5. How Go Handles Concurrency
6. Understanding the Go Scheduler
7. The G-M-P Model Explained
8. Goroutines
9. Goroutine Lifecycle
10. Goroutine Stack Growth
11. Scheduler Preemption

## Part II — Synchronization

12. WaitGroup
13. Mutex
14. RWMutex
15. Atomic Operations
16. sync.Once
17. sync.Cond
18. sync.Map

## Part III — Channels

19. Introduction to Channels
20. Unbuffered Channels
21. Buffered Channels
22. Channel Directions
23. Closing Channels
24. Iterating with range
25. Nil Channels
26. Select Statement
27. Default Case
28. Timeouts
29. Tickers
30. Timers

## Part IV — Context

31. context.Background()
32. context.TODO()
33. WithCancel
34. WithTimeout
35. WithDeadline
36. WithValue
37. Graceful Shutdown
38. Request Cancellation

## Part V — Concurrency Patterns

39. Worker Pool
40. Fan-Out
41. Fan-In
42. Pipeline
43. Producer-Consumer
44. Semaphore Pattern
45. Rate Limiter
46. Batch Processing
47. Job Queue
48. Publish/Subscribe

## Part VI — Production Engineering

49. Race Conditions
50. Deadlocks
51. Livelocks
52. Starvation
53. Memory Model
54. Happens-Before Relationship
55. Escape Analysis
56. Garbage Collection
57. Performance Optimization
58. Benchmarking
59. Profiling
60. Debugging Concurrent Programs

## Part VII — Real Projects

61. Concurrent HTTP Downloader
62. Image Processing Pipeline
63. CSV Processor
64. Log Processing Service
65. Background Job Scheduler
66. Notification Service
67. API Gateway Request Fan-Out
68. Concurrent Web Crawler
69. Distributed Worker Pool
70. File Synchronization Tool

## Part VIII — Interview Preparation

71. Beginner Questions
72. Intermediate Questions
73. Advanced Questions
74. Google-Style Questions
75. Common Pitfalls
76. Best Practices

## Part IX — Exercises

77. Hands-on Labs
78. Coding Challenges
79. Mini Projects
80. Final Capstone Project

---

# Chapter Template

Every chapter should follow the same learning structure:

* Learning Objectives
* Prerequisites
* Real-World Analogy
* Concept Overview
* Internal Working
* Visual Diagram
* Step-by-Step Execution
* Syntax
* Beginner Example
* Intermediate Example
* Production Example
* Memory Flow
* Common Mistakes
* Performance Notes
* Best Practices
* Interview Questions
* Practice Exercises
* Summary
* Key Takeaways
* Further Reading

---

# Example Production Projects

* Concurrent File Processing Engine
* HTTP Request Aggregator
* Email Notification Worker Pool
* Image Thumbnail Generator
* Batch Database Importer
* Distributed Task Queue
* Log Processing Pipeline
* Background Cache Refresher
* Metrics Collector
* Real-Time Event Processor

---

# Learning Outcome

After completing this handbook, you should be able to:

* Understand how the Go scheduler works.
* Use goroutines safely and efficiently.
* Design communication using channels.
* Prevent race conditions and deadlocks.
* Build scalable worker pools and pipelines.
* Apply cancellation using context.
* Debug concurrent applications.
* Optimize concurrency performance.
* Build production-ready backend services.
* Confidently answer Go concurrency interview questions.
