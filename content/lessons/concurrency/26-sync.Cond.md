# sync.Cond (Condition Variables)

`sync.Cond` is the most misunderstood and rarely used primitive in the `sync` package. 

Imagine you have 10 worker Goroutines that are waiting for a specific condition to be met (e.g., "Wait until the queue has at least 5 items").

### The Bad Solution (Busy Waiting)
```go
// Terrible! Spikes CPU usage to 100% just to wait!
for len(queue) < 5 {
    time.Sleep(10 * time.Millisecond) 
}
```

### The OK Solution (Channels)
You could use a channel. But what if you want to wake up **all 10** Goroutines at the exact same time? If you send `ch <- true`, only ONE Goroutine wakes up. (You would have to `close(ch)` to wake them all up, but then you can never reuse the channel!).

This is where `sync.Cond` shines. It allows you to **broadcast** a wakeup signal to multiple sleeping Goroutines!

## 1. Syntax

A `sync.Cond` is always paired with a `sync.Mutex`.

```go
var mu sync.Mutex
// Create a Cond variable bound to the Mutex
cond := sync.NewCond(&mu)
queue := []int{}

// --- WORKER GOROUTINES ---
for i := 0; i < 3; i++ {
    go func(id int) {
        mu.Lock() // 1. Lock the Mutex first!
        
        // 2. We MUST use a 'for' loop, not an 'if' statement!
        for len(queue) == 0 {
            // 3. Wait instantly UNLOCKS the Mutex and puts the Goroutine to sleep.
            // When it wakes up, it automatically RE-LOCKS the Mutex!
            cond.Wait()
        }
        
        fmt.Printf("Worker %d got item: %d\n", id, queue[0])
        queue = queue[1:] // Consume item
        mu.Unlock()
    }(i)
}

// --- MAIN GOROUTINE (PRODUCER) ---
time.Sleep(1 * time.Second)
mu.Lock()
queue = append(queue, 42)
mu.Unlock()

// 4. Wake up EXACTLY ONE sleeping worker
cond.Signal() 

// OR: Wake up ALL sleeping workers simultaneously!
// cond.Broadcast()
```

## 2. The `Wait()` Loop Rule

Look closely at the worker code:
```go
for len(queue) == 0 {
    cond.Wait()
}
```
Why is this a `for` loop? Why not an `if` statement?

Because of **Spurious Wakeups**. In complex operating systems, a sleeping thread can occasionally wake up without a signal being sent! Furthermore, if `Broadcast()` wakes up 3 workers, but there is only 1 item in the queue, the first worker will consume it. When the second worker wakes up, the queue is empty again!

By using a `for` loop, the Goroutine re-checks the condition immediately after waking up. If the condition is no longer true, it goes safely right back to sleep.

## 3. Why it is rarely used

While `sync.Cond` is powerful, it is incredibly difficult to read and reason about. In 99% of Go codebases, developers prefer using Channels or Context cancellation to signal state changes, because they are more idiomatic and less prone to Deadlocks.
