# Functional Options Pattern

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

The **Functional Options** pattern is arguably the most famous and universally used Go-specific design pattern. Popularized by Dave Cheney and Rob Pike, it elegantly solves the problem of creating complex objects with many optional configuration parameters.

Unlike Java or Python, Go does not support function overloading or default arguments. Functional Options leverage Go's first-class functions and variadic arguments (`...`) to create clean, readable, and highly extensible APIs for object initialization.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand why standard constructor functions fail to scale in Go.
* Identify the drawbacks of using Config structs.
* Implement the Functional Options pattern from scratch.
* Design clean APIs that allow sensible defaults and flexible overrides.

---

# Prerequisites

Before reading this chapter you should know:

* First-class functions in Go.
* Variadic function arguments (`...Type`).
* Struct initialization.

---

# Why This Topic Exists

Imagine building an HTTP Server struct. A server needs a Port, a Timeout, a Max Connections limit, and a TLS certificate.
If you write a constructor:
`func NewServer(port int, timeout time.Duration, maxConn int, tls bool) *Server`

What if someone only wants to specify the Port, and leave the rest as default? In Go, you can't just omit the arguments. You'd have to pass `NewServer(8080, 0, 0, false)`. This is ugly and confusing. What does `0` mean? Is it infinite timeout or default timeout?

What if we add a new parameter later? We break the API for everyone who was calling `NewServer`! The Functional Options pattern solves this permanently.

---

# Real-World Analogy

### Ordering a Burger

* **Constructor Function**: The cashier asks you 20 questions in a row (Cheese? Pickles? Onions? Ketchup? Mayo? Bacon?). You must answer every single one, even if you just want a standard burger.
* **Config Struct**: You are handed a massive checklist with 20 boxes to fill out before you can hand it back to the cashier.
* **Functional Options**: You walk up and say, "I'll take a standard Burger." The cashier starts making it. As an afterthought, you add, "Oh, and add Bacon... and no onions." The cashier seamlessly applies your specific modifications to the sensible default burger.

---

# Core Concepts

* **Sensible Defaults**: The constructor creates the object with the most common, safe defaults.
* **The Option Type**: A function signature that accepts a pointer to the object being built (e.g., `type Option func(*Server)`).
* **Variadic Constructor**: The constructor accepts a variadic list of these Option functions (`func NewServer(opts ...Option)`).
* **Execution**: The constructor loops over the provided options, executing each one to modify the object.

---

# Architecture Diagram

```mermaid
flowchart TD
    Default[Create Default Server]
    Opt1[Apply Option 1: WithPort(9090)]
    Opt2[Apply Option 2: WithTLS()]
    Opt3[Apply Option N...]
    Return[Return configured *Server]
    
    Default --> Opt1
    Opt1 --> Opt2
    Opt2 --> Opt3
    Opt3 --> Return
```

---

# Step-by-Step Implementation

1. Define your struct (e.g., `Server`).
2. Define the Option type: `type Option func(*Server)`.
3. Write functions that return this Option type. Name them starting with `With...` (e.g., `WithPort(port int) Option`).
4. Inside the `With...` function, return a closure that mutates the pointer to the struct.
5. Write your constructor: `func NewServer(opts ...Option) *Server`.
6. Inside the constructor, initialize the struct with defaults.
7. Loop over the `opts` slice and call `opt(server)` for each one.
8. Return the server.

---

# Syntax

```go
type Server struct { port int }

type Option func(*Server)

func WithPort(p int) Option {
    return func(s *Server) { s.port = p }
}

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080} // Default
    for _, opt := range opts {
        opt(s) // Apply modification
    }
    return s
}
```

---

# Beginner Example

Building a simple Coffee struct.

