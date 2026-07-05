# The Go Scheduler

The Go Scheduler is the crown jewel of the Go language. It is an $O(1)$ algorithmic engine responsible for distributing hundreds of thousands of Goroutines across the limited number of physical OS Threads available to the application.

## 1. M:N Scheduling

Most languages use a **1:1 Scheduler**. 
* In Java (prior to virtual threads) or C++, if you spawn 100 threads, you are asking the Operating System for 100 OS Threads. 
* The OS Scheduler must manage them.

Go uses an **M:N Scheduler**.
* `M` represents the number of Goroutines (e.g., 100,000).
* `N` represents the number of OS Threads (usually equal to the number of physical CPU cores, e.g., 4).
* The Go Scheduler is responsible for multiplexing the `M` Goroutines onto the `N` OS Threads.

## 2. Cooperative vs Preemptive Scheduling

How does the Scheduler know when to pause a running Goroutine and let another one run?

### Pre Go 1.14 (Cooperative)
Before Go 1.14, the scheduler was primarily **Cooperative**. It relied on the Goroutine to yield control back to the scheduler. A Goroutine would yield during:
* Network I/O (HTTP requests)
* Channel operations (`ch <- val`)
* System calls (Reading files)
* Mutex locks
* Garbage Collection pauses

**The Bug**: If a developer wrote an infinite math loop (`for { x++ }`), the Goroutine would *never* make a network call. It would hog the CPU core forever, starving all other Goroutines and deadlocking the system.

### Go 1.14+ (Asynchronously Preemptive)
The Go team introduced **Asynchronous Preemption**. 
The Go Runtime runs a background thread called `sysmon` (System Monitor). `sysmon` watches all running Goroutines. If it sees a Goroutine hogging a CPU core for more than **10 milliseconds**, it sends a UNIX `SIGURG` signal to that thread, forcibly pausing the Goroutine and injecting a call to `runtime.gosched()`.

This guarantees that no infinite loop can ever crash a Go web server.

## 3. The `GOMAXPROCS` Variable

By default, Go creates one OS Thread for every physical core on the host machine. If you deploy your Go app to an 8-core Kubernetes pod, it uses 8 OS Threads.

You can manually control this using the `runtime.GOMAXPROCS(n)` function or the `GOMAXPROCS` environment variable. 

*Setting `GOMAXPROCS=1` disables parallelism entirely, forcing all Goroutines to run concurrently on a single core.*
