# Accept Interfaces, Return Structs

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

One of the most famous Go proverbs is: **"Accept interfaces, return structs."**

In traditional Object-Oriented languages (like Java), you typically define large, sprawling interfaces and return them from your factory methods. Go flips this on its head. In Go, interfaces are *implicit* and *small*. The best practice is to return concrete types (structs) from your constructors, and define small, specific interfaces exactly where you need to consume them as arguments.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand the "Accept Interfaces, Return Structs" proverb.
* Design APIs that are easy to test and mock.
* Prevent the "Package Dependency Cycle" nightmare.
* Define interfaces at the consumer level rather than the producer level.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Methods.
* Basic Go Interfaces (implicit satisfaction).

---

# Why This Topic Exists

Imagine you write a `UserService` package. If your `NewUserService()` function returns an interface `IUserService` containing 20 methods, any developer who wants to mock your service in their tests now has to create a mock object that implements all 20 methods, even if they only need to test the `GetUser` method!

Worse, if you return an interface, the consumer is tightly coupled to *your* definition of what that interface should be. By returning a concrete struct, the consumer can define their *own* 1-method interface (e.g., `UserFetcher`) locally, and your struct will implicitly satisfy it. This decoupling is the superpower of Go.

---

# Real-World Analogy

### The Universal Power Adapter

* **Returning an Interface**: An electronics manufacturer (Producer) builds a TV and decrees, "This TV must only be plugged into a 'Standardized 120V US Wall Outlet' interface." If you travel to Europe, you can't plug it in. The producer was too restrictive.
* **Accepting an Interface (The Go Way)**: The TV manufacturer returns a concrete `TV` object. However, the TV itself only *accepts* something that provides power: `type PowerSource interface { SupplyPower() }`. Now, the user (Consumer) can plug the TV into a US wall outlet, a European adapter, or a solar battery. As long as it has a `SupplyPower()` method, the TV accepts it.

---

# Core Concepts

* **Producer**: The package that defines and creates a struct (e.g., a Database client).
* **Consumer**: The package that uses the struct to do work (e.g., an HTTP handler).
* **Implicit Interfaces**: In Go, you don't use the `implements` keyword. If a struct has the methods, it satisfies the interface automatically.
* **Consumer-Defined Interfaces**: Interfaces should be defined in the package that *uses* them, not the package that *implements* them.

---

# Architecture Diagram

```mermaid
flowchart LR
    subgraph "Producer Package (e.g., /database)"
    Struct[Concrete Struct<br/>type DB struct{}]
    Method[Method: GetUser()]
    Struct --> Method
    end

    subgraph "Consumer Package (e.g., /http)"
    Interface[Consumer Interface<br/>type UserGetter interface {<br/>GetUser()<br/>}]
    Handler[HTTP Handler]
    Interface -.-> Handler
    end

    Struct -- Satisfies --> Interface
    note over Struct,Interface: The Consumer defines what it needs.<br/>The Producer implicitly provides it.
```

---

# Step-by-Step Implementation

1. **In the Producer Package**: Write your concrete struct and its methods.
2. **In the Producer Package**: Write a constructor that returns a pointer to the concrete struct `*MyStruct`. Do not define an interface here!
3. **In the Consumer Package**: Figure out exactly which methods you need to call on the struct.
4. **In the Consumer Package**: Define a small interface containing *only* those methods.
5. **In the Consumer Package**: Write your functions to accept that small interface as an argument.
6. **In `main.go`**: Pass the concrete struct into the consumer function.

---

# Syntax

```go
// --- PRODUCER PACKAGE ---
// Return a concrete struct
func NewClient() *Client { return &Client{} }


// --- CONSUMER PACKAGE ---
// Define a tiny interface right where you need it
type Reader interface {
    Read() []byte
}

// Accept the interface
func ProcessData(r Reader) {
    data := r.Read()
}
```

---

# Beginner Example

The Anti-Pattern vs The Go Pattern.

### The Anti-Pattern (Java Style)
```go
package store

// BAD: The producer defines a massive interface
type IDataStore interface {
	SaveUser()
	DeleteUser()
	GetMetrics()
	ClearCache()
}

type MySQLStore struct{}

func (m *MySQLStore) SaveUser() {}
func (m *MySQLStore) DeleteUser() {}
func (m *MySQLStore) GetMetrics() {}
func (m *MySQLStore) ClearCache() {}

// BAD: Returning the massive interface
func NewStore() IDataStore {
	return &MySQLStore{}
}
```

