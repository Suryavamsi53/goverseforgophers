# Facade Pattern

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

The **Facade Pattern** is a Structural Design Pattern that provides a simplified, higher-level interface to a complex body of code (like a framework, a third-party library, or a cluster of microservices).

In Go, packages naturally act as Facades. You expose a few simple, well-documented functions to the public API, while hiding the complex, interwoven struct logic within unexported (lowercase) types inside the package. The Facade pattern is all about reducing cognitive load for the developer using your code.

---

# Learning Objectives

After completing this chapter you will be able to:

* Hide complex subsystem initialization behind a simple API.
* Understand the difference between an Adapter and a Facade.
* Use Go packages to naturally enforce the Facade pattern using exported vs. unexported types.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Methods.
* Exported (Uppercase) vs Unexported (Lowercase) identifiers in Go.

---

# Why This Topic Exists

Imagine you want to upload a video to a cloud service. Under the hood, this requires:
1. Initializing an Authentication client and getting a Token.
2. Initializing a Storage client.
3. Compressing the video file.
4. Uploading the file in 10MB chunks.
5. Updating a Database with the video metadata.

If you force every developer on your team to write all 5 steps every time they want to upload a video, the codebase will become bloated, fragile, and prone to copy-paste errors. Instead, you create a `VideoUploader` Facade with a single method: `Upload(filepath string)`. 

---

# Real-World Analogy

### The Customer Service Agent

* **The Subsystem**: A massive retail company has a Warehouse Department, a Billing Department, a Shipping Department, and a Technical Support team.
* **The Problem**: A customer wants to return a broken laptop. If they had to navigate the subsystem themselves, they would have to call Tech Support for a RMA number, call Billing for a refund, and call Shipping for a return label.
* **The Facade**: The company hires a Customer Service Agent. The customer dials one number (the API). They say "I want to return this." The Agent (the Facade) handles coordinating all the complex internal departments on the customer's behalf.

---

# Core Concepts

* **The Facade**: A struct (or a package) that provides a simple interface.
* **The Subsystem**: A complex network of structs, functions, or external services that perform the actual work.
* **Loose Coupling**: The client code only interacts with the Facade. If the underlying subsystem changes (e.g., swapping AWS for Google Cloud), the client code doesn't break, because the Facade hides the transition.

---

# Architecture Diagram

```mermaid
flowchart TD
    Client[Client Code]
    Facade[Facade API<br/>UploadVideo()]
    
    subgraph Complex Subsystem
    Auth[Auth Client]
    Compress[Video Compressor]
    Store[Cloud Storage]
    DB[Database]
    end

    Client -- "Calls 1 simple method" --> Facade
    
    Facade -- "1. Get Token" --> Auth
    Facade -- "2. Compress File" --> Compress
    Facade -- "3. Upload Chunks" --> Store
    Facade -- "4. Save Metadata" --> DB
```

---

# Step-by-Step Implementation

