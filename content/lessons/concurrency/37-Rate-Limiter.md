# Rate Limiting

A Semaphore limits how many tasks can run **at the exact same time** (Concurrency). 
A Rate Limiter limits how many tasks can run **over a period of time** (Throughput).

If you call the Stripe API 100 times in 1 second, Stripe will ban your IP address, even if you only ran 5 requests concurrently using a Semaphore. You need a Rate Limiter to enforce a speed limit (e.g., max 10 requests per second).

## 1. The Token Bucket Algorithm

The most robust rate-limiting algorithm is the **Token Bucket**.

Imagine a physical bucket that holds 10 tokens.
* Every time a user makes an HTTP request, they must take 1 token from the bucket.
* If the bucket is empty, the user's request is rejected (HTTP 429 Too Many Requests).
* A background timer automatically drops a new token into the bucket every 100 milliseconds.

This allows for short bursts of traffic (instantly grabbing 10 tokens), but strictly enforces a long-term average (10 requests per second).

## 2. Implementing with `time.Ticker`

We can build a basic Rate Limiter in Go using a Buffered Channel (the bucket) and a Ticker (the token generator).

```go
func main() {
    // 1. The Bucket (Capacity: 5 tokens for burst)
    rateLimiter := make(chan time.Time, 5)
    
    // Fill the bucket initially
    for i := 0; i < 5; i++ {
        rateLimiter <- time.Now()
    }

    // 2. The Refiller
    // Drops a new token into the bucket every 200ms (5 req/sec average)
    go func() {
        for t := range time.Tick(200 * time.Millisecond) {
            // Non-blocking send! If the bucket is already full at 5, drop the token.
            select {
            case rateLimiter <- t:
            default:
            }
        }
    }()

    // 3. The Clients
    requests := make(chan int, 100)
    for i := 1; i <= 20; i++ { requests <- i }
    close(requests)

    for req := range requests {
        // Wait for a token to become available
        <-rateLimiter 
        fmt.Println("Processed request", req, "at", time.Now())
    }
}
```

## 3. The `golang.org/x/time/rate` Package

While building it yourself is a great learning exercise, you should never write your own rate limiter for production.

Google provides the highly-optimized `golang.org/x/time/rate` package. It implements a mathematically perfect Token Bucket and handles Context cancellation effortlessly.

```go
import "golang.org/x/time/rate"

// Limit: 5 requests per second. Burst: up to 10 instantly.
limiter := rate.NewLimiter(5, 10)

func handleAPIRequest(ctx context.Context) error {
    // Blocks until a token is available. 
    // Automatically aborts if the context times out first!
    err := limiter.Wait(ctx)
    if err != nil {
        return err // Rate limit exceeded or context cancelled
    }
    
    // Safe to hit the API!
    return nil
}
```
