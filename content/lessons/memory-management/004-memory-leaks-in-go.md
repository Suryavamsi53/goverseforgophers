# Memory Leaks in Go

Because Go is garbage collected, many developers believe it is impossible to create a Memory Leak. 
This is completely false. While Go prevents you from "forgetting to free" memory (like in C), you can easily create logical memory leaks where the Garbage Collector refuses to delete memory because your code is secretly holding a reference to it!

## 1. The Goroutine Leak (The Most Common)

A Goroutine takes ~2KB of memory. If a Goroutine is blocked and cannot exit, it lives forever.

```go
func LeakyFunction() {
    ch := make(chan int) // Unbuffered channel!

    go func() {
        // This Goroutine blocks forever waiting to send!
        ch <- 42
        fmt.Println("Done")
    }()

    // The main function returns without ever reading from 'ch'!
}
```

If you call this function 100,000 times, you will spawn 100,000 Goroutines. None of them can exit. The Garbage Collector sees that they are still "running" and refuses to clean them up. Your server crashes with an OOM.

**The Fix:** Always use `context.Context` to send a cancellation signal, or ensure channels are properly buffered and closed.

## 2. The Subslice Leak (The Hidden Danger)

This is the most notorious memory leak in Go.

Imagine you download a massive 100 Megabyte JSON file into a `[]byte`.
You only need the first 10 bytes (maybe an ID).

```go
func GetID() []byte {
    massiveData := download100MB() // Allocates 100MB on the Heap
    
    // We slice the first 10 bytes and return it!
    smallSlice := massiveData[0:10]
    return smallSlice
}
```

**The Bug:** In Go, slicing a slice (`[0:10]`) does NOT copy the data! It creates a new slice header that points to the exact same 100MB backing array in memory!
Because `smallSlice` is returned, it stays alive. Because `smallSlice` points to the 100MB array, the Garbage Collector **cannot delete the 100MB array**, even though you only care about 10 bytes!

**The Fix:** You must physically copy the data to a brand new array to sever the connection!

```go
func GetID() []byte {
    massiveData := download100MB()
    
    // Create a brand new, isolated 10-byte slice!
    safeSlice := make([]byte, 10)
    
    // Physically copy the bytes over
    copy(safeSlice, massiveData[0:10])
    
    // Return the safe copy! The 100MB array has no more pointers, the GC deletes it!
    return safeSlice
}
```
*(Note: As of Go 1.22, `slices.Clone()` does this automatically!)*

## 3. The `time.Ticker` Leak

If you use `time.Tick(1 * time.Second)`, it creates an infinite background timer. 

```go
func DoWork() {
    for range time.Tick(time.Second) {
        fmt.Println("Tick!")
        return // We exit early!
    }
}
```
If you return early, the underlying OS timer is never stopped. It ticks in the background forever, leaking memory.

**The Fix:** Always use `time.NewTicker()` and explicitly call `.Stop()` using `defer`!

```go
func DoWorkSafe() {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop() // Guarantees the OS timer is destroyed!

    for range ticker.C {
        // ...
    }
}
```
