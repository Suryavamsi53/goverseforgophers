# The Functional Options Pattern

In languages like Java or Python, if a class has 10 optional configuration parameters, you solve it using **Method Overloading** or **Default Arguments** (e.g., `def create_server(host, port=8080, timeout=5)`).

Go does not support method overloading, and it does not support default arguments. 

If you want to create a `Server` struct with optional `Timeout` and `TLS` settings, the naïve approach is a massive, ugly config struct:

```go
// The Ugly Way
type ServerConfig struct {
    Host    string
    Port    int
    Timeout time.Duration
    TLS     bool
}

func NewServer(cfg ServerConfig) *Server { ... }

// Usage requires passing a bunch of empty defaults:
srv := NewServer(ServerConfig{Host: "localhost", Port: 80})
```

To solve this elegantly, Rob Pike (co-creator of Go) popularized the **Functional Options Pattern**.

## 1. The Implementation

The core idea is that you pass a variadic list of *functions* to the constructor. Each function takes a pointer to the configuration struct and modifies a single field.

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

// 1. Define the Option signature
type Option func(*Server)

// 2. The Constructor takes a variadic list of Options
func NewServer(opts ...Option) *Server {
    // Set absolute defaults
    s := &Server{
        host:    "127.0.0.1",
        port:    8080,
        timeout: 30 * time.Second,
    }

    // Apply all provided options
    for _, opt := range opts {
        opt(s)
    }

    return s
}
```

## 2. Defining the Options

You define functions that return an `Option` closure.

```go
func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(t time.Duration) Option {
    return func(s *Server) {
        s.timeout = t
    }
}
```

## 3. The Usage (Beautiful APIs)

Now, the API for your library is perfectly clean and self-documenting.

```go
// Default server
srv1 := NewServer()

// Custom server
srv2 := NewServer(
    WithHost("0.0.0.0"),
    WithPort(443),
    WithTimeout(5 * time.Second),
)
```

## 4. Why use this in Production?

1. **Backwards Compatibility**: If you add a new `WithMaxConnections` option in v2.0 of your library, none of the v1.0 code breaks because the `NewServer()` variadic signature hasn't changed.
2. **Encapsulation**: The internal `Server` fields remain private (lowercase). External packages cannot mutate the server after it is built.
3. **Immutability**: Once the options are applied and the constructor returns, the struct is finalized.

The Functional Options Pattern is the gold standard for Go library design, heavily used in the standard library and popular open-source tools (like gRPC's `grpc.DialOption`).
