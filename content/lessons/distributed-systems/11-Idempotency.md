# Idempotency

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
* Beginner Example
* Intermediate Example (Idempotency Keys)
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

In mathematics and computer science, **Idempotency** is the property of an operation whereby it can be applied multiple times without changing the result beyond the initial application.

In distributed systems, idempotency is the ultimate defense against network instability. Because networks drop packets, clients will inevitably retry requests. If your server is not idempotent, a retried request might result in charging a customer's credit card twice, or sending an email twice. Idempotency guarantees that no matter how many times a request is duplicated, the business outcome remains safe and identical to a single execution.

---

# Learning Objectives

After completing this chapter you will be able to:

* Identify which HTTP methods are naturally idempotent and which are not.
* Understand the deadly combination of Network Timeouts and Retries.
* Implement an Idempotency Key system using Go and Redis.
* Safely handle duplicated messages from Message Queues.

---

# Prerequisites

Before reading this chapter you should know:

* Retries and Exponential Backoff (`08-Retries-and-Exponential-Backoff.md`).
* Message Queues (`06-Message-Queues.md`).

---

# Why This Topic Exists

A user clicks the "Pay $100" button on your website. 
The Mobile App sends a POST request to your Go backend. 
Your Go backend talks to Stripe, successfully charges the card $100, and saves the receipt in the database.

Now, your Go backend sends the HTTP 200 OK response back to the Mobile App. 
**But right at that millisecond, the user's cell connection drops.** 
The Mobile App never receives the response. It assumes the request failed. Because the app was programmed well, it automatically executes a **Retry**. 

Your Go backend receives a brand new POST request to "Pay $100". Without idempotency, you will charge the user another $100. The user sues you.

---

# Real-World Analogy

### The Elevator Button

* **Not Idempotent (The Volume Dial)**: If you turn a volume dial up, and then turn it up again, the music gets louder and louder. Every action compounds.
* **Idempotent (The Elevator Button)**: You press the button to call the elevator to floor 5. The light turns on. You are impatient, so you press the button 10 more times. The elevator doesn't come 10 times. It doesn't go to floor 50. Pressing the button once has the exact same effect as pressing it 11 times. The operation is idempotent.

---

# Core Concepts

### HTTP Methods and Idempotency
By design, HTTP defines the semantics of idempotency for its verbs:
* **GET**: Idempotent. Fetching data 10 times doesn't change the database.
* **PUT**: Idempotent. Updating a user's name to "Alice" 10 times results in the name being "Alice".
* **DELETE**: Idempotent. Deleting User 5 multiple times results in User 5 being gone (though subsequent calls might return 404, the *state of the system* remains the same).
* **POST**: **NOT IDEMPOTENT**. Sending a POST request to `/orders` 10 times will create 10 distinct orders.

### Idempotency Keys
To make POST requests idempotent, clients generate a unique string (UUID) called an **Idempotency Key** and include it in the request header. The server stores this key. If the server sees a request with a key it has already processed, it skips the business logic and simply returns the cached success response.

---

# Architecture Diagram

```mermaid
flowchart TD
    Client[Mobile App]
    API[Go Backend]
    DB[(Redis Cache)]
    Stripe[Stripe API]

    Client -- "1. POST /charge<br/>Header: Idempotency-Key: X" --> API
    
    API -- "2. Check Key 'X'" --> DB
    DB -- "Not Found" --> API
    
    API -- "3. Charge Card" --> Stripe
    Stripe -- "Success" --> API
    
    API -- "4. Save Key 'X' -> Success" --> DB
    API -.-x |"5. Network Drops"| Client
    
    note left of Client: Timeout! Retrying...
    
    Client -- "6. POST /charge<br/>Header: Idempotency-Key: X" --> API
    API -- "7. Check Key 'X'" --> DB
    DB -- "Found! Return cached response" --> API
    
    API -- "8. Return Success (NO STRIPE CALL!)" --> Client
```

---

# Step-by-Step Implementation

