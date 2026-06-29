# Chain of Responsibility Pattern

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

---

# Introduction

The **Chain of Responsibility Pattern** is a Behavioral Design Pattern that lets you pass requests along a chain of handlers. Upon receiving a request, each handler decides either to process the request or to pass it to the next handler in the chain.

If this sounds exactly like the HTTP Middleware (Decorator) pattern from Chapter 10, you are correct! In Go, the Chain of Responsibility is most frequently implemented using Middlewares. However, the classic object-oriented implementation uses a linked list of Handler structs, which is highly useful for processing sequential validation or escalating support tickets.

---

# Learning Objectives

After completing this chapter you will be able to:

* Build a linked list of handler objects.
* Understand the difference between the Decorator pattern (execution wraps back around) and Chain of Responsibility (execution stops or flows one way).
* Implement sequential validation pipelines.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Interfaces.
* Linked Lists (Pointers to the next node).
* The Decorator Pattern (`10-Decorator.md`).

---

# Why This Topic Exists

Imagine building a hospital triage system. A patient arrives with a symptom.
1. The Receptionist checks if the patient has insurance. If not, they are rejected.
2. The Nurse checks if the symptom is minor. If yes, the nurse treats them and the chain stops.
3. If the symptom is severe, the Nurse passes the patient to the General Doctor.
4. If it's a specific organ failure, the General Doctor passes them to a Surgeon.

If you code this in a single `ProcessPatient()` function, you will have deeply nested `if/else` blocks that violate the Single Responsibility Principle. 

By using the Chain of Responsibility, you encapsulate each staff member into their own Handler struct. The patient is passed down the chain. Each handler independently decides: "Can I handle this? If yes, stop. If no, pass it to the `next` handler."

---

# Real-World Analogy

### Tech Support Escalation

* **Level 1 (Chatbot)**: Tries to answer basic FAQs. If it can't, it forwards the request.
* **Level 2 (Human Agent)**: Can handle billing and account resets. If it's a severe bug, they forward the request.
* **Level 3 (Senior Engineer)**: Receives the escalated ticket, finds the bug, and patches the database. The chain ends here.

The customer doesn't know who fixed their problem. They just submitted the ticket to the start of the chain.

---

# Core Concepts

* **Handler Interface**: Defines a method to process the request, and a method to set the `next` handler in the chain.
* **Concrete Handlers**: Structs that implement the interface. They contain the logic to evaluate the request.
* **The `next` Pointer**: Every Concrete Handler holds a pointer to the next handler in the chain (forming a Linked List).

---

# Architecture Diagram

```mermaid
flowchart LR
    Client[Client Request]
    H1[Handler 1<br/>(Auth)]
    H2[Handler 2<br/>(Validation)]
    H3[Handler 3<br/>(Database Save)]

    Client -- "Submits" --> H1
    
    H1 -- "Fails" --> Reject1[Reject Request]
    H1 -- "Passes" --> H2
    
    H2 -- "Fails" --> Reject2[Reject Request]
    H2 -- "Passes" --> H3
    
    H3 -- "Succeeds" --> Finish[Done]
```

---

# Step-by-Step Implementation