### The Go Pattern (Accept Interfaces, Return Structs)
```go
package store

// GOOD: Return the concrete struct
type MySQLStore struct{}

func NewStore() *MySQLStore {
	return &MySQLStore{}
}

func (m *MySQLStore) SaveUser() {}
func (m *MySQLStore) DeleteUser() {}
// ... other methods
```
```go
package api

// GOOD: The consumer defines exactly what it needs locally
type UserSaver interface {
	SaveUser()
}

// GOOD: Accept the small interface
func RegisterUser(s UserSaver) {
	s.SaveUser() // We don't care if it's MySQL, Postgres, or a Mock!
}
```

---

# Intermediate Example

Mocking becomes incredibly easy when you accept small interfaces.

```go
package main

import "fmt"

// 1. The Consumer defines what it needs
type Notifier interface {
	SendEmail(address, message string) error
}

// 2. The Consumer accepts the interface
func WelcomeNewUser(n Notifier, email string) {
	fmt.Println("Registering user...")
	n.SendEmail(email, "Welcome to our app!")
}

// --- IN PRODUCTION ---

type SendGridClient struct { APIKey string }
func (s *SendGridClient) SendEmail(addr, msg string) error {
	fmt.Printf("Sending real email to %s via SendGrid...\n", addr)
	return nil
}
func NewSendGrid() *SendGridClient { return &SendGridClient{} }

// --- IN YOUR TESTS ---

// A mock struct for testing that doesn't actually send emails
type MockNotifier struct {
	EmailsSent int
}
func (m *MockNotifier) SendEmail(addr, msg string) error {
	m.EmailsSent++
	return nil
}

func main() {
	// Production
	prodClient := NewSendGrid()
	WelcomeNewUser(prodClient, "user@example.com")

	// Testing
	mockClient := &MockNotifier{}
	WelcomeNewUser(mockClient, "test@example.com")
	fmt.Printf("Mock sent %d emails.\n", mockClient.EmailsSent)
}
```

---

# Advanced Example

Using Interface Composition to accept exactly what you need. Sometimes a function needs to Read *and* Write, but another function only needs to Read.

```go
package main

import "fmt"

// Concrete Struct
type FileSystem struct{}
func (fs *FileSystem) Read(path string) string { return "data" }
func (fs *FileSystem) Write(path, data string) {}
func (fs *FileSystem) Delete(path string) {}
func NewFileSystem() *FileSystem { return &FileSystem{} }

// Consumer Interfaces
type Reader interface {
	Read(path string) string
}

type Writer interface {
	Write(path, data string)
}

// Composition: ReadWriter requires BOTH methods
type ReadWriter interface {
	Reader
	Writer
}

// Function 1 only needs to read. We restrict its power.
func PrintFile(r Reader) {
	fmt.Println(r.Read("test.txt"))
	// r.Delete() // Compile error! It can't delete!
}

// Function 2 needs to read and write.
func CopyFile(rw ReadWriter) {
	data := rw.Read("a.txt")
	rw.Write("b.txt", data)
}

func main() {
	fs := NewFileSystem()
	
	// 'fs' implicitly satisfies Reader, Writer, and ReadWriter!
	PrintFile(fs) 
	CopyFile(fs)
}
```

---

# Production Use Cases

### 1. The `io.Reader` and `io.Writer`
The Go standard library is built entirely on this principle. `os.Open()` returns a concrete `*os.File` struct. However, the `json.NewDecoder()` function accepts an `io.Reader` interface. Because `*os.File` has a `Read()` method, it implicitly satisfies `io.Reader`. This allows the JSON decoder to read from files, network sockets, or memory buffers seamlessly.

### 2. Dependency Injection in HTTP APIs
When building a REST API, your HTTP Handlers need a Database connection. If you pass a concrete `*sql.DB` to your handlers, you can never unit-test them without a real database. By defining an interface (e.g., `UserRepository`) in your `handlers` package, you can inject a mock database during `go test`.

---

# Performance Analysis

Calling a method on an interface is slightly slower than calling a method on a concrete struct due to dynamic dispatch (the runtime has to look up a table to find the actual function pointer).
However, this overhead is measured in single-digit nanoseconds. It is entirely negligible compared to database queries or network requests. The architectural decoupling gained by using interfaces massively outweighs the nanosecond performance cost.

---

# Best Practices

* **Keep interfaces small**: Interfaces with 1 or 2 methods are ideal. Rob Pike says, "The bigger the interface, the weaker the abstraction."
* **Define interfaces where they are used**: Don't define a `UserService` interface in the `db` package. Define it in the `http` package where the HTTP handler actually calls it.
* **Return concrete structs**: Your constructors (`New()`) should return `*ConcreteStruct`. If you return an interface, you force everyone to use your exact abstraction.

