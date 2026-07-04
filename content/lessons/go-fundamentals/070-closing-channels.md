# Closing Channels

Unlike opening a file or a database connection, **you do not have to close every channel you create**. 

Channels do not tie up system resources like file descriptors. If a channel is no longer referenced, the Garbage Collector will destroy it automatically, even if you never close it.

You only close a channel when you explicitly need to tell the receiver: *"I am done sending data, you can stop waiting."*

## 1. The `close` Function

You close a channel using the built-in `close()` function.

```go
func main() {
    ch := make(chan int, 3)
    
    ch <- 1
    ch <- 2
    
    // Signal that no more data will be sent
    close(ch) 
}
```

### ⚠️ The Golden Rule of Closing
**Only the sender should ever close a channel, never the receiver.** 
Sending data to a closed channel causes a fatal `panic`. If a receiver closes the channel while the sender is still working, the sender will crash the entire program the next time it tries to send!

## 2. Reading from a Closed Channel

What happens if you try to receive data from a closed channel?

1. If there is still data inside the buffer, you will receive that data normally.
2. If the channel is empty and closed, you will immediately receive the **zero value** of the channel's type, without blocking.

```go
ch := make(chan int, 1)
ch <- 42
close(ch)

fmt.Println(<-ch) // Prints 42
fmt.Println(<-ch) // Prints 0 (Instantly, does not block)
fmt.Println(<-ch) // Prints 0
```

## 3. The Comma Ok Idiom

If a channel is closed and returning `0`, how do you know if the sender actually sent a `0`, or if the channel is just closed?

We use the "comma ok" idiom (just like Maps and Type Assertions!).

```go
val, ok := <-ch
if !ok {
    fmt.Println("Channel is permanently closed!")
} else {
    fmt.Println("Received:", val)
}
```

## 4. Ranging over Channels

The most idiomatic way to consume all data from a channel until it closes is using a `for range` loop. The loop will automatically block and wait for new data, and will immediately exit the moment the channel is closed.

```go
func main() {
    ch := make(chan string, 2)
    
    ch <- "Task A"
    ch <- "Task B"
    close(ch)
    
    // This loop safely empties the channel and breaks automatically
    for task := range ch {
        fmt.Println("Processing:", task)
    }
}
```
