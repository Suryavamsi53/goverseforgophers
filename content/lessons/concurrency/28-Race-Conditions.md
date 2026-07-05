# Race Conditions

A **Race Condition** is the most insidious bug in computer science. It occurs when two Goroutines access the same variable concurrently, and at least one of the accesses is a Write.

Because the Go Scheduler can pause a Goroutine at *any* nanosecond, the order of execution is non-deterministic. If your application works perfectly on your laptop 99 times, it might crash on the 100th time in Production.

## 1. The Classic Counter Bug

```go
var counter int

func main() {
    var wg sync.WaitGroup
    
    // Spawn 1000 Goroutines
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            // DANGER: Data Race!
            counter++ 
        }()
    }
    
    wg.Wait()
    fmt.Println(counter) 
}
```

If you run this code, it might print `943`. The next time it might print `988`. It will almost never print `1000`. Why?

Because `counter++` is actually three separate assembly instructions:
1. `READ` the value from RAM into the CPU Register.
2. `ADD` 1 to the CPU Register.
3. `WRITE` the CPU Register back to RAM.

If Goroutine A reads `0` and pauses. Then Goroutine B reads `0`, adds 1, and writes `1`. Then Goroutine A resumes, adds 1, and writes `1`. 
Two increments occurred, but the value is only `1`! Data has been lost forever.

## 2. The Solution

As we learned in the `sync` module, you fix Race Conditions by establishing a "Happens-Before" guarantee using one of three tools:
1. `sync.Mutex` (Locking)
2. `sync/atomic` (Hardware locks)
3. **Channels** (Passing data instead of sharing it)

## 3. The Race Detector (The Ultimate Weapon)

Because Race Conditions are invisible to the naked eye, Google built an incredible tool directly into the Go compiler: **The Race Detector**.

When compiling or running your code, simply add the `-race` flag.

```bash
go run -race main.go
go test -race ./...
```

The compiler will inject special instrumentation code into your binary that tracks every single memory access during runtime. If it detects two Goroutines touching the same memory address without a Mutex/Channel in between, it will throw a massive red warning directly into your terminal:

```text
==================
WARNING: DATA RACE
Write at 0x00c0000a6010 by goroutine 7:
  main.main.func1()
      /app/main.go:12 +0x4c

Previous read at 0x00c0000a6010 by goroutine 8:
  main.main.func1()
      /app/main.go:12 +0x38
==================
```

### The Production Rule
The Race Detector makes your code run ~10x slower and use ~10x more memory. 
**Never run the `-race` flag in Production environments.** 
It should be used exclusively during local development and in your CI/CD Pipeline (GitHub Actions) to block bad code from being merged.
