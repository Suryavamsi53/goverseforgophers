# errgroup (Error Group)

While `sync.WaitGroup` is great for waiting on multiple Goroutines, it has two massive architectural flaws:
1. **Error Handling**: A Goroutine cannot easily return an `error` back to the main thread. You have to pass a complex slice or channel just to collect errors.
2. **Context Cancellation**: If you spawn 10 Goroutines, and the 1st one fails immediately, the other 9 will blindly continue running, wasting CPU cycles and database connections!

To solve this, the Go team provides an official extension package: `golang.org/x/sync/errgroup`.

## 1. Syntax and Error Collection

`errgroup.Group` acts as a drop-in replacement for `sync.WaitGroup`, but it automatically handles error propagation.

```go
import "golang.org/x/sync/errgroup"

func main() {
    var g errgroup.Group
    
    // We want to fetch 3 users concurrently
    userIDs := []int{1, 2, 3}
    
    for _, id := range userIDs {
        // We MUST capture the loop variable!
        userID := id 
        
        // g.Go() automatically calls Add(1) and Done() for us!
        g.Go(func() error {
            return fetchUser(userID) // Returns an error!
        })
    }
    
    // g.Wait() blocks until all Goroutines finish.
    // If ANY of the goroutines returned an error, g.Wait() returns the FIRST error it received!
    if err := g.Wait(); err != nil {
        fmt.Println("At least one fetch failed:", err)
    } else {
        fmt.Println("All users fetched successfully!")
    }
}
```

## 2. Automatic Context Cancellation

The true superpower of `errgroup` is `WithContext`. 

If you pass a parent Context into `errgroup.WithContext()`, it creates a new child context. 
If **any** Goroutine returns an error, the `errgroup` automatically calls `cancel()` on that child context! 

If the other 9 Goroutines are passing that context down to their database queries, those queries will be instantly aborted, saving massive amounts of system resources.

```go
func main() {
    // 1. Create the Group with a Context!
    g, ctx := errgroup.WithContext(context.Background())
    
    for i := 1; i <= 3; i++ {
        id := i
        g.Go(func() error {
            // 2. Pass the ctx down! 
            // If Worker 1 fails, this ctx will be cancelled, instantly aborting Worker 2 and 3!
            return fetchUserWithContext(ctx, id)
        })
    }
    
    if err := g.Wait(); err != nil {
        fmt.Println("Process aborted due to error:", err)
    }
}
```

## 3. Limiting Concurrency

In Go 1.20, `errgroup` introduced the `SetLimit()` method. 

If you have 10,000 URLs to scrape, you do not want to spawn 10,000 Goroutines instantly (you will get rate-limited by the API). You want a maximum of 10 concurrent workers.

```go
g.SetLimit(10) // Only 10 Goroutines will execute at a time!

for _, url := range urls {
    u := url
    // If 10 are running, the 11th call to g.Go() will BLOCK here until one finishes!
    g.Go(func() error {
        return download(u)
    })
}
```
`errgroup` is the ultimate, enterprise-grade replacement for `sync.WaitGroup`.
