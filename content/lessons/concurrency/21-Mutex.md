# Mutex (Mutual Exclusion)

While Channels (CSP) are the idiomatic way to pass data *between* Goroutines, there are times when sharing memory is unavoidable. 

If you have a global Cache Map or a shared `counter` variable, and 1,000 Goroutines try to read and write to it simultaneously, the memory will become corrupted. Go's runtime is so strict about this that if two Goroutines write to a Go `map` concurrently, the application will instantly **Fatal Panic**.

To protect shared memory, we use a **Mutex**.

## 1. What is a Mutex?

A Mutex (Mutual Exclusion lock) is a synchronization primitive. When a Goroutine calls `Lock()`, it gains exclusive access to the code block. 
If a second Goroutine calls `Lock()` while the first one is still inside, the second Goroutine is instantly put to sleep by the Scheduler until the first one calls `Unlock()`.

## 2. Basic Syntax

```go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    // The defer guarantees Unlock is called even if a Panic occurs below!
    defer c.mu.Unlock() 
    
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    return c.count
}
```

## 3. The Defer Rule

Always use `defer mu.Unlock()` immediately after `mu.Lock()`. 

If you write a complex function with 5 different `return` statements, and you forget to manually call `Unlock()` before just one of those returns, the Mutex will remain permanently locked. Every other Goroutine in your application that tries to access that data will permanently freeze (Deadlock).

## 4. The Lock Copying Trap

A `sync.Mutex` contains internal state (a semaphore). It must **never be copied**.

```go
// DANGER: Passing Mutex by Value!
func updateCache(mu sync.Mutex, data map[string]string) {
    mu.Lock() // This locks the COPY, not the original!
    defer mu.Unlock()
    data["key"] = "value" // Fatal Crash: Concurrent Map Write!
}
```
If you pass a Mutex into a function, you must pass it as a pointer (`*sync.Mutex`). Better yet, embed the Mutex inside a struct (like `SafeCounter` above) and only pass pointers to the struct!
