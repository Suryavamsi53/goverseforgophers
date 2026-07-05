# Starvation

Starvation is the quiet, invisible cousin of a Deadlock. 

In a Deadlock, execution stops entirely. In **Starvation**, the system continues to run, but one or more Goroutines are systematically denied access to resources (CPU or Mutexes), causing their execution to be infinitely delayed.

## 1. Mutex Starvation

Imagine you have a `sync.Mutex`. 
* **Goroutine A** (a heavy background worker) locks the Mutex, does 10 milliseconds of work, unlocks it, and instantly tries to lock it again in a loop.
* **Goroutine B** (a fast HTTP handler) tries to acquire the Mutex.

Because Goroutine A is constantly re-locking the Mutex immediately after unlocking it, Goroutine B might wait in line for a very long time, unable to squeeze in.

### Go's Built-in Solution (Barging vs Starvation Mode)
Historically, the Go Scheduler allowed "Barging"—if a new Goroutine arrived just as the Mutex was unlocked, it could steal the lock from Goroutines that had been waiting in line. This made the system fast, but caused severe starvation.

To fix this, Go updated the `sync.Mutex` internal architecture to include **Starvation Mode**.
1. If a Goroutine waits for a Mutex for more than **1 millisecond**, the Mutex enters Starvation Mode.
2. In Starvation Mode, Barging is completely disabled.
3. The Mutex is handed off strictly in FIFO (First-In, First-Out) order to the sleeping Goroutines in the queue.
4. Once the queue is empty, or the wait time drops below 1ms, it returns to Normal Mode for speed.

## 2. CPU Starvation (The Preemption Problem)

Before Go 1.14, if you wrote a Goroutine containing an infinite math loop without any network calls, channel operations, or Mutex locks, the Go Scheduler could not pause it. 

If you had a 1-core machine, that infinite math loop would completely monopolize the CPU. The web server would freeze, and all other Goroutines would starve to death.

As discussed in the Scheduler lesson, Go 1.14 solved this with **Asynchronous Preemption**. The `sysmon` background thread now shoots UNIX `SIGURG` signals at greedy Goroutines every 10 milliseconds to forcibly pause them and cure CPU starvation.

## 3. RWMutex Writer Starvation

The `sync.RWMutex` presents a unique starvation risk. 

If you have a massive stream of Readers constantly acquiring `RLock()`, and a Writer calls `Lock()`, the Writer must wait for all Readers to finish. But what if new Readers keep arriving and acquiring `RLock()` before the old ones finish? The Writer would starve forever!

To prevent this, the Go standard library implements a strict rule: **When a Writer calls `Lock()`, no *new* Readers are allowed to acquire `RLock()`.** 
New Readers are blocked and queued up behind the Writer. The Writer waits for the existing Readers to finish, executes, and then unblocks the new Readers.
