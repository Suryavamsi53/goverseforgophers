# Fan-Out Pattern

**Fan-Out** is the exact opposite of Fan-In. 

You have a single channel producing a massive stream of data (e.g., a websocket feeding 10,000 live stock prices per second). Processing this stream with a single Goroutine is too slow. 

The Fan-Out pattern takes the single input channel and distributes its data across multiple worker Goroutines.

*(Note: If this sounds exactly like the Worker Pool pattern we already covered, you are correct! The Worker Pool is the most common implementation of Fan-Out).*

## 1. Fan-Out vs Broadcast

It is critical to understand that Fan-Out is a **Load Balancing** pattern. 

If you have 1 input channel and 5 workers reading from it using `<-ch`, the Go Scheduler will randomly assign each piece of data to exactly **one** worker. 
* Message 1 goes to Worker A.
* Message 2 goes to Worker B.

This is NOT a Broadcast. Worker A will never see Message 2. 

## 2. The Broadcast Pattern (Fan-Out to All)

If you actually want a **Broadcast** (where every single worker receives a copy of every single message), you cannot just use multiple readers on a single channel.

Instead, the orchestrator must maintain a list of specific channels for each worker, and iterate through them.

```go
func Broadcast(input <-chan int, workers []chan<- int) {
    for msg := range input {
        // Send a copy to every single worker
        for _, workerCh := range workers {
            // DANGER: What if workerCh is full/blocked?
            workerCh <- msg 
        }
    }
}
```

## 3. The Slow Consumer Problem

The Broadcast implementation above has a fatal flaw: **The Slow Consumer Problem**.

If you are broadcasting to 5 workers, and Worker 3 is performing a heavy database insert and hasn't read from its channel yet, the `workerCh <- msg` line will **block**. 
The Orchestrator will completely freeze. It will stop reading from the `input` channel, and it will stop sending messages to Workers 1, 2, 4, and 5. 

One slow worker brings down the entire system!

### The Solution: Non-Blocking Sends
To build a resilient Fan-Out/Broadcast system, you must use the `select` statement with a `default` case to drop messages for slow consumers (or route them to a dead-letter queue).

```go
for _, workerCh := range workers {
    select {
    case workerCh <- msg:
        // Delivered successfully
    default:
        // Worker is too slow! Buffer is full!
        // We drop the message (or log it) and instantly move on 
        // to keep the system flowing!
        log.Println("Dropped message for slow worker")
    }
}
```
