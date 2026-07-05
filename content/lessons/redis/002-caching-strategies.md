# Caching Strategies & The Stampede

Redis is primarily used as a caching layer in front of a slow primary database (like PostgreSQL) to protect it from extreme read traffic.

There are specific architectural patterns you must follow when synchronizing data between Redis and Postgres.

## 1. The Cache-Aside Pattern (Lazy Loading)

This is the most common pattern in the industry. The Go application acts as the mediator between Redis and Postgres.

1. **Read Request**:
   * The Go app asks Redis: `GET user_42`.
   * **Cache Hit**: Redis returns the data. We return `200 OK`. (Postgres is never touched).
   * **Cache Miss**: Redis returns `redis.Nil`. The Go app queries Postgres. The Go app saves the result to Redis (`SET user_42 "..." EX 3600`), and returns `200 OK`.
2. **Write Request**:
   * The Go app `UPDATE`s the data in Postgres.
   * The Go app `DEL`s the key in Redis (Invalidation). The next Read will trigger a Cache Miss and fetch the fresh data!

## 2. The Cache Stampede (Thundering Herd)

Imagine an E-Commerce site on Black Friday. The "Home Page Products" key is cached in Redis for 10 minutes.
At exactly 12:10 PM, the key expires (TTL = 0) and Redis deletes it.

In the next 100 milliseconds, 5,000 users hit the Home Page.
All 5,000 Go Goroutines check Redis simultaneously. They all get a **Cache Miss**.
All 5,000 Goroutines instantly send a massive, complex `SELECT` query to PostgreSQL at the exact same millisecond. 
PostgreSQL instantly crashes.

This is the **Cache Stampede**.

### The Solution: Mutex Locking (Singleflight)
When a Cache Miss occurs, you cannot let 5,000 Goroutines hit the database. You must force 4,999 of them to wait, while exactly 1 Goroutine fetches the data.

In Go, we use the `golang.org/x/sync/singleflight` package!

```go
var g singleflight.Group

func GetHomePageData() string {
    // 1. Check Redis
    val, err := rdb.Get(ctx, "home_page").Result()
    if err == nil { return val } // Cache Hit!

    // 2. Cache Miss! Use Singleflight!
    // The first Goroutine executes the function. Any other Goroutine that hits 
    // this block with the key "fetch_home" will BLOCK and wait for the first one to finish!
    v, err, _ := g.Do("fetch_home", func() (interface{}, error) {
        
        // 3. Query Postgres (Only happens ONCE!)
        data := queryPostgres()
        
        // 4. Save back to Redis
        rdb.Set(ctx, "home_page", data, time.Minute*10)
        return data, nil
    })
    
    return v.(string)
}
```

## 3. Write-Through and Write-Behind Caching

* **Write-Through**: When the Go app receives an `UPDATE`, it synchronously updates Postgres, and then synchronously updates Redis in the same request. Data is always perfectly in sync, but writes are slower.
* **Write-Behind (Extreme Speed)**: The Go app receives an `UPDATE`, updates Redis, and immediately returns `200 OK`. An asynchronous Goroutine (or Kafka worker) updates Postgres 5 seconds later. Writes are blazingly fast, but if the Redis server crashes before the background worker runs, the update is lost permanently!