```go
package main

import "fmt"

type Coffee struct {
	Size   string
	Shots  int
	Syrup  string
	Milk   bool
}

// 1. Define the Option type
type CoffeeOption func(*Coffee)

// 2. Define functional options
func WithSize(size string) CoffeeOption {
	return func(c *Coffee) { c.Size = size }
}

func WithExtraShots(shots int) CoffeeOption {
	return func(c *Coffee) { c.Shots = shots }
}

func WithSyrup(syrup string) CoffeeOption {
	return func(c *Coffee) { c.Syrup = syrup }
}

func WithMilk() CoffeeOption {
	return func(c *Coffee) { c.Milk = true }
}

// 3. The Constructor
func NewCoffee(opts ...CoffeeOption) *Coffee {
	// Start with a sensible default
	coffee := &Coffee{
		Size:  "Medium",
		Shots: 1,
		Syrup: "None",
		Milk:  false,
	}

	// Apply any overrides the user provided
	for _, opt := range opts {
		opt(coffee)
	}

	return coffee
}

func main() {
	// A default coffee
	defaultCoffee := NewCoffee()
	fmt.Printf("Default: %+v\n", defaultCoffee)

	// A custom coffee
	customCoffee := NewCoffee(
		WithSize("Large"),
		WithExtraShots(3),
		WithMilk(),
	)
	fmt.Printf("Custom: %+v\n", customCoffee)
}
```

---

# Intermediate Example

Using an interface for the Option type. Some libraries (like gRPC) prefer using an interface instead of a raw function type. This allows for slightly more complex behavior inside the option, but the fundamental pattern is identical.

```go
package main

import "fmt"

type Database struct {
	Host     string
	MaxConns int
}

// 1. Define an Interface instead of a func type
type DBOption interface {
	apply(*Database)
}

// 2. We need a concrete type that implements the interface.
// Using a function type that implements its own interface is a neat Go trick!
type optionFunc func(*Database)

func (f optionFunc) apply(db *Database) {
	f(db)
}

// 3. Option constructors
func WithHost(host string) DBOption {
	return optionFunc(func(db *Database) {
		db.Host = host
	})
}

func WithMaxConnections(conns int) DBOption {
	return optionFunc(func(db *Database) {
		db.MaxConns = conns
	})
}

func NewDatabase(opts ...DBOption) *Database {
	db := &Database{
		Host:     "localhost:5432",
		MaxConns: 10,
	}
	
	for _, opt := range opts {
		opt.apply(db)
	}
	return db
}

func main() {
	db := NewDatabase(WithMaxConnections(100))
	fmt.Printf("DB: %+v\n", db)
}
```

---

# Advanced Example

Handling errors during option application. What if an option receives invalid data? We can change the signature to `type Option func(*Server) error`.

```go
package main

import (
	"errors"
	"fmt"
)

type Worker struct {
	Timeout int
}

type WorkerOption func(*Worker) error

func WithTimeout(seconds int) WorkerOption {
	return func(w *Worker) error {
		if seconds <= 0 {
			return errors.New("timeout must be greater than zero")
		}
		w.Timeout = seconds
		return nil
	}
}

func NewWorker(opts ...WorkerOption) (*Worker, error) {
	w := &Worker{Timeout: 30} // default
	
	for _, opt := range opts {
		if err := opt(w); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}
	
	return w, nil
}

func main() {
	// Success
	w1, _ := NewWorker(WithTimeout(60))
	fmt.Println("Worker 1 Timeout:", w1.Timeout)

	// Failure
	_, err := NewWorker(WithTimeout(-5))
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

---

# Production Use Cases

### 1. gRPC Connections
The official Go gRPC library uses this pattern heavily. When you call `grpc.Dial()`, you pass in options like `grpc.WithInsecure()`, `grpc.WithBlock()`, or `grpc.WithTimeout()`.
`conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())`

### 2. AWS SDK
The AWS SDK for Go v2 uses functional options to configure everything from AWS credentials to specific DynamoDB query parameters. This allows the SDK to add new features to APIs without breaking backwards compatibility.

---

# Performance Analysis

* **Memory/CPU**: Executing a few closures during object initialization takes microscopic amounts of CPU time (nanoseconds). It is a one-time cost paid during initialization.
* **Readability vs Config Structs**: Passing a `Config{}` struct is slightly faster because it doesn't allocate closures, but it clutters the API. For singletons, servers, and clients, Functional Options are vastly superior. For extremely high-throughput, per-request object creation, a `Config` struct might be better.

---

# Best Practices

* **Prefix with `With`**: Name your option functions starting with `With` (e.g., `WithLogger()`, `WithTimeout()`). This makes your API discoverable via autocomplete in IDEs.
* **Keep defaults sensible**: If the user calls `New()` with 0 options, the returned object should work perfectly for 80% of use cases.
* **Keep the struct private**: The struct fields being modified should usually be unexported (lowercase) so they cannot be mutated after the object is created, making it thread-safe.

---

# Common Mistakes

### The Empty Config Struct Anti-Pattern
```go
// BAD: The user MUST pass a struct, even if it's empty.
client := NewClient(Config{}) 

