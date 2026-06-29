# sync.Map

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

We know that standard Go maps (`map[string]int`) are fundamentally **not thread-safe**. If Goroutine A writes to a map while Goroutine B reads from it, the Go runtime instantly crashes the application with a `fatal error: concurrent map read and map write`.

While wrapping a standard map in a `sync.RWMutex` (Chapter 22) is the standard solution, Go provides a specialized, highly optimized alternative: `sync.Map`. This data structure is designed for very specific, extreme-concurrency use cases.

---

# Learning Objectives

After completing this chapter you will be able to:

* Explain when to use `sync.Map` vs `map + RWMutex`.
* Use the core methods: `Load`, `Store`, `LoadOrStore`, and `Delete`.
* Iterate over a `sync.Map` safely using `Range`.
* Understand the internal "read" and "dirty" map architecture.

---

# Prerequisites

Before reading this chapter you should know:

* `sync.RWMutex` (`22-RWMutex.md`)
* `sync/atomic` (`23-Atomic.md`)
* Empty Interfaces (`any` / `interface{}`)

---

# Why This Topic Exists

If you have a map wrapped in an `RWMutex` running on a massive 64-core server, and all 64 CPU cores are constantly doing `mu.RLock()` to read the same map, the hardware memory bus becomes saturated updating the internal Mutex reader-count. This causes a phenomenon called **Cache Contention**.

To solve this for Google-scale applications, the Go team created `sync.Map`. It uses `atomic` operations and a clever two-map system to allow lock-free reads that scale linearly across infinite CPU cores.

---

# Real-World Analogy

### The Library Index Card System

* **Standard Map + Mutex**: There is one master index book. When a librarian wants to add a book, they lock the book. When a student wants to read the book, they lock it. If 100 students want to read, they bump elbows.
* **sync.Map (Two Maps)**: There is a **Read-Only Display Board** (lock-free) and a **Private Desk Ledger** (locked). 
  - When students want to find a book, they look at the Display Board. No locks needed!
  - When a librarian adds a new book, they lock the Private Ledger, write it down, and occasionally take the ledger and photocopy it to become the new Display Board (promoted).

---

# Core Concepts

* **Thread-Safe**: Safe for concurrent use by multiple Goroutines without external locks.
* **Type Unsafe**: Internally uses `any` (empty interface), meaning you must type-assert values when you load them (`val.(string)`).
* **LoadOrStore**: A specialized atomic operation that fetches a key, but if it doesn't exist, inserts a default value safely without race conditions.

---

# Internal Runtime Explanation

Internally, `sync.Map` maintains TWO maps:
1. `read`: An `atomic.Value` containing a read-only map. Lookups here are 100% lock-free and blazing fast.
2. `dirty`: A standard map protected by a Mutex. 

**Reads**: It first checks the `read` map atomically. If the key isn't there, it locks the Mutex and checks the `dirty` map. If it finds it in `dirty`, it increments a "miss" counter.
**Writes**: It locks the Mutex and writes to the `dirty` map.
**Promotion**: When the "miss" counter gets too high, it takes the `dirty` map, atomically replaces the `read` map with it, and sets `dirty` to nil. This is called a "promotion".

---

# Memory Layout

```text
Heap Memory

+-------------------------------------------------+
| sync.Map                                        |
|                                                 |
| read (atomic.Value) ---> [ Map: A=1, B=2 ]      | <-- Lock-Free Readers
|                                                 |
| mu (sync.Mutex)                                 | <-- Lock for Writers
| dirty (map) -----------> [ Map: A=1, B=2, C=3 ] | 
+-------------------------------------------------+
```

---

# Architecture Diagram

```mermaid
flowchart TD
    Req[Goroutine calls Load(K)]
    ReadMap[Check atomic 'read' map]
    FoundRead{Found?}
    
    Mutex[Acquire Mutex]
    DirtyMap[Check 'dirty' map]
    Miss[Increment Misses]
    Promote[Promote 'dirty' to 'read']
    
    Req --> ReadMap
    ReadMap --> FoundRead
    FoundRead -->|Yes| Return[Return Value]
    FoundRead -->|No| Mutex
    
    Mutex --> DirtyMap
    DirtyMap --> Miss
    Miss --> Promote
    Promote --> Return
```

