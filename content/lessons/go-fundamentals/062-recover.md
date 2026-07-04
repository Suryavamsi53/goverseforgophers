# Recover (Deep Dive)

While we touched on `recover()` in the Panics lesson, it has strict rules about *where* and *how* it must be called to actually work.

## 1. The Execution Rule

The `recover()` built-in function **must be called directly inside a deferred function**. 

If you call `recover()` in normal code, it does nothing and returns `nil`. If you call it inside a nested function inside the defer, it also fails.

**❌ Bad (Called normally - Fails to catch panic):**
```go
func badCatch() {
    recover() // Does nothing here
    panic("boom")
}
```

**❌ Bad (Nested too deep - Fails to catch panic):**
```go
func nestedCatch() {
    defer func() {
        func() {
            recover() // Fails! Must be called directly by the deferred function.
        }()
    }()
    panic("boom")
}
```

**✅ Good (Directly inside defer):**
```go
func goodCatch() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Caught:", r)
        }
    }()
    panic("boom")
}
```

## 2. Returning the Panic as an Error

If a function panics but you successfully recover from it, how do you tell the caller that something went wrong?

If you don't do anything, the caller will assume the function succeeded and returned its default zero-values! 

To properly convert a Panic into an Error, you **must use Named Return Variables**.

```go
// We name the return variable 'err'
func parseData() (err error) {
    defer func() {
        // If a panic occurs, we catch it
        if r := recover(); r != nil {
            // We modify the named return variable directly!
            err = fmt.Errorf("system panicked: %v", r)
        }
    }()

    // Imagine this slice lookup goes out of bounds and panics
    data := []int{1}
    fmt.Println(data[10]) 

    return nil
}

func main() {
    err := parseData()
    if err != nil {
        fmt.Println("Handled safely:", err)
    }
}
```

This pattern is heavily used in Go libraries (like JSON parsers) to allow internal recursive logic to `panic` on bad data for speed, while the top-level exported function wraps it in a `defer recover()` and hands a clean `error` back to the user.
