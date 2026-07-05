# gRPC and Protocol Buffers

In the previous lesson, we learned *why* gRPC is replacing REST for internal microservice communication. Now we will look at how it actually works.

## 1. The Protocol Buffer Contract (`.proto`)

Everything in gRPC starts with a language-agnostic `.proto` file. This file defines the exact structure of your data and the functions your API exposes.

```proto
syntax = "proto3";
package order;

// 1. Define the Data Structures (Messages)
message CreateOrderRequest {
  string user_id = 1;
  float total_amount = 2;
  repeated string item_ids = 3; // 'repeated' means array/slice
}

message OrderResponse {
  string order_id = 1;
  bool success = 2;
}

// 2. Define the API (Service)
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (OrderResponse);
}
```

**The Field Tags (1, 2, 3):**
Notice the `= 1`, `= 2` tags. These are the secret to Protobuf's speed. In JSON, if you send `{"user_id": "bob"}`, it transmits the 7 bytes of the string "user_id" over the wire. Protobuf doesn't send the field names! It only sends the integer `1` and the data `"bob"`. The receiving client uses the `.proto` file to decode `1` back into `user_id`.

## 2. Code Generation

Once you write the `.proto` file, you run a compiler (the `protoc` CLI tool).
If your Backend is Go and your Frontend is TypeScript, you run `protoc` twice. 

It will instantly generate hundreds of lines of perfectly typed Go code and TypeScript code, handling all the network serialization, TCP connections, and Context timeout propagation for you. 

```go
// The generated Go Client code usage:
client := pb.NewOrderServiceClient(grpcConn)

// It looks like a local function call, but it's making a network request!
resp, err := client.CreateOrder(context.Background(), &pb.CreateOrderRequest{
    UserId: "bob",
    TotalAmount: 99.99,
})
```

## 3. Streaming (The Killer Feature)

REST is strictly Request/Response. If you need a continuous feed of data, you have to implement websockets manually. 

gRPC supports **Streaming** natively over HTTP/2.

* **Server Streaming**: Client sends 1 request, Server returns a continuous stream of responses (e.g., Live Stock Ticker).
* **Client Streaming**: Client streams gigabytes of data to the Server in chunks, Server returns 1 response when done (e.g., File Upload).
* **Bidirectional Streaming**: Both stream data simultaneously (e.g., Multiplayer Game, Chat App).

```proto
// Example of Server Streaming (notice the 'stream' keyword)
rpc SubscribeToStockPrices(StockRequest) returns (stream StockPrice);
```

## 4. Backwards Compatibility

If you change an API in REST, you often break the clients.
Protobuf makes backwards compatibility incredibly easy:
1. **Never change the integer tags of existing fields.**
2. **Never delete a tag number** (If you remove a field, mark it as `reserved` so no future developer accidentally reuses that integer).
3. **You can add new fields safely.** Old clients will simply ignore fields they don't recognize, and new clients will handle them.
