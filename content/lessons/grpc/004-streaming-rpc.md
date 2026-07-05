# Streaming RPCs (WebSockets Alternative)

Standard RPCs (Unary) are just like HTTP: The client sends 1 Request, the server sends 1 Response.

But because gRPC is built on top of HTTP/2, it natively supports **Streaming** without the need for WebSockets! 
You can stream gigabytes of data endlessly in either direction over a single open TCP connection.

## 1. The 3 Types of Streams

In your `.proto` file, you simply add the `stream` keyword to define the endpoints:

```protobuf
service DataService {
    // 1. Server Streaming: Client sends 1 request, Server returns a stream of responses.
    rpc DownloadLargeFile(DownloadReq) returns (stream Chunk);

    // 2. Client Streaming: Client sends a stream of data, Server returns 1 response.
    rpc UploadTelemetry(stream Metric) returns (Summary);

    // 3. Bidirectional: Both sides stream data simultaneously! (Like WebSockets!)
    rpc ChatRoom(stream ChatMsg) returns (stream ChatMsg);
}
```

## 2. Server Streaming Implementation

Let's implement a Server Stream. The client asks for a large file, and the server chunks it into pieces and streams it back.

**The Server (Go):**
```go
// The generated interface gives us a specialized stream object!
func (s *server) DownloadLargeFile(req *pb.DownloadReq, stream pb.DataService_DownloadLargeFileServer) error {
    
    // We break the massive file into chunks
    for i := 0; i < 10; i++ {
        chunk := &pb.Chunk{Data: []byte("file chunk data...")}
        
        // Push the chunk down the open TCP connection!
        if err := stream.Send(chunk); err != nil {
            return err
        }
        time.Sleep(100 * time.Millisecond) // Simulate slow read
    }
    
    // Returning nil closes the stream successfully!
    return nil 
}
```

**The Client (Go):**
```go
// The client makes the initial request
stream, err := client.DownloadLargeFile(context.Background(), &pb.DownloadReq{FileId: 1})

for {
    // Block and wait for the next chunk from the server
    chunk, err := stream.Recv()
    
    // When the server returns nil, Recv() returns io.EOF!
    if err == io.EOF {
        fmt.Println("Download complete!")
        break
    }
    if err != nil {
        log.Fatalf("Stream broke: %v", err)
    }
    
    fmt.Println("Received chunk size:", len(chunk.Data))
}
```

## 3. Bidirectional Streaming (The Ultimate Power)

In a Chat Application, both sides need to push data simultaneously. 

To achieve this in Go, the Client must spin up a dedicated **Goroutine** to run the `stream.Recv()` loop in the background, while the main thread runs the `stream.Send()` loop! 

Because gRPC is perfectly thread-safe, you can effortlessly push and pull data from the exact same `stream` object concurrently across multiple Goroutines!
