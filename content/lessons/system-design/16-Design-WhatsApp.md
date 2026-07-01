# System Design Interview: Design WhatsApp (Real-time Chat)

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* System Requirements
* Back-of-the-Envelope Estimation
* High-Level Design
* Deep Dive: Managing Connections (WebSockets)
* Deep Dive: Message Routing & Delivery
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

Designing a real-time chat application like WhatsApp, Facebook Messenger, or Discord requires a fundamental shift in architecture compared to traditional REST APIs. In a standard web app, the client initiates all communication (HTTP Request). But in a chat app, the server must push messages to the client instantly, without the client asking. This requires persistent, stateful connections, presenting unique challenges in routing, scaling, and load balancing.

---

# Learning Objectives

After completing this chapter you will be able to:

* Understand when and why to use WebSockets over HTTP Long Polling.
* Design a highly concurrent connection management layer in Go.
* Explain how messages are routed between different backend servers.
* Design a delivery acknowledgement system (single tick vs. double tick).

---

# Prerequisites

* Network Protocols (`03-Network-Protocols.md`)
* Databases (`07-Databases.md`)
* Caching (`06-Caching.md`)

---

# System Requirements

### Functional Requirements
1. **1-on-1 Chat**: Users can send real-time text messages to each other.
2. **Delivery Status**: Senders should see when a message is Sent, Delivered, and Read.
3. **Online Status**: Users can see if their contacts are online.
4. **Message History**: Users can view their past conversations across multiple devices (Assuming a cloud-synced model like Messenger/Telegram, rather than WhatsApp's strict local-device model, for architectural breadth).

### Non-Functional Requirements
1. **Low Latency**: Messages must be delivered in under 500ms.
2. **High Concurrency**: Millions of concurrent active connections.
3. **High Availability**: The system cannot lose messages.

---

# Back-of-the-Envelope Estimation

* **Daily Active Users (DAU)**: 500 Million.
* **Concurrent Connections**: Assume 10% of DAU are online at the same time -> `50 Million concurrent connections`.
* **Message Throughput**: 40 messages per user per day -> `20 Billion messages/day`.
  * `20 Billion / 86400 = ~230,000 Messages Per Second (Average)`.
* **Storage**: 20 Billion messages * 100 Bytes = `2 TB / day`. 
  * Over 5 years: `~3.6 Petabytes`. We will need a highly scalable NoSQL database (like Cassandra or HBase) optimized for fast writes and time-series reads.

**The Bottleneck**: Managing 50 Million persistent TCP connections. A standard HTTP API server is stateless and closes connections immediately. A chat server must hold millions of connections open indefinitely. Go is uniquely perfectly suited for this due to lightweight goroutines.

---

# High-Level Design

The architecture is split into stateless API servers (for historical data) and stateful Chat Servers (for real-time messaging).

1. **Stateless API Servers**: Handle user registration, profile updates, and fetching old message history via standard HTTP REST.
2. **Chat Servers (WebSockets)**: Maintain persistent TCP connections with online users. They do nothing but route messages.
3. **Session Store (Redis)**: Tracks *which* Chat Server a specific user is currently connected to.
4. **Message Queue / PubSub (Kafka/Redis)**: Routes messages between different Chat Servers.
5. **Database (Cassandra/ScyllaDB)**: Stores the permanent history of the chats.

---

# Deep Dive: Managing Connections (WebSockets)

Traditional HTTP is unidirectional (Client -> Server). To get real-time messages, clients used to use **Long Polling** (keeping an HTTP request hanging until the server has data). This is incredibly inefficient.

**WebSockets (WS)** provide a persistent, bi-directional, full-duplex TCP connection. The client can push to the server, and the server can push to the client, instantly.

When Alice opens WhatsApp:
1. She connects to the Load Balancer.
2. The Load Balancer routes her to `Chat Server A`.
3. An active WebSocket connection is established.
4. `Chat Server A` updates the global Redis Session Store: `[User: Alice, Server: Chat Server A]`.

---

# Deep Dive: Message Routing & Delivery

What happens when Alice sends a message to Bob?

1. **Send**: Alice sends a message payload over her WebSocket to `Chat Server A`.
2. **Acknowledge**: `Chat Server A` assigns an ID, writes the message to the Database, and sends a "Sent" (single tick) ack back to Alice.
3. **Lookup**: `Chat Server A` queries Redis to find out where Bob is connected. 
    * Redis says: `[User: Bob, Server: Chat Server B]`.
4. **Route**: `Chat Server A` publishes the message to a message broker (e.g., Redis Pub/Sub or Kafka) targeting `Chat Server B`.
5. **Deliver**: `Chat Server B` receives the message from the broker and pushes it down Bob's active WebSocket.
6. **Delivered Ack**: Bob's app sends a "Delivered" ack back to `Chat Server B`, which routes it back to Alice (double tick).

*(If Bob is offline, the message remains in the Database. When Bob comes online and connects, his app queries the API for undelivered messages).*

---

# Code Examples & Good Principles

### Principle: A Basic WebSocket Hub in Go

Go is famous for chat servers because it can handle millions of connections by dedicating one lightweight goroutine to each connection, rather than an expensive OS thread.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configure the Upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for this example
	},
}

