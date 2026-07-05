# TLS and mTLS (Encryption in Transit)

If your Go application sends JSON payloads over plain HTTP, the data is sent as raw text across the internet. 
Any router, ISP, or hacker sitting between the client and your server can read the credit card numbers in plain text (A Man-in-the-Middle Attack).

To prevent this, you must encrypt the network connection using **TLS (Transport Layer Security)**, which gives us HTTPS.

## 1. How TLS Works (The Handshake)

TLS uses a brilliant combination of Asymmetric (Public/Private) and Symmetric cryptography.

1. **The Certificate**: Your Go server holds a Private Key and an SSL Certificate (which contains the Public Key, digitally signed by a trusted Certificate Authority like Let's Encrypt).
2. **Asymmetric Handshake**: The client connects and downloads the Certificate. The client generates a random "Master Secret", encrypts it using the server's Public Key, and sends it back. Only the server's Private Key can decrypt this Master Secret!
3. **Symmetric Encryption**: Both sides now possess the exact same Master Secret. They use it to generate Symmetric Keys (like AES-GCM) to encrypt all further HTTP traffic. (Symmetric encryption is used because it is 10,000x faster than Asymmetric).

## 2. Terminating TLS in Go

While it is standard practice to let an API Gateway (like NGINX or Kubernetes Ingress) handle TLS termination, you can easily boot an HTTPS server directly in Go:

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Encrypted Hello!"))
    })

    // ListenAndServeTLS requires the Certificate and the Private Key files!
    log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", mux))
}
```

## 3. Mutual TLS (mTLS) - Zero Trust Architecture

Standard TLS only authenticates the Server. The client asks, *"Are you really Google?"* and the server proves it with a certificate. The server has no idea who the client is.

In a Microservice architecture, what if a hacker breaches your Kubernetes cluster and tries to make an API call directly to the `Billing Service`? The Billing Service needs to know exactly which microservice is calling it.

This requires **mTLS (Mutual TLS)**.
In mTLS, both the Server AND the Client must present TLS Certificates to each other! 

1. `Order Service` sends its Client Certificate to `Billing Service`.
2. `Billing Service` verifies the Client Certificate.
3. `Billing Service` sends its Server Certificate to `Order Service`.
4. `Order Service` verifies the Server Certificate.

### Istio and the Service Mesh
Managing and rotating thousands of TLS certificates for 50 Go microservices is impossible. 
Enterprise teams use a **Service Mesh** (like Istio or Linkerd). 
The Service Mesh injects an Envoy sidecar proxy into every Pod. The Envoy proxies automatically handle the mTLS handshake, certificate rotation, and encryption seamlessly. Your Go code just sends plain HTTP to `localhost`, completely unaware that Envoy is transparently encrypting it before it leaves the Pod! This is the foundation of **Zero Trust Architecture**.
