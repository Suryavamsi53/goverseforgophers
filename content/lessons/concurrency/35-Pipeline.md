# The Pipeline Pattern

The **Pipeline** is an advanced architectural pattern used to process streams of data through a series of sequential stages. It is famously used in ETL (Extract, Transform, Load) jobs and video processing engines.

## 1. What is a Pipeline?

Imagine a factory assembly line.
1. Station A (Extract) pulls unpainted car frames from a warehouse.
2. Station B (Transform) paints the cars.
3. Station C (Load) installs the engines and drives them to the dealership.

Instead of one giant function doing all three tasks, a Pipeline splits the tasks into three independent Goroutines, connected by Channels. 

**The incredible benefit**: While Station C is installing the engine in Car 1, Station B is already painting Car 2, and Station A is already fetching Car 3! The system is constantly flowing.

## 2. Implementation

Each stage of the pipeline is a function that takes an `<-chan` (input) and returns a `<-chan` (output).

### Stage 1: Generator
```go
// Generates a stream of numbers
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out) // Crucial: Closes when generator finishes
    }()
    return out
}
```

### Stage 2: Transformer
```go
// Squares the numbers
func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out) // Crucial: Closes when input stream dries up
    }()
    return out
}
```

### Stage 3: The Consumer (main)
```go
func main() {
    // Wire the pipeline together!
    // gen() -> sq() -> print
    
    in := gen(2, 3, 4)
    out := sq(in)
    
    for result := range out {
        fmt.Println(result) // Prints 4, 9, 16
    }
}
```

## 3. Graceful Shutdown (The Context Injection)

The code above works perfectly, but it hides a massive memory leak risk.

What if the `main` loop decides to exit early (e.g., it only wanted the first 2 results)?
If `main` exits early, the `in` and `out` channels are never closed. The `gen` and `sq` Goroutines will sit permanently blocked on `out <- n`, resulting in a global deadlock or a Ghost Goroutine leak.

To build a production-grade Pipeline, **every single stage must accept a `context.Context`**.

```go
func sq(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
                // Success
            case <-ctx.Done():
                // The main orchestrator cancelled the pipeline early!
                // Instantly exit the Goroutine to prevent memory leaks!
                return
            }
        }
    }()
    return out
}
```
