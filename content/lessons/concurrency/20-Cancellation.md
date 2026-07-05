# Context Cancellation

In the previous lesson, we learned *why* we need `Context`. Now we will learn how to trigger and intercept the cancellation signals.

## 1. WithCancel

The most common way to make a context cancellable is using `context.WithCancel`. This returns a child context and a `cancel` function.

```go
// Create a cancellable child context from the background
ctx, cancel := context.WithCancel(context.Background())

// ALWAYS defer the cancel function!
// If you don't, the child context will leak memory in the parent tree.
defer cancel()

// Spawn a worker
go worker(ctx)

// Manually trigger the cancellation signal
cancel()
```

## 2. Intercepting the Signal (Done Channel)

How does the `worker` Goroutine know that `cancel()` was called? 

The Context object provides a `.Done()` method, which returns a read-only channel `<-chan struct{}`. 
When `cancel()` is called, the Go Runtime **closes** this channel. 

As we learned in the Channel Closing lesson, receiving from a closed channel instantly returns without blocking. We intercept this using a `select` statement.

```go
func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            // The Done() channel was closed! The signal arrived!
            fmt.Println("Worker shutting down due to cancellation:", ctx.Err())
            return
            
        default:
            // Do normal work here
            fmt.Println("Worker is processing...")
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

## 3. WithTimeout (The Lifesaver)

In HTTP handlers, you rarely use `WithCancel` manually. Instead, you use `context.WithTimeout`.

If your microservice calls a 3rd party API, and you want to guarantee it never hangs for more than 2 seconds, you wrap it in a timeout context.

```go
// This context will automatically call cancel() after 2 seconds
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel() // Still defer to release resources if it finishes early!

// Pass the context into the HTTP Request
req, _ := http.NewRequestWithContext(ctx, "GET", "https://slow-api.com", nil)

// If the API takes 3 seconds, the HTTP Client will detect the Context cancellation
// at 2 seconds, instantly drop the TCP connection, and return an error!
res, err := http.DefaultClient.Do(req)
if err != nil {
    fmt.Println("Request failed or timed out:", err)
}
```

## 4. How the Standard Library Uses Context

When you use `context.WithTimeout`, how does `http.DefaultClient.Do` actually know to cancel the TCP connection?

Deep inside the Go standard library, in the `net` package, the OS socket operations are wrapped in a `select` statement! 

```go
// Inside the Go Standard Library (Conceptual)
select {
case result := <-readSyscall():
    return result
case <-ctx.Done():
    // The timeout fired! Force close the TCP socket!
    closeTCPConnection()
    return context.DeadlineExceeded
}
```
This is why you must pass the Context down through every single function in your architecture. If you fail to pass it to the final database query or HTTP request, the standard library cannot execute the `select` statement, and the cancellation signal will be completely ignored!
