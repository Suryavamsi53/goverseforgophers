# Strategy Pattern

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

The **Strategy Pattern** is a Behavioral Design Pattern that lets you define a family of algorithms, put each of them into a separate class/struct, and make their objects interchangeable at runtime.

Because Go uses implicit interfaces and first-class functions, the Strategy pattern is one of the most natural and heavily used patterns in the language. It completely eliminates massive `if/else` or `switch` statements that bloat your core business logic, replacing them with clean, interchangeable interfaces.

---

# Learning Objectives

After completing this chapter you will be able to:

* Identify code smells (massive switch statements) that indicate the need for a Strategy pattern.
* Implement interchangeable algorithms using Go interfaces.
* Implement the Strategy pattern dynamically using first-class functions instead of structs.
* Swap an object's behavior at runtime without changing its class.

---

# Prerequisites

Before reading this chapter you should know:

* Go Interfaces (`02-Accept-Interfaces-Return-Structs.md`).
* First-class functions.

---

# Why This Topic Exists

Imagine you are building an E-Commerce application. You need a function to calculate the final price of an order.
* First, you add a flat discount.
* Then, marketing adds a percentage discount.
* Then, they add a "Buy One Get One Free" discount.

Your `CalculatePrice()` function now has a massive 50-line `switch` statement checking `order.DiscountType`. Every time a new discount type is invented, you have to modify this core function, risking breaking the entire checkout system (violating the Open/Closed Principle).

The Strategy pattern solves this. You extract each discount calculation into its own struct (the Strategy). The order simply holds an interface `DiscountStrategy`, and calls `strategy.Apply()`. The core checkout code never changes again.

---

# Real-World Analogy

### Google Maps Navigation

* **The Context**: You want to get from Point A to Point B.
* **The Problem**: Depending on your situation, you might want to Walk, Drive, or take Transit.
* **The Strategy**: Google Maps doesn't have one massive, jumbled function for routing. It has a generic `RouteStrategy` interface. 
  * If you click the Car icon, it swaps the strategy to the `DriveRouting` struct. 
  * If you click the Train icon, it swaps the strategy to the `TransitRouting` struct.
  * The map UI just calls `strategy.BuildRoute()`. It doesn't care how the route is built.

---

# Core Concepts

* **Strategy Interface**: A common interface that all concrete strategies must implement.
* **Concrete Strategies**: Various structs (or functions) that implement the interface, each containing a different algorithm.
* **Context**: The object that holds a reference to a Strategy interface. It delegates the work to the linked strategy instead of executing it itself.

---

# Architecture Diagram

```mermaid
flowchart TD
    Client[Client Code]
    Context[Order Struct<br/>Holds Strategy Interface]
    Interface((DiscountStrategy<br/>Apply()))
    
    Flat[FlatDiscount Struct]
    Percent[PercentDiscount Struct]
    BOGO[BOGODiscount Struct]

    Client -- "1. Sets Strategy" --> Context
    Context -- "2. Delegates Work" --> Interface
    
    Flat -. implements .-> Interface
    Percent -. implements .-> Interface
    BOGO -. implements .-> Interface
```

---

# Step-by-Step Implementation

1. Identify an algorithm that changes frequently based on context (e.g., sorting, compression, payment processing).
2. Declare the **Strategy Interface** containing the method(s) required to execute the algorithm.
3. Extract the `if/else` logic into separate **Concrete Strategy** structs that implement the interface.
4. Add a field to the **Context** struct to hold the Strategy interface.
5. Create a `SetStrategy()` method on the Context to allow changing the strategy at runtime.
6. The Context delegates the work by calling the interface method.

---

# Syntax

```go
// 1. Interface
type Strategy interface { Execute(data string) string }

// 2. Concrete Strategy
type ConcreteStrategyA struct{}
func (s *ConcreteStrategyA) Execute(d string) string { return d + " A" }

// 3. Context
type Context struct {
    strategy Strategy
}
func (c *Context) SetStrategy(s Strategy) { c.strategy = s }
func (c *Context) DoWork() { c.strategy.Execute("data") }
```

---

# Beginner Example

A Payment Processing system.

