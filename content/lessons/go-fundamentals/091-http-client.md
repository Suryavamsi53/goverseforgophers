# HTTP Client

Go was built for the modern web. The `net/http` package provides a remarkably robust HTTP client right out of the box, with built-in support for HTTP/2, connection pooling, and timeouts.

## 1. The DefaultClient Trap

The easiest way to make a network request is `http.Get()`.

```go
resp, err := http.Get("https://api.github.com")
```

**⚠️ NEVER USE THIS IN PRODUCTION.**

`http.Get` uses the global `http.DefaultClient`. The DefaultClient has **no timeout configured**. If the server you are calling goes offline and drops your packets instead of closing the connection, your goroutine will hang there, blocking forever. If this happens to 10,000 requests, your server will run out of file descriptors and crash.

## 2. Building a Custom Client

You should always instantiate your own `http.Client` with a strict `Timeout`.

```go
func main() {
    client := &http.Client{
        Timeout: 5 * time.Second, // Hard deadline for the entire request
    }

    resp, err := client.Get("https://api.github.com")
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Status Code:", resp.StatusCode)
}
```

## 3. The Connection Leak Trap

Notice the `defer resp.Body.Close()` in the code above? 

Under the hood, Go maintains a pool of open TCP connections to remote servers (Keep-Alive). If you read the response body but forget to call `.Close()`, Go cannot return that TCP socket back to the connection pool. 
This is a **Socket Leak**, and it will rapidly exhaust your server's network ports. 

*Rule: Always defer `Close()` immediately after checking `err != nil`.*

## 4. Advanced Requests (`http.NewRequest`)

If you need to set custom headers (like Authorization tokens) or send a JSON payload, you cannot use `.Get()`. You must construct a Request object first.

```go
import (
    "bytes"
    "net/http"
)

func main() {
    client := &http.Client{Timeout: 5 * time.Second}

    // 1. Prepare JSON payload
    payload := []byte(`{"name": "Alice"}`)

    // 2. Build the Request
    req, err := http.NewRequest("POST", "https://api.example.com/users", bytes.NewBuffer(payload))
    if err != nil {
        panic(err)
    }

    // 3. Add Headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer my-secret-token")

    // 4. Execute the Request
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
}
```

### 🧠 Performance Insight (Transport Layer)
If you are building a high-throughput microservice that talks to another service 10,000 times a second, you can customize the client's `http.Transport`. You can increase `MaxIdleConnsPerHost` from the default of 2 up to 1,000, ensuring Go keeps hundreds of TCP sockets warm and ready to fire, bypassing the costly TCP/TLS handshake latency entirely!
