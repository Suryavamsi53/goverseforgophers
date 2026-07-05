# Observer Pattern

The Observer Pattern is a behavioral design pattern that defines a one-to-many dependency between objects. When the subject changes state, all its registered dependents (Observers) are notified and updated automatically.

This is the foundational pattern behind Event-Driven Architecture (Pub/Sub) and Reactive UI frameworks (like React or Vue).

## 1. The Synchronous Implementation

In a classic Go implementation, the `Subject` maintains a slice of `Observer` interfaces. When something happens, it iterates through the slice and calls `.Update()` on each one.

```go
// 1. The Observer Interface
type Observer interface {
    Update(eventID string)
}

// 2. The Subject
type EventBroker struct {
    observers []Observer
    mu        sync.RWMutex // Protects the slice from data races!
}

// 3. Subscription Management
func (b *EventBroker) Subscribe(o Observer) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.observers = append(b.observers, o)
}

// 4. The Trigger
func (b *EventBroker) Publish(eventID string) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    
    // Notify everyone!
    for _, obs := range b.observers {
        obs.Update(eventID)
    }
}
```

## 2. The Danger of Synchronous Observers

The code above has a fatal flaw. 
If `Observer 1` is a logging struct that takes 5 seconds to write to a slow hard drive, the `Publish` function will **block** for 5 seconds. `Observer 2` and `Observer 3` will have to wait in line. The entire application slows to a crawl!

## 3. The Asynchronous Go Solution (Channels)

Because Go has Channels and Goroutines, we can implement a blazingly fast, non-blocking Observer pattern.

Instead of passing an interface with an `.Update()` method, the Observers simply provide a Go `chan`.

```go
type AsyncBroker struct {
    subscribers []chan string
    mu          sync.RWMutex
}

func (b *AsyncBroker) Subscribe(ch chan string) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.subscribers = append(b.subscribers, ch)
}

func (b *AsyncBroker) Publish(eventID string) {
    b.mu.RLock()
    defer b.mu.RUnlock()
    
    for _, ch := range b.subscribers {
        // We use a non-blocking select!
        // If the Observer's channel is full/slow, we drop the message 
        // to protect the broker from freezing!
        select {
        case ch <- eventID:
        default:
            fmt.Println("Warning: Observer is too slow, dropping event!")
        }
    }
}
```

Now, the `Publish` method executes in nanoseconds, regardless of how slow the Observers are!

## 4. Production Pub/Sub

While building an in-memory Observer using Channels is great for a single server, it doesn't work if you have 10 servers in a Kubernetes cluster (an event published on Server 1 will never reach the Observers on Server 2).

In a distributed system, you move the Observer Pattern out of Go code entirely, and use a dedicated Message Broker (like **Apache Kafka**, **RabbitMQ**, or **Redis Pub/Sub**). The Go application simply becomes an API client that pushes and pulls from the external broker.
