# Ranging over Channels

If a worker Goroutine needs to process thousands of jobs from a channel, manually writing `val, ok := <-ch` inside a `for` loop is tedious and error-prone.

The most idiomatic way to consume data from a channel in Go is to use the `for range` loop.

## 1. Syntax

```go
func Consumer(jobs <-chan int) {
    // This loop automatically receives data.
    // It will BLOCK and sleep if the channel is empty.
    // It will AUTOMATICALLY EXIT the loop when the channel is closed.
    for job := range jobs {
        fmt.Printf("Processing Job: %d\n", job)
    }
    
    fmt.Println("All jobs processed, worker shutting down.")
}
```

## 2. The Deadlock Trap

The `range` loop is extremely powerful, but it introduces a massive risk if you do not strictly control the channel lifecycle.

If the Sender forgets to `close(ch)` when it is done, the Consumer's `range` loop will never exit. It will sit permanently blocked, waiting for more data. 
If this was the only running Goroutine, the Go Runtime will detect that all Goroutines are asleep forever and crash the application with a **Fatal Deadlock**.
If this was a background Goroutine in a web server, it will just sit in memory forever, causing a massive **Goroutine Memory Leak**.

## 3. The Worker Pool Application

Combining `for range` with Buffered Channels creates the most famous Go architecture: the **Worker Pool**.

```go
func main() {
    jobs := make(chan int, 100)
    var wg sync.WaitGroup

    // 1. Spawn 3 Worker Goroutines
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            // The 3 workers will compete to pull jobs off the channel using range!
            for j := range jobs {
                fmt.Printf("Worker %d processing job %d\n", workerID, j)
            }
        }(i)
    }

    // 2. Dump 10 jobs into the queue
    for j := 1; j <= 10; j++ {
        jobs <- j
    }

    // 3. CRITICAL: Close the channel! 
    // This signals the `range` loops in all 3 workers to exit when the buffer is empty.
    close(jobs)

    // 4. Wait for all workers to exit their range loops and call Done()
    wg.Wait()
}
```

Because `range` automatically handles the `comma-ok` idiom under the hood, all 3 workers safely consume the queue and gracefully exit exactly when the buffer hits zero.
