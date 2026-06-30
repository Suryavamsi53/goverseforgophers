# Load Balancers: L4 vs L7, Nginx, and HAProxy

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* What is a Load Balancer?
* L4 vs L7 Load Balancing
* Common Routing Algorithms
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

In `02-Scalability.md`, we learned that horizontal scaling (adding more servers) is the key to handling massive traffic. However, if you have 100 stateless web servers, how do clients know which one to talk to? You can't ask a user to manually pick an IP address.

This is where the **Load Balancer (LB)** comes in. It sits between your clients and your backend servers, acting as a traffic cop. It accepts incoming requests and distributes them across your pool of healthy servers, ensuring no single server is overwhelmed while others sit idle.

---

# Learning Objectives

After completing this chapter you will be able to:

* Explain the core function of a Load Balancer in a distributed system.
* Differentiate between Layer 4 (Transport) and Layer 7 (Application) load balancing.
* Understand common load-balancing algorithms (Round Robin, Least Connections, IP Hash).
* Implement robust health check endpoints in your Go services.

---

# Prerequisites

Before reading this chapter you should know:

* Vertical vs Horizontal Scaling (`02-Scalability.md`).
* The OSI Model Layers (`03-Network-Protocols.md`).

---

# Why This Topic Exists

During a system design interview, placing a Load Balancer in front of your web tier is usually the first box you draw. However, interviewers will drill down: *How* is it routing traffic? Is it doing SSL termination? What happens if a backend server crashes mid-request? 

Understanding the mechanics of load balancing—specifically the differences between L4 and L7—allows you to optimize for either raw throughput (L4) or complex routing logic based on HTTP headers (L7).

---

# What is a Load Balancer?

A Load Balancer is either a hardware device (legacy) or a software application (like Nginx, HAProxy, or AWS ALB) that reverse-proxies traffic. 

Its primary jobs are:
1. **Traffic Distribution**: Spreading incoming requests across multiple backend servers.
2. **High Availability**: Automatically detecting when a backend server fails (via Health Checks) and stopping traffic to it.
3. **SSL Termination**: Decrypting HTTPS traffic at the load balancer so backend servers don't have to waste CPU cycles on cryptography.

---

# L4 vs L7 Load Balancing

Load balancers operate at different layers of the OSI model, primarily Layer 4 (Transport) and Layer 7 (Application).

### Layer 4 (Transport Layer) Load Balancing
* **How it works**: It routes traffic based purely on network information: the IP address and the TCP/UDP port. It does not look at the contents of the packet.
* **Pros**: Extremely fast. Because it doesn't parse HTTP headers or inspect the payload, it uses very little CPU and can handle millions of connections per second.
* **Cons**: "Dumb" routing. It cannot route traffic based on URL paths (e.g., sending `/api` to one server and `/images` to another).

### Layer 7 (Application Layer) Load Balancing
* **How it works**: It fully terminates the TCP connection, parses the HTTP payload (headers, URL, cookies), and makes routing decisions based on that content.
* **Pros**: Highly intelligent. You can route traffic based on the URL path, the user's language header, or a session cookie (Sticky Sessions). 
* **Cons**: Slower and more CPU intensive. Parsing HTTP traffic takes significantly more processing power than simple IP/Port routing.

---

# Common Routing Algorithms

How does the Load Balancer choose *which* server gets the next request?

1. **Round Robin**: Requests are distributed across the group of servers sequentially (Server 1, then Server 2, then Server 3, then back to 1). Best for servers with identical specs.
2. **Least Connections**: Sends the request to the server with the fewest active connections. Best if requests take varying amounts of time to process.
3. **IP Hash**: The client's IP address is mathematically hashed to determine the server. Ensures that a specific user always gets routed to the same server.
4. **Weighted Round Robin**: You can assign "weights" (e.g., Server A is twice as powerful as Server B, so it gets 2 requests for every 1 that B gets).

---

# Code Examples & Good Principles

For a Load Balancer to know if a server is alive, the server must provide a **Health Check** endpoint. The LB will ping this endpoint every few seconds. If it fails, the server is removed from the rotation.

### Implementing a Robust Health Check in Go (Good Principle)

**Bad Practice**: A health check that just returns `200 OK` regardless of internal state. If your Go app has lost its database connection, it shouldn't receive traffic!

**Good Practice**: A "Deep" health check that verifies critical dependencies (DB, Redis) before returning healthy.

