# Anonymous Functions

An anonymous function is simply a function that does not have a name. Because functions are first-class citizens in Go, you can define them inline, assign them to variables, or execute them immediately.

## 1. Assigning to Variables

You can assign a nameless function directly to a variable and execute it later.

```go
func main() {
    greet := func(name string) {
        fmt.Println("Hello,", name)
    }

    greet("Alice") // Outputs: Hello, Alice
    greet("Bob")
}
```

## 2. IIFE (Immediately Invoked Function Expressions)

You can define a function and execute it immediately on the exact same line by appending `()` to the end of the definition.

```go
func main() {
    func(msg string) {
        fmt.Println("Executing Immediately:", msg)
    }("System Booting...")
}
```

## 3. Real-World Use Case: The `defer` Wrapper

The most critical and frequent use of anonymous functions in production Go code is wrapping logic inside `defer` statements.

`defer` can only execute a single function call. If you need to execute multiple lines of cleanup code before a function returns, you wrap them in an anonymous IIFE.

```go
func processTransaction() {
    db := connectToDB()
    
    // Use an anonymous function to group cleanup logic
    defer func() {
        fmt.Println("Closing database connection...")
        db.Close()
        fmt.Println("Transaction complete.")
    }()

    fmt.Println("Processing...")
    // When this function returns, the entire anonymous defer block runs.
}
```

## 4. Higher-Order Functions

Anonymous functions are heavily used when passing logic into other functions (like callbacks). 

```go
// execute takes a callback function as an argument
func execute(callback func(int) int, value int) {
    result := callback(value)
    fmt.Println("Result:", result)
}

func main() {
    // Passing an anonymous function directly inline
    execute(func(x int) int {
        return x * x
    }, 5) // Prints 25
}
```
