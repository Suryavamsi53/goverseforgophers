# Timeouts and Deadlines

In Lesson 19 (Context), we learned the mechanics of using `context.WithTimeout` to abort an HTTP request if it takes too long. 

In a distributed architecture, timeouts are not just a convenient feature; they are the ultimate fail-safe that prevents your entire cluster from running out of memory.

## 1. The Endless Wait

If you do not specify a timeout on an HTTP Client or a Database Connection, Go will use the OS default TCP timeout (which can be several minutes, or sometimes infinite).

If a downstream service silently hangs (e.g., its router drops the packet but doesn't send a TCP `RST` back), your Goroutine will sit in memory forever. If this happens to 100,000 requests, your server crashes.

**Rule: Every single network call in your application MUST have an explicit timeout.**

## 2. Setting Timeouts in Go

### The Bad Way (Client Timeout)
```go
client := &http.Client{
    Timeout: 5 * time.Second,
}
client.Do(req)
```
This is better than nothing, but it is dangerous. The `Client.Timeout` covers the *entire* exchange (DNS lookup, TCP Handshake, TLS Handshake, writing the body, and reading the body). If you need to read a massive 5GB file, the timeout will trigger and kill the download midway.

### The Correct Way (Context Timeouts)
Always use Contexts to control timeouts. 

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
http.DefaultClient.Do(req)
```

## 3. Global Deadlines (Distributed Tracing)

Imagine this architecture:
`Mobile App -> API Gateway -> Order Service -> Billing Service -> Stripe`

The Mobile App sets a timeout of 5 seconds. If the request isn't done in 5 seconds, the Mobile App gives up and shows an error to the user.

If the API Gateway and Order Service take 4 seconds, and then the Billing Service calls Stripe... the Billing Service should **not** use a hardcoded 5-second timeout for Stripe! The user has only 1 second left before they drop off! If the Billing service waits 5 seconds, it is wasting 4 seconds of compute time doing work for a user who already closed the app!

### Deadline Propagation
Instead of creating a *new* timeout at every microservice, you propagate a **Global Deadline**.

The API Gateway creates a context with a Deadline of 5 seconds.
`ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))`

When the API Gateway calls the Order Service (via gRPC or HTTP headers), it passes the Deadline over the network.
When the Order Service receives the request, it extracts the Context. By the time it calls the Billing Service, the Context automatically knows only 1 second remains! 

If Stripe takes 1.1 seconds, the Context triggers a cancellation signal all the way back up the chain instantly, saving resources across all 4 microservices simultaneously. gRPC handles this Context propagation automatically over the wire!
