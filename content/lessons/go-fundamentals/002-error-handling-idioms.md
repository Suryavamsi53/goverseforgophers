# Error Handling Idioms (`err != nil`)

If you come from Java, Python, or JavaScript, you are used to `try/catch` blocks. When a function fails, it throws an Exception, which halts the normal flow of the program and leaps up the call stack until a `catch` block intercepts it.

Go completely rejects this paradigm. **In Go, Errors are just standard values.**

## 1. The `error` Interface

An error in Go is nothing special. It is simply any struct that implements the built-in `error` interface:

```go
type error interface {
    Error() string
}
```

When a function can fail, it simply returns the error as the last return value.

```go
// Returns the data AND the error!
func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil
}
```

## 2. The `if err != nil` Pattern

Because errors are just values, you must explicitly check them immediately after the function returns.

```go
result, err := Divide(10, 0)
if err != nil {
    // Handle the error! (Log it, return it, or crash)
    log.Printf("Math failed: %v", err)
    return
}
fmt.Println("Result:", result)
```

Many developers complain that writing `if err != nil` hundreds of times is tedious. However, it forces you to consciously consider the failure state of every single operation, rather than blindly wrapping a massive block of code in a `try/catch` and praying. It makes Go code incredibly robust and predictable.

## 3. Error Wrapping (Context)

If a Database query fails, the low-level SQL driver returns `connection refused`.
If you return that raw error to the HTTP handler, the logs will just say `connection refused`. You have no idea *which* query failed!

You must add Context to your errors using **Wrapping** (`fmt.Errorf` with the `%w` verb).

```go
func GetUser(id int) (*User, error) {
    err := db.Query("...")
    if err != nil {
        // We wrap the original error with context!
        return nil, fmt.Errorf("failed to fetch user %d from db: %w", id, err)
    }
}
```
Now, the log output becomes a beautiful chain: `failed to fetch user 42 from db: connection refused`.

## 4. Inspecting Wrapped Errors (`errors.Is` and `errors.As`)

If you wrap an error, how does the HTTP handler check if the underlying root cause was a `sql.ErrNoRows` so it can return a 404 instead of a 500?

You cannot use standard equality (`err == sql.ErrNoRows`) because the error is buried inside the wrapper!

You must use `errors.Is()`:

```go
err := GetUser(99)

// errors.Is recursively unwraps the error chain to see if the target exists anywhere inside it!
if errors.Is(err, sql.ErrNoRows) {
    http.Error(w, "User Not Found", 404)
    return
}
```

If you need to extract a specific Custom Error Struct (e.g., to read a specific HTTP Status Code field embedded in the error), you use `errors.As()`.

## 5. What about `panic`?

Go does have a mechanism similar to throwing an Exception: `panic()`.
When you call `panic("boom")`, the entire Goroutine crashes instantly.

**Enterprise Rule:** You should almost NEVER use `panic`. It is strictly reserved for truly unrecoverable, catastrophic states (e.g., the application boots up, but the required configuration file is missing). For all standard business logic failures, you must return an `error`.
