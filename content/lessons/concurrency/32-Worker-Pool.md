# The Worker Pool Pattern

The Worker Pool is the most famous and widely implemented concurrency pattern in Go. 

If you need to process 1,000,000 images, spawning 1,000,000 Goroutines is technically possible (they only cost 2KB of RAM each). However, if all 1,000,000 Goroutines try to open a file or write to a database simultaneously, you will exhaust your server's File Descriptors and TCP Ports, crashing the entire OS.

You must throttle concurrency. The Worker Pool pattern achieves this by spawning a fixed number of workers that continually pull jobs from a shared queue.

## 1. The Architecture

The architecture consists of:
1. A **Buffered Channel** to hold the jobs.
2. A fixed number of **Worker Goroutines** running a `for range` loop.
3. An **Orchestrator** (usually `main`) that dumps jobs into the channel and then closes it.

## 2. Implementation

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// The worker function. It ONLY accepts receive-only channels!
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // The range loop will automatically exit when the jobs channel is closed!
    for j := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, j)
        time.Sleep(100 * time.Millisecond) // Simulate heavy work
        results <- j * 2
    }
}

func main() {
    const numJobs = 100
    const numWorkers = 5
    
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)
    var wg sync.WaitGroup
    
    // 1. Boot up the Worker Pool
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }
    
    // 2. Dump all 100 jobs into the Buffered Channel
    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    
    // 3. CRITICAL: Close the jobs channel! 
    // This tells the workers to exit their range loops once the queue is empty.
    close(jobs)
    
    // 4. Wait for all workers to finish
    wg.Wait()
    close(results) // Safe to close results now
    
    fmt.Println("All jobs processed.")
}
```

## 3. The Power of the Pattern

With exactly 5 workers processing 100 jobs, the application will only ever maintain 5 active database connections, and consume a perfectly stable amount of CPU and RAM. 

If you deploy this to Kubernetes, you can trivially scale the `numWorkers` variable using an environment variable based on the physical size of the Pod.
