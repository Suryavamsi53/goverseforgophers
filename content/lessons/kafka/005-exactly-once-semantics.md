# Exactly-Once Semantics (Transactions)

Distributed systems are governed by the Two Generals' Problem. When sending messages over an unreliable network, you must choose between two delivery guarantees:

1. **At-Most-Once**: The message is sent. If the network fails, the message is lost forever. (No duplicates, but data loss).
2. **At-Least-Once**: The message is sent. If an ACK is not received, the message is retried. (No data loss, but guarantees duplicates).

For decades, engineers accepted that **Exactly-Once Delivery** was mathematically impossible over a network. 
However, Apache Kafka achieves "Exactly-Once Semantics" (EOS) through a brilliant implementation of Distributed Transactions.

## 1. The Duplicate Problem (At-Least-Once)

Imagine a Bank Microservice:
1. The Go Consumer reads a `Deposit_$100` message from Kafka.
2. The Go Consumer runs `UPDATE accounts SET balance = balance + 100` in Postgres.
3. The Go Consumer attempts to commit the Kafka Offset.
4. **CRASH!** The network dies before the Offset reaches Kafka.

When the Go Consumer reboots, Kafka says: *"You never committed that message!"* and hands the Go app the exact same `Deposit_$100` message again. The user gets $200!

## 2. The Idempotency Solution (Manual EOS)

Historically, engineers solved this manually using **Idempotency** (as discussed in the Distributed Systems module).

You must maintain an `idempotency_keys` table in PostgreSQL.
1. Read the Kafka message (ID: 42).
2. Open a Postgres Transaction.
3. Check `SELECT * FROM idempotency_keys WHERE id = 42`. If it exists, abort!
4. `UPDATE accounts...`
5. `INSERT INTO idempotency_keys (id) VALUES (42)`.
6. Commit Postgres Transaction.
7. Commit Kafka Offset.

If the app crashes at step 7, Kafka will replay the message. But step 3 will catch it and safely ignore it! You achieved Exactly-Once processing!

## 3. Kafka Transactions (Native EOS)

What if your Go application reads a message from `Topic A`, processes it, and writes the result to `Topic B`? (A pure stream-processing application, like a Kafka Streams or Flink job).

You can't use a Postgres idempotency table here, because you are only talking to Kafka!

Kafka introduced **Transactions** to solve the `Consume -> Process -> Produce` loop.

```go
// 1. Initialize a Transactional Producer
producer, _ := kafka.NewProducer(&kafka.ConfigMap{
    "transactional.id": "my-go-processor", // Required for EOS!
})
producer.InitTransactions(nil)

// 2. Start the Transaction
producer.BeginTransaction()

// 3. Produce the transformed data to Topic B
producer.Produce(&kafka.Message{TopicPartition: topicB, Value: []byte("processed")}, nil)

// 4. THE MAGIC: Commit the Producer writes AND the Consumer Offset AT THE SAME TIME!
// Kafka guarantees this is an atomic operation!
producer.SendOffsetsToTransaction(consumer.Position(), consumer.Assignment())
producer.CommitTransaction()
```

### How it works under the hood
Kafka uses a **Transaction Coordinator** and a Two-Phase Commit (2PC) protocol. 
It ensures that the output messages written to Topic B, and the Offset committed on Topic A, are written atomically. If the Go app crashes mid-transaction, Kafka automatically rolls back the produced messages in Topic B, making them completely invisible to downstream consumers!

**Rule of Thumb:**
* If your Go app reads from Kafka and writes to a Database -> Use **Idempotency Keys** in the Database.
* If your Go app reads from Kafka and writes back to Kafka -> Use **Kafka Transactions**.