// Hub maintains the set of active clients on THIS server instance
type Hub struct {
	clients map[*websocket.Conn]string // Maps connection to UserID
}

var serverHub = Hub{
	clients: make(map[*websocket.Conn]string),
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// 1. Upgrade initial GET request to a WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Principle: In a real app, authenticate the user via JWT before upgrading
	userID := r.URL.Query().Get("user")
	
	// Register client to this server's hub
	serverHub.clients[ws] = userID
	fmt.Printf("User %s connected.\n", userID)
	
	// Principle: In a distributed system, you would now update Redis: 
	// Redis.Set("session:userID", "ServerIP_1")

	// 2. Listen for incoming messages in an infinite loop
	for {
		var msg struct {
			To   string `json:"to"`
			Text string `json:"text"`
		}
		
		// Block until a message is received from this client
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("User %s disconnected.\n", userID)
			delete(serverHub.clients, ws)
			break // Break the loop, close the connection
		}

		fmt.Printf("Received message from %s to %s: %s\n", userID, msg.To, msg.Text)
		
		// Principle: Here, we would look up the 'To' user in Redis to find their server.
		// If they are on THIS server, we can send it directly.
		for clientConn, clientID := range serverHub.clients {
			if clientID == msg.To {
				clientConn.WriteJSON(map[string]string{
					"from": userID,
					"text": msg.Text,
				})
			}
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Chat server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

# Architecture Diagram

```mermaid
flowchart TD
    Alice([Alice])
    Bob([Bob])
    
    subgraph Edge
        LB[Load Balancer]
    end
    
    subgraph Chat Servers (Stateful)
        CS_A[Chat Server A]
        CS_B[Chat Server B]
    end
    
    subgraph Infrastructure
        Redis[(Redis Session Store)]
        PubSub[[Redis Pub/Sub]]
        Cassandra[(Cassandra DB)]
    end

    %% Connections
    Alice <--"Persistent WebSocket"--> LB
    LB <--> CS_A
    
    Bob <--"Persistent WebSocket"--> LB
    LB <--> CS_B
    
    %% Session Updates
    CS_A -. "Update Session (Alice)" .-> Redis
    CS_B -. "Update Session (Bob)" .-> Redis
    
    %% Message Flow
    Alice -- "1. Send Msg to Bob" --> CS_A
    CS_A -- "2. Save to DB" --> Cassandra
    CS_A -- "3. Where is Bob?" --> Redis
    CS_A -- "4. Publish to CS_B" --> PubSub
    PubSub -- "5. Deliver" --> CS_B
    CS_B -- "6. Push to Bob" --> Bob
```

---

# Real-World Analogy

* **HTTP REST**: Sending a letter via the postal service. You drop it in a box, and eventually, the recipient gets it. To check for replies, you have to constantly walk to your mailbox (polling).
* **WebSockets**: Making a phone call. You dial a number, establish a persistent connection, and both parties can speak and listen instantly until one hangs up.
* **The Redis Session Store**: The phone company's switchboard directory. If Chat Server A wants to connect to Bob, it has to ask the switchboard which specific telephone wire (Server B) Bob is currently plugged into.

---

# Interview Questions

## Beginner
**Q**: Why is a relational database (like PostgreSQL) generally a bad choice for storing chat messages at scale?
*Answer*: Chat applications generate an astronomical volume of writes (billions per day). Relational databases struggle to scale writes horizontally. A NoSQL database like Cassandra is optimized for massive write throughput and time-series data (which matches the chronological nature of chat history).

## Intermediate
**Q**: How do you handle Group Chats in this architecture?
*Answer*: In a 1-on-1 chat, Server A routes to Server B. In a group chat of 50 people, Server A looks up the 49 other users in the group. If they are spread across 10 different Chat Servers, Server A publishes the message to the Message Broker (Kafka/PubSub) 10 times, once for each target server. The target servers then push the message to the connected WebSockets.

## Advanced
**Q**: If a Chat Server suddenly crashes, what happens to the 1 million users connected to it?
*Answer*: Their TCP connections immediately drop. The clients' mobile apps will detect the dropped connection and automatically attempt to reconnect. The Load Balancer will route them to surviving Chat Servers. However, this creates a **"Thundering Herd"**. 1 million clients reconnecting simultaneously can crash the surviving servers. To prevent this, the client apps MUST implement **Exponential Backoff and Jitter** (`13-Resiliency.md`) when attempting to reconnect.

---

# Quiz

## Multiple Choice Questions
**1. What is the primary purpose of the Redis Session Store in this architecture?**
A) To cache recent messages for fast retrieval.
B) To track which Chat Server each active user is currently connected to.
C) To permanently store user profile pictures.
*Answer*: B. Without it, Chat Server A wouldn't know which server to forward Alice's message to.

