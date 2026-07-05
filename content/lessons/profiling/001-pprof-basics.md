# pprof Basics (CPU and Memory Profiling)

When your Go application uses 100% CPU in production, you cannot guess what is causing it. You need mathematical proof.

Go has an incredibly powerful profiling tool built directly into the standard library called **pprof**.

## 1. Exposing pprof over HTTP

To profile a live Go application, you simply import the `net/http/pprof` package. 
It automatically registers a set of debugging endpoints on your default HTTP router!

```go
import (
    "net/http"
    _ "net/http/pprof" // The blank identifier executes the init() function!
    "log"
)

func main() {
    // Start the server (usually on a private internal port!)
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // ... start your real application logic
}
```

*Security Warning: Never expose port 6060 to the public internet! pprof dumps the entire state of your application's memory and CPU, which is a massive security vulnerability.*

## 2. CPU Profiling

If your application is lagging, you want to see exactly which functions are consuming the most CPU cycles.

From your terminal (on your laptop), you connect to the live server using the `go tool pprof` command.

```bash
# This samples the CPU for 30 seconds!
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

Once the profile finishes, you enter the interactive pprof shell.
Type `top10`:

```text
(pprof) top10
Showing nodes accounting for 9.20s, 95.83% of 9.60s total
      flat  flat%   sum%        cum   cum%
     4.10s 42.71% 42.71%      4.10s 42.71%  runtime.cgocall
     2.80s 29.17% 71.88%      2.80s 29.17%  encoding/json.Unmarshal
     1.50s 15.62% 87.50%      1.50s 15.62%  main.calculateHeavyMath
```

This output instantly proves that `encoding/json.Unmarshal` is consuming 29% of your entire CPU! You now know exactly where to optimize.

## 3. Memory Profiling (Heap)

If your application is crashing with Out-Of-Memory (OOM) errors, or the Garbage Collector is working overtime, you need to see exactly where RAM is being allocated.

```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

In the pprof shell, you can use the `list` command to see the exact line of code causing the allocations!

```text
(pprof) list main.generateMassiveData
Total: 500MB
ROUTINE ======================== main.generateMassiveData in main.go
      50MB      500MB (flat, cum) 100.00% of Total
         .          .     12: func generateMassiveData() {
         .          .     13:     for i := 0; i < 100000; i++ {
      50MB      500MB     14:         data := make([]byte, 5000) // THIS LINE IS THE PROBLEM!
         .          .     15:         process(data)
         .          .     16:     }
         .          .     17: }
```
pprof gives you X-ray vision into the physical hardware utilization of your Go binary.
