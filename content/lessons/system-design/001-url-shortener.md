# System Design: URL Shortener

## 1. Learning Objectives
* **What you'll learn**: How to design a scalable URL Shortener (like bit.ly) focusing on API design, Base62 encoding, database indexing, and caching.
* **Why it matters**: It is the most common System Design interview question. It perfectly demonstrates your ability to handle read-heavy architectures and ID generation at scale.
* **Where it's used**: Bitly, TinyURL, Twitter (t.co), and internal deep-linking services at major corporations.

---

## 2. Real-world Story
Imagine a massive library with 10 million books. Every book has a unique, 500-character description of exactly where it is located ("Third Floor, Section A, Row 9, Shelf 2..."). This is a Long URL.
Instead of forcing visitors to memorize that, you give them a tiny ticket: `A1`. 
When they hand `A1` to the librarian (The Go Server), the librarian instantly looks at a Rolodex (Database), finds the 500-character description, and walks them to the book.

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[User (Browser)] -->|POST /api/shorten?url=...| B(API Gateway)
    B --> C[Go Shortener Service]
    
    C -->|1. Generate ID| D{KGS: Key Generation Service}
    C -->|2. Save Mapping| E[(PostgreSQL)]
    
    A2[User (Browser)] -->|GET /A1b2C| B
    B --> C
    C -->|3. Check Cache| F[(Redis Cache)]
    F -.->|Cache Miss| E
    C -->|4. HTTP 301 Redirect| A2
    
    style E fill:#3b82f6,color:#fff
    style F fill:#dc2626,color:#fff
```

---

## 4. Internal Working (Under the Hood)
The core of a URL shortener is a highly efficient Key-Value mapping. 
* **Write Path**: Accept a long URL, generate a unique 7-character string, and save `[ShortHash -> LongURL]` to the database.
* **Read Path**: Accept the 7-character string, query the database, and return an `HTTP 301 Moved Permanently` (or `302 Found`).

---

## 5. Compiler Behavior
* **Base62 Encoding**: Base62 consists of `[a-z, A-Z, 0-9]`. Using 7 characters of Base62 yields `62^7 = 3.5 Trillion` possible URLs. In Go, generating this string is extremely fast because you are doing mathematical modulo division on a 64-bit integer (the DB auto-increment ID) and looking up the characters from a predefined byte array constant without allocating dynamic memory.

---

## 6. Memory Management
* **The Redis Cache**: URL Shorteners are insanely Read-Heavy (100 Reads for every 1 Write). A viral tweet might send 50,000 requests per second to the exact same short URL. Hitting PostgreSQL for this will melt the database. You MUST implement a Redis Cache (Cache-Aside pattern). Go's memory footprint remains tiny because the hot URLs are fetched from Redis in 1ms.

---

## 7. Code Examples

### 🔹 Example 1: Base62 Encoder in Go
```go
const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const base = uint64(len(alphabet))

// Converts a sequential Database ID (e.g., 1000000000) into a Base62 string (e.g., "aUKYO")
func EncodeBase62(id uint64) string {
    if id == 0 { return string(alphabet[0]) }
    
    var bytes []byte
    for id > 0 {
        rem := id % base
        bytes = append(bytes, alphabet[rem])
        id = id / base
    }
    
    // Reverse the slice (optional, for aesthetics)
    for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
        bytes[i], bytes[j] = bytes[j], bytes[i]
    }
    return string(bytes)
}
```

### 🔹 Example 2: The Write API (Shortening)
```go
func ShortenURL(db *sql.DB, longURL string) (string, error) {
    // 1. Insert the Long URL and get the Auto-Incremented ID
    var id uint64
    err := db.QueryRow("INSERT INTO urls (long_url) VALUES ($1) RETURNING id", longURL).Scan(&id)
    if err != nil { return "", err }
    
    // 2. Convert ID to Base62 Short URL
    shortKey := EncodeBase62(id)
    
    // 3. (Optional) Save the shortKey back to the DB for faster lookups
    return "https://go.ly/" + shortKey, nil
}
```

### 🔹 Example 3: Advanced (Key Generation Service - KGS)
```go
// Relying on Postgres Auto-Increment is a bottleneck at 10,000 Writes/Sec.
// Instead, build a standalone Go KGS that pre-generates 1 million Base62 keys
// and stores them in memory. The API just pops a key from a Go Channel instantly!
var keyPool chan string = make(chan string, 1000)

