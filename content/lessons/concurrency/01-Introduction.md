# Introduction to Concurrency

Before we write a single line of concurrent Go code, we must fundamentally understand why concurrency exists and the architectural problems it solves. 

## 1. The Bottleneck of Sequential Execution

In a traditional, single-threaded application (like Node.js or a simple Python script), execution is strictly **sequential**. 
If your web server receives a request that requires:
1. Fetching a user from the Database (takes 50ms)
2. Fetching their recent orders from an API (takes 100ms)
3. Calculating their discount (takes 5ms)

A sequential program will execute these one after the other. Total time: **155ms**.

During the 150ms spent waiting on the network (I/O), the CPU is literally doing nothing. It is sitting idle, wasting precious compute cycles.

## 2. Enter Concurrency

Concurrency is the ability of an application to deal with multiple things at once. 

Instead of waiting for the database to return before calling the external API, what if we fired **both** network requests at the exact same time?

1. Fire DB request & API request simultaneously.
2. The DB returns in 50ms.
3. The API returns in 100ms.
4. CPU calculates the discount in 5ms.

Total time: **105ms**. 

We just made the endpoint 33% faster without writing a single line of optimization logic. We simply stopped wasting idle CPU time.

## 3. Go's Concurrency Philosophy

Most languages handle concurrency using OS Threads (Java, C++) or Async/Await Event Loops (JavaScript, Python, Rust).

Go took a radically different approach. It introduced the **Goroutine**.

> "Concurrency is not Parallelism." — Rob Pike

In Go, you don't worry about thread pools, event loops, or callbacks. You simply write synchronous-looking code, prefix it with the `go` keyword, and let the Go Runtime's internal **Scheduler** figure out how to multiplex it across your CPU cores.

In this module, we will dive deep into the internals of the Go Scheduler, Goroutines, Channels, and the synchronization primitives that make Go the king of modern backend engineering.
