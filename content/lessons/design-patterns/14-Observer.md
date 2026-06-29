# Observer Pattern

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Real-World Analogy
* Core Concepts
* Architecture Diagram
* Step-by-Step Implementation
* Syntax
* Beginner Example
* Intermediate Example (Using Channels)
* Advanced Example
* Production Use Cases
* Performance Analysis
* Best Practices
* Common Mistakes
* Debugging Guide
* Exercises
* Quiz
* Interview Questions
* Cheat Sheet
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

The **Observer Pattern** is a Behavioral Design Pattern that defines a one-to-many dependency between objects so that when one object changes state, all its dependents are notified and updated automatically.

In the wider industry, this is often called the **Pub/Sub (Publish/Subscribe)** or **Event Listener** pattern. In Go, while you can implement this using classic Interfaces and Structs, it is incredibly common to implement the Observer pattern using Go's native concurrency primitives: **Channels** and **Goroutines**.

---

# Learning Objectives

After completing this chapter you will be able to:

* Implement a classic Observer pattern using Interfaces and Mutexes.
* Implement a highly concurrent Observer pattern using Channels.
* Safely manage the registration and deregistration of dynamic subscribers.

---

# Prerequisites

Before reading this chapter you should know:

* Structs and Interfaces.
* Concurrency, Goroutines, and Channels (`08-Goroutines.md`).
* Mutexes (`29-Mutexes.md`).

---

# Why This Topic Exists

Imagine building an E-Commerce system. When a user places an order, three things must happen:
1. An Email receipt must be sent.
2. The Inventory system must deduct the items.
3. The Analytics system must record the sale.

If you hardcode these three function calls inside the `PlaceOrder()` method, your core checkout logic becomes tightly coupled to Email, Inventory, and Analytics packages. If the Email server is down, the entire checkout process might fail!

The Observer pattern decouples them. The `PlaceOrder()` function simply yells "AN ORDER HAPPENED!" (Publishes an Event). The Email, Inventory, and Analytics systems are registered as "Subscribers". They hear the event and process it independently.

---

# Real-World Analogy

### The YouTube Channel

* **The Subject (Publisher)**: A YouTube Channel.
* **The Observers (Subscribers)**: Millions of users who clicked the "Subscribe" button.
* **The Event**: The YouTuber uploads a new video.
* **The Notification**: The YouTube server loops through the list of subscribers and sends a push notification to every single one of them. The YouTuber doesn't know who the subscribers are; they just know they clicked "Publish".

---

# Core Concepts

* **Subject (Publisher)**: The object that holds the state. It maintains a list of Observers and provides methods to `Subscribe` and `Unsubscribe`.
* **Observer (Subscriber)**: An interface with an `Update()` method (or a Channel) that reacts when the Subject sends a notification.
* **Event**: The data or state change that triggers the notification.

---

# Architecture Diagram

```mermaid
flowchart TD
    Subject[Publisher<br/>Order System]
    
    Sub1[Subscriber 1<br/>Email Service]
    Sub2[Subscriber 2<br/>Inventory Service]
    Sub3[Subscriber 3<br/>Analytics Service]

    Subject -- "1. Registers" <-- Sub1
    Subject -- "2. Registers" <-- Sub2
    Subject -- "3. Registers" <-- Sub3

    Subject -- "Event: Order Placed!<br/>NotifyAll()" --> Sub1
    Subject -- "NotifyAll()" --> Sub2
    Subject -- "NotifyAll()" --> Sub3
```

---

# Step-by-Step Implementation (Classic Interface Approach)

1. Define the **Observer Interface** (e.g., `type Observer interface { Update(data string) }`).
2. Create concrete structs that implement the Observer interface.
3. Create the **Subject Struct**. Give it a slice to hold the registered observers: `observers []Observer`.
4. Create `Subscribe(o Observer)` and `Unsubscribe(o Observer)` methods on the Subject.
5. Create a `NotifyAll(data string)` method that loops through the slice and calls `o.Update(data)` on every observer.

---

# Syntax (Classic)

```go
type Observer interface { Update(msg string) }

type Subject struct {
    observers []Observer
}

func (s *Subject) Subscribe(o Observer) {
    s.observers = append(s.observers, o)
}

func (s *Subject) Notify(msg string) {
    for _, o := range s.observers {
        o.Update(msg)
    }
}
```

---

# Beginner Example