1. Define the **Handler Interface** containing `Execute(*Request)` and `SetNext(Handler)`.
2. Create **Concrete Handlers**. Give each one a `next Handler` field.
3. Implement `SetNext()`, which assigns the `next` field and usually returns the passed handler to allow for method chaining during setup.
4. Implement `Execute()`. Inside, write the specific logic. 
5. If the handler cannot fully process the request (or if it's a sequential pipeline), call `if h.next != nil { h.next.Execute(request) }`.
6. In `main`, instantiate the handlers, link them together, and pass the request to the first one.

---

# Syntax

```go
type Handler interface {
    Execute(data string)
    SetNext(h Handler) Handler
}

type BaseHandler struct {
    next Handler
}
func (b *BaseHandler) SetNext(h Handler) Handler {
    b.next = h
    return h
}
func (b *BaseHandler) ExecuteNext(data string) {
    if b.next != nil { b.next.Execute(data) }
}
```

---

# Beginner Example

A Patient Triage system (Escalation Chain).

```go
package main

import "fmt"

// 1. Interface
type Department interface {
	Execute(patient *Patient)
	SetNext(Department) Department
}

type Patient struct {
	Name      string
	Condition string
}

// 2. Concrete Handler 1
type Reception struct {
	next Department
}
func (r *Reception) SetNext(next Department) Department {
	r.next = next
	return next
}
func (r *Reception) Execute(p *Patient) {
	fmt.Println("Reception registering patient...")
	// Pass to next
	if r.next != nil {
		r.next.Execute(p)
	}
}

// 2. Concrete Handler 2
type Nurse struct {
	next Department
}
func (n *Nurse) SetNext(next Department) Department {
	n.next = next
	return next
}
func (n *Nurse) Execute(p *Patient) {
	if p.Condition == "Paper Cut" {
		fmt.Println("Nurse treated the paper cut. Chain ends here.")
		return // DO NOT pass to next!
	}
	fmt.Println("Nurse cannot treat condition. Escalating...")
	if n.next != nil {
		n.next.Execute(p)
	}
}

// 2. Concrete Handler 3
type Surgeon struct {
	next Department
}
func (s *Surgeon) SetNext(next Department) Department {
	s.next = next
	return next
}
func (s *Surgeon) Execute(p *Patient) {
	fmt.Println("Surgeon is operating on:", p.Condition)
}

func main() {
	// Build the chain
	reception := &Reception{}
	nurse := &Nurse{}
	surgeon := &Surgeon{}

	// Link them: Reception -> Nurse -> Surgeon
	reception.SetNext(nurse).SetNext(surgeon)

	// Scenario 1: Minor issue
	p1 := &Patient{Name: "Alice", Condition: "Paper Cut"}
	fmt.Println("--- Patient 1 Arrives ---")
	reception.Execute(p1)

	// Scenario 2: Major issue
	p2 := &Patient{Name: "Bob", Condition: "Appendicitis"}
	fmt.Println("\n--- Patient 2 Arrives ---")
	reception.Execute(p2)
}
```

---

# Intermediate Example

A Request Validation Pipeline. If any handler fails, the chain is broken and an error is returned.

```go
package main

import (
	"errors"
	"fmt"
)

type Request struct {
	Token    string
	Role     string
	DataSize int
}

type Handler interface {
	Handle(req *Request) error
	SetNext(h Handler) Handler
}

// Embed a Base struct to save us from writing SetNext() on every handler!
type BaseHandler struct {
	next Handler
}
func (b *BaseHandler) SetNext(h Handler) Handler {
	b.next = h
	return h
}
func (b *BaseHandler) handleNext(req *Request) error {
	if b.next != nil {
		return b.next.Handle(req)
	}
	return nil
}

// Concrete Handlers
type AuthHandler struct { BaseHandler }
func (h *AuthHandler) Handle(req *Request) error {
	if req.Token != "valid_token" {
		return errors.New("Auth Failed")
	}
	fmt.Println("Auth Passed.")
	return h.handleNext(req) // Pass to next
}

type RoleHandler struct { BaseHandler }
func (h *RoleHandler) Handle(req *Request) error {
	if req.Role != "admin" {
		return errors.New("Insufficient Permissions")
	}
	fmt.Println("Role Passed.")
	return h.handleNext(req)
}

type SizeHandler struct { BaseHandler }
func (h *SizeHandler) Handle(req *Request) error {
	if req.DataSize > 100 {
		return errors.New("Payload Too Large")
	}
	fmt.Println("Size Passed.")
	return h.handleNext(req)
}

func main() {
	// Build Chain
	pipeline := &AuthHandler{}
	pipeline.SetNext(&RoleHandler{}).SetNext(&SizeHandler{})

	// Test 1: Succeeds
	req1 := &Request{Token: "valid_token", Role: "admin", DataSize: 50}
	fmt.Println("Testing Req1:")
	err := pipeline.Handle(req1)
	fmt.Println("Result:", err)

	// Test 2: Fails at Role
	fmt.Println("\nTesting Req2:")
	req2 := &Request{Token: "valid_token", Role: "user", DataSize: 50}
	err2 := pipeline.Handle(req2)
	fmt.Println("Result:", err2)
}
```

---

# CoR vs Decorator Pattern

In Go, these two patterns are often confused because they both involve wrapping/chaining execution.
* **Chain of Responsibility**: Linear flow. Handler A executes, and then calls Handler B. Once B finishes, it does *not* return data back to A for further processing. The chain either breaks early, or hits the end and stops.
* **Decorator (HTTP Middleware)**: Onion flow (Outside-In, Inside-Out). Decorator A executes "Before" logic, calls B. B executes, and *returns back to A*. Decorator A then executes "After" logic.

If you don't need "After" logic, you are implementing a Chain of Responsibility.

---

# Production Use Cases

### 1. Web Framework Validation
Many web frameworks use a Chain of Responsibility for their request validation pipelines. The request passes through a CSRF checker, a JSON payload size checker, and a Rate Limiter. If any fail, an HTTP Error is immediately returned and the chain stops.

### 2. Standard Logger Fallbacks
In complex logging systems (like `zap` or `logrus`), you can set up a chain of sinks. A log event is generated. Handler 1 tries to write to Elasticsearch. If the network is down, it catches the error and passes the log to Handler 2, which writes to a local file. 

---

# Performance Analysis

The Chain of Responsibility adds overhead proportional to the length of the chain. If a chain has 50 handlers, and a request requires reaching the 50th handler, you incur 50 interface method calls. While Go handles this easily, long chains can impact the latency of ultra-high-throughput systems. 

---

# Best Practices

* **Use Embedded Base Structs**: As shown in the Intermediate example, always create a `BaseHandler` struct that implements the `SetNext` and `handleNext` methods. Embed this base struct into your Concrete Handlers to eliminate repetitive boilerplate code.
* **Return the Handler**: Make `SetNext` return a `Handler` interface. This enables a fluent API for setup: `h1.SetNext(h2).SetNext(h3)`.

---

# Common Mistakes

### Forgetting to Call Next
If a handler successfully processes its portion of the validation but the programmer forgets to type `h.handleNext(req)`, the chain will silently stop halfway through, and the core business logic will never execute.

### Unhandled Requests
If a request traverses the entire chain (e.g., searching for a Support Agent who can fix the bug) and hits the end of the chain without being resolved, you must ensure your system gracefully handles this (e.g., returning a "No Handlers Available" error).

---

# Debugging Guide

* **"Request randomly stops"**: One of your Concrete Handlers evaluated a condition to true, but failed to return an error OR failed to call the `next` handler. Ensure all logical branches either explicitly error out, or call `handleNext`.

---

# Exercises

## Beginner
Create a Chain of `Loggers`. Interface: `Log(level, msg)`. Handlers: `InfoLogger`, `WarningLogger`, `ErrorLogger`. If `Log` is called with level "ERROR", the `InfoLogger` and `WarningLogger` should ignore it and pass it down until the `ErrorLogger` catches it and prints it.

## Intermediate
Implement the `BaseHandler` embedding pattern. Create a processing pipeline that modifies a string. `TrimSpaceHandler` -> `ToLowerHandler` -> `RemovePunctuationHandler`. Print the final string.

---

# Quiz

## Multiple Choice Questions
**1. What is the primary difference between Chain of Responsibility and Decorator?**
A) Decorator uses structs, CoR uses functions.
B) Decorator wraps execution (allows "before" and "after" logic), while CoR generally passes execution linearly forward without resuming previous handlers.
C) CoR is faster.
*Answer*: B

