# Level 4 — Querying

*(Based on official PostgreSQL documentation, current stable series: PostgreSQL 18.)*

## 1. Learning Objectives

* **What you'll learn**: How to filter rows precisely (`WHERE`, `AND`/`OR`, `BETWEEN`, `IN`, `LIKE`/`ILIKE`, `EXISTS`), control result ordering (`ORDER BY`, `NULLS FIRST/LAST`), and aggregate data (`GROUP BY`, `HAVING`) — plus how each of these is actually executed internally, so you can predict performance instead of guessing.
* **Why it matters**: Filtering, sorting, and grouping are the operations that turn raw storage into the answers your application actually needs. Getting the *logic* right is only half the job — getting the *shape* of the query right (which filter form, which subquery pattern) is what determines whether a query takes 2 milliseconds or 20 seconds at production scale.

---

## 2. Topics Covered

* Filtering (`WHERE`, `AND`, `OR`, `BETWEEN`, `IN`, `LIKE`, `ILIKE`, `EXISTS`)
* Sorting (`ORDER BY`, `NULLS FIRST`/`LAST`)
* Grouping (`GROUP BY`, `HAVING`)

---

## 3. Query Execution Overview (Diagram)

```mermaid
graph LR
    A[FROM / JOIN] --> B[WHERE]
    B --> C[GROUP BY]
    C --> D[HAVING]
    D --> E[SELECT]
    E --> F[ORDER BY]
    F --> G[LIMIT / OFFSET]

    style A fill:#2b6cb0,color:#fff
    style B fill:#4299e1,color:#fff
    style C fill:#4299e1,color:#fff
    style D fill:#4299e1,color:#fff
    style E fill:#4299e1,color:#fff
    style F fill:#4299e1,color:#fff
    style G fill:#4299e1,color:#fff
```

This is PostgreSQL's **logical** clause evaluation order — notably *different* from the order you type the clauses in SQL. This is why `HAVING` can filter on an aggregate that `WHERE` cannot (rows haven't been grouped yet when `WHERE` runs), and why you can't reference a `SELECT`-list column alias inside `WHERE` (the alias doesn't exist yet at that stage — though Postgres *does* allow it in `GROUP BY`/`ORDER BY`, which run after `SELECT`).

---

## 4. Deep Dive: Filtering (WHERE)

### 4.1 Basic boolean logic — `AND` / `OR`

```sql
SELECT * FROM customers WHERE status = 'active' AND age > 18;

SELECT * FROM customers WHERE status = 'active' OR status = 'trial';
```
`AND` binds tighter than `OR` — always parenthesize mixed conditions explicitly (`WHERE (status = 'active' OR status = 'trial') AND age > 18`) rather than relying on precedence rules that are easy to misremember and easy to misread months later.

### 4.2 `BETWEEN`

```sql
SELECT * FROM orders WHERE placed_at BETWEEN '2026-01-01' AND '2026-01-31';
```
`BETWEEN x AND y` is inclusive on both ends (`>= x AND <= y`). For half-open date ranges (a very common real bug source — "January" accidentally including part of February), prefer explicit `>=` / `<` bounds: `WHERE placed_at >= '2026-01-01' AND placed_at < '2026-02-01'`.

### 4.3 `IN`

```sql
SELECT * FROM orders WHERE status IN ('pending', 'processing', 'shipped');

-- IN with a subquery
SELECT * FROM customers WHERE id IN (SELECT customer_id FROM orders WHERE total > 1000);
```
A literal-list `IN` is just sugar for chained `OR`s and is fine at any scale. An `IN (subquery)` is where care is needed — see the `IN` vs `EXISTS` comparison below.

### 4.4 `LIKE` / `ILIKE`

```sql
SELECT * FROM products WHERE name LIKE 'Pro%';      -- case-sensitive
SELECT * FROM products WHERE name ILIKE 'pro%';      -- case-insensitive
```
`%` matches any sequence of characters, `_` matches exactly one. A **leading** wildcard (`LIKE '%phone%'`) cannot use a standard B-tree index and forces a sequential scan; a **trailing**-only wildcard (`LIKE 'iphone%'`) *can* use a standard B-tree index. For real leading-wildcard/full-text search at scale, use a `pg_trgm` GIN index or PostgreSQL's built-in full-text search (`tsvector`/`tsquery`) instead of `LIKE '%...%'`.

