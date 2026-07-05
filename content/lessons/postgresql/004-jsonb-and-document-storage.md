# JSONB and Document Storage

Many teams choose MongoDB because they want the flexibility of a "Schema-less" Document database. They want to store massive JSON payloads without having to define 50 columns in a SQL table.

In 2014, PostgreSQL introduced the **`JSONB`** data type, effectively neutralizing MongoDB's primary advantage. You can use PostgreSQL as a world-class NoSQL database, while still retaining full ACID compliance and relational Joins!

## 1. JSON vs JSONB

PostgreSQL actually has two JSON data types:
* `JSON`: Stores an exact copy of the raw text string you send it. (Fast to insert, incredibly slow to query).
* `JSONB`: Converts the JSON text into a highly-compressed, binary tree format on the hard drive. (Slightly slower to insert, blazingly fast to query).

**Rule: You should never use the `JSON` type. Always use `JSONB`.**

## 2. Using JSONB in Go

If your Go application accepts Webhooks from Stripe, the JSON payload might have 300 fields. You don't want to create a SQL column for every field. You just want to dump the entire payload into a single column.

```sql
CREATE TABLE stripe_events (
    id SERIAL PRIMARY KEY,
    event_type VARCHAR(50),
    payload JSONB
);
```

In your Go application, you simply map the JSONB column to a standard `map[string]interface{}` or a custom struct using the `database/sql` driver and the `encoding/json` package!

```go
// 1. Marshal the struct into JSON bytes
payloadBytes, _ := json.Marshal(stripeEvent)

// 2. Insert directly into Postgres! Postgres natively understands it!
db.ExecContext(ctx, "INSERT INTO stripe_events (event_type, payload) VALUES ($1, $2)", 
    "payment_intent.succeeded", payloadBytes)
```

## 3. Querying JSONB (The Operators)

Because the data is stored in binary, you can write native SQL queries that dive deep into the JSON tree!

**Extracting a text value (`->>`):**
```sql
-- Find all events where the nested 'currency' key is 'usd'
SELECT id FROM stripe_events 
WHERE payload->'data'->'object'->>'currency' = 'usd';
```

**Checking for key existence (`?`):**
```sql
-- Find all events where the payload contains a top-level key called 'error'
SELECT id FROM stripe_events WHERE payload ? 'error';
```

**Checking for key-value containment (`@>`):**
```sql
-- Find all events where the JSON contains this exact sub-object!
SELECT id FROM stripe_events WHERE payload @> '{"status": "failed"}';
```

## 4. Indexing JSONB

Running the queries above on a 10-million row table will cause a Sequential Scan and take 5 seconds. 
Because `JSONB` is binary, we can index it!

As covered in the Advanced Indexing lesson, we cannot use a B-Tree. We must use a **GIN Index**.

```sql
-- Creates an index on EVERY key and value inside the JSON payload!
CREATE INDEX idx_stripe_payload ON stripe_events USING GIN (payload);
```

With this GIN index, the `@>` (containment) query above will execute in 1 millisecond. You have successfully built a massive, schema-less, highly indexed NoSQL database directly inside PostgreSQL!
