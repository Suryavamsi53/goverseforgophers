# RabbitMQ & Message Queues

## 1. Learning Objectives
* **What you'll learn**: The mechanics of AMQP (Advanced Message Queuing Protocol), Exchanges, Routing Keys, and how RabbitMQ differs fundamentally from Kafka.
* **Why it matters**: RabbitMQ excels at complex routing, task queues, and low-latency message delivery, ensuring unreliable background jobs (like video processing or sending emails) never get lost.
* **Where it's used**: Celery workers, background job processing, legacy enterprise integration, and targeted event routing.

---

## 2. Real-world Story
Imagine a massive Post Office. 
Kafka is like a public bulletin board—you pin a message, and anyone can come read it forever.
RabbitMQ is like a smart mail-sorting facility. You hand a letter (Message) to the Postmaster (Exchange). The Postmaster looks at the zip code (Routing Key), and physically routes the letter into exactly one specific Mailbox (Queue). Once the recipient reads the letter, it is taken out of the mailbox and destroyed forever. 

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Go Producer] -->|Message + Routing Key 'pdf.create'| B(Exchange)
    
    B -->|Direct Match| C[Queue: pdf_tasks]
    B -->|Wildcard Match| D[Queue: all_tasks]
    
    C -->|Pushed to| E[Go PDF Worker 1]
    C -->|Pushed to| F[Go PDF Worker 2]
    
    style B fill:#8b5cf6,color:#fff
```

---

## 4. Internal Working (Under the Hood)
RabbitMQ uses a "Smart Broker, Dumb Consumer" model.
1. **Producer**: Sends messages to an Exchange.
2. **Exchange**: Evaluates rules (Bindings) and pushes the message into Queues.
3. **Queue**: A physical buffer in RAM/Disk holding messages until consumed.
4. **Consumer**: Connects to a Queue. RabbitMQ *pushes* the message to the Go consumer. 
5. **ACK (Acknowledgment)**: The Go consumer must explicitly tell RabbitMQ "I finished processing this". Only then will RabbitMQ delete the message.

---

## 5. Compiler Behavior
* **Goroutine Multiplexing**: The official Go client (`github.com/rabbitmq/amqp091-go`) heavily utilizes channels. A single TCP connection can be multiplexed into multiple lightweight "Channels" (an AMQP concept, not a Go channel), allowing hundreds of concurrent Goroutines to publish messages without paying the TCP handshake penalty.

---

## 6. Memory Management
* **Queue Memory Limits**: Unlike Kafka which safely stores terabytes of data on disk, RabbitMQ primarily relies on RAM for extreme low-latency. If your Go consumers crash and a queue fills up with 50 million unread messages, RabbitMQ will trigger a memory alarm and completely block all Producers from sending new data!

---

## 7. Code Examples

### 🔹 Example 1: Simple (Publishing)
```go
import amqp "github.com/rabbitmq/amqp091-go"

func PublishTask() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
    defer conn.Close()
    ch, _ := conn.Channel() // The logical session
    defer ch.Close()

    // Declare the Queue to ensure it exists
    q, _ := ch.QueueDeclare("pdf_tasks", true, false, false, false, nil)

    // Publish to the default Exchange, routing directly to the Queue
    ch.PublishWithContext(context.Background(),
        "",     // exchange
        q.Name, // routing key (matches queue name)
        false, false,
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte("Generate Invoice #123"),
        })
}
```

### 🔹 Example 2: Intermediate (Consuming)
```go
func ConsumeTasks() {
    // ... setup conn and ch ...
    msgs, _ := ch.Consume(
        "pdf_tasks", // queue
        "",          // consumer name
        false,       // auto-ack (FALSE means manual ack is required!)
        false, false, false, nil,
    )

    // msgs is a native Go channel!
    for d := range msgs {
        log.Printf("Received a task: %s", d.Body)
        
        err := GeneratePDF(d.Body)
        if err == nil {
            d.Ack(false) // Tell RabbitMQ to delete the message!
        } else {
            d.Nack(false, true) // Re-queue the message to try again later!
        }
    }
}
```

### 🔹 Example 3: Advanced (Fanout Exchange)
```go
// Broadcasting one message to multiple queues!
// 1. Create a Fanout Exchange
ch.ExchangeDeclare("logs_exchange", "fanout", true, false, false, false, nil)