---

# Common Mistakes

### The Preemptive Interface
```go
// BAD: The developer creates an interface before they even have a second implementation, "just in case".
package auth

type Authenticator interface {
    Login()
}
type CognitoAuth struct{}
func NewAuth() Authenticator { return &CognitoAuth{} } // Returning interface!
```
*Fix: Do not create an interface until you actually have two concrete types that need it, or until you need to mock it in a test. YAGNI (You Aren't Gonna Need It).*

---

# Debugging Guide

* **"Cannot use 'x' as type 'y' in argument to 'z': missing method"**: You passed a struct into a function that requires an interface, but your struct is missing one of the required methods, or the method signature (arguments/return types) is slightly wrong.
* **Pointer vs Value Receiver**: If you define `func (m *MyStruct) Read()`, the *pointer* `*MyStruct` satisfies the interface, but the *value* `MyStruct` does not! Always pass the pointer (e.g., `&MyStruct{}`).

---

# Exercises

## Beginner
Create a concrete struct `Dog` with a method `Speak() string` that returns `"Woof"`. 
Create a function `MakeSound(??? speaker)` that accepts an interface. Pass a `*Dog` into it and print the sound.

## Intermediate
Examine the standard library `io.Writer` interface. Create a concrete struct `FakeDB` with a `Write(p []byte) (n int, err error)` method. Pass it to `fmt.Fprintln()`. (Yes, `fmt.Fprintln` can write directly to your database struct!).

---

# Quiz

## Multiple Choice Questions
**1. Where is the most idiomatic place to define an interface in Go?**
A) In the package that implements the methods (Producer).
B) In the package that calls the methods (Consumer).
C) In a global `interfaces` folder.
*Answer*: B

## True or False
**Constructors (`New...` functions) in Go should generally return interfaces to hide the implementation details.**
*Answer*: False. They should return concrete structs ("return structs"). Returning an interface restricts the consumer and goes against the proverb.

---

# Interview Questions

## Beginner
**Q**: Explain "Accept Interfaces, Return Structs."
*Answer*: It means functions and methods should take interfaces as parameters (making them flexible and mockable), but factory functions should return concrete types (structs), allowing the caller to decide how they want to abstract the returned object.

## Intermediate
**Q**: What is the difference between explicit interfaces (Java) and implicit interfaces (Go)?
*Answer*: In Java, a class must explicitly state `implements InterfaceName`. This tightly couples the class to the interface definition. In Go, interfaces are satisfied implicitly just by having the matching methods. This allows developers to define interfaces *after* the concrete structs have already been written, without modifying the original struct's code.

## Advanced
**Q**: If "Accept interfaces, return structs" is the rule, why does the standard library `error` type violate this? (e.g., `errors.New()` returns an `error` interface, not a struct).
*Answer*: `error` is a deliberate exception. It is a built-in, globally defined interface that the entire language agrees upon. For standardizing fundamental concepts (like `error` or `io.Reader`), returning interfaces is sometimes necessary to guarantee cross-package compatibility. However, for application-level business logic, returning structs is the correct pattern.

---

# Cheat Sheet

* **The Proverb**: Accept Interfaces (Consumer side), Return Structs (Producer side).
* **Consumer Side**: `func DoWork(doer TheInterface) { doer.Do() }`
* **Producer Side**: `func NewDoer() *ConcreteDoer { return &ConcreteDoer{} }`
* **Small Interfaces**: `type Reader interface { Read() []byte }`

---

# Summary

"Accept interfaces, return structs" is the secret to writing decoupled, highly testable Go code. By abandoning the massive, monolithic interfaces of traditional OOP and embracing small, consumer-defined interfaces, your Go programs will remain flexible and easy to refactor as they grow to massive scales.

---

# Key Takeaways

* ✔ Interfaces in Go are satisfied implicitly (no `implements` keyword).
* ✔ Consumers should define exactly what methods they need in small interfaces.
* ✔ Producers should return concrete structs.
* ✔ This pattern is the key to easy unit-testing and mocking.

---

# Further Reading
* [Go Proverbs (Video)](https://www.youtube.com/watch?v=PAAkCSZUG1c)
* [SOLID Go Design](https://dave.cheney.net/2015/08/20/solid-go-design)

---

# Next Chapter
➡️ **Next:** `03-Context-Propagation.md`
