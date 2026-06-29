# Go Concurrency Masterclass

Welcome to the ultimate guide to Concurrency in Go. Concurrency is arguably Go's most powerful and distinguishing feature. It was designed from the ground up to take advantage of modern multi-core processors, making it incredibly simple to write highly performant, scalable, and concurrent code.

## 🎯 Learning Outcome

By the end of this course, you will transition from a beginner to a highly capable backend engineer who can:
- Confidently design, implement, debug, and optimize concurrent Go applications for production systems.
- Deeply understand the internal Go scheduler and GPM model.
- Prevent race conditions, deadlocks, and Goroutine leaks.
- Implement advanced enterprise patterns (Worker Pools, Rate Limiters, Fan-In/Fan-Out).
- Ace Google-style concurrency interview questions.

---

## 📚 Course Outline

### Phase 1: Foundations & The Runtime
* `01-Introduction`: The `go` keyword and Goroutines.
* `02-Why-Concurrency`: Why do we need it? CPU vs I/O bounds.
* `03-Concurrency-vs-Parallelism`: Understanding the difference.
* `04-Process-vs-Thread-vs-Goroutine`: OS level vs User-space threads.
* `05-Go-Runtime`: How Go manages execution.
* `06-Go-Scheduler`: The brains behind Goroutines.
* `07-GPM-Model`: Goroutines, Processors, and Machine threads.
* `08-Goroutines`: In-depth look at Goroutine lifecycle and stack.

### Phase 2: Core Primitives
* `09-WaitGroup`: Waiting for tasks to finish.
* `10-Channels`: Safely communicating between Goroutines.
* `11-Buffered-Channels`: Asynchronous message passing.
* `12-Unbuffered-Channels`: Synchronous message passing.
* `13-Channel-Directions`: Send-only and Receive-only channels.
* `14-Channel-Closing`: Signaling completion.
* `15-Range`: Iterating over channels.
* `16-Select`: Handling multiple channel operations.
* `17-Timers` & `18-Tickers`: Scheduling and intervals.
* `19-Context` & `20-Cancellation`: Managing request lifecycles.

### Phase 3: Memory Synchronization & Pitfalls
* `21-Mutex` & `22-RWMutex`: Locking shared memory.
* `23-Atomic`: Lock-free synchronization.
* `24-sync.Once`, `25-sync.Map`, `26-sync.Cond`: Advanced sync primitives.
* `27-errgroup`: Synchronizing error handling.
* `28-Race-Conditions`, `29-Deadlocks`, `30-Starvation`, `31-Livelock`: Debugging concurrency bugs.

### Phase 4: Advanced Patterns
* `32-Worker-Pool`: Limiting resource usage.
* `33-Fan-In` & `34-Fan-Out`: Multiplexing data streams.
* `35-Pipeline`: Streaming data processing.
* `36-Semaphore` & `37-Rate-Limiter`: Throttling systems.

### Phase 5: Production Engineering & Capstone
* `38-Performance`, `39-Benchmarking`, `40-Profiling`: Tuning your code.
* `41-Capstone`: The final distributed systems project.

Let's get started with **01-Introduction**!
