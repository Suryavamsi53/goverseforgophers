# Pagination, Filtering, and Sorting

If you have 10 million users in your PostgreSQL database, and a React frontend makes a request to `GET /users`, what happens?

The Go server queries `SELECT * FROM users`, loads 10 million rows into RAM, serializes a 5-Gigabyte JSON string, and attempts to send it over the network. The Go server crashes with an OOM (Out of Memory) error, and the user's browser freezes.

APIs must restrict the amount of data returned using **Pagination**.

## 1. Offset-Based Pagination (The Standard)

The most common way to paginate is using Query Parameters: `?limit=X&offset=Y`.

* `GET /users?limit=50&offset=0` (Returns users 1 to 50)
* `GET /users?limit=50&offset=50` (Returns users 51 to 100)

**The Go / SQL Implementation:**
```sql
-- This translates perfectly to standard SQL!
SELECT * FROM users ORDER BY created_at DESC LIMIT 50 OFFSET 50;
```

**The Flaw of Offset Pagination:**
If you have a massive table, `OFFSET 1000000` is incredibly slow in PostgreSQL. The database has to physically scan and discard the first 1 million rows before returning the next 50. 
Additionally, if new users are inserted into the database while the client is clicking "Next Page", items will physically shift in the database, causing the client to see duplicate users across pages!

## 2. Cursor-Based Pagination (Extreme Scale)

For massive feeds (like Twitter or Facebook), you cannot use Offsets. You must use a **Cursor**.

A cursor is a unique identifier (usually a timestamp or a Snowflake ID) pointing to the exact last item the client saw.

* `GET /users?limit=50` (Returns the first 50 users). The API includes a `next_cursor` in the JSON response (e.g., `next_cursor=1700000000`).
* `GET /users?limit=50&cursor=1700000000`

**The Go / SQL Implementation:**
```sql
-- Blazingly fast! Uses a B-Tree index on created_at to instantly jump to the spot!
SELECT * FROM users 
WHERE created_at < 1700000000 
ORDER BY created_at DESC 
LIMIT 50;
```
Cursor pagination guarantees $O(log N)$ performance, and prevents duplicate items if the database shifts. However, you cannot implement "Jump to Page 5", you can only implement "Next/Previous".

## 3. Filtering and Searching

To filter resources, you append standard query parameters.
* `GET /users?role=admin&active=true`

If you need complex search operations (like searching across multiple fields), it is an industry standard to use a `q` parameter.
* `GET /users?q=suryavamsi`

## 4. Sorting

To allow the client to sort the data, use a `sort` parameter. The industry standard is to use a minus sign (`-`) prefix to denote descending order.

* `GET /users?sort=-created_at` (Sort by newest first).
* `GET /users?sort=age,-created_at` (Sort by age ascending, then created_at descending).

**Security Warning:**
In Go, NEVER directly concatenate the `sort` query parameter into your SQL `ORDER BY` clause! This is a massive SQL Injection vulnerability. 
You must validate the `sort` parameter against a hardcoded whitelist of allowed columns in your Go code before appending it to the SQL string!
