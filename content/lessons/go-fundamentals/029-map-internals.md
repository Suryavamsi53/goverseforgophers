# Map Internals (Deep Dive)

A map in Go is not a simple array; it is a highly optimized Hash Table implemented in the Go runtime. Understanding its internals helps you optimize for memory and speed.

## 1. The `hmap` Struct

When you create a map, Go allocates an `hmap` (hash map) struct. This struct holds metadata about the map, including a pointer to an array of **Buckets**.

```mermaid
graph TD
    subgraph hmap [hmap Struct]
        C[Count: Number of elements]
        B[B: log2 of number of buckets]
        P[Buckets Pointer]
    end

    subgraph Memory [Array of bmap (Buckets)]
        B0[Bucket 0]
        B1[Bucket 1]
        B2[Bucket 2]
    end

    P --> Memory
    
    subgraph BucketDesign [Inside a Single Bucket]
        K[8 Keys] --- V[8 Values] --- O[Overflow Pointer]
    end
    
    B0 -.-> BucketDesign
```

## 2. Buckets and the Rule of 8

In Go, each bucket (`bmap`) is designed to hold exactly **8 key-value pairs**. 

When you insert a key:
1. Go hashes the key.
2. The *low-order* bits of the hash determine which bucket the data goes into.
3. The *high-order* bits of the hash are stored inside the bucket to quickly distinguish keys during lookups without comparing the full string.

If a bucket gets full (more than 8 items hash to the same bucket), Go creates an **Overflow Bucket** and chains them together.

## 3. Evacuation (Resizing)

As your map grows, the buckets fill up. If the map's "Load Factor" (average items per bucket) gets too high (typically > 6.5), Go triggers a **growth phase**.

1. Go allocates a new array of buckets that is **double** the size of the old one.
2. It does not pause your program to copy everything at once (which would cause massive latency spikes).
3. Instead, it performs **Incremental Evacuation**. Every time you write or read from the map, Go secretly moves a small chunk of data from the old buckets to the new buckets until the migration is complete.

## 4. Why Map Iteration is Randomized

If you iterate over a map, the order is **never** guaranteed.

```go
m := map[int]string{1: "A", 2: "B", 3: "C"}
for k, v := range m {
    fmt.Println(k, v) // Order changes every time you run the program!
}
```

Why? Because the Go team intentionally randomized map iteration! 
In early versions of Go, the order was technically stable until a resize occurred. Developers began accidentally relying on that undocumented order, causing production bugs when the map resized. To force developers to write correct code, the Go runtime now explicitly selects a random starting bucket every time you use `range` on a map.
