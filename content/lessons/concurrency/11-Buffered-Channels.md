# Buffered Channels

In the previous lesson, we learned that standard channels are **Unbuffered** (they have a capacity of 0). If you send data into an unbuffered channel, the Goroutine will permanently block until another Goroutine is ready to receive it.

But what if you want a worker to quickly dump 5 items into a queue and immediately move on to other work, without waiting for the consumer to catch up? 

You need a **Buffered Channel**.

## 1. Syntax and Capacity

You create a buffered channel by providing a second argument to the `make` function: the **capacity**.

```go
// Creates a channel that can hold up to 3 strings without blocking the sender
ch := make(chan string, 3)

// Sender does NOT block. It drops the data into the buffer instantly.
ch <- "Job 1"
ch <- "Job 2"
ch <- "Job 3"

// Sender BLOCKS here! The buffer is full (capacity 3 reached).
// It will sleep until a receiver pulls at least one item out.
ch <- "Job 4" 
```

## 2. The Internal Architecture (Ring Buffer)

Under the hood, a buffered channel is implemented as a **Circular Queue (Ring Buffer)** using an array.

When you create `make(chan int, 3)`, the Go Runtime allocates an array of size 3 in the Heap, along with two pointers:
* `sendx`: The index where the next item will be inserted.
* `recvx`: The index where the next item will be extracted.
* `qcount`: The current number of items in the buffer.

Because it uses a pre-allocated array, sending to a buffered channel that is not full does not require a context switch. It is incredibly fast.

## 3. Asynchronous vs Synchronous

* **Unbuffered (Synchronous)**: Guarantees delivery. When `ch <- x` finishes, the sender knows with 100% certainty that the receiver has received the data.
* **Buffered (Asynchronous)**: No delivery guarantee. When `ch <- x` finishes, the sender only knows the data is sitting in the queue. If the receiver crashes before reading it, the data is lost.

## 4. When to use Buffered Channels

Senior Go engineers strictly avoid Buffered Channels unless they mathematically need them. Unbuffered channels force you to design cleaner, highly-synchronized systems.

However, Buffered Channels are absolutely required for two specific Enterprise Patterns:
1. **Rate Limiting**: Creating a channel with a capacity of 100 to ensure only 100 HTTP requests are processed simultaneously.
2. **The Worker Pool Pattern**: A central orchestrator dumps 1,000 URLs into a buffered channel, and 10 worker Goroutines pull from the buffer at their own pace to download the files.
