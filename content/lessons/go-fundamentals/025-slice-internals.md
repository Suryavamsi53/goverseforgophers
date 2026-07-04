# Slice Internals (Deep Dive)

To truly master Go and prevent subtle memory bugs, you must understand exactly how Slices are implemented inside the Go compiler runtime.

## 1. The `SliceHeader` Struct

A slice is not actually an array. It is a tiny, 24-byte struct (on 64-bit architectures) called a `SliceHeader`. 

If you look at the Go source code (`reflect.SliceHeader`), a slice is defined as:

```go
type SliceHeader struct {
    Data uintptr // Pointer to the underlying array (8 bytes)
    Len  int     // Length of the slice (8 bytes)
    Cap  int     // Capacity of the underlying array (8 bytes)
}
```

### 📊 Architecture Visualization

```mermaid
graph LR
    subgraph SliceHeader [Slice Struct (24 Bytes)]
        P[Data Pointer]
        L[Len: 3]
        C[Cap: 5]
    end

    subgraph Memory [Heap Allocated Backing Array]
        A0[ 0 ]
        A1[ 1 ]
        A2[ 2 ]
        A3[ 3 ]
        A4[ 4 ]
    end

    P -->|Points to index 1| A1
    
    style SliceHeader fill:#f9f,stroke:#333,stroke-width:2px
    style Memory fill:#bbf,stroke:#333,stroke-width:2px
```

When you write `slice := arr[1:4]`:
1. `Data` points to the memory address of `arr[1]`.
2. `Len` is set to `3`.
3. `Cap` is set to `4` (assuming `arr` had a length of 5).

## 2. Why Passing Slices is Fast

When you pass a slice to a function in Go, it is passed by value (copied). But because the slice is just a 24-byte struct, passing it is virtually instantaneous.

```go
func process(data []byte) { ... }
```
Even if the backing array contains 10 Gigabytes of video data, calling `process(data)` only copies 24 bytes (the pointer, length, and capacity). The pointer inside the copied struct still points to the exact same 10GB array in memory!

## 3. The Re-Slicing Danger

Because slices are just windows into arrays, multiple slices can point to the exact same backing array.

```go
func main() {
    original := []int{1, 2, 3, 4, 5}
    
    // Create a new window looking at the same memory
    window := original[1:3] // [2, 3]
    
    // Modify the new window
    window[0] = 99
    
    // The original is mutated!
    fmt.Println(original) // [1, 99, 3, 4, 5]
}
```

**Memory Leak Trap**: If you load a massive 1GB file into a slice, and then create a tiny slice representing just the first 10 bytes (`tiny := massive[:10]`), the Garbage Collector **cannot** free the 1GB array because `tiny` still holds a pointer to it! 
