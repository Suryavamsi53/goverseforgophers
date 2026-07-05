# Builder Pattern

The Builder Pattern is a creational design pattern that allows you to construct complex objects step by step. 

It is incredibly useful when an object requires a massive amount of configuration, or when the configuration requires complex validation before the object can be created.

*(Note: The Builder pattern achieves the same goal as the Functional Options pattern from Lesson 1, but uses a fluent, method-chaining approach instead of closures).*

## 1. The Implementation

The core idea is to create a separate `Builder` struct that temporarily holds the configuration. Once all the steps are chained together, you call a final `Build()` method which validates the configuration and returns the actual object.

```go
package server

import "errors"

// 1. The Complex Object we want to build
type Server struct {
    Host string
    Port int
    TLS  bool
}

// 2. The Builder Struct
type ServerBuilder struct {
    server *Server
}

// 3. The Constructor for the Builder
func NewServerBuilder() *ServerBuilder {
    return &ServerBuilder{
        server: &Server{ // Default values
            Host: "localhost",
            Port: 8080,
        },
    }
}

// 4. The Chaining Methods (Notice they return the Builder!)
func (b *ServerBuilder) Host(host string) *ServerBuilder {
    b.server.Host = host
    return b
}

func (b *ServerBuilder) Port(port int) *ServerBuilder {
    b.server.Port = port
    return b
}

func (b *ServerBuilder) EnableTLS() *ServerBuilder {
    b.server.TLS = true
    return b
}

// 5. The Final Build Method (with validation!)
func (b *ServerBuilder) Build() (*Server, error) {
    if b.server.Port < 1 || b.server.Port > 65535 {
        return nil, errors.New("invalid port number")
    }
    return b.server, nil
}
```

## 2. The Usage (Method Chaining)

Because each configuration method returns the Builder itself (`*ServerBuilder`), you can chain the methods together into a highly readable, fluent interface.

```go
func main() {
    // Build a complex server in one readable chain
    srv, err := server.NewServerBuilder().
        Host("192.168.1.100").
        Port(443).
        EnableTLS().
        Build()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Server running on %s:%d (TLS: %v)\n", srv.Host, srv.Port, srv.TLS)
}
```

## 3. Builder vs Functional Options

Both patterns solve the exact same problem: handling massive, complex constructors.

* **Use Builder when:**
  1. The construction process requires strict step-by-step validation.
  2. You want a Fluent API (method chaining) which is very common in SQL Query Builders (e.g., `db.Select().Where().Limit()`).
* **Use Functional Options when:**
  1. You are building a standard Go library (it is the idiomatic standard).
  2. The configuration is mostly setting static properties, rather than executing complex algorithmic steps.
