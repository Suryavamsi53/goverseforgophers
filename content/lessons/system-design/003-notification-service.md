# System Design: Notification Service

## 1. Learning Objectives
* **What you'll learn**: How to design a high-throughput, reliable Notification System that can blast millions of Emails, SMS, and Push Notifications without dropping a single message or sending duplicates.
* **Why it matters**: Notifications are the primary driver of user engagement. If the service goes down, engagement drops. If the service sends the same email 5 times, users will flag you as spam and destroy your company's domain reputation.
* **Where it's used**: E-commerce order updates, social media alerts, and marketing blasts.

---

## 2. Real-world Story
Imagine a massive Post Office on December 23rd.
Millions of packages arrive simultaneously. If the Post Office tries to instantly hire 1 million delivery drivers, they will go bankrupt.
Instead, they use **Queues**. They accept all the packages, place them in massive warehouses safely (Message Brokers), and the delivery drivers (Workers) pick them up at a steady, reliable pace. If a driver crashes, the package goes back into the warehouse to be delivered by someone else.

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Billing Service] -->|1. HTTP POST| B(Notification API)
    
    B -->|2. Validates & Enqueues| C((Kafka / RabbitMQ))
    
    C -->|3. Consumes| D[Email Worker (Go)]
    C -->|3. Consumes| E[SMS Worker (Go)]
    
    D -->|4. API Call| F(SendGrid / AWS SES)
    E -->|4. API Call| G(Twilio)
    
    style C fill:#f97316,color:#fff
