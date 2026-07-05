# Zero-Allocation Techniques

If your Go microservice processes 50,000 requests per second, and every request allocates a 10KB struct on the Heap, you are generating 500 Megabytes of garbage per second.

The Garbage Collector will have to run constantly, consuming 50% of your CPU just to clean up the trash.

To achieve extreme performance, Enterprise Go engineers use **Zero-Allocation** techniques to completely eliminate Heap allocations.

## 1. The `sync.Pool` (Object Recycling)

Instead of creating a new 10KB struct, letting the GC delete it, and creating a new one... what if you just kept a box of pre-built structs and reused them?

This is exactly what `sync.Pool` does. It is a thread-safe object pool.

```go
// 1. Define the Pool globally
var bufferPool = sync.Pool{
    // The New function tells the pool how to create an item if the box is empty!
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func HandleRequest() {
    // 2. Grab an object from the Pool (Type assert it!)
    // If the pool is empty, it calls the New function. 
    // If the pool has an old item, it returns it instantly!
    buf := bufferPool.Get().(*bytes.Buffer)
    
    // 3. CRITICAL: Clean the object before using it! It might have old data!
    buf.Reset()
    
    // 4. Use the object
    buf.WriteString("hello world")
    
    // 5. Put the object back in the box for the next Goroutine to use!
    // We generated ZERO GARBAGE!
    bufferPool.Put(buf)
}
```
Libraries like `encoding/json` and web frameworks like `Fiber` use `sync.Pool` aggressively to achieve zero-allocation performance.

## 2. Pre-Allocating Slices

If you append to a slice, and the slice runs out of capacity, Go has to:
1. Allocate a brand new array on the Heap (double the size).
2. Copy all the old elements over.
3. Garbage collect the old array.

This is massively expensive. If you know the final size, you must pre-allocate the capacity!

```go
// BAD: Allocates multiple times as it grows from 0 to 1000!
var badSlice []int
for i := 0; i < 1000; i++ {
    badSlice = append(badSlice, i)
}

// GOOD: Allocates exactly ONCE. Zero reallocation overhead!
goodSlice := make([]int, 0, 1000) // Length 0, Capacity 1000
for i := 0; i < 1000; i++ {
    goodSlice = append(goodSlice, i)
}
```

## 3. String to Byte Conversion (Unsafe)

Converting a `string` to a `[]byte` forces a Heap allocation, because strings are immutable and bytes are mutable. Go is forced to make a copy to prevent you from mutating the original string.

```go
s := "hello"
b := []byte(s) // Allocates memory!
```

If you only need to *read* the bytes, and you promise not to mutate them, you can use the `unsafe` package to bypass the allocation entirely. (This is how high-performance JSON parsers work).

```go
import "unsafe"

// Zero allocations! It physically manipulates the Slice Header pointer!
b := unsafe.Slice(unsafe.StringData(s), len(s))
```

## 4. The Arena Allocator (Experimental)

Google is currently testing a revolutionary feature in Go called **Arenas** (available via `GOEXPERIMENT=arenas`).

An Arena allows you to allocate thousands of objects in a specific chunk of memory. When you are done, you call `arena.Free()`. 
Instead of the Garbage Collector scanning the objects individually, the entire chunk of memory is instantly deleted in $O(1)$ time! This allows Go to rival C++ and Rust in extreme high-frequency trading scenarios.
