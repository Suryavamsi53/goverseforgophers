# Time and Clocks in Distributed Systems

If you write a log entry on your laptop, `time.Now()` is perfectly accurate. But in a distributed system with 1,000 servers across the globe, what does `time.Now()` actually mean?

## 1. The Clock Drift Problem

Every server has a physical quartz crystal clock on its motherboard. These crystals vibrate at slightly different frequencies due to temperature and manufacturing defects. Over the course of a day, Server A's clock might drift 5 milliseconds faster than Server B.

If you rely on `time.Now()` to order events, you will experience impossible bugs:
1. Server B creates an Order at `12:00:00.005`.
2. Server A marks the Order as Paid at `12:00:00.002`.

According to the database timestamps, the Order was paid **before** it was created!

### Network Time Protocol (NTP)
Servers use NTP to sync their clocks with atomic clocks over the internet. However, network latency makes NTP imprecise. NTP can keep clocks synchronized within a few milliseconds, but it can never guarantee exact nanosecond ordering.

## 2. Logical Clocks (Lamport Clocks)

Because Physical Time (`time.Now()`) is unreliable, Leslie Lamport invented **Logical Clocks** in 1978.

Instead of tracking the *time* an event happened, a Lamport Clock tracks the *causality* (the order) of events using a simple integer counter.

1. Every server maintains a local counter (`clock = 0`).
2. Before a server does an action, it increments its clock: `clock++`.
3. When Server A sends a network message to Server B, it attaches its current clock value (e.g., `[Message, clock=5]`).
4. When Server B receives the message, it updates its own clock to `max(local_clock, received_clock) + 1`.

This mathematically guarantees that if Event X caused Event Y, the logical clock of Y will always be greater than X, regardless of physical quartz drift!

## 3. Vector Clocks (Conflict Resolution)

Lamport Clocks can order events, but they cannot detect **Concurrent Conflicts** (when two users update the same record on two different servers at the exact same time).

**Vector Clocks** solve this. Instead of a single integer, a Vector Clock is an array of integers, representing the clock state of *every* node in the system.

*(Example of a Vector Clock: `[NodeA: 2, NodeB: 1, NodeC: 0]`)*

Databases like Amazon DynamoDB and Riak use Vector Clocks. If Node A and Node B accept a write at the exact same time, their Vector Clocks will diverge. When the database tries to sync the nodes, it compares the arrays, realizes a concurrent conflict occurred, and forces the application (or the user) to resolve the merge conflict (just like Git!).

## 4. Google TrueTime (The Hardware Solution)

Google got tired of dealing with Logical Clocks, so they solved the problem with hardware. 

For their globally distributed database, **Google Spanner**, they installed physical GPS receivers and Atomic Clocks into every single server rack in their datacenters. This API, called TrueTime, guarantees that clock drift between any two servers worldwide is bounded to exactly 7 milliseconds. Spanner literally waits 7 milliseconds before committing a transaction to guarantee absolute, physical global ordering.
