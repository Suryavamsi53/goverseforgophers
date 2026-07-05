# The 8 Fallacies of Distributed Computing

In 1994, L. Peter Deutsch and James Gosling (the inventor of Java) formulated the "8 Fallacies of Distributed Computing". These are the dangerous assumptions that developers make when moving from a Monolithic architecture to Microservices.

If you architect a Go system based on any of these fallacies, your system will inevitably collapse in Production.

## 1. "The network is reliable."
* **The Reality**: Networks drop packets, routers reboot, and DNS servers fail. 
* **The Fix**: You must architect your Go HTTP clients to handle connection drops. Never assume a `200 OK` is guaranteed. Use the `Circuit Breaker` pattern.

## 2. "Latency is zero."
* **The Reality**: A local function call takes 1 nanosecond. An HTTP request across a datacenter takes 5,000,000 nanoseconds. If your microservice makes 10 sequential HTTP requests to render a single page, the user will experience massive lag.
* **The Fix**: Use `x/sync/errgroup` or Goroutines to fire those 10 HTTP requests concurrently (Fan-Out).

## 3. "Bandwidth is infinite."
* **The Reality**: If 10,000 users request a 5MB JSON payload simultaneously, you will saturate the network interface card on your Kubernetes node.
* **The Fix**: Always use Pagination (`LIMIT` / `OFFSET`), compress payloads (gzip/brotli), and use gRPC (which uses ultra-compact binary Protobuf instead of bloated JSON).

## 4. "The network is secure."
* **The Reality**: If your microservices communicate over raw HTTP inside a Kubernetes cluster, a compromised pod can sniff the traffic and steal passwords.
* **The Fix**: Implement mTLS (Mutual TLS) using a Service Mesh like Istio. Ensure every internal request is encrypted.

## 5. "Topology doesn't change."
* **The Reality**: IP addresses are completely ephemeral. In Kubernetes, a pod's IP address might change 15 times a day as it gets killed and rescheduled.
* **The Fix**: Never hardcode IP addresses. Always use Service Discovery (e.g., querying `http://billing-service.default.svc.cluster.local`).

## 6. "There is one administrator."
* **The Reality**: In a monolithic app, the DB admin knows exactly what schema changes happen. In Microservices, Team A might update the `Order` JSON schema without telling Team B, completely breaking Team B's parser.
* **The Fix**: Strict API Versioning (`/api/v1/orders`), API Gateways, and strongly-typed Protobuf contracts.

## 7. "Transport cost is zero."
* **The Reality**: Cloud providers charge you for data transfer. Sending huge JSON payloads between AWS Availability Zones can cost thousands of dollars a month.
* **The Fix**: Use edge caching (Redis/CDN) to prevent data from traversing the network in the first place.

## 8. "The network is homogeneous."
* **The Reality**: Your system might have Go on Linux, Node on Windows, and Swift on iOS. Relying on language-specific serialization (like Go's `encoding/gob`) will break interoperability.
* **The Fix**: Stick to universal standards (JSON, Protobuf, gRPC).
