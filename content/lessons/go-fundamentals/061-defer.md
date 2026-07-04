# The `defer` Keyword

The `defer` statement pushes a function call onto a list. The list of saved calls is executed automatically after the surrounding function returns.

## 1. Syntax and the LIFO Stack

When you defer multiple functions, they are executed in **LIFO (Last-In, First-Out)** order. Think of it like stacking plates: the last plate you put on top is the first one you take off.

```go
func countdown() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    
    fmt.Println("Blastoff!")
}
// Output: Blastoff!, 3, 2, 1
```

## 2. Argument Evaluation (The Trap)

This is the most common bug developers encounter with `defer`. 

When you defer a function, the function itself doesn't execute until the end. **However, the arguments passed to the function are evaluated IMMEDIATELY.**

```go
func main() {
    x := 10
    
    // The value of 'x' is evaluated right here! 
    // It captures '10', not what 'x' becomes later.
    defer fmt.Println("Deferred:", x) 
    
    x = 99
    fmt.Println("Normal:", x)
}
// Output:
// Normal: 99
// Deferred: 10
```
If you need the deferred function to evaluate the *final* state of the variable, you must wrap it in an anonymous closure, because closures capture variables by reference!

```go
defer func() {
    fmt.Println("Deferred Closure:", x) // Prints 99!
}()
```

## 3. Best Practices (Cleanup)

`defer` is almost exclusively used for cleanup tasks to ensure you don't leak resources. By putting the cleanup logic immediately after the initialization, you guarantee it runs no matter how many `if` statements or `return` blocks occur below it.

```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    
    // GUARANTEED to close the file, even if we panic or return early below!
    defer file.Close() 

    // ... 100 lines of complex processing logic ...
    
    return nil
}
```

## 4. Under the Hood: Performance

In early versions of Go, `defer` added significant CPU overhead because it had to dynamically allocate a linked list of function pointers on the Heap.

**In Go 1.14+, the compiler introduced "Open-Coded Defers".** 
If your function has fewer than 8 defers and doesn't use loops, the compiler completely removes the `defer` keyword and statically injects the cleanup code directly before every single `return` statement in your function at compile time. 

Result? `defer` is now virtually instantaneous (zero-cost abstraction).
