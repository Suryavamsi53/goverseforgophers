# Redis Distributed Caching

## 1. Learning Objectives
* **What you'll learn**: The architecture of Redis as a high-performance in-memory data store, caching strategies (Cache-Aside), and distributed locking in Go.
* **Why it matters**: Database disk I/O is the #1 bottleneck in backend systems. Redis bypasses the disk entirely, serving data from RAM in under 1 millisecond, allowing your Go app to scale to millions of requests.
* **Where it's used**: Session stores, leaderboards, rate limiting, and caching massive database queries at companies like Twitter, GitHub, and StackOverflow.

---

## 2. Real-world Story
Imagine a brilliant math professor (PostgreSQL). Solving complex equations takes the professor 5 minutes. If 1,000 students ask the professor the exact same equation, the professor spends 5,000 minutes doing the same work.
Redis is the chalkboard outside the classroom. The first time the professor solves the equation, they write the answer on the chalkboard. The next 999 students just read the chalkboard in 1 second. 

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Go App] -->|1. GET user:42| B{Redis}
    B -->|2. Cache Hit!| A
    B -.->|3. Cache Miss (Not Found)| C[(PostgreSQL)]
    C -.->|4. Fetch heavy data| A
    A -.->|5. SET user:42 (Expire 1 hr)| B
    
    style B fill:#dc2626,color:#fff
```

---

## 4. Internal Working (Under the Hood)
Redis is fundamentally **Single-Threaded**. 
It executes all commands sequentially in a single core event loop. This sounds slow, but because it operates entirely in RAM, it can process 100,000+ commands per second! Because it is single-threaded, atomic operations (like `INCR`) are mathematically guaranteed to be free of race conditions without needing complex Mutex locks.

---

## 5. Compiler Behavior
* **Network I/O in Go**: The industry standard driver `github.com/redis/go-redis/v9` perfectly integrates with Go's `context` package. It utilizes connection pooling under the hood, multiplexing thousands of concurrent Goroutines onto a small pool of persistent TCP connections to the Redis server.

---

## 6. Memory Management
* **Eviction Policies**: RAM is expensive. If you cache too much data, Redis will run out of memory. You must configure an Eviction Policy (like `allkeys-lru` - Least Recently Used). When RAM hits 100%, Redis automatically deletes the oldest, least accessed keys to make room for new data.

---

## 7. Code Examples

### 🔹 Example 1: Simple (The Cache-Aside Pattern)
```go
import "github.com/redis/go-redis/v9"

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

func GetUserProfile(userID string) string {
    // 1. Try Redis first
    val, err := rdb.Get(ctx, "user:"+userID).Result()
    if err == nil { return val } // Cache Hit!

    // 2. Cache Miss! Fetch from Postgres (Simulated)
    data := "heavy_db_data_for_" + userID

    // 3. Save to Redis with a 1-hour expiration
    rdb.Set(ctx, "user:"+userID, data, time.Hour)
    return data
}
```

### 🔹 Example 2: Intermediate (Distributed Locks)
```go
// Preventing Race Conditions across 50 Go Servers
// SetNX = Set if Not eXists
success, _ := rdb.SetNX(ctx, "lock:process_invoice_1", "locked", 10*time.Second).Result()
if !success {
    return errors.New("Another Go server is currently processing this invoice!")
}
// Do heavy processing...
rdb.Del(ctx, "lock:process_invoice_1")
```

### 🔹 Example 3: Advanced (Pipelining)
```go
// Executing 100 commands in a single TCP network hop!
pipe := rdb.Pipeline()

for i := 0; i < 100; i++ {
    pipe.Incr(ctx, "page_views:home")
}
// Sends all 100 commands instantly, bypassing 99 TCP round-trips!
_, err := pipe.Exec(ctx)
```

### 🔹 Example 4: Production
```go
// Complex Data Structures: Redis isn't just Strings!
// Adding a user to a Leaderboard (ZSET) sorted by score.
rdb.ZAdd(ctx, "leaderboard:global", redis.Z{
    Score:  9500.5,
    Member: "Suryavamsi",
})

