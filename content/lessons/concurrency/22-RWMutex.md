# RWMutex (Readers-Writer Lock)

A standard `sync.Mutex` is brutally restrictive: it only allows exactly **one** Goroutine to access the locked code at a time. 

If you have a global configuration map that is updated once a day, but read 10,000 times a second by HTTP handlers, a standard Mutex will destroy your performance. 10,000 Goroutines will line up in single-file, waiting to read data that isn't even changing!

To solve this, we use the `sync.RWMutex` (Readers-Writer Lock).

## 1. How it works

An `RWMutex` distinguishes between two types of locks:
* **Read Lock (`RLock()`)**: Multiple Goroutines can hold a Read Lock simultaneously. If 10,000 Goroutines call `RLock()`, they all proceed instantly.
* **Write Lock (`Lock()`)**: Only **one** Goroutine can hold a Write Lock. If a Goroutine calls `Lock()`, it must wait for all active Readers to finish. Once the Writer acquires the lock, all new Readers are blocked until the Writer finishes.

## 2. Syntax

```go
type ConfigStore struct {
    mu     sync.RWMutex
    config map[string]string
}

// 1. READ-ONLY METHOD
func (c *ConfigStore) Get(key string) string {
    // 10,000 Goroutines can execute this simultaneously!
    c.mu.RLock() 
    defer c.mu.RUnlock()
    
    return c.config[key]
}

// 2. WRITE METHOD
func (c *ConfigStore) Set(key, value string) {
    // Blocks ALL other Readers and Writers until finished!
    c.mu.Lock() 
    defer c.mu.Unlock()
    
    c.config[key] = value
}
```

## 3. The Lock Upgrade Trap

A very common mistake is attempting to "upgrade" a Read Lock into a Write Lock.

```go
func (c *ConfigStore) UpdateIfEmpty(key, val string) {
    c.mu.RLock() // Acquire Read Lock
    
    if c.config[key] == "" {
        // DANGER: We are trying to acquire a Write Lock 
        // while we still hold the Read Lock!
        c.mu.Lock() 
        c.config[key] = val
        c.mu.Unlock()
    }
    
    c.mu.RUnlock()
}
```
**This will instantly Deadlock your application.** 
The `Lock()` call will wait forever for the active Reader (which is the current Goroutine itself!) to call `RUnlock()`. But `RUnlock()` can't be reached because `Lock()` is blocking!

To fix this, you must release the Read Lock *before* attempting to acquire the Write Lock.

## 4. When to avoid RWMutex

An `RWMutex` is heavier and slower than a standard `Mutex`. It contains more complex internal accounting to track the number of active readers.

If your workload is heavily write-biased (e.g., writing logs), or if the map is only read rarely, stick to a standard `sync.Mutex`. Only use `RWMutex` when Reads vastly outnumber Writes (like a 99% read / 1% write ratio).
