# Semaphores

We've covered the Worker Pool as a way to limit concurrency (e.g., only 5 workers running at a time). 
However, setting up a Worker Pool requires boilerplate: spawning Goroutines, managing a WaitGroup, and orchestrating a job channel.

For simpler use cases, you can achieve concurrency limiting using a **Semaphore**.

## 1. What is a Semaphore?

A Semaphore is a synchronization primitive invented by Edsger Dijkstra. Think of it as a bouncer at a nightclub. The club has a strict capacity of 10 people. 
* If 10 people are inside, the bouncer blocks the door. 
* When 1 person leaves, the bouncer allows exactly 1 person from the line to enter.

## 2. Implementing a Semaphore in Go

Go does not have a built-in `sync.Semaphore` struct. Why? Because you can build a perfect Semaphore using a **Buffered Channel**!

By creating a buffered channel with a capacity of `N`, you create a Semaphore with a capacity of `N`.

```go
func main() {
    // 1. Create a Semaphore with a strict capacity of 3
    sem := make(chan struct{}, 3)
    var wg sync.WaitGroup

    // We have 100 jobs to run
    for i := 1; i <= 100; i++ {
        wg.Add(1)
        
        go func(jobID int) {
            defer wg.Done()

            // 2. ACQUIRE LOCK (The Bouncer)
            // We push an empty struct into the channel.
            // If 3 structs are already inside, this send will BLOCK!
            sem <- struct{}{}
            
            // 3. DO HEAVY WORK
            // Only 3 Goroutines will ever be executing this section simultaneously.
            fmt.Printf("Processing job %d\n", jobID)
            time.Sleep(1 * time.Second)
            
            // 4. RELEASE LOCK
            // We pull our struct out of the channel, freeing up 1 slot for the next Goroutine.
            <-sem 
            
        }(i)
    }

    wg.Wait()
}
```

## 3. The `x/sync/semaphore` Package

While the Buffered Channel trick is idiomatic and beautiful, it has one limitation: you can only acquire a lock with a weight of `1`.

What if you have a massive video processing job that requires 3 "slots" of CPU power, and a tiny audio job that requires 1 "slot"?

For weighted concurrency, Google provides the `golang.org/x/sync/semaphore` package.

```go
import "golang.org/x/sync/semaphore"

// Create a semaphore with a total weight of 10
var sem = semaphore.NewWeighted(10)

func main() {
    ctx := context.Background()
    
    // Acquire 3 slots of capacity. 
    // Blocks if the remaining capacity is < 3.
    if err := sem.Acquire(ctx, 3); err != nil {
        log.Fatal(err)
    }
    
    go func() {
        defer sem.Release(3)
        // ... process 4K Video ...
    }()
}
```
If you pass a timeout Context to `sem.Acquire`, it will automatically abort the acquisition if it waits in line too long!
