# Singleton Pattern

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
* Best Practices (And Anti-Patterns)
* Common Mistakes
* Debugging Guide
* Exercises
* Quiz
* Interview Questions
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

The **Singleton Pattern** is a Creational Design Pattern that ensures a class has only one instance, while providing a global access point to this instance.

In Go, we don't have classes, but we can ensure a specific `struct` is only instantiated exactly once across the entire lifecycle of the application. Go provides an incredibly elegant, thread-safe mechanism to achieve this out of the box: the `sync.Once` primitive.

*Warning: The Singleton is widely considered an anti-pattern in modern software engineering because it introduces global state. We will explore how to implement it safely, and when to avoid it entirely.*

---

# Learning Objectives

After completing this chapter you will be able to:

* Implement a thread-safe Singleton using `sync.Once`.
* Understand the difference between `init()` initialization and Lazy Initialization.
* Identify why global state ruins unit testing.
* Refactor a Singleton into Dependency Injection.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Pointers.
* Goroutines and Race Conditions (`28-Race-Conditions.md`).
* The `sync` package.

---

# Why This Topic Exists

Imagine your application needs to load a massive, 500MB configuration file from disk. You have 100 different Goroutines handling HTTP requests, and they all need to read this configuration.

If every Goroutine calls `LoadConfig()`, your server will try to load 500MB of data 100 times, run out of memory, and crash. You need a way to guarantee that `LoadConfig()` is executed *exactly once*, and that all 100 Goroutines share a pointer to that single, globally available configuration object.

---

# Real-World Analogy

### The National Constitution

* **The Object**: The original, signed Constitution of a country.
* **The Problem**: If a judge needs to reference the law, they shouldn't write a brand new Constitution. There can only be one.
* **The Singleton**: The government places the original document in a highly secure, global archive. Any judge in the country can request access to *view* it (global access point), but no one is allowed to create a second one.

---

# Core Concepts

* **Global Variable**: A variable declared at the package level that holds the single instance.
* **Unexported Constructor**: Preventing users from calling `new(Struct)` directly.
* **Lazy Initialization**: Delaying the creation of the Singleton until the very first time it is requested.
* **Thread-Safety**: Ensuring that if 1,000 Goroutines request the Singleton at the exact same millisecond, only 1 instance is created.
* **`sync.Once`**: Go's built-in tool that guarantees a function is executed exactly one time, regardless of how many Goroutines call it simultaneously.

---

# Architecture Diagram

```mermaid
flowchart TD
    G1[Goroutine 1]
    G2[Goroutine 2]
    G3[Goroutine 3]
    
    GetInst[GetInstance() Function]
    SyncOnce{sync.Once.Do}
    Init[Create Instance]
    Return[Return Global Pointer]
    
    G1 --> GetInst
    G2 --> GetInst
    G3 --> GetInst
    
    GetInst --> SyncOnce
    
    SyncOnce -- "First Caller Only" --> Init
    Init --> Return
    
    SyncOnce -- "Subsequent Callers" --> Return
```

---

# Step-by-Step Implementation

1. Declare a private global variable to hold the pointer to your struct: `var instance *Config`.
2. Declare a private global `sync.Once` variable: `var once sync.Once`.
3. Create an exported function `GetInstance() *Config`.
4. Inside `GetInstance`, call `once.Do(func() { ... })`.
5. Inside the closure, instantiate the struct and assign it to the global variable.
6. Return the global variable.

---

# Syntax

```go
var (
    instance *Database
    once     sync.Once
)

func GetInstance() *Database {
    // This closure will only execute ONCE, even under heavy concurrency.
    once.Do(func() {
        instance = &Database{}
    })
    return instance
}
```

---

# Beginner Example

A simple thread-safe Configuration loader.

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

type Config struct {
	LogLevel string
}

var (
	configInstance *Config
	once           sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		fmt.Println("--- LOADING CONFIGURATION FROM DISK ---")
		time.Sleep(1 * time.Second) // Simulate expensive loading
		configInstance = &Config{LogLevel: "DEBUG"}
	})
	return configInstance
}

