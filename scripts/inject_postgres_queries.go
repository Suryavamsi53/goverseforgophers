package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var queries = map[string]string{
	"001-level-1-postgresql-fundamentals.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Connect and check version\nSELECT version();\n\n-- Create a new database\nCREATE DATABASE fin_app;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Create schema for tenant isolation\nCREATE SCHEMA IF NOT EXISTS tenant_a;\n\n-- Set search path to prioritize tenant schema\nSET search_path TO tenant_a, public;\n```" + `
`,
	"002-level-2-crud-operations.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Standard Insert\nINSERT INTO users (email, name) VALUES ('test@test.com', 'John');\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Bulk Insert with RETURNING (Avoids a 2nd SELECT query)\nINSERT INTO logs (level, message)\nVALUES \n  ('INFO', 'Server started'),\n  ('ERROR', 'DB disconnected')\nRETURNING id, created_at;\n\n-- Update and return the old data for queue processing\nUPDATE orders SET status = 'processed'\nWHERE status = 'pending'\nRETURNING id, customer_id;\n```" + `
`,
	"003-level-3-constraints.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Basic Constraints\nCREATE TABLE accounts (\n    id SERIAL PRIMARY KEY,\n    email VARCHAR(255) UNIQUE NOT NULL,\n    balance NUMERIC CHECK (balance >= 0)\n);\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Exclusion Constraint: Prevent double-booking a hotel room in the same time range\nCREATE EXTENSION IF NOT EXISTS btree_gist;\n\nCREATE TABLE room_bookings (\n    room_id INT,\n    booking_period TSTZRANGE,\n    EXCLUDE USING GIST (\n        room_id WITH =,\n        booking_period WITH &&\n    )\n);\n```" + `
`,
	"004-level-4-querying.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Standard filtering\nSELECT * FROM customers WHERE status = 'active' AND age > 18;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- EXISTS vs IN for high-performance subqueries\n-- This stops scanning the 'orders' table the millisecond it finds 1 match.\nSELECT id, name \nFROM customers c\nWHERE EXISTS (\n    SELECT 1 \n    FROM orders o \n    WHERE o.customer_id = c.id \n    AND o.total > 1000\n);\n```" + `
`,
	"005-level-5-joins.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Standard Inner Join\nSELECT users.name, orders.total\nFROM users\nINNER JOIN orders ON users.id = orders.user_id;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- LATERAL JOIN: Get the Top 3 most recent purchases FOR EACH customer\nSELECT c.name, recent_purchases.*\nFROM customers c\nCROSS JOIN LATERAL (\n    SELECT total, created_at\n    FROM orders o\n    WHERE o.customer_id = c.id\n    ORDER BY created_at DESC\n    LIMIT 3\n) AS recent_purchases;\n```" + `
`,
	"006-level-6-aggregate-functions.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Standard Aggregation\nSELECT department_id, AVG(salary), COUNT(*)\nFROM employees\nGROUP BY department_id;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- JSON_AGG: Have PostgreSQL build nested JSON directly (eliminates ORM N+1 issues)\nSELECT \n    p.id as post_id,\n    p.title,\n    JSON_AGG(\n        JSON_BUILD_OBJECT('id', c.id, 'body', c.body)\n    ) as comments\nFROM posts p\nLEFT JOIN comments c ON p.id = c.post_id\nGROUP BY p.id;\n```" + `
`,
	"007-level-7-built-in-functions.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Basic String and Math\nSELECT UPPER(name), ROUND(price, 2) FROM products;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- COALESCE for Fallbacks and DATE_TRUNC for time-series aggregation\nSELECT \n    DATE_TRUNC('hour', created_at) as signup_hour,\n    COUNT(COALESCE(referral_code, 'organic')) as total_signups\nFROM users\nGROUP BY 1 \nORDER BY 1 DESC;\n```" + `
`,
	"008-level-8-views.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Standard View (Logical alias, runs query every time)\nCREATE VIEW active_users AS \nSELECT * FROM users WHERE status = 'active';\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Materialized View (Caches result to disk for reporting)\nCREATE MATERIALIZED VIEW monthly_revenue_report AS\nSELECT DATE_TRUNC('month', created_at) as month, SUM(amount) \nFROM transactions \nGROUP BY 1;\n\n-- Refresh it concurrently via pg_cron (zero downtime for readers)\nREFRESH MATERIALIZED VIEW CONCURRENTLY monthly_revenue_report;\n```" + `
`,
	"009-level-9-indexing.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Standard B-Tree Index\nCREATE INDEX idx_users_email ON users(email);\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Partial Index: Extremely small and fast index for a tiny subset of data\nCREATE INDEX idx_pending_orders ON orders(created_at) \nWHERE status = 'pending';\n\n-- GIN Index: O(1) searches inside JSONB documents\nCREATE INDEX idx_users_metadata ON users USING GIN (metadata);\n```" + `
`,
	"010-level-10-transactions.md": `
### 🔹 Basic Implementation
` + "```sql\nBEGIN;\nUPDATE accounts SET balance = balance - 100 WHERE id = 1;\nUPDATE accounts SET balance = balance + 100 WHERE id = 2;\nCOMMIT;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Serializable Isolation: Prevents all race conditions and anomalies\nBEGIN;\nSET TRANSACTION ISOLATION LEVEL SERIALIZABLE;\n\n-- If another transaction modifies this balance simultaneously,\n-- PostgreSQL will mathematically detect the conflict and abort one of them.\nSELECT balance FROM accounts WHERE id = 1;\nUPDATE accounts SET balance = balance - 100 WHERE id = 1;\nCOMMIT;\n```" + `
`,
	"011-level-11-normalization.md": `
### 🔹 Basic Implementation
` + "```sql\n-- 3NF Database Structure\nCREATE TABLE customers (id SERIAL, name TEXT);\nCREATE TABLE orders (id SERIAL, customer_id INT);\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Intentional Denormalization for extreme read performance\n-- Instead of COUNT() joining millions of likes on every API call,\n-- store the aggregate directly on the post.\nALTER TABLE posts ADD COLUMN cached_like_count INT DEFAULT 0;\n\n-- Updated via Trigger when a like is inserted.\n```" + `
`,
	"012-level-12-advanced-sql.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Common Table Expression (CTE)\nWITH ActiveUsers AS (\n    SELECT id FROM users WHERE last_login > NOW() - INTERVAL '7 days'\n)\nSELECT * FROM orders WHERE user_id IN (SELECT id FROM ActiveUsers);\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Window Functions: Calculate Month-over-Month Growth without Self-Joins\nSELECT \n    month,\n    revenue,\n    LAG(revenue) OVER (ORDER BY month) as previous_month_revenue,\n    revenue - LAG(revenue) OVER (ORDER BY month) as growth\nFROM monthly_stats;\n```" + `
`,
	"013-level-13-json-jsonb.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Select a specific key from a JSON object\nSELECT metadata->>'theme' as user_theme FROM users;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- JSONB Containment Operator (@>) utilizing GIN Index\n-- Lightning fast search for all users who have dark mode enabled\nSELECT id, email \nFROM users \nWHERE settings @> '{\"theme\": \"dark\"}'::jsonb;\n\n-- Updating a nested JSON key in-place\nUPDATE users \nSET settings = jsonb_set(settings, '{notifications,email}', 'false'::jsonb)\nWHERE id = 1;\n```" + `
`,
	"014-level-14-arrays.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Creating and querying an array\nCREATE TABLE posts (id SERIAL, tags TEXT[]);\nSELECT * FROM posts WHERE 'golang' = ANY(tags);\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Overlap Operator (&&) with GIN Indexing\n-- Find posts that have ANY of the provided tags\nCREATE INDEX idx_posts_tags ON posts USING GIN (tags);\n\nSELECT * FROM posts \nWHERE tags && ARRAY['microservices', 'distributed-systems'];\n```" + `
`,
	"015-level-15-uuid.md": `
### 🔹 Basic Implementation
` + "```sql\nCREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";\nSELECT uuid_generate_v4();\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Using UUID as Primary Key for Distributed Systems\nCREATE TABLE distributed_events (\n    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- Built-in PG 13+\n    tenant_id UUID NOT NULL,\n    payload JSONB\n);\n```" + `
`,
	"016-level-16-sequences.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Under the hood of SERIAL\nCREATE SEQUENCE order_id_seq;\nSELECT nextval('order_id_seq');\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Fixing a sequence desync after manual DB imports\n-- If max(id) is 5000, next insert will fail if sequence is at 100\nSELECT setval('users_id_seq', (SELECT COALESCE(MAX(id), 1) FROM users));\n```" + `
`,
	"017-level-17-stored-procedures.md": `
### 🔹 Basic Implementation
` + "```sql\nCREATE FUNCTION get_active_user_count() RETURNS INT AS $$\n    SELECT count(*)::INT FROM users WHERE status = 'active';\n$$ LANGUAGE sql;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Stored Procedure (Supports internal Transactions!)\nCREATE OR REPLACE PROCEDURE transfer_funds(sender INT, receiver INT, amount NUMERIC) \nLANGUAGE plpgsql AS $$\nBEGIN\n    UPDATE accounts SET balance = balance - amount WHERE id = sender;\n    UPDATE accounts SET balance = balance + amount WHERE id = receiver;\n    COMMIT; -- Procedures can commit mid-execution, Functions cannot.\nEND;\n$$;\n```" + `
`,
	"018-level-18-triggers.md": `
### 🔹 Basic Implementation
` + "```sql\nCREATE TRIGGER update_timestamp\nBEFORE UPDATE ON users\nFOR EACH ROW\nEXECUTE FUNCTION trigger_set_timestamp();\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Bulletproof Audit Logging Trigger\nCREATE OR REPLACE FUNCTION audit_log_changes() RETURNS trigger AS $$\nBEGIN\n    INSERT INTO audit_logs (table_name, record_id, old_data, new_data)\n    VALUES (TG_TABLE_NAME, OLD.id, row_to_json(OLD), row_to_json(NEW));\n    RETURN NEW;\nEND;\n$$ LANGUAGE plpgsql;\n\nCREATE TRIGGER users_audit\nAFTER UPDATE ON users\nFOR EACH ROW EXECUTE FUNCTION audit_log_changes();\n```" + `
`,
	"019-level-19-plpgsql.md": `
### 🔹 Basic Implementation
` + "```sql\nDO $$\nDECLARE\n    user_count INT;\nBEGIN\n    SELECT count(*) INTO user_count FROM users;\n    RAISE NOTICE 'Total users: %', user_count;\nEND $$;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Dynamic SQL Generation for Table Partitioning via Cron\nDO $$\nDECLARE\n    next_month text := to_char(NOW() + interval '1 month', 'YYYY_MM');\nBEGIN\n    EXECUTE format(\n        'CREATE TABLE IF NOT EXISTS metrics_%s PARTITION OF metrics FOR VALUES FROM (%L) TO (%L)', \n        next_month, \n        DATE_TRUNC('month', NOW() + interval '1 month'),\n        DATE_TRUNC('month', NOW() + interval '2 months')\n    );\nEND $$;\n```" + `
`,
	"020-level-20-security.md": `
### 🔹 Basic Implementation
` + "```sql\nCREATE ROLE api_worker LOGIN PASSWORD 'secret';\nGRANT SELECT, INSERT ON users TO api_worker;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Row Level Security (RLS) for Multi-Tenant Data Isolation\nALTER TABLE tenant_data ENABLE ROW LEVEL SECURITY;\n\nCREATE POLICY isolate_tenants ON tenant_data\nUSING (tenant_id = current_setting('app.current_tenant')::UUID);\n\n-- Go backend sets context before query:\n-- SET LOCAL app.current_tenant = 'uuid...';\n-- DB physically blocks reading other tenants even if Go has a bug.\n```" + `
`,
	"021-level-21-backup-restore.md": `
### 🔹 Basic Implementation
` + "```bash\n# Logical Backup (SQL Dump)\npg_dump -U postgres -d mydb > backup.sql\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```bash\n# Physical Backup with WAL-G / pgBackRest for PITR (Point in Time Recovery)\n# In postgresql.conf:\narchive_mode = on\narchive_command = 'wal-g wal-push %p'\n\n# To restore DB to exact millisecond before DROP TABLE:\nrestore_command = 'wal-g wal-fetch %f %p'\nrecovery_target_time = '2026-07-08 14:15:00'\n```" + `
`,
	"022-level-22-performance-tuning.md": `
### 🔹 Basic Implementation
` + "```sql\n-- View the execution plan without running the query\nEXPLAIN SELECT * FROM orders WHERE status = 'shipped';\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- EXPLAIN ANALYZE actually runs the query and compares Planner estimates vs Reality\nEXPLAIN (ANALYZE, BUFFERS) \nSELECT * FROM large_table WHERE indexed_col = 123;\n\n-- Watch for \"Seq Scan\" (Disk Read) vs \"Index Only Scan\" (RAM Read)\n-- Watch for \"Buffers: shared hit=500 read=1000\" (Read from disk means lack of RAM)\n```" + `
`,
	"023-level-23-partitioning.md": `
### 🔹 Basic Implementation
` + "```sql\nCREATE TABLE logs (\n    id SERIAL,\n    created_at TIMESTAMP\n) PARTITION BY RANGE (created_at);\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Declarative Range Partitioning for High-Volume Time-Series\nCREATE TABLE logs_2026_07 PARTITION OF logs\nFOR VALUES FROM ('2026-07-01') TO ('2026-08-01');\n\n-- Deleting old data? NEVER use DELETE. Use DROP TABLE.\n-- DROP TABLE logs_2024_01; -> instantly frees 500GB with ZERO disk IO.\n```" + `
`,
	"024-level-24-replication.md": `
### 🔹 Basic Implementation
` + "```sql\n-- In postgresql.conf for Master\nwal_level = replica\nmax_wal_senders = 10\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Logical Replication (CDC: Change Data Capture)\n-- Master (Publisher): Only replicate specific tables to the Data Warehouse\nCREATE PUBLICATION analytics_pub FOR TABLE users, orders;\n\n-- Analytics DB (Subscriber)\nCREATE SUBSCRIPTION analytics_sub \nCONNECTION 'host=master port=5432 user=rep_user password=secret dbname=prod'\nPUBLICATION analytics_pub;\n```" + `
`,
	"025-level-25-high-availability.md": `
### 🔹 Basic Implementation
` + "```ini\n# PgBouncer Config\n[databases]\nprod_db = host=127.0.0.1 port=5432 dbname=prod\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```ini\n# Transaction-level connection pooling for Go/Serverless scaling\n[pgbouncer]\npool_mode = transaction\nmax_client_conn = 10000\ndefault_pool_size = 50\n\n# 10,000 incoming Lambda/Go connections multiplexed onto just 50 physical Postgres sockets.\n```" + `
`,
	"026-level-26-full-text-search.md": `
### 🔹 Basic Implementation
` + "```sql\nSELECT * FROM articles \nWHERE body ILIKE '%postgres%'; -- Slow, O(N)\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Native Full Text Search (Replaces Elasticsearch for medium datasets)\n-- Creates GIN index on lexemes (stemmed words)\nCREATE INDEX idx_articles_fts ON articles USING GIN (to_tsvector('english', body));\n\n-- Ranks results by relevance mathematically\nSELECT title, ts_rank(to_tsvector('english', body), to_tsquery('english', 'fast & database')) as rank\nFROM articles\nWHERE to_tsvector('english', body) @@ to_tsquery('english', 'fast & database')\nORDER BY rank DESC;\n```" + `
`,
	"027-level-27-extensions.md": `
### 🔹 Basic Implementation
` + "```sql\nCREATE EXTENSION IF NOT EXISTS pg_stat_statements;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- PostGIS: Geospatial calculations in milliseconds\nCREATE EXTENSION postgis;\n\n-- Find all delivery drivers within 5km of a restaurant's coordinates\nSELECT id, name \nFROM drivers\nWHERE ST_DWithin(\n    location_geom, \n    ST_MakePoint(-122.4194, 37.7749)::geography, \n    5000 -- meters\n);\n```" + `
`,
	"028-level-28-monitoring.md": `
### 🔹 Basic Implementation
` + "```sql\n-- View all currently executing queries\nSELECT pid, query, state FROM pg_stat_activity WHERE state = 'active';\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Find the top 5 most CPU/Disk intensive queries in your cluster\nSELECT \n    query, \n    calls, \n    total_exec_time / 1000 as total_seconds, \n    mean_exec_time as avg_ms\nFROM pg_stat_statements\nORDER BY total_exec_time DESC\nLIMIT 5;\n```" + `
`,
	"029-level-29-concurrency.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Row-level locking (Pessimistic)\nSELECT * FROM wallets WHERE user_id = 1 FOR UPDATE;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Advisory Locks: App-level distributed locks managed by PostgreSQL\n-- E.g., Ensure only ONE cron job worker processes payouts globally.\nSELECT pg_try_advisory_lock(9999); -- Returns TRUE if lock acquired\n\n-- In Go:\n-- if acquired { processPayouts(); pg_advisory_unlock(9999); }\n```" + `
`,
	"030-level-30-advanced-storage.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Standard table size checking\nSELECT pg_size_pretty(pg_total_relation_size('events'));\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- TOAST (The Oversized-Attribute Storage Technique)\n-- 8KB page limits mean big JSON is compressed and stored out-of-line automatically.\n-- You can force external storage without compression for extremely fast writes of large payloads:\nALTER TABLE heavy_payloads ALTER COLUMN raw_json SET STORAGE EXTERNAL;\n```" + `
`,
	"031-level-31-postgresql-internals.md": `
### 🔹 Basic Implementation
` + "```sql\nSHOW shared_buffers;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Internals tuning for massive scale\n-- postgresql.conf\nshared_buffers = 16GB       # Cache frequently accessed data (25% of RAM)\nwork_mem = 64MB             # RAM allocated PER sort/hash operation\nmaintenance_work_mem = 2GB  # Speeds up VACUUM and CREATE INDEX\neffective_cache_size = 48GB # Hints to Planner how much RAM OS has for caching\n```" + `
`,
	"032-level-32-distributed-postgresql.md": `
### 🔹 Basic Implementation
` + "```sql\n-- Single node logic\nCREATE TABLE user_events (id SERIAL, tenant_id INT, data JSONB);\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Citus Extension for Horizontal Sharding\nCREATE TABLE user_events (id BIGSERIAL, tenant_id INT, data JSONB);\n\n-- Distribute the table across 10+ worker nodes transparently\nSELECT create_distributed_table('user_events', 'tenant_id');\n\n-- Queries filtering by tenant_id are routed to the exact node instantly.\n```" + `
`,
	"033-level-33-postgresql-with-go.md": `
### 🔹 Basic Implementation
` + "```go\n// database/sql + lib/pq\ndb, err := sql.Open(\"postgres\", uri)\nrows, err := db.Query(\"SELECT name FROM users\")\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```go\n// pgxpool for native Postgres protocol and connection pooling\npool, _ := pgxpool.New(context.Background(), \"postgres://...\")\n\n// Bulk Copy (Millions of rows in seconds, bypassing INSERT overhead)\nrows := [][]interface{}{{\"Alice\", 25}, {\"Bob\", 30}}\npool.CopyFrom(\n    context.Background(),\n    pgx.Identifier{\"users\"},\n    []string{\"name\", \"age\"},\n    pgx.CopyFromRows(rows),\n)\n```" + `
`,
	"034-level-34-interview-topics.md": `
### 🔹 Basic Implementation
` + "```sql\n-- N+1 Query Problem\n-- Query 1: SELECT * FROM users;\n-- Query 2-100: SELECT * FROM posts WHERE user_id = X;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- MVCC and Dead Tuples explanation via SQL\n-- After 1M UPDATEs, physical disk is full of \"dead\" invisible rows.\n-- Check bloat:\nSELECT n_dead_tup, n_live_tup FROM pg_stat_user_tables WHERE relname = 'users';\n\n-- Fix bloat (Autovacuum usually handles this, but MANUAL vacuum forces it):\nVACUUM ANALYZE users;\n```" + `
`,
	"035-level-35-real-world-backend-patterns.md": `
### 🔹 Basic Implementation
` + "```sql\n-- OFFSET Pagination (Gets exponentially slower as offset grows)\nSELECT * FROM feed ORDER BY created_at DESC LIMIT 20 OFFSET 50000;\n```" + `

### 🔹 Advanced / Optimized Implementation
` + "```sql\n-- Cursor (Keyset) Pagination: O(1) Time Complexity\nSELECT * FROM feed \nWHERE (created_at, id) < ('2026-07-08 14:00:00', 98765)\nORDER BY created_at DESC, id DESC \nLIMIT 20;\n\n-- Upsert Pattern (Insert if new, Update if exists)\nINSERT INTO user_stats (user_id, logins)\nVALUES (1, 1)\nON CONFLICT (user_id) DO UPDATE \nSET logins = user_stats.logins + 1;\n```" + `
`,
}

