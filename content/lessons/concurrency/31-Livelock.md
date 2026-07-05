# Livelock

We have covered Deadlocks (where Goroutines permanently sleep) and Starvation (where Goroutines are denied access). The final, and arguably most frustrating concurrency bug, is the **Livelock**.

## 1. What is a Livelock?

A Livelock occurs when two or more Goroutines are actively executing code and changing their state in response to each other, but they are doing so in a way that prevents any actual progress from being made.

Unlike a Deadlock, the CPU is not idle. A Livelock will consume **100% of your CPU**, spinning wildly while accomplishing absolutely nothing.

### The Real-World Analogy
Imagine two people walking toward each other in a narrow hallway. 
1. Person A steps to their right to let Person B pass.
2. At the exact same time, Person B steps to their left to let Person A pass.
3. They are still blocking each other.
4. Person A steps to their left.
5. Person B steps to their right.
6. They are still blocking each other.

They will repeat this dance forever. They are not "dead" (sleeping), they are very much alive and moving, but they are trapped in a loop.

## 2. Livelock in Code

Livelocks usually happen when developers try to be "smart" and implement their own lock recovery mechanisms using polling or timeouts, instead of relying on standard `Mutex` or `Channel` synchronization.

```go
func workerA(lock1, lock2 *int32) {
    for {
        // Try to grab Lock 1
        if atomic.CompareAndSwapInt32(lock1, 0, 1) {
            
            // Try to grab Lock 2
            if atomic.CompareAndSwapInt32(lock2, 0, 1) {
                fmt.Println("Worker A finished!")
                return // Success
            }
            
            // Failed to get Lock 2! We must release Lock 1 and try again.
            atomic.StoreInt32(lock1, 0)
        }
    }
}
```

If `workerA` and `workerB` execute simultaneously:
* A grabs `lock1`. B grabs `lock2`.
* A tries to grab `lock2` (fails). B tries to grab `lock1` (fails).
* A releases `lock1`. B releases `lock2`.
* They both instantly restart the loop.

This will run billions of times per second, melting your server's CPU, but the `fmt.Println` will never execute.

## 3. The Solution: Random Jitter

If you absolutely must implement a polling retry system (e.g., retrying an HTTP request that is failing due to rate limits), you must implement **Random Jitter**.

Going back to the hallway analogy: if Person A waits 1 second before stepping, and Person B waits a random time between 1 and 3 seconds, they will naturally desynchronize and successfully pass each other.

```go
// Adding Jitter to prevent Livelock sync loops
sleepTime := time.Duration(rand.Intn(100)) * time.Millisecond
time.Sleep(sleepTime)
```
