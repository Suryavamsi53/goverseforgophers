# Distributed Transactions

As we discussed in the Microservices module (Saga Pattern), rolling back an action across multiple independent databases is the hardest problem in distributed systems.

If the `OrderService` (Postgres) succeeds, but the `InventoryService` (Redis) fails, you cannot execute a single SQL `ROLLBACK` command. 

There are two primary ways to solve this.

## 1. Two-Phase Commit (2PC)

Two-Phase Commit is a highly pessimistic, blocking protocol that guarantees absolute Consistency (CP system). It relies on a central Coordinator.

**Phase 1: The Prepare Phase (Voting)**
1. The Coordinator asks `OrderDB` and `InventoryDB`: "Are you ready to commit?"
2. `OrderDB` checks its constraints, locks the rows, and replies "Yes".
3. `InventoryDB` checks its constraints, locks the rows, and replies "Yes".

**Phase 2: The Commit Phase**
1. If EVERY database replied "Yes", the Coordinator sends the "Commit!" command to all of them simultaneously.
2. If ANY database replied "No" (or timed out), the Coordinator sends the "Rollback!" command to all of them.

### Why 2PC is Terrible for Microservices
If `InventoryDB` replies "Yes" during Phase 1, it must physically lock those rows to guarantee it can commit them in Phase 2. 
If the Coordinator crashes before it sends the Phase 2 command, `InventoryDB` is stuck waiting forever, with its rows permanently locked! 2PC creates massive single points of failure and terrible latency.

## 2. The Saga Pattern (Eventual Consistency)

Instead of pessimistic locking, modern microservices use the **Saga Pattern**, which embraces Eventual Consistency (AP system).

A Saga is a sequence of local, independent transactions. There are no global locks.

1. `OrderService` creates the order (Status: PENDING) in its own Postgres DB. It publishes an `OrderCreated` event.
2. `InventoryService` hears the event, deducts the stock in its own DB, and publishes `InventoryReserved`.
3. `OrderService` hears the event and updates the order (Status: COMPLETE).

### Compensating Actions
What if Step 2 fails (out of stock)?
You cannot "rollback" the `OrderService`, because it already committed the transaction locally!
Instead, the `InventoryService` publishes an `InventoryFailed` event.
The `OrderService` hears this, and executes a **Compensating Action**: a brand new, independent SQL transaction that updates the order (Status: CANCELLED).

**Sagas are infinitely faster and more scalable than 2PC, but they are much harder to debug.**

## 3. Managing Sagas (Temporal.io)

Because Sagas rely on a chain of asynchronous events, a network drop can leave your system in a zombie state (e.g., the Order is PENDING forever).

To solve this, Enterprise Go teams use workflow engines like **Temporal**.

Temporal allows you to write Go code that acts as a flawless, crash-proof Orchestrator.

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    // 1. Execute Order
    err := workflow.ExecuteActivity(ctx, CreateOrderActivity, order).Get(ctx, nil)
    
    // 2. Execute Inventory
    err = workflow.ExecuteActivity(ctx, ReserveInventoryActivity, order).Get(ctx, nil)
    if err != nil {
        // INVENTORY FAILED! Execute the Compensation!
        // Temporal guarantees this will run, even if this workflow server 
        // loses power and reboots!
        workflow.ExecuteActivity(ctx, CancelOrderActivity, order).Get(ctx, nil)
        return err
    }
    
    return nil
}
```
Temporal tracks the state of every workflow in its own database. If your Go server crashes mid-execution, when it reboots, Temporal seamlessly resumes the Go function from the exact line of code it died on!