// 2. Publish to the Exchange (Routing key is ignored)
ch.PublishWithContext(ctx, "logs_exchange", "", false, false, amqp.Publishing{Body: []byte("Server crashed")})
```

### 🔹 Example 4: Production
```go
// Prefetch Count (QoS)
// NEVER let RabbitMQ push 100,000 messages into your Go consumer's RAM at once.
// Set QoS to 10: RabbitMQ will only give you 10 messages. It won't give you 
// an 11th until you Ack one of the first 10.
ch.Qos(10, 0, false)
```

### 🔹 Example 5: Interview
```go
// Q: What happens if your Go consumer crashes halfway through processing a message?
// A: Because auto-ack is false, the TCP connection drops before an Ack is sent. 
// RabbitMQ instantly notices the dropped connection and re-queues the message for another worker. No data loss!
```

---

## 8. Production Examples
1. **Work Queues**: Distributing heavy tasks (like image resizing or sending batch emails) across 50 separate Go worker containers.
2. **Complex Routing (Topic Exchanges)**: Publishing an event with routing key `error.billing.db`. RabbitMQ can smartly route this to the `billing_errors` queue AND the `global_db_errors` queue based on wildcard bindings (`*.billing.*`).

---

## 9. Performance & Benchmarking
* **Kafka vs RabbitMQ**: Kafka is for streaming massive historical datasets (Millions of msg/sec). RabbitMQ is for complex routing and task queues with extreme low-latency (50K msg/sec). RabbitMQ messages are deleted once read; Kafka messages persist.

---

## 10. Best Practices
* ✅ **Do**: Use persistent queues and durable messages (`DeliveryMode: amqp.Persistent`) if you cannot afford to lose tasks during a RabbitMQ server reboot.
* ❌ **Don't**: Use RabbitMQ to transfer 50MB video files. Put the video in an S3 Bucket, and put the S3 URL in the RabbitMQ message.
* 🏢 **Google / Uber / Netflix Style**: Decouple the HTTP API from the Workers. The HTTP API only publishes to RabbitMQ (takes 2ms) and returns `202 Accepted`. Background Go Workers pull from RabbitMQ and do the 5-second heavy lifting.

---

## 11. Common Mistakes
1. **The Infinite Nack Loop**: If a message causes your Go code to panic or error predictably, and you return `d.Nack(requeue=true)`, RabbitMQ will instantly hand it back to you. You will process it, fail, and Nack it 1,000 times a second, burning 100% CPU. (Send it to a Dead Letter Exchange instead!).
2. **Connection Churn**: Opening a new TCP `amqp.Dial` connection for every single published message. Keep a long-lived connection pool in your Go server!

---

## 12. Debugging
How to troubleshoot RabbitMQ in production:
* **The Management UI**: RabbitMQ ships with a brilliant web dashboard (port 15672). You can visually see the exact number of messages waiting in every queue, the publish rates, and which Go consumers are currently connected.

---

## 13. Exercises
1. **Easy**: Write a Publisher that sends a "Hello World" task.
2. **Medium**: Write two Go Consumers that listen to the same queue. Watch how RabbitMQ round-robins the tasks between them automatically.
3. **Hard**: Implement a Dead Letter Exchange (DLX). If a task fails 3 times, route it to a `failed_tasks` queue so it doesn't block the system.
4. **Expert**: Implement an AMQP reconnection loop in Go. If the RabbitMQ server restarts, your Go app should detect the closed channel and automatically re-dial without crashing.

---

## 14. Quiz
1. **MCQ**: What RabbitMQ component determines which queues receive a message?
   * (A) The Producer (B) The Exchange (C) The Consumer. *(Answer: B. The Producer only knows about the Exchange. The Exchange handles all routing).*
2. **Code Review**: `msgs, _ := ch.Consume("q", "", true, ...)` What is the massive danger of `auto-ack = true`? *(RabbitMQ considers the message "processed" the millisecond it sends it over TCP. If your Go app crashes 1 millisecond later, the message is permanently deleted and lost!).*

---

## 15. FAANG Interview Questions
* **Beginner**: Explain the Publisher/Subscriber pattern.
* **Intermediate**: How do you guarantee a message is processed Exactly Once in a distributed task queue? (Hint: You can't. You must use At-Least-Once delivery combined with Idempotent consumers).
* **Senior (Google/Meta)**: Architect a high-throughput priority queue system. Since RabbitMQ strictly processes FIFO, how do you handle "VIP Customer" tasks that must jump the line?

---

## 16. Mini Project
**The Resilient Thumbnail Generator**
* Build an HTTP endpoint `POST /upload` that simulates saving an image and publishes an `image_id` to RabbitMQ.
* Build a separate Go Worker that consumes `image_id`s, sleeps for 2 seconds (simulating processing), and prints "Done".
* Kill the Worker mid-sleep. Verify the task remains in RabbitMQ and is instantly picked up when you restart the Worker.

---

## 17. Enterprise Features & Observability
* **Quorum Queues**: For true High Availability, enterprise clusters use Quorum Queues, which replicate the queue data across 3 different physical servers using the Raft consensus algorithm.

---

## 18. Source Code Reading
Walkthrough of `github.com/rabbitmq/amqp091-go`.
* **Asynchronous Networking**: Study how the Go client maps AMQP frames to Go channels. When the server pushes a message over the TCP socket, a background Goroutine decodes the binary frame and pushes it directly into your `range msgs` loop.

---

## 19. Architecture
* **RPC over RabbitMQ**: You can actually use RabbitMQ for Request/Response (RPC). The Go Publisher includes a `ReplyTo` queue name and a `CorrelationID`. The Go Worker processes the task and publishes the result back to that specific `ReplyTo` queue!

---

## 20. Summary & Cheat Sheet
* **Exchange**: Sorts the mail (Direct, Fanout, Topic).
* **Queue**: Stores the mail.
* **Routing Key**: The address label.
* **Ack**: Manual confirmation that work is finished safely.
