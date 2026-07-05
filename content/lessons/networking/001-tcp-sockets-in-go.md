# TCP Sockets in Go

While `net/http` is fantastic for building web servers, there are times when you need to drop down to the raw networking layer and build a custom protocol (e.g., a multiplayer game server, an IoT device coordinator, or a custom database like Redis).

To do this, you must interact directly with the **TCP (Transmission Control Protocol)** layer using the `net` package.

## 1. Creating a TCP Server

A TCP server does not understand URLs, Headers, or JSON. It only understands a continuous stream of raw bytes.

```go
package main

import (
    "fmt"
    "net"
)

func main() {
    // 1. Open a raw TCP socket on port 9000
    listener, err := net.Listen("tcp", ":9000")
    if err != nil { panic(err) }
    defer listener.Close()
    
    fmt.Println("Listening on TCP port 9000...")

    // 2. The infinite accept loop!
    for {
        // Blocks until a client establishes a 3-Way Handshake!
        conn, err := listener.Accept()
        if err != nil { continue }
        
        // 3. Spawn a Goroutine to handle this specific connection!
        go handleConnection(conn)
    }
}
```

## 2. Handling the Connection Stream

Unlike HTTP (which has a distinct Request and Response), a TCP `conn` is a permanent, bi-directional pipe. You can read and write to it simultaneously until someone closes it.

```go
func handleConnection(conn net.Conn) {
    defer conn.Close() // ALWAYS close the connection when done!
    
    // Create a 1KB buffer to store the incoming bytes
    buffer := make([]byte, 1024)
    
    for {
        // Read blocks until the client sends data
        n, err := conn.Read(buffer)
        if err != nil {
            fmt.Println("Client disconnected!")
            return
        }
        
        // Convert the raw bytes to a string
        msg := string(buffer[:n])
        fmt.Printf("Received: %s", msg)
        
        // Write data back down the pipe!
        conn.Write([]byte("Message received!\n"))
    }
}
```

## 3. The Framing Problem

If you build a custom protocol, you will encounter the **TCP Framing Problem**.

TCP is a *Streaming* protocol, not a *Message* protocol. 
If the client sends "Hello" and then instantly sends "World", TCP does NOT guarantee that the server's `conn.Read()` will receive them separately!
The server might receive "Hello", or it might receive "HelloWorld" in a single read, or it might receive "Hel" and then "loWorld" in two reads! TCP only guarantees the *order* of the bytes, not the grouping.

**The Solution:**
You must implement a Framing Protocol. 
1. **Delimiters**: The easiest way is to mandate that every message ends with a newline (`\n`). The server uses `bufio.NewReader(conn).ReadString('\n')` to guarantee it only processes one full message at a time.
2. **Length-Prefixed**: Enterprise protocols (like gRPC) send the exact length of the payload in the first 4 bytes. The server reads the first 4 bytes, realizes the payload is exactly 50 bytes long, and then reads exactly 50 bytes before stopping!
