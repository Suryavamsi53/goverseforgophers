# Mechanical Sympathy & Memory Allocations

## 1. Learning Objectives
* **What you'll learn**: The difference between Stack and Heap allocations in Go, what Escape Analysis is, and how to minimize Garbage Collection (GC) pauses.
* **Why it matters**: In high-performance systems (like trading engines or massive web scrapers), CPU time spent pausing the application to clean up memory (GC) directly translates to increased latency and decreased throughput.
* **Where it's used**: Any hot-path execution loop where microsecond latencies matter.

---

## 2. Real-world Story
Imagine a fast-food restaurant. 
If a customer eats in the store on a plastic tray (The Stack), the moment they leave, the tray is instantly wiped and ready for the next customer. Zero overhead.
If a customer takes the food to go in a cardboard box (The Heap), someone eventually has to drive around the city, find that cardboard box in the trash, and take it to the dump (Garbage Collection). This takes immense time and resources.

In Go, variables allocated on the **Stack** are instantly cleaned up the nanosecond a function returns. Variables allocated on the **Heap** survive after the function returns, meaning the Go Runtime has to periodically scan memory and clean them up.

---

## 3. Escape Analysis in Action

How does Go decide whether a variable goes on the fast Stack or the slow Heap? It uses a compiler phase called **Escape Analysis**. 

Look at this code:

```go
package main

type User struct {
    ID   int
    Name string
}

// Example 1: Stays on the Stack
func CreateUserValue() User {
    u := User{ID: 1, Name: "Alice"}
    return u // We return a COPY of the struct
}

// Example 2: Escapes to the Heap
func CreateUserPointer() *User {
    u := User{ID: 2, Name: "Bob"}
    return &u // We return a POINTER to the struct
}
```

In `Example 1`, `u` is copied and returned. The original `u` dies when the function ends. It is allocated on the Stack. Extremely fast.

In `Example 2`, we return a **pointer** to `u`. This means `u` must survive *after* the function returns so the caller can read it. Go says, "Ah, this variable has *escaped* the function!" and moves it to the Heap.

### Measuring it
You can ask the Go compiler what it's doing by passing `-gcflags="-m"` during build:

```bash
$ go build -gcflags="-m" main.go
./main.go:15:2: moved to heap: u
```

---

## 4. When to use Pointers vs Values
Many developers think, *"I should use pointers to avoid copying large structs!"* 
This is a **myth** in Go. 
Copying a 64-byte struct on the CPU cache is incredibly fast (nanoseconds). However, allocating a pointer to the Heap requires a lock on the memory allocator and eventually triggers a Garbage Collection pause (microseconds to milliseconds).

**Rule of Thumb for High Performance Go:**
1. Default to passing by Value (returns copies).
2. Only use Pointers if:
   - You **need** to mutate the original struct.
   - The struct is massively large (e.g., thousands of bytes).
   - The struct contains a `sync.Mutex` (Mutexes must never be copied).

---

## 5. Quiz

1. **MCQ**: What flag do you pass to the Go compiler to see which variables are being allocated on the heap?
   * (A) `go build -race`
   * (B) `go build -gcflags="-m"` *(Answer: B)*
   * (C) `go build -tags=heap`

2. **System Design Follow-up**: Why does generating a UUID inside a highly-concurrent HTTP handler often cause massive CPU spikes in Go?
   * *(Because many UUID libraries allocate a `[]byte` slice for the UUID on the heap for every single request, forcing the Garbage Collector to work constantly under high load.)*