```go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

type HealthStatus struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Time     string `json:"time"`
}

// Principle: Health checks must verify dependencies, not just the HTTP listener.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	status := HealthStatus{
		Status:   "ok",
		Database: "ok",
		Time:     time.Now().Format(time.RFC3339),
	}
	statusCode := http.StatusOK

	// Ping the database with a short timeout
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Printf("Health check failed: DB down: %v", err)
		status.Status = "error"
		status.Database = "disconnected"
		// Principle: Return 503 Service Unavailable so the LB stops sending traffic
		statusCode = http.StatusServiceUnavailable 
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(status)
}

func main() {
	// Initialize dummy DB (will fail in real life without credentials)
	var err error
	db, err = sql.Open("postgres", "postgres://user:pass@localhost/db")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/healthz", healthCheckHandler)
	
	log.Println("Server starting on :8080 with deep health checks")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

# Architecture Diagram

```mermaid
flowchart TD
    Client((Client)) --> LB{L7 Load Balancer<br>e.g., Nginx / AWS ALB}
    
    subgraph Target Group: API
        LB -- "/api/*" --> API1[API Server 1]
        LB -- "/api/*" --> API2[API Server 2]
    end
    
    subgraph Target Group: Static Assets
        LB -- "/images/*" --> CDN1[Static File Server]
    end
    
    API1 -. "Ping /healthz" .-> LB
    API2 -. "Ping /healthz" .-> LB
    CDN1 -. "Ping /healthz" .-> LB
    
    style LB fill:#f96,stroke:#333,stroke-width:2px
```

---

# Real-World Analogy

* **The System**: A massive hospital.
* **The Load Balancer**: The receptionist at the front desk.
* **Round Robin**: The receptionist assigns new patients to doctors in a strict circle, regardless of how busy the doctor is.
* **Least Connections**: The receptionist looks at which doctor has the fewest people in their waiting room and sends the new patient there.
* **L4 Routing**: The receptionist only looks at your ZIP code (IP address) to direct you to a wing.
* **L7 Routing**: The receptionist reads your intake form (HTTP Payload) and sees you have a broken leg, so they route you specifically to the Orthopedics wing (`/orthopedics`).

---

# Interview Questions

## Beginner
**Q**: What is the purpose of a health check?
*Answer*: To allow the Load Balancer to actively monitor the status of backend servers. If a server fails the health check, the LB stops routing traffic to it, preventing users from experiencing errors.

## Intermediate
**Q**: When would you use L4 load balancing instead of L7?
*Answer*: You use L4 when you need extreme performance and low latency, and you don't care about the content of the traffic. For example, balancing raw TCP traffic for a multiplayer game server or a database cluster where HTTP headers don't exist.

## Advanced
**Q**: What is "Sticky Sessions" (Session Affinity) and why is it considered an anti-pattern in modern distributed systems?
*Answer*: Sticky Sessions force the LB to always route a specific user to the same specific backend server (often by injecting a cookie). This is required if the backend server stores session state in local RAM. It's an anti-pattern because it breaks horizontal scaling—if that specific server crashes, the user loses their session. The modern solution is to make servers stateless and store sessions in a shared cache like Redis.

---

# Quiz

## Multiple Choice Questions
**1. Which load balancing algorithm routes a client to the same server every time based on their IP address?**
A) Round Robin
B) IP Hash
C) Least Connections
*Answer*: B

## True or False
**An L7 Load Balancer can inspect the HTTP URL path to make routing decisions.**
*Answer*: True. Because it operates at the application layer, it can parse the HTTP payload and route traffic based on URLs, headers, or cookies.

---

# Exercises

## Beginner
Look up the documentation for Nginx. Find the configuration block used to define an `upstream` block (a pool of backend servers).

## Intermediate
Modify the Go health check code example to also check the availability of a hypothetical Redis cache, in addition to the PostgreSQL database. If either fails, the health check should return a `503 Service Unavailable`.

---

# Summary

Load Balancers are the unsung heroes of high availability. They distribute traffic, mask server failures from users, and can intelligently route requests based on application logic (L7). However, a load balancer is only as smart as the health checks it relies upon. Always ensure your Go applications expose deep, meaningful `/healthz` endpoints so the infrastructure knows exactly when a service is struggling.

---

# Key Takeaways

* ✔ L4 Load Balancing is fast and routes based on IP/Port.
* ✔ L7 Load Balancing is smarter and routes based on HTTP URLs/Headers.
* ✔ Never use "Sticky Sessions" if you can avoid it; keep your backend stateless.
* ✔ Implement "Deep" health checks in Go to verify critical database/cache connections.

---

# Further Reading
* [NGINX Load Balancing Documentation](https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/)
* [AWS Application Load Balancer vs Network Load Balancer](https://aws.amazon.com/elasticloadbalancing/features/)

---

# Next Chapter
➡️ **Next:** `06-Caching.md`