### 4.5 `EXISTS`

```sql
SELECT id, name
FROM customers c
WHERE EXISTS (
    SELECT 1 FROM orders o
    WHERE o.customer_id = c.id AND o.total > 1000
);
```
`EXISTS` only cares whether the subquery returns **at least one row** — the planner can stop scanning the instant it finds a match, and doesn't care what columns the subquery selects (`SELECT 1` is idiomatic precisely to signal "we don't care about the value").

### 4.6 `IN` vs `EXISTS` — Comparison Table

| Aspect | `IN (subquery)` | `EXISTS (subquery)` |
|---|---|---|
| Stops early on first match? | Historically no (materializes full list); modern planner can transform simple cases into a semi-join, but not guaranteed | Yes — semantically a semi-join, short-circuits on first match |
| NULL handling | Dangerous: `NOT IN` with any `NULL` in the subquery result silently returns zero rows | Safe: `NOT EXISTS` handles `NULL`s correctly regardless |
| Best for | Small, fixed, or literal lists; simple subqueries the planner can optimize well | Correlated subqueries checking "does at least one related row exist" |
| Correlated subquery support | Supported, but less idiomatic | Natural fit — designed for correlated existence checks |
| Typical production advice | Use for literal lists (`IN ('a','b','c')`) | Prefer for "does a related row exist" checks, and **always** prefer `NOT EXISTS` over `NOT IN` |

**Golden rule**: `WHERE id NOT IN (SELECT customer_id FROM orders)` is a classic production bug — if even one row in the subquery has a `NULL` `customer_id`, the entire `NOT IN` returns **zero rows**, silently, with no error. `NOT EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id)` doesn't have this trap.

---

## 5. Deep Dive: Sorting

### 5.1 `ORDER BY`

```sql
SELECT * FROM orders ORDER BY placed_at DESC;

SELECT * FROM orders ORDER BY status ASC, placed_at DESC;   -- multi-column sort
```

### 5.2 `NULLS FIRST` / `NULLS LAST`

```sql
SELECT * FROM customers ORDER BY last_login_at DESC NULLS LAST;
```

| Sort direction | Postgres default NULL placement | Why it matters |
|---|---|---|
| `ASC` (default) | `NULLS LAST` | NULLs are treated as "larger than any value" by default in ascending order |
| `DESC` | `NULLS FIRST` | NULLs are treated as "larger than any value," so they sort first when descending |

This default surprises people constantly — "customers who never logged in" (`last_login_at IS NULL`) will show up at the **top** of a naive `ORDER BY last_login_at DESC` (most-recent-first) query unless you explicitly add `NULLS LAST`.

**Performance note**: `ORDER BY` on an indexed column in the same direction as the index (or its exact reverse) can be satisfied directly by an **Index Scan**, avoiding a separate sort step entirely. `ORDER BY` on an unindexed column (or a computed expression) requires an explicit **Sort** node, which uses `work_mem` and spills to disk (`external sort` in `EXPLAIN ANALYZE` output) if the result set is larger than `work_mem` allows.

---

## 6. Deep Dive: Grouping

### 6.1 `GROUP BY`

```sql
SELECT status, COUNT(*) AS order_count, SUM(total) AS revenue
FROM orders
GROUP BY status;
```

### 6.2 `HAVING`

```sql
SELECT customer_id, SUM(total) AS lifetime_value
FROM orders
GROUP BY customer_id
HAVING SUM(total) > 5000;
```

### 6.3 `WHERE` vs `HAVING` — Comparison Table

| Aspect | `WHERE` | `HAVING` |
|---|---|---|
| Runs relative to `GROUP BY` | **Before** grouping (filters individual rows) | **After** grouping (filters groups) |
| Can reference aggregate functions (`SUM`, `COUNT`, ...)? | No | Yes |
| Performance impact | Reduces rows *before* the (expensive) grouping/aggregation step — always prefer when the condition doesn't need an aggregate | Necessarily processes the aggregation for all groups before discarding any — cannot skip ungrouped rows early |
| Typical mistake | Using `HAVING` for a plain per-row condition that `WHERE` could have handled, wasting effort by aggregating rows that get discarded afterward | Using `WHERE` and expecting it to filter on `SUM()`/`COUNT()` — Postgres will raise a syntax error since aggregates don't exist yet at `WHERE`'s evaluation stage |

