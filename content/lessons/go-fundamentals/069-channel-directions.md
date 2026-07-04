# Channel Directions

By default, when you create a channel, it is **bidirectional**. You can both send data into it and receive data from it.

However, when passing channels as arguments to functions, you can restrict their direction to increase type safety and enforce strict architectural boundaries.

## 1. Syntax

You define a directional channel by placing the arrow `<-` relative to the `chan` keyword.

* `chan T`: Bidirectional (Send and Receive)
* `chan<- T`: **Send Only** (You can write to it, but you cannot read from it)
* `<-chan T`: **Receive Only** (You can read from it, but you cannot write to it)

## 2. Enforcing Boundaries

Imagine a system where one function produces data (Ping) and another function consumes it (Pong). 

```go
// The Producer function is strictly allowed to SEND data.
// If it tries to read from 'out', the compiler will throw an error!
func ping(out chan<- string, msg string) {
    out <- msg
}

// The Consumer function is strictly allowed to RECEIVE data.
// If it tries to send data to 'in', the compiler will throw an error!
func pong(in <-chan string) {
    msg := <-in
    fmt.Println(msg)
}

func main() {
    // We create a standard Bidirectional channel
    ch := make(chan string)
    
    // When we pass it, Go automatically downgrades the permissions 
    // based on the function signatures!
    go ping(ch, "Hello")
    pong(ch)
}
```

### 🧠 Architectural Insight
Why restrict channels? In massive codebases, if a worker function accidentally reads from a channel it was only supposed to write to, it can steal data intended for another worker, leading to impossible-to-reproduce race conditions. By locking down the direction in the function signature, the Go compiler mathematically guarantees that data only flows in one direction, preventing these bugs entirely.
