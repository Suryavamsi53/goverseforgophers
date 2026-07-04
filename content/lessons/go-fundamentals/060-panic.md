# Panics and Recover

While Go doesn't use Exceptions for normal control flow, it does have a mechanism for catastrophic, unrecoverable failures: **Panic**.

## 1. What is a Panic?

A panic instantly halts the normal execution of the current goroutine. It runs any deferred functions, and then crashes the program, printing a stack trace.

```go
func main() {
    fmt.Println("Starting...")
    panic("critical system failure")
    fmt.Println("This will never print") // Unreachable
}
```
Go will automatically panic if you do something illegal, such as:
1. Dereferencing a `nil` pointer.
2. Accessing an array index out of bounds.
3. Writing to a closed channel.

## 2. When to Panic (Crash Early, Crash Hard)

In Java, it's common to throw exceptions for simple things like "file not found". 
**In Go, you should almost NEVER use panic.** You should return an `error`.

You only use `panic` when a bug is so severe that it is unsafe for the program to continue running.
* **Example:** Your application requires a `.env` configuration file to start up. If the file is missing, the app cannot function. This is a valid reason to `panic()` during initialization.

## 3. Recover (The Safety Net)

If you are writing a web server, a bug in one user's HTTP request (like a nil pointer) shouldn't panic and crash the entire server for everyone else!

To catch a panic and restore normal execution, you use the built-in `recover()` function inside a `defer` block.

```go
func safeExecute() {
    // 1. Setup the recovery safety net
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()

    // 2. Do dangerous work
    fmt.Println("Doing work...")
    panic("Unexpected bug!")
    
    // 3. Execution stops at the panic, jumps to the defer, and then safely exits the function.
}

func main() {
    safeExecute()
    fmt.Println("Application is still running!")
}
```

### 🧠 Middleware Architecture
In production Go servers (like Gin or Echo), there is always a global "Recovery Middleware". This middleware wraps every single HTTP request in a `defer recover()` block. If a specific request panics due to a developer bug, the middleware catches it, logs the stack trace, returns a `500 Internal Server Error` to that specific user, and keeps the server alive!