// GOOD: The user passes nothing for defaults.
client := NewClient() 
```

### Modifying state after creation
Functional options should only be applied *during* creation. Do not expose a method like `server.Apply(WithPort(9090))` after the server is already running, as this introduces race conditions.

---

# Debugging Guide

* **"Option has no effect"**: Ensure your `With...` function returns a closure that mutates a **pointer** to the struct `func(s *Server)`. If it takes a value `func(s Server)`, it modifies a copy!
* **Compiler Error: cannot use Option (type) as type Option**: Usually happens when mixing interface-based options and func-based options. Ensure your types match exactly.

---

# Exercises

## Beginner
Create a `Car` struct with `Color` (string) and `Speed` (int). Write a constructor `NewCar` that defaults to a "Red" car going "100" mph. Implement `WithColor` and `WithSpeed` functional options.

## Intermediate
Refactor the Beginner exercise so that `WithSpeed` returns an error if the speed is over 200. Update the `NewCar` constructor to return `(*Car, error)`.

---

# Quiz

## Multiple Choice Questions
**1. Why are Functional Options preferred over long constructor parameter lists in Go?**
A) Go doesn't support pointers.
B) Go doesn't support default arguments or function overloading.
C) Functional options execute faster.
*Answer*: B

## True or False
**Functional options should be used to mutate an object's state *while* the program is running.**
*Answer*: False. They are an initialization pattern. Using them to mutate state later can cause thread-safety issues (data races).

---

# Interview Questions

## Beginner
**Q**: What problem does the Functional Options pattern solve?
*Answer*: It provides a clean, scalable way to configure complex objects with sensible defaults, avoiding massive constructor signatures and maintaining backwards compatibility if new configuration fields are added.

## Intermediate
**Q**: Compare the Functional Options pattern to passing a `Config` struct.
*Answer*: A `Config` struct is simpler to write, but if you add a new field, you might need to handle empty/zero values awkwardly (e.g., distinguishing between a user wanting a 0-second timeout vs not specifying a timeout). Functional Options clearly define intent, encapsulate validation, and read much closer to natural language (`New(WithTimeout(5))`), though they require slightly more boilerplate to set up.

## Advanced
**Q**: In the advanced Functional Options pattern using an interface (like `grpc.DialOption`), why do we define a `func` type that implements the interface?
*Answer*: Defining `type optionFunc func(*Config)` and giving it a method `func (f optionFunc) apply(c *Config) { f(c) }` allows us to write simple closure functions that automatically satisfy the interface. It's an elegant Go idiom that prevents us from having to define a completely new concrete struct for every single option we want to create.

---

# Cheat Sheet

* **The Boilerplate**:
```go
type Server struct { port int }

type Option func(*Server)

func WithPort(p int) Option {
    return func(s *Server) { s.port = p }
}

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080}
    for _, opt := range opts { opt(s) }
    return s
}
```

---

# Summary

The Functional Options pattern perfectly encapsulates the Go ethos: it leverages the language's specific strengths (first-class functions, variadic arguments) to solve a problem without needing complex OOP features like classes or inheritance. Mastering this pattern is a rite of passage for Go developers building libraries and APIs.

---

# Key Takeaways

* ✔ Replaces messy constructors and config structs.
* ✔ Extensible without breaking backwards compatibility.
* ✔ Always start option function names with `With...`.
* ✔ Provide sensible defaults before applying options.

---

# Further Reading
* [Dave Cheney: Functional options for friendly APIs](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)

---

# Next Chapter
➡️ **Next:** `02-Accept-Interfaces-Return-Structs.md`