```

---

## 4. Internal Working (Under the Hood)
The golden rule of Notifications is **Asynchronous Processing**.
When a user registers, the Auth Service sends an event to the Notification Service. The Notification API does NOT call SendGrid directly! External APIs (SendGrid, Twilio) have strict Rate Limits and high latency (500ms).
Instead, the API drops the payload into a Message Broker (Kafka) and instantly returns `202 Accepted`. Background Go Workers pull from the queue and deal with the slow 3rd-party APIs.

---

## 5. Compiler Behavior
* **Goroutine Throttling**: If a Go worker pulls 10,000 emails from Kafka, it shouldn't spawn 10,000 Goroutines to call SendGrid instantly, as SendGrid will block your IP for Rate Limiting. You must use a Go **Worker Pool** (e.g., exactly 50 Goroutines) to throttle the outbound network requests to a highly controlled, mathematically predictable rate.

---

## 6. Memory Management
* **In-Memory vs Persistent Queues**: Never use Go Channels as the primary queue for a Notification System! If the Go server restarts for a deployment, all un-sent notifications in the Channel are permanently erased from RAM. You MUST use a persistent disk-backed broker like RabbitMQ, AWS SQS, or Kafka.

---

## 7. Code Examples

### 🔹 Example 1: The API Handler (Fast)
```go
// The API accepts the request and pushes to the Queue instantly.
func SendNotificationHandler(w http.ResponseWriter, r *http.Request) {
    var req NotificationRequest
    json.NewDecoder(r.Body).Decode(&req)
    
    // Save to RabbitMQ/Kafka
    err := broker.Publish("email_queue", req)
    if err != nil {
        http.Error(w, "Broker down", http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusAccepted) // 202 - We got it, we'll do it later!
}
```

### 🔹 Example 2: The Go Worker Pool (Throttled)
```go
func StartEmailWorkers(numWorkers int) {
    jobs := make(chan NotificationRequest, 1000)
    
    // Boot exactly 'numWorkers' Goroutines (e.g., 50)
    for i := 0; i < numWorkers; i++ {
        go func(workerID int) {
            for req := range jobs {
                // Call SendGrid
                err := sendgrid.Send(req.To, req.Body)
                if err != nil {
                    // Re-queue the job for retry!
                }
            }
        }(i)
    }
    
    // Main loop pulls from RabbitMQ and feeds the internal Go Channel
    for msg := range rabbitmq.Consume("email_queue") {
        jobs <- parseMsg(msg)
    }
}
```

### 🔹 Example 3: Advanced (Idempotency Key)
```go
// Preventing Duplicate Emails!
func ProcessEmail(req NotificationRequest) {
    // Check Redis: Have we successfully sent this specific NotificationID before?
    exists, _ := redis.SetNX(ctx, "sent:"+req.ID, "1", 24*time.Hour).Result()
    if !exists {
        log.Println("Duplicate detected! Skipping.")
        return 
    }
    
    sendgrid.Send(req.To, req.Body)
}
```

### 🔹 Example 4: Production (Rate Limiter Integration)
```go
// Respecting User Preferences (Don't spam them!)
func ShouldSend(userID string) bool {
    // Check Database: Did the user unsubscribe?
    if db.IsUnsubscribed(userID) { return false }
    
    // Check Redis: Have we sent them more than 3 marketing emails today?
    count := redis.Incr(ctx, "daily_limit:"+userID)
    if count > 3 { return false }
    
    return true
}
```

### 🔹 Example 5: Interview
```go
// Q: What happens if SendGrid is completely down for 2 hours?
// A: The Go Workers will fail. They must place the failed messages into a 
// "Dead Letter Queue" (DLQ) or implement Exponential Backoff retries. 
// The messages are safely preserved in the Broker until SendGrid recovers!
```

---

## 8. Production Examples
1. **Template Engine**: Instead of Microservices sending raw HTML to the Notification Service, they send `TemplateID=123` and `Data={"name": "Alice"}`. The Notification Service uses Go's `html/template` to render the beautiful HTML securely in one centralized location.
2. **Third-Party Failover**: If Twilio's SMS API goes down, the Go Worker detects 5 consecutive timeouts, trips a Circuit Breaker, and instantly dynamically routes all SMS traffic to a backup vendor (like Amazon SNS) automatically.

---

## 9. Performance & Benchmarking
* **Batching**: Sending 1,000 separate HTTP requests to SendGrid takes 1,000 network handshakes (Slow). SendGrid provides Batch APIs. The Go Worker pulls 1,000 messages from Kafka, constructs ONE massive JSON payload, and sends it to SendGrid in a single HTTP request (Insanely Fast).

---

## 10. Best Practices
* ✅ **Do**: Make sure Notification IDs are uniquely generated by the caller. If the caller crashes and retries, the ID remains the same, allowing your Idempotency check to block the duplicate email!
* ❌ **Don't**: Build your own SMTP server. Email deliverability (avoiding the Spam Folder, handling DKIM/SPF records) is an entire industry. Pay for SendGrid, Mailgun, or AWS SES.
* 🏢 **Google / Uber / Netflix Style**: Use a **Priority Queue**. Password Reset emails must be delivered in 2 seconds (High Priority Queue). Weekly Marketing blasts can take 6 hours (Low Priority Queue). The Go Workers prioritize the High queue first.

---

## 11. Common Mistakes
1. **Unbounded Retries**: If an email bounces because the address is `test@test.com`, the worker fails and puts it back in the queue. It tries again 1 second later. It will infinitely retry forever, burning CPU and Queue storage (Poison Pill). Always set a `Max_Retries = 3` and then route to a Dead Letter Queue!
2. **Hardcoding Vendor Logic**: Hardcoding Twilio SDK calls directly into the Order Service. If you switch to AWS SNS, you have to rewrite the Order Service. The Notification Service acts as a Facade to hide vendor details from the rest of the company.

---

## 12. Debugging
How to troubleshoot a Notification Service:
* **Message Tracking**: Every notification MUST be saved to a database (Cassandra/Postgres) with a `Status` column (`PENDING`, `SENT`, `FAILED`, `BOUNCED`). When Customer Support gets a ticket, they can search the database and say: "Ah, it says here the email Hard Bounced because your inbox is full."

---

## 13. Exercises
1. **Easy**: Write a Go script that sends an email using an external provider's SDK (e.g., SendGrid).
2. **Medium**: Build an HTTP API that validates a payload and pushes it into a local Go Channel (simulating a Queue) and returns `202 Accepted`.
3. **Hard**: Build a Worker Pool that reads from the channel and prints the emails to the console at a max rate of 5 per second.
4. **Expert**: Replace the Go Channel with RabbitMQ, implement Retry logic, and a Redis Idempotency check.

---

## 14. Quiz
1. **MCQ**: What is the primary purpose of separating the API from the Workers using a Message Queue?
   * (A) To encrypt the emails (B) To decouple burst traffic from slow 3rd-party APIs (C) To compress the payload. *(Answer: B)*
2. **System Design Follow-up**: How do you handle Apple/Android Push Notifications? *(Apple APNS requires maintaining long-lived HTTP/2 or TCP connections. It is highly specialized. Usually, companies use Firebase Cloud Messaging (FCM) as a unified wrapper, so the Go Worker just makes a standard HTTP call to Firebase).*

---

## 15. FAANG Interview Questions
* **Beginner**: Why is an asynchronous architecture necessary for this system?
* **Intermediate**: Explain the exact mechanism of a Dead Letter Queue (DLQ).
* **Senior (Google/Meta)**: Architect a "Notification Preferences" microservice. A user wants SMS for Security alerts, but Push for Marketing. How do you design the database schema and query it in <10ms before sending the blast to 10 million users?

---

## 16. Mini Project
**The Smart Retry Worker**
* Build a Go application with a local SQLite Database (`notifications` table).
* Build a Worker that scans the table for `status = 'PENDING'`.
* Simulate a 3rd Party API that randomly fails 50% of the time.
* If it fails, increment a `retry_count` column.
* If it succeeds, set `status = 'SENT'`.
* If `retry_count > 3`, set `status = 'FAILED'`.
* Watch the worker self-heal the system automatically!

---

## 17. Enterprise Features & Observability
* **Deliverability Webhooks**: SendGrid will send an HTTP POST Webhook back to your Go server when the user actually *opens* the email or clicks a link. You ingest these webhooks to provide Analytics to the marketing team.

---

## 18. Source Code Reading
Walkthrough of `github.com/hibiken/asynq`.
* **Go Task Queues**: Study `asynq`. It is an incredible Go library that wraps Redis to provide robust queues, worker pools, retries, and scheduled tasks (e.g., "Send this email in exactly 3 hours") with minimal code!

---

## 19. Architecture
* **Data Privacy**: The Notification Service is highly sensitive. It touches emails and phone numbers. Ensure all payloads in the Message Queue are encrypted at rest, and Logs are strictly scrubbed of PII.

---

## 20. Summary & Cheat Sheet
* **Core Rule**: Async processing via Message Queues.
* **Workers**: Throttled Worker Pools via Go channels.
* **Resiliency**: Circuit Breakers and DLQs for 3rd-party outages.
* **Idempotency**: Block duplicates using Redis SetNX.
