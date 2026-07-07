# Concurrency Optimizations (Lock Contention)

## 1. Learning Objectives
* **What you'll learn**: How to identify lock contention, when to use `sync.RWMutex`, and how to use Lock-Free data structures via `sync/atomic`.
* **Why it matters**: Goroutines are incredibly lightweight, but if 10,000 Goroutines all try to acquire the same `sync.Mutex` at the same time, your concurrent program effectively becomes single-threaded. This is called Lock Contention.
* **Where it's used**: High-concurrency state management, custom in-memory caches, and concurrent counters.

---

## 2. The Cost of Mutexes
A `sync.Mutex` ensures that only one Goroutine can access a block of memory at a time. While necessary for data integrity, acquiring a lock requires OS-level coordination (futexes on Linux).

### The Bottleneck: sync.Mutex
```go
var counter int
var mu sync.Mutex

// If 10,000 Goroutines call this, 9,999 of them are put to sleep,
// waiting in a queue. This causes massive context-switching overhead.
func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}
```

---

## 3. Optimization 1: sync.RWMutex
If your workload is "Read-Heavy" (e.g., an in-memory cache where 99% of requests are reading data, and 1% are writing data), a standard `sync.Mutex` is too restrictive. It blocks readers from reading while another reader is reading!

**Solution**: Use `sync.RWMutex`.
* `RLock()`: Allows infinite Goroutines to read at the same time.
* `Lock()`: Blocks everything (readers and writers) for exclusive write access.

```go
var cache map[string]string
var rwMu sync.RWMutex

func getCache(key string) string {
    rwMu.RLock()         // 10,000 Goroutines can do this simultaneously!
    defer rwMu.RUnlock()
    return cache[key]
}

func setCache(key, value string) {
    rwMu.Lock()          // Blocks everyone else until finished
    defer rwMu.Unlock()
    cache[key] = value
}
```

---

## 4. Optimization 2: Lock-Free Atomic Operations
For simple counters or boolean flags, a Mutex is extreme overkill. Modern CPUs have special hardware instructions that allow you to modify a variable concurrently without needing an OS-level lock.

In Go, this is exposed via the `sync/atomic` package.

### The Lock-Free Way
```go
import "sync/atomic"

var counter int64

// No Mutex! The CPU hardware ensures this is thread-safe.
// This is easily 10x to 50x faster than using sync.Mutex.
func incrementAtomic() {
    atomic.AddInt64(&counter, 1)
}

func getCounter() int64 {
    return atomic.LoadInt64(&counter)
}
```

---

## 5. Quiz

1. **MCQ**: When is it inappropriate to use `sync/atomic` instead of `sync.Mutex`?
   * (A) When incrementing a hit counter.
   * (B) When you need to update two different variables together in a single transaction. *(Answer: B, atomic operations only protect a single memory address at a time).*
   * (C) When tracking a boolean flag.

2. **System Design Follow-up**: If `sync.RWMutex` allows infinite concurrent readers, why shouldn't you just use it everywhere instead of `sync.Mutex`?
   * *(Because `sync.RWMutex` is a more complex struct. The act of acquiring an `RLock` involves more internal bookkeeping and atomic operations than a standard `sync.Mutex`. If your workload is Write-Heavy, an `RWMutex` is actually slower than a normal `Mutex`!)*
