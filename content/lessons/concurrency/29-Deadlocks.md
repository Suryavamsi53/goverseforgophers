# Deadlocks

A **Deadlock** occurs when a group of Goroutines are all waiting for each other to release a lock or send a message, forming a circular dependency. Because they are all waiting, none of them can make progress, and the application freezes permanently.

## 1. The Single Goroutine Deadlock

The simplest deadlock occurs when a single Goroutine blocks itself permanently.

```go
func main() {
    ch := make(chan int)
    
    // The main goroutine tries to receive from an empty channel.
    // It goes to sleep, waiting for another goroutine to send data.
    // But there are no other goroutines! It sleeps forever.
    <-ch 
}
```

The Go Runtime is incredibly smart. If it detects that **every single Goroutine** in the application is asleep, it knows the application can never wake up. It will instantly crash the app with: `fatal error: all goroutines are asleep - deadlock!`.

## 2. The Ghost Deadlock (Resource Leaks)

The Go Runtime can only detect a global deadlock. It cannot detect a partial deadlock.

```go
func main() {
    ch := make(chan int)
    
    // Worker goes to sleep waiting for data
    go func() {
        <-ch
        fmt.Println("Done")
    }()
    
    // Main thread does other things and exits
    time.Sleep(1 * time.Second)
    fmt.Println("Server running fine")
}
```

If you run this code, the Go Runtime will NOT crash. The `main` Goroutine is awake and doing work! But the background worker is deadlocked. It will sit in RAM forever. If you spawn 100 of these buggy background workers, you have a massive **Memory Leak**.

## 3. The Classic Mutex Deadlock (AB-BA)

In systems with multiple Mutexes, deadlocks occur if two Goroutines try to acquire the same locks in a different order.

```go
var lockA sync.Mutex
var lockB sync.Mutex

func worker1() {
    lockA.Lock()
    time.Sleep(1 * time.Millisecond) // Simulate work
    lockB.Lock()                     // Tries to get B, but Worker 2 has it!
    
    lockB.Unlock()
    lockA.Unlock()
}

func worker2() {
    lockB.Lock()
    time.Sleep(1 * time.Millisecond) // Simulate work
    lockA.Lock()                     // Tries to get A, but Worker 1 has it!
    
    lockA.Unlock()
    lockB.Unlock()
}
```

* Worker 1 grabs Lock A and waits for Lock B.
* Worker 2 grabs Lock B and waits for Lock A.

Neither can proceed. This is the **AB-BA Deadlock**.

### The Solution: Lock Ordering
To prevent this, you must establish a strict architectural rule for **Lock Ordering**. If your system requires multiple locks, every single Goroutine must acquire them in the exact same alphabetical or hierarchical order. (e.g., Always grab Lock A *before* Lock B, everywhere in the codebase).