---

# Step-by-Step Execution

1. `var sm sync.Map`.
2. `sm.Store("key", "val")`: Acquires lock, writes to `dirty` map.
3. `sm.Load("key")`: Checks `read` map (miss). Acquires lock, checks `dirty`, increments miss counter.
4. After many misses, `dirty` map is promoted to `read` map.
5. `sm.Load("key")`: Checks `read` map. Hits instantly. Lock-free!

---

# Syntax

```go
import "sync"

var sm sync.Map

// Store a value
sm.Store("name", "GoVerse")

// Load a value (Returns value, ok)
val, ok := sm.Load("name")

// Delete a value
sm.Delete("name")
```

---

# Beginner Example

Basic thread-safe storage and type assertion.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map

	// Store (Accepts any, any)
	sm.Store("port", 8080)
	sm.Store("host", "localhost")

	// Load
	if val, ok := sm.Load("port"); ok {
		// You MUST type-assert the value because it is returned as `any`
		port := val.(int)
		fmt.Printf("Running on port %d\n", port)
	}

	// Delete
	sm.Delete("host")
}
```

---

# Intermediate Example

Using `LoadOrStore`. This is incredibly useful for caching, where you want to fetch an item, but if it doesn't exist, calculate it and store it, ensuring two Goroutines don't calculate it twice.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map
	var wg sync.WaitGroup

	// 5 Goroutines trying to initialize the same cache key simultaneously
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			
			// LoadOrStore atomically checks if "token" exists. 
			// If not, it stores "abc-123".
			// 'loaded' is true if the value was already in the map.
			actual, loaded := sm.LoadOrStore("token", "abc-123")
			
			if loaded {
				fmt.Printf("G%d: Token already existed. Read: %v\n", id, actual)
			} else {
				fmt.Printf("G%d: I was the one who stored the token! %v\n", id, actual)
			}
		}(i)
	}

	wg.Wait()
}
// OUTPUT: Only one Goroutine will print "I was the one who stored..."
```

---

# Advanced Example

