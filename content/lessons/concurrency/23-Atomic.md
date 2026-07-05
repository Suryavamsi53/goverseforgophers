# Atomic Operations

While a `sync.Mutex` is safe, it requires interacting with the Go Scheduler. When a Goroutine hits a locked Mutex, the Scheduler puts that Goroutine to sleep, switches to another Goroutine, and wakes the original one up later. 

This **Context Switching** takes valuable time (hundreds of nanoseconds). 

If you just want to increment a simple integer counter (`count++`), putting Goroutines to sleep is massive overkill. We can achieve this much faster using Hardware-Level Lock-Free programming: the `sync/atomic` package.

## 1. What is an Atomic Operation?

An atomic operation bypasses the Go Scheduler entirely. It compiles down to a single, specialized CPU instruction (like `LOCK XADD` on x86 processors). 

Because it is a single hardware instruction, it is physically impossible for two CPU cores to interrupt each other or create a Data Race.

## 2. Basic Syntax

```go
import "sync/atomic"

type Metrics struct {
    // Must be a precise integer size (int32, int64, uint64)
    requestCount int64 
}

func (m *Metrics) RecordRequest() {
    // 10x faster than a Mutex! No Context Switching!
    atomic.AddInt64(&m.requestCount, 1)
}

func (m *Metrics) GetCount() int64 {
    // Safely reads the value without tearing
    return atomic.LoadInt64(&m.requestCount)
}
```

## 3. The `atomic.Value` (Go 1.19+)

Historically, `sync/atomic` only worked on integers and pointers. If you wanted to atomically swap an entire Struct or a Map, you had to use weird unsafe pointer math.

Go introduced `atomic.Value` and later generic `atomic.Pointer[T]` to allow lock-free swapping of complex data structures.

This is famously used for **Hot Configuration Reloading**.

```go
import "sync/atomic"

type Config struct {
    APIKey string
    Retries int
}

// 1. Create a generic atomic pointer
var currentConfig atomic.Pointer[Config]

func init() {
    // Set initial config
    currentConfig.Store(&Config{APIKey: "default", Retries: 3})
}

// 2. This runs in a background Goroutine every 5 minutes
func reloadConfigFromDatabase() {
    newCfg := fetchConfig() // Expensive operation
    
    // Atomically swap the pointer. All future reads instantly see the new config!
    currentConfig.Store(newCfg)
}

// 3. This is called by 10,000 HTTP handlers per second safely!
func HandleRequest() {
    cfg := currentConfig.Load()
    fmt.Println("Using API Key:", cfg.APIKey)
}
```

## 4. Compare-And-Swap (CAS)

The most advanced atomic function is `CompareAndSwap`. It acts as an optimistic lock. 
It says to the CPU: *"Update this variable to X, but ONLY if its current value is Y."*

This is the exact primitive used under the hood to build standard Mutexes and complex lock-free data structures!
