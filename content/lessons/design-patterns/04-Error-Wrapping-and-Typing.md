# Error Wrapping and Typing

In early versions of Go (pre-1.13), if a database query failed, developers would use `fmt.Errorf` to add context:
`return fmt.Errorf("failed to get user: %v", err)`

The problem with this approach is that `fmt.Errorf` concatenates the errors into a giant string and completely destroys the original error type. If the caller wanted to know if the error was a `sql.ErrNoRows`, they had to perform string parsing (`strings.Contains(err.Error(), "no rows")`), which is incredibly brittle.

Go 1.13 introduced the **Error Wrapping Pattern** to solve this.

## 1. The `%w` Verb

To add context to an error without destroying its type, you use the `%w` verb in `fmt.Errorf`. This creates a wrapped error chain (like a Russian nesting doll).

```go
func GetUser(id int) (*User, error) {
    err := db.QueryRow("...").Scan(...)
    if err != nil {
        if err == sql.ErrNoRows {
            // We use %w to wrap the original sql.ErrNoRows!
            return nil, fmt.Errorf("user %d not found: %w", id, err)
        }
        return nil, fmt.Errorf("database query failed: %w", id, err)
    }
    return user, nil
}
```

## 2. Unwrapping with `errors.Is`

Because we used `%w`, the caller can now inspect the chain to see if `sql.ErrNoRows` exists anywhere inside it, completely bypassing brittle string parsing!

```go
user, err := GetUser(42)
if err != nil {
    // Is checks the entire wrapped chain!
    if errors.Is(err, sql.ErrNoRows) {
        http.Error(w, "User does not exist", 404)
        return
    }
    
    http.Error(w, "Internal Server Error", 500)
    return
}
```

## 3. Custom Error Types (`errors.As`)

Sometimes, standard errors aren't enough. You want an error that carries metadata (like an HTTP Status Code). You create a Custom Error Type by implementing the `Error()` method.

```go
// 1. Define a custom struct
type HTTPError struct {
    StatusCode int
    Message    string
}

// 2. Implement the error interface
func (e *HTTPError) Error() string {
    return fmt.Sprintf("status %d: %s", e.StatusCode, e.Message)
}

// 3. Return it from a function
func Validate() error {
    return &HTTPError{StatusCode: 400, Message: "Invalid email"}
}
```

How does the caller extract the `StatusCode` from the generic `error` interface? You use `errors.As`.

```go
err := Validate()

// We create an empty pointer of the type we are looking for
var httpErr *HTTPError

// As() searches the error chain. If it finds an HTTPError, 
// it injects it into our httpErr variable!
if errors.As(err, &httpErr) {
    // We now have full access to the struct fields!
    w.WriteHeader(httpErr.StatusCode)
    w.Write([]byte(httpErr.Message))
}
```

## 4. Sentinel Errors

A Sentinel Error is a pre-defined, global error variable. It is a common pattern in the standard library (e.g., `io.EOF`, `sql.ErrNoRows`).

```go
var ErrUserNotFound = errors.New("user not found")
var ErrInsufficientFunds = errors.New("insufficient funds")
```

**The Rule of Thumb**:
* Use **Sentinel Errors** (`errors.Is`) for simple, static errors (like "Not Found").
* Use **Custom Error Types** (`errors.As`) when you need dynamic metadata (like "Timeout after X seconds" or "HTTP Status Code Y").
* Use **`%w`** everywhere to build the chain and preserve the stack!