Iterating over a `sync.Map` using `Range`. Because standard `for k, v := range map` does not work on a `sync.Map`, you must use a callback function.

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var sm sync.Map
	sm.Store("A", 1)
	sm.Store("B", 2)
	sm.Store("C", 3)

	// Range takes a callback function. 
	// It executes this function for every key-value pair in the map.
	// If the callback returns false, the iteration stops (breaks).
	sm.Range(func(key, value any) bool {
		k := key.(string)
		v := value.(int)
		
		fmt.Printf("Key: %s, Value: %d\n", k, v)
		
		// Return true to continue to the next item
		return true 
	})
}
```

---

# Production Use Cases

### 1. Append-Only Caches
When building a global cache of IP addresses blocklists. The list is updated once an hour but read millions of times a minute by incoming HTTP requests. Because `sync.Map` is lock-free for existing keys, it will not bottleneck the server's CPU cores.

### 2. Connection Pools
Tracking active WebSocket connections. `LoadOrStore` is used to ensure that if a user connects via two tabs simultaneously, only one connection object is initialized in the map.

---

# Performance Analysis

* **The Golden Rule**: Use `sync.Map` ONLY when you have keys that are written once and read many times (Append-Only), OR when multiple Goroutines read, write, and overwrite entries for disjoint (different) keys.
* **The Penalty**: If you constantly insert *new* keys into a `sync.Map`, it performs terribly. Every new key triggers a Mutex lock on the `dirty` map, and constant inserts prevent the `dirty` map from ever being promoted efficiently. A standard `map + Mutex` is much faster for write-heavy workloads.

---

# Best Practices

* **Default to `map + RWMutex`**: The Go authors explicitly state that standard maps with Mutexes should be your default choice. They provide type-safety (no `any` casting) and predictable performance. Use `sync.Map` only when profiling proves Cache Contention is occurring.
* **Wrap the sync.Map**: To avoid littering your codebase with `val.(string)` type assertions, wrap `sync.Map` in a struct with strongly typed methods.

---

# Common Mistakes

### Using Range for Snapshots
```go
// BAD: Assuming Range is atomic across the whole map.
sm.Range(func(key, value any) bool {
    // If another Goroutine deletes a key while Range is running, 
    // you might not see it, or you might see it. 
    // Range does NOT freeze the map in time!
    return true
})
```

---

# Debugging Guide

* **Panics on Assertion**: The biggest runtime risk with `sync.Map` is storing an `int` and accidentally type-asserting it as a `string` (`val.(string)`). This will panic instantly. Always write strongly-typed wrappers!

---

# Exercises

## Beginner
Create a `sync.Map`. Store three string keys with boolean values. Use `Load()` to retrieve one and print it, ensuring you type-assert the boolean correctly.

## Intermediate
Write a struct `type UserCache struct { internal sync.Map }`. Write two methods: `AddUser(id string, name string)` and `GetUser(id string) string`. This strongly-typed wrapper hides the empty interfaces from the rest of your application.

---

# Quiz

## Multiple Choice Questions
**1. Which internal map does `sync.Map` check first during a `Load()` operation?**
A) The dirty map
B) The read map
C) The mutex map
*Answer*: B

## True or False
**`sync.Map` is faster than `map + Mutex` in all scenarios.**
*Answer*: False. If you are constantly inserting new keys, `sync.Map` is significantly slower due to dirty map promotion overhead. It is only faster for append-once, read-many workloads.

---

# Interview Questions

## Beginner
**Q**: Why can't you use a normal `for k,v := range myMap` loop on a `sync.Map`?
*Answer*: Because `sync.Map` is a struct, not a built-in map type. You must use the `.Range()` method and pass a callback function to iterate over its elements.

## Intermediate
**Q**: Explain the utility of `LoadOrStore`.
*Answer*: It combines reading a key and writing a default value into a single, atomic operation. It prevents a race condition where two Goroutines check if a key exists, both see it missing, and both attempt to perform the heavy initialization to store it.

## Google-Level Questions
**Q**: Explain the "Amortized Constant Time" performance of `sync.Map` reads. How does the dirty map promotion work to achieve lock-free reads?
*Answer*: When a key is requested that isn't in the atomic `read` map, the `sync.Map` acquires a Mutex and checks the `dirty` map, incrementing a miss counter. When `misses == len(dirty)`, the `dirty` map is atomically swapped to become the new `read` map. From that point on, all future reads for those keys hit the atomic `read` map instantly without acquiring the Mutex. Thus, the cost of the Mutex locks is "amortized" (spread out) over the billions of lock-free reads that follow the promotion.

---

# Mini Project

**Requirement**: The Active User Tracker
Use a `sync.Map` to track active users on a server.
1. When a user connects, `Store` their UserID (string) and a LastSeen timestamp (`time.Time`).
2. Write a background Goroutine that runs a `time.Ticker` every 5 seconds.
3. When the ticker fires, use `sm.Range()` to iterate over all users. If their LastSeen time is older than 10 seconds, `Delete` them from the map and print "User Disconnected".

---

# Cheat Sheet

* **Store**: `sm.Store(k, v)`
* **Load**: `v, ok := sm.Load(k)`
* **Delete**: `sm.Delete(k)`
* **Iterate**: `sm.Range(func(k, v any) bool { return true })`
* **Atomic Insert**: `actual, loaded := sm.LoadOrStore(k, defaultV)`

---

# Summary

`sync.Map` is a hyper-optimized data structure built to solve Mutex contention on massively multi-core systems. While you shouldn't use it everywhere, understanding its internal `read`/`dirty` promotion mechanics provides profound insight into how the Go authors approach lock-free scaling.

---

# Key Takeaways

* ✔ Use for append-once, read-many workloads.
* ✔ Do not use for write-heavy continuous inserts.
* ✔ Requires type assertions (`any`).
* ✔ `LoadOrStore` prevents concurrent duplicate initializations.

---

# Further Reading
* [Go documentation for sync.Map](https://pkg.go.dev/sync#Map)

---

# Next Chapter
➡️ **Next:** `26-sync.Cond.md`
