# Level 1 — PostgreSQL Fundamentals

*(Based on official PostgreSQL documentation, current stable series: PostgreSQL 18. PostgreSQL 19 is in beta as of mid‑2026.)*

## 1. Learning Objectives

* **What you'll learn**: The core mechanics of PostgreSQL — its architecture, how it compares to MySQL and MongoDB, the ACID and CAP models, how to install and connect to it, the building blocks of a database (schemas, tables, tuples, relations), and the full native data type system (numeric, character, boolean, date/time, UUID, JSONB, arrays).
* **Why it matters**: Every higher-level PostgreSQL skill — indexing, query tuning, replication, partitioning — sits on top of these fundamentals. Getting the mental model of "how Postgres actually stores and processes a row" right at this stage prevents a huge class of bugs and performance problems later (wrong data types, N+1 queries, broken transaction assumptions, schema designs that don't scale).

---

## 2. Topics Covered

* Introduction (Architecture, vs MySQL/MongoDB, ACID, CAP)
* Installation (Docker, pgAdmin, psql)
* Database Basics (Schema, Tables, Tuples, Relations)
* Data Types (Numeric, Character, Boolean, Date/Time, UUID, JSONB, Arrays)

---

## 3. Deep Dive: Concepts

### 3.1 Introduction — Architecture, vs MySQL/MongoDB, ACID, CAP

**Process architecture.** PostgreSQL uses a **process-per-connection** model (not threads). When a client connects, the `postmaster` (the main server process) forks a dedicated **backend process** to serve that connection for its lifetime. This is different from MySQL, which historically used a similar process/thread-per-connection model but has since layered thread-pooling on top in some configurations, and very different from MongoDB, which is a single multi-threaded `mongod` process handling many connections internally.

Key server-side processes you'll see in `ps aux | grep postgres`:
* `postmaster` — the supervisor process, listens for connections.
* **backend processes** — one per client connection, executes queries.
* `checkpointer` — periodically flushes dirty shared-buffer pages to disk.
* `background writer` — writes dirty buffers gradually to avoid I/O spikes.
* `WAL writer` — flushes Write-Ahead Log records to disk.
* `autovacuum launcher/workers` — reclaim space from dead tuples (see MVCC below).
* `stats collector` / `logical replication launcher` — auxiliary processes.

**Memory architecture**: `shared_buffers` (Postgres's own page cache, distinct from OS cache), `work_mem` (per-operation sort/hash memory), `maintenance_work_mem` (for VACUUM, CREATE INDEX), and `wal_buffers`.

**PostgreSQL vs MySQL vs MongoDB**

| Aspect | PostgreSQL | MySQL | MongoDB |
|---|---|---|---|
| Data model | Relational, strongly typed, extensible (custom types, extensions) | Relational | Document (BSON) |
| Concurrency control | MVCC (multiversion concurrency control) via row versions | MVCC in InnoDB engine | MVCC (WiredTiger engine) |
| Transactions | Full ACID, multi-statement, multi-table by default | ACID with InnoDB (not with MyISAM) | ACID within/across documents since 4.0, multi-document since 4.0/6.0 for sharded clusters |
| Joins | Native, deeply optimized | Native | Manual (`$lookup`) — not the primary access pattern |
| Schema | Rigid, enforced at the DB level | Rigid | Flexible / schema-less by default |
| Extensibility | Custom types, operators, index types (GIN/GiST/BRIN), extensions (PostGIS, pg_trgm, pgvector) | Limited | Limited |
| Best fit | Complex queries, data integrity-critical systems, analytics + OLTP hybrid | Simple, high-read web apps, well-understood ops tooling | Rapidly evolving schemas, document-shaped data, horizontal scale-out |

**ACID**
* **Atomicity** — a transaction's operations either all happen or none do (`BEGIN` / `COMMIT` / `ROLLBACK`).
* **Consistency** — constraints (`CHECK`, `FOREIGN KEY`, `UNIQUE`) are never violated at transaction boundaries.
* **Isolation** — concurrent transactions don't see each other's uncommitted changes. PostgreSQL implements this via **MVCC**: every row version (tuple) carries `xmin`/`xmax` transaction IDs, so readers never block writers and writers never block readers.
* **Durability** — once committed, data survives a crash, guaranteed by the **Write-Ahead Log (WAL)**: changes are written to the WAL and `fsync`'d before the transaction is acknowledged as committed.

Isolation levels supported: `READ UNCOMMITTED` (treated as `READ COMMITTED`), `READ COMMITTED` (default), `REPEATABLE READ`, `SERIALIZABLE` (true serializable snapshot isolation, SSI).

**CAP theorem context.** CAP (Consistency, Availability, Partition tolerance — pick 2 of 3 under a network partition) is primarily a *distributed systems* framework. A single PostgreSQL primary is a CA system in the classic sense (it doesn't tolerate partitions gracefully — a network split just isolates it). When you add **streaming replication**, PostgreSQL becomes tunable: synchronous replication favors consistency (a commit waits for a standby to confirm) over availability; asynchronous replication favors availability over strict consistency (a standby may lag). MongoDB, as a sharded, natively distributed system, is usually discussed as CP with tunable consistency via write/read concerns. This is why "Postgres vs Mongo" is often really "vertical, strongly-consistent relational engine" vs "horizontally distributed, tunably-consistent document engine."

### 3.2 Installation — Docker, pgAdmin, psql

**Docker** (fastest path for local dev):
```bash
docker run --name pg-dev \
  -e POSTGRES_PASSWORD=devpass \
  -e POSTGRES_DB=fin_app \
  -p 5432:5432 \
  -d postgres:18
```
This pulls the official image (maintained per the Docker Official Images program) and starts a server with a default `postgres` superuser.

**psql** (the official command-line client):
```bash
psql -h localhost -U postgres -d fin_app
```
Useful meta-commands: `\l` (list databases), `\dt` (list tables), `\d tablename` (describe table), `\du` (list roles), `\x` (expanded output toggle), `\timing` (show query timing).

**pgAdmin** is the official graphical administration tool, available as a desktop app or a web app (also dockerizable) for browsing schemas, running queries, viewing query plans, and managing roles/backups through a UI rather than the CLI.

### 3.3 Database Basics — Schema, Tables, Tuples, Relations

* **Cluster** — one running Postgres server instance manages a *database cluster*: a collection of databases sharing the same set of background processes and one data directory.
* **Database** — an isolated namespace of schemas; you cannot query across databases within a single SQL statement (without `dblink`/`postgres_fdw`).
* **Schema** — a namespace *inside* a database that holds tables, views, functions, etc. Every database has a default `public` schema. Schemas enable **multi-tenancy** (one schema per tenant) and namespace collision avoidance.
* **Table (relation)** — in Postgres terminology, "relation" is the formal term covering tables, views, indexes, and sequences; informally "table" and "relation" are used interchangeably for actual data tables.
* **Tuple** — Postgres's internal term for a physical row version on disk (the *row* is the logical concept; a *tuple* is one physical, versioned instance of it under MVCC — an `UPDATE` creates a new tuple rather than overwriting in place).
* **Relations between tables** — expressed via `PRIMARY KEY`, `FOREIGN KEY`, and join tables for many-to-many relationships. Referential integrity is enforced by the engine itself, not application code.

### 3.4 Data Types

PostgreSQL's type system is one of its biggest differentiators — it's strict (no silent implicit type coercion the way MySQL historically allowed) and extensible.

| Category | Types | Notes |
|---|---|---|
| Numeric | `smallint`, `integer`, `bigint`, `decimal`/`numeric`, `real`, `double precision`, `serial`/`bigserial` | `numeric` is arbitrary-precision (exact) — use for money; `real`/`double precision` are IEEE 754 floats (inexact). `serial` is sugar for an integer column + owned sequence. |
| Character | `char(n)`, `varchar(n)`, `text` | `text` has no length limit and is the idiomatic default; `varchar(n)` enforces a limit; `char(n)` pads with spaces. Performance is essentially identical between them internally. |
| Boolean | `boolean` | Stores `TRUE`/`FALSE`/`NULL`; literals `'t'`, `'true'`, `'yes'`, `'1'` are all accepted on input. |
| Date/Time | `date`, `time`, `timestamp`, `timestamptz`, `interval` | `timestamptz` (timestamp with time zone) is almost always the right default for application data — it stores UTC internally and converts on display based on session `TimeZone`. |
| UUID | `uuid` | 128-bit identifier, commonly generated via `gen_random_uuid()` (built-in since PG13) or the `uuid-ossp` extension. Good for distributed ID generation without coordination. |
| JSONB | `json`, `jsonb` | `json` stores an exact text copy; `jsonb` stores a decomposed binary format — slightly slower to input, much faster to query, and supports indexing (GIN). `jsonb` is recommended for almost all use cases. |
| Arrays | `integer[]`, `text[]`, etc. | Any base type can be made into a (possibly multi-dimensional) array column, queryable with array operators (`@>`, `&&`, `ANY`). |

---

## 4. Production Usage Scenarios (Real-World Examples)

### Scenario: Introduction (Architecture, vs MySQL/MongoDB, ACID, CAP)
**Context**: You're building a fintech ledger service that must process concurrent transfers between accounts without ever losing money or double-crediting.
**The Problem**: Under concurrent load, two transfers touching the same account could race — read-modify-write cycles can clobber each other, and a crash mid-transfer could leave the ledger inconsistent.
**The PostgreSQL Solution**: Wrap each transfer in a single transaction using `SELECT ... FOR UPDATE` to lock the affected account rows, rely on MVCC so concurrent readers (e.g., a reporting dashboard) never block or get blocked by the writers, and trust WAL-backed durability so a crash mid-commit never leaves a half-applied transfer. `SERIALIZABLE` isolation is used for the transfer logic specifically to eliminate write-skew anomalies.

### Scenario: Installation (Docker, pgAdmin, psql)
**Context**: Your team needs identical local dev environments and a CI pipeline that spins up a throwaway database per test run.
**The Problem**: "Works on my machine" version drift between developers, and slow, manual DB provisioning for CI.
**The PostgreSQL Solution**: Pin an exact image tag (`postgres:18.4`) in `docker-compose.yml` so every developer and CI runner gets byte-identical server behavior; use `psql` scripts in CI to apply migrations and seed data non-interactively; use pgAdmin only for humans doing ad hoc investigation, never for anything that needs to be reproducible.

### Scenario: Database Basics (Schema, Tables, Tuples, Relations)
**Context**: You're building a SaaS product where each customer's data must be logically isolated, and you have hundreds of customers.
**The Problem**: A single shared table with a `tenant_id` column risks a forgotten `WHERE tenant_id = ?` clause leaking one customer's data to another — a serious breach.
**The PostgreSQL Solution**: Use one **schema per tenant** (`tenant_a`, `tenant_b`, ...) with identical table structures, and set `search_path` per connection to the correct tenant schema. Queries naturally resolve to the right tenant's tables without a `tenant_id` filter to forget. For very large tenant counts, this is combined with **Row-Level Security (RLS)** as a second line of defense.

### Scenario: Data Types (Numeric, Character, Boolean, Date/Time, UUID, JSONB, Arrays)
**Context**: You're storing monetary amounts, flexible product attributes, and globally-unique order IDs generated by multiple independent services.
**The Problem**: Using `float`/`double precision` for money introduces rounding errors that compound over millions of transactions; a rigid product-attributes table can't keep up with catalog changes; auto-incrementing IDs collide across services generating IDs independently.
**The PostgreSQL Solution**: Store all currency amounts as `numeric(12,2)` for exact decimal arithmetic; store variable product attributes in a `jsonb` column (indexed with a GIN index for fast attribute lookups) instead of a rigid column-per-attribute schema; use `uuid` primary keys (`gen_random_uuid()`) for anything generated by multiple services so IDs never collide without a central coordinator.

---

## 5. Code & Query Implementation

### 🔹 Basic Implementation
```sql
-- Connect and check version
SELECT version();

-- Create a new database
CREATE DATABASE fin_app;

-- Connect to it, then create a simple table
\c fin_app

CREATE TABLE accounts (
    id          BIGSERIAL PRIMARY KEY,
    owner_name  TEXT NOT NULL,
    balance     NUMERIC(12,2) NOT NULL DEFAULT 0,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO accounts (owner_name, balance) VALUES ('Asha Rao', 1000.00);

SELECT * FROM accounts;
```

### 🔹 Advanced / Optimized Implementation
```sql
-- Create schema for tenant isolation
CREATE SCHEMA IF NOT EXISTS tenant_a;

-- Set search path to prioritize tenant schema
SET search_path TO tenant_a, public;

-- A table using UUID, JSONB, and arrays together
CREATE TABLE tenant_a.orders (
    order_id     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id   BIGINT NOT NULL REFERENCES public.accounts(id),
    line_items   JSONB NOT NULL,          -- flexible order payload
    tags         TEXT[] DEFAULT '{}',      -- e.g. {'priority','gift-wrap'}
    placed_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- GIN index so jsonb attribute lookups are fast
CREATE INDEX idx_orders_line_items ON tenant_a.orders USING GIN (line_items);

-- Atomic, isolation-safe transfer between two accounts
BEGIN;
  SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
  SELECT * FROM accounts WHERE id = 2 FOR UPDATE;

  UPDATE accounts SET balance = balance - 100 WHERE id = 1;
  UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;

-- Query orders tagged 'priority' whose line_items contain a given SKU
SELECT order_id
FROM tenant_a.orders
WHERE tags @> ARRAY['priority']
  AND line_items @> '{"sku": "ABC123"}';
```

---

## 6. Internals & Under the Hood

**How the engine processes a query — Parser → Rewriter → Planner → Executor**
1. **Parser** — checks SQL syntax and produces a raw parse tree.
2. **Rewriter** — applies rules (e.g., expands views into their underlying queries).
3. **Planner/Optimizer** — generates possible query plans (sequential scan, index scan, various join strategies) and picks the cheapest one based on table statistics (`pg_statistic`, refreshed by `ANALYZE`).
4. **Executor** — walks the chosen plan tree bottom-up, pulling rows through nodes (scan → join → sort → aggregate) using a Volcano-style iterator model.

**Storage internals**
* **Heap** — the default table storage: an unordered collection of fixed-size pages (default 8 KB) containing tuples. There is no clustering by primary key by default (unlike MySQL/InnoDB, where the table *is* the primary key's B-tree).
* **MVCC tuples** — every tuple has hidden `xmin` (creating transaction ID) and `xmax` (deleting/superseding transaction ID) system columns. An `UPDATE` is implemented as an `INSERT` of a new tuple + marking the old tuple's `xmax`, which is why **VACUUM** is needed to reclaim dead tuple space.
* **WAL (Write-Ahead Log)** — every modification is logged sequentially before the corresponding heap page is flushed, enabling crash recovery and forming the basis for streaming replication.
* **TOAST** (The Oversized-Attribute Storage Technique) — values too large for a page (e.g., a big `jsonb` blob or long `text`) are automatically compressed and/or moved to a separate TOAST table, transparent to queries.

---

## 7. Performance & Benchmarking

**`EXPLAIN` / `EXPLAIN ANALYZE`**
```sql
EXPLAIN ANALYZE
SELECT * FROM accounts WHERE owner_name = 'Asha Rao';
```
* `EXPLAIN` alone shows the *planned* execution strategy and estimated cost without running the query.
* `EXPLAIN ANALYZE` actually **executes** the query and adds real timing/row-count data per plan node — essential for spotting cases where the planner's row estimates are wrong (a common root cause of slow queries).
* Watch for `Seq Scan` on large tables where an `Index Scan` was expected — usually means a missing index or stale statistics (run `ANALYZE tablename;`).

**Memory vs. disk trade-offs**
* `shared_buffers` too small → frequent disk reads for hot data; too large → starves OS cache and other processes. A common starting guideline is ~25% of system RAM, tuned from there.
* `work_mem` too small → sorts and hash joins spill to disk (`temp_files`, visible in logs), which is dramatically slower; too large, multiplied across many concurrent connections doing sorts, can exhaust server memory.
* Sequential scans are disk-throughput-bound; index scans on cold data are disk-latency-bound (random I/O) — on spinning disks this is a much bigger penalty than on SSD/NVMe, which is part of why `random_page_cost` is tunable.

---

## 8. Best Practices & Common Mistakes

* ✅ **Do**: Use `timestamptz` instead of `timestamp` for anything user-facing or cross-timezone.
* ✅ **Do**: Default to `jsonb` over `json` unless you specifically need to preserve exact input formatting/whitespace.
* ✅ **Do**: Run `ANALYZE` (or rely on autovacuum's analyze) after large bulk loads so the planner has fresh statistics.
* ✅ **Do**: Index foreign key columns explicitly — Postgres does **not** create an index on a `FOREIGN KEY` column automatically (unlike the referenced primary key side).
* ❌ **Don't**: Use `float`/`double precision` for money or any value requiring exact decimal arithmetic.
* ❌ **Don't**: Ignore index write overhead — every index speeds up reads but slows down every `INSERT`/`UPDATE`/`DELETE` on that table, and consumes disk/cache space.
* ❌ **Don't**: Leave long-running idle transactions open — they hold back the point up to which VACUUM can clean dead tuples, causing table/index bloat.
* ⚠️ **Common Mistake**: N+1 queries — issuing one query per row of a parent result instead of a single `JOIN`. Fine in isolation, catastrophic under load.
* ⚠️ **Common Mistake**: Missing indexes on foreign keys, leading to slow cascading deletes and slow joins in the "many" direction of a one-to-many relationship.
* ⚠️ **Common Mistake**: Treating `SERIAL`/`BIGSERIAL` sequences as gap-free — they are **not** transactional-safe against gaps (a rolled-back `INSERT` still consumes a sequence value).

---

## 9. Interview Questions

1. **Beginner**: Explain PostgreSQL's process-per-connection architecture to a junior developer, and why it means one client's crash doesn't take down the whole server.
2. **Beginner**: What's the difference between `char(n)`, `varchar(n)`, and `text`, and which should you default to?
3. **Intermediate**: Why does Postgres need `VACUUM`, and what would happen if autovacuum were disabled on a heavily-updated table?
4. **Intermediate**: How would you design a multi-tenant schema in Postgres, and what are the trade-offs between schema-per-tenant and a shared table with a `tenant_id` column plus Row-Level Security?
5. **Intermediate**: A query that used to use an index scan suddenly switches to a sequential scan after a large data load — what's your diagnostic process?
6. **Expert**: Describe what happens internally, at the tuple level, when you run an `UPDATE` on a row that's part of a long-running concurrent transaction — how do `xmin`/`xmax` and MVCC visibility rules interact with a second transaction's snapshot?
7. **Expert**: Explain how `SELECT ... FOR UPDATE` and `SERIALIZABLE` isolation each prevent write-skew anomalies, and describe a concrete scenario where `FOR UPDATE` alone is insufficient but `SERIALIZABLE` is.
8. **Expert**: Walk through the low-level locking behavior when two concurrent transactions attempt to update the same row versus two different rows in the same page — what locks are taken, at what granularity, and what determines who waits versus who proceeds?

---

*Primary sources: PostgreSQL Official Documentation (postgresql.org/docs), PostgreSQL Versioning Policy, and PostgreSQL release notes (postgresql.org/docs/release).*