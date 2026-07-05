# Retries and Exponential Backoff

In a distributed system, network packets will randomly drop, and routers will briefly restart. This means that an HTTP request to another microservice might randomly fail with a `502 Bad Gateway`, but if you try it again exactly 1 millisecond later, it succeeds.

These are called **Transient Failures**.

The solution to Transient Failures is simple: retry the request. However, if you implement retries incorrectly, you will accidentally launch a self-inflicted DDoS attack against your own servers.

## 1. The Retries Storm (DDoS)

Imagine your database goes offline for exactly 2 seconds. 
During those 2 seconds, 5,000 HTTP requests hit your web server. 
The web server's database queries fail. 

If your web server has a naïve retry loop (`while err != nil { retry() }`), all 5,000 Goroutines will start frantically pinging the database at the speed of light. They might execute 100,000 queries per second while the database is rebooting. 
When the database finally comes back online, it is instantly hit by millions of backlogged retries and crashes again!

## 2. Exponential Backoff

To prevent a retry storm, you must space out the retries. Instead of retrying instantly, you wait. 
Crucially, you must wait *exponentially longer* after each failure to give the struggling downstream system time to recover.

* **Attempt 1 Fails**: Wait 100ms, then retry.
* **Attempt 2 Fails**: Wait 200ms, then retry.
* **Attempt 3 Fails**: Wait 400ms, then retry.
* **Attempt 4 Fails**: Wait 800ms, then retry.

This ensures that if the downstream service is offline for a long time, your application gradually stops hammering it.

## 3. The Thundering Herd (Jitter)

Exponential Backoff solves the rapid-fire problem, but it introduces a new one.

If the database drops offline for 2 seconds, and 5,000 requests fail at the exact same nanosecond...
* All 5,000 requests will wait exactly 100ms.
* At exactly 100ms, all 5,000 requests will retry at the exact same nanosecond!
* They fail, wait exactly 200ms, and hit it simultaneously again.

This synchronized wave of traffic is called the **Thundering Herd**. 

To solve this, we add **Random Jitter**. Instead of waiting exactly 100ms, a Goroutine waits `100ms ± random(50ms)`. 
Now, some wait 75ms, some wait 120ms. The retry traffic is perfectly smoothed out over time, preventing the herd effect.

## 4. Implementing in Go

Again, do not write this math yourself. Use the industry standard `github.com/cenkalti/backoff` package.

```go
import "github.com/cenkalti/backoff/v4"

func FetchData() error {
    // Creates an Exponential Backoff object that automatically adds Random Jitter!
    b := backoff.NewExponentialBackOff()
    b.MaxElapsedTime = 5 * time.Second // Give up entirely after 5 seconds

    operation := func() error {
        // The HTTP request you want to retry
        resp, err := http.Get("https://api.example.com")
        if err != nil {
            return err // Returning an error triggers a backoff retry
        }
        if resp.StatusCode >= 500 {
            return fmt.Errorf("server error") // Also trigger a retry for 5xx errors
        }
        
        // Success!
        return nil 
    }

    // Execute the operation using the backoff algorithm
    return backoff.Retry(operation, b)
}
```

## 5. The Golden Rule of Retries

**Never retry a timeout.** 
If an HTTP request takes 10 seconds and times out, retrying it is incredibly dangerous. The original request might still be processing on the remote server! If you retry it, the remote server is now processing the exact same heavy task *twice*, slowing it down even further. Only retry fast network drops (e.g., `Connection Refused` or `502 Bad Gateway`).
