# Goroutines (Syntax and Execution)

Now that we understand the deep internals of the Scheduler and the G-P-M model, how do we actually write concurrent code?

## 1. The `go` Keyword

Spawning a Goroutine is as simple as prefixing a function call with the `go` keyword.

```go
package main

import (
    "fmt"
    "time"
)

func fetchUser(id int) {
    time.Sleep(100 * time.Millisecond) // Simulating DB query
    fmt.Printf("Fetched user %d\n", id)
}

func main() {
    go fetchUser(1) // Runs concurrently!
    go fetchUser(2) // Runs concurrently!
    
    fmt.Println("Main function finished")
}
```

## 2. The Main Exit Trap

If you run the code above, the output will look like this:
```text
Main function finished
```

Where are the "Fetched user" logs? They never printed. 
Why? Because the `main` function is itself a Goroutine (the "Main Goroutine"). 

**The Golden Rule:** When the Main Goroutine finishes executing and returns, the Go Runtime instantly forcefully terminates every single other Goroutine in the application, regardless of whether they were finished or not.

Because `main()` takes 0.001 milliseconds to run, it exited before the `fetchUser` goroutines could finish their 100ms sleep.

## 3. The Ugly Fix: `time.Sleep`

The worst possible way to solve this is to force the Main Goroutine to sleep, giving the background workers time to finish.

```go
func main() {
    go fetchUser(1)
    go fetchUser(2)
    
    time.Sleep(1 * time.Second) // DO NOT DO THIS!
    fmt.Println("Main function finished")
}
```

This is terrible architecture. If the database takes 1.1 seconds, the program crashes early. If the database takes 0.1 seconds, your program wastes 0.9 seconds of CPU time doing absolutely nothing.

We need a way to mathematically synchronize Goroutines. That is where `sync.WaitGroup` comes in.
