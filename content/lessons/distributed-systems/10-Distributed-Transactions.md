# Distributed Transactions

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* The Monolith vs Microservices
* Two-Phase Commit (2PC)
* The Saga Pattern
* Architecture Diagram: Saga Choreography
* Compensating Transactions
* Best Practices
* Common Mistakes
* Debugging Guide
* Exercises
* Quiz
* Interview Questions
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

In a monolithic application with a single relational database, ensuring data integrity is easy: you use an **ACID Transaction**. If you deduct money from Account A and add it to Account B, and the server crashes halfway through, the database automatically rolls back both changes. 

In a microservice architecture, Account A lives in the Billing Service's database, and Account B lives in the Payment Service's database. Standard database transactions cannot span across networks. This introduces the incredibly difficult problem of **Distributed Transactions**.

---

# Learning Objectives

After completing this chapter you will be able to:

* Explain why ACID transactions are impossible across microservices without massive performance penalties.
* Understand the Two-Phase Commit (2PC) protocol and why modern architectures avoid it.
* Design robust distributed workflows using the **Saga Pattern**.
* Implement Compensating Transactions to handle rollbacks.

---

# Prerequisites

Before reading this chapter you should know:

* Database ACID Properties (Atomicity, Consistency, Isolation, Durability).
* Message Queues (`06-Message-Queues.md`).
* The Command Pattern (`15-Command.md` in Design Patterns).

---

# Why This Topic Exists

Imagine an E-Commerce application with three microservices: `OrderService`, `InventoryService`, and `PaymentService`.

A user places an order:
1. `OrderService` creates the order. (Success)
2. `InventoryService` deducts 1 item from stock. (Success)
3. `PaymentService` attempts to charge the credit card. (**FAILS: Card Declined!**)

The system is now in an inconsistent state. The order exists, the item is removed from inventory, but no money was collected. You must "roll back" the previous two steps. But because they are separate databases, there is no magic `db.Rollback()` command. You must architect the rollback manually.

---

# Two-Phase Commit (2PC)

The traditional way to solve this in enterprise software (Java EE, older RDBMS) was the **Two-Phase Commit protocol**.

It relies on a central **Coordinator**.
* **Phase 1 (Prepare)**: The Coordinator asks the Order, Inventory, and Payment databases: "Are you ready to commit this transaction? Lock your rows!" All databases reply "Yes" and lock the data.
* **Phase 2 (Commit)**: The Coordinator says "Go ahead and Commit!" The databases write the data and release the locks.

### Why 2PC is Terrible for the Cloud
1. **Performance Killer**: Locking rows across multiple network boundaries for long periods destroys throughput. 
2. **Single Point of Failure**: If the Coordinator crashes between Phase 1 and Phase 2, the databases are stuck holding locks forever (or until complex timeout resolution).
3. **No Polyglot Support**: MongoDB, Redis, and DynamoDB do not support standard 2PC protocols. 

---

# The Saga Pattern

Because 2PC doesn't scale, modern microservices use the **Saga Pattern**. 

A Saga is a sequence of local database transactions. Each service updates its own database and then publishes an event (to a Message Queue) to trigger the next local transaction in the sequence.

If a local transaction fails, the Saga executes a series of **Compensating Transactions** that undo the changes made by the preceding local transactions.

### Two Ways to Coordinate a Saga:
1. **Choreography (Event-Driven)**: No central controller. Services just listen to events and react. (Great for simple workflows).
2. **Orchestration (Command-Driven)**: A central Coordinator Service tells each service what to do. (Great for complex workflows with many steps).

---

# Architecture Diagram: Saga Choreography

```mermaid
flowchart TD
    Order[Order Service]
    Inv[Inventory Service]
    Pay[Payment Service]
    
    Order -- "1. OrderCreated Event" --> Inv
    
    Inv -- "2. InventoryReserved Event" --> Pay
    
    Pay -- "3. PaymentFailed Event" --> Inv
    
    note right of Pay: Card Declined! Must Rollback.
    
    Inv -- "4. InventoryRestored (Compensating Action)" --> Order
    Order -- "5. OrderCancelled (Compensating Action)" --> End((End))
```

---

# Compensating Transactions

A Compensating Transaction is a *semantic* rollback. It is a brand new transaction that reverses the business logic of the original transaction.

* **Original Action**: `UPDATE inventory SET stock = stock - 1 WHERE id = X`
* **Compensating Action**: `UPDATE inventory SET stock = stock + 1 WHERE id = X`

### The Catch
Because Sagas rely on *local* transactions, there is no "Isolation" (The 'I' in ACID). 
If the `InventoryService` deducts stock (Step 2), another user can see the reduced stock. If the Payment fails (Step 3) and the Saga restores the stock (Step 4), the other user saw temporary, dirty data. Sagas trade strict Consistency (ACID) for Eventual Consistency (BASE) to achieve high availability.

---

# Step-by-Step Implementation (Conceptual Choreography in Go)

While a full Saga requires a robust message queue (like RabbitMQ) and state persistence, the underlying logic is just event handlers.

**Inventory Service Handler:**
```go
func HandleOrderCreatedEvent(orderID int, itemID int) {
    err := DeductInventory(itemID)
    
    if err != nil {
        // We failed! Tell the OrderService to rollback.
        PublishEvent("InventoryFailed", orderID)
        return
    }
    
    // We succeeded! Tell the PaymentService to continue the Saga.
    PublishEvent("InventoryReserved", orderID)
}
```

