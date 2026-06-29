# Adapter Pattern

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

The **Adapter Pattern** is a Structural Design Pattern that allows objects with incompatible interfaces to collaborate. 

In Go, the Adapter pattern is one of the most frequently used patterns. Because Go encourages the use of small, consumer-defined interfaces (`02-Accept-Interfaces-Return-Structs.md`), you will frequently encounter third-party libraries or legacy structs that *almost* satisfy your interface, but the method names or signatures are slightly off. The Adapter solves this by wrapping the incompatible struct and mapping its methods to the required interface.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand how to bridge incompatible interfaces without modifying the original code.
* Use struct embedding to create clean adapters.
* Implement the Adapter pattern to integrate third-party SDKs into your business logic.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Interfaces.
* Composition (Struct Embedding).

---

# Why This Topic Exists

Imagine your business logic requires a `Logger` interface with a `Log(message string)` method.
You decide to install a popular third-party logging library (like `zap` or `logrus`). However, the third-party library uses a method named `Info(msg string)` instead of `Log(message string)`. 

You cannot edit the third-party library's source code. You do not want to change your entire business logic to match the third-party library (because you might swap it out later). The Adapter pattern is the glue that binds them together.

---

# Real-World Analogy

### The Universal Travel Adapter

* **The Target Interface**: A standard 120V US wall outlet. Your laptop expects to plug into this.
* **The Adaptee (Incompatible)**: A 220V European wall outlet. It provides power, but the prongs don't fit.
* **The Adapter**: A physical block you buy at the airport. You plug the European outlet into one side, and your US laptop plugs into the other. The block translates the physical shape and the voltage. The laptop has no idea it's actually pulling power from a European grid.

---

# Core Concepts

* **Target Interface**: The interface your core application expects to use.
* **Adaptee**: The existing, incompatible struct (usually legacy code or a third-party SDK).
* **Adapter**: A new struct that implements the Target Interface. It holds a reference to the Adaptee and translates calls from the Target Interface to the Adaptee's specific methods.

---

# Architecture Diagram

```mermaid
flowchart LR
    Client[Client Code]
    Target((Target Interface<br/>Log()))
    Adapter[Adapter Struct<br/>Log()]
    Adaptee[3rd Party Struct<br/>Info()]

    Client -- calls --> Target
    Adapter -. implements .-> Target
    Adapter -- "Translates Log() to Info()" --> Adaptee
```

---

# Step-by-Step Implementation

1. Identify the **Target Interface** your application requires.
2. Identify the **Adaptee** struct that has the functionality but the wrong signature.
3. Create an **Adapter** struct. Give it a field that holds a pointer to the Adaptee.
4. Write methods on the Adapter to satisfy the Target Interface.
5. Inside those methods, call the corresponding methods on the Adaptee, translating arguments or return values as necessary.
6. In `main.go`, wrap the Adaptee inside the Adapter and pass it to your application.

---

# Syntax

```go
// 1. Target
type Target interface { Request() string }

// 2. Adaptee (Incompatible)
type Adaptee struct{}
func (a *Adaptee) SpecificRequest() string { return "Data" }

// 3. Adapter
type Adapter struct {
    adaptee *Adaptee
}

// 4. Translate
func (a *Adapter) Request() string {
    return a.adaptee.SpecificRequest()
}
```

---

# Beginner Example

Adapting a legacy logger to a modern interface.

```go
package main

import "fmt"

// 1. Target Interface (What our app expects)
type Logger interface {
	Log(message string)
}

// 2. The Adaptee (Legacy or 3rd Party Code we cannot change)
type LegacyPrinter struct{}
func (l *LegacyPrinter) PrintToConsole(text string) {
	fmt.Println("LEGACY:", text)
}

// 3. The Adapter
type PrinterAdapter struct {
	printer *LegacyPrinter
}

// 4. Translation
func (a *PrinterAdapter) Log(message string) {
	// Translate the Log call into a PrintToConsole call
	a.printer.PrintToConsole(message)
}

// Client Code
func ProcessPayment(logger Logger) {
	logger.Log("Payment processed successfully.")
}

func main() {
	legacy := &LegacyPrinter{}
	
	// We cannot do: ProcessPayment(legacy) because it lacks the Log() method.
	// So we wrap it in the adapter!
	adapter := &PrinterAdapter{printer: legacy}
	
	ProcessPayment(adapter)
}
```