**Rule of thumb**: filter as much as possible in `WHERE` (row-level, pre-aggregation) and reserve `HAVING` strictly for conditions that genuinely require the aggregated value.

---

## 7. Production Usage Scenarios (Real-World Examples)

### Scenario: Filtering (WHERE, AND, OR, BETWEEN, IN, LIKE, ILIKE, EXISTS)
**Context**: A customer support dashboard needs to find all high-value customers who have placed at least one large order, out of tens of millions of customer rows, refreshed on every dashboard load.
**The Problem**: A naive correlated `IN (SELECT customer_id FROM orders WHERE total > 1000)` on a huge `orders` table can force the planner into materializing a very large intermediate result, and a careless `NOT IN` variant elsewhere in the same dashboard silently returns zero rows the moment a `NULL` sneaks into the subquery.
**The PostgreSQL Solution**: Rewrite the existence check as `WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id AND o.total > 1000)` so the planner can short-circuit per customer row instead of building a full match list, back it with an index on `orders(customer_id, total)`, and standardize on `NOT EXISTS` instead of `NOT IN` anywhere the team needs a "no matching row" check to avoid the `NULL` trap entirely.

### Scenario: Sorting (ORDER BY, NULLS FIRST/LAST)
**Context**: An admin panel lists customers by most-recent activity, and support agents specifically want customers who've *never* logged in surfaced at the bottom of the list, not the top.
**The Problem**: `ORDER BY last_login_at DESC` on its own puts `NULL` (never-logged-in) customers at the very top by Postgres's default `NULLS FIRST` behavior for descending sorts — the opposite of what the support team actually wants, and a subtle bug that's easy to ship without noticing in a small test dataset.
**The PostgreSQL Solution**: Explicitly specify `ORDER BY last_login_at DESC NULLS LAST`, and back the query with a descending index (`CREATE INDEX ON customers (last_login_at DESC NULLS LAST)`) so the correct ordering is served directly by an index scan instead of requiring a separate sort step on every page load.

### Scenario: Grouping (GROUP BY, HAVING)
**Context**: A finance team needs a daily report of which customers exceeded $5,000 in lifetime order value, computed against a fast-growing `orders` table with hundreds of millions of rows.
**The Problem**: Computing `SUM(total)` per customer across the *entire* orders table every time the report runs, then filtering with `HAVING`, forces Postgres to aggregate every single customer's full order history — including the vast majority of customers who will never come close to the $5,000 threshold — before discarding them.
**The PostgreSQL Solution**: Push any row-level filtering possible into `WHERE` first (e.g., restrict to the current fiscal year with `WHERE placed_at >= '2026-01-01'` before grouping), keep `HAVING SUM(total) > 5000` strictly for the aggregate-level threshold that genuinely requires the sum, and back the query with an index on `orders(customer_id, placed_at)` so the pre-aggregation filtering itself is fast.

---

## 8. Code & Query Implementation

### 🔹 Basic Implementation
```sql
-- Standard filtering
SELECT * FROM customers WHERE status = 'active' AND age > 18;

-- BETWEEN, IN, LIKE
SELECT * FROM orders WHERE placed_at BETWEEN '2026-01-01' AND '2026-01-31';
SELECT * FROM orders WHERE status IN ('pending', 'shipped');
SELECT * FROM products WHERE name ILIKE 'pro%';

-- Sorting with NULL handling
SELECT * FROM customers ORDER BY last_login_at DESC NULLS LAST;

-- Grouping with HAVING
SELECT status, COUNT(*) FROM orders GROUP BY status HAVING COUNT(*) > 100;
```