func main() {
	dir := "content/lessons/postgresql"
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	replaceTarget := "## 4. Code & Query Implementation\n### 🔹 Basic Implementation\n```sql\n-- Standard query example\nSELECT * FROM table;\n```\n\n### 🔹 Advanced / Optimized Implementation\n```sql\n-- Optimized query with indexes or advanced features\n```\n"

	successCount := 0

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".md") {
			continue
		}

		queryReplacement, exists := queries[file.Name()]
		if !exists {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		contentBytes, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading %s: %v\n", file.Name(), err)
			continue
		}

		content := string(contentBytes)
		
		// The template generated in the scaffold has a specific format.
		// Let's replace the Code & Query Implementation section safely.
		if strings.Contains(content, replaceTarget) {
			newContent := strings.Replace(content, replaceTarget, "## 4. Code & Query Implementation\n" + queryReplacement, 1)
			err = os.WriteFile(filePath, []byte(newContent), 0644)
			if err != nil {
				fmt.Printf("Error writing %s: %v\n", file.Name(), err)
			} else {
				successCount++
			}
		} else {
            // Wait, let's try a fallback replacement just in case spaces mismatch.
            // Split by "## 4. Code & Query Implementation" and "## 5. Internals & Under the Hood"
            parts := strings.Split(content, "## 4. Code & Query Implementation")
            if len(parts) == 2 {
                subParts := strings.Split(parts[1], "## 5. Internals & Under the Hood")
                if len(subParts) == 2 {
                    newContent := parts[0] + "## 4. Code & Query Implementation\n" + queryReplacement + "\n## 5. Internals & Under the Hood" + subParts[1]
                    os.WriteFile(filePath, []byte(newContent), 0644)
                    successCount++
                }
            }
        }
	}

	fmt.Printf("Successfully injected queries into %d markdown files.\n", successCount)
}
