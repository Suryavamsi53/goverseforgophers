# Retries and Exponential Backoff

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Real-World Analogy
* Core Concepts
* The Problem with Simple Retries (Thundering Herd)
* Exponential Backoff and Jitter
* Architecture Diagram
* Step-by-Step Implementation
* Syntax
* Beginner Example
* Intermediate Example (Adding Jitter)
* Advanced Example (Library usage)
* Production Use Cases
* Performance Analysis
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

In distributed systems, transient failures are a mathematical certainty. A network packet gets dropped, a router restarts, or a database experiences a split-second locking conflict.

When a transient failure occurs, the simplest and most effective solution is to just try again. **Retries** increase system reliability exponentially. However, if implemented poorly, retries will accidentally perform a Distributed Denial of Service (DDoS) attack on your own infrastructure. This chapter explores how to retry safely using **Exponential Backoff** and **Jitter**.

---

# Learning Objectives

After completing this chapter you will be able to:

* Distinguish between transient errors (worth retrying) and terminal errors (worth failing fast).
* Understand the "Thundering Herd" problem.
* Implement Exponential Backoff mathematically.
* Add random "Jitter" to prevent synchronized retry spikes.

---

# Prerequisites

Before reading this chapter you should know:

* Fallacies of Distributed Computing (`02-Fallacies-of-Distributed-Computing.md`).
* Circuit Breaker Pattern (`07-Circuit-Breaker-Pattern.md`).

---

# Why This Topic Exists

Imagine your Authentication Microservice goes offline for exactly 3 seconds to deploy a patch.
During those 3 seconds, 10,000 users try to log in. Their requests fail.

If your frontend is programmed with a "Simple Retry", all 10,000 apps will instantly retry at the exact same millisecond. 
When the Auth service boots up, it is immediately slammed with 10,000 simultaneous requests. It runs out of memory and crashes again. The apps retry again. The Auth service crashes again. You have just built an automated DDoS attack against your own company.

You must space the retries out. You must use Backoff and Jitter.

---

# Real-World Analogy

### Calling Customer Support

* **The Transient Failure**: You call the IRS on tax day. You get a busy signal.
* **Simple Retry (Bad)**: You instantly hit redial. Still busy. You hit redial 100 times in 10 seconds. You are wasting your energy and blocking the phone lines.
* **Exponential Backoff (Good)**: You decide to wait 1 minute before calling back. Busy. You wait 2 minutes. Busy. You wait 4 minutes. Busy. You wait 8 minutes. You finally get through. You saved your sanity by exponentially increasing the delay.
* **Jitter (Great)**: If a TV commercial tells 10,000 people to call at 5:00 PM, and they all use exact Exponential Backoff (1m, 2m, 4m), they will all call at exactly 5:01, 5:03, and 5:07, crashing the phone lines every time. If they add *Jitter* (randomness), one person waits 1m12s, another waits 54s. The load is smoothed out evenly.

---

# Core Concepts

* **Transient Error**: A temporary error (e.g., HTTP 503 Service Unavailable, Network Timeout). You *should* retry these.
* **Terminal Error**: A permanent error (e.g., HTTP 400 Bad Request, 401 Unauthorized). If the user provided the wrong password, retrying the wrong password 5 times will not magically make it correct. Do *not* retry these.
* **Exponential Backoff**: Multiplying the wait time by a factor (usually 2) after each failed attempt (e.g., 1s, 2s, 4s, 8s).
* **Jitter**: Adding a random number of milliseconds to the backoff duration to prevent synchronized spikes.

---

# Architecture Diagram

```mermaid
flowchart TD
    Req[Network Request]
    Check{Is Error?}
    
    Req --> Check
    Check -- "No (200 OK)" --> Success[Return Data]
    Check -- "Yes (400 Bad Request)" --> Terminal[Terminal! Abort.]
    Check -- "Yes (503 Timeout)" --> Trans{Attempts < Max?}
    
    Trans -- "No" --> Abort[Return Error to User]
    Trans -- "Yes" --> Math[Calculate Delay:<br/>(Base * 2^Attempt) + RandomJitter]
    
    Math --> Sleep[time.Sleep]
    Sleep --> Req
```