func main() {
	var wg sync.WaitGroup

	// 10 Goroutines asking for the config at the exact same time
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := GetConfig()
			fmt.Printf("Worker %d got LogLevel: %s\n", id, cfg.LogLevel)
		}(i)
	}

	wg.Wait()
}
```
*Output: Notice that "LOADING CONFIGURATION" prints exactly once, and all 10 workers instantly receive the shared pointer.*

---

# Intermediate Example

The `init()` alternative (Eager Initialization).
If your Singleton doesn't require complex logic or error handling, you can use Go's built-in `init()` function. `init()` is guaranteed by the Go runtime to run exactly once, synchronously, before `main()` starts.

```go
package main

import "fmt"

type GlobalCache struct {
	store map[string]string
}

// 1. Declare the unexported global variable
var cacheInstance *GlobalCache

// 2. The runtime calls init() before main()
func init() {
	fmt.Println("Eagerly initializing the cache...")
	cacheInstance = &GlobalCache{
		store: make(map[string]string),
	}
}

// 3. Simple getter
func GetCache() *GlobalCache {
	return cacheInstance
}

func main() {
	fmt.Println("Main started.")
	
	c := GetCache()
	c.store["user1"] = "Alice"
	
	fmt.Println("Cache length:", len(c.store))
}
```

---

# Advanced Example

Handling Errors in a Singleton initialization.
`sync.Once.Do` does not return an error. If the initialization fails (e.g., the DB is unreachable), and you use a standard `sync.Once`, subsequent calls will return a `nil` pointer and panic! 
To handle errors gracefully, we must track the error state.

*(Note: Go 1.21 introduced `sync.OnceValues` which simplifies this, but we will look at the manual implementation to understand the mechanics).*

```go
package main

import (
	"errors"
	"fmt"
	"sync"
)

type DBClient struct{}

var (
	dbInstance *DBClient
	dbError    error
	once       sync.Once
)

func GetDB() (*DBClient, error) {
	once.Do(func() {
		fmt.Println("Attempting to connect to DB...")
		// Simulate a connection failure
		success := false 
		
		if !success {
			dbError = errors.New("connection refused")
			return
		}
		dbInstance = &DBClient{}
	})

	// Return both the instance and the cached error
	return dbInstance, dbError
}