---

# Intermediate Example

Using struct embedding to create a cleaner Adapter. 
Sometimes you only need to adapt *one* method out of ten. By embedding the Adaptee, you automatically inherit all the matching methods, and you only need to explicitly write the adapter function for the incompatible one!

```go
package main

import "fmt"

// Target Interface (Needs 2 methods)
type Animal interface {
	Move()
	Speak()
}

// Adaptee (Matches Move, but fails Speak)
type Fish struct{}
func (f *Fish) Move() { fmt.Println("Swimming...") }
func (f *Fish) Blub() { fmt.Println("Blub blub...") } // Incompatible!

// Adapter using Embedding
type FishAdapter struct {
	*Fish // Embed the fish pointer!
}

// We only need to adapt Speak! Move() is automatically promoted from the embedded Fish.
func (a *FishAdapter) Speak() {
	a.Blub() // Translate
}

func MakeAnimalAct(a Animal) {
	a.Move()
	a.Speak()
}

func main() {
	fish := &Fish{}
	adapter := &FishAdapter{Fish: fish}
	
	MakeAnimalAct(adapter)
}
```

---

# Advanced Example

Data translation. The Adapter doesn't just rename functions; it translates data formats. Here we adapt a system that returns XML to an interface that expects JSON.

```go
package main

import (
	"encoding/json"
	"fmt"
)

// Target Interface (Expects JSON)
type AnalyticsEngine interface {
	AnalyzeData(jsonPayload []byte)
}

type ModernAnalytics struct{}
func (m *ModernAnalytics) AnalyzeData(jsonPayload []byte) {
	fmt.Printf("Analyzing JSON: %s\n", string(jsonPayload))
}

// Adaptee (Produces XML)
type LegacySensor struct{}
func (s *LegacySensor) ReadXML() string {
	return `<data><temp>72</temp><status>ok</status></data>`
}

// Adapter
type XMLToJSONAdapter struct {
	sensor *LegacySensor
	engine AnalyticsEngine
}

// The core business logic
func (a *XMLToJSONAdapter) Run() {
	// 1. Get XML
	xmlData := a.sensor.ReadXML()
	
	// 2. Translate XML to JSON (Simulated here)
	// In reality, you'd unmarshal XML and marshal to JSON
	jsonMap := map[string]string{"temp": "72", "status": "ok"}
	jsonData, _ := json.Marshal(jsonMap)
	
	// 3. Send to Target
	a.engine.AnalyzeData(jsonData)
}

func main() {
	sensor := &LegacySensor{}
	engine := &ModernAnalytics{}
	
	adapter := &XMLToJSONAdapter{
		sensor: sensor,
		engine: engine,
	}
	
	adapter.Run()
}
```

---

# Production Use Cases

### 1. Standard Library `http.HandlerFunc`
In Go, the `http.Handler` is an interface with a `ServeHTTP(w, r)` method. If you write a standard function `func MyHandler(w, r)`, it doesn't satisfy the interface because it's just a function, not an object with a method!
The standard library provides an adapter: `http.HandlerFunc(MyHandler)`. This takes your regular function and adapts it into a type that implements `ServeHTTP`.

### 2. AWS SDK Wrappers
If you are migrating an app from AWS DynamoDB to MongoDB, your core logic expects a `SaveUser()` method. DynamoDB uses `PutItem()`. MongoDB uses `InsertOne()`. You write a `DynamoAdapter` and a `MongoAdapter` that both satisfy your application's `Database` interface.

---

# Performance Analysis

The Adapter pattern has virtually zero overhead. It is simply a struct that makes a direct function call to another struct. Unless the Adapter is doing heavy data translation (like XML to JSON parsing), the performance impact of the extra function call is negligible (measured in fractions of a nanosecond).