// Fetch the top 3 players instantly!
topPlayers, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard:global", 0, 2).Result()
```

### 🔹 Example 5: Interview
```go
// Q: Why must you always put an Expiration (TTL) on cache keys?
// A: If the underlying Database data changes, the cache becomes "Stale". 
// A TTL guarantees that the cache will eventually delete itself, forcing the app to fetch fresh data.
```

---

## 8. Production Examples
1. **Rate Limiting**: Limiting a user to 100 requests per minute by using Redis `INCR` and `EXPIRE`. It is globally accurate across all your microservices.
2. **Session Storage**: Instead of JWTs, storing a secure random `session_id` cookie in the browser, and keeping the actual JSON session data in Redis. If you need to instantly ban a user, you just `DEL` their key in Redis!

---

## 9. Performance & Benchmarking
* **The Cache Stampede (Thundering Herd)**: If a heavy database query is cached in Redis, and the cache expires... 10,000 concurrent Goroutines might suffer a Cache Miss at the exact same millisecond. They will ALL query PostgreSQL simultaneously, crashing the database! (Mitigated by Go's `golang.org/x/sync/singleflight`).

---

## 10. Best Practices
* ✅ **Do**: Use descriptive key namespaces with colons (e.g., `user:101:profile`).
* ❌ **Don't**: Use the Redis `KEYS *` command in production! Because Redis is single-threaded, scanning 10 million keys will freeze the entire server for 5 seconds, causing a massive global outage for your app. Use `SCAN` instead.
* 🏢 **Google / Uber / Netflix Style**: For hyper-scale, use Redis Cluster, which shards your keys automatically across multiple master nodes using Hash Slots.

---

## 11. Common Mistakes
1. **JSON Serialization Overhead**: Fetching a huge JSON string from Redis and running `json.Unmarshal` in Go is extremely CPU-heavy. For massive datasets, consider Redis Hashes (`HGETALL`), or caching the pre-rendered HTML/Protobuf directly.
2. **Lack of Context Timeouts**: If the Redis server hangs, and you don't use `context.WithTimeout`, your Go Goroutines will freeze indefinitely waiting for a response.

---

## 12. Debugging
How to troubleshoot Redis in production:
* **redis-cli MONITOR**: Run this command to see a live stream of every single command hitting the Redis server in real-time. (Do not run this for long on a heavily loaded server, as it halves throughput!).
* **redis-cli INFO memory**: Check the exact RAM usage and cache hit/miss ratio.

---

## 13. Exercises
1. **Easy**: Connect to Redis and `SET` a string.
2. **Medium**: Implement a Cache-Aside wrapper function that fetches from a mock DB if the key doesn't exist in Redis.
3. **Hard**: Build a Global Rate Limiter middleware in Go using Redis `INCR` that blocks IPs exceeding 10 requests per minute.
4. **Expert**: Implement `singleflight` around your Cache Miss logic to mathematically prove that under 1,000 concurrent requests, the mock Database is only hit exactly 1 time.

---

## 14. Quiz
1. **MCQ**: Because Redis is single-threaded, what is true about the `INCR` command?
   * (A) It is slow (B) It requires a Go Mutex (C) It is mathematically atomic and thread-safe. *(Answer: C)*
2. **Code Review**: `rdb.Set(ctx, "key", "val", 0)`. What is the danger of setting a TTL of 0? *(It means the key NEVER expires. It will sit in RAM forever until explicitly deleted or evicted).*

---

## 15. FAANG Interview Questions
* **Beginner**: Explain the difference between Redis and PostgreSQL.
* **Intermediate**: What is a Cache Stampede and how do you prevent it?
* **Senior (Google/Meta)**: Explain the Redlock algorithm. Why is a simple `SETNX` distributed lock dangerous if the Redis Master node crashes and fails over to a Replica before replicating the lock?

---

## 16. Mini Project
**Real-Time Leaderboard API**
* Build a Go API with `POST /score` (updates a player's score).
* Build `GET /top10` (returns the top players).
* Use the Redis `ZSET` (Sorted Set) data structure.
* Benchmark it using `wrk`. You should easily achieve 20,000+ requests per second on a basic laptop!

---

## 17. Enterprise Features & Observability
* **Persistence (RDB vs AOF)**: Redis is in-memory, but it can save to disk! RDB takes periodic snapshots. AOF (Append Only File) logs every single write command to disk, guaranteeing zero data loss on restart, at the cost of slight disk I/O overhead.

---

## 18. Source Code Reading
Walkthrough of `github.com/redis/go-redis`.
* **The Connection Pool**: Study how `go-redis` natively handles severed connections, retries, and backoff algorithms completely transparently to the developer.

---

## 19. Architecture
* **Write-Through vs Cache-Aside**: In Cache-Aside, the Go app writes to the DB, then deletes the cache. In Write-Through, the Go app writes ONLY to the Cache, and an asynchronous process syncs the Cache to the DB.

---

## 20. Summary & Cheat Sheet
* **Speed**: ~1ms latency, pure RAM.
* **Architecture**: Single-threaded, atomic.
* **Mandatory**: Always set a TTL (Expiration).
* **Go Pattern**: Cache-Aside + Singleflight.
