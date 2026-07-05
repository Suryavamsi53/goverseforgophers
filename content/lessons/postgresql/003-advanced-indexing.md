# Advanced Indexing (GIN, BRIN, Partial)

We already know that a standard B-Tree index changes a query from $O(N)$ (Sequential Scan) to $O(log N)$. 

However, B-Trees are designed for simple equality (`=`) and range (`>`, `<`) checks on standard columns (integers, short strings). When your Go application starts doing full-text search, filtering by JSON keys, or querying massive 100-gigabyte time-series tables, a B-Tree will fail you.

PostgreSQL offers specialized index types for these enterprise scenarios.

## 1. GIN (Generalized Inverted Index)

If you have a `tags` column stored as a PostgreSQL Array (e.g., `['go', 'backend', 'api']`), and you want to find all articles tagged with "go", a B-Tree is useless. A B-Tree can only compare the *entire* array, not the individual elements inside it.

You must use a **GIN Index**.

```sql
CREATE INDEX idx_articles_tags ON articles USING GIN (tags);
```

An Inverted Index maps individual elements to their rows (exactly like the index at the back of a textbook). 
It is the only way to efficiently query **Arrays**, **Full-Text Search (tsvector)**, and **JSONB**.

*Note: GIN indexes are extremely fast for Reading, but very slow for Writing, as inserting one array requires updating multiple separate index entries!*

## 2. BRIN (Block Range Index)

Imagine an `audit_logs` table with 50 billion rows. Every log has a `created_at` timestamp.
If you create a standard B-Tree index on `created_at`, the index file itself will be 100 Gigabytes! It won't fit in RAM, rendering it useless.

Because `created_at` naturally increases over time (it is perfectly physically ordered on the hard drive), you can use a **BRIN Index**.

```sql
CREATE INDEX idx_logs_time ON audit_logs USING BRIN (created_at);
```

Instead of tracking every single row, a BRIN index divides the table into 1 Megabyte "Blocks". For each block, it simply records the Minimum and Maximum timestamp. 
When you query `WHERE created_at = '2026-01-01'`, Postgres checks the BRIN index, instantly skips the 99% of blocks that don't overlap with that date, and only scans the block that does.

**The Superpower:** A BRIN index on a 50 Billion row table is only 5 Megabytes in size! It is the ultimate weapon for Time-Series data.

## 3. Partial Indexes

If you have a `users` table with 10 million rows, but only 5,000 of them are `active = true`, creating an index on the `active` column is a massive waste of space. 99% of the index is tracking `false` values that you will never query!

You can create an index that only tracks a subset of the data: a **Partial Index**.

```sql
CREATE INDEX idx_active_users ON users(email) WHERE active = true;
```

This index will only contain 5,000 entries. It takes up almost zero RAM, updates instantly, and makes queries like `SELECT email FROM users WHERE active = true` blazingly fast. 
Use Partial Indexes heavily for soft-deletes (`deleted_at IS NULL`) or boolean flags!
