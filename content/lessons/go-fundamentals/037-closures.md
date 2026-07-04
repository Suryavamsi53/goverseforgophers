# Closures and State Capture

A **Closure** is a special type of anonymous function that "closes over" and captures variables from the outside scope where it was defined. 

Even after the outer function finishes executing and its variables should technically be destroyed, the closure retains a persistent reference to them.

## 1. Basic State Capture

In this example, the `counter` function returns an anonymous closure function.

```go
func counter() func() int {
    count := 0 // Local variable
    
    // Return a closure that captures the 'count' variable
    return func() int {
        count++
        return count
    }
}

func main() {
    increment := counter()
    
    fmt.Println(increment()) // 1
    fmt.Println(increment()) // 2
    fmt.Println(increment()) // 3
}
```

### 🧠 Escape Analysis (Under the Hood)
Wait—when `counter()` finishes executing, shouldn't its local variable `count` be destroyed and popped off the stack? 
Yes, normally! But the Go compiler runs an **Escape Analysis** pass. It detects that the inner closure still needs access to `count` after the function returns. Therefore, the compiler silently moves `count` from the Stack to the **Heap**, keeping it alive in memory!

## 2. The Loop Variable Capture Trap

The single most notorious bug in Go history involves closures and loop variables. 

**Prior to Go 1.22**, loop variables were reused for every iteration. If you spawned a goroutine (which takes time to start) inside a loop, it captured the reference to the loop variable, not a copy of its value.

```go
// ⚠️ THE HISTORICAL TRAP (Before Go 1.22)
func main() {
    funcs := []func(){}
    
    for i := 0; i < 3; i++ {
        funcs = append(funcs, func() {
            fmt.Println(i) // Captures a pointer to 'i'
        })
    }

    for _, f := range funcs {
        f() // Output was 3, 3, 3! (Because 'i' ended at 3)
    }
}
```
To fix this, developers historically had to create a local copy inside the loop: `val := i`.

### ✅ The Fix in Go 1.22+
The Go core team recognized this was causing too many production bugs. **In Go 1.22, the language spec was fundamentally changed.** 

Loop variables are now newly instantiated on every single iteration. If you run the exact code above in Go 1.22+, it correctly outputs `0, 1, 2` because the closure captures a fresh, distinct memory address on every loop!
