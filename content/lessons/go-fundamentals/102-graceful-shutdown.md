# Graceful Shutdowns

When you deploy a new version of your Go application, Kubernetes (or Docker) sends a signal to the old application to terminate. 

If your application instantly crashes upon receiving that signal, any users who were in the middle of downloading a file or processing a credit card payment will be abruptly disconnected. 

A production server must implement a **Graceful Shutdown**.

## 1. Catching OS Signals

Operating Systems communicate via Signals. 
* `SIGINT`: Sent when you press `Ctrl+C` in the terminal.
* `SIGTERM`: Sent by Kubernetes/Docker asking the app to shut down.
* `SIGKILL`: Sent by the OS to forcibly murder the app (cannot be caught).

Go's `os/signal` package allows us to intercept these signals so we can control the shutdown sequence.

```go
import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    // 1. Create a channel to receive OS signals
    quit := make(chan os.Signal, 1)
    
    // 2. Route SIGINT and SIGTERM to our channel
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    fmt.Println("Server running. Press Ctrl+C to stop.")
    
    // 3. Block the main thread until a signal arrives
    <-quit 
    
    fmt.Println("\nShutdown signal received! Cleaning up...")
    // ... Close databases, flush logs ...
}
```

## 2. Shutting Down the HTTP Server

When a shutdown signal arrives, we need to tell the `http.Server` to stop accepting *new* connections, but wait for all *active* connections to finish their work before exiting.

The `server.Shutdown(ctx)` method does exactly this!

```go
func main() {
    server := &http.Server{Addr: ":8080"}

    // 1. Start the server in a Goroutine so it doesn't block main()
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            panic(err)
        }
    }()

    // 2. Wait for the termination signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    fmt.Println("Shutting down server...")

    // 3. Create a 10-second timeout context
    // If active requests take longer than 10 seconds to finish, we forcefully kill them.
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 4. Gracefully drain connections
    if err := server.Shutdown(ctx); err != nil {
        fmt.Println("Server forced to shutdown:", err)
    }

    fmt.Println("Server gracefully stopped.")
}
```

### Architecture Insight
By combining OS signal trapping, a 10-second timeout Context, and `server.Shutdown()`, you ensure zero-downtime deployments. Kubernetes routes new traffic to your new container, while the old container safely drains its remaining users before quietly going to sleep.
