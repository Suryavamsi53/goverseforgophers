# Garbage Collection Mechanics

Languages like C and C++ require manual memory management (`malloc` and `free`). If you forget to `free` memory, your server crashes with an Out-Of-Memory (OOM) error.

Go handles this automatically using a **Concurrent Garbage Collector (GC)**.

## 1. The Tricolor Mark-and-Sweep Algorithm

The Go GC operates in two main phases: Mark and Sweep.
It uses a "Tricolor" algorithm to find which variables on the Heap are still being used by your program.

1. **White**: Objects that the GC has not checked yet (or objects that are dead and ready to be deleted).
2. **Grey**: Objects the GC knows are alive, but it hasn't checked their children yet.
3. **Black**: Objects the GC knows are alive, AND it has checked all their children.

### The Marking Phase
When the GC starts, every object on the Heap is White.
The GC looks at the "Roots" (global variables and variables currently on the Stack). It colors them Grey.
Then, it pulls a Grey object, looks at all the pointers inside that object, colors the children Grey, and colors the parent Black.
It repeats this loop until there are no Grey objects left!

### The Sweeping Phase
Once Marking is done, any object that is still White is mathematically proven to be "unreachable" by your Go program. The Sweeping phase simply deletes all the White objects and frees the RAM!

## 2. Stop-The-World (STW) Pauses

In Java, the GC traditionally pauses your entire application to do the Marking. If you have a 10GB Heap, your web server might freeze for 2 entire seconds!

Go's GC is optimized for **Ultra-Low Latency**.
Go performs the Tricolor Marking *concurrently* while your application is still running! 
It only requires two incredibly tiny "Stop-The-World" pauses to synchronize states. In modern Go, these STW pauses are consistently under **1 millisecond**, even on a 50GB Heap!

## 3. The Pacer and GOGC

How often does the Garbage Collector run? 
Go uses a "Pacer" algorithm controlled by a single environment variable: `GOGC` (default is 100).

`GOGC=100` means the GC will trigger when the Heap size *doubles*.
* If your application uses 50MB of alive memory, the GC will sit idle until the Heap reaches 100MB. Then it runs and cleans it back down to 50MB.

**Tuning the Pacer:**
* `GOGC=50`: The GC runs more frequently (when memory grows by 50%). Uses less RAM, but wastes more CPU.
* `GOGC=200`: The GC runs less frequently (when memory grows by 200%). Uses a lot of RAM, but saves CPU.
* `GOGC=off`: Turns the GC off completely! (Dangerous, but used in CLI tools that run for 2 seconds and exit).

## 4. The Memory Limit (Go 1.19+)

Prior to Go 1.19, if your Kubernetes Pod had a memory limit of 500MB, and the Go application spiked to 501MB before the Pacer decided to run, the Linux OOM-Killer would ruthlessly assassinate your Pod!

Go 1.19 introduced the `GOMEMLIMIT` environment variable.

```bash
export GOMEMLIMIT=450MiB
```
This tells the Go Pacer: "I don't care what GOGC is set to. If my RAM approaches 450MB, you must trigger an emergency Garbage Collection immediately to save my life!"
This single feature eradicated 90% of Kubernetes OOM crashes in Go.
