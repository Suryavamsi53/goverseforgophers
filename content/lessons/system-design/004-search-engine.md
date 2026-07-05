# System Design: Search Engine (Elasticsearch / Inverted Indexes)

## 1. Learning Objectives
* **What you'll learn**: How to design a scalable Search Engine architecture (like e-commerce product search) by understanding the mechanics of Inverted Indexes, Tokenization, and Elasticsearch integration in Go.
* **Why it matters**: `SELECT * FROM products WHERE description LIKE '%apple%'` is a Full Table Scan. It takes 10 seconds on a 10-million row database. Users expect search results in 50 milliseconds. Traditional databases cannot do this.
* **Where it's used**: Amazon Product Search, Twitter Keyword Search, Log Analysis (Kibana), and any application requiring fuzzy text matching.

---

## 2. Real-world Story
Imagine a traditional database like a book. To find every page that mentions the word "apple", you must read the entire book from Page 1 to Page 500 (`Full Table Scan`). It is agonizingly slow.
A Search Engine is the **Index at the back of the book**. You flip to the back, look up the word "apple" alphabetically, and it instantly tells you: "Pages 12, 45, 99". You found the data in 1 second without reading the whole book. This is called an **Inverted Index**.

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Primary DB: PostgreSQL] -->|Change Data Capture (Kafka)| B(Go Sync Worker)
    
    B -->|Indexes Document| C[(Elasticsearch Cluster)]
    
    D[User Web Client] -->|GET /search?q=iphone| E[Go Search API]
    
    E -->|O 1 Lookup Query| C
    C -.->|Returns Doc IDs| E
    E -->|Transforms to JSON| D
    
    style C fill:#0f172a,color:#fff
```

---

## 4. Internal Working (Under the Hood)
When you insert a string like `"The Quick Brown Fox"` into Elasticsearch:
1. **Analyzer**: Lowercases the text (`the quick brown fox`).
2. **Tokenizer**: Splits it into words (`[the, quick, brown, fox]`).
3. **Filter**: Removes stop words (`the`) and stems words (e.g., `running` -> `run`).
4. **Inverted Index**: Stores the mapping:
   * `quick` -> Document ID 1
   * `brown` -> Document ID 1
When a user searches `"brown"`, the engine does an O(1) hash map lookup on the word "brown" and instantly returns Document ID 1!

---

## 5. Compiler Behavior
* **Memory Mapping (mmap)**: Lucene (the engine inside Elasticsearch) relies heavily on the Linux Kernel's page cache. It maps the index files directly into virtual memory. When a Go API queries Elasticsearch, the data is served directly from RAM (page cache) at blistering speeds, entirely bypassing disk I/O!

---

## 6. Memory Management
* **Denormalization**: Elasticsearch does NOT support SQL `JOIN`s efficiently. If a Product has an Author, you cannot store them in two different tables and join them at search time. You must flatten the data into a massive JSON object (Denormalization) before the Go worker sends it to Elasticsearch. This uses more Disk Space but optimizes CPU for extreme read speeds.

---

## 7. Code Examples

### 🔹 Example 1: The Go Sync Worker (Indexing)
```go
import "github.com/elastic/go-elasticsearch/v8"

