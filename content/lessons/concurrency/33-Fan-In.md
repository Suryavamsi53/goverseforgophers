# Fan-In Pattern

When building distributed systems, you often have multiple separate data streams that need to be aggregated into a single pipeline for processing.

For example, a monitoring system might have 3 different Goroutines tailing 3 different log files. You want a central `Logger` Goroutine to write all of them to a central database.

The process of taking multiple input channels and multiplexing them into a single output channel is called **Fan-In**.

## 1. The Architecture

The Fan-In function takes an arbitrary number of input channels (`...<-chan int`). 
It creates a single output channel. 
It spawns a Goroutine for *every* input channel to forward the data to the output channel.

## 2. Implementation

```go
// merge takes multiple input channels and fans them into a single output channel
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // The forwarder function reads from a specific input channel and sends to out
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            out <- n
        }
    }

    // Spawn a forwarder for each input channel
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Spawn a background monitor to close the output channel 
    // once all forwarders are finished
    go func() {
        wg.Wait()
        close(out)
    }()

    return out
}
```

## 3. Why the Background Monitor?

Look closely at this block:
```go
go func() {
    wg.Wait()
    close(out)
}()
```
Why did we spawn a *new* Goroutine just to call `wg.Wait()`? Why didn't we just call `wg.Wait()` directly inside the `merge` function before returning?

If we called `wg.Wait()` directly, the `merge` function would **block**. It would wait for all data to finish processing before it ever returned the `out` channel to the caller! 
The caller wouldn't be able to start reading from `out`, meaning the forwarders would instantly deadlock trying to write to an unbuffered channel that no one is reading from!

By wrapping `wg.Wait()` in a background Goroutine, `merge` instantly returns the `out` channel. The caller begins reading, the forwarders begin writing, and when the streams naturally dry up, the background monitor quietly closes the `out` channel, gracefully terminating the caller's `range` loop.
