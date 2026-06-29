# Error Wrapping and Typing

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

In Go, errors are just values (specifically, they are any type that implements the `error` interface, which has a single method: `Error() string`). This makes error handling explicit and highly visible. 

However, as an error bubbles up the call stack from the database layer, to the service layer, to the HTTP layer, it loses context. A simple `"connection refused"` error doesn't tell you *what* failed to connect. **Error Wrapping** allows you to attach context at each layer (e.g., `"failed to fetch user: failed to connect to database: connection refused"`) while preserving the original error so that the top layer can still inspect it.

---

# Learning Objectives

After completing this chapter you will be able to:

* Wrap errors to provide contextual stack traces.
* Use `errors.Is` to check for specific sentinel errors.
* Use `errors.As` to extract structured data from custom error types.
* Design a robust error-handling strategy for a layered application.

---

# Prerequisites

Before reading this chapter you should know:

* Basic Go error handling (`if err != nil`).
* Structs and Interfaces.

---

# Why This Topic Exists

Prior to Go 1.13, developers used external libraries (like `pkg/errors`) to wrap errors. They would concatenate error strings using `fmt.Errorf("user error: %v", err)`. But string concatenation destroys the original error type! If the database returned a `sql.ErrNoRows`, concatenating it into a string means the HTTP handler can no longer check `if err == sql.ErrNoRows` to return a 404.

Go 1.13 introduced the `%w` verb in `fmt.Errorf` and the `errors.Is`/`errors.As` functions, creating a standardized, built-in design pattern for creating layered, inspectable errors.

---

# Real-World Analogy

### The Russian Nesting Doll (Matryoshka)

An error in Go is like a Russian Nesting Doll.
* **The Core (Original Error)**: The innermost doll is the root cause (e.g., `sql.ErrNoRows`).
* **The Wrappers**: As the doll is passed up the chain, each person puts it inside a slightly larger doll with a label (e.g., `"failed to query database"`, then `"failed to fetch user"`).
* **Inspection**: The person at the very top receives the largest doll. They can read the outer label, but if they want to know the *exact* cause, they must open all the dolls (using `errors.Is` or `errors.As`) to inspect the tiny doll at the center.

---

# Core Concepts

* **Sentinel Errors**: Predefined, package-level error variables (e.g., `var ErrNotFound = errors.New("not found")`). Used for simple equality checks.
* **Custom Error Types**: A struct that implements the `error` interface, allowing you to attach metadata (e.g., HTTP status codes).
* **Wrapping (`%w`)**: Using `fmt.Errorf("context: %w", err)` to wrap an error inside another error.
* **Unwrapping**: `errors.Unwrap(err)` peels off one layer.
* **errors.Is**: Checks if *any* error in the chain matches a specific Sentinel Error.
* **errors.As**: Checks if *any* error in the chain matches a specific Custom Error Type, and if so, extracts it.

---

# Architecture Diagram

```mermaid
flowchart TD
    DB[Database Layer<br/>Returns sql.ErrNoRows]
    Service[Service Layer<br/>Wraps: fmt.Errorf]
    HTTP[HTTP Layer<br/>Inspects with errors.Is]
    
    DB -- "sql.ErrNoRows" --> Service
    Service -- "failed to get user: %w" --> HTTP
    
    note right of HTTP: errors.Is(err, sql.ErrNoRows)<br/>Returns TRUE -> Sends HTTP 404
```

---

# Step-by-Step Implementation

1. **At the source**: Return a base error (either a sentinel error or a custom struct).
2. **At intermediate layers**: Catch the error. Do not return it raw. Wrap it with context using `fmt.Errorf("doing X failed: %w", err)`.
3. **At the top layer (Consumer)**: Receive the heavily wrapped error.
4. **Inspection**: Use `errors.Is(err, ErrSpecific)` to check for logic flow, or `errors.As(err, &myStruct)` to extract data.
5. **Response**: Log the full wrapped string, and return a sanitized response to the user.

---

# Syntax

```go
import (
    "errors"
    "fmt"
)

var ErrNotFound = errors.New("item not found")

// Wrap an error
err := fmt.Errorf("fetching failed: %w", ErrNotFound)

// Check if a specific error exists anywhere in the chain
if errors.Is(err, ErrNotFound) {
    // Handle 404
}
```

---

# Beginner Example