1. Identify the complex sequence of actions required to perform a common task.
2. Create a new struct (the Facade) to orchestrate this task.
3. Initialize the Facade with references to the required subsystem components (or initialize them inside the Facade's constructor).
4. Create a single, high-level method on the Facade (e.g., `Execute()`).
5. Move the complex sequence of actions inside the `Execute()` method.
6. The client code now simply instantiates the Facade and calls `Execute()`.

---

# Syntax

```go
type Facade struct {
    subA *subsystemA
    subB *subsystemB
}

func NewFacade() *Facade {
    return &Facade{
        subA: &subsystemA{},
        subB: &subsystemB{},
    }
}

// The single, simple API method
func (f *Facade) DoWork() {
    data := f.subA.GetData()
    f.subB.ProcessData(data)
}
```

---

# Beginner Example

A Smart Home System. Turning off the house when you go to bed requires controlling the lights, the thermostat, and the security system.

```go
package main

import "fmt"

// --- THE COMPLEX SUBSYSTEM ---

type Lights struct{}
func (l *Lights) TurnOff() { fmt.Println("Lights: OFF") }

type Thermostat struct{}
func (t *Thermostat) SetTemperature(degrees int) { fmt.Printf("Thermostat: %d°C\n", degrees) }

type SecuritySystem struct{}
func (s *SecuritySystem) Arm() { fmt.Println("Security: ARMED") }


// --- THE FACADE ---

type SmartHomeFacade struct {
	lights   *Lights
	thermo   *Thermostat
	security *SecuritySystem
}

func NewSmartHome() *SmartHomeFacade {
	return &SmartHomeFacade{
		lights:   &Lights{},
		thermo:   &Thermostat{},
		security: &SecuritySystem{},
	}
}

// The simple API for the user
func (s *SmartHomeFacade) Goodnight() {
	fmt.Println("--- Executing Goodnight Routine ---")
	s.lights.TurnOff()
	s.thermo.SetTemperature(18)
	s.security.Arm()
	fmt.Println("--- Goodnight! ---")
}

func main() {
	// The client only needs to know about the Facade!
	home := NewSmartHome()
	home.Goodnight()
}
```

---

# Intermediate Example

The Facade as a Go Package. In Go, you often don't even need a `Facade` struct. You can use a package to act as a Facade. You export the simple functions, and hide the complex structs by making them unexported (lowercase).

File: `payment/payment.go`
```go
package payment

import "fmt"

// Unexported subsystem structs
type fraudChecker struct{}
func (f *fraudChecker) check(amount float64) bool { return amount < 10000 }

type bankAPI struct{}
func (b *bankAPI) transfer(amount float64) { fmt.Println("Transferred:", amount) }

type receiptGenerator struct{}
func (r *receiptGenerator) generate() { fmt.Println("Receipt emailed.") }

// Exported Facade Function
// The user only sees this one function!
func ProcessTransaction(amount float64) error {
	f := &fraudChecker{}
	if !f.check(amount) {
		return fmt.Errorf("fraud detected")
	}

	b := &bankAPI{}
	b.transfer(amount)

	r := &receiptGenerator{}
	r.generate()

	return nil
}
```

File: `main.go`
```go
package main

import "payment"

func main() {
	// The client code is beautiful and simple.
	payment.ProcessTransaction(50.00)
}
```

---

# Production Use Cases

### 1. Cloud SDKs (AWS / GCP)
When you use a high-level library like `aws-sdk-go` to upload a file to S3 via a `Uploader` struct, that struct is a Facade. Under the hood, it is managing HTTP clients, calculating MD5 checksums, splitting the file into multipart chunks, managing exponential backoff retries, and handling concurrency. You just call `Upload()`.

### 2. Database Migrations
A database migration tool provides a simple command: `MigrateUp()`. Behind the scenes, it connects to the database, reads a `schema_migrations` table to find the current version, parses SQL files from the filesystem, runs them in an ACID transaction, and updates the version table.

---

# Performance Analysis

The Facade pattern has strictly zero performance impact. It is purely an architectural reorganization of code. You are simply moving 10 lines of code out of `main()` and into a helper function/struct.

---

# Best Practices

* **Don't force the Facade**: If an advanced user wants to bypass the Facade and interact with the Subsystem directly to achieve a highly customized result, let them! A Facade is a convenience tool, not a mandatory prison. (Unless you specifically hide the subsystem using unexported types).
* **Keep it Stateless**: Ideally, a Facade struct should not hold complex state. It should merely act as an orchestrator that passes data between the subsystem components.

---

# Common Mistakes

### The God Object Anti-Pattern
A Facade should simplify a *specific* workflow (like "Uploading a Video"). Do not create a single `AppFacade` that contains 50 methods for every single operation in your entire application. That creates a bloated "God Object" that becomes impossible to maintain.

---

# Debugging Guide

* **"Facade is too restrictive"**: If developers are constantly asking you to add new flags and options to your Facade method (`Execute(bool, bool, int, string)`), your Facade is leaking abstraction. At that point, it may be better to expose the underlying subsystem so developers can construct their own workflows.

---

# Exercises

## Beginner
Create a `ComputerFacade` with a method `TurnOn()`. Inside `TurnOn()`, it should call methods on three unexported structs: `cpu.Freeze()`, `memory.Load()`, and `cpu.Jump()`. Call the Facade from `main()`.

## Intermediate
Refactor a complex sequence into a package-level Facade. Create a package `downloader`. Create unexported functions for `resolveDNS`, `openTCP`, and `downloadBytes`. Expose a single exported function `Download(url string)` that orchestrates them.

---

# Quiz

## Multiple Choice Questions
**1. What is the primary difference between an Adapter and a Facade?**
A) An Adapter changes the interface of an existing object to match a new interface. A Facade defines a new, simpler interface to wrap a complex subsystem of multiple objects.
B) An Adapter is used for software, a Facade is used for hardware.
C) They are exactly the same thing.
*Answer*: A

## True or False
**A client is strictly forbidden from bypassing the Facade and talking to the subsystem directly.**
*Answer*: False. A Facade provides a convenient shortcut for the 90% use case. Advanced users are often allowed to bypass it to interact with the raw subsystem for the 10% edge cases (unless the subsystem is deliberately hidden via unexported variables).

---

# Interview Questions

## Beginner
**Q**: What problem does the Facade pattern solve?
*Answer*: It solves the problem of tight coupling and high cognitive load. It prevents the client code from having to memorize and orchestrate dozens of complex, interdependent steps just to perform a common task.

## Intermediate
**Q**: How does Go's package visibility (`Exported` vs `unexported`) naturally lend itself to the Facade pattern?
*Answer*: In Go, any type or function starting with a lowercase letter is completely hidden from external packages. By writing all your complex subsystem logic using unexported structs, and only exporting a single Uppercase function or struct, you physically force the consumer to use your Facade. They literally cannot access the subsystem.

## Advanced
**Q**: Contrast the Facade pattern with the Builder pattern.
*Answer*: The Builder pattern simplifies the *creation* of a complex object step-by-step, ultimately returning a data structure. The Facade pattern simplifies the *execution* of a complex process, orchestrating behavior across multiple subsystems, usually returning a simple result or error.

---

# Cheat Sheet

* **The Subsystem**: `type a struct{}`, `type b struct{}`
* **The Facade**:
```go
type Facade struct { a *a; b *b }
func (f *Facade) SimpleAction() {
    f.a.DoA()
    f.b.DoB()
}
```

---

# Summary

The Facade pattern is the essence of good API design. Whenever you feel overwhelmed by the number of steps required to do something simple, you build a Facade. It encapsulates the messy reality of the underlying system, presenting a clean, beautiful, and foolproof interface to the world.

---

# Key Takeaways

* ✔ Facades hide complexity behind a single simple API.
* ✔ They decouple client code from the fragile internal subsystem.
* ✔ In Go, you can build Facades using packages and unexported types.
* ✔ Facades organize behavior, while Builders organize creation.

---

# Further Reading
* [Refactoring.guru: Facade Pattern](https://refactoring.guru/design-patterns/facade)

---

# Next Chapter
➡️ **Next:** `12-Proxy.md`
