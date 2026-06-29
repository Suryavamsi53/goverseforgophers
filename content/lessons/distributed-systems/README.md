# Distributed Systems in Go Curriculum

Welcome to the **Distributed Systems in Go** module. Building on your knowledge of Concurrency and Design Patterns, this curriculum will teach you how to architect software that runs across multiple machines, scaling horizontally to handle massive loads while remaining resilient to network failures.

Go was explicitly designed by Google to build distributed systems. From Kubernetes to Docker, Terraform to Prometheus, Go is the language of the cloud.

## Curriculum Overview

### Part 1: Foundations of Distributed Systems
* **01-CAP-Theorem.md**: Understanding Consistency, Availability, and Partition Tolerance, and why you can only pick two.
* **02-Fallacies-of-Distributed-Computing.md**: The 8 false assumptions every developer makes about the network.
* **03-Time-and-Clocks.md**: Why `time.Now()` is dangerous across multiple servers (NTP, Vector Clocks, Logical Clocks).

### Part 2: Communication Protocols
* **04-RPC-vs-REST.md**: Moving beyond CRUD to Remote Procedure Calls.
* **05-gRPC-and-Protobuf.md**: The industry standard for high-performance microservice communication.
* **06-Message-Queues.md**: Asynchronous communication, decoupling, and buffering (Kafka, RabbitMQ, Redis Pub/Sub).

### Part 3: Reliability and Resiliency
* **07-Circuit-Breaker-Pattern.md**: Protecting your system from cascading failures.
* **08-Retries-and-Exponential-Backoff.md**: Handling transient network blips safely.
* **09-Timeouts-and-Deadlines.md**: Propagating cancellation signals across network boundaries using Context.

### Part 4: Consistency and State
* **10-Distributed-Transactions.md**: The Saga Pattern and Two-Phase Commit (2PC).
* **11-Idempotency.md**: Ensuring safe retries in payment systems and APIs.
* **12-Consensus-Algorithms.md**: A high-level overview of Raft and Paxos (how etcd and Consul work).

### Part 5: Observability
* **13-Distributed-Tracing.md**: Tracking a single request as it jumps across 10 different microservices (OpenTelemetry, Jaeger).
* **14-Centralized-Logging.md**: Aggregating logs in a distributed environment.
* **15-Metrics-and-Monitoring.md**: Exporting system health data (Prometheus).

---

Let's begin the journey into the cloud!