### 🔹 Advanced / Optimized Implementation
```sql
-- EXISTS vs IN for high-performance subqueries
-- This stops scanning the 'orders' table the millisecond it finds 1 match.
SELECT id, name
FROM customers c
WHERE EXISTS (
    SELECT 1
    FROM orders o
    WHERE o.customer_id = c.id
    AND o.total > 1000
);

-- Safe NOT EXISTS instead of the NULL-unsafe NOT IN
SELECT c.id, c.name
FROM customers c
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.customer_id = c.id
);

-- Trigram-indexed search to make a leading-wildcard LIKE fast
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_products_name_trgm ON products USING GIN (name gin_trgm_ops);
SELECT * FROM products WHERE name ILIKE '%phone%';   -- now index-assisted

-- Push filtering before aggregation, reserve HAVING for aggregate-only conditions
SELECT customer_id, SUM(total) AS lifetime_value
FROM orders
WHERE placed_at >= '2026-01-01'
GROUP BY customer_id
HAVING SUM(total) > 5000
ORDER BY lifetime_value DESC;
```

---

## 9. Internals & Under the Hood

**How PostgreSQL engine processes this (Parser → Planner → Executor)**
* **Parser**: validates the SQL and builds a parse tree reflecting the clauses as written.
* **Planner**: reorders operations into an efficient physical plan following logical evaluation order (`FROM`/`JOIN` → `WHERE` → `GROUP BY` → `HAVING` → `SELECT` → `ORDER BY` → `LIMIT`), choosing scan methods (Seq Scan vs Index Scan vs Bitmap Index Scan) for `WHERE` predicates, a `HashAggregate` or `GroupAggregate` strategy for `GROUP BY`, and a `Sort` node (or an index-satisfied ordering) for `ORDER BY` — all based on table statistics from `pg_statistic`.
* **Executor**: pulls rows bottom-up through the chosen plan tree; for `EXISTS`, the executor uses a **semi-join** strategy that literally stops evaluating the inner side after the first matching row per outer row.

**Storage impact: WAL, Heap, and TOAST**
* Pure `SELECT` queries (filtering, sorting, grouping) are read-only and generate **no WAL** — WAL is a write-path concept. The exception is when a `SELECT` triggers automatic actions like hint-bit updates on tuples (a minor internal optimization, not user-visible WAL of consequence) or when statistics are updated by autovacuum's analyze phase as a side effect of activity.
* `TOAST` becomes relevant when filtering/sorting on large `text`/`jsonb` columns — Postgres must detoast (decompress/fetch out-of-line) the value to evaluate a `LIKE`/`ILIKE` predicate or to sort by it, which is significantly more expensive than filtering/sorting on a small in-line column. Filtering on a smaller indexed column first, when possible, avoids unnecessary detoasting.

---

## 10. Performance & Benchmarking

**`EXPLAIN ANALYZE` impacts**
```sql
EXPLAIN ANALYZE
SELECT status, COUNT(*), SUM(total)
FROM orders
WHERE placed_at >= '2026-01-01'
GROUP BY status;
```
* Look for `Seq Scan` vs `Index Scan`/`Bitmap Heap Scan` feeding the `WHERE` clause — a sequential scan on a large, selective filter usually means a missing or unused index.
* Look for `HashAggregate` (builds an in-memory hash table per group — fast, but memory-bound by `work_mem`) vs `GroupAggregate` (requires pre-sorted input — used when the planner expects too many distinct groups to hash efficiently, or when the input is already sorted for free via an index).
* A `Sort` node with `Sort Method: external merge  Disk: ...` in the output means the sort spilled to disk because it exceeded `work_mem` — a strong signal to either raise `work_mem` for that workload or add an index that lets `ORDER BY` be satisfied without a separate sort.

### Performance Improvement Tips