---

# Step-by-Step Implementation

1. Define a `MaxRetries` constant (e.g., 3 or 5).
2. Define a `BaseDelay` (e.g., 100ms).
3. Start a `for` loop.
4. Execute the network request.
5. If success, break the loop and return.
6. If error, check if it is a Terminal Error. If yes, break and return the error.
7. If Transient Error, calculate the backoff: `delay = BaseDelay * (2 ^ attempt)`.
8. Calculate Jitter: `jitter = rand.Intn(maxJitter)`.
9. `time.Sleep(delay + jitter)`.
10. Loop again.

---

# Syntax (The Math)

```go
// attempt is 0-indexed
delay := baseDelay * time.Duration(1<<attempt) // 1<<attempt is 2^attempt

// Add randomness between 0 and 50% of the delay
jitter := time.Duration(rand.Int63n(int64(delay / 2)))

sleepDuration := delay + jitter
time.Sleep(sleepDuration)
```

---

# Beginner Example

A standard loop with Exponential Backoff (No Jitter yet).

```go
package main

import (
	"errors"
	"fmt"
	"time"
)

// Simulate a flaky network call
func FlakyNetworkCall() (string, error) {
	return "", errors.New("503 Service Unavailable")
}

func main() {
	maxRetries := 4
	baseDelay := 100 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		fmt.Printf("Attempt %d...\n", attempt+1)
		
		data, err := FlakyNetworkCall()
		
		if err == nil {
			fmt.Println("Success:", data)
			return
		}

		// Calculate Exponential Backoff: 100ms, 200ms, 400ms, 800ms
		// 1<<attempt is bitwise shifting, equivalent to 2^attempt
		backoff := baseDelay * time.Duration(1<<attempt)
		
		fmt.Printf("Failed: %v. Backing off for %v\n\n", err, backoff)
		time.Sleep(backoff)
	}

	fmt.Println("All retries exhausted. Operation failed.")
}
```

---

# Intermediate Example

Adding **Jitter**. In a highly concurrent system, this is absolutely mandatory.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

