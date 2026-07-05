# The Circuit Breaker Pattern

In a distributed system, relying on external APIs (like Stripe, SendGrid, or other internal microservices) is dangerous. 

If the Stripe API goes down entirely, your HTTP requests will return a `500 Error` instantly. This is bad, but it's fast.

The true nightmare scenario is when Stripe does not go down, but instead becomes **extremely slow**.

## 1. The Cascading Failure

If Stripe normally takes `50ms`, but suddenly degrades and takes `10 seconds`, what happens to your Go Web Server?

1. User 1 hits checkout. The Goroutine calls Stripe. It blocks for 10 seconds.
2. Within those 10 seconds, 5,000 other users hit checkout.
3. You now have 5,000 Goroutines blocked in memory, all waiting for Stripe.
4. Your server runs out of RAM or exhausts all Database Connection limits.
5. Your entire web server crashes. 

Because Stripe was slow, your *entire* infrastructure collapsed. This is a **Cascading Failure**.

## 2. The Circuit Breaker

To prevent this, we borrow a concept from electrical engineering: The Circuit Breaker.

If a power surge hits your house, the physical circuit breaker trips (opens), cutting off electricity to prevent your house from burning down. In software, a Circuit Breaker wraps external API calls and cuts them off if they start failing or slowing down.

### The 3 States
1. **CLOSED**: Normal operation. Requests flow through to the API. The Breaker counts successes and failures.
2. **OPEN (Tripped)**: If the failure rate (or latency) exceeds a threshold (e.g., 50% of requests fail in a 10-second window), the Breaker trips to the OPEN state. 
   * *Crucial*: While OPEN, any new request instantly returns an error (e.g., `ErrCircuitOpen`) **without** actually making the HTTP request! This instantly relieves pressure on the slow API and prevents your Goroutines from blocking!
3. **HALF-OPEN**: After a cooldown period (e.g., 30 seconds), the Breaker lets a *single* test request through. If it succeeds, the API is healthy again, and the Breaker resets to CLOSED. If it fails, the API is still broken, and it snaps back to OPEN.

## 3. Implementing in Go

You should never write this yourself. The industry standard Go library for this is `github.com/sony/gobreaker` (written by Sony).

```go
import "github.com/sony/gobreaker"

var cb *gobreaker.CircuitBreaker

func init() {
    settings := gobreaker.Settings{
        Name:        "StripeAPI",
        MaxRequests: 1,                // Requests allowed in Half-Open state
        Interval:    10 * time.Second, // Clearing interval for counts
        Timeout:     30 * time.Second, // Cooldown before trying Half-Open
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            // Trip if more than 50% of requests fail, AND we have at least 10 requests
            failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
            return counts.Requests >= 10 && failureRatio >= 0.5
        },
    }
    cb = gobreaker.NewCircuitBreaker(settings)
}

func ChargeCard() error {
    // Wrap the dangerous HTTP call inside the Breaker
    _, err := cb.Execute(func() (interface{}, error) {
        return http.Get("https://api.stripe.com/charge")
    })
    
    return err
}
```

If Stripe degrades, `cb.Execute` will instantly return a circuit open error, saving your Goroutines and keeping your server alive.
