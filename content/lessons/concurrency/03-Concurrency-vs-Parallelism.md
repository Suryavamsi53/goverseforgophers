# Concurrency vs Parallelism

One of the most famous talks in Go's history was given by Rob Pike (co-creator of Go), titled: **"Concurrency is not Parallelism"**.

Many developers use these terms interchangeably, but they are fundamentally different concepts in computer science.

## 1. Concurrency (Dealing with many things at once)

Concurrency is a **Design Pattern**. It is how you structure your code so that multiple tasks can make progress over overlapping time periods.

Imagine a single barista working at a coffee shop:
1. They take Customer A's order.
2. They start grinding the espresso beans for Customer A.
3. While the machine is brewing the espresso (I/O bound waiting), the barista turns to Customer B and takes their order.
4. The barista then turns back, finishes Customer A's drink, and hands it to them.

There is only **one** barista (one CPU core). They are not doing two things at the exact same physical millisecond. But by quickly switching back and forth, they are *dealing* with multiple customers concurrently.

## 2. Parallelism (Doing many things at once)

Parallelism is **Hardware Execution**. It is the physical act of doing multiple things at the exact same nanosecond.

Imagine two baristas working at a coffee shop:
1. Barista 1 takes Customer A's order and makes it.
2. Barista 2 takes Customer B's order and makes it.

Because there are two physical workers (two CPU cores), the work is executing truly in parallel.

## 3. The Go Paradigm

In Go, you design your application to be **Concurrent** by splitting it into independent Goroutines. 

Whether those Goroutines actually run in **Parallel** depends entirely on the hardware the application runs on!
* If you run your Go web server on a cheap 1-Core AWS EC2 instance, your Goroutines run concurrently (rapidly context-switching), but never in parallel.
* If you run that exact same Go binary on a 16-Core server, the Go Scheduler will automatically distribute those Goroutines across the 16 cores, achieving true Parallelism.

**You write concurrent code. The Go Runtime gives you parallelism for free if the hardware supports it.**
