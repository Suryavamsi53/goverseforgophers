# Go Design Patterns Curriculum

Welcome to the **Design Patterns in Go** module. In this curriculum, we will explore how to apply classic Gang of Four (GoF) design patterns to Go, as well as discover Go-specific idioms that leverage its unique features like implicit interfaces, first-class functions, and concurrency.

Go is not a traditional object-oriented language. It lacks classes and inheritance. Therefore, implementing classic design patterns often looks very different (and simpler!) in Go.

## Curriculum Overview

### Part 1: Go-Specific Idioms (The "Go Way")
* **01-Functional-Options.md**: The idiomatic way to handle complex, optional configurations using higher-order functions.
* **02-Accept-Interfaces-Return-Structs.md**: The golden rule of Go interfaces and dependency injection.
* **03-Context-Propagation.md**: Managing request-scoped data, timeouts, and cancellation signals cleanly.
* **04-Error-Wrapping-and-Typing.md**: Designing robust error handling APIs using `errors.Is` and `errors.As`.

### Part 2: Creational Patterns
* **05-Factory-Method.md**: Creating instances of structs without exposing the underlying implementation.
* **06-Builder.md**: Constructing complex objects step-by-step (often superseded by Functional Options in Go).
* **07-Singleton.md**: Enforcing a single instance globally (using `sync.Once`).
* **08-Object-Pool.md**: Reusing expensive allocations to reduce Garbage Collection pressure (using `sync.Pool`).

### Part 3: Structural Patterns
* **09-Adapter.md**: Bridging incompatible interfaces (very common when migrating systems or using third-party SDKs).
* **10-Decorator.md**: Dynamically adding behavior to an object (heavily used in HTTP Middlewares).
* **11-Facade.md**: Providing a simplified, unified interface to a complex subsystem.
* **12-Proxy.md**: Controlling access to an object (e.g., caching, lazy loading, or access control).

### Part 4: Behavioral Patterns
* **13-Strategy.md**: Encapsulating interchangeable algorithms (perfectly suited for Go interfaces).
* **14-Observer.md**: Implementing Pub/Sub mechanics (often implemented via Channels in Go).
* **15-Command.md**: Encapsulating requests as objects (useful for undo/redo or queuing tasks).
* **16-State.md**: Allowing an object to alter its behavior when its internal state changes.
* **17-Chain-of-Responsibility.md**: Passing a request through a chain of handlers (another common HTTP Middleware pattern).

---

Let's begin mastering Design Patterns in Go!
