# Profiling (pprof)

Benchmarking is great for testing tiny, isolated functions. But what if your massive production web server starts randomly consuming 100% CPU, and you have no idea which function is causing it?

You use **pprof** (Performance Profiler).

## 1. Enabling `pprof` in Production

Go's profiler has remarkably low overhead. The Go team officially recommends leaving it enabled in production environments so you can diagnose live servers without restarting them.

To enable it, simply import `net/http/pprof` into your web server. It will automatically inject debugging endpoints into the `http.DefaultServeMux`.

```go
package main

import (
    "net/http"
    // The underscore prevents the compiler from removing the unused import.
    // It triggers the package's init() function, which mounts the profiler routes!
    _ "net/http/pprof" 
)

func main() {
    // Start your application
    go runHeavyWorkload()

    // Start a dedicated port specifically for the profiler
    http.ListenAndServe("localhost:6060", nil)
}
```

## 2. Capturing a Profile

While your server is experiencing heavy load, open your terminal and run the `go tool pprof` command. 

This command reaches out to your server, commands it to capture a 30-second snapshot of exactly what the CPU is doing, downloads the data, and opens an interactive terminal interface.

```bash
$ go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

Inside the interactive prompt, type `top` to see the top 10 functions burning CPU cycles:

```text
(pprof) top
Showing nodes accounting for 9.85s, 95% of 10.30s total
      flat  flat%   sum%        cum   cum%
     4.10s 39.81% 39.81%      4.10s 39.81%  syscall.Syscall
     2.30s 22.33% 62.14%      2.30s 22.33%  runtime.mcall
     1.10s 10.68% 72.82%      3.80s 36.89%  main.HeavyJSONParser
```
*In this snapshot, we instantly see that `main.HeavyJSONParser` is consuming a massive 36% of the CPU!*

## 3. Flame Graphs

Reading text output is helpful, but visualizing it is better. If you add the `-http` flag, `pprof` will generate a stunning, interactive web dashboard with a **Flame Graph**.

```bash
$ go tool pprof -http=":8000" http://localhost:6060/debug/pprof/profile?seconds=30
```

A Flame Graph shows your entire call stack visually. The wider a box is, the more CPU time it consumed. This allows you to instantly spot performance bottlenecks and the exact line of code that triggered them.

## 4. Memory Profiling

If your server isn't burning CPU, but it's slowly running out of RAM (a Memory Leak), you can profile the Heap instead of the CPU.

```bash
$ go tool pprof http://localhost:6060/debug/pprof/heap
```
This tells you exactly which lines of code allocated memory that the Garbage Collector was unable to clean up!
