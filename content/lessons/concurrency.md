# Go Concurrency: Goroutines and Channels

Concurrency in Go is a first-class citizen. Unlike threads in languages like Java or C++, Go uses **goroutines**, which are lightweight, user-space threads managed by the Go runtime.

## Goroutines

To start a new goroutine, you simply use the `go` keyword followed by a function call:

```go
package main

import (
	"fmt"
	"time"
)

func sayHello() {
	fmt.Println("Hello from goroutine!")
}

func main() {
	go sayHello() // Starts a new goroutine
	
	// We need to wait a bit, otherwise main might exit before the goroutine runs
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Hello from main!")
}
```

## Channels

Goroutines communicate via **channels**. Channels provide a way for two goroutines to synchronize execution and communicate by passing a value of a specified element type.

```go
package main

import "fmt"

func main() {
	messages := make(chan string) // Create a new channel of strings

	go func() {
		messages <- "ping" // Send a value into the channel
	}()

	msg := <-messages // Receive a value from the channel
	fmt.Println(msg)
}
```

### Best Practices
- **Don't communicate by sharing memory; share memory by communicating.**
- Always close channels if the receiver needs to know that no more values will be sent.
- Be careful with unbuffered channels, they can lead to deadlocks if not handled properly.
