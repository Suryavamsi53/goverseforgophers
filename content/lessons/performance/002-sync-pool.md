# Object Reuse with sync.Pool

## 1. Learning Objectives
* **What you'll learn**: How to eliminate memory allocations in hot paths by reusing objects with `sync.Pool`.
* **Why it matters**: Garbage Collection (GC) is the enemy of high throughput. If your HTTP server handles 10,000 requests per second and each request allocates a new `[]byte` buffer for JSON encoding, you are generating 10,000 objects per second that the GC must clean up. `sync.Pool` solves this.
* **Where it's used**: Inside high-performance libraries like `encoding/json`, `fmt`, and popular web frameworks like `gin` and `fiber`.

---

## 2. The Mechanics of sync.Pool
A `sync.Pool` is a thread-safe cache for temporary objects.
Instead of asking the OS for new memory (which is slow), you ask the Pool: *"Do you have an old object I can use?"*
When you are done with the object, instead of letting it die (and causing GC), you put it back in the Pool.

**Critical Rule**: Objects in a `sync.Pool` are automatically cleared during Garbage Collection. You cannot use it to store persistent data like database connections. It is strictly for *temporary* memory reuse.

---

## 3. Example: Zero-Allocation JSON Encoding

Imagine a high-traffic API that writes JSON logs for every request.

### The Bad Way (High Allocations)
```go
package main

import "bytes"

func handleRequest() {
    // A brand new buffer is allocated on the heap EVERY time
    buf := new(bytes.Buffer)
    buf.WriteString(`{"status": "success"}`)
    // Do something with buf...
    
    // buf dies here. The GC must clean it up later.
}
```

### The Fast Way (sync.Pool)
```go
package main

import (
    "bytes"
    "sync"
)

// 1. Create the pool and define how to make a NEW object if the pool is empty
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func handleRequest() {
    // 2. Get a buffer from the pool
    buf := bufferPool.Get().(*bytes.Buffer)
    
    // 3. ALWAYS reset the object before using it! It might have old data.
    buf.Reset() 
    
    buf.WriteString(`{"status": "success"}`)
    // Do something with buf...
    
    // 4. Put the buffer back so the next request can use it
    bufferPool.Put(buf)
}
```

In the Fast Way, if 10,000 concurrent requests hit your server, the pool will grow to 10,000 buffers. When the traffic drops, the GC will automatically clean up the unused buffers. You've eliminated constant allocation and deallocation!

---

## 4. Common Pitfalls

1. **Forgetting to Reset**: If you put a buffer back into the pool without resetting it, the next function that pulls it out will read your old data! Always call `.Reset()` (or clear the slice).
2. **Pooling Small Objects**: Do not put tiny structs in a `sync.Pool`. The overhead of the mutex inside `sync.Pool` is actually slower than just letting Go allocate a tiny struct on the Stack. Only pool large objects, large `[]byte` slices, or `bytes.Buffer`.

---

## 5. Quiz

1. **MCQ**: What happens to the items inside a `sync.Pool` when the Go Garbage Collector runs?
   * (A) They are locked and preserved.
   * (B) They are wiped clean and the pool becomes empty. *(Answer: B)*
   * (C) They are written to swap memory on disk.

2. **System Design Follow-up**: Why does the popular Go web framework `Fiber` use `sync.Pool` to reuse the `Context` object for every HTTP request?
   * *(To achieve zero memory allocations per request. Fiber pulls a Context from a pool, populates it with the HTTP data, passes it to your handler, and then returns it to the pool when your handler finishes. This is why you cannot start a Goroutine inside a Fiber handler and access the Context later—the Context will have been wiped and reused by another request!)*
