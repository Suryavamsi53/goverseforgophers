# Message Queues: Event-Driven Architecture with Kafka & RabbitMQ

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* Synchronous vs Asynchronous Communication
* Message Queues vs Event Streaming
* RabbitMQ (Message Broker)
* Apache Kafka (Event Streaming)
* Code Examples & Good Principles
* Architecture Diagram
* Real-World Analogy
* Interview Questions
* Quiz
* Exercises
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

In a monolithic architecture, when a user registers, the system might directly call a function to send a welcome email. If the email server is slow, the user is forced to stare at a loading spinner. In a microservices environment, direct synchronous calls between services create a fragile web of dependencies. 

Message Queues and Event Streams decouple services, allowing them to communicate asynchronously. This ensures that a failure in one service doesn't cascade and take down the entire platform.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand the difference between Synchronous and Asynchronous communication.
* Explain the Pub/Sub (Publish/Subscribe) pattern.
* Differentiate between traditional Message Queues (RabbitMQ) and Event Streaming Platforms (Kafka).
* Identify when to use message queues for load leveling and fault tolerance.
* Write basic producer and consumer code in Go.

---

# Prerequisites

Before reading this chapter you should know:

* Microservice Basics & Network Protocols (`03-Network-Protocols.md`)
* Vertical vs. Horizontal Scaling (`02-Scalability.md`)

---

# Why This Topic Exists

When designing systems like Uber or Amazon, many actions don't need to happen instantly. When you order a book, the inventory service, billing service, and shipping service don't need to coordinate in real-time. The order service simply announces "An order was placed," and the other services react. Mastering message queues allows you to build highly resilient, loosely coupled systems that can handle massive traffic spikes (load leveling).

---

# Synchronous vs Asynchronous Communication

### Synchronous (REST / gRPC)
* Service A sends a request to Service B and **blocks** (waits) until it receives a response.
* **Pros**: Simple to reason about, immediate feedback.
* **Cons**: Tight coupling. If Service B is down, Service A fails. If Service B is slow, Service A becomes slow.

### Asynchronous (Message Queues)
* Service A sends a message to a Queue/Broker and immediately returns success to the user. Service B reads the message from the Queue at its own pace.
* **Pros**: Loose coupling. Service B can be down for maintenance, and Service A can still accept orders (they just queue up). Handles sudden traffic spikes gracefully.
* **Cons**: Complex debugging, eventual consistency (the user might have to wait a few seconds to see the result).

---

# Message Queues vs Event Streaming

While often used interchangeably, there is a fundamental architectural difference between traditional message queues and event streaming.

### Traditional Message Queues (RabbitMQ, SQS)
* **Smart Broker, Dumb Consumer**: The broker keeps track of which consumer read which message.
* **Ephemeral**: Once a message is successfully processed and acknowledged (ACKed) by the consumer, it is **deleted** from the queue.
* **Use Case**: Task queues (e.g., sending emails, processing a video, generating a PDF).

### Event Streaming (Apache Kafka, AWS Kinesis)
* **Dumb Broker, Smart Consumer**: The broker acts like a massive distributed append-only log. The consumer is responsible for keeping track of its own "offset" (which message it read last).
* **Persistent**: Messages are **not deleted** after being read. They stay in the log for a configured retention period (e.g., 7 days, or even forever). Multiple different consumer groups can read the exact same message independently.
* **Use Case**: Event sourcing, real-time analytics, auditing (e.g., calculating surge pricing based on a stream of GPS coordinates).

---

# RabbitMQ (Message Broker)

* **Protocol**: AMQP (Advanced Message Queuing Protocol).
* **Architecture**: Producers send messages to an **Exchange**. The Exchange routes the messages to specific **Queues** based on routing keys (e.g., topic, direct, fanout). Consumers read from the Queues.
* **Strengths**: Extremely flexible routing logic. Built-in priority queues and delayed delivery.

---

# Apache Kafka (Event Streaming)

* **Architecture**: Producers publish events to **Topics**. Topics are partitioned across multiple servers (Brokers) for massive horizontal scalability. Consumers read from **Partitions**.
* **Performance**: Kafka writes directly to the filesystem cache and uses zero-copy network transmission. It can handle millions of messages per second.
* **Strengths**: High throughput, replayability (you can rewind the offset and re-process past events), and strict ordering (within a single partition).

---

# Code Examples & Good Principles

### 1. RabbitMQ Producer in Go (Good Principle)

