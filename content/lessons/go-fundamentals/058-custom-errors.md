# Custom Errors

Because `error` is just an interface requiring an `Error() string` method, you can easily create your own custom error structs that contain rich, contextual data.

## 1. Why `errors.New()` is not enough

`errors.New("user not found")` returns a simple string. But what if the calling function needs to know *why* the user wasn't found? Did the database timeout? Or does the user literally not exist? Parsing a string to figure this out is terrible practice.

## 2. Building a Custom Error Struct

Let's build an `HTTPError` that holds both a message and an HTTP Status Code.

```go
// 1. Define the struct with contextual data
type HTTPError struct {
    StatusCode int
    Message    string
}

// 2. Implement the Error() string method to satisfy the error interface
func (e *HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}
```

## 3. Returning and Asserting Custom Errors

Now, a function can return this rich error struct:

```go
func fetchUser() error {
    // Return our custom error as an interface
    return &HTTPError{
        StatusCode: 404,
        Message:    "User not found in database",
    }
}
```

When the caller receives the error, they can use a **Type Assertion** (or Type Switch) to extract the contextual data (like the Status Code)!

```go
func main() {
    err := fetchUser()
    
    if err != nil {
        // Assert that the error is specifically an *HTTPError
        if httpErr, ok := err.(*HTTPError); ok {
            fmt.Printf("Sending Status %d to client.\n", httpErr.StatusCode)
        } else {
            fmt.Println("Generic error:", err)
        }
    }
}
```
### Architecture Insight
By injecting contextual data (like `Retryable bool`, `Timeout time.Duration`, or `StatusCode int`) directly into custom error structs, your routing layers can automatically decide how to handle failures (e.g., retrying the request or returning a 500 API response) without hardcoding logic!