```go
package main

import "fmt"

// 1. The Strategy Interface
type PaymentStrategy interface {
	Pay(amount float64)
}

// 2. Concrete Strategy A
type CreditCard struct {
	CardNumber string
}
func (c *CreditCard) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using Credit Card ending in %s\n", amount, c.CardNumber[len(c.CardNumber)-4:])
}

// 2. Concrete Strategy B
type PayPal struct {
	Email string
}
func (p *PayPal) Pay(amount float64) {
	fmt.Printf("Paid $%.2f using PayPal account %s\n", amount, p.Email)
}

// 3. The Context
type ShoppingCart struct {
	amount   float64
	strategy PaymentStrategy
}

func (s *ShoppingCart) SetStrategy(strategy PaymentStrategy) {
	s.strategy = strategy
}

func (s *ShoppingCart) Checkout() {
	if s.strategy == nil {
		fmt.Println("Error: No payment strategy selected!")
		return
	}
	s.strategy.Pay(s.amount)
}

func main() {
	cart := &ShoppingCart{amount: 100.50}

	// Pay with Credit Card
	cart.SetStrategy(&CreditCard{CardNumber: "1234567890123456"})
	cart.Checkout()

	// Change mind, pay with PayPal instead! (Runtime Swap)
	cart.SetStrategy(&PayPal{Email: "user@example.com"})
	cart.Checkout()
}
```

---

# Intermediate Example

The **Functional Strategy Pattern**. Because Go has first-class functions, we don't always need to declare structs and interfaces! We can just use a function signature as the strategy type. This significantly reduces boilerplate for simple algorithms.

```go
package main

import "fmt"

// 1. Define a Function Signature as the Strategy
type DiscountStrategy func(price float64) float64

// 2. Define Concrete Strategies as normal functions
func FlatDiscount(price float64) float64 {
	return price - 10.0 // $10 off
}

func PercentageDiscount(price float64) float64 {
	return price * 0.8 // 20% off
}

// 3. The Context
type Order struct {
	Price    float64
	// Store the function
	Strategy DiscountStrategy 
}

func (o *Order) CalculateFinalPrice() float64 {
	if o.Strategy == nil {
		return o.Price
	}
	// Execute the function
	return o.Strategy(o.Price) 
}

func main() {
	order := &Order{Price: 100.0}

	// Swap strategies at runtime just by passing functions!
	order.Strategy = FlatDiscount
	fmt.Printf("Flat Discount Price: $%.2f\n", order.CalculateFinalPrice())

	order.Strategy = PercentageDiscount
	fmt.Printf("Percentage Discount Price: $%.2f\n", order.CalculateFinalPrice())
}
```

---

# Advanced Example

Combining the Factory Method and Strategy patterns. In a production REST API, you usually receive a string payload (e.g., `{"compression": "gzip"}`). You use a Factory to instantiate the correct Strategy, and then inject it into your Context.

```go
package main

import (
	"fmt"
	"strings"
)

// --- STRATEGIES ---
type CompressionStrategy interface {
	Compress(data string) string
}

type GzipStrategy struct{}
func (g *GzipStrategy) Compress(data string) string { return "GZIP[" + data + "]" }

type SnappyStrategy struct{}
func (s *SnappyStrategy) Compress(data string) string { return "SNAPPY[" + data + "]" }


// --- THE FACTORY ---
func GetCompressionStrategy(alg string) CompressionStrategy {
	switch strings.ToLower(alg) {
	case "gzip":
		return &GzipStrategy{}
	case "snappy":
		return &SnappyStrategy{}
	default:
		return nil
	}
}


// --- THE CONTEXT ---
type Archiver struct {
	strategy CompressionStrategy
}
func (a *Archiver) Archive(data string) {
	if a.strategy == nil {
		fmt.Println("Error: Invalid compression strategy")
		return
	}
	compressed := a.strategy.Compress(data)
	fmt.Println("Saved to disk:", compressed)
}


// --- THE CLIENT ---
func main() {
	// 1. Receive user input (e.g., from an HTTP Request)
	userInput := "snappy"
	
	// 2. Use Factory to resolve the strategy
	strategy := GetCompressionStrategy(userInput)
	
	// 3. Inject into Context
	archiver := &Archiver{strategy: strategy}
	
	// 4. Execute
	archiver.Archive("My Important File Data")
}
```

---

# Production Use Cases

### 1. The `sort` Package
The Go standard library `sort` package is a perfect example of the Strategy pattern. `sort.Sort()` takes an interface `sort.Interface` (which requires `Len`, `Less`, and `Swap` methods). The sorting algorithm inside `sort.Sort` doesn't know what it's sorting. You provide the Strategy by implementing those 3 methods on your custom slice.

