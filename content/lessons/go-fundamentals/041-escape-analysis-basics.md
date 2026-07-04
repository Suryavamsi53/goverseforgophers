# Escape Analysis Basics

We briefly touched on Escape Analysis in the Memory Layout lesson. Here, we will dive deeply into how the compiler makes these decisions, and how you can actually audit your own code to find hidden performance bottlenecks.

## 1. What is Escape Analysis?

Escape Analysis is a phase during the Go compilation process. The compiler reads your code and determines if a variable can safely be allocated on the ultra-fast Stack, or if it must "escape" to the Garbage-Collected Heap.

### The Rule of Thumb
If the compiler can definitively prove that a variable will not be referenced after the function returns, it stays on the **Stack**. If it cannot prove this, the variable goes to the **Heap**.

## 2. Common Triggers for Heap Escapes

Here are the three most common reasons a variable escapes to the Heap:

1. **Returning a Pointer**: If you create a struct and return a pointer to it (`return &user`), the struct must survive the function's death. It escapes to the Heap.
2. **Assigning to an Interface**: When you assign a concrete value to an interface (like `fmt.Println(myVar)`), the interface requires a dynamic heap allocation under the hood. 
3. **Closures**: If an anonymous function captures an outside variable and is executed later, the variable escapes.

## 3. Auditing Your Code (`gcflags`)

You don't have to guess where your memory is going. The Go compiler will literally tell you if you ask it.

By passing the `-gcflags="-m"` flag to the `go build` or `go run` command, the compiler prints out its escape analysis decisions.

```go
package main
import "fmt"

func calculate() *int {
    x := 42
    return &x
}

func main() {
    val := calculate()
    fmt.Println(*val)
}
```

Run this in your terminal:
```bash
$ go run -gcflags="-m" main.go
```

**The Output:**
```text
./main.go:5:2: moved to heap: x
./main.go:11:13: ... argument does not escape
./main.go:11:14: *val escapes to heap
```
*(Notice how the compiler caught that `x` was moved to the heap because we returned a pointer to it!)*

## 4. The Performance Paradox

Many developers think returning pointers is *always* faster because it avoids copying structs. 

**This is a myth.**

Because returning a pointer forces the struct onto the Heap, it creates work for the Garbage Collector. If the struct is small (like a 3D coordinate with 3 floats), copying it on the Stack is actually **10x to 100x faster** than putting it on the Heap.

Only return pointers for large structs, or when mutation is absolutely necessary!
