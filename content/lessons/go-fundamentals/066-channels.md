# Channels

A WaitGroup allows you to *wait* for goroutines to finish, but what if you need to actually pass data between them while they are running?

> **"Do not communicate by sharing memory; instead, share memory by communicating."** — Effective Go

In languages like C++ or Java, threads share data by accessing the exact same variables in memory, requiring complex, error-prone Mutex locks to prevent race conditions. 

Go provides **Channels**: typed conduits that safely pipe data from one goroutine to another.

## 1. Syntax

You create a channel using `make(chan Type)`. 

The `<-` operator is used to send and receive data. The arrow indicates the direction of data flow.

```go
func main() {
    // 1. Create a channel of integers
    ch := make(chan int)

    // 2. Spawn a background worker
    go func() {
        // SEND data INTO the channel
        ch <- 42 
    }()

    // 3. Block and wait to RECEIVE data FROM the channel
    val := <-ch 
    
    fmt.Println("Received:", val)
}
```

## 2. Under the Hood: The `hchan` Struct

Channels feel like magic, but they are incredibly complex data structures managed by the Go runtime. 

When you call `make(chan int)`, Go allocates an `hchan` struct on the Heap.

```mermaid
graph TD
    subgraph hchan [hchan Struct (Heap)]
        M[sync.Mutex]
        Q[Circular Queue / Ring Buffer]
        S[SendQ: WaitList of Blocked Senders]
        R[RecvQ: WaitList of Blocked Receivers]
    end
    
    G1((Goroutine A)) -->|Attempts to Send| M
    G2((Goroutine B)) -->|Attempts to Receive| M
```

**Why are channels thread-safe?**
Because the `hchan` struct contains a hidden `sync.Mutex` lock! Every time a goroutine attempts to read or write to a channel, the Go runtime locks the channel, performs the data transfer, and unlocks it. 

Furthermore, if Goroutine A tries to read from an empty channel, the Go Scheduler puts Goroutine A to sleep and places it into the `RecvQ` waiting list. Once Goroutine B writes data to the channel, the Scheduler wakes Goroutine A back up!
