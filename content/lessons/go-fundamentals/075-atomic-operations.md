# Atomic Operations

While a `sync.Mutex` is safe, it is heavy. When a goroutine hits a locked Mutex, the Go Scheduler puts that goroutine to sleep, switches to another goroutine, and wakes the original one up later. This context switching takes valuable CPU cycles.

If all you need to do is increment a simple integer counter (like tracking total HTTP requests), a Mutex is massive overkill.

## 1. The `sync/atomic` Package

The `atomic` package bypasses the Go Scheduler entirely and leverages hardware-level CPU instructions (like **Compare-And-Swap**). 

Atomic operations are mathematically guaranteed to execute in a single, uninterruptible CPU clock cycle. No other thread can see the variable in a half-written state.

```go
import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var counter int64 // The variable to mutate
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // Atomic Addition: Blazing fast, no Mutex required!
            atomic.AddInt64(&counter, 1) 
        }()
    }

    wg.Wait()
    
    // Atomic Load: Safely read the value
    final := atomic.LoadInt64(&counter)
    fmt.Println("Total:", final)
}
```

## 2. Common Atomic Functions

The package provides primitives for `int32`, `int64`, `uint32`, `uint64`, and `Pointer` types.

* `AddT(addr *T, delta T)`: Safely adds a number.
* `LoadT(addr *T) T`: Safely reads a value.
* `StoreT(addr *T, val T)`: Safely overwrites a value.
* `SwapT(addr *T, new T) old T`: Overwrites a value and returns the old one.
* `CompareAndSwapT(addr *T, old, new T) bool`: The foundation of lock-free programming. It only updates the value if it currently equals `old`.

## 3. The `atomic.Value` Type

What if you need to atomically load and store a complex struct (like a massive application configuration object) without a Mutex?

You can use `atomic.Value`. It acts as a thread-safe container for any type (using `any` under the hood).

```go
var config atomic.Value

// Store a complex map atomically
config.Store(map[string]string{"db": "localhost:5432"})

// Load the map safely from 10,000 concurrent goroutines
// (Requires a Type Assertion to extract the data)
data := config.Load().(map[string]string)
```

**Performance Insight**: Atomic operations are an order of magnitude faster than Mutexes, but they should only be used for simple counters or lock-free data structures. For complex state logic, stick to `sync.Mutex`.
