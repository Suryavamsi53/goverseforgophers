# Performance Optimization

Go is naturally fast, but poorly written Go can generate massive Garbage Collection (GC) pauses that ruin performance. 

Here are the three most critical techniques for optimizing Go applications.

## 1. Slice Pre-allocation

When you `append()` to a slice, and the slice runs out of capacity, Go creates a brand new, larger array on the Heap, copies the data over, and abandons the old array for the Garbage Collector.

If you know roughly how many items you need, **always pre-allocate capacity.**

**❌ BAD (Triggers multiple memory re-allocations):**
```go
var users []User
for _, id := range userIDs {
    // Slice doubles in size repeatedly, trashing the Heap
    users = append(users, fetchUser(id)) 
}
```

**✅ GOOD (Zero re-allocations):**
```go
// Create a slice with Length 0, but Capacity equal to the exact size we need!
users := make([]User, 0, len(userIDs)) 
for _, id := range userIDs {
    // Appends perfectly without ever resizing the underlying array!
    users = append(users, fetchUser(id))
}
```

## 2. Escape Analysis and Structs

As covered in the Memory Layout chapters, passing large structs by pointer `*User` avoids copying data, but forces the struct onto the Heap (triggering GC). Passing by value `User` keeps it on the ultra-fast Stack, but copies memory.

* **Rule of Thumb**: For small structs (under ~64 bytes, like a UUID or a Coordinate), pass by Value. The CPU can copy it faster than the GC can clean it up.
* **Rule of Thumb**: For massive structs, or when mutation is required, pass by Pointer.

*Always use `go build -gcflags="-m"` to verify your assumptions!*

## 3. Bypassing GC: `sync.Pool`

If your web server allocates a 32KB `bytes.Buffer` for every single incoming HTTP request, a high-traffic server will allocate gigabytes of RAM every second, causing the Garbage Collector to consume 100% of your CPU.

You can bypass the Garbage Collector entirely using an Object Pool (`sync.Pool`). 

A Pool keeps a cache of instantiated objects. When you need a buffer, you "borrow" it from the pool. When you are done, you "return" it to the pool to be reused by the next web request.

```go
var bufferPool = sync.Pool{
    // Define how to create a new object if the pool is empty
    New: func() any {
        return new(bytes.Buffer) 
    },
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // 1. Borrow a buffer from the pool (Zero Allocation!)
    buf := bufferPool.Get().(*bytes.Buffer)
    
    // 2. CRITICAL: Reset the buffer before returning it so the next user doesn't see your data!
    defer func() {
        buf.Reset()
        bufferPool.Put(buf) // Return it to the pool
    }()

    // ... Use the buffer ...
}
```
`sync.Pool` is the secret weapon used by the standard library's `fmt` and `json` packages to achieve their blazing speeds.
