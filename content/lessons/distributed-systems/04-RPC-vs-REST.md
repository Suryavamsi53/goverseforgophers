# RPC vs REST

When building microservices, you must choose how they communicate. For the last decade, REST (Representational State Transfer) over HTTP/1.1 has been the dominant standard. 

However, modern backend engineering is rapidly shifting towards RPC (Remote Procedure Call).

## 1. The REST Architecture

REST treats the system as a collection of **Resources** (e.g., `/users`, `/orders`). You manipulate these resources using standard HTTP verbs:
* `POST /users` (Create)
* `GET /users/1` (Read)
* `PUT /users/1` (Update)
* `DELETE /users/1` (Delete)

### The Pros of REST
* Universal: Every language and browser natively understands HTTP and JSON.
* Easy to debug: You can use `curl` or Postman to instantly read the payloads.

### The Cons of REST
* **Bloat**: JSON is a text format. Sending the number `123456789` in JSON takes 9 bytes of ASCII text. In binary, it takes 4 bytes.
* **No Contracts**: A REST endpoint returns JSON. If the backend team accidentally changes the `userId` field to `user_id`, the client's Go struct parser will silently break. (Swagger/OpenAPI attempts to fix this, but it is often out-of-sync).
* **Action Routing**: How do you map a non-CRUD action to a REST endpoint? (e.g., "Calculate Tax" or "Restart Server"). You usually end up violating REST principles by creating endpoints like `POST /users/1/calculate-tax`.

## 2. The RPC Architecture

RPC ignores "Resources". It treats remote servers as if they were local Go functions.

Instead of hitting `POST /orders`, your client literally calls a function like `OrderService.CreateOrder()`. The RPC framework serializes the function arguments, sends them over the network, executes the function on the remote server, and returns the result.

### The Rise of gRPC
Google created **gRPC**, which uses **HTTP/2** for transport and **Protocol Buffers (Protobuf)** for serialization.

### The Pros of gRPC
* **Strict Contracts**: You define your API in a `.proto` file. Both the Client and the Server auto-generate their Go code from this single source of truth. It is mathematically impossible for the frontend to expect `userId` if the backend changed it to `user_id`. The compiler will catch it.
* **Speed**: Protobuf is binary. It is ~5x faster to serialize and ~50% smaller over the network than JSON.
* **HTTP/2 Multiplexing**: gRPC allows you to send thousands of concurrent requests over a single, long-lived TCP connection, completely eliminating the TCP Handshake overhead.

## 3. When to use Which?

* **Use REST**: For public-facing APIs (like the Twitter or Stripe API) where you want any developer to easily connect using standard tools. For Web UI frontends (though gRPC-Web is gaining traction).
* **Use gRPC**: For **Internal Microservice-to-Microservice** communication. If your Go `BillingService` needs to talk to your Go `OrderService`, using JSON/REST is an unacceptable waste of CPU and bandwidth.
