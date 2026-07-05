# Decorator Pattern

The Decorator Pattern allows you to attach new behaviors to an object dynamically by placing it inside a wrapper object that contains the new behaviors.

In Go, this is heavily used to implement **Middleware** for HTTP servers, gRPC interceptors, and adding caching/logging layers to databases without modifying the original code.

## 1. The Problem

You have a `UserService` that fetches a user from a database.

```go
type UserService interface {
    GetUser(id int) string
}

type PostgresUser struct{}
func (p *PostgresUser) GetUser(id int) string {
    return "User Data from DB"
}
```

Your manager asks you to add **Logging** (so we know when a user is fetched) and **Caching** (so we don't hit the DB too often).

The bad approach is to modify the `PostgresUser` struct and hardcode a Redis cache and a Logger directly into the SQL query logic. This violates the Single Responsibility Principle.

## 2. The Solution (The Decorator)

Instead, we create new structs that *wrap* the `UserService`.

### The Logging Decorator
```go
type LoggingDecorator struct {
    service UserService // The inner object we are wrapping
}

func (l *LoggingDecorator) GetUser(id int) string {
    fmt.Printf("[LOG] Fetching user %d\n", id)
    
    // Call the inner object
    result := l.service.GetUser(id)
    
    fmt.Printf("[LOG] Fetched user %d successfully\n", id)
    return result
}
```

### The Caching Decorator
```go
type CachingDecorator struct {
    service UserService
    cache   map[int]string
}

func (c *CachingDecorator) GetUser(id int) string {
    // 1. Try to return from cache
    if val, ok := c.cache[id]; ok {
        fmt.Println("[CACHE] Returning cached data!")
        return val
    }
    
    // 2. Cache miss! Call the inner object
    result := c.service.GetUser(id)
    
    // 3. Save to cache and return
    c.cache[id] = result
    return result
}
```

## 3. The Usage (Russian Nesting Dolls)

Because both Decorators implement the `UserService` interface, and both *accept* a `UserService` interface, you can chain them together like Russian nesting dolls!

```go
func main() {
    // 1. The core implementation
    var core UserService = &PostgresUser{}
    
    // 2. Wrap it in a Cache
    var cached UserService = &CachingDecorator{
        service: core,
        cache:   make(map[int]string),
    }
    
    // 3. Wrap the Cache in a Logger!
    var fullyDecorated UserService = &LoggingDecorator{
        service: cached,
    }
    
    // The caller has no idea this is decorated. It just calls the interface!
    // Execution: Logger -> Cache -> Postgres -> Cache -> Logger
    fmt.Println(fullyDecorated.GetUser(1))
    fmt.Println(fullyDecorated.GetUser(1)) // Second call hits the cache!
}
```

## 4. HTTP Middleware (The Ultimate Go Decorator)

If you have ever written a web server in Go using `http.HandlerFunc`, you have used the Decorator pattern! 
Middleware is literally just a function that takes an `http.Handler`, wraps it with logging/auth logic, and returns a new `http.Handler`.
