# Proxy Pattern

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
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

The **Proxy Pattern** is a Structural Design Pattern that provides a surrogate or placeholder for another object to control access to it.

A Proxy sits exactly in the middle between the Client and the Real Object. Because the Proxy implements the exact same interface as the Real Object, the Client has absolutely no idea it is talking to a Proxy. This allows the Proxy to intercept calls and perform actions—like caching, lazy loading, or access control—*before* or *after* passing the request to the real object.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand the difference between a Proxy and a Decorator.
* Implement a **Caching Proxy** to dramatically reduce database load.
* Implement a **Virtual Proxy** to delay the expensive initialization of an object (Lazy Loading).
* Implement a **Protection Proxy** to enforce Role-Based Access Control (RBAC).

---

# Prerequisites

Before reading this chapter you should know:

* Interfaces (`02-Accept-Interfaces-Return-Structs.md`).
* Struct Embedding.

---

# Why This Topic Exists

Imagine your Go server has a function that queries a massive Database for a User Profile. This query takes 2 seconds and uses a lot of CPU. 
If 1,000 users request that profile simultaneously, your database will crash.

You could modify the Database code to add Redis caching. But modifying core, tested code violates the Open-Closed Principle. 

Instead, you write a `ProxyDB` struct that implements the same `Database` interface. You give the Proxy to the HTTP handler. When the handler asks for the profile, the Proxy checks its cache. If the data is there, it returns it instantly (0 seconds). If not, the Proxy forwards the request to the Real Database, caches the result, and returns it. You just saved your database without changing a single line of the original database code!

---

# Real-World Analogy

### The Credit Card

* **The Real Object**: A stack of physical cash in your bank vault.
* **The Problem**: Carrying heavy stacks of cash is annoying and dangerous.
* **The Proxy**: The Credit Card. It implements the exact same interface as cash (it can be used to buy things). 
* **The Control**: When you swipe the card, the Proxy intercepts the purchase. It checks your balance (Access Control). If approved, it communicates with the bank (Remote Proxy) to transfer the money. The cashier (Client) doesn't care if you use cash or a card; both satisfy the "Payment" interface.

---

# Core Concepts

* **Subject Interface**: The interface defining the methods both the Proxy and Real Object must implement.
* **Real Subject**: The heavy, slow, or sensitive object that does the actual work.
* **Proxy**: The wrapper object that holds a reference to the Real Subject and controls access to it.

### Types of Proxies:
1. **Virtual Proxy (Lazy Initialization)**: Delays creating the Real Subject until the first time it's actually needed.
2. **Protection Proxy (Auth)**: Checks permissions before allowing the call to reach the Real Subject.
3. **Caching Proxy**: Returns cached results for identical requests instead of hitting the Real Subject.
4. **Remote Proxy**: The Real Subject lives on another server. The Proxy handles the network communication (e.g., gRPC stubs).

---

# Architecture Diagram

```mermaid
flowchart LR
    Client[Client Code]
    Interface((Database Interface<br/>Query()))
    
    Proxy[Proxy Cache<br/>Query()]
    Real[Real Database<br/>Query()]
    
    Client -- "Calls Query()" --> Interface
    Proxy -. implements .-> Interface
    Real -. implements .-> Interface
    
    Client -- "Thinks it's talking to Real" --> Proxy
    
    Proxy -- "Cache Miss: Forwards Call" --> Real
    Real -- "Returns Data" --> Proxy
    Proxy -- "Cache Hit: Returns Instantly" --> Client
```

---

# Step-by-Step Implementation

1. Identify the **Subject Interface** (e.g., `type Downloader interface { Download(url string) }`).
2. Identify the **Real Subject** that implements the interface.
3. Create the **Proxy** struct. It should contain a pointer to the Real Subject.
4. Make the Proxy implement the Subject Interface.
5. Inside the Proxy's methods, add your control logic (caching, auth).
6. If the control logic passes, manually call the corresponding method on the Real Subject pointer.
7. Inject the Proxy into the Client code instead of the Real Subject.

---

# Syntax

```go
type Subject interface { Do() string }

type RealSubject struct{}
func (r *RealSubject) Do() string { return "Real Data" }

type Proxy struct {
    real *RealSubject
}
func (p *Proxy) Do() string {
    // 1. Intercept! Do proxy logic (auth, cache, logging)
    // 2. Forward to real subject
    return p.real.Do() 
}
```

---

# Beginner Example

A **Protection Proxy** that checks a password before allowing access to a highly classified document.

