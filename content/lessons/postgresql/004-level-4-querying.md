# Level 4 - Querying

## 1. Learning Objectives
* **What you'll learn**: Master the core mechanics of Level 4 - Querying in PostgreSQL.
* **Why it matters**: Crucial for building scalable, high-performance, and robust backend systems.

---

## 2. Topics Covered
* Filtering (WHERE, AND, OR, BETWEEN, IN, LIKE, ILIKE, EXISTS)
* Sorting (ORDER BY, NULLS FIRST/LAST)
* Grouping (GROUP BY, HAVING)


---

## 3. Production Usage Scenarios (Real-world Examples)
For each concept, here is how we use it in a real production environment at scale:

### Scenario: Filtering (WHERE, AND, OR, BETWEEN, IN, LIKE, ILIKE, EXISTS)
**Context**: Imagine you are building a highly concurrent backend system.
**The Problem**: (Define the engineering challenge)
**The PostgreSQL Solution**: (How Level 4 - Querying solves it effectively in production)

### Scenario: Sorting (ORDER BY, NULLS FIRST/LAST)
**Context**: Imagine you are building a highly concurrent backend system.
**The Problem**: (Define the engineering challenge)
**The PostgreSQL Solution**: (How Level 4 - Querying solves it effectively in production)

### Scenario: Grouping (GROUP BY, HAVING)
**Context**: Imagine you are building a highly concurrent backend system.
**The Problem**: (Define the engineering challenge)
**The PostgreSQL Solution**: (How Level 4 - Querying solves it effectively in production)



---

## 4. Code & Query Implementation

### 🔹 Basic Implementation
```sql
-- Standard filtering
SELECT * FROM customers WHERE status = 'active' AND age > 18;
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
```

---

## 5. Internals & Under the Hood
* **How PostgreSQL engine processes this**: (Parser -> Planner -> Executor)
* **Storage impact**: WAL logs, Heap, and TOAST considerations.

---

## 6. Performance & Benchmarking
* **EXPLAIN ANALYZE impacts**
* **Memory vs Disk Trade-offs**

---

## 7. Best Practices & Common Mistakes
* ✅ **Do**: Follow standard PostgreSQL conventions.
* ❌ **Don't**: Ignore index overhead or transaction locking.
* ⚠️ **Common Mistake**: N+1 queries, missing indexes on foreign keys.

---

## 8. Interview Questions
1. **Beginner**: Explain Level 4 - Querying to a junior dev.
2. **Intermediate**: How would you optimize queries involving Level 4 - Querying?
3. **Expert**: Describe the low-level locking and memory behavior of Level 4 - Querying in high-concurrency environments.
