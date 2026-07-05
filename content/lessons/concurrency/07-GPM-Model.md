# The G-P-M Model (Architecture Deep Dive)

To understand exactly how the Scheduler multiplexes 100,000 Goroutines onto 4 CPU cores, you must understand the architecture of the **G-P-M Model**.

The Go Runtime defines three core C-structs in its internal source code: `G`, `M`, and `P`.

## 1. The 'G' (Goroutine)
A `G` represents a Goroutine. It contains:
* Its current state (Runnable, Running, Waiting).
* Its stack memory (starting at 2KB).
* A pointer to the function it is supposed to execute.
* Hardware registers (saved when the Goroutine is paused).

## 2. The 'M' (Machine / OS Thread)
An `M` represents a physical OS Thread created by the Operating System.
* An `M` is dumb. It only knows how to execute native machine code.
* It does not know what a Goroutine is. It relies on a `P` to feed it work.

## 3. The 'P' (Processor / Context)
A `P` represents a Logical Processor. 
* By default, the Runtime creates exactly one `P` for every CPU core on your machine (e.g., 4 CPU cores = 4 `P`s).
* The `P` is the brain. It holds a **Local Run Queue (LRQ)** of up to 256 `G`s that are waiting to be executed.
* To execute Go code, an `M` (OS Thread) **must** acquire a `P`.

## 4. The Global Run Queue (GRQ)
If a `P`'s Local Run Queue fills up (exceeds 256 `G`s), it dumps the excess Goroutines into the **Global Run Queue (GRQ)**.

## 5. Work Stealing (The Genius of Go)

Imagine a 4-core machine.
* `P1`, `P2`, `P3`, and `P4` all have 10 Goroutines in their Local Run Queues.
* `P1` encounters a Goroutine that runs a heavy `for` loop. `P1` is busy for 10 milliseconds.
* Meanwhile, `P2` finishes all 10 of its Goroutines in 1 millisecond.

Now, `P2` has an empty queue, while `P1` is backed up. The CPU core for `P2` is sitting idle!

To prevent idle CPU cores, the Go Scheduler uses a **Work Stealing Algorithm**:
1. `P2` checks its Local Queue. Empty.
2. `P2` checks the Global Queue.
3. If the Global Queue is empty, `P2` will randomly pick another `P` (like `P1`) and literally **steal half of its Goroutines** from its Local Run Queue!

This guarantees that load is perfectly balanced across all CPU cores automatically, with zero configuration required by the developer.

## 6. The Netpoller (Handling I/O)

What happens if a Goroutine on `P1` makes a database query? Does the OS Thread (`M1`) go to sleep?
If `M1` goes to sleep, we lose 25% of our processing power on a 4-core machine!

To prevent this, Go uses the **Netpoller** (built on top of Linux `epoll` or macOS `kqueue`).
1. `G1` makes a network request.
2. The Runtime moves `G1` off the `P` and puts it into the asynchronous **Netpoller**.
3. `M1` is now free! It immediately grabs `G2` from the Local Run Queue and continues executing.
4. When the database replies, the Netpoller moves `G1` back into a Local Run Queue to finish formatting the JSON.
