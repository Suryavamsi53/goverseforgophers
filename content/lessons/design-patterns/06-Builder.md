# Builder Pattern

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

The **Builder Pattern** is a Creational Design Pattern that lets you construct complex objects step by step. It allows you to produce different types and representations of an object using the same construction code.

In Go, the Builder pattern is heavily overshadowed by the **Functional Options** pattern for configuring structs. However, the Builder pattern still shines brightly in specific scenarios: when the construction process requires complex validation, order-dependent steps, or when constructing abstract syntax trees, complex database queries, or nested JSON payloads.

---

# Learning Objectives

After completing this chapter you will be able to:

* Differentiate between the Builder Pattern and Functional Options.
* Implement a fluent, chainable API (Method Chaining).
* Use the Director pattern to encapsulate common build sequences.
* Safely construct complex objects while ensuring all required fields are validated.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Methods.
* Pointers (to modify the builder's state).
* Functional Options (`01-Functional-Options.md`).

---

# Why This Topic Exists

If a struct has 15 fields, a standard constructor function becomes unreadable: `NewHouse(4, 2, true, false, "red", "tile")`. 
While Functional Options solve this nicely, what if the creation of a `House` requires strict validation? (e.g., "A house cannot have a roof if it doesn't have walls"). 

A Builder allows you to define a dedicated `HouseBuilder` struct. It accumulates state step-by-step through method calls (`builder.AddWalls().AddRoof()`). Finally, you call a `Build()` method. The `Build()` method performs all complex validations at once and returns the finalized `House` object (or an error).

---

# Real-World Analogy

### Custom Computer Assembly

* **The Object**: A custom Gaming PC.
* **The Builder**: The website configurator tool. You select the CPU, then the Motherboard, then the RAM. 
* **Validation**: If you select an Intel CPU, the Builder restricts your Motherboard choices to only Intel-compatible ones. You cannot proceed if the build is invalid.
* **The Build Method**: Once you are finished, you click "Checkout (Build)". The configurator does a final compatibility check and produces the final receipt (the Object).

---

# Core Concepts

* **The Product**: The complex object being built (e.g., `House`, `SQLQuery`).
* **The Builder**: A struct with methods to configure parts of the Product.
* **Method Chaining**: Builder methods return the Builder instance itself (`*Builder`), allowing calls to be chained: `b.Step1().Step2().Build()`.
* **The Build Method**: The final method that validates the accumulated state and returns the finalized Product.
* **The Director (Optional)**: A struct that accepts a Builder and executes a predefined series of steps to construct a popular version of the Product.

---

# Architecture Diagram

```mermaid
flowchart LR
    Client[Client Code]
    Builder[Builder Struct]
    Product[Complex Product]
    Director[Director (Optional)]

    Client -- "1. Init Builder" --> Builder
    Client -- "2. .SetA().SetB()" --> Builder
    Client -- "3. .Build()" --> Builder
    Builder -- "Validates & Returns" --> Product
    
    Director -. "Automates steps<br/>for the Client" .-> Builder
```

---

# Step-by-Step Implementation

1. Define the complex **Product** struct.
2. Define the **Builder** struct containing fields necessary to track the configuration.
3. Write setter methods on the Builder. Have them return `*Builder` to allow method chaining.
4. Write a `Build() (*Product, error)` method on the Builder.
5. Inside `Build()`, perform all necessary validation logic. If valid, instantiate and return the Product.

---

# Syntax

```go
type Builder struct {
    name string
}

// Return the pointer to allow chaining
func (b *Builder) SetName(n string) *Builder {
    b.name = n
    return b
}

// The final construction step
func (b *Builder) Build() (*Product, error) {
    if b.name == "" { return nil, errors.New("name required") }
    return &Product{Name: b.name}, nil
}
```

---

# Beginner Example

Building a SQL Query. This is where the Builder pattern shines in Go.

```go
package main

import (
	"fmt"
	"strings"
)

// The Builder
type QueryBuilder struct {
	table   string
	columns []string
	where   []string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// Chaining methods
func (q *QueryBuilder) Select(cols ...string) *QueryBuilder {
	q.columns = append(q.columns, cols...)
	return q
}

func (q *QueryBuilder) From(table string) *QueryBuilder {
	q.table = table
	return q
}

func (q *QueryBuilder) Where(condition string) *QueryBuilder {
	q.where = append(q.where, condition)
	return q
}

// The Final Build Method
func (q *QueryBuilder) Build() (string, error) {
	if q.table == "" {
		return "", fmt.Errorf("table name is required")
	}
	if len(q.columns) == 0 {
		q.columns = append(q.columns, "*") // Default to all
	}

	query := fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.columns, ", "), q.table)

	if len(q.where) > 0 {
		query += " WHERE " + strings.Join(q.where, " AND ")
	}

	return query + ";", nil
}

func main() {
	// Fluent API (Method Chaining)
	query, err := NewQueryBuilder().
		Select("id", "name", "email").
		From("users").
		Where("age > 18").
		Where("status = 'active'").
		Build()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(query)
	// Output: SELECT id, name, email FROM users WHERE age > 18 AND status = 'active';
}
```

---

# Intermediate Example

The Director pattern. Imagine we are building HTTP Responses. A Director can quickly assemble common responses using our Builder.

```go
package main

import "fmt"

// 1. The Product
type HTTPResponse struct {
	Code int
	Body string
}

// 2. The Builder
type ResponseBuilder struct {
	response *HTTPResponse
}

func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{response: &HTTPResponse{}}
}

func (b *ResponseBuilder) SetCode(code int) *ResponseBuilder {
	b.response.Code = code
	return b
}

func (b *ResponseBuilder) SetBody(body string) *ResponseBuilder {
	b.response.Body = body
	return b
}

func (b *ResponseBuilder) Build() *HTTPResponse {
	return b.response
}

// 3. The Director (Encapsulates common recipes)
type Director struct {
	builder *ResponseBuilder
}

func NewDirector(b *ResponseBuilder) *Director {
	return &Director{builder: b}
}

// Recipe 1
func (d *Director) Construct404() *HTTPResponse {
	return d.builder.SetCode(404).SetBody("Not Found").Build()
}

// Recipe 2
func (d *Director) Construct500() *HTTPResponse {
	return d.builder.SetCode(500).SetBody("Internal Server Error").Build()
}

func main() {
	builder := NewResponseBuilder()
	director := NewDirector(builder)

	// Using the director for standard builds
	resp1 := director.Construct404()
	fmt.Printf("Response: %d - %s\n", resp1.Code, resp1.Body)

	// Using the builder manually for a custom build
	resp2 := NewResponseBuilder().SetCode(200).SetBody("{\"status\":\"ok\"}").Build()
	fmt.Printf("Response: %d - %s\n", resp2.Code, resp2.Body)
}
```

---

# Advanced Example

Using interfaces for the Builder to allow swapping implementations. Let's say we want to build a Document, but we want the option to output either HTML or Markdown.

```go
package main

import "fmt"

// 1. The Builder Interface
type DocumentBuilder interface {
	SetTitle(title string) DocumentBuilder
	AddParagraph(text string) DocumentBuilder
	Build() string
}

// 2. Concrete Builder A: HTML
type HTMLBuilder struct {
	content string
}
func (h *HTMLBuilder) SetTitle(t string) DocumentBuilder {
	h.content += fmt.Sprintf("<h1>%s</h1>\n", t)
	return h
}
func (h *HTMLBuilder) AddParagraph(t string) DocumentBuilder {
	h.content += fmt.Sprintf("<p>%s</p>\n", t)
	return h
}
func (h *HTMLBuilder) Build() string { return h.content }

// 3. Concrete Builder B: Markdown
type MarkdownBuilder struct {
	content string
}
func (m *MarkdownBuilder) SetTitle(t string) DocumentBuilder {
	m.content += fmt.Sprintf("# %s\n", t)
	return m
}
func (m *MarkdownBuilder) AddParagraph(t string) DocumentBuilder {
	m.content += fmt.Sprintf("%s\n\n", t)
	return m
}
func (m *MarkdownBuilder) Build() string { return m.content }

// 4. The Director
func CreateStandardDoc(b DocumentBuilder) string {
	// The Director doesn't care if it's HTML or Markdown!
	return b.SetTitle("Welcome").
		AddParagraph("This is the first paragraph.").
		Build()
}

func main() {
	htmlBuilder := &HTMLBuilder{}
	fmt.Println("--- HTML ---")
	fmt.Print(CreateStandardDoc(htmlBuilder))

	mdBuilder := &MarkdownBuilder{}
	fmt.Println("--- MARKDOWN ---")
	fmt.Print(CreateStandardDoc(mdBuilder))
}
```

---

# Production Use Cases

### 1. ORMs and SQL Query Builders
Popular Go libraries like `squirrel` (a SQL query builder) use this pattern exclusively. Because SQL queries require a strict grammar (SELECT must come before FROM), the Builder pattern allows the library to accumulate the parts of the query and validate the syntax during the final `ToSql()` call.

### 2. Complex Test Fixtures
When writing unit tests for systems with complex nested data structures (like a Kubernetes Custom Resource Definition), developers often write a `PodBuilder` to easily construct valid Mock objects for testing without writing massive inline JSON/Struct literals.

---

# Performance Analysis

The Builder pattern requires creating a separate Builder struct before creating the actual Product struct. This slightly increases memory allocations. However, because Builders are usually allocated on the stack (thanks to escape analysis, if the builder doesn't escape the function), the performance penalty is practically zero in Go.

---

# Best Practices

* **Return Pointers for Chaining**: Always return `*Builder` from your configuration methods so the user can chain them (`.MethodA().MethodB()`).
* **Delay Initialization**: Do not instantiate the heavy Product struct until the `Build()` method is called. Accumulate the state in the Builder struct first. This saves memory if the build process fails validation halfway through.
* **Immutable Results**: Once `Build()` is called, the returned Product should ideally be immutable (no exported fields) to prevent tampering after the complex validation has passed.

---

# Common Mistakes

### Confusing Builder with Functional Options
* **Use Functional Options** when you just need to configure a simple struct (e.g., setting a timeout on an HTTP Client).
* **Use Builder** when the construction process requires multiple steps, strict ordering, or complex cross-field validation (e.g., if `A` is set, `B` must also be set). Using Functional Options for cross-field validation is incredibly messy.

---

# Debugging Guide

* **"Panic: nil pointer dereference"**: You probably forgot to return `*Builder` from one of your setter methods, or you forgot to initialize the underlying map/slice inside the `NewBuilder()` constructor.

---

# Exercises

## Beginner
Create a `BurgerBuilder`. Give it methods `AddCheese()`, `AddPatty()`, and `AddLettuce()`. Add a `Build() string` method that returns a string describing the burger. Use method chaining to create a double cheeseburger.

## Intermediate
Add validation to your `BurgerBuilder`. In the `Build() (*Burger, error)` method, return an error if the user tries to build a burger with 0 patties. Test the error handling.

---

# Quiz

## Multiple Choice Questions
**1. What is the primary benefit of having the builder methods return a pointer to the builder (`*Builder`)?**
A) It uses less memory.
B) It allows for Method Chaining (a Fluent API) like `b.SetX().SetY()`.
C) It prevents data races.
*Answer*: B