Wrapping an error and checking it with `errors.Is`.

```go
package main

import (
	"errors"
	"fmt"
)

// 1. Define a Sentinel Error
var ErrDatabaseDown = errors.New("database connection refused")

// The lowest layer
func connectDB() error {
	// Something goes wrong
	return ErrDatabaseDown
}

// The middle layer
func fetchUser() error {
	err := connectDB()
	if err != nil {
		// 2. Wrap the error using %w
		return fmt.Errorf("fetchUser failed: %w", err)
	}
	return nil
}

func main() {
	err := fetchUser()
	
	// The full error string contains all layers
	fmt.Println("Log:", err) 
	// Output: Log: fetchUser failed: database connection refused

	// 3. Inspect the chain using errors.Is
	if errors.Is(err, ErrDatabaseDown) {
		fmt.Println("Action: Alert the Ops team, DB is down!")
	} else {
		fmt.Println("Action: Return 500 Internal Error")
	}
}
```

---

# Intermediate Example

Using **Custom Error Types** and `errors.As`. Sentinel errors are just strings. If you want to attach metadata (like an HTTP status code), you must build a custom error struct.

```go
package main

import (
	"errors"
	"fmt"
)

// 1. Define a Custom Error Type
type HTTPError struct {
	StatusCode int
	Message    string
}

// It must implement the error interface
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

func performAction() error {
	// Create the custom error
	rootErr := &HTTPError{StatusCode: 403, Message: "user lacks admin role"}
	
	// Wrap it!
	return fmt.Errorf("action aborted: %w", rootErr)
}

func main() {
	err := performAction()

	// 2. Prepare a variable of the target type
	var httpErr *HTTPError

	// 3. Use errors.As to search the chain. 
	// If found, it populates the httpErr variable and returns true!
	if errors.As(err, &httpErr) {
		fmt.Printf("Detected HTTP Error!\n")
		fmt.Printf("Send status code %d to user.\n", httpErr.StatusCode)
	} else {
		fmt.Println("Unknown error")
	}
}
```

---

# Advanced Example

Joining multiple errors. In Go 1.20, `errors.Join` was introduced to combine multiple independent errors into a single error. This is incredibly useful for validating a struct and returning all validation failures at once, rather than stopping at the first one.

```go
package main

import (
	"errors"
	"fmt"
)

func validateForm(username, password string) error {
	var errs error

	if len(username) < 3 {
		// Join appends the new error to the chain
		errs = errors.Join(errs, errors.New("username too short"))
	}
	if len(password) < 8 {
		errs = errors.Join(errs, errors.New("password too short"))
	}

	return errs
}

func main() {
	err := validateForm("al", "pass")
	
	if err != nil {
		fmt.Println("Validation Failed:")
		// Printing the joined error automatically puts each error on a new line!
		fmt.Println(err)
	}
}
```

---

# Production Use Cases

### 1. The Global Error Handler
In a REST API, you don't want to write `if err == sql.ErrNoRows { write(404) }` in every single HTTP handler. 
Instead, handlers return a wrapped error up to a centralized Global Middleware. The Middleware uses `errors.As` to check if the error is a custom `APIError` containing a status code. If it is, it extracts the code and sends it. If it isn't, it defaults to sending a 500 Internal Server Error.

### 2. Retries and Idempotency
If a background worker makes a network request and fails, it inspects the error using `errors.Is`. If `errors.Is(err, ErrNetworkTimeout)` is true, the worker knows it is safe to retry. If `errors.Is(err, ErrUnauthorized)` is true, the worker aborts instantly because retrying will never fix an expired token.

---

# Performance Analysis

Wrapping errors and using `errors.Is`/`errors.As` is slightly slower than direct equality checks (`err == ErrNotFound`) because it requires iterating through the linked list of wrapped errors and using reflection (`errors.As`). However, the performance cost is negligible outside of microscopic, high-frequency tight loops.

---

# Best Practices

* **Use `%w`**: Never use `%v` in `fmt.Errorf` unless you intentionally want to *destroy* the error chain (sometimes done for security to prevent internal DB errors from leaking to the frontend).
* **Don't repeat yourself**: If the innermost error says "database query failed", don't wrap it with "failed to query database". Wrap it with *context*: "failed to update user profile: %w".
* **Export Sentinel Errors**: If you expect callers of your package to inspect your errors, define them as exported variables (`var ErrInvalidInput = errors.New(...)`).

