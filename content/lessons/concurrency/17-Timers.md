# Timers

In the previous lesson, we used `time.After(2 * time.Second)` to create a timeout inside a `select` block. Under the hood, `time.After` is a wrapper around a more powerful primitive: the `time.Timer`.

## 1. What is a Timer?

A Timer represents a single event in the future. You tell the Timer how long you want to wait, and it provides a channel that will receive the current time once that duration has expired.

```go
// Create a timer that will fire in 2 seconds
timer := time.NewTimer(2 * time.Second)

fmt.Println("Waiting...")
// Block until the timer fires data into its internal channel (timer.C)
<-timer.C
fmt.Println("Timer fired!")
```

## 2. Why use time.NewTimer over time.Sleep?

If `<-timer.C` blocks the Goroutine, why not just use `time.Sleep(2 * time.Second)`?

Because a Timer gives you **Control**. 
If you use `time.Sleep`, the Goroutine is put to sleep by the OS scheduler, and there is absolutely no way to wake it up early or cancel the sleep.

If you use a Timer, you can stop it!

### Stopping a Timer
Imagine you start a 5-second timer to timeout a database query, but the database replies in 1 second. You no longer need the timer!

```go
timer := time.NewTimer(5 * time.Second)

go func() {
    <-timer.C
    fmt.Println("Timer expired")
}()

// Stop the timer before it fires!
stop := timer.Stop()
if stop {
    fmt.Println("Timer was successfully stopped early.")
}
```

## 3. The `time.After` Leak

While `time.After` is incredibly convenient for `select` statements, it has a massive architectural flaw if used inside a `for` loop.

```go
// DANGER: Memory Leak
for {
    select {
    case data := <-ch:
        fmt.Println(data)
    case <-time.After(1 * time.Minute):
        fmt.Println("Timeout")
    }
}
```
Every single time this loop iterates, it creates a brand new 1-minute Timer in the Go Heap. 
If `ch` receives 10,000 messages per second, this loop will allocate 10,000 Timers in the Heap every second, and none of them will be Garbage Collected until their 1-minute duration expires! This will cause your server to run out of RAM instantly.

### The Correct Way (Reusing a Timer)
To fix this, you instantiate **one** Timer outside the loop, and `Reset()` it inside the loop.

```go
timer := time.NewTimer(1 * time.Minute)
defer timer.Stop()

for {
    // Crucial: Reset the timer before looping
    timer.Reset(1 * time.Minute) 
    
    select {
    case data := <-ch:
        fmt.Println(data)
    case <-timer.C:
        fmt.Println("Timeout")
    }
}
```