---

# Best Practices

* **Adhere to the Single Responsibility Principle**: The Adapter should *only* translate data and method calls. Do not put heavy business logic inside the Adapter.
* **Use Embedding for Partial Adaptations**: If the Target interface has 10 methods, and the Adaptee already perfectly matches 9 of them, embed the Adaptee in the Adapter struct so you only have to write code for the 1 missing method.

---

# Common Mistakes

### Two-Way Adapters
Trying to create a single Adapter struct that can adapt `A` to `B`, and also adapt `B` to `A`. This creates bloated, confusing code. Write two separate adapters: `AToBAdapter` and `BToAAdapter`.

---

# Debugging Guide

* **"Cannot use adapter as type Target"**: You missed implementing one of the methods defined in the Target Interface on your Adapter struct, or the signature (return types/parameters) doesn't perfectly match.

---

# Exercises

## Beginner
Create a Target interface `Speaker` with method `SayHello()`. Create an Adaptee struct `FrenchPerson` with method `DireBonjour()`. Write a `TranslatorAdapter` that makes the FrenchPerson satisfy the `Speaker` interface.

## Intermediate
Examine the `http.Handler` interface. Write a standard function `func HelloWorld(w http.ResponseWriter, r *http.Request)`. Use the built-in `http.HandlerFunc` adapter to wrap your function and pass it to `http.ListenAndServe`.

---

# Quiz

## Multiple Choice Questions
**1. When should you use the Adapter pattern?**
A) When you want to add new features to an existing object dynamically.
B) When you need an existing class to work with a target interface, but you cannot change the existing class's source code.
C) When you want to ensure only one instance of an object exists.
*Answer*: B

## True or False
**In Go, an Adapter must use the `implements` keyword to declare that it satisfies the target interface.**
*Answer*: False. Go interfaces are satisfied implicitly. As long as the Adapter struct has the required methods, it works.

---

# Interview Questions

## Beginner
**Q**: What is the Adapter pattern?
*Answer*: It is a structural pattern that acts as a wrapper, allowing an object with an incompatible interface to be used by a system that expects a different interface.

## Intermediate
**Q**: How can struct embedding make writing an Adapter easier in Go?
*Answer*: If the Target interface requires 5 methods, and the Adaptee already has 4 of them with the correct names and signatures, you can embed the Adaptee inside the Adapter struct. Go will promote those 4 methods automatically. You then only need to write the wrapper function for the 1 incompatible method.

## Advanced
**Q**: Explain how `http.HandlerFunc` acts as an Adapter in the Go standard library.
*Answer*: The `http.ListenAndServe` function expects objects that satisfy the `http.Handler` interface (which requires a `ServeHTTP` method). Most developers just write standard functions. `http.HandlerFunc` is a defined function type with a `ServeHTTP` method attached to it. By casting your standard function to `http.HandlerFunc(myFunc)`, you are adapting a raw function into an object that satisfies the interface, allowing the router to call it.

---

# Cheat Sheet

* **Target Interface**: `type Target interface { Do() }`
* **Adaptee**: `type Legacy struct{}` / `func (l *Legacy) DoLegacy() {}`
* **Adapter**: 
```go
type Adapter struct { L *Legacy }
func (a *Adapter) Do() { a.L.DoLegacy() }
```

---

# Summary

The Adapter pattern is the ultimate integration tool. By wrapping third-party SDKs, legacy code, and external APIs in Adapters, you ensure that your core business logic remains pristine, completely decoupled from the specific naming conventions and data structures of the outside world.

---

# Key Takeaways

* ✔ Use Adapters to bridge incompatible interfaces without touching source code.
* ✔ Adapters translate method calls and data formats.
* ✔ Use struct embedding to save boilerplate when adapting large interfaces.

---

# Further Reading
* [Refactoring.guru: Adapter Pattern](https://refactoring.guru/design-patterns/adapter)

---

# Next Chapter
➡️ **Next:** `10-Decorator.md`
