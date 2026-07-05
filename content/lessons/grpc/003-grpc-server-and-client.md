# gRPC Server and Client in Go

Once you have defined your `.proto` file and run the `protoc` compiler, generating the Go code, it's time to actually build the Server and the Client.

## 1. Implementing the Server

The generated Go code contains an Interface (e.g., `UserServiceServer`) that dictates exactly what methods your server must implement.

```go
package main

import (
    "context"
    "net"
    "google.golang.org/grpc"
    pb "myapp/internal/pb" // The generated Protobuf package!
)

// 1. Create a struct that will implement the generated interface
type server struct {
    // Required for forward compatibility
    pb.UnimplementedUserServiceServer 
}

// 2. Implement the core business logic!
func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    // Extract the ID from the strongly-typed generated struct
    userID := req.GetId()
    
    // Simulate database lookup
    if userID != 42 {
        // gRPC uses specific Status Codes, not HTTP 404!
        return nil, status.Errorf(codes.NotFound, "user not found")
    }

    // Return the strongly-typed generated response struct
    return &pb.GetUserResponse{
        Id:       userID,
        Email:    "test@test.com",
        IsActive: true,
    }, nil
}

func main() {
    // 3. Open a standard TCP port
    lis, _ := net.Listen("tcp", ":50051")
    
    // 4. Create a new gRPC Server instance
    s := grpc.NewServer()
    
    // 5. Register our implementation with the gRPC Server
    pb.RegisterUserServiceServer(s, &server{})
    
    // 6. Start listening!
    s.Serve(lis)
}
```

## 2. Implementing the Client

If you have ever written an HTTP Client using `net/http` to parse a massive JSON REST API, you know it takes dozens of lines of code.

With gRPC, the `protoc` compiler generated the *entire* client for you!

```go
package main

import (
    "context"
    "log"
    "google.golang.org/grpc"
    pb "myapp/internal/pb"
)

func main() {
    // 1. Open a permanent HTTP/2 TCP connection to the server
    // (In production, use grpc.WithTransportCredentials for TLS!)
    conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
    defer conn.Close()

    // 2. Instantiate the auto-generated Client
    client := pb.NewUserServiceClient(conn)

    // 3. Make the RPC call! It looks exactly like a local function call!
    req := &pb.GetUserRequest{Id: 42}
    res, err := client.GetUser(context.Background(), req)
    
    if err != nil {
        log.Fatalf("RPC failed: %v", err)
    }

    // Access the strongly-typed fields! No JSON decoding required!
    log.Printf("User Email: %s", res.GetEmail())
}
```

## 3. Connection Multiplexing (The Golden Rule)

In the Client code above, we called `grpc.Dial()`.
**You must never call `grpc.Dial()` on every request!**

In REST, you often open and close HTTP connections constantly. 
gRPC relies on HTTP/2. The entire purpose of HTTP/2 is that it opens ONE single TCP connection and holds it open forever. You multiplex thousands of concurrent RPC calls over that one connection using Goroutines.

You should call `grpc.Dial()` exactly *once* when your Go application boots up (in `main.go`), and share that `conn` object across all your HTTP Handlers!
