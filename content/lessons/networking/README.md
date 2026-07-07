# Networking & Socket Programming in Go

Welcome to the **Networking** curriculum module.
At its core, almost every Go backend application is a network server. Understanding how TCP, UDP, and HTTP protocols operate at the socket level is essential for building load balancers, proxies, real-time gaming servers, and custom database drivers.

## Curriculum

1. [Lesson 1: TCP Sockets in Go](001-tcp-sockets-in-go.md)
   - `net.Listen` and `net.Dial`
   - Managing connection lifecycles and timeouts
   - Handling stream fragmentation

2. [Lesson 2: UDP and WebSockets](002-udp-websockets.md)
   - Connectionless protocols vs Streaming protocols
   - Broadcasting real-time state with WebSockets

3. [Lesson 3: Building a Reverse Proxy](003-reverse-proxy.md)
   - `httputil.ReverseProxy`
   - Intercepting and modifying HTTP requests/responses

## Why this matters?
Frameworks like Gin or Fiber hide the complexity of the network. But when a TCP connection drops mid-stream, or you suffer from socket exhaustion (TIME_WAIT), you need to understand the underlying OS networking primitives to debug the issue.
