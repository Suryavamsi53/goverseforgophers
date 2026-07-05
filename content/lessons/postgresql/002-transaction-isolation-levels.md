# Transaction Isolation Levels

In the Database module, we learned how to use `db.BeginTx()` in Go to wrap multiple SQL statements into a single atomic Transaction.

However, if you have two independent Go servers executing Transactions at the exact same time on the exact same data, strange phenomena can occur. PostgreSQL allows you to choose exactly how strictly you want to isolate these concurrent transactions.

## 1. The Three Concurrency Anomalies

Before we discuss the Isolation Levels, you must understand the bugs they prevent:

1. **Dirty Read**: Transaction A updates a row, but hasn't committed yet. Transaction B reads the uncommitted data. Transaction A rolls back. Transaction B now has data in memory that never existed!
2. **Non-Repeatable Read**: Transaction A reads a row. Transaction B updates that row and commits. Transaction A reads the exact same row again, but the data has magically changed in the middle of the transaction!
3. **Phantom Read**: Transaction A runs `SELECT COUNT(*) WHERE age > 18` (returns 5). Transaction B inserts a new 20-year-old and commits. Transaction A runs the exact same query again, and the count is magically 6! A "phantom" row appeared.

## 2. Read Committed (The Postgres Default)

By default, every transaction you start in PostgreSQL operates at the **Read Committed** level.

* **Prevents**: Dirty Reads. (You will never see uncommitted data).
* **Allows**: Non-Repeatable Reads and Phantom Reads.

If your Go application runs a report that takes 10 seconds, and it queries the `users` table twice, the data might change between the first and second query! For 95% of web applications, this is perfectly fine and provides phenomenal concurrent performance.

## 3. Repeatable Read

If you are building a financial ledger, the data *cannot* change underneath you.

```go
opts := &sql.TxOptions{Isolation: sql.LevelRepeatableRead}
tx, _ := db.BeginTx(ctx, opts)
```

At this level, Postgres takes a "Snapshot" of the entire database the millisecond your transaction starts. 
* **Prevents**: Dirty Reads, Non-Repeatable Reads. (If you query a row 100 times, it will always return the exact same data, even if another transaction modified it!)
* **Postgres Bonus**: In Postgres, Repeatable Read *also* prevents Phantom Reads! (This is not required by the SQL standard, but Postgres's MVCC implementation is that good).

**The Catch (Serialization Failures):**
If Transaction A and Transaction B both read Row 1 (value=100), and both try to `UPDATE Row 1 SET value = value + 10`, Transaction B will instantly fail with a `could not serialize access due to concurrent update` error. Your Go application MUST catch this error and retry the transaction from scratch!

## 4. Serializable

The ultimate level of data integrity.

* **Prevents**: Everything. It guarantees that the result of running transactions concurrently is mathematically identical to running them one by one, sequentially.

```go
opts := &sql.TxOptions{Isolation: sql.LevelSerializable}
```

Postgres achieves this by maintaining complex "Predicate Locks" in memory, watching every single row you touch. If it detects even a hint of a cross-dependency between two transactions, it mercilessly kills one of them.

**When to use Serializable:**
Never, unless you are building a nuclear reactor control system or a highly complex double-entry accounting ledger. The performance overhead is massive, and your Go application will have to retry 50% of its transactions due to Serialization failures. Stick to `Read Committed`!
