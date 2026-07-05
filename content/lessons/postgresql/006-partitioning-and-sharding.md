# Partitioning and Sharding

When a PostgreSQL table hits 500 million rows, performance degrades rapidly. 
* B-Tree indexes become larger than RAM, forcing the database to swap to disk (Cache Thrashing).
* Auto-Vacuum takes days to finish, causing table bloat.

To fix this, you must split the massive table into smaller, manageable chunks.

## 1. Table Partitioning (Single Server)

PostgreSQL natively supports **Declarative Partitioning**. This allows you to split one massive table into dozens of smaller physical tables, while still querying it as if it were a single table!

The most common strategy is **Range Partitioning** by Date (e.g., for Audit Logs or Metrics).

```sql
-- 1. Create the Master Table (It holds no data itself!)
CREATE TABLE audit_logs (
    id SERIAL,
    user_id INT,
    action VARCHAR(50),
    created_at DATE NOT NULL
) PARTITION BY RANGE (created_at);

-- 2. Create the physical Partitions (One for each month)
CREATE TABLE audit_logs_jan2026 PARTITION OF audit_logs 
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

CREATE TABLE audit_logs_feb2026 PARTITION OF audit_logs 
    FOR VALUES FROM ('2026-02-01') TO ('2026-03-01');
```

**The Superpower (Partition Pruning):**
When your Go app runs `SELECT * FROM audit_logs WHERE created_at = '2026-01-15'`, the Postgres Query Planner detects the date, entirely ignores the `feb2026` table, and only scans the `jan2026` table! 
This keeps your indexes tiny, ensuring they always fit in RAM.

## 2. Archiving (Dropping Partitions)

If you have a strict compliance rule that logs must be deleted after 30 days, running `DELETE FROM audit_logs WHERE created_at < NOW() - INTERVAL '30 days'` on a 1-billion row table will lock the database and destroy performance.

With Partitioning, you do not use `DELETE`. You simply drop the entire table!
```sql
DROP TABLE audit_logs_jan2026;
```
This executes in 0.001 seconds and frees up 100 Gigabytes of disk space instantly, without generating a single byte of WAL or triggering a Vacuum!

## 3. Sharding (Multiple Servers)

Partitioning splits data across multiple files on the *same* hard drive.
What happens if the hard drive itself fills up, or the CPU hits 100%? You must split the data across multiple physical servers. This is called **Sharding**.

Native PostgreSQL does not support Sharding out-of-the-box (though features like Foreign Data Wrappers exist). 

If you need true distributed Sharding in Postgres, Enterprise Go teams use an extension called **Citus** (acquired by Microsoft).

### How Citus Works
Citus turns Postgres into a Distributed System.
1. You have 1 Coordinator Node and 5 Worker Nodes.
2. Your Go application connects to the Coordinator Node just like normal Postgres.
3. You run a Citus command: `SELECT create_distributed_table('users', 'company_id');`
4. Citus transparently shreds the `users` table, distributing the rows across the 5 Worker Nodes based on a hash of the `company_id` (Consistent Hashing!).
5. When Go runs `SELECT * FROM users WHERE company_id = 42`, the Coordinator Node instantly routes the query to the specific Worker Node holding that data.

With Citus, Postgres scales horizontally, capable of handling Petabytes of data and millions of transactions per second!