---

# Common Mistakes

### Using `==` instead of `errors.Is`
```go
err := fmt.Errorf("wrap: %w", sql.ErrNoRows)

// BAD: This is FALSE because 'err' is a wrapped object, not the exact memory address of sql.ErrNoRows!
if err == sql.ErrNoRows { } 

// GOOD: This searches the entire chain and correctly returns TRUE.
if errors.Is(err, sql.ErrNoRows) { }
```

### Passing a value to `errors.As` instead of a pointer
```go
var customErr CustomErrorStruct // Value!
// BAD: Will panic! errors.As requires a pointer to the target variable so it can mutate it.
errors.As(err, &customErr) 

var customErrPtr *CustomErrorStruct // Pointer!
// GOOD
errors.As(err, &customErrPtr)
```

---

# Debugging Guide

* **`errors.Is` returns false unexpectedly**: You probably wrapped the error using `%v` instead of `%w`. Change `fmt.Errorf("...: %v", err)` to `%w`.
* **`errors.As` panics at runtime**: You passed a non-pointer to `errors.As`. Ensure the second argument is a pointer to the variable you want populated.

---

# Exercises

## Beginner
Create a sentinel error `ErrPermissionDenied`. Write a function that returns this error wrapped inside another error using `%w`. In `main`, verify that `errors.Is` successfully detects `ErrPermissionDenied`.

## Intermediate
Create a custom error struct `ValidationError` containing a `Field` string and a `Reason` string. Write a function that returns this error wrapped. In `main`, use `errors.As` to extract the `Field` and `Reason` and print them.

---

# Quiz

## Multiple Choice Questions
**1. What is the difference between `%v` and `%w` in `fmt.Errorf`?**
A) `%w` adds a newline.
B) `%v` formats as a string and destroys the original error type. `%w` formats as a string but preserves the original error inside a linked list for `errors.Is` to find.
C) `%w` is faster.
*Answer*: B

## True or False
**You should define an `error` as a struct when you need to attach extra data fields (like HTTP codes) to the error.**
*Answer*: True. Sentinel errors (`errors.New`) are just static strings. Custom structs implementing the `error` interface allow rich metadata.

---

# Interview Questions

## Beginner
**Q**: Why is `errors.Is(err, target)` better than `err == target`?
*Answer*: `err == target` only works if the error has never been wrapped. `errors.Is` unwraps the error layer by layer, checking every nested error to see if any of them match the target.

## Intermediate
**Q**: How do you implement a custom error type in Go?
*Answer*: You define a struct, and then you attach an `Error() string` method to that struct. Because it implements the `error` interface implicitly, it can now be returned anywhere an `error` is expected.

## Advanced
**Q**: Explain how `errors.As` works under the hood.
*Answer*: `errors.As` takes the error chain and the target pointer. It uses Go reflection to determine the type of the target pointer. It then loops through the error chain, checking if each error can be type-asserted to the target type. If it finds a match, it assigns the matched error to the target pointer and returns true.

---

# Cheat Sheet

* **Sentinel Error**: `var ErrBoom = errors.New("boom")`
* **Wrapping**: `fmt.Errorf("context: %w", err)`
* **Checking Sentinel**: `if errors.Is(err, ErrBoom) {}`
* **Extracting Custom**:
```go
var target *MyCustomError
if errors.As(err, &target) {
    fmt.Println(target.MyCustomField)
}
```
* **Joining Errors**: `errs = errors.Join(err1, err2)`

---

# Summary

Error handling in Go is not just an afterthought; it is a core architectural design pattern. By wrapping errors with `%w` and inspecting them with `errors.Is` and `errors.As`, you can maintain highly detailed logs (for developers) while simultaneously providing structured, programmatic error routing (for your application logic).

---

# Key Takeaways

* ✔ Use `%w` to wrap errors and add context without destroying the root cause.
* ✔ Use `errors.Is` to check for specific Sentinel Errors.
* ✔ Use `errors.As` to extract Custom Error structs.
* ✔ Use `errors.Join` to accumulate multiple validation errors.

---

# Further Reading
* [Go Blog: Working with Errors in Go 1.13](https://go.dev/blog/go1.13-errors)

---

# Next Chapter
➡️ **Next:** `05-Factory-Method.md`
