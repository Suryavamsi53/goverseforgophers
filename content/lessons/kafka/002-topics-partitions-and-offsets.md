# Topics, Partitions, and Offsets

Kafka stores events in a massive, immutable, append-only log on the hard drive. 
To organize these logs, Kafka uses a specific hierarchy.

## 1. The Topic

A **Topic** is a logical category for events. 
(e.g., `user_clicks`, `order_created`, `gps_locations`).

If a Go application wants to publish an event, it publishes it to a specific Topic. If a Go application wants to read data, it subscribes to a specific Topic.

## 2. The Partition (The Secret to Scale)

If millions of users are buying items simultaneously, a single `order_created` Topic on a single hard drive will instantly bottleneck and crash the physical server.

To solve this, Kafka physically splits every Topic into **Partitions**.

* You configure the `order_created` Topic to have **50 Partitions**.
* Kafka distributes these 50 Partitions evenly across 10 physical servers (Brokers).
* Now, when Go writes an order, the write load is distributed across 10 different hard drives simultaneously!

**The Partition Key:**
When you publish a message to Kafka, you provide a Key (e.g., the `UserID`) and a Value (the JSON payload). 
Kafka runs a hash function on the Key: `Hash("user_123") % 50 = Partition 14`.

Because the same Key *always* hashes to the same Partition, Kafka mathematically guarantees that all events for `user_123` are physically written to Partition 14 in the exact chronological order they occurred!

## 3. The Offset

Unlike RabbitMQ, which tracks exactly which messages have been processed and deletes them, Kafka puts the burden of tracking progress entirely on the Consumer.

Every message written to a Partition is assigned an incremental ID called an **Offset**.
* Message 1: Offset 0
* Message 2: Offset 1
* Message 3: Offset 2

The Go Consumer reads Offset 0, processes it, and then explicitly tells Kafka: *"I have successfully processed up to Offset 0"*. This is called **Committing the Offset**.

If the Go application crashes, when it reboots, it asks Kafka: *"Where did I leave off?"* Kafka replies *"You last committed Offset 0"*. The Go application immediately resumes reading from Offset 1.

## 4. The Ordering Guarantee

This is the most critical concept in Kafka:
**Kafka only guarantees strict chronological ordering WITHIN a single Partition!**

If `Event A` and `Event B` are written to the same Topic, but `Event A` goes to Partition 1, and `Event B` goes to Partition 2, there is absolutely no guarantee which event the Go consumer will read first!

This is why the **Partition Key** is so vital. If `Event A` (Create User) and `Event B` (Update User) belong to the same user, you MUST use the `UserID` as the Partition Key. This forces both events into the exact same Partition, mathematically guaranteeing that the Go consumer will read `Create User` before it reads `Update User`.
