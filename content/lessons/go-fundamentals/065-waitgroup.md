# WaitGroup

In the previous lesson, our `main` function exited before our background goroutines could finish their work. 

To solve this, we use a `sync.WaitGroup`. It acts as a concurrent counter that allows the main thread to wait until a specific number of workers have finished their tasks.

## 1. Syntax: Add, Done, Wait

A WaitGroup has exactly three methods:
1. `Add(int)`: Increases the counter. (Called *before* starting the goroutine).
2. `Done()`: Decreases the counter by 1. (Called *inside* the goroutine).
3. `Wait()`: Blocks the current thread until the counter hits 0.

```go
import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    // 3. Ensure Done() is called even if the worker panics
    defer wg.Done() 
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Second) // Simulate work
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        // 1. Increment the counter BEFORE spawning the goroutine
        wg.Add(1) 
        go worker(i, &wg)
    }

    // 2. Block the main thread until the counter drops to 0
    fmt.Println("Main thread waiting...")
    wg.Wait() 
    
    fmt.Println("All workers finished!")
}
```

## 2. The Pass-by-Value Trap

Notice how we passed `&wg` (a pointer) into the worker function?

What happens if we pass the WaitGroup by value? `func worker(id int, wg sync.WaitGroup)`

**Deadlock!**
If you pass a WaitGroup by value, Go creates a complete *copy* of the WaitGroup struct. The worker function will call `Done()` on the copy, meaning the original WaitGroup in `main()` will never decrement. The `Wait()` function will wait forever, and the Go runtime will crash the program with a `fatal error: all goroutines are asleep - deadlock!`.

## 3. Alternative: Closures

To avoid the pointer pass-by-value trap entirely, it is highly idiomatic in Go to use closures so the goroutine can capture the original WaitGroup by reference automatically.

```go
func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 3; i++ {
        wg.Add(1)
        
        go func(id int) {
            defer wg.Done() // Captures 'wg' safely!
            fmt.Println("Working:", id)
        }(i)
    }

    wg.Wait()
}
```