### 2. File Uploaders
If your application can upload images to AWS S3, Google Cloud, or Azure, you define an `UploaderStrategy` interface. Based on the user's configuration, you inject the correct concrete struct into your core image processing service.

---

# Performance Analysis

Using interface-based strategies adds a microscopic amount of overhead due to dynamic dispatch (method lookup at runtime). Using functional strategies adds a tiny amount of overhead due to closure execution. Both are entirely negligible in standard applications. The massive boost in code readability and maintainability is well worth it.

---

# Best Practices

* **Use Functional Strategies for simple logic**: If the strategy only requires a single method and no state (like the Discount example), use a function signature (`type Strategy func()`). It is significantly more idiomatic in Go than writing single-method structs.
* **Use Struct Strategies for stateful logic**: If the strategy requires a database connection, API keys, or multiple methods (like the Payment example), use an interface and concrete structs.

---

# Common Mistakes

### Strategies that know too much
A Strategy should be entirely self-contained. It should receive all the data it needs via parameters. Do not pass a pointer to the entire `Context` struct into the Strategy, as this creates a tight circular dependency between the Strategy and the Context.

---

# Debugging Guide

* **"panic: runtime error: invalid memory address or nil pointer dereference"**: The Context tried to execute the strategy, but `SetStrategy` was never called (the strategy field is `nil`). Always add a nil check `if c.strategy == nil` before executing.

---

# Exercises

## Beginner
Create a `RouteStrategy` interface with a method `BuildRoute(A, B string)`. Create two structs: `WalkStrategy` and `DriveStrategy`. Create a `Navigator` context struct that accepts the strategy and prints the route.

## Intermediate
Refactor the Beginner exercise to use the Functional Strategy pattern. Define `type RouteStrategy func(A, B string)`. Write `Walk` and `Drive` as standard functions and pass them to the `Navigator`.

---

# Quiz

## Multiple Choice Questions
**1. What architectural principle does the Strategy pattern heavily enforce?**
A) The Singleton Principle.
B) The Open/Closed Principle (Open for extension, closed for modification).
C) The Liskov Substitution Principle.
*Answer*: B. By extracting algorithms into strategies, you can add new algorithms (extension) without touching the context's core code (modification).

## True or False
**In Go, the Strategy pattern MUST be implemented using Interfaces and Structs.**
*Answer*: False. Because Go has first-class functions, simple strategies are often implemented purely using function signatures.

---

# Interview Questions

## Beginner
**Q**: What is the Strategy pattern?
*Answer*: It is a behavioral pattern that allows you to define a family of algorithms, encapsulate each one, and make them interchangeable at runtime without altering the core context that uses them.

## Intermediate
**Q**: How does the Strategy pattern differ from the State pattern?
*Answer*: Structurally, they look identical (a Context holding an interface). However, in the Strategy pattern, the Client explicitly chooses the strategy and sets it on the Context. In the State pattern, the State itself controls when and how to transition to the next State, entirely hidden from the Client.

## Advanced
**Q**: Explain how the `sort` package in the Go standard library utilizes the Strategy pattern.
*Answer*: The `sort.Sort(data Interface)` function contains a core sorting algorithm (like QuickSort). However, it delegates the strategy for comparing and swapping elements back to the caller. The caller implements the `sort.Interface` (`Len`, `Less`, `Swap`) on their custom data type. The sorting algorithm doesn't care what the data is, it just executes the strategy provided by the caller.

---

# Cheat Sheet

* **Interface Strategy**:
```go
type Strategy interface { Execute() }
type Context struct { s Strategy }
func (c *Context) Do() { c.s.Execute() }
```
* **Functional Strategy**:
```go
type Strategy func()
type Context struct { s Strategy }
func (c *Context) Do() { c.s() }
```

---

# Summary

The Strategy Pattern is the antidote to "spaghetti code." By replacing bloated conditional statements with clean, isolated, interchangeable strategies, your code becomes incredibly easy to unit test and endlessly extensible. When combined with Go's first-class functions, it becomes one of the most powerful and lightweight tools in your architectural toolbelt.

---

# Key Takeaways

* ✔ Use Strategy to eliminate massive `switch` statements.
* ✔ Defines a family of interchangeable algorithms.
* ✔ In Go, you can use Interfaces or First-Class Functions.
* ✔ The Client determines which strategy the Context should use at runtime.

---

# Further Reading
* [Refactoring.guru: Strategy Pattern](https://refactoring.guru/design-patterns/strategy)

---

# Next Chapter
➡️ **Next:** `14-Observer.md`
