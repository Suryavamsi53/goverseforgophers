# Mutex (Mutual Exclusion)

Channels are great for passing data between goroutines, but what if you have a single piece of state (like a global cache or a user's bank balance) that 1,000 goroutines need to read and update at the same time?

Using channels for state management can be overly complex. For shared state, Go provides the traditional `sync.Mutex`.

## 1. The Race Condition

If multiple goroutines access the same variable simultaneously, and at least one of them is writing to it, you have a **Data Race**. 

```go
// ❌ DANGEROUS: Race Condition
var counter int

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter++ // BUG: 1000 goroutines modifying 'counter' simultaneously!
        }()
    }
    wg.Wait()
    fmt.Println(counter) // Might print 982, 991, or 1000. It's unpredictable!
}
```
If you do this with a `map`, Go will instantly crash with `fatal error: concurrent map read and map write`.

## 2. Locking with `sync.Mutex`

A Mutex (Mutual Exclusion lock) acts like a bathroom key. Only one goroutine can hold the key at a time. If another goroutine wants to modify the variable, it must wait in line for the key to be returned.

```go
var counter int
var mu sync.Mutex // 1. Declare the Mutex

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            mu.Lock()   // 2. Acquire the key
            counter++   // 3. Safely modify the data
            mu.Unlock() // 4. Return the key
        }()
    }
    wg.Wait()
    fmt.Println(counter) // Guaranteed to be exactly 1000
}
```

## 3. The `defer` Unlock Pattern

If your function has complex `if` statements or could potentially panic, it is very easy to forget to call `mu.Unlock()`. If you forget, every other goroutine will wait forever, causing a fatal Deadlock.

**Best Practice:** Always use `defer mu.Unlock()` immediately after locking.

```go
type SafeCache struct {
    mu   sync.Mutex
    data map[string]string
}

func (c *SafeCache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock() // Guaranteed to unlock when the function exits
    
    c.data[key] = value
}
```

## 4. Channels vs. Mutexes

The Go Wiki provides a golden rule for choosing between the two:
* **Use Mutexes** for protecting shared internal state (caches, configurations, simple counters).
* **Use Channels** for orchestrating control flow, passing ownership of data, or managing worker pools.