## True or False
**In the Chain of Responsibility, every handler in the chain MUST process the request.**
*Answer*: False. A handler can decide to completely ignore the request and simply pass it to the next handler, or it can process it and stop the chain entirely.

---

# Interview Questions

## Beginner
**Q**: What is the Chain of Responsibility pattern?
*Answer*: It is a pattern that links multiple handler objects into a chain. A request is passed to the first handler, which can either process it (and stop), reject it, or pass it to the next handler in the chain.

## Intermediate
**Q**: How does embedding a struct simplify the implementation of this pattern in Go?
*Answer*: Every handler requires a `next` field and a `SetNext()` method to build the linked list. By defining these once in a `BaseHandler` struct and embedding it into every concrete handler, you adhere to DRY (Don't Repeat Yourself) principles and keep the concrete handlers focused solely on business logic.

## Advanced
**Q**: In a high-throughput Go web server, would you implement request validation using a linked-list Chain of Responsibility, or an array-based loop?
*Answer*: While the classic linked-list CoR is elegant, modern Go web frameworks (like Gin) actually use an array (slice) of handler functions and loop through them. Array iteration is significantly more cache-friendly and faster than chasing pointers through a linked list on the heap. So for critical path performance, array-based pipelines are preferred.

---

# Cheat Sheet

* **Base Handler (Embedding)**:
```go
type Base struct { next Handler }
func (b *Base) SetNext(h Handler) Handler { b.next = h; return h }
func (b *Base) Next(req string) { if b.next != nil { b.next.Do(req) } }
```
* **Concrete Handler**:
```go
type Concrete struct { Base }
func (c *Concrete) Do(req string) {
    if req == "bad" { return } // Reject
    c.Next(req) // Pass
}
```

---

# Summary

The Chain of Responsibility is a brilliant way to decouple the sender of a request from its receivers. By breaking massive, monolithic validation or processing functions into small, independent handlers, you create a system that is incredibly easy to test, extend, and reconfigure on the fly.

---

# Key Takeaways

* ✔ Forms a linked list of Handler objects.
* ✔ Handlers can stop the chain, or pass the request forward.
* ✔ Perfect for validation pipelines and escalation systems.
* ✔ Use embedded structs to share `SetNext` logic.

---

# Further Reading
* [Refactoring.guru: Chain of Responsibility Pattern](https://refactoring.guru/design-patterns/chain-of-responsibility)

---

# Conclusion of Design Patterns Curriculum
🎉 Congratulations! You have completed the comprehensive guide to Design Patterns in Go. You now possess the architectural vocabulary to build clean, maintainable, and highly scalable software.

**Next Module:** *Distributed Systems in Go*
