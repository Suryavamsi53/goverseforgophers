# Chain of Responsibility

The Chain of Responsibility is a behavioral design pattern that lets you pass requests along a chain of handlers. Upon receiving a request, each handler decides either to process the request or to pass it to the next handler in the chain.

*(Note: We already touched on this in Lesson 10 (Decorator). The HTTP Middleware pattern is essentially a Chain of Responsibility!)*

## 1. The Implementation

The core idea is that each Handler contains a pointer to the `next` Handler in the chain.

```go
// 1. The Request Object
type Request struct {
    Path          string
    IsAuth        bool
    IsRateLimited bool
}

// 2. The Handler Interface
type Handler interface {
    SetNext(Handler)
    Execute(*Request) error
}
```

## 2. The Concrete Handlers

We create specific handlers. If a handler detects an error, it instantly returns (breaking the chain). If it passes, it calls `next.Execute()`.

```go
type AuthHandler struct {
    next Handler
}
func (h *AuthHandler) SetNext(next Handler) { h.next = next }
func (h *AuthHandler) Execute(r *Request) error {
    if !r.IsAuth {
        return errors.New("unauthorized: missing token")
    }
    fmt.Println("Auth passed.")
    if h.next != nil {
        return h.next.Execute(r)
    }
    return nil
}

type RateLimitHandler struct {
    next Handler
}
func (h *RateLimitHandler) SetNext(next Handler) { h.next = next }
func (h *RateLimitHandler) Execute(r *Request) error {
    if r.IsRateLimited {
        return errors.New("429: Too Many Requests")
    }
    fmt.Println("Rate limit passed.")
    if h.next != nil {
        return h.next.Execute(r)
    }
    return nil
}
```

## 3. Wiring the Chain

In the `main` function (or factory), we wire the objects together in the correct sequence.

```go
func main() {
    // 1. Create Handlers
    auth := &AuthHandler{}
    rate := &RateLimitHandler{}
    
    // 2. Build the Chain! (Auth -> RateLimit -> Final Logic)
    auth.SetNext(rate)
    
    // 3. Execute!
    req := &Request{IsAuth: true, IsRateLimited: false}
    err := auth.Execute(req)
    if err != nil {
        fmt.Println("Request Failed:", err)
    }
}
```

## 4. Why use this?

The Chain of Responsibility is vital for **Decoupling**.

If the `AuthHandler` directly hardcoded a call to `RateLimitHandler`, they would be permanently fused together. 
By using an Interface and `SetNext`, you can dynamically rearrange the chain at runtime! 

* Maybe the `/login` route uses: `RateLimit -> Auth`
* Maybe the `/admin` route uses: `Auth -> RateLimit -> AdminCheck -> AuditLogger`

You construct entirely different pipelines out of the same reusable Lego blocks without modifying a single line of the internal handler code.