The classic Interface-based Observer. We must use a `sync.Mutex` because in a real application, subscribers might be added or removed from different goroutines while a notification is happening!

```go
package main

import (
	"fmt"
	"sync"
)

// 1. The Observer Interface
type Observer interface {
	OnNotify(event string)
}

// 2. Concrete Observers
type EmailService struct{}
func (e *EmailService) OnNotify(event string) {
	fmt.Println("EmailService received:", event)
}

type AnalyticsService struct{}
func (a *AnalyticsService) OnNotify(event string) {
	fmt.Println("AnalyticsService logging:", event)
}

// 3. The Subject (Publisher)
type OrderSystem struct {
	mu        sync.RWMutex
	observers []Observer
}

func (o *OrderSystem) Subscribe(obs Observer) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.observers = append(o.observers, obs)
}

func (o *OrderSystem) NotifyAll(event string) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	for _, obs := range o.observers {
		obs.OnNotify(event)
	}
}

// Core Business Logic
func (o *OrderSystem) PlaceOrder(id int) {
	fmt.Printf("\n--- Placing Order #%d ---\n", id)
	// Do DB work...
	
	// Notify everyone asynchronously? No, this is synchronous right now.
	o.NotifyAll(fmt.Sprintf("Order #%d was completed!", id))
}

func main() {
	system := &OrderSystem{}

	// Register subscribers
	system.Subscribe(&EmailService{})
	system.Subscribe(&AnalyticsService{})

	// Trigger the event
	system.PlaceOrder(101)
	system.PlaceOrder(102)
}
```

---

# Intermediate Example (Using Channels)

In Go, the classic interface approach is synchronous (the `PlaceOrder` function halts while the EmailService runs). This is bad. Instead, we use Go Channels to create an asynchronous Pub/Sub system!

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// The Subject
type EventBroker struct {
	mu          sync.RWMutex
	// The list of subscribers is now a slice of channels!
	subscribers []chan string 
}

// Subscribe returns a channel that the subscriber can listen to
func (b *EventBroker) Subscribe() <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	// Create a buffered channel so the publisher doesn't block
	ch := make(chan string, 10) 
	b.subscribers = append(b.subscribers, ch)
	return ch
}

func (b *EventBroker) Publish(event string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	for _, ch := range b.subscribers {
		// Non-blocking send!
		select {
		case ch <- event:
		default:
			fmt.Println("Warning: Subscriber channel full, dropping event!")
		}
	}
}

