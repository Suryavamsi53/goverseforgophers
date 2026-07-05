# Producer ACKs and Durability

When your Go application (the Producer) sends a message to Kafka, how do you know it was actually saved to the hard drive? What if the network drops a packet, or the Kafka server loses power the millisecond the message arrives?

In enterprise Go applications, you must configure the Producer's **ACK (Acknowledgment) Level**.

## 1. Replication in Kafka

First, you must understand that Kafka doesn't just store data on one server. 
For High Availability, Kafka replicates every Partition across multiple physical Brokers.

If your Topic has a `ReplicationFactor=3`:
* Broker 1 holds the **Leader** Partition. (All Go Producers write to the Leader).
* Broker 2 and Broker 3 hold **Follower** Partitions. (They constantly pull data from the Leader to keep a backup copy).

## 2. The `acks` Configuration

When you configure your Go Kafka Producer (e.g., using `confluent-kafka-go` or `segmentio/kafka-go`), you must set the `acks` parameter. This defines how much confirmation your Go app requires before it considers the write "Successful".

### `acks = 0` (Fire and Forget)
* **How it works**: The Go app blasts the message over the TCP socket and instantly assumes success. It does not wait for a response from Kafka.
* **Pros**: Impossibly fast. Highest possible throughput.
* **Cons**: Massive data loss. If the Kafka server is offline, or the network is broken, the Go app silently drops the message into the void and has no idea.
* **Use Case**: High-volume, low-value telemetry (e.g., tracking mouse movements on a webpage).

### `acks = 1` (Leader Acknowledgment)
* **How it works**: The Go app waits for the Leader Broker to save the message to its own hard drive. Once saved, the Leader replies "Success".
* **Pros**: Very fast, much safer than `0`.
* **Cons**: Edge-case data loss! Imagine the Leader saves the data and replies "Success". One millisecond later, the Leader Broker is struck by lightning and explodes. The Follower Brokers didn't have time to copy the data yet! The data is permanently lost, even though the Go app was told it was successful!
* **Use Case**: Most standard web applications and logs.

### `acks = all` (Highest Durability)
* **How it works**: The Go app sends the message to the Leader. The Leader waits until ALL Follower Brokers have safely copied the data to their hard drives. Only then does the Leader reply "Success" to the Go app.
* **Pros**: Mathematically impossible to lose data (as long as one Follower survives).
* **Cons**: High latency. The Go app has to wait for 3 separate servers to perform disk I/O before continuing.
* **Use Case**: Financial transactions, Billing systems, and core Domain Events.

## 3. Retries and Idempotence

If `acks=all` times out, your Go Producer will automatically retry sending the message.

**The Danger**: What if the Leader safely wrote the message, but the network failed *while sending the "Success" ACK back to Go*? 
The Go app thinks it failed, so it retries. Now Kafka has the exact same message written twice (Duplicate Data)!

To fix this, modern Kafka Producers enable **Idempotence** (`enable.idempotence=true`). 
The Go Producer secretly attaches a unique Sequence ID to every message. If a retry occurs, the Kafka Broker sees the duplicate Sequence ID and safely ignores the second message, guaranteeing Exactly-Once writing!
