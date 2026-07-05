# Tickers

A Timer fires exactly once. But what if you need to execute a task repeatedly on a strict interval? 

For example, checking a database queue for new jobs every 5 seconds, or pinging a WebSocket connection every 10 seconds to keep it alive.

For this, you use a `time.Ticker`.

## 1. Syntax

A Ticker provides a channel (`ticker.C`) that receives the current time repeatedly at a specified interval.

```go
// Creates a ticker that fires every 500 milliseconds
ticker := time.NewTicker(500 * time.Millisecond)

// A separate channel to shut it down
done := make(chan bool)

go func() {
    for {
        select {
        case <-done:
            fmt.Println("Ticker stopped!")
            return
        case t := <-ticker.C:
            fmt.Println("Tick at", t)
        }
    }
}()

// Let it run for 2 seconds (it will tick ~4 times)
time.Sleep(2 * time.Second)
ticker.Stop()
done <- true
```

## 2. The Stop() Leak

Look closely at the code above:
```go
ticker.Stop()
done <- true
```

Why did we need the `done <- true` signal? Doesn't `ticker.Stop()` kill the Goroutine?

**NO!** 

Calling `ticker.Stop()` stops the Go Runtime from sending new ticks into the `ticker.C` channel. However, **it does NOT close the channel**. 

If we didn't send the `done` signal, the Goroutine's `select` statement would just sit there forever, permanently blocked, waiting for a tick on a channel that will never tick again. This causes a massive Goroutine leak. 

Always design a secondary cancellation mechanism (like a `done` channel or a `Context`) when using Tickers inside background Goroutines.

## 3. The Dropped Tick (Backpressure)

What happens if the Ticker is set to 100ms, but the task inside the loop (like a database query) takes 500ms?

Will the Ticker queue up 5 ticks inside the channel?
**No.** 

The `ticker.C` channel only has a buffer capacity of `1`. 
If the Ticker tries to send a tick, but the Goroutine is busy and hasn't consumed the previous tick, the Go Runtime **drops the new tick entirely**. 

This is an incredible feature for system stability (built-in Backpressure). It guarantees your system will not be overwhelmed by a backlog of queued ticks if your database suddenly slows down. It simply skips the ticks until the worker is free again.
