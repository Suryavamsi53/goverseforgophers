# Flame Graphs and Tracing

The `top10` command in `pprof` is great, but looking at raw text output makes it difficult to understand the complex call chain of a deeply nested microservice.

To truly understand performance, you must visualize it using **Flame Graphs** and the **Go Tracer**.

## 1. The Web UI and Flame Graphs

Instead of using the interactive terminal shell, you can launch a beautiful web UI directly from the `go tool pprof` command.

```bash
# Add the -http flag to instantly boot a local web server!
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=10
```

This opens a browser window. If you click **View -> Flame Graph**, you will see a massive, colorful interactive chart.

* **The X-Axis (Width)**: Represents the amount of CPU time (or Memory) consumed. The wider the box, the more expensive it is!
* **The Y-Axis (Depth)**: Represents the Call Stack. If `FuncA` calls `FuncB`, `FuncB` will be stacked directly on top of `FuncA`.

If you see a massive, wide box at the top of the Flame Graph (e.g., `regexp.Compile`), you have instantly found your bottleneck! You shouldn't be compiling Regular Expressions inside a `for` loop!

## 2. Go Execution Tracer (`go tool trace`)

`pprof` operates by *Sampling*. It wakes up 100 times a second, checks what the CPU is doing, and goes back to sleep. This is great for finding heavy functions, but it is terrible for finding **Latency Spikes** (e.g., "Why did this specific HTTP request take 500ms?").

To diagnose latency, concurrency bottlenecks, or Garbage Collection pauses, you use the **Go Execution Tracer**.

The Tracer does not sample. It records *every single event* that occurs in the Go runtime!

```bash
# Download a 5-second trace from a live application
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5

# Launch the Trace visualizer!
go tool trace trace.out
```

### What does the Tracer reveal?

The Web UI for the Tracer shows a timeline of all your CPU cores.
1. **Goroutine Blocking**: You can click on a specific Goroutine and see exactly *why* it was paused. (e.g., "Blocked on Channel Send", "Blocked on Mutex Lock", "Waiting for Network I/O").
2. **Garbage Collection**: You will see bright red blocks showing exactly when the Stop-The-World GC pauses occurred, and exactly how many microseconds they lasted.
3. **Core Utilization**: If you have an 8-core CPU, but the Tracer shows that 7 cores are completely empty and 1 core is doing 100% of the work, you know your Concurrency model is flawed (e.g., you are bottlenecking on a single global Mutex!).

*Warning: Because the Tracer records every single event, it imposes a massive performance penalty (~20% overhead). Do not leave tracing enabled continuously in production!*
