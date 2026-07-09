# Level 2 — CRUD Operations

*(Based on official PostgreSQL documentation, current stable series: PostgreSQL 18.)*

## 1. Learning Objectives

* **What you'll learn**: How PostgreSQL executes Data Definition Language (database/table/column operations) and the four core Data Manipulation Language operations — Create (`INSERT`), Read (`SELECT`), Update (`UPDATE`), Delete (`DELETE`) — including the `RETURNING` clause that makes Postgres's DML uniquely efficient.
* **Why it matters**: CRUD operations are the vast majority of what any backend application actually does against its database. Small inefficiencies here (an extra round-trip, a missing index, a badly-scoped `UPDATE`) get multiplied by every request your system serves, so getting these fundamentals right has an outsized effect on real-world performance and correctness.

---

## 2. Topics Covered

* Database Operations
* Table Operations
* Column Operations
* Insert Operations (`RETURNING`)
* Read Operations
* Update
* Delete

---

## 3. Deep Dive: Concepts

### 3.1 Database Operations

Database-level DDL manages the top-level namespace of a cluster:
```sql
CREATE DATABASE fin_app
  WITH OWNER = app_user
       ENCODING = 'UTF8'
       TEMPLATE = template0;

ALTER DATABASE fin_app SET timezone TO 'UTC';

DROP DATABASE IF EXISTS fin_app_old;
```
Note: `CREATE DATABASE` cannot run inside a transaction block (it's one of a handful of commands that manage their own internal transaction), and by default it's cloned from `template1` — using `template0` avoids inheriting any custom objects that might have been added to `template1`.

### 3.2 Table Operations

```sql
CREATE TABLE users (
    id         BIGSERIAL PRIMARY KEY,
    email      TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE users RENAME TO app_users;
ALTER TABLE app_users ADD CONSTRAINT chk_email CHECK (email LIKE '%@%');

TRUNCATE TABLE app_users;      -- fast wipe, resets storage, not transaction-safe row-by-row
DROP TABLE IF EXISTS app_users;
```
`TRUNCATE` is dramatically faster than `DELETE FROM table` for wiping an entire table because it deallocates the table's data pages directly instead of marking every row dead one at a time — but it doesn't fire per-row triggers and requires an `ACCESS EXCLUSIVE` lock.

### 3.3 Column Operations

```sql
ALTER TABLE users ADD COLUMN phone TEXT;
ALTER TABLE users ALTER COLUMN phone SET NOT NULL;   -- requires no existing NULLs
ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(100);
ALTER TABLE users DROP COLUMN phone;
ALTER TABLE users ADD COLUMN status TEXT DEFAULT 'active';
```
Since PostgreSQL 11, adding a column with a constant `DEFAULT` is a fast, metadata-only operation (it doesn't rewrite the whole table). Adding `NOT NULL`, changing a column's type, or adding a `DEFAULT` that isn't a simple constant, however, can still require a full table rewrite and an `ACCESS EXCLUSIVE` lock — something to plan around on large, live tables.

### 3.4 Insert Operations (`RETURNING`)

`INSERT ... RETURNING` lets you get back generated values (serial IDs, default timestamps, computed columns) from the same round-trip as the insert, instead of doing a second `SELECT` afterward:
```sql
INSERT INTO users (email, name)
VALUES ('test@test.com', 'John')
RETURNING id, created_at;
```
This is one of PostgreSQL's most-loved features versus databases that require a separate `LAST_INSERT_ID()`-style follow-up query — it removes an entire network round trip and closes a race window where another session could insert between your `INSERT` and your `SELECT`.

### 3.5 Read Operations

```sql
SELECT id, email FROM users WHERE email = 'test@test.com';

SELECT * FROM users
ORDER BY created_at DESC
LIMIT 20 OFFSET 40;              -- pagination (see note below)

SELECT u.id, u.email, o.id AS order_id
FROM users u
JOIN orders o ON o.user_id = u.id
WHERE o.status = 'pending';
```
`OFFSET`-based pagination gets progressively slower on large offsets because Postgres still has to scan and discard all the skipped rows; **keyset pagination** (`WHERE created_at < :last_seen ORDER BY created_at DESC LIMIT 20`) scales far better for deep pagination.

### 3.6 Update

```sql
UPDATE users SET name = 'Jonathan' WHERE id = 42;

UPDATE orders SET status = 'processed'
WHERE status = 'pending'
RETURNING id, customer_id;
```
Recall from Level 1: Postgres never modifies a tuple in place. An `UPDATE` inserts a brand-new tuple version and marks the old one's `xmax`, so an `UPDATE` is really "insert + mark old row dead" under the hood — this is why heavy update workloads need active autovacuum to reclaim space, and why updating *every* column of a wide table costs the same as updating one (a full new tuple is written either way).

### 3.7 Delete

```sql
DELETE FROM users WHERE id = 42;

DELETE FROM sessions
WHERE expires_at < now()
RETURNING id;                     -- audit/log what was removed
```
Like `UPDATE`, a `DELETE` doesn't immediately free space — it marks the tuple's `xmax` so it's invisible to future transactions, and the space is reclaimed later by `VACUUM`. For deleting an entire table's worth of data, `TRUNCATE` (Section 3.2) is far cheaper than a bare `DELETE FROM table`.

---

## 4. Production Usage Scenarios (Real-World Examples)

### Scenario: Database Operations
**Context**: You're standing up a new microservice and need its own isolated database as part of an automated provisioning pipeline.
**The Problem**: Manually creating databases per environment (dev/staging/prod) leads to drift — different encodings, collations, or owners across environments causing subtle bugs that only show up in production.
**The PostgreSQL Solution**: Codify `CREATE DATABASE` (encoding, owner, template) inside migration/provisioning scripts (e.g., run via Terraform, Flyway, or a bootstrap SQL file) so every environment is created identically, always from `template0`, never by hand.

### Scenario: Table Operations
**Context**: You need to add a new feature flag column to a `users` table with 50 million rows in a live, high-traffic system.
**The Problem**: A naive `ALTER TABLE ... ADD COLUMN ... NOT NULL DEFAULT false` on an older Postgres version (or with a non-constant default) can trigger a full table rewrite, taking an `ACCESS EXCLUSIVE` lock and blocking all reads/writes for the duration — unacceptable on a live system.
**The PostgreSQL Solution**: On modern PostgreSQL (11+), adding a column with a constant default is metadata-only and instant. For riskier changes (type changes, adding `NOT NULL` to an existing column), the safe pattern is: add the column nullable first, backfill in small batches, then add the `NOT NULL` constraint using `NOT VALID` + `VALIDATE CONSTRAINT` to avoid a single giant table-locking scan.

### Scenario: Column Operations
**Context**: A `varchar(50)` column for product names turns out to be too short as international product names come in.
**The Problem**: Widening a column type can require Postgres to rewrite the entire table to verify every existing value fits, locking the table for the duration on very large tables.
**The PostgreSQL Solution**: Widening `varchar(n)` to a larger `n` (or to unconstrained `text`) is a fast, metadata-only change in modern Postgres because it only *relaxes* a constraint rather than needing to re-validate data — the team leans on this to make schema evolution non-disruptive, while planning maintenance windows only for genuinely disruptive changes (e.g., `text` → `integer`).

### Scenario: Insert Operations (`RETURNING`)
**Context**: A checkout service inserts a new order and immediately needs the generated order ID to enqueue a fulfillment job.
**The Problem**: Doing `INSERT` followed by a separate `SELECT ... ORDER BY id DESC LIMIT 1` to fetch the new ID is both slower (two round trips) and unsafe under concurrency (another session's insert could sneak in between the two statements, returning the wrong ID).
**The PostgreSQL Solution**: Use `INSERT ... RETURNING id` to get the exact generated ID atomically, in the same statement and round trip, with zero risk of grabbing another transaction's row.

### Scenario: Read Operations
**Context**: A mobile app's activity feed lets users scroll back through months of history.
**The Problem**: `OFFSET`-based pagination gets slower page by page as the offset grows, because the database still has to scan and throw away all preceding rows — by page 500 it's scanning tens of thousands of unwanted rows per request.
**The PostgreSQL Solution**: Switch to keyset (cursor-based) pagination using an indexed `created_at`/`id` column as the seek key (`WHERE (created_at, id) < (:last_created_at, :last_id) ORDER BY created_at DESC, id DESC LIMIT 20`), which uses the index to jump directly to the right spot regardless of how deep the user has scrolled.

### Scenario: Update
**Context**: A background worker needs to safely claim a batch of "pending" jobs from a shared `jobs` table without two workers processing the same job.
**The Problem**: A naive `SELECT` to find pending jobs followed by a separate `UPDATE` to mark them claimed creates a race condition — two workers can both select the same "pending" job before either has updated it.
**The PostgreSQL Solution**: Use a single atomic statement: `UPDATE jobs SET status='claimed', worker_id=:id WHERE status='pending' AND id IN (SELECT id FROM jobs WHERE status='pending' LIMIT 10 FOR UPDATE SKIP LOCKED) RETURNING id, payload;` — `FOR UPDATE SKIP LOCKED` lets concurrent workers each grab a disjoint batch instead of blocking on or double-claiming the same rows.

### Scenario: Delete
**Context**: A GDPR-style data-retention job needs to purge expired user sessions every night, and compliance requires an audit trail of what was deleted.
**The Problem**: A bare `DELETE FROM sessions WHERE expires_at < now();` removes the rows but leaves no record of exactly what was deleted for the audit log, and running it as one giant statement on a huge table can hold locks and bloat the table.
**The PostgreSQL Solution**: Use `DELETE ... RETURNING id, user_id, expires_at` to capture exactly what was removed for the audit log in the same statement, batch the delete (`... WHERE id IN (SELECT id FROM sessions WHERE expires_at < now() LIMIT 5000)`) to avoid one enormous long-held lock, and rely on autovacuum (or a scheduled `VACUUM`) afterward to reclaim the space.

---

## 5. Code & Query Implementation

### 🔹 Basic Implementation
```sql
-- Standard Insert
INSERT INTO users (email, name) VALUES ('test@test.com', 'John');

-- Standard Read
SELECT * FROM users WHERE email = 'test@test.com';

-- Standard Update
UPDATE users SET name = 'Jonathan' WHERE email = 'test@test.com';

-- Standard Delete
DELETE FROM users WHERE email = 'test@test.com';
```

### 🔹 Advanced / Optimized Implementation
```sql
-- Bulk Insert with RETURNING (avoids a 2nd SELECT query)
INSERT INTO logs (level, message)
VALUES
  ('INFO',  'Server started'),
  ('ERROR', 'DB disconnected')
RETURNING id, created_at;

-- Update and return the old/new data for queue processing
UPDATE orders SET status = 'processed'
WHERE status = 'pending'
RETURNING id, customer_id;

-- Upsert: insert or update on conflict (common CRUD pattern)
INSERT INTO users (email, name)
VALUES ('test@test.com', 'John')
ON CONFLICT (email)
DO UPDATE SET name = EXCLUDED.name
RETURNING id;

-- Concurrency-safe job claiming with SKIP LOCKED
UPDATE jobs SET status = 'claimed', worker_id = 7
WHERE id IN (
    SELECT id FROM jobs
    WHERE status = 'pending'
    ORDER BY id
    LIMIT 10
    FOR UPDATE SKIP LOCKED
)
RETURNING id, payload;

-- Batched delete with audit trail
DELETE FROM sessions
WHERE id IN (
    SELECT id FROM sessions
    WHERE expires_at < now()
    LIMIT 5000
)
RETURNING id, user_id, expires_at;
```

---

## 6. Internals & Under the Hood

**How the engine processes this (Parser → Planner → Executor)**
* `INSERT`/`UPDATE`/`DELETE` go through the same Parser → Rewriter → Planner → Executor pipeline as `SELECT` (see Level 1), but the Executor's final step differs: instead of (or in addition to) returning rows to the client, it invokes the storage layer to write new tuples (`INSERT`), write new tuples + mark old ones dead (`UPDATE`), or mark tuples dead (`DELETE`).
* When a `RETURNING` clause is present, the Executor captures the just-written (or just-deleted) tuple's values and streams them back to the client in the same response — no extra scan is needed since the data is already in hand from the write itself.
* `ON CONFLICT` (upsert) is planned as a regular `INSERT` with a fallback path: if a unique/exclusion constraint would be violated, execution branches into the `DO UPDATE`/`DO NOTHING` clause instead of raising an error.

**Storage impact: WAL, Heap, and TOAST**
* Every `INSERT`/`UPDATE`/`DELETE` generates WAL records *before* the corresponding heap page changes are considered durable — this is true even for a single-row `INSERT`, which is why extremely high-frequency single-row inserts benefit from batching (fewer WAL flushes, fewer round trips) versus one insert per statement.
* An `UPDATE` that changes even one column of a wide row still writes an entirely new heap tuple containing every column's value (unless HOT — Heap-Only Tuple — optimization applies, which is possible when no indexed columns changed and the new tuple fits on the same page, avoiding index updates entirely).
* Large values (long `text`, big `jsonb` payloads) inserted or updated may be routed through **TOAST**, transparently compressing and/or storing them out-of-line — this affects the actual amount of WAL/heap I/O generated by what looks like "one row's" insert or update.

---

## 7. Performance & Benchmarking

**`EXPLAIN ANALYZE` impacts**
```sql
EXPLAIN ANALYZE
UPDATE orders SET status = 'processed' WHERE status = 'pending';
```
* For DML, `EXPLAIN ANALYZE` actually performs the write (wrap in a transaction and `ROLLBACK` if you just want to measure, not commit, the effect).
* Look for the scan strategy feeding the `Update`/`Delete` node — a `Seq Scan` filtering `WHERE status = 'pending'` on a huge `orders` table means every row is being read to find the small pending subset; a partial index (`CREATE INDEX ON orders (id) WHERE status = 'pending'`) can turn that into an efficient index scan.

**Memory vs. disk trade-offs**
* Bulk `INSERT`s benefit from `maintenance_work_mem` when any indexes need updating and from batching multiple rows into a single `INSERT ... VALUES (...), (...), (...)` statement rather than one statement per row, since each round trip and each WAL flush has fixed overhead.
* Large batched `DELETE`/`UPDATE` statements risk long transaction times and heavier `work_mem`/lock footprint; chunking them (as in the `LIMIT 5000` pattern above) keeps individual transactions small, avoids long lock waits for other sessions, and keeps replication lag on standbys lower.
* Frequent updates to the same rows generate frequent dead tuples — autovacuum working memory (`autovacuum_work_mem`, falling back to `maintenance_work_mem`) and its scheduling thresholds directly affect how quickly that space is reclaimed versus how much table bloat accumulates.

---

## 8. Best Practices & Common Mistakes

* ✅ **Do**: Use `RETURNING` whenever you need the result of a write — it removes an entire extra query and a race-condition window.
* ✅ **Do**: Use `FOR UPDATE SKIP LOCKED` for any "claim a job from a shared queue" pattern instead of hand-rolled locking.
* ✅ **Do**: Batch large `UPDATE`/`DELETE` operations into chunks rather than one giant statement, on live production tables.
* ❌ **Don't**: Rely on a separate `SELECT` after `INSERT` to fetch a generated ID — it's slower and race-prone versus `RETURNING`.
* ❌ **Don't**: Ignore index overhead or transaction locking — every index on a table adds write cost to every `INSERT`/`UPDATE`/`DELETE`, and a long-running `UPDATE`/`DELETE` transaction holds row locks that can block other sessions.
* ⚠️ **Common Mistake**: N+1 queries — issuing an `INSERT`/`UPDATE` per item in a loop from the application, instead of one batched, multi-row statement.
* ⚠️ **Common Mistake**: Missing indexes on foreign key columns, causing slow `DELETE`s/`UPDATE`s on the "one" side of a relationship (Postgres must scan the child table for rows to cascade against).
* ⚠️ **Common Mistake**: Assuming `DELETE`/`UPDATE` immediately frees disk space — dead tuples linger until vacuumed, so heavy churn without healthy autovacuum leads to table and index bloat.

---

## 9. Interview Questions

1. **Beginner**: Explain the four CRUD operations in PostgreSQL to a junior dev, and what the `RETURNING` clause adds on top of standard SQL `INSERT`/`UPDATE`/`DELETE`.
2. **Beginner**: What's the difference between `DELETE FROM table` and `TRUNCATE TABLE`, and when would you pick one over the other?
3. **Intermediate**: How would you optimize a query workload with deep `OFFSET`-based pagination, and what would you replace it with?
4. **Intermediate**: How does `ON CONFLICT ... DO UPDATE` (upsert) work internally, and how is it different from a manual "check if exists, then insert or update" application-level pattern?
5. **Intermediate**: Why might adding a `NOT NULL` column with a non-constant default to a huge live table be dangerous, and how would you do it safely?
6. **Expert**: Describe the low-level locking behavior when two concurrent transactions run `UPDATE` on the same row versus different rows — what row-level and page-level locks are taken, and how does this differ under `READ COMMITTED` versus `SERIALIZABLE`?
7. **Expert**: Explain Heap-Only Tuple (HOT) updates — what conditions must hold for an `UPDATE` to qualify, and why does this matter for index maintenance overhead in high-concurrency environments?
8. **Expert**: Walk through why `FOR UPDATE SKIP LOCKED` is safe for concurrent job-queue claiming while a plain `FOR UPDATE` (without `SKIP LOCKED`) would cause workers to serialize on each other, and what visibility/locking mechanics make the difference.

---

*Primary sources: PostgreSQL Official Documentation (postgresql.org/docs) — chapters on Data Definition, Data Manipulation, DML, and the `INSERT`/`UPDATE`/`DELETE` reference pages.*