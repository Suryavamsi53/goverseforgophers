# HTTP Server

Writing a production-grade web server in Go does not require massive frameworks like Spring Boot (Java) or Django (Python). The standard library is powerful enough to handle millions of concurrent connections natively.

## 1. The `http.Handler` Interface

Everything in Go's web ecosystem revolves around a single interface:

```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

Any struct that implements this method can serve web traffic. However, for simplicity, Go provides `http.HandleFunc`, which allows you to use standard functions instead of structs.

```go
import (
    "fmt"
    "net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
    // w is the io.Writer where we send data back to the client!
    fmt.Fprintf(w, "Welcome to the Go Web Server!")
}

func main() {
    // Register the route
    http.HandleFunc("/", homeHandler)

    // Start the server on port 8080 (Blocks forever)
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## 2. The `ServeMux` (Router)

In the example above, `ListenAndServe` uses `nil` for the handler, which defaults to the global `http.DefaultServeMux`. 

**Best Practice:** Never use the global ServeMux in production. Third-party packages you import can secretly register malicious routes to the global Mux! Always create your own isolated router.

```go
func main() {
    // Create an isolated router
    mux := http.NewServeMux()
    
    // Go 1.22+ supports HTTP Methods directly in the route string!
    mux.HandleFunc("GET /users", getUsers)
    mux.HandleFunc("POST /users", createUser)

    http.ListenAndServe(":8080", mux)
}
```

## 3. The `http.Server` Struct (Timeouts!)

Just like the `http.Client`, the basic `http.ListenAndServe` function provides **no timeouts**. 

If a malicious client (like a Slowloris attack) connects to your server and sends 1 byte of data every 10 seconds, it will tie up a Goroutine forever, eventually crashing your server.

To secure your server, you must instantiate a custom `http.Server` struct.

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /", homeHandler)

    // A secure, production-ready server configuration
    server := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,  // Max time to read the request payload
        WriteTimeout: 10 * time.Second, // Max time to process and write the response
        IdleTimeout:  120 * time.Second,// Max time to keep a TCP connection warm
    }

    server.ListenAndServe()
}
```

### 🧠 Concurrency Architecture
Node.js uses a single-threaded event loop. Python uses process forks. 
Go uses the GMP Scheduler. 
When a new user connects to your Go server, the `http.Server` automatically spawns a brand new **Goroutine** exclusively for that user. This means your handlers can execute blocking code (like heavy database queries) without ever slowing down the other users on the server!
