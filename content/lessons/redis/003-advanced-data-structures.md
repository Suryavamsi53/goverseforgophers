# Advanced Data Structures

If you only use Redis for `SET` and `GET` (Strings), you are drastically underutilizing the database. Redis ships with incredible native Data Structures that execute complex operations in RAM in $O(1)$ time.

## 1. Hashes (`HSET`, `HGET`)

If you want to cache a User struct, you could convert it to JSON and `SET` it as a String. But if you only want to update the user's `age`, you have to download the entire JSON string, decode it in Go, update the age, encode it to JSON, and `SET` it back. This is slow and prone to race conditions.

Instead, use a **Redis Hash**. It acts like a native Go `map[string]string` stored directly in RAM.

```go
// Set individual fields on the Hash
rdb.HSet(ctx, "user:42", "name", "Alice", "age", 30)

// Increment the age field atomically! (No JSON decoding required!)
rdb.HIncrBy(ctx, "user:42", "age", 1) 

// Get a single field
name, _ := rdb.HGet(ctx, "user:42", "name").Result()
```

## 2. Sorted Sets (`ZADD`, `ZRANGE`)

How do you build a real-time global Leaderboard for a video game with 10 million players? Sorting a 10-million row table in Postgres takes 5 seconds. In Redis, it takes 1 millisecond.

A **Sorted Set (ZSET)** automatically maintains perfect sorting order using a Skip List data structure in RAM.

```go
// Add players with their Score!
rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 500, Member: "Player1"})
rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 900, Member: "Player2"})

// Instantly fetch the Top 10 players! (O(log N))
top10, _ := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, 9).Result()
for _, z := range top10 {
    fmt.Printf("%s: %f\n", z.Member, z.Score)
}
```
*Note: We also used Sorted Sets in the System Design module to build a Sliding Window Rate Limiter (using the UNIX timestamp as the Score)!*

## 3. Sets (`SADD`, `SISMEMBER`)

Sets store unique strings. They are mathematically optimized for Intersections, Unions, and Memberships.

* **Use Case (Mutual Friends)**:
  * User A's friends: `SADD friends:A "Bob" "Charlie" "Dave"`
  * User B's friends: `SADD friends:B "Charlie" "Eve" "Frank"`
  * Instantly find mutual friends: `SINTER friends:A friends:B` -> Returns `"Charlie"`.

## 4. HyperLogLog (`PFADD`, `PFCOUNT`)

Imagine counting the exact number of Unique Visitors to your website today.
If 100 million unique IP addresses visit, storing them all in a standard Redis Set (`SADD`) would consume 12 Gigabytes of RAM.

**HyperLogLog** is a probabilistic data structure. It estimates the unique count with an error margin of ~0.8%. 
The magic? It uses a maximum of **12 Kilobytes of RAM**, regardless of whether you count 1 thousand or 1 billion unique IPs!

```go
// Add IP addresses to the HyperLogLog
rdb.PFAdd(ctx, "unique_visitors_today", "192.168.1.1", "10.0.0.5")

// Get the approximate count!
count, _ := rdb.PFCount(ctx, "unique_visitors_today").Result()
```
Enterprise companies use HyperLogLog for analytics, massive metric aggregations, and unique view counts on YouTube videos!