func IndexProductToElastic(es *elasticsearch.Client, product Product) {
    // Flatten the data into a simple JSON Document!
    body, _ := json.Marshal(map[string]interface{}{
        "title":       product.Title,
        "description": product.Description,
        "category":    product.Category.Name, // Denormalized!
    })

    // Upsert the document into the "products" index
    req := esapi.IndexRequest{
        Index:      "products",
        DocumentID: product.ID,
        Body:       bytes.NewReader(body),
        Refresh:    "true", // Make it searchable instantly
    }
    req.Do(context.Background(), es)
}
```

### 🔹 Example 2: The Search API (Fuzzy Matching)
```go
// Allows users to misspell "Iphone" as "Iphon" and still find it!
func SearchProducts(es *elasticsearch.Client, query string) {
    // Construct the Elasticsearch JSON Query DSL
    searchQuery := map[string]interface{}{
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                "title": map[string]interface{}{
                    "query":     query,
                    "fuzziness": "AUTO", // The magic of Fuzzy search!
                },
            },
        },
    }
    
    body, _ := json.Marshal(searchQuery)
    es.Search(
        es.Search.WithContext(context.Background()),
        es.Search.WithIndex("products"),
        es.Search.WithBody(bytes.NewReader(body)),
    )
}
```

### 🔹 Example 3: Advanced (Relevance Scoring)
```json
// Boosting fields! A match in the "Title" is worth 3x more than a match in the "Description"
{
  "query": {
    "multi_match": {
      "query": "iphone",
      "fields": ["title^3", "description"] 
    }
  }
}
```

### 🔹 Example 4: Production (Pagination)
```go
// Elasticsearch pagination (From and Size). 
// NEVER use From > 10,000 (Deep Pagination) as it will crash the cluster! 
// Use 'search_after' for deep scrolling.
es.Search(
    es.Search.WithFrom(0),
    es.Search.WithSize(20),
    // ...
)
```

### 🔹 Example 5: Interview
```go
// Q: How do you keep the primary PostgreSQL database and Elasticsearch in sync?
// A: Avoid dual-writes in Go. Use the Outbox Pattern or Change Data Capture (Debezium). 
// Debezium tails the Postgres WAL and pushes changes to Kafka. A Go worker reads Kafka 
// and pushes to Elasticsearch, mathematically guaranteeing eventual consistency.
```

---

## 8. Production Examples
1. **Typeahead / Autocomplete**: When you type "Macb..." in the search bar, the API returns suggestions instantly. This is powered by an `Edge N-Gram` tokenizer. It breaks the word into `[m, ma, mac, macb]` and indexes them all, so a prefix search is lightning fast.
2. **Geospatial Search**: "Find restaurants near me". Elasticsearch natively supports `geo_point` data types and can instantly calculate distances using BKD trees on a map index.

---

## 9. Performance & Benchmarking
* **Sharding**: A single search index can grow to 10 Terabytes. Elasticsearch splits the index into "Shards". Shard A lives on Server 1, Shard B on Server 2. When the Go API sends a query, Elasticsearch runs the search on Server 1 and Server 2 simultaneously (Scatter-Gather Phase) and merges the top results!

---

## 10. Best Practices
* ✅ **Do**: Use Bulk APIs. When migrating 1 million rows from Postgres to Elasticsearch, don't send 1 million HTTP requests. Use the `_bulk` endpoint in Go to send 5,000 documents in a single HTTP payload.
* ❌ **Don't**: Store large binary blobs (like images) in Elasticsearch. Store them in S3, and just store the S3 URL string in the Search Index.
* 🏢 **Google / Uber / Netflix Style**: Use an API Gateway to route write requests (`POST /products`) to a Write Service (Postgres), and route read requests (`GET /search`) to a dedicated Read Service backed entirely by Elasticsearch (CQRS Pattern).

---

## 11. Common Mistakes
1. **The Deep Pagination Crash**: A malicious user requests `?page=1000000`. To find the millionth result, Elasticsearch has to load all 1,000,000 results into RAM, sort them, and return the last 10. This will OOM the cluster. Always enforce a hard limit (`max_result_window = 10000`).
2. **Mapping Explosions**: Allowing dynamic field mapping. If users insert documents with 10,000 unique JSON keys, Elasticsearch creates an index for every key and crashes. Enforce strict schemas (Mappings) for your indexes!

---

## 12. Debugging
How to troubleshoot Search queries:
* **The `_analyze` API**: If a product isn't showing up in search, you can ask Elasticsearch exactly how it tokenized the text. It will reveal if a Stop Word filter accidentally deleted the crucial search keyword!

---

## 13. Exercises
1. **Easy**: Spin up Elasticsearch locally using Docker.
2. **Medium**: Write a Go script to index 3 JSON documents (Products).
3. **Hard**: Write a Go API that accepts a `?q=` query parameter, performs a `multi_match` fuzzy search on the indexed documents, and returns the JSON.
4. **Expert**: Implement an Edge N-Gram analyzer in your Index Mapping to support real-time "Typeahead" autocomplete functionality.

---

## 14. Quiz
1. **MCQ**: What is the core data structure that enables fast full-text search?
   * (A) B-Tree (B) Hash Map (C) Inverted Index. *(Answer: C)*
2. **System Design Follow-up**: Why does indexing data in Elasticsearch take significant CPU power compared to inserting a row in Postgres? *(Because the Analyzer has to parse strings, remove punctuation, stem verbs, and update massive inverted index structures in RAM before flushing to disk).*

---

## 15. FAANG Interview Questions
* **Beginner**: Why can't we just use SQL `LIKE '%term%'`?
* **Intermediate**: Explain the concept of Term Frequency / Inverse Document Frequency (TF-IDF or BM25) for Relevance Scoring.
* **Senior (Google/Meta)**: Architect a global scale search index. If a node containing a primary shard crashes, how does the cluster maintain High Availability and prevent data loss during searches?

---

## 16. Mini Project
**The Book Finder**
* Download a CSV of 1,000 books.
* Write a Go script using `goroutines` and the Elasticsearch Bulk API to ingest the 1,000 books in under 1 second.
* Write a small Go HTML server with a Search Box.
* Connect the Search Box to an Elasticsearch backend query. Prove that searching for a misspelled title still returns the correct book!

---

## 17. Enterprise Features & Observability
* **Synonyms**: You can configure a Synonyms list (`sneakers => shoes`). If a user searches for "sneakers", the Go API doesn't need to change. Elasticsearch intercepts the query and automatically searches for "shoes" as well!

---

## 18. Source Code Reading
Walkthrough of `github.com/elastic/go-elasticsearch`.
* **The Typed vs Untyped API**: The V8 Go client provides both a raw JSON builder and a Strictly Typed API (`typedapi`). Study how the Typed API uses complex Go structs to statically validate your massive Elasticsearch JSON queries at compile time!

---

## 19. Architecture
* **The Write-Ahead Log (Translog)**: Like Postgres, Elasticsearch is durable. Before it builds the complex Inverted Index in RAM, it instantly writes the raw JSON to a sequential disk log (Translog). If the server loses power, it replays the Translog on boot.

---

## 20. Summary & Cheat Sheet
* **Core Tech**: Inverted Index.
* **Process**: Analyze -> Tokenize -> Index.
* **Strengths**: Lightning-fast fuzzy text and geospatial search.
* **Weaknesses**: Deep pagination and relational JOINs.
