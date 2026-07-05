# Idempotency

In Lesson 8, we learned that we must retry network requests when they fail. 

But what if the network fails *after* the remote server has already processed the request, but *before* the remote server can send the `200 OK` response back to us?

If we retry a "Charge Credit Card" request, and the original request actually succeeded, we will charge the user twice! This is illegal. 

To solve this, all distributed API endpoints must be **Idempotent**.

## 1. What is Idempotency?

An operation is idempotent if making the same request multiple times has the exact same effect as making it once.

* `GET /users/1` is naturally idempotent. You can read the user 100 times, the database doesn't change.
* `PUT /users/1` is idempotent. Setting the user's name to "John" 100 times just overwrites it with "John".
* `POST /orders` is **NOT** idempotent. Hitting it 100 times creates 100 new orders!

## 2. The Idempotency Key (Stripe Pattern)

To make a `POST` request idempotent, the client must generate a unique identifier (usually a UUID v4) called an **Idempotency Key**.

The client sends this key in the HTTP Headers:
`Idempotency-Key: a1b2c3d4-5678-90ab-cdef-1234567890ab`

### The Backend Implementation
When the Go backend receives the request, it checks a fast Key-Value store (like Redis):

1. **Check Redis**: Does `Idempotency-Key: a1b2...` exist?
2. **If NO**: 
   * Store the key in Redis with a status of `PROCESSING`.
   * Charge the card via Stripe.
   * Update the key in Redis to `DONE` and save the Stripe Receipt URL.
   * Return the `200 OK` to the client.
3. **If YES (Status: PROCESSING)**:
   * The client double-clicked the checkout button!
   * Return a `409 Conflict` (or ask them to wait).
4. **If YES (Status: DONE)**:
   * The client is retrying a dropped connection. The card was already charged!
   * Do NOT hit Stripe again. Simply return the cached `200 OK` and the saved Stripe Receipt URL from Redis.

```go
func HandleCheckout(w http.ResponseWriter, r *http.Request) {
    idempKey := r.Header.Get("Idempotency-Key")
    if idempKey == "" {
        http.Error(w, "Idempotency-Key required", 400)
        return
    }

    // 1. Try to set the key in Redis (SETNX)
    acquired := redis.SetNX(ctx, idempKey, "PROCESSING", 24*time.Hour)
    if !acquired {
        // We have seen this request before! Handle safely!
        status := redis.Get(ctx, idempKey)
        if status == "DONE" {
            w.Write([]byte("Order was already processed successfully!"))
            return
        }
    }

    // 2. Actually process the payment
    processPayment()

    // 3. Mark as DONE
    redis.Set(ctx, idempKey, "DONE", 24*time.Hour)
}
```

## 3. At-Least-Once Delivery

Idempotency is absolutely mandatory when working with Apache Kafka or RabbitMQ.

Message Queues guarantee **At-Least-Once Delivery**. This means Kafka promises it will deliver your message... but occasionally, due to network partitions, it will accidentally deliver the exact same message twice to your Go consumer. 

Your Go consumer must use an Idempotency Key (like the `OrderID`) to ensure it doesn't process the Kafka event a second time.