1. The Client generates a unique UUID (Idempotency Key) for the specific user action.
2. The Client sends the request, including `Idempotency-Key: <UUID>` in the HTTP Headers.
3. The Server extracts the header.
4. The Server queries a fast database (like Redis or Postgres) to check if the key exists.
5. **If it exists**: The server immediately returns the cached HTTP response stored with that key.
6. **If it does not exist**: The server processes the payment/order.
7. Before returning the response to the client, the server saves the Idempotency Key and the HTTP Response payload to the database.

---

# Beginner Example

A naturally idempotent function (PUT/Upsert). You don't always need a caching layer if your database logic is designed to be idempotent.

```go
package main

import "fmt"

var database = make(map[string]int)

// NOT IDEMPOTENT (Like a standard POST)
func AddPointsBad(userID string, amount int) {
	database[userID] += amount // Compounding mutation!
	fmt.Printf("User %s has %d points.\n", userID, database[userID])
}

// IDEMPOTENT (Like a PUT)
func SetPointsGood(userID string, targetAmount int) {
	database[userID] = targetAmount // Absolute mutation!
	fmt.Printf("User %s has %d points.\n", userID, database[userID])
}

func main() {
	fmt.Println("--- BAD (Retries cause data corruption) ---")
	AddPointsBad("alice", 10) // Expected: 10
	AddPointsBad("alice", 10) // Retry! Now Alice has 20. Bug!

	fmt.Println("\n--- GOOD (Retries are safe) ---")
	SetPointsGood("bob", 10) // Expected: 10
	SetPointsGood("bob", 10) // Retry! Bob still has 10. Safe!
}
```

---

# Intermediate Example (Idempotency Keys)

Handling non-idempotent operations (like charging a credit card) using a Cache.

```go
package main

import (
	"fmt"
	"net/http"
)

// Simulate Redis Database
var idempotencyCache = make(map[string]string)

func ChargeCreditCard(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the Idempotency Key
	key := r.Header.Get("Idempotency-Key")
	if key == "" {
		http.Error(w, "Idempotency-Key header required", 400)
		return
	}

	// 2. Check the Cache
	if cachedResponse, exists := idempotencyCache[key]; exists {
		fmt.Println("[SERVER] Duplicate request detected! Returning cached response.")
		w.Write([]byte(cachedResponse))
		return
	}

	// 3. Process the actual heavy/sensitive work
	fmt.Println("[SERVER] Processing payment via Stripe...")
	// ... Stripe API Call ...
	successResponse := "Payment of $100 successful. Receipt ID: 999"

	// 4. Save the response to the cache for future retries
	idempotencyCache[key] = successResponse

	// 5. Return response
	w.Write([]byte(successResponse))
}

func main() {
	// Simulating the Client
	req, _ := http.NewRequest("POST", "/charge", nil)
	
	// The client generates a unique UUID for this specific checkout attempt
	uniqueCheckoutID := "uuid-1234-5678" 
	req.Header.Set("Idempotency-Key", uniqueCheckoutID)

	fmt.Println("--- Client Attempt 1 ---")
	ChargeCreditCard(nil, req) // Uses dummy ResponseWriter for example

	fmt.Println("\n--- Client Attempt 2 (Network dropped, Retrying!) ---")
	ChargeCreditCard(nil, req) 
}
```

*Output:*
```text
--- Client Attempt 1 ---
[SERVER] Processing payment via Stripe...

--- Client Attempt 2 (Network dropped, Retrying!) ---
[SERVER] Duplicate request detected! Returning cached response.
```
Notice how Stripe was only called once, despite the client retrying the exact same request.

---

# Production Use Cases

### 1. Stripe API
The Stripe API is famous for popularizing the `Idempotency-Key` header. If you are integrating with Stripe, you must always provide this header. If your Go server crashes while waiting for Stripe to reply, you can safely restart your Go server and send the exact same request to Stripe with the exact same Idempotency-Key. Stripe will recognize it and simply hand you back the receipt from the first successful charge.

### 2. Message Queue Consumers (RabbitMQ / Kafka)
Message brokers guarantee "At-Least-Once" delivery. This means because of network jitter, a broker might deliver the same message to your Go worker twice. Your worker MUST be idempotent. It must extract the `MessageID`, check the database if it has already processed that `MessageID`, and if so, safely ignore the duplicate message.

