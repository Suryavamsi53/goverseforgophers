# The Request/Response Cycle (TCP/IP)

When your Go web server uses the `net/http` package to listen on port 8080, it abstracts away layers of immense network complexity.

To write high-performance Go applications, you must understand exactly what happens under the hood when a user opens `http://your-go-app.com`.

## 1. The OSI Model

The internet operates in layers (The OSI Model).
1. **Layer 3 (Network - IP)**: Gets the data from Computer A to Computer B using IP Addresses.
2. **Layer 4 (Transport - TCP)**: Ensures the data arrives in the correct order, without missing any packets.
3. **Layer 7 (Application - HTTP)**: The actual text data (Headers, Body) being sent.

## 2. The TCP 3-Way Handshake

Before a single byte of HTTP data can be sent, the client and the Go server must establish a TCP connection. This is the **3-Way Handshake**:

1. **SYN**: Client says "I want to connect."
2. **SYN-ACK**: Go Server says "I acknowledge your request, let's connect."
3. **ACK**: Client says "I acknowledge your acknowledgement."

*Performance Impact*: This handshake takes 1 full Round Trip Time (RTT). If the user is in Tokyo and the Go server is in New York, the handshake takes ~200 milliseconds before the HTTP request even begins!

## 3. The HTTP Request

Once the TCP connection is open, the Client sends plain text across the wire:

```text
GET /users/42 HTTP/1.1
Host: your-go-app.com
User-Agent: Mozilla/5.0
Accept: application/json

```
*(Notice the double line break at the end. That tells the Go server the headers are finished and the body is starting!)*

## 4. The Go `net/http` Architecture

How does Go handle this incoming text?

1. The `net/http` server runs an infinite `for` loop, calling `listener.Accept()`.
2. When the TCP connection is accepted, **Go spawns a brand new Goroutine specifically for this connection!**
3. Inside the Goroutine, Go parses the raw text into the `http.Request` struct.
4. Go passes this struct to your `ServeHTTP` handler function.
5. Your handler writes data to the `http.ResponseWriter`.
6. Go serializes that data back into plain text and sends it over the TCP socket.

**The Go Superpower:** In Node.js or Python, handling 10,000 concurrent requests requires complex asynchronous callbacks or Event Loops. In Go, because every request gets its own lightweight Goroutine, you write simple, synchronous, blocking code, and the Go runtime handles the concurrency magically under the hood!

## 5. Keep-Alive (Connection Pooling)

In the old days of HTTP/1.0, the Go server would close the TCP connection immediately after sending the response.

If the browser needed to fetch an image 10 milliseconds later, it had to perform the expensive 3-Way Handshake all over again!

In HTTP/1.1, the `Connection: keep-alive` header is the default. The Go server leaves the TCP connection open! The browser can send 50 sequential HTTP requests over the exact same TCP socket, saving massive amounts of network latency!
