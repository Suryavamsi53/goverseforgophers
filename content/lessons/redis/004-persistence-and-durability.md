# Persistence and Durability (RDB & AOF)

Redis is an In-Memory database. If the server loses power, RAM is erased, and all your data is destroyed.

If Redis is strictly used as a Cache, data loss is acceptable (the Go app will just experience a Cache Stampede and fetch the data from Postgres again). But if you are using Redis as your Primary Database, or for critical Rate Limiting counters, you need **Persistence**.

Redis offers two mechanisms to save RAM data to the physical hard drive.

## 1. RDB (Redis Database Backup) - Snapshots

RDB takes point-in-time snapshots of your entire dataset and saves it as a highly compressed binary file (`dump.rdb`) on the hard drive.

You configure RDB in `redis.conf`:
`save 60 1000` (Save to disk if at least 1,000 keys changed in the last 60 seconds).

* **How it works**: Redis is single-threaded! If it stops to write 10GB of data to the hard drive, it will freeze. Instead, Redis uses the Linux `fork()` system call. It clones the main process. The child process writes the memory to disk in the background, while the parent process continues to serve Go API requests instantly!
* **Pros**: Incredibly fast restarts. Tiny file sizes. Zero performance impact on the main thread.
* **Cons (The Flaw)**: Data Loss. If you take a snapshot every 60 seconds, and the server crashes at second 59, you permanently lose 59 seconds of data!

## 2. AOF (Append Only File) - Durability

If you cannot afford to lose 59 seconds of data, you must use **AOF**.

Instead of taking massive snapshots, AOF simply logs every single Write command (`SET`, `HINCRBY`) to a text file on the hard drive, exactly like PostgreSQL's Write-Ahead Log (WAL).

You configure AOF fsync policies in `redis.conf`:
* `appendfsync always`: Redis waits for the hard drive to confirm the write before returning `OK` to Go. (Mathematically safe, but brutally slow).
* `appendfsync everysec`: Redis returns `OK` to Go instantly, and a background thread flushes the log to the hard drive once per second. (The industry standard!).

* **Pros**: Incredible durability. At `everysec`, you will only ever lose 1 second of data maximum.
* **Cons**: The file grows infinitely! If you run `INCR counter` 1 million times, the AOF file contains 1 million lines of text, even though the actual data is just the number `1,000,000`.

### AOF Rewriting
To fix the infinite file growth, Redis automatically runs an **AOF Rewrite** in the background. It reads the current state of RAM, and rewrites the AOF file from scratch using the absolute minimum number of commands required to rebuild the state!

## 3. The Enterprise Standard (RDB + AOF)

Modern Redis deployments use BOTH simultaneously.

* Redis uses **AOF** for high-durability recovery upon reboot.
* Redis periodically generates **RDB** snapshots and ships them to AWS S3 for long-term disaster recovery backups!