```go
package main

import "fmt"

// 1. The Interface
type Server interface {
	HandleRequest(url string, method string) (int, string)
}

// 2. The Real Subject
type Application struct{}
func (a *Application) HandleRequest(url, method string) (int, string) {
	return 200, "Sensitive application data"
}

// 3. The Proxy
type NginxProxy struct {
	app         *Application
	maxReqCount int
	rateLimiter map[string]int
}
func NewNginxProxy() *NginxProxy {
	return &NginxProxy{
		app:         &Application{},
		maxReqCount: 2,
		rateLimiter: make(map[string]int),
	}
}

// 4. Implement the Interface with Protection Logic
func (n *NginxProxy) HandleRequest(url, method string) (int, string) {
	allowed := n.checkRateLimit(url)
	if !allowed {
		return 403, "Rate Limit Exceeded"
	}
	
	// Forward to the real application
	return n.app.HandleRequest(url, method)
}

func (n *NginxProxy) checkRateLimit(url string) bool {
	if n.rateLimiter[url] >= n.maxReqCount {
		return false
	}
	n.rateLimiter[url]++
	return true
}

func main() {
	// Client only interacts with the Proxy
	proxy := NewNginxProxy()

	fmt.Println(proxy.HandleRequest("/api/status", "GET")) // 200
	fmt.Println(proxy.HandleRequest("/api/status", "GET")) // 200
	fmt.Println(proxy.HandleRequest("/api/status", "GET")) // 403! Blocked by Proxy
}
```

---

# Intermediate Example

A **Virtual Proxy** (Lazy Loading). We have a massive Machine Learning Model that takes 10 seconds and 5GB of RAM to load into memory. We only want to load it if the user *actually* requests a prediction.

```go
package main

import (
	"fmt"
	"time"
)

// The Interface
type MLModel interface {
	Predict(input string) string
}

// The Real Subject (Extremely heavy)
type DeepLearningModel struct{}
func NewDeepLearningModel() *DeepLearningModel {
	fmt.Println("--> LOADING 5GB MODEL INTO RAM (Takes 3 seconds)...")
	time.Sleep(3 * time.Second)
	return &DeepLearningModel{}
}
func (m *DeepLearningModel) Predict(input string) string {
	return "Prediction for: " + input
}

// The Virtual Proxy
type LazyModelProxy struct {
	realModel *DeepLearningModel
}

func (p *LazyModelProxy) Predict(input string) string {
	// LAZY INITIALIZATION: Only create the real model on the very first request
	if p.realModel == nil {
		fmt.Println("[Proxy] Instantiating the real model just in time...")
		p.realModel = NewDeepLearningModel()
	}
	
	// Forward the request
	return p.realModel.Predict(input)
}

func main() {
	fmt.Println("Server started.")
	
	// The Proxy instantiates instantly (0 seconds, 0MB RAM)
	var model MLModel = &LazyModelProxy{}
	
	fmt.Println("User is browsing the site...")
	time.Sleep(1 * time.Second)
	
	fmt.Println("User clicked Predict!")
	// The Proxy intercepts, loads the model, and forwards the request
	result := model.Predict("Cat Image")
	fmt.Println(result)
	
	fmt.Println("User clicked Predict again!")
	// The model is already loaded, returns instantly
	result2 := model.Predict("Dog Image")
	fmt.Println(result2)
}
```

---

# Advanced Example

A **Caching Proxy** for a Database.

```go
package main

import (
	"fmt"
	"time"
)

// Interface
type Database interface {
	GetUser(id int) string
}

// Real Subject
type MySQL struct{}
func (m *MySQL) GetUser(id int) string {
	fmt.Println("[MySQL] Executing slow query...")
	time.Sleep(1 * time.Second) // Simulate network/disk latency
	return fmt.Sprintf("User_%d_Data", id)
}

// Caching Proxy
type RedisProxy struct {
	db    *MySQL
	cache map[int]string
}
func NewRedisProxy(db *MySQL) *RedisProxy {
	return &RedisProxy{
		db:    db,
		cache: make(map[int]string),
	}
}

func (r *RedisProxy) GetUser(id int) string {
	// 1. Check Cache
	if data, exists := r.cache[id]; exists {
		fmt.Println("[Redis] Cache Hit! Returning instantly.")
		return data
	}

	// 2. Cache Miss: Hit Real Database
	data := r.db.GetUser(id)

	// 3. Save to Cache
	r.cache[id] = data
	return data
}

func main() {
	realDB := &MySQL{}
	proxyDB := NewRedisProxy(realDB)

	start := time.Now()
	fmt.Println(proxyDB.GetUser(1))
	fmt.Printf("Time taken: %v\n\n", time.Since(start)) // Takes 1s

	start2 := time.Now()
	fmt.Println(proxyDB.GetUser(1))
	fmt.Printf("Time taken: %v\n", time.Since(start2)) // Takes 0s!
}
```

---

# Production Use Cases

### 1. ORM Lazy Loading
When using a Database ORM (like GORM), if you fetch a `User` struct, the ORM might return a Proxy object for the `User.Posts` array. The posts aren't actually loaded from the DB until you explicitly loop over `User.Posts` (Virtual Proxy), saving massive amounts of memory.