---

# Best Practices

* **Scope the Key**: An idempotency key should usually be scoped to a specific User ID. If two different users happen to randomly generate the same UUID (astronomically rare, but possible), User B shouldn't receive User A's payment receipt. Store the cache key as `userID + idempotencyKey`.
* **Expire the Keys**: Don't store idempotency keys in Redis forever. Usually, network retries happen within minutes. Set a TTL (Time To Live) on the Redis key for 24 hours. After that, it gets automatically deleted to save RAM.
* **Concurrency Protection**: If a client has a bug and fires two identical requests at the *exact same millisecond*, both might check the cache, see it's empty, and both charge the card! You must use Redis `SETNX` (Set if Not Exists) or a Database Unique Constraint to ensure atomic locking during the check phase.

---

# Common Mistakes

### Using the Request Payload as the Key
Do not hash the JSON payload and use it as the idempotency key! If a user intentionally wants to buy two identical cups of coffee in a row, the payloads are identical. The hashing approach would block the second legitimate purchase. The Client must explicitly generate a new, unique UUID for the second purchase.

---

# Debugging Guide

* **"Database Unique Constraint Violation"**: If you are inserting a row with a Unique constraint, and a retry occurs, the database will throw an error on the second insert. This is actually a crude but effective form of idempotency! Catch the unique constraint error, realize it's a retry, and return a 200 OK to the client instead of a 500 Server Error.

---

# Quiz

## Multiple Choice Questions
**1. Which HTTP verb requires the use of an Idempotency-Key header to safely support retries in an API?**
A) GET
B) PUT
C) POST
*Answer*: C. POST operations are typically non-idempotent (e.g., creating a new record, charging a card). PUT is naturally idempotent (replacing a record).

## True or False
**If a client needs to retry a request, they should generate a brand new Idempotency-Key for the retry attempt to ensure it goes through.**
*Answer*: False. That defeats the entire purpose! The client MUST send the exact same Idempotency-Key on all retries of the *same* operation. They only generate a new key when the user initiates a *new*, separate action.

---

# Interview Questions

## Beginner
**Q**: What does it mean for an API endpoint to be "Idempotent"?
*Answer*: It means that calling the endpoint once has the exact same effect on the system as calling it 100 times sequentially.

## Intermediate
**Q**: Why is Idempotency critical when working with Message Queues like RabbitMQ?
*Answer*: Message Queues generally provide "At-Least-Once" delivery guarantees. Due to network failures or consumer crashes before sending an ACK, the broker may re-deliver a message that has already been processed. If the consumer is not idempotent, processing the duplicate message will cause data corruption (like deducting inventory twice for the same order).

## Advanced
**Q**: Explain how to prevent a race condition when implementing Idempotency Keys in a highly concurrent environment.
*Answer*: A naive "Check Cache -> Process -> Save Cache" approach suffers from race conditions. If two identical requests hit the server simultaneously, both will check the cache, find it empty, and process the heavy workload. To prevent this, you must use atomic locks. In Redis, you use the `SETNX` command. In a SQL database, you insert the Idempotency Key into a table with a Unique constraint. The first request grabs the lock, and the second request immediately fails or blocks until the first finishes and populates the cache.

---

# Summary

Idempotency is the final puzzle piece that allows you to safely use Timeouts and Retries. Without idempotency, retries are a game of Russian Roulette. By enforcing Idempotency Keys across your APIs and Message Queues, you guarantee that your system's state remains perfectly consistent, no matter how chaotic the network becomes.

---

# Key Takeaways

* ✔ Idempotency means multiple identical requests yield the same result.
* ✔ GET, PUT, and DELETE are naturally idempotent.
* ✔ POST is not. Protect it with Idempotency Keys.
* ✔ Crucial for safely handling automatic client retries.
* ✔ Essential for protecting Message Queue workers from duplicate events.

---

# Further Reading
* [Stripe API Design: Designing Robust and Predictable APIs with Idempotency](https://stripe.com/blog/idempotency)

---

# Next Chapter
➡️ **Next:** `12-Consensus-Algorithms.md`
