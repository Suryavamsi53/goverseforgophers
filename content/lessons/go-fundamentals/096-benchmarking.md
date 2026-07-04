# Benchmarking

How do you scientifically prove that `strconv.Itoa` is faster than `fmt.Sprintf`? 
You benchmark it.

The `testing` package provides a powerful benchmarking tool that will execute your function millions of times to calculate its exact nanosecond performance and memory allocations.

## 1. Writing a Benchmark

Benchmarks live in `_test.go` files alongside your unit tests.
* The function name must start with `Benchmark`.
* It must accept `(b *testing.B)`.
* It must contain a `for` loop that runs exactly `b.N` times.

```go
// performance_test.go
package main

import (
    "fmt"
    "strconv"
    "testing"
)

func BenchmarkItoa(b *testing.B) {
    for i := 0; i < b.N; i++ {
        strconv.Itoa(42)
    }
}

func BenchmarkSprintf(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("%d", 42)
    }
}
```

## 2. Running the Benchmark

By default, `go test` only runs Unit Tests. To run benchmarks, use the `-bench` flag with a regex matching the functions you want to run (or `.` for all).

```bash
$ go test -bench=.
```

**The Output:**
```text
BenchmarkItoa-10        100000000       4.15 ns/op
BenchmarkSprintf-10      15000000      78.4  ns/op
```
* `BenchmarkItoa-10`: Ran using 10 CPU cores.
* `100000000`: The Go testing engine dynamically increased `b.N` to 100 million until the test ran long enough to get a statistically valid average.
* `4.15 ns/op`: It took 4.15 nanoseconds per operation. 

**Conclusion**: `Itoa` is mathematically 18x faster than `Sprintf`!

## 3. Auditing Memory Allocations

To see *why* `Sprintf` is slower, we can ask the benchmarking tool to track Heap escapes and memory allocations by adding the `-benchmem` flag.

```bash
$ go test -bench=. -benchmem
```

**The Output:**
```text
BenchmarkItoa-10        100000000    4.15 ns/op    0 B/op    0 allocs/op
BenchmarkSprintf-10      15000000    78.4 ns/op   16 B/op    2 allocs/op
```
Here is the architectural proof! `Itoa` uses exactly 0 bytes of heap memory and triggers 0 allocations (it stays entirely on the fast Stack). `Sprintf` requires 16 bytes of Heap memory and 2 distinct allocations per call, which triggers Garbage Collection lag!

## 4. The Compiler Deletion Trap

The Go compiler is incredibly aggressive at optimization. If you write a benchmark that calculates a value but never uses it, the compiler might realize the code is "dead" and silently delete the entire loop during compilation! 

This will result in a fake benchmark of `0.00 ns/op`.

To prevent this, you must assign the result to a global variable (creating a side-effect).

```go
var result string // Global

func BenchmarkItoa(b *testing.B) {
    var r string
    for i := 0; i < b.N; i++ {
        r = strconv.Itoa(42) // Local assignment
    }
    result = r // Escape to global to prevent compiler deletion!
}
```
