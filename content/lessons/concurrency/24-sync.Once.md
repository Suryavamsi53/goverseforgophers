# sync.Once

Imagine you have a Go application that connects to a database. 
You only want to initialize that database connection exactly **once** when the application starts. However, you might have 10 different HTTP handlers running concurrently that all try to fetch the database connection at the exact same time.

If all 10 handlers detect that `db == nil`, they might all try to initialize the database simultaneously, opening 10 redundant connection pools!

We can solve this perfectly with `sync.Once`.

## 1. Syntax

`sync.Once` guarantees that a given function will only ever be executed exactly one time, regardless of how many Goroutines call it simultaneously.

```go
package main

import (
    "fmt"
    "sync"
)

var (
    once sync.Once
    db   *Database
)

func GetDatabase() *Database {
    // If 1,000 Goroutines hit this line at the exact same nanosecond,
    // sync.Once will block 999 of them. 
    // It allows 1 to execute the function, and then instantly returns for the other 999!
    once.Do(func() {
        fmt.Println("Initializing Database Connection...")
        db = ConnectToPostgres()
    })
    
    return db
}
```

## 2. How it works (Under the hood)

You might think `sync.Once` is just a Mutex wrapped around a boolean flag:
```go
// NAÏVE IMPLEMENTATION (Slow!)
mu.Lock()
if !initialized {
    initialize()
    initialized = true
}
mu.Unlock()
```
The problem with the naïve approach is that *every single time* someone calls `GetDatabase()` for the entire lifespan of the application, it requires a Mutex lock! This ruins performance.

Instead, the Go Runtime implements `sync.Once` using **Atomic Operations** (Fast Path) and **Mutexes** (Slow Path).

1. **Fast Path**: It uses `atomic.Load` to check a `done` integer flag. (10x faster than a Mutex).
2. If `done == 0`, it falls back to the **Slow Path**: it acquires a Mutex, double-checks the flag, runs the function, and then uses `atomic.Store` to set the flag to `1`.
3. For the rest of the application's lifespan, all 10,000 Goroutines will only hit the atomic Fast Path, experiencing zero lock contention!

## 3. The Deadlock Trap

Because `sync.Once` holds a Mutex during the execution of the function, you must NEVER call `once.Do` again from *inside* the `once.Do` function.

```go
var once sync.Once

func Setup() {
    once.Do(func() {
        // FATAL CRASH: Deadlock!
        // The first once.Do holds the lock.
        // The second once.Do tries to acquire the lock and freezes forever.
        once.Do(func() {
            fmt.Println("Nested call")
        })
    })
}
```

## 4. sync.OnceValues (Go 1.21+)

Historically, `sync.Once` didn't return anything, forcing developers to use global variables (like `var db *Database`).

In Go 1.21, the Go team introduced `sync.OnceValue` and `sync.OnceValues` which return the initialized data and an error!

```go
var getDB = sync.OnceValues(func() (*Database, error) {
    return ConnectToPostgres()
})

func HandleRequest() {
    // Clean, localized, and thread-safe!
    db, err := getDB()
}
```
