# PostgreSQL Architecture & Storage Engine

When you execute a SQL query from your Go application, PostgreSQL doesn't just write a string to a text file. It utilizes a highly optimized, complex storage engine designed to guarantee ACID compliance without sacrificing concurrency.

To write performant Go applications, you must understand how PostgreSQL physically stores data on disk.

## 1. The Page (8KB Blocks)

PostgreSQL does not read or write individual rows to the hard drive. The fundamental unit of storage in Postgres is the **Page** (which is strictly 8 Kilobytes).

Every table and index is divided into 8KB pages. 
If your Go application requests a user with `ID = 5`, Postgres goes to the hard drive, reads the entire 8KB page that contains that user, loads the 8KB page into RAM (the Shared Buffers), and then extracts the single row for you.

**The Performance Impact:**
If a user row is 1KB, one Page can fit 8 users.
If you run `SELECT * FROM users LIMIT 800`, Postgres has to read 100 Pages from the hard drive.
If you reduce the size of your rows (e.g., stopping the use of `VARCHAR(MAX)` for tiny strings, or dropping unused columns), you can fit 16 users per Page. Suddenly, that exact same query only requires 50 Page reads. Your database just became 100% faster!

## 2. Write-Ahead Logging (WAL)

If you insert a new Order, Postgres writes it to the 8KB Page currently sitting in RAM. 
If the server loses power 2 seconds later, RAM is erased, and your Order is lost forever.

To prevent data loss (Durability), Postgres uses the **Write-Ahead Log (WAL)**.
Before Postgres modifies the 8KB Page in RAM, it appends a tiny string to a sequential log file on the hard drive describing the change (e.g., "Insert Order 42"). 

* Writing to an append-only log file on disk is blazingly fast.
* If the server crashes, when it reboots, Postgres simply replays the WAL file to reconstruct the RAM state exactly as it was.

Because WAL is so fast, Postgres can return a `200 OK` to your Go application almost instantly, guaranteeing data safety without waiting to write the massive 8KB Page to the actual data file!

## 3. MVCC (Multi-Version Concurrency Control)

If your Go application runs a massive `SELECT` query that takes 5 seconds, and another Goroutine runs an `UPDATE` on one of those rows... does the `SELECT` query crash? Does it block the `UPDATE`?

PostgreSQL uses **MVCC** to ensure that "Readers never block Writers, and Writers never block Readers."

When you `UPDATE` a row, Postgres does **not** overwrite the old data. 
1. It physically creates a brand new row with the new data.
2. It marks the old row as "Dead" (a Tuple).
3. The 5-second `SELECT` query continues to read the "Dead" row! 
4. Any new queries that start *after* the update will read the new row.

### The VACUUM Process
Because MVCC leaves "Dead" rows sitting on the hard drive, your database will slowly bloat and run out of disk space. 
Postgres runs a background daemon called the **Auto-Vacuum**. It constantly sweeps the 8KB pages, looking for Dead rows that no active queries are reading anymore, and physically deletes them to free up space.

**Enterprise Trap**: If you run a massive transaction in Go that lasts for 24 hours, you will block the Vacuum process from deleting *any* dead rows during that time, causing catastrophic database bloat! Keep your transactions small and fast.
