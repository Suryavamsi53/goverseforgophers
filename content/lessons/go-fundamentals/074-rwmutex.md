# RWMutex (Readers-Writer Lock)

A standard `sync.Mutex` is highly restrictive: it only allows **one** goroutine to access the data at a time, regardless of whether that goroutine is writing new data or just reading existing data.

If you have a global configuration struct that is read 10,000 times a second, but only updated once a day, a standard Mutex will bottleneck your entire application because the 10,000 readers have to wait in a single-file line.

## 1. Enter the `sync.RWMutex`

A Readers-Writer Mutex solves this by splitting the lock into two types:
* `Lock()` / `Unlock()`: The **Writer** Lock. Exclusive. No one else can read or write while this is held.
* `RLock()` / `RUnlock()`: The **Reader** Lock. Shared. Infinite goroutines can hold this lock simultaneously, *as long as no one holds the Writer lock.*

```mermaid
gantt
    title RWMutex Execution Timeline
    dateFormat  s
    axisFormat %S
    
    section RLock (Readers)
    Goroutine A (Reads) :a1, 0, 3s
    Goroutine B (Reads) :a2, 1, 3s
    Goroutine C (Reads) :a3, 1, 4s
    
    section Lock (Writers)
    Goroutine D (Writes) :crit, w1, 5s, 2s
    
    section RLock (Readers)
    Goroutine E (Reads) :a4, 7s, 2s
```
*Notice how Goroutines A, B, and C all execute concurrently! But Goroutine D (the Writer) must wait for them all to finish, and block Goroutine E until it is done.*

## 2. Implementation

```go
type ConfigStore struct {
    mu     sync.RWMutex
    config map[string]string
}

// Write Operation (Exclusive)
func (s *ConfigStore) Set(key, val string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.config[key] = val
}

// Read Operation (Shared / Highly Concurrent)
func (s *ConfigStore) Get(key string) string {
    s.mu.RLock() // Use RLock for reading!
    defer s.mu.RUnlock()
    return s.config[key]
}
```

## 3. The Upgrade Trap

A common mistake is attempting to "upgrade" a reader lock to a writer lock. 

```go
// ❌ FATAL DEADLOCK
func (s *ConfigStore) UpdateIfEmpty(key, val string) {
    s.mu.RLock()
    if s.config[key] == "" {
        // You cannot call Lock() while holding an RLock()! 
        // This will instantly deadlock the application.
        s.mu.Lock() 
        s.config[key] = val
        s.mu.Unlock()
    }
    s.mu.RUnlock()
}
```
If you need to read and then write, you must completely release the `RLock` before acquiring the `Lock`.
