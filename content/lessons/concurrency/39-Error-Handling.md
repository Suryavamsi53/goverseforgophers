# Concurrency Error Handling

In a sequential Go program, error handling is simple: you return `(val, err)` and check `if err != nil`.

In a concurrent Go program, Goroutines do not return values to the caller! The `go` keyword executes the function asynchronously and instantly discards any return values. 

```go
// DANGER: The error returned by this function is completely lost!
go func() error {
    return errors.New("database connection failed")
}()
```

How do you safely capture errors from 100 background workers?

## 1. The Error Channel

The most idiomatic way to handle concurrent errors is to treat `error` as just another piece of data, and pipe it through a channel.

```go
func main() {
    errs := make(chan error, 3) // Buffer size matches the number of workers!
    var wg sync.WaitGroup

    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            err := fetchFromDatabase()
            if err != nil {
                // Pipe the error out!
                errs <- fmt.Errorf("worker %d failed: %w", id, err)
            }
        }(i)
    }

    // Close the error channel once all workers finish
    go func() {
        wg.Wait()
        close(errs)
    }()

    // Range over the errors (if any)
    for err := range errs {
        fmt.Println("Caught error:", err)
    }
}
```

### The Unbuffered Trap
Notice that we created a buffered channel `make(chan error, 3)`. 
If we used an unbuffered channel `make(chan error)`, and Worker 1 failed and sent an error, it would block forever until the main thread reached the `for err := range errs` loop! 
By buffering the channel to the exact number of workers, we guarantee the workers can instantly dump their errors and exit gracefully without deadlocking.

## 2. Using `errgroup`

As covered in Lesson 27, if you only care about the **first** error that occurs, building channels and WaitGroups manually is unnecessary boilerplate.

You should immediately reach for `golang.org/x/sync/errgroup`.

```go
var g errgroup.Group

g.Go(func() error {
    return fetchFromDatabase()
})

if err := g.Wait(); err != nil {
    // This will cleanly capture and return the exact error from the Goroutine!
    log.Fatal(err) 
}
```

## 3. Panic Recovery in Goroutines

If a panic occurs inside the `main` Goroutine, you can use `defer recover()` to catch it.
However, **panics do not cross Goroutine boundaries**.

If a background Goroutine panics, and it does not have its own `recover()`, the entire application will instantly crash, regardless of whether `main` has a recover block!

```go
func main() {
    // This defer will NEVER catch the panic below!
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Main caught a panic!") 
        }
    }()

    go func() {
        // FATAL CRASH: This will instantly kill the entire server!
        panic("Database corrupted!") 
    }()

    time.Sleep(1 * time.Second)
}
```

### The Safe Goroutine Wrapper
If you are building an Enterprise system, you must wrap all background Goroutines in a recovery block to prevent a single buggy worker from taking down the entire web server.

```go
func SafeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Recovered from panic in background worker: %v\n", r)
            }
        }()
        fn() // Execute the actual work
    }()
}
```
