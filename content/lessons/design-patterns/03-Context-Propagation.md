# Context Propagation Pattern

In modern Go architectures, no function operates in a vacuum. A simple database query is actually part of a massive HTTP request that has a strict 5-second timeout and a unique Trace ID.

To guarantee that timeouts, cancellations, and observability traces survive the journey from the API Gateway down to the deepest database query, you must master the **Context Propagation Pattern**.

## 1. The Golden Rule of Context

The rule is non-negotiable: **`context.Context` must always be the first parameter of any function that performs I/O or takes significant time.**

```go
// BAD: Missing Context. This function can never be cancelled!
func FetchUser(id int) (*User, error)

// BAD: Context is inside a struct. Hard to track and overrides easily.
type Worker struct { ctx context.Context }

// PERFECT: Context is the explicit first parameter.
func FetchUser(ctx context.Context, id int) (*User, error)
```

## 2. Propagating the Chain

When an HTTP request enters your Go server, the `net/http` package automatically attaches a Context to the request. You must extract it and pass it down.

```go
func UserHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Extract the Context from the HTTP Request
    ctx := r.Context()
    
    // 2. Add a strict 2-second timeout to it
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // 3. Propagate it to the Service layer
    userService.UpdateProfile(ctx, r.Body)
}
```

Inside the Service layer, you propagate it further to the Repository layer.

```go
func (s *UserService) UpdateProfile(ctx context.Context, data []byte) error {
    // Do business logic...
    
    // 4. Propagate it to the Repository layer
    return s.repo.Save(ctx, data)
}
```

Inside the Repository layer, you finally hand the Context to the standard library or database driver!

```go
func (r *UserRepository) Save(ctx context.Context, data []byte) error {
    // 5. The Context reaches its final destination!
    // If the 2-second timeout expired, ExecContext will instantly 
    // sever the TCP connection and return context.DeadlineExceeded.
    _, err := r.db.ExecContext(ctx, "UPDATE users...", data)
    return err
}
```

## 3. Metadata Propagation (Context Values)

Context is also used to propagate request-scoped metadata (Trace IDs, Auth Tokens).

**The String Key Trap:**
Never use a primitive `string` as a Context key. If Package A injects `ctx.WithValue("id", 1)`, and Package B injects `ctx.WithValue("id", 2)`, Package B will silently overwrite Package A's data!

**The Custom Type Solution:**
Always define a private, unexported custom type for your Context keys. This makes collisions mathematically impossible.

```go
package auth

// 1. Define an unexported type
type contextKey string

// 2. Instantiate a constant using that type
const userIDKey = contextKey("user_id")

// 3. Helper to inject the value
func WithUserID(ctx context.Context, id int) context.Context {
    return context.WithValue(ctx, userIDKey, id)
}

// 4. Helper to extract the value
func GetUserID(ctx context.Context) (int, bool) {
    id, ok := ctx.Value(userIDKey).(int)
    return id, ok
}
```
Now, other packages can safely call `auth.GetUserID(ctx)` without ever worrying about key collisions!
