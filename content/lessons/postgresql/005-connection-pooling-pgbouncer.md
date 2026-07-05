# Connection Pooling and PgBouncer

In the Database module, we learned how to configure the internal Go `sql.DB` connection pool (`SetMaxOpenConns`).

However, if you scale your architecture to massive heights, the Go internal connection pool is not enough. You will hit a terrifying bottleneck inherent to PostgreSQL's process architecture.

## 1. The PostgreSQL Process Model

Unlike MySQL or Redis (which use lightweight threads to handle connections), PostgreSQL uses a **Process-Per-Connection** model.

Every single time a Go server opens a TCP connection to Postgres, the Postgres operating system forks a completely new, heavy OS process to handle that connection.
* Forking a process takes ~10 milliseconds (very slow).
* Each idle process consumes ~10 Megabytes of RAM.

If you scale your Go API to 50 Kubernetes Pods, and each Pod has a `SetMaxOpenConns(100)`, your cluster will attempt to open 5,000 simultaneous connections to Postgres.
Postgres will try to fork 5,000 processes, instantly consume 50 Gigabytes of RAM, and crash.

## 2. The Solution: PgBouncer

To solve this, we introduce a piece of infrastructure called a **Connection Multiplexer** (the most famous being **PgBouncer**).

PgBouncer sits between your Go applications and the PostgreSQL database.

1. **Client to PgBouncer**: PgBouncer accepts thousands of lightweight connections from your 50 Go Pods.
2. **PgBouncer to Postgres**: PgBouncer only maintains a very small, fixed number of heavy connections to the actual PostgreSQL database (e.g., 50 connections).
3. **The Multiplexing**: When Go Pod #1 executes a SQL query, PgBouncer grabs one of the 50 real Postgres connections, runs the query, and instantly returns the connection to the pool.

Postgres never sees 5,000 connections. It only sees 50 stable, permanent connections. RAM usage remains flat, and performance skyrockets.

## 3. PgBouncer Pooling Modes

When you install PgBouncer, you must configure its Pooling Mode. If you choose the wrong one, your Go application will break!

1. **Session Pooling**: The Go app grabs a connection, and PgBouncer locks that connection to the Go app until the Go app formally calls `db.Close()`. (This defeats the purpose of multiplexing).
2. **Transaction Pooling (The Standard)**: The Go app only locks the connection for the duration of a `BEGIN ... COMMIT` block. Once the transaction commits, PgBouncer instantly steals the connection back, even if the Go app hasn't closed the TCP socket! (This is what you must use!).
3. **Statement Pooling**: PgBouncer steals the connection back after *every single SQL statement*. (You cannot use multi-statement Transactions in this mode!).

## 4. Go Prepared Statement Danger

If you use **Transaction Pooling** (which you should), you will encounter a massive bug in Go.

The Go standard library `database/sql` heavily uses **Prepared Statements** under the hood. 
* Go asks Postgres: `PREPARE stmt1 AS SELECT * FROM users WHERE id=$1`
* Go then calls: `EXECUTE stmt1(5)`

If PgBouncer steals the connection between the `PREPARE` and the `EXECUTE`, the `EXECUTE` command might be sent to a completely different Postgres backend process that has no idea what `stmt1` is! The query crashes!

**The Fix:** 
If you are using PgBouncer in Transaction Mode, you MUST disable Prepared Statements in your Go PostgreSQL driver. 

If you are using the popular `pgx` driver, you simply append a flag to the connection string:
`postgres://user:pass@pgbouncer:6432/db?statement_cache_capacity=0&default_query_exec_mode=exec`
