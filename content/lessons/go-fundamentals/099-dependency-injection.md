# Dependency Injection (Constructor Pattern)

As we learned in the Interfaces lesson, Dependency Injection (DI) is the act of passing dependencies (like database connections) into a struct, rather than hardcoding them inside the struct.

While languages like Java rely on massive frameworks (like Spring) using Reflection and `@Autowired` annotations to magically inject dependencies, **Go prefers explicit, manual injection via Constructors.**

## 1. The Constructor Pattern

Go does not have a built-in `constructor` keyword. Instead, the idiomatic pattern is to create a function named `New<Type>()`.

```go
package user

// 1. Define the dependency interface
type Database interface {
    Save(name string) error
}

// 2. Define the service
type Service struct {
    db Database // Relies on the interface, not a concrete struct!
}

// 3. The Constructor Function (Injects the dependency)
func NewService(db Database) *Service {
    return &Service{
        db: db,
    }
}
```

## 2. Wiring it together in `main.go`

Your `cmd/server/main.go` file is known as the **Composition Root**. It is the only place in your application that knows about the concrete implementations. It constructs the real database, and passes it into the services.

```go
package main

import (
    "myapp/internal/database"
    "myapp/internal/user"
)

func main() {
    // 1. Initialize concrete dependencies
    realDB := database.NewPostgresDB("localhost:5432")

    // 2. Inject them into the services
    userService := user.NewService(realDB)

    // 3. Start application
    userService.Register("Alice")
}
```

## 3. The Options Pattern (Variadic Functions)

What if `Service` has 1 mandatory dependency, but 5 optional configurations (like timeouts, retries, or custom loggers)? Passing 6 arguments to `NewService()` becomes incredibly messy.

Senior Go engineers solve this using the **Functional Options Pattern**.

```go
type Service struct {
    db      Database
    timeout time.Duration // Optional
    retries int           // Optional
}

// Define a function type that modifies the Service
type Option func(*Service)

// Option builder functions
func WithTimeout(t time.Duration) Option {
    return func(s *Service) { s.timeout = t }
}

func WithRetries(r int) Option {
    return func(s *Service) { s.retries = r }
}

// Constructor with variadic Options
func NewService(db Database, opts ...Option) *Service {
    s := &Service{
        db:      db,
        timeout: 5 * time.Second, // Default
        retries: 3,               // Default
    }
    
    // Apply all provided options
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```
Now, initializing the service is incredibly clean and readable:
```go
// Default configuration
svc := NewService(realDB)

// Custom configuration!
customSvc := NewService(realDB, WithTimeout(10*time.Second), WithRetries(5))
```