**Payment Service Handler:**
```go
func HandleInventoryReservedEvent(orderID int, amount float64) {
    err := ChargeCreditCard(amount)
    
    if err != nil {
        // We failed! Start the chain of compensating transactions!
        PublishEvent("PaymentFailed", orderID)
        return
    }
    
    PublishEvent("PaymentSuccess", orderID)
}
```

**Inventory Service Compensating Handler (The Rollback):**
```go
func HandlePaymentFailedEvent(orderID int, itemID int) {
    // This is the COMPENSATING TRANSACTION
    fmt.Println("Payment failed. Restoring inventory...")
    RestoreInventory(itemID)
    
    // Continue passing the failure back down the chain
    PublishEvent("InventoryRestored", orderID)
}
```

---

# Production Use Cases

### 1. Travel Booking (Orchestrator Saga)
Booking a vacation requires Flights, Hotels, and Rental Cars. An `Orchestrator` microservice receives the request. It issues commands to the Flight service, then the Hotel service. If the Rental Car API returns an error, the Orchestrator knows exactly which services succeeded, and issues "Cancel Booking" commands to the Flight and Hotel services to refund the user.

### 2. Banking (Outbox Pattern)
When using Sagas, you must guarantee that when a local transaction commits, the Event is reliably published to the Message Queue. If the DB commits but the network crashes before sending to RabbitMQ, the Saga is permanently stuck. This is solved by the **Transactional Outbox Pattern**: saving the Event into an `outbox` table in the *same* database transaction, and using a background worker to reliably push the table contents to RabbitMQ.

---

# Best Practices

* **Design for Reversibility**: When designing a microservice, always ensure every API endpoint has a counterpart. If you have an `AddFunds` endpoint, you MUST build a `Refund` endpoint to support Saga compensation.
* **Idempotency is Mandatory**: Because compensating events might be delivered multiple times by the message queue, your compensation logic must be idempotent. Restoring inventory twice because you received the "PaymentFailed" message twice is a critical bug. (Covered in the next chapter).
* **Use Orchestration for Complex Sagas**: If a Saga has more than 3 or 4 steps, Choreography becomes a nightmare to debug (a "spiderweb" of events). Use a central Orchestrator to manage the state machine of the distributed transaction.

---

# Common Mistakes

### The "Point of No Return"
If your Saga includes sending a physical package via FedEx or sending an Email to a customer, those are actions that *cannot* be compensated (you cannot un-send an email). 
**Rule**: Put all non-compensatable actions at the very end of the Saga. Once you send the email, the Saga must never fail.

---

# Quiz

## Multiple Choice Questions
**1. Why is the Two-Phase Commit (2PC) protocol generally avoided in modern cloud microservices?**
A) It requires writing too much Go code.
B) It forces strict synchronous locking across multiple databases over a network, creating a massive performance bottleneck and a single point of failure.
C) It is impossible to implement securely.
*Answer*: B

## True or False
**In the Saga Pattern, if Step 3 fails, the system uses a special database command to instantly rollback the data saved in Steps 1 and 2.**
*Answer*: False. Sagas do not use database rollbacks across network boundaries. They use Compensating Transactions—brand new local transactions that semantically reverse the business logic (e.g., adding money back to an account after it was deducted).

---

# Interview Questions

## Beginner
**Q**: What is a Distributed Transaction?
*Answer*: It is a transaction that spans multiple physical databases or microservices. Because standard ACID database locks cannot span networks reliably, ensuring all databases commit or rollback together requires advanced architectural patterns.

## Intermediate
**Q**: What is the difference between Saga Choreography and Saga Orchestration?
*Answer*: In Choreography, services are completely decentralized; they listen for events from a message queue, perform local work, and emit new events. In Orchestration, a central controller service (the Orchestrator) explicitly commands each service to execute tasks or compensating transactions, tracking the state of the overall workflow in a database.

## Advanced
**Q**: Sagas lack Isolation (the 'I' in ACID). What anomalies can this cause, and how do you mitigate them?
*Answer*: Because intermediate local transactions are committed immediately, other services can read that intermediate state before the full Saga completes (a "Dirty Read"). If the Saga later rolls back, the data read by the other service is now invalid. To mitigate this, developers use Semantic Locks (adding a status column like `status=PENDING` to rows being acted upon by a Saga) to warn other services not to trust or modify that data until the Saga finishes.

---

# Summary

Distributed Transactions are the hardest problem in microservice architecture. Moving from the comforting safety of an ACID monolithic database to the eventual consistency of the Saga pattern requires a complete paradigm shift. By designing services with explicit Compensating Transactions, you can build systems that span the globe while maintaining reliable, albeit eventual, data integrity.

---

# Key Takeaways

* ✔ Standard ACID transactions cannot span microservices.
* ✔ 2PC is too slow and fragile for the cloud.
* ✔ The Saga pattern strings together local transactions via Events/Commands.
* ✔ Handle failures by executing Compensating Transactions to reverse the business logic.

---

# Further Reading
* [Microservices.io: Saga Pattern](https://microservices.io/patterns/data/saga.html)
* [The Transactional Outbox Pattern](https://microservices.io/patterns/data/transactional-outbox.html)

---

# Next Chapter
➡️ **Next:** `11-Idempotency.md`
