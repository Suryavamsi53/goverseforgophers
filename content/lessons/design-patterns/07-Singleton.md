# Singleton Pattern

The Singleton Pattern ensures that a class has only one instance globally, and provides a global point of access to it.

In enterprise systems, Singletons are frequently used for database connection pools, global configuration structs, or logger instances. If 10,000 HTTP requests try to open a new database connection simultaneously, the database will crash. You must ensure they all share the exact same database connection instance.

## 1. The Thread-Safety Problem

If you implement a Singleton naïvely, it will cause a catastrophic **Data Race**.

```go
// DANGER: NOT THREAD-SAFE!
var instance *Database

func GetDatabase() *Database {
    // If 10 Goroutines hit this line at the exact same nanosecond,
    // they will all see instance == nil, and they will ALL initialize a new Database!
    if instance == nil {
        instance = ConnectToDB() 
    }
    return instance
}
```

## 2. The Mutex Solution (Slow)

To fix the Data Race, you can wrap the check in a `sync.Mutex`.

```go
var mu sync.Mutex

func GetDatabase() *Database {
    mu.Lock()
    defer mu.Unlock()
    
    if instance == nil {
        instance = ConnectToDB()
    }
    return instance
}
```
**The Flaw**: This works, but it destroys performance. Every single time *any* Goroutine in your application wants to run a SQL query, it must wait in line to acquire this Mutex lock just to get the pointer to the database! 

## 3. The `sync.Once` Solution (The Go Standard)

As we covered in the Concurrency module, the idiomatic and highly optimized way to implement a Singleton in Go is using `sync.Once`.

```go
package database

import "sync"

type Database struct { }

var (
    instance *Database
    once     sync.Once
)

// GetInstance guarantees thread-safe, one-time initialization.
func GetInstance() *Database {
    // once.Do uses a fast-path atomic check under the hood.
    // If the instance already exists, it bypasses the Mutex entirely and returns instantly!
    once.Do(func() {
        instance = &Database{}
        // Perform expensive connection logic here...
    })
    
    return instance
}
```

## 4. Why Singletons are Anti-Patterns

While `sync.Once` makes Singletons fast and safe, the Singleton pattern itself is widely considered an **Anti-Pattern** in modern software engineering.

Why? Because it relies on **Global State**.

If your `OrderService` calls `database.GetInstance()` directly inside its functions, it is tightly coupled to the global database. 
* You cannot easily swap the database for a different one.
* **Unit Testing is impossible**: You cannot easily Mock the database, because the service reaches out to the global state instead of accepting it as a dependency.

### The Solution: Dependency Injection
Instead of using the Singleton pattern, instantiate the database exactly once in your `main.go` file, and pass it explicitly as a pointer to the services that need it.

```go
func main() {
    // Instantiate exactly once here
    db := database.Connect()
    
    // Pass it explicitly! No global Singletons needed!
    orderService := NewOrderService(db)
}
```