## True or False
**In Go, the Builder Pattern is the best way to configure an HTTP Server struct with an optional port and timeout.**
*Answer*: False. For simple, unordered optional configuration, the Functional Options pattern (`01-Functional-Options.md`) is significantly more idiomatic in Go. Builder is for complex, step-by-step construction.

---

# Interview Questions

## Beginner
**Q**: What is the purpose of the `Build()` method in the Builder Pattern?
*Answer*: The `Build()` method signals the end of the configuration steps. It is responsible for validating all the accumulated data inside the Builder, instantiating the final Product object, and returning it (along with any validation errors).

## Intermediate
**Q**: When would you choose to use the Builder Pattern over the Functional Options pattern in Go?
*Answer*: Functional Options are great for unordered, independent configurations. I would choose the Builder pattern when constructing an object requires strict validation based on the combination of fields, when the construction requires a specific sequence of steps, or when building a grammar/tree structure (like a SQL query or XML document).

## Advanced
**Q**: What is the purpose of the Director in the Builder pattern?
*Answer*: The Director is a struct or function that takes a generic Builder interface as an argument. It encapsulates predefined recipes or sequences of builder calls. This hides the complex construction steps from the client, allowing the client to simply ask the Director for a "Standard" or "Premium" object, while still allowing the underlying Builder implementation to be swapped out (e.g., generating HTML vs Markdown).

---

# Cheat Sheet

* **The Builder**:
```go
type QueryBuilder struct { table string }

// Constructor
func NewQB() *QueryBuilder { return &QueryBuilder{} }

// Chaining Method
func (q *QueryBuilder) Table(name string) *QueryBuilder {
    q.table = name
    return q
}

// Final Step
func (q *QueryBuilder) Build() (string, error) {
    if q.table == "" { return "", errors.New("missing table") }
    return "SELECT * FROM " + q.table, nil
}
```

---

# Summary

While Functional Options dominate simple struct configuration in Go, the Builder pattern remains the undisputed champion for constructing complex, highly-validated objects and Domain-Specific Languages (like SQL or HTML generators). By enabling method chaining, Builders provide a fluent, readable API that makes complex instantiation a joy.

---

# Key Takeaways

* ✔ Use Builders for complex, multi-step object creation.
* ✔ Return `*Builder` from setter methods to enable method chaining.
* ✔ Perform all cross-field validation inside the final `Build()` method.
* ✔ Use a Director to encapsulate common construction recipes.

---

# Further Reading
* [Refactoring.guru: Builder Pattern](https://refactoring.guru/design-patterns/builder)

---

# Next Chapter
➡️ **Next:** `07-Singleton.md`