## True or False
**WebSockets allow the server to send data to the client without the client initiating an HTTP request.**
*Answer*: True. Once the initial WebSocket handshake is completed, the connection is fully bi-directional.

---

# Exercises

## Beginner
Review the Go WebSocket example. What happens to the `serverHub.clients` map if a user closes their browser window suddenly? (Hint: look at the `ReadJSON` error handling).

## Intermediate
Research End-to-End Encryption (E2EE) used by WhatsApp (the Signal Protocol). In E2EE, does the Cassandra database store plain text messages? If not, how do users see their message history on a new device?

---

# Summary

Designing a real-time chat application shifts the architectural focus from stateless request/response to stateful connection management. WebSockets provide the real-time pipeline, but managing millions of concurrent connections requires lightweight concurrency models (like Go routines). By decoupling the persistent connections from the message routing logic via a Redis Session Store and a Pub/Sub broker, the system can scale horizontally to handle global traffic.

---

# Key Takeaways

* ✔ WebSockets are mandatory for low-latency, bi-directional real-time communication.
* ✔ Chat Servers must be stateful (holding TCP connections) but scalable via a central Session Store (Redis).
* ✔ Go is the industry standard for WebSocket servers due to its lightweight Goroutines.
* ✔ A fast-write NoSQL database (Cassandra) is required to handle the massive write-heavy workload of message history.

---

# Further Reading
* [Discord Engineering: How Discord handles 2.5 Million Concurrent Voice Users](https://discord.com/blog/how-discord-handles-two-and-half-million-concurrent-voice-users-using-webrtc)
* [Cassandra Data Modeling for Chat Apps](https://www.datastax.com/blog/data-modeling-messaging-system)

---

# Next Chapter
➡️ **Next:** `17-Design-Uber.md`
