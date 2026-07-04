# Goroutines

A Goroutine is a lightweight thread managed by the Go runtime. 

## 1. Syntax

To start a new goroutine, you simply put the `go` keyword in front of a function or method call.

```go
func printNumbers() {
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // 1. Spawns a new concurrent goroutine
    go printNumbers() 
    
    // 2. The main function continues executing immediately!
    fmt.Println("Main function finished")
}
```

### The Main Trap
If you run the code above, it will likely only print `"Main function finished"` and exit. The numbers will never print!

**Why?**
When the `main()` function finishes execution, the Go program terminates immediately. It **does not wait** for background goroutines to finish. Because `printNumbers` had a 100ms sleep delay, the program died before it ever woke up!

To fix this, we need synchronization tools (like `sync.WaitGroup` or Channels), which we will learn in the next lessons.

## 2. Anonymous Goroutines

You don't have to define a named function to use a goroutine. You can spawn them instantly using anonymous IIFEs (Immediately Invoked Function Expressions).

```go
func main() {
    message := "Processing Data"
    
    go func(msg string) {
        fmt.Println("Worker:", msg)
    }(message)
    
    // Hacky wait so we can see the output (never use time.Sleep for sync in production!)
    time.Sleep(time.Second) 
}
```

## 3. The Closure Trap (Go < 1.22)

Just like we learned in the Closures lesson, spawning anonymous goroutines inside a loop used to be the #1 cause of bugs in Go.

```go
for i := 0; i < 5; i++ {
    go func() {
        // Because the goroutine takes a millisecond to start, 
        // the loop already finished and 'i' is 5 for EVERY goroutine!
        fmt.Println(i) 
    }()
}
```

**The old fix:** Pass the variable into the goroutine as an argument.
```go
for i := 0; i < 5; i++ {
    go func(val int) {
        fmt.Println(val) // Captures a safe copy!
    }(i)
}
```
*(Note: As of Go 1.22, the compiler fixes this automatically by re-allocating `i` on every loop iteration, but passing the value explicitly is still considered excellent practice for readability).*