### 2. gRPC Stubs
When you write a microservice that communicates with another service via gRPC, you interact with an autogenerated "Stub" object. This stub acts as a **Remote Proxy**. You call `stub.GetUser()`. The stub serializes the request, sends it over the network to the real service, waits for the response, and deserializes it. Your code thinks it called a local function.

---

# Proxy vs Decorator

They look identical in code (both wrap an object and implement the same interface). The difference is **Intent**:
* **Decorator**: Adds *new* behaviors (like logging) to the object and always passes the request to the core object.
* **Proxy**: *Controls access* to the object. It might deny the request, return a cached response, or instantiate the object itself, completely bypassing the core object.

---

# Performance Analysis

The overhead of the Proxy itself is an interface method call (nanoseconds). However, the performance *gains* provided by Virtual Proxies (saving RAM/Startup time) and Caching Proxies (saving Database I/O) are usually the most massive performance optimizations you can make in a backend system.

---

# Best Practices

* **Identical Interfaces**: The Proxy MUST implement the exact same interface as the Real Subject, otherwise it breaks the pattern. The Client should never know if they hold a Real object or a Proxy.
* **Dependency Injection**: Pass the Real Subject into the Proxy's constructor, so the Proxy is decoupled from the concrete implementation of the Real Subject.

---

# Common Mistakes

### State Desynchronization in Caching Proxies
If you implement a Caching Proxy, you must implement a mechanism to invalidate the cache if the underlying Real Database is updated. Otherwise, the Proxy will permanently return stale data.

---

# Debugging Guide

* **"Data is stale"**: Your Caching Proxy intercepted the call and returned data from its internal map, but the real database was updated by a different process. You need a Cache Invalidation strategy (e.g., TTL timeouts).

---

# Exercises

## Beginner
Write an interface `Downloader` with method `Download(url string)`. Write a Real Subject that simulates a 2-second download. Write a Caching Proxy that stores downloaded URLs in a map. If the URL is already in the map, return it instantly.

## Intermediate
Write a Protection Proxy for a `BankService`. The `Withdraw(amount int)` method should only execute on the Real Subject if the Proxy verifies that the `currentUserRole == "admin"`. If not, return an error.

---

# Quiz

## Multiple Choice Questions
**1. Which type of Proxy delays the creation of an expensive object until it is explicitly requested?**
A) Protection Proxy
B) Virtual Proxy
C) Remote Proxy
*Answer*: B

## True or False
**A Proxy and a Decorator have the exact same structural code, but different architectural intents.**
*Answer*: True. Both wrap a real object and implement its interface. The Decorator intends to *add features*, while the Proxy intends to *control access* (and might prevent the real object from executing entirely).

---

# Interview Questions

## Beginner
**Q**: What is the primary purpose of the Proxy pattern?
*Answer*: To provide a surrogate object that controls access to a real object, allowing you to execute logic (caching, lazy loading, auth) before or after the real object is called, without modifying the client or the real object.

## Intermediate
**Q**: Explain how a Remote Proxy works in distributed systems (like gRPC).
*Answer*: A Remote Proxy (often called a client stub) acts as a local representative for an object that lives in a different memory space (a different server). The client calls a method on the proxy as if it were a local object. The proxy handles the complex network communication, serialization, and deserialization behind the scenes.

## Advanced
**Q**: What is the danger of using a Virtual Proxy (Lazy Loading) for Database relations in an N+1 scenario?
*Answer*: If you fetch 100 Users, and each User's `Profile` is a Virtual Proxy, looping through the 100 Users and calling `User.Profile.GetBio()` will trigger the Virtual Proxy to lazily load the profile 100 separate times, resulting in 100 individual database queries (the N+1 query problem). This destroys performance compared to a single JOIN query.

---

# Cheat Sheet

* **The Interface**: `type DB interface { Query() string }`
* **The Real Subject**: `type RealDB struct{}`
* **The Proxy**:
```go
type CacheProxy struct {
    real *RealDB
    cache string
}
func (p *CacheProxy) Query() string {
    if p.cache == "" {
        p.cache = p.real.Query() // Hit Real DB
    }
    return p.cache // Return Cache
}
```

---

# Summary

The Proxy pattern is the guardian of your expensive resources. Whether you are protecting a database from a stampede of requests using a Cache, protecting memory using Lazy Loading, or protecting a network via an API Gateway, the Proxy pattern ensures that your system remains robust, fast, and secure.

---

# Key Takeaways

* ✔ Proxies control access to a Real Subject.
* ✔ Both the Proxy and the Real Subject implement the same interface.
* ✔ Virtual Proxies save memory via Lazy Initialization.
* ✔ Caching Proxies save CPU/Network by intercepting redundant requests.

---

# Further Reading
* [Refactoring.guru: Proxy Pattern](https://refactoring.guru/design-patterns/proxy)

---

# Next Chapter
➡️ **Next:** `13-Strategy.md` (Beginning of Part 4: Behavioral Patterns)