```go
package main

import (
	"log"
	"github.com/streadway/amqp"
)

func main() {
	// 1. Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 2. Open a Channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	// 3. Declare a Queue (idempotent)
	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable (survives broker restart)
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	// 4. Publish a message
	body := "Generate PDF Report"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // Principle: Ensure messages survive crashes
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	
	log.Printf(" [x] Sent %s", body)
}
```

### 2. Kafka Consumer in Go (Good Principle)

```go
package main

import (
	"context"
	"fmt"
	"log"
	"github.com/segmentio/kafka-go"
)

func main() {
	// Principle: Use Consumer Groups to allow horizontal scaling of consumers
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		GroupID:   "analytics-group",
		Topic:     "user-clicks",
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	defer r.Close()

	fmt.Println("Listening for events...")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("could not read message: %v", err)
		}
		
		fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		// Note: In Kafka, reading advances the offset for the GroupID automatically.
	}
}
```

---

# Architecture Diagram

```mermaid
flowchart TD
    subgraph Microservices
        A[Order Service]
    end

    subgraph Message Broker / Event Stream
        B[Kafka Topic: orders]
    end

    subgraph Consumers
        C[Inventory Service]
        D[Billing Service]
        E[Notification Service]
    end

    A -- "Publish (Event: OrderCreated)" --> B
    B -- "Subscribe" --> C
    B -- "Subscribe" --> D
    B -- "Subscribe" --> E
```

---

# Real-World Analogy

* **Synchronous (REST)**: Calling someone on the phone. You have to wait for them to pick up. If they are busy, you fail to communicate.
* **Message Queue (RabbitMQ)**: Sending a task to a team of workers. The boss puts a task on a clipboard. The next available worker takes the paper, does the job, and throws the paper away.
* **Event Streaming (Kafka)**: A corporate bulletin board. The boss pins an announcement to the board. The Marketing team, Sales team, and HR team can all read the same announcement whenever they want. The announcement stays on the board for a month.

---

# Interview Questions

## Beginner
**Q**: What is the Pub/Sub pattern?
*Answer*: Publish/Subscribe is a messaging pattern where publishers don't send messages directly to specific subscribers. Instead, they publish messages to a topic or exchange, and any number of subscribers can independently read those messages.

## Intermediate
**Q**: Why would you choose Kafka over RabbitMQ?
*Answer*: I would choose Kafka if I need event replayability, massive throughput (millions of events/sec), or if multiple different microservices need to react to the exact same event stream (Event Sourcing). I would choose RabbitMQ for complex routing or simple task queues where messages should be deleted after processing.

## Advanced
**Q**: In Kafka, how is message ordering guaranteed, and how does it affect scalability?
*Answer*: Kafka only guarantees strict ordering **within a single partition**. If you need all events for a specific UserID to be processed in order, you must use the UserID as the message Key. Kafka will hash the Key and send all messages for that UserID to the same partition. However, this means all messages for that user will be processed by a single consumer thread, creating a bottleneck if that user is highly active.

---

# Quiz

## Multiple Choice Questions
**1. Which of the following automatically deletes a message once a consumer acknowledges it?**
A) Apache Kafka
B) RabbitMQ
C) AWS Kinesis
*Answer*: B

## True or False
**Synchronous communication is generally more resilient to sudden traffic spikes than asynchronous communication.**
*Answer*: False. Synchronous systems get overwhelmed and crash during spikes. Asynchronous queues act as a buffer, allowing the system to process the spike at a steady, manageable pace.

---

# Exercises

## Beginner
Modify the RabbitMQ Producer Go code to publish a JSON payload instead of a plain text string.

## Intermediate
Create a Kafka Producer in Go that publishes 100 messages in a loop. Then start two instances of the Kafka Consumer code (using the same `GroupID`). Observe how Kafka balances the load between the two consumers.

---

# Summary

Message Queues and Event Streams are the nervous system of a distributed architecture. They provide the ultimate decoupling, allowing services to scale, fail, and recover independently. Use RabbitMQ for worker queues and complex routing, and leverage Kafka for high-throughput event sourcing and analytics.

---

# Key Takeaways

* ✔ Asynchronous communication prevents cascading failures and provides load buffering.
* ✔ RabbitMQ is a message broker that deletes messages after they are processed.
* ✔ Kafka is an event streaming platform that persists messages for replayability and multiple independent readers.
* ✔ Kafka guarantees ordering only at the partition level.

---

# Further Reading
* [RabbitMQ Official Tutorials (Go)](https://www.rabbitmq.com/tutorials/tutorial-one-go.html)
* [Kafka: The Definitive Guide](https://www.oreilly.com/library/view/kafka-the-definitive/9781491936153/)

---

# Next Chapter
➡️ **Next:** `09-Storage.md`