func main() {
	broker := &EventBroker{}

	// Subscriber 1 (Runs in background)
	emailChan := broker.Subscribe()
	go func() {
		for event := range emailChan {
			fmt.Println("[Email Worker] processing:", event)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Subscriber 2 (Runs in background)
	analyticsChan := broker.Subscribe()
	go func() {
		for event := range analyticsChan {
			fmt.Println("[Analytics Worker] processing:", event)
		}
	}()

	// Publisher
	fmt.Println("Publishing event...")
	broker.Publish("User Signed Up!")
	
	time.Sleep(1 * time.Second) // Wait for async workers to finish
}
```

---

# Production Use Cases

### 1. WebSockets and Chat Servers
When a user sends a chat message in a Discord channel, they don't send it to the other 50 users directly. They send it to the Server (Subject). The Server looks up the slice of WebSocket connections (Observers) registered to that chat room, and ranges over them, writing the message to every connection.

### 2. Event-Driven Microservices
At a massive scale, the Observer pattern moves out of code and into infrastructure using Message Brokers like **Kafka**, **RabbitMQ**, or **AWS SNS**. The core concept is identical: Microservice A publishes a message to a Kafka Topic. Microservices B and C are subscribed to that Topic and react accordingly.

---

# Performance Analysis

* **Synchronous Observer**: Calling `obs.OnNotify()` blocks the publisher until the subscriber finishes. If an observer does a slow DB query, the publisher hangs.
* **Channel Observer**: Using buffered channels allows the publisher to instantly fire off the events and return. However, if the buffer fills up, you must decide whether to block (which risks a system deadlock) or drop the event.

---

# Best Practices

* **Protect the Slice**: If subscribers can be added dynamically at runtime, you MUST protect the `[]Observer` slice with a `sync.RWMutex`. Slices are not thread-safe.
* **Use Asynchronous Handlers**: Observers should never block the Publisher. Either use channels, or have the Publisher wrap the `o.Update()` call in a Goroutine: `go o.Update(msg)`.
* **Handle Unsubscribing gracefully**: Removing an item from a slice in Go is slightly annoying. If you have many dynamic subscribers (like WebSockets disconnecting), use a `map[Observer]struct{}` instead of a slice for `O(1)` deletion.

---

# Common Mistakes

### Memory Leaks (The Lapsed Listener Problem)
If an Observer object goes out of scope and the client no longer uses it, the Garbage Collector *cannot* delete it! Why? Because the Subject's internal slice still holds a pointer to it. You must explicitly `Unsubscribe()` observers before discarding them, or they will live in memory forever.

---

# Debugging Guide

* **"Publisher is hanging forever"**: If you used unbuffered channels for your subscribers, and one subscriber crashes or stops reading from its channel, the publisher will block forever trying to send to it. Always use buffered channels and a `select` statement with a `default` case to drop/log stalled subscribers.

---

# Exercises

## Beginner
Implement a classic Interface-based Observer. Subject: `WeatherStation`. Observers: `PhoneDisplay` and `WindowDisplay`. When the `WeatherStation` calls `SetTemperature(t int)`, it should notify both displays to print the new temperature.

## Intermediate
Refactor the Beginner exercise. Use a `map[Observer]bool` instead of a slice in the Subject so you can easily implement a `RemoveObserver(o Observer)` method. 

---

# Quiz

## Multiple Choice Questions
**1. What is the biggest danger of a synchronous Observer pattern?**
A) It uses too much memory.
B) If one Observer is very slow (e.g., waiting on a network call), it forces the Publisher and all other Observers to wait, freezing the system.
C) The interfaces are difficult to implement.
*Answer*: B

## True or False
**In a Channel-based Observer pattern in Go, the publisher should ideally block if a subscriber's channel is full to ensure no data is lost.**
*Answer*: False! If a publisher blocks, it risks freezing the entire core application. It is much safer to use a `select/default` block to drop the event (or log an error) rather than halting the core Publisher because a single subscriber is slow.

---

# Interview Questions

## Beginner
**Q**: What is the purpose of the Observer Pattern?
*Answer*: It defines a one-to-many relationship where a Publisher broadcasts events to multiple decoupled Subscribers, removing the need for the Publisher to know exactly who depends on its data.

## Intermediate
**Q**: Explain the "Lapsed Listener" memory leak.
*Answer*: It occurs when an Observer registers with a long-lived Subject, but the application finishes using the Observer and forgets to unregister it. Because the Subject maintains a strong pointer reference to the Observer in its internal list, the Garbage Collector cannot reclaim the Observer's memory, leading to a permanent memory leak.

## Advanced
**Q**: How would you build a highly concurrent Pub/Sub system in Go that guarantees the Publisher never blocks?
*Answer*: I would use channels. The Subject maintains a list of buffered channels. When publishing, it ranges over the channels and uses a non-blocking `select` statement. If a channel's buffer is full, the `default` case is executed, allowing the Publisher to instantly log a warning and move on to the next subscriber without ever blocking.

---

# Cheat Sheet

* **Synchronous Publisher**: `go func(){ obs.Notify() }() `(Wrap in goroutine to prevent blocking)
* **Non-Blocking Channel Send**:
```go
select {
case subChan <- event:
default:
    // Drop event, subscriber is too slow
}
```
* **Map for `O(1)` Deletion**:
```go
type Subject struct {
    obs map[Observer]struct{}
}
func (s *Subject) Remove(o Observer) { delete(s.obs, o) }
```

---

# Summary

The Observer pattern is the backbone of event-driven architecture. Whether implemented via classic interfaces for simple GUI updates, or via Go channels for high-throughput backend systems, it allows you to build incredibly resilient and decoupled systems where components react to the world without being hard-wired to it.

---

# Key Takeaways

* ✔ Observer decouples Publishers from Subscribers.
* ✔ Classic implementations use Interfaces.
* ✔ Modern Go implementations use Channels for async processing.
* ✔ Always protect the subscriber list with a Mutex.
* ✔ Beware of memory leaks caused by forgetting to Unsubscribe.

---

# Further Reading
* [Refactoring.guru: Observer Pattern](https://refactoring.guru/design-patterns/observer)

---

# Next Chapter
➡️ **Next:** `15-Command.md`
