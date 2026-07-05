# Interceptors and Metadata

In standard HTTP REST APIs, if you want to add Authentication (JWTs), Logging, or Rate Limiting to every single endpoint, you use **Middleware**.

In gRPC, Middleware is called **Interceptors**.

## 1. Unary Interceptors (The Middleware)

An Interceptor intercepts the RPC call *before* it reaches your actual handler logic.

```go
// A standard logging Interceptor
func LoggingInterceptor(
    ctx context.Context, 
    req interface{}, 
    info *grpc.UnaryServerInfo, // Contains the name of the RPC method!
    handler grpc.UnaryHandler,  // The actual business logic function
) (interface{}, error) {
    
    start := time.Now()
    log.Printf("--> Calling RPC: %s", info.FullMethod)
    
    // 1. Pass control to the real handler
    resp, err := handler(ctx, req)
    
    // 2. Logic to run AFTER the handler finishes
    log.Printf("<-- Finished RPC: %s (took %v)", info.FullMethod, time.Since(start))
    
    return resp, err
}
```

To use it, you inject it into the Server constructor in `main.go`:
```go
s := grpc.NewServer(
    grpc.UnaryInterceptor(LoggingInterceptor),
)
```
*Note: To chain multiple interceptors together (Auth -> Logger -> Metrics), you must use the `github.com/grpc-ecosystem/go-grpc-middleware` library, which provides the `grpc.ChainUnaryInterceptor` helper!*

## 2. Metadata (The gRPC Headers)

In HTTP REST, you pass the JWT Authentication token using an HTTP Header (`Authorization: Bearer <token>`).
Because gRPC abstracts away HTTP, you cannot access HTTP headers directly!

gRPC replaces HTTP Headers with **Metadata** (a simple key-value map).

**The Client (Sending Metadata):**
```go
// Create a map of metadata
md := metadata.Pairs("authorization", "Bearer my-secret-jwt")

// Inject the metadata into a NEW Context!
ctx := metadata.NewOutgoingContext(context.Background(), md)

// Make the RPC call using the new Context!
res, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 42})
```

**The Server (Reading Metadata in the Interceptor):**
```go
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    
    // 1. Extract the Metadata from the incoming Context
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
    }

    // 2. Read the specific key (Notice it returns an array of strings!)
    tokens := md["authorization"]
    if len(tokens) == 0 || tokens[0] != "Bearer my-secret-jwt" {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token")
    }

    // 3. Auth passed! Continue to the handler!
    return handler(ctx, req)
}
```

## 3. Context Propagation (The Enterprise Standard)

Because gRPC strictly enforces that `context.Context` is the first parameter of every single method, passing Tracing IDs or Metadata through a 5-layer deep microservice architecture is mathematically guaranteed to work.

If the API Gateway receives a Trace ID, it injects it into the Context. When the Gateway makes a gRPC call to the Billing Service, the gRPC library *automatically* serializes the Context Metadata, ships it over the wire, and injects it into the Billing Service's Context!
