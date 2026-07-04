# Middleware

Middleware is a crucial architectural pattern for web servers. It allows you to intercept incoming HTTP requests, execute common logic (like Authentication, Logging, or CORS headers), and then pass the request down to the actual handler.

In Go, Middleware is implemented using the **Decorator Pattern**.

## 1. The Decorator Pattern

A middleware in Go is simply a function that takes an `http.Handler` as an argument, and returns a new `http.Handler` that wraps the original one.

```go
// 1. The Middleware Function
func LoggingMiddleware(next http.Handler) http.Handler {
    
    // Return a new Handler (using HandleFunc for convenience)
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        start := time.Now()
        
        // Execute the next handler in the chain
        next.ServeHTTP(w, r)
        
        // This runs AFTER the handler finishes!
        duration := time.Since(start)
        fmt.Printf("[%s] %s took %v\n", r.Method, r.URL.Path, duration)
    })
}
```

## 2. Applying Middleware

To apply middleware, you simply wrap your handlers when registering them with the router.

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome!"))
}

func main() {
    mux := http.NewServeMux()
    
    // Wrap the handler!
    finalHandler := LoggingMiddleware(http.HandlerFunc(homeHandler))
    
    mux.Handle("/", finalHandler)
    http.ListenAndServe(":8080", mux)
}
```

## 3. Passing Data using Context

The most powerful use of middleware is Authentication. If a user provides a valid JWT Token, the middleware needs to extract their `UserID` and pass it down to the final handler so the handler knows who is logged in.

How do we pass data through `next.ServeHTTP(w, r)` without changing the function signature? 
**We use the `r.Context()`!**

```go
// Create a custom type for the context key to avoid collisions
type contextKey string
const userIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        
        token := r.Header.Get("Authorization")
        if token != "valid-token" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return // Abort the request! (Do NOT call next)
        }

        // 1. Inject data into a new context
        ctx := context.WithValue(r.Context(), userIDKey, "user_123")
        
        // 2. Clone the request with the new context
        reqWithCtx := r.WithContext(ctx)
        
        // 3. Pass the new request down the chain
        next.ServeHTTP(w, reqWithCtx)
    })
}
```

Inside the final handler, we extract the data:

```go
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
    // Extract the User ID from the context!
    userID := r.Context().Value(userIDKey).(string)
    
    fmt.Fprintf(w, "Welcome to your dashboard, %s!", userID)
}
```
