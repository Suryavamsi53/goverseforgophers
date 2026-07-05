# WaitGroup

To solve the "Main Exit Trap", the Go standard library provides the `sync` package. The most fundamental synchronization primitive is the `sync.WaitGroup`.

A WaitGroup acts as a concurrent counter that allows a Goroutine (usually `main`) to wait until a specific number of background Goroutines have completed their work.

## 1. The Three Methods

A `sync.WaitGroup` has three methods:
1. `Add(int)`: Increments the counter. (e.g., "I am starting 1 new task").
2. `Done()`: Decrements the counter by 1. (e.g., "I finished my task").
3. `Wait()`: Blocks the current Goroutine, completely putting it to sleep until the counter hits exactly `0`.

## 2. Basic Implementation

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    // 3. Guarantee Done() is called when the function exits
    defer wg.Done() 
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(100 * time.Millisecond)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        // 1. Increment the counter BEFORE spawning the goroutine!
        wg.Add(1)
        go worker(i, &wg)
    }

    // 2. Block the main thread until the counter hits 0
    wg.Wait()
    fmt.Println("All workers completed. Safe to exit.")
}
```

## 3. The Fatal Pointer Trap (Pass-By-Value)

Look closely at the `worker` function signature: `func worker(id int, wg *sync.WaitGroup)`.
We passed the WaitGroup as a **Pointer** `&wg`.

If you forget the `*` and pass it by value (`wg sync.WaitGroup`), Go will create a completely new **copy** of the WaitGroup inside the `worker` function. 
When the worker calls `wg.Done()`, it decrements the *copy*, not the original! 
The Main Goroutine will sit at `wg.Wait()` forever, waiting for the original counter to hit 0, resulting in a **Fatal Deadlock** crash.

## 4. The Loop Closure Bug

A very common pattern is wrapping the Goroutine in an anonymous closure instead of a named function.

```go
func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        // Anonymous closure
        go func(workerID int) {
            defer wg.Done() // We don't need a pointer here because closures capture state!
            fmt.Printf("Worker %d done\n", workerID)
        }(i) // Pass 'i' as an argument!
    }

    wg.Wait()
}
```
**Warning**: If you do not pass `i` as an argument into the closure, and instead just reference `i` directly inside the closure, you will hit the famous "Loop Variable Capture Bug". All 3 Goroutines will likely print "Worker 3 done", because the loop will have finished executing `i++` before the Goroutines even start!