func main() {
	// First attempt triggers the connection
	db1, err1 := GetDB()
	fmt.Printf("Attempt 1: DB=%v, Err=%v\n", db1 != nil, err1)

	// Second attempt instantly returns the cached error
	db2, err2 := GetDB()
	fmt.Printf("Attempt 2: DB=%v, Err=%v\n", db2 != nil, err2)
}
```

---

# Production Use Cases

### 1. In-Memory Caches
A global Redis client or a local LRU cache usually needs to be accessible from multiple decoupled packages (e.g., the `auth` package and the `billing` package). A Singleton provides an easy global access point.

### 2. Standard Logger
Go's built-in `log` package is technically a Singleton! It instantiates a standard `*log.Logger` under the hood, allowing you to call `log.Println()` from anywhere in your codebase without passing a logger object around.

---

# Performance Analysis

`sync.Once` is heavily optimized. It uses lock-free atomic reads (`sync/atomic`) on the fast path to check if the function has already been executed. If it has, the overhead is roughly 1 nanosecond. It only acquires a Mutex lock on the slow path (the very first execution). It is perfectly safe for ultra-high-throughput systems.

---

# Best Practices (And Anti-Patterns)

* **Prefer Dependency Injection**: The Singleton is an anti-pattern because it creates **Global State**. If `ServiceA` calls `GetDB()`, it is tightly coupled to the concrete database. It is incredibly difficult to mock global variables in Unit Tests. 
  * *Fix*: Instead of `GetDB()`, instantiate the DB once in `main()`, and pass the pointer (Dependency Injection) to `ServiceA` as an interface.
* **Keep Singletons Stateless (or Thread-Safe)**: If your Singleton is a configuration struct that is read-only, it is perfectly safe. If your Singleton has a `map` that is being written to, you MUST protect that map with a `sync.RWMutex`, otherwise 100 Goroutines will create a fatal data race.

---

# Common Mistakes

### The Double-Checked Locking (DCL) Anti-Pattern
Developers coming from Java often try to implement the DCL pattern manually to avoid locks:
```go
// BAD: Manual Double-Checked Locking in Go is prone to data races!
if instance == nil {
    mu.Lock()
    defer mu.Unlock()
    if instance == nil {
        instance = &Config{}
    }
}
```
*Fix: Never do this in Go. Just use `sync.Once`. It implements DCL safely and efficiently under the hood at the assembly level.*

---

# Debugging Guide

* **"panic: assignment to entry in nil map"**: You successfully retrieved the Singleton pointer, but you forgot to initialize a map or slice *inside* the Singleton struct during the `once.Do` block.
* **Flaky Unit Tests**: If Test A modifies the Singleton, and Test B reads the Singleton, Test B will randomly fail depending on execution order. This is why Singletons are dangerous. *Fix: Reset global variables in a `setup/teardown` block in your tests, or better yet, refactor to Dependency Injection.*

---

# Exercises

## Beginner
Create a Singleton `Logger` struct using `sync.Once`. It should have a method `Log(msg string)`. Write a loop that spawns 50 Goroutines, all calling `GetLogger().Log("hello")`. Prove that the `Logger` is only instantiated once.

## Intermediate
Refactor a Singleton into Dependency Injection.
Take this code: `func HandleUser() { db := GetDB(); db.Save() }`
Change it so `HandleUser` accepts an interface `Saver` as an argument, completely removing the need for `HandleUser` to know about the Singleton.

---

# Quiz

## Multiple Choice Questions
**1. Why is `sync.Once` preferred over a standard `sync.Mutex` for Singletons in Go?**
A) `sync.Once` is just an alias for `sync.Mutex`.
B) A standard Mutex forces all callers to lock and unlock the Mutex every single time they fetch the instance, killing performance. `sync.Once` uses atomic reads to bypass the lock entirely after the first execution.
C) `sync.Once` prevents deadlocks.
*Answer*: B

## True or False
**Global singletons make unit testing easier because you don't have to pass mock objects into functions.**
*Answer*: False. Global singletons make testing a nightmare because tests running in parallel will overwrite each other's global state (test pollution).

---

# Interview Questions

## Beginner
**Q**: How do you implement a thread-safe Singleton in Go?
*Answer*: By declaring a package-level variable for the instance, and using a `sync.Once.Do()` block inside a getter function to instantiate it exactly once.

## Intermediate
**Q**: What is the difference between Lazy Initialization (using `sync.Once`) and Eager Initialization (using `init()`)?
*Answer*: Eager initialization (`init()`) runs at startup before the application begins, regardless of whether the object is ever used. Lazy initialization (`sync.Once`) delays the creation of the object until the exact moment `GetInstance()` is called for the first time, saving memory if the object is never requested.

## Advanced
**Q**: Why is the Singleton considered an anti-pattern in modern software engineering, and how does Dependency Injection solve it?
*Answer*: Singletons introduce hidden, global state that tightly couples code and ruins unit test isolation (test pollution). Dependency Injection solves this by moving the responsibility of instantiation to the top level (`main()`). The instance is created once, and then passed explicitly (injected) down the call stack as an interface. This makes the code decoupled, easily mockable, and removes hidden global state.

---

# Cheat Sheet

* **Thread-Safe Singleton**:
```go
var (
    instance *MyStruct
    once     sync.Once
)
func GetInstance() *MyStruct {
    once.Do(func() { instance = &MyStruct{} })
    return instance
}
```
* **Eager Singleton**:
```go
var instance *MyStruct
func init() { instance = &MyStruct{} }
```

---

# Summary

The Singleton Pattern guarantees a single instance of an object exists. Go's `sync.Once` provides the most elegant, thread-safe implementation of this pattern of any language. However, with great power comes great responsibility: use Singletons sparingly for stateless configurations or caches, and avoid them entirely for core business logic to preserve the testability of your codebase.

---

# Key Takeaways

* ✔ Use `sync.Once` for thread-safe Lazy Initialization.
* ✔ Use `init()` for Eager Initialization.
* ✔ Singletons introduce global state, making unit testing difficult.
* ✔ Prefer Dependency Injection over Singletons whenever possible.

---

# Further Reading
* [Go documentation for sync.Once](https://pkg.go.dev/sync#Once)
* [Refactoring.guru: Singleton](https://refactoring.guru/design-patterns/singleton)

---

# Next Chapter
➡️ **Next:** `08-Object-Pool.md`
