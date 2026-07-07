# Performance Engineering in Go

Welcome to the **Performance Engineering** curriculum module.
Once you know how to profile an application to find bottlenecks, you need to know how to fix them. This module covers advanced techniques for squeezing every last drop of performance out of the Go runtime.

## Curriculum

1. [Lesson 1: Mechanical Sympathy & Memory Allocations](001-memory-allocations.md)
   - The Stack vs The Heap
   - Escape Analysis (`go build -gcflags="-m"`)
   - Minimizing Garbage Collection (GC) pressure

2. [Lesson 2: Object Reuse with sync.Pool](002-sync-pool.md)
   - When to use `sync.Pool`
   - Pooling bytes buffers and JSON encoders

3. [Lesson 3: Concurrency Optimizations](003-concurrency-optimizations.md)
   - False sharing and cache lines
   - Lock contention vs Lock-free data structures (`sync/atomic`)

## The Golden Rule of Performance
> "Premature optimization is the root of all evil." — Donald Knuth

Always write clean, readable code first. Then measure it. Only optimize the code paths that your profiling tools prove are actual bottlenecks.
