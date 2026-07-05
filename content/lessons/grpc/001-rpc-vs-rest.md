# RPC vs REST (The API Evolution)

For the last 15 years, **REST** (Representational State Transfer) over HTTP/1.1 has been the undisputed standard for building web APIs.

However, in massive microservice architectures (where an API Gateway might trigger a chain of 15 internal server-to-server calls just to render one page), the overhead of REST becomes a critical bottleneck.

Enterprise Go teams use **gRPC** for internal server-to-server communication.

## 1. The Problems with REST

1. **JSON is Slow and Heavy**: JSON is a text format. It must be serialized by the sender and parsed by the receiver. This CPU overhead is massive. Furthermore, the payload is bloated with repeated keys: `{"id": 1, "name": "Bob"}, {"id": 2, "name": "Alice"}`.
2. **HTTP/1.1 is Inefficient**: HTTP/1.1 opens a TCP connection, sends the request, waits for the response, and then closes the connection (or keeps it alive for sequential use). It suffers from Head-of-Line blocking. You cannot send 10 requests simultaneously over a single HTTP/1.1 connection.
3. **No Strict Contract**: If the `User-Service` team changes the JSON field from `name` to `first_name`, the `Billing-Service` team won't know until their Go code parses the JSON and crashes in Production. OpenAPI/Swagger helps, but it is often out-of-sync with the actual code.

## 2. What is RPC?

**RPC (Remote Procedure Call)** is a design paradigm. 
Instead of making an HTTP `GET /users/42` request and parsing JSON, RPC attempts to make the remote network call look exactly like a local function call in your code!

```go
// In REST, you write complex network code:
resp, _ := http.Get("http://user-service/users/42")
var u User
json.NewDecoder(resp.Body).Decode(&u)

// In RPC, the network is abstracted away! It looks like a local function!
user, _ := userServiceClient.GetUser(ctx, &GetUserRequest{ID: 42})
```

## 3. The gRPC Revolution

Google created **gRPC** (gRPC Remote Procedure Calls) to solve the bottlenecks of REST.

It provides three massive upgrades:

1. **Protocol Buffers (Protobuf)**: Instead of JSON text, gRPC uses Protobuf. The data is compiled into a highly-compressed, binary format before being sent over the wire. It is up to 10x faster to serialize/deserialize than JSON.
2. **HTTP/2**: gRPC exclusively uses HTTP/2. HTTP/2 maintains a single, permanent TCP connection and multiplexes thousands of binary requests concurrently over that single connection without Head-of-Line blocking!
3. **Strict Contracts (Code Generation)**: You define your API in a language-agnostic `.proto` file. The gRPC compiler automatically generates the Go structs, the Client network code, and the Server interface! If a team changes a field name, your Go code fails to compile! It is mathematically impossible for the API contract to drift from the codebase.

## 4. When to use REST vs gRPC?

* **Public APIs (External)**: Always use REST! Web browsers (React/Angular) and external third-party developers expect standard JSON over HTTP/1.1. Browsers have terrible native support for gRPC.
* **Internal APIs (Server-to-Server)**: Always use gRPC! When your Go API Gateway talks to your Go internal microservices, there is zero reason to use slow JSON. gRPC reduces CPU usage, slashes network latency, and enforces strict type-safety across your engineering teams.
