# UDP & WebSockets in Go

## 1. Learning Objectives
* **What you'll learn**: The difference between TCP (stateful) and UDP (stateless), and how to use WebSockets to upgrade a standard HTTP connection into a bi-directional real-time stream.
* **Why it matters**: Standard HTTP requests are unidirectional (Client asks, Server responds). Real-time apps like multiplayer games, chat apps, and live dashboards require the Server to push data to the Client without being asked.
* **Where it's used**: Video streaming, VoIP, live cryptocurrency price tickers, and collaborative editing tools (like Google Docs).

---

## 2. UDP (User Datagram Protocol)
While TCP guarantees that packets arrive in order and without corruption, UDP makes **no guarantees**. It fires the packet into the network and forgets about it. 

Why use UDP? Because dropping the overhead of handshakes and ordering makes it incredibly fast. It is heavily used in FPS multiplayer gaming (if you drop a packet containing player movement, you don't care, because a new packet with updated movement will arrive 10 milliseconds later anyway).

### Go UDP Example

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	// Listen for UDP packets on port 8081
	addr := net.UDPAddr{
		Port: 8081,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		// ReadFromUDP does not block waiting for a handshake
		n, clientAddr, _ := conn.ReadFromUDP(buffer)
		fmt.Printf("Received %s from %s\n", string(buffer[:n]), clientAddr)
		
		// Fire a response back instantly
		conn.WriteToUDP([]byte("Packet received!"), clientAddr)
	}
}
```

---

## 3. WebSockets (Upgrading HTTP)
WebSockets solve the problem of real-time communication in the browser. They start as a standard HTTP `GET` request with an `Upgrade: websocket` header. If the Go server agrees, the connection is kept alive indefinitely, and both sides can push messages to each other.

To use WebSockets in Go, the community standard is the `github.com/gorilla/websocket` package.

### Go WebSocket Upgrade Example

```go
package main

import (
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

// The Upgrader dictates buffer sizes and CORS policies
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for this example
	},
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 1. Upgrade the HTTP connection to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 2. Enter a continuous loop to read and write messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read failed:", err)
			break
		}
		
		log.Printf("Received: %s", message)
		
		// Echo the message back to the client
		err = conn.WriteMessage(messageType, []byte("Server says: "+string(message)))
		if err != nil {
			log.Println("Write failed:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 4. Quiz

1. **MCQ**: What happens to a standard HTTP connection when a WebSocket upgrade is successful?
   * (A) It is closed immediately and a new UDP port is opened.
   * (B) It is kept alive indefinitely and hijacked by the WebSocket protocol. *(Answer: B)*
   * (C) It times out after 30 seconds.

2. **System Design Follow-up**: If you have 100,000 users connected via WebSockets, and you want to broadcast a chat message to all of them, why shouldn't you loop through `conn.WriteMessage()` for all 100,000 connections sequentially?
   * *(Because writing to a socket involves blocking OS system calls. A sequential loop would take seconds to finish, creating massive latency for the users at the end of the array. You must use Goroutines or worker pools to write to connections concurrently).*
