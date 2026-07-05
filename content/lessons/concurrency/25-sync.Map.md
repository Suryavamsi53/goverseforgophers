# sync.Map

As mentioned earlier, standard Go maps (`map[string]int`) are NOT thread-safe. If two Goroutines try to write to a standard map simultaneously, the Go Runtime will instantly kill the application with a `fatal error: concurrent map writes`.

The standard solution is to wrap the map in a `sync.RWMutex`. 
However, for highly specific enterprise workloads, wrapping a map in a Mutex is too slow. For these workloads, Go provides `sync.Map`.

## 1. What is sync.Map?

`sync.Map` is an optimized, lock-free (mostly) concurrent map built directly into the standard library.

Unlike a standard map, you do not declare its types. It uses `any` (empty interfaces) for both Keys and Values, meaning you must use Type Assertions when reading data out of it.

## 2. Syntax

```go
var m sync.Map

// 1. Write Data
m.Store("user_42", "John Doe")

// 2. Read Data
val, ok := m.Load("user_42")
if ok {
    // Must Type Assert from `any` back to `string`
    fmt.Println(val.(string)) 
}

// 3. Delete Data
m.Delete("user_42")

// 4. Load or Store (Atomic Upsert)
// If the key exists, it returns the existing value.
// If it doesn't exist, it stores the new value.
actual, loaded := m.LoadOrStore("user_99", "Jane Doe")
```

## 3. Iteration

You cannot use a standard `for range` loop on a `sync.Map`. Instead, you pass a callback function to the `Range` method.

```go
m.Range(func(key, value any) bool {
    fmt.Printf("Key: %v, Value: %v\n", key, value)
    
    // Return true to continue iterating. 
    // Return false to break the loop.
    return true 
})
```

## 4. When to actually use sync.Map

The official Go documentation explicitly states that you should **NOT** use `sync.Map` for general-purpose programming! 

A standard map wrapped in a `sync.RWMutex` is generally much faster and provides compile-time type safety. 

**You should ONLY use `sync.Map` in two highly specific scenarios:**

1. **Append-Only Caches**: When the map is written to once, but read millions of times by thousands of Goroutines (e.g., caching a static configuration). `sync.Map` internally uses `sync/atomic` for reads, completely avoiding Mutex bottlenecks.
2. **Disjoint Key Sets**: When multiple Goroutines are reading/writing to the map, but they are all dealing with *completely different keys*. (e.g., Goroutine A only manages `user_1`, Goroutine B only manages `user_2`).

If you have multiple Goroutines constantly updating the *same* key, a standard `sync.Mutex` is significantly faster!