func main() {
	maxRetries := 5
	baseDelay := 1 * time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 1. Calculate standard exponential backoff
		backoff := baseDelay * time.Duration(1<<attempt)

		// 2. Calculate Jitter (Random value between 0 and half the backoff)
		jitterAmount := rand.Int63n(int64(backoff / 2))
		jitter := time.Duration(jitterAmount)

		// 3. Add them together
		finalSleep := backoff + jitter

		fmt.Printf("Attempt %d failed. Sleeping for %v (Base: %v, Jitter: %v)\n", 
			attempt+1, finalSleep, backoff, jitter)
			
		time.Sleep(finalSleep)
	}
}
```
*Sample Output:*
```text
Attempt 1 failed. Sleeping for 1.2s (Base: 1s, Jitter: 200ms)
Attempt 2 failed. Sleeping for 2.8s (Base: 2s, Jitter: 800ms)
Attempt 3 failed. Sleeping for 4.1s (Base: 4s, Jitter: 100ms)
Attempt 4 failed. Sleeping for 11.5s (Base: 8s, Jitter: 3.5s)
```

---

# Production Use Cases

### 1. Cloud SDKs (AWS / GCP)
If you use the `aws-sdk-go` to read a file from S3, and the network drops for a microsecond, the SDK does not instantly return an error to you. Under the hood, the SDK automatically executes an Exponential Backoff and Jitter algorithm. It will retry silently 3 or 4 times before finally bubbling the error up to your application.

### 2. Database Connection Pools
When a Go application starts, the database might still be booting up. Instead of crashing the application instantly with a `Fatal` error, a robust application will try to `db.Ping()` the database in a loop with exponential backoff until the database is ready to accept connections.

---

# Performance Analysis

Retries inherently increase the latency of an operation. If a user clicks a button, and the backend retries 3 times taking a total of 7 seconds, the user will likely have closed the app. 
Therefore, **Backend-to-Backend** retries (e.g., Queue Workers) can have long, generous backoffs (minutes). **Frontend-facing** APIs should have very short backoffs (max 2-3 seconds total) and then fail fast, informing the user.

---

# Best Practices

* **Maximum Cap**: Exponential backoff grows incredibly fast (2, 4, 8, 16, 32...). If you have a bug that causes infinite retries, you will quickly be waiting years. Always enforce a `MaxDelay` cap (e.g., "never sleep for more than 30 seconds, regardless of the attempt number").
* **Distinguish Errors**: Never retry a 4xx HTTP error. Only retry 5xx errors, `context.DeadlineExceeded`, or `io.EOF` (connection drops).
* **Idempotency**: If you are sending a POST request to create a payment, and the network drops *while* waiting for the response, the server might have successfully charged the card. If you retry, you charge the card twice! (We will solve this in `11-Idempotency.md`).

---

# Common Mistakes

### Retrying at Multiple Layers
If Service A calls Service B, and Service B calls Service C.
Service C goes down. 
Service B retries 3 times. 
Service A times out waiting for B, so Service A retries 3 times. 
You now have $3 \times 3 = 9$ retries amplifying through the system. This creates a combinatorial explosion of traffic. 
**Rule of thumb**: Only retry at the *edges* of your system (or rely entirely on Circuit Breakers internally).

---

# Quiz

## Multiple Choice Questions
**1. What is the primary purpose of adding "Jitter" to an exponential backoff algorithm?**
A) To make the retries happen faster.
B) To prevent the "Thundering Herd" problem by de-synchronizing thousands of clients who all failed at the exact same moment.
C) To encrypt the retry payload.
*Answer*: B

## True or False
**If a REST API returns an HTTP 403 Forbidden error, you should retry the request with exponential backoff.**
*Answer*: False. A 403 Forbidden is a Terminal Error. Your token is invalid or lacks permissions. Retrying a billion times will never change the result; it just wastes server resources. Only retry Transient Errors (500s, network drops).

---

# Interview Questions

## Beginner
**Q**: What is the difference between a Transient Error and a Terminal Error?
*Answer*: A transient error is temporary (like a network timeout or a 503 Gateway overload). If you try again in a few seconds, it will likely succeed. A terminal error is a permanent logic failure (like a 400 Bad Request or invalid JSON). Retrying a terminal error will never succeed.

## Intermediate
**Q**: Explain Exponential Backoff and why it is superior to a constant delay (e.g., waiting exactly 1 second every time).
*Answer*: If a server is down due to extreme load, 10,000 clients retrying every 1 second will sustain the extreme load, preventing the server from ever recovering. Exponential backoff (1s, 2s, 4s, 8s) exponentially reduces the rate of incoming traffic over time, giving the struggling server the breathing room it needs to recover.

## Advanced
**Q**: How does the Circuit Breaker pattern interact with Retry logic? Should you use both?
*Answer*: They solve the same problem from opposite ends. A Circuit Breaker protects the *Server* from being hammered by shutting down outbound traffic. Retries try to protect the *Request* from failing by trying again. If you wrap a Retry loop *inside* a Circuit Breaker, the Retry loop will instantly fail fast if the circuit is open, preventing the Thundering Herd without actually executing network calls. They complement each other perfectly.

---

# Summary

Retries are the simplest way to add "9s" of reliability to your distributed systems. However, physics and math dictate that synchronized retries are a weapon of mass destruction. By strictly implementing Exponential Backoff and Jitter, you ensure that your services gracefully decay and recover during outages, rather than contributing to their own demise.

---

# Key Takeaways

* ✔ Only retry Transient errors; fail fast on Terminal errors.
* ✔ Multiply the delay by 2 on every failure.
* ✔ Add Jitter (randomness) to spread out the load.
* ✔ Be wary of combinatorial explosions when retrying across multiple microservice layers.

---

# Further Reading
* [AWS Architecture Blog: Exponential Backoff And Jitter](https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)

---

# Next Chapter
➡️ **Next:** `09-Timeouts-and-Deadlines.md`