func GetPreGeneratedKey() string {
    return <-keyPool // Nanosecond latency, completely bypasses the DB on Write!
}
```

### 🔹 Example 4: Production (The Redirect Handler)
```go
func RedirectHandler(w http.ResponseWriter, r *http.Request) {
    shortKey := chi.URLParam(r, "key")
    
    // 1. Check Redis First!
    longURL, err := redisClient.Get(ctx, shortKey).Result()
    if err == redis.Nil {
        // 2. Cache Miss: Check Postgres
        longURL = FetchFromDB(shortKey)
        // 3. Save to Redis for 1 Hour
        redisClient.Set(ctx, shortKey, longURL, time.Hour)
    }
    
    // HTTP 302: Tells the browser to go to the Long URL, but don't cache the redirect!
    // (If you use 301, the browser caches it forever, and you can't track analytics!)
    http.Redirect(w, r, longURL, http.StatusFound)
}
```

### 🔹 Example 5: Interview
```go
// Q: What happens if two users submit the exact same Long URL? 
// A: It depends on product requirements! If we want to save space, we check the DB first and 
// return the existing short URL. If we want per-user analytics (e.g. tracking click rates 
// for Alice vs Bob), we MUST generate two unique short URLs!
```

---

## 8. Production Examples
1. **Analytics Engine**: When a user clicks a short link, the Go Redirect Handler asynchronously drops an event into an Apache Kafka topic before responding with the `302 Redirect`. A separate Go worker processes the Kafka events to update real-time click charts.
2. **Malware Scanning**: Before saving a Long URL to the database, the API makes a gRPC call to a Security Microservice (using the Google Safe Browsing API) to block phishing links.

---

## 9. Performance & Benchmarking
* **Data Storage Math**: If you generate 100 million URLs per month, and each record is 500 Bytes (ID + LongURL + ShortURL + CreatedAt). That is `100M * 500B = 50 GB / month`. Over 5 years, the database will grow to 3 TB. PostgreSQL can easily handle a 3TB table provided you place a strict B-Tree Index on the `short_url` column!

---

## 10. Best Practices
* ✅ **Do**: Use a NoSQL database (like DynamoDB or Cassandra) if you need extreme horizontal write scaling and don't care about complex relational joins.
* ❌ **Don't**: Use MD5 or SHA-256 to hash the long URL. Hashing generates a 64-character string. If you truncate it to 7 characters, you will suffer massive Hash Collisions. Base62 Encoding from a unique Counter mathematically guarantees zero collisions!
* 🏢 **Google / Uber / Netflix Style**: Use ZooKeeper or `etcd` to manage the distributed counter. If you have 5 Go API servers, Server A gets IDs 1-1,000,000. Server B gets 1,000,001-2,000,000. They can encode Base62 purely in local RAM without ever causing database write conflicts!

---

## 11. Common Mistakes
1. **Missing Rate Limiting**: URL Shorteners are public APIs. Malicious actors will write a script to shorten 10 million URLs, filling up your database and costing you thousands of dollars. You must implement IP-based rate limiting (e.g., Redis Token Bucket).
2. **HTTP 301 vs 302**: Using `301 Moved Permanently`. The user's browser will permanently cache the redirect. The next time they click the short link, their browser will route them directly to the destination without hitting your Go server. You will lose all click analytics! Use `302 Found`.

---

## 12. Debugging
How to troubleshoot a URL Shortener:
* **Cache Stampede**: If a viral URL expires from Redis, 10,000 Go Goroutines might query Postgres simultaneously. You must use Go's `golang.org/x/sync/singleflight` to ensure only ONE Goroutine queries the database, and the other 9,999 wait for the result!

---

## 13. Exercises
1. **Easy**: Write the `EncodeBase62(id uint64)` function in Go and test it with ID `999999`.
2. **Medium**: Spin up a local SQLite database and write the `/shorten` and `/{key}` endpoints using `net/http`.
3. **Hard**: Integrate a Redis container via Docker Compose and implement the Cache-Aside pattern on the Redirect endpoint.
4. **Expert**: Implement the `singleflight` package to protect the database during a Cache Miss on a highly concurrent viral link.

---

## 14. Quiz
1. **MCQ**: What is the most robust way to guarantee the 7-character string is 100% unique across a distributed system?
   * (A) Truncate an MD5 hash (B) Generate a random string and query the DB to see if it exists (C) Use a centralized auto-incrementing ID/Counter and convert to Base62. *(Answer: C)*
2. **System Design Follow-up**: How would you expire links after 30 days? *(If using Cassandra/DynamoDB, use native TTL features. If using Postgres, run a background Go worker every night that runs `DELETE FROM urls WHERE created_at < NOW() - INTERVAL '30 days'` in chunks).*

---

## 15. FAANG Interview Questions
* **Beginner**: Calculate the number of possible combinations for a 7-character Base62 string.
* **Intermediate**: Discuss the pros and cons of using an RDBMS (Postgres) vs NoSQL (Cassandra) for this specific problem.
* **Senior (Google/Meta)**: Architect the Key Generation Service (KGS). How does KGS guarantee it doesn't give the same ID to two different API servers if the KGS node crashes mid-assignment?

---

## 16. Mini Project
**The Go.ly Service**
* Write a Go HTTP server using `chi` router.
* Create a `sync.Mutex` protected global `counter uint64` to simulate the DB ID generator.
* Implement `/api/shorten` (returns JSON).
* Implement `/{key}` (returns HTTP 302 Redirect).
* Benchmark the redirect endpoint using `wrk -c 100 -d 10s http://localhost:8080/A1b2C`. You should hit >50,000 req/sec!

---

## 17. Enterprise Features & Observability
* **Custom Aliases**: Enterprise customers want `https://go.ly/MyBrand`. You must update the DB schema to accept a `custom_alias` field. When generating a link, check if `custom_alias` is provided; if so, validate its uniqueness before saving!

---

## 18. Source Code Reading
Walkthrough of `math/rand` vs `crypto/rand`.
* **Predictability**: If you decide to use random string generation instead of a Counter, you MUST use `crypto/rand`. `math/rand` is seeded predictably. Hackers can figure out the seed and guess every single private shortened URL generated by your system!

---

## 19. Architecture
* **Read Heavy Scaling**: Because the system is 100:1 Reads:Writes, you deploy 1 Postgres Primary (for writes) and 10 Postgres Read Replicas. The Go redirect handlers only ever query the Read Replicas!

---

## 20. Summary & Cheat Sheet
* **Encoding**: Base62 (a-z, A-Z, 0-9).
* **ID Generation**: Distributed Counter (ZooKeeper/etcd) -> Base62.
* **Database**: High read throughput (NoSQL or Read Replicas).
* **Caching**: Redis is mandatory for hot links.
* **Redirection**: HTTP 302 to track analytics.
