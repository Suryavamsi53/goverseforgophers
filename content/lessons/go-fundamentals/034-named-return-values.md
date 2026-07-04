# Named Return Values

Go allows you to name your return variables directly in the function signature. When you do this, Go automatically creates these variables and initializes them to their zero values before the function even begins.

## 1. Basic Syntax and the "Naked Return"

```go
// We declare 'x' and 'y' in the return signature
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    
    // "Naked return" - automatically returns the current values of x and y
    return 
}
```

Because `x` and `y` are pre-declared, we don't use `:=` inside the function, we just assign to them using `=`. 
When we type `return` without specifying variables, it is called a "naked return".

## 2. Why use Named Returns?

Named returns are controversial. While naked returns can make small functions look clean, they make long functions incredibly difficult to read, because the developer has to scroll back up to the function signature to remember what is actually being returned.

However, named returns have two massive benefits:

### Benefit A: Self-Documenting Code
If a function returns two floats, it's hard to know what they represent. Named returns act as perfect, built-in documentation.

```go
// Confusing
func getCoordinates() (float64, float64) { ... }

// Perfect clarity
func getCoordinates() (lat float64, lng float64) { ... }
```

### Benefit B: Manipulation via `defer`
This is the true architectural superpower of named returns. 
In Go, `defer` blocks execute *after* the function completes, but *before* the return values are actually handed back to the caller. 
**If you use named returns, a `defer` block can intercept and modify the return values!**

```go
func doWork() (err error) {
    defer func() {
        if r := recover(); r != nil {
            // We intercepted the panic and injected an error into the named return variable!
            err = fmt.Errorf("system panicked: %v", r)
        }
    }()
    
    // Simulating a crash
    panic("database exploded")
    
    return nil
}
```
Without the named return `(err error)`, the defer block would have no way to access the return variable to modify it.

## 3. The Shadowing Trap

Be extremely careful using the `:=` operator if you have named returns.

```go
func calculate() (result int) {
    // ❌ BUG: This creates a NEW local variable named 'result', shadowing the named return!
    result := 10 
    
    // This naked return will return 0 (the zero value of the original named return) 
    // because the local shadowed variable is ignored!
    return 
}
```
