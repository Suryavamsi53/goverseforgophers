# Channel Directions

By default, when you create a channel using `ch := make(chan int)`, it is **bidirectional**. You can both send data into it (`ch <- 1`) and receive data from it (`<-ch`).

However, in Enterprise Go Architecture, passing bidirectional channels into functions is considered a security and design risk. If you pass a channel to a `Worker` function so it can receive jobs, you do not want that `Worker` function accidentally sending data *back* into the jobs queue!

To enforce architectural boundaries, Go allows you to restrict channel directions in function signatures.

## 1. Send-Only Channels (`chan<-`)

A send-only channel is denoted by the arrow pointing **into** the `chan` keyword. A function accepting this can only produce data.

```go
// This function is mathematically proven to ONLY send data. 
// If it tries to read `<-ch`, the Go compiler will instantly fail.
func Producer(ch chan<- string) {
    ch <- "Job 1"
    ch <- "Job 2"
}
```

## 2. Receive-Only Channels (`<-chan`)

A receive-only channel is denoted by the arrow pointing **out of** the `chan` keyword. A function accepting this can only consume data.

```go
// This function is mathematically proven to ONLY receive data.
// If it tries to send `ch <- "Fail"`, the Go compiler will instantly fail.
func Consumer(ch <-chan string) {
    job := <-ch
    fmt.Println("Processing:", job)
}
```

## 3. Implicit Conversion

You do not need to cast the channel. The Go compiler automatically handles the conversion from a Bidirectional channel to a Unidirectional channel when you pass it as an argument.

```go
func main() {
    // 1. Create a bidirectional channel
    jobs := make(chan string, 10)

    // 2. Pass it to Producer (automatically converted to chan<-)
    go Producer(jobs)

    // 3. Pass it to Consumer (automatically converted to <-chan)
    go Consumer(jobs)
    
    // ...
}
```

## 4. The Architectural Benefit

Using unidirectional channels guarantees the **Principle of Least Privilege**. It makes your functions self-documenting. 

If a senior engineer sees `func AuditLog(ch chan<- LogEntry)`, they immediately know that `AuditLog` is responsible for generating logs, not processing them, without having to read a single line of the function body.