| Tip | Why It Helps |
|---|---|
| Prefer `EXISTS`/`NOT EXISTS` over `IN`/`NOT IN` for correlated subqueries | Enables semi-join short-circuiting and avoids the `NOT IN` + `NULL` silent-zero-rows trap |
| Avoid leading-wildcard `LIKE '%x%'` on large tables without a `pg_trgm` GIN index | A leading wildcard cannot use a standard B-tree index and forces a full scan |
| Filter in `WHERE`, not `HAVING`, whenever the condition doesn't need an aggregate | Reduces the row count *before* the expensive grouping/aggregation step |
| Create indexes matching your `ORDER BY` column(s) and direction (including `NULLS FIRST/LAST` when relevant) | Lets Postgres serve sorted results directly from the index, skipping a separate `Sort` node entirely |
| Use composite indexes covering both the `WHERE` filter and `GROUP BY`/`ORDER BY` columns together | A single index can satisfy filtering, grouping, and ordering in one pass instead of three separate operations |
| Watch `work_mem` against `GROUP BY`/`ORDER BY` on large result sets | Too small forces disk-based sort/hash spills (`external merge`), which is dramatically slower than in-memory |
| Use `ILIKE` only when case-insensitivity is genuinely required | `ILIKE` cannot use a standard B-tree index at all (needs `citext` or a functional/trigram index) — `LIKE` can, for prefix matches |

**Memory vs. disk trade-offs**
* `GROUP BY` via `HashAggregate` trades memory for speed — it's fast as long as the hash table of all distinct groups fits in `work_mem`; once it doesn't, Postgres (since PG13) can spill groups to disk incrementally rather than failing outright, but this is slower than an in-memory hash.
* `ORDER BY` without index support similarly trades `work_mem` for speed — an in-memory quicksort/heapsort is fast; a disk-spilled external merge sort is an order of magnitude slower, and is a strong signal to either add a supporting index or paginate/limit the result set.

---

## 11. Best Practices & Common Mistakes

* ✅ **Do**: Follow standard PostgreSQL conventions — prefer `EXISTS` for correlated existence checks, filter in `WHERE` before aggregating, and always specify `NULLS FIRST`/`LAST` explicitly when NULL ordering matters to the feature.
* ❌ **Don't**: Ignore index overhead or transaction locking — every filter/sort/group performance problem should be diagnosed with `EXPLAIN ANALYZE` before reaching for a "just add more indexes" instinct, since more indexes also cost more on every write.
* ⚠️ **Common Mistake**: N+1 queries — issuing a query per row instead of a single query using `WHERE ... IN (...)` or a `JOIN`; missing indexes on foreign keys compound this by making each of those N queries slow individually too.
* ⚠️ **Common Mistake**: Using `NOT IN (subquery)` where the subquery can return `NULL` — silently returns zero rows for the entire query, with no error raised.
* ⚠️ **Common Mistake**: Relying on default `NULLS FIRST`/`LAST` behavior instead of specifying it explicitly — the default flips between `ASC` and `DESC`, which is a common source of "why are the NULLs at the top now" bugs when someone flips a sort direction.

---

## 12. Interview Questions

1. **Beginner**: Explain the difference between `LIKE` and `ILIKE` to a junior dev, and when a `LIKE` pattern can and can't use an index.
2. **Beginner**: What's the default `NULL` placement for `ORDER BY ... ASC` versus `ORDER BY ... DESC`, and how do you override it?
3. **Intermediate**: How would you optimize a query that filters, groups, and sorts a large table, and what would you look for in `EXPLAIN ANALYZE` to confirm the optimization worked?
4. **Intermediate**: Why is `NOT IN (subquery)` dangerous in PostgreSQL, and what's the safe alternative?
5. **Intermediate**: What's the difference between `WHERE` and `HAVING`, and why can't `WHERE` reference an aggregate function like `SUM()`?
6. **Expert**: Describe the low-level execution difference between `HashAggregate` and `GroupAggregate` for a `GROUP BY` query, and explain the conditions under which the planner would choose one over the other.
7. **Expert**: Explain how `EXISTS` is executed as a semi-join internally, and why this allows short-circuiting in a way that a naive `IN (subquery)` historically could not always guarantee.
8. **Expert**: In a high-concurrency environment, describe the memory behavior of a `GROUP BY` query whose distinct group count is much larger than expected — what happens as `work_mem` is exceeded, how does PostgreSQL's spill-to-disk behavior work (post-PG13), and what locking (if any) is involved for a pure read query versus a concurrent write to the same table?

---

*Primary sources: PostgreSQL Official Documentation (postgresql.org/docs) — chapters on Queries (SELECT), Aggregate Functions, Pattern Matching, Subquery Expressions, and the `pg_trgm` contrib module.*