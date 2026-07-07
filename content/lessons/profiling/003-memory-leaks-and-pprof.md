# Memory Leaks & Advanced pprof

## 1. Learning Objectives
* **What you'll learn**: How to track down memory leaks in Go applications using `pprof` heap profiles.
* **Why it matters**: Go is a garbage-collected language, which means you don't have to manually `free()` memory. However, memory leaks are still very possible in Go! If a server leaks memory, it will eventually hit OOM (Out Of Memory) and crash in production.
* **Where it's used**: Long-running background workers, custom caching implementations, and WebSocket managers.

---

## 2. How does Go leak memory?
If Go has a Garbage Collector, how can memory leak? 
A memory leak in Go happens when you maintain a **reference** to an object that you no longer need. Because the reference exists, the Garbage Collector says, *"Ah, they are still using this, I cannot delete it!"*

### Common Leak 1: The Global Map
If you build an in-memory cache using a global `map[string]string` but forget to `delete(map, key)` when the data expires, the map will grow infinitely until the server crashes.

### Common Leak 2: The Abandoned Goroutine
Goroutines take up memory (at least 2KB for their stack). If you start a Goroutine that blocks forever waiting on a channel that will never receive data, that Goroutine is "leaked." It will never die, and its memory is locked forever.

---

## 3. Detecting Leaks with pprof

When your server's RAM usage is constantly climbing, you can use `pprof` to find exactly which line of code is allocating the memory that isn't being freed.

### Step 1: Get the Heap Profile
Run this terminal command while your server is running (assuming `net/http/pprof` is imported and running on port 8080):

```bash
go tool pprof -web http://localhost:8080/debug/pprof/heap
```

### Step 2: In-use Space vs Allocated Space
By default, `pprof` shows you `inuse_space`. This is exactly what you want for memory leaks. It shows you memory that has been allocated and is **currently sitting in RAM**.

If you want to see where memory is being allocated *frequently* (even if it's being cleaned up by the GC properly, which causes CPU overhead), you can switch to `alloc_space`:

```bash
go tool pprof -alloc_space -web http://localhost:8080/debug/pprof/heap
```

---

## 4. Reading the Graph
When the `-web` interface opens your browser, you will see a massive graph of boxes.
* **Large Boxes**: Indicate functions holding the most memory.
* **Arrows**: Show the call stack (Function A called Function B).

If you see a box for your `func startWorkerQueue()` holding 4GB of memory, you immediately know where to look.

### The CLI Alternative: `top`
If you prefer the terminal, you can just run `go tool pprof` without the `-web` flag, and then type `top` to see the worst offenders:

```bash
(pprof) top
Showing nodes accounting for 512MB, 95% of total
      flat  flat%   sum%        cum   cum%
     400MB 75.00% 75.00%      400MB 75.00%  main.badCacheStore
      50MB 10.00% 85.00%      450MB 85.00%  main.processData
```

---

## 5. Quiz

1. **MCQ**: You have a Goroutine that executes `time.Sleep(100 * time.Hour)`. Is its memory eligible for Garbage Collection?
   * (A) Yes, because it's sleeping.
   * (B) No, because the Goroutine is still alive and has an active stack. *(Answer: B)*
   * (C) Only the variables inside it are collected.

2. **System Design Follow-up**: If you want to take a snapshot of the heap memory and compare it to another snapshot 10 minutes later to definitively prove a slow memory leak, how do you do it?
   * *(You use `pprof` to download the profile to a file (`wget http://localhost:8080/debug/pprof/heap -O heap1.out`), wait 10 minutes, get `heap2.out`, and then run `go tool pprof -base heap1.out heap2.out`. This shows you the **diff**, revealing exactly what memory grew between the two snapshots!)*
