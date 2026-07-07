# Building a Reverse Proxy in Go

## 1. Learning Objectives
* **What you'll learn**: How to use Go's `httputil.ReverseProxy` to intercept, modify, and forward HTTP traffic.
* **Why it matters**: Reverse proxies are the backbone of modern cloud architecture. Load balancers (Nginx, HAProxy), API Gateways (Kong), and Service Meshes (Envoy) are all fundamentally reverse proxies.
* **Where it's used**: Microservice routing, rate limiters, authentication gateways, and A/B testing routers.

---

## 2. What is a Reverse Proxy?
A standard Proxy (Forward Proxy) sits in front of a **Client** (like a corporate VPN hiding your IP from the internet).
A **Reverse Proxy** sits in front of a **Server**. 
When a client makes a request to `api.goverse.com`, the request hits the Reverse Proxy first. The proxy looks at the request and decides:
1. *Is this user authenticated?*
2. *Should I route this to the Go Backend or the Node.js backend?*
3. *Should I block this IP?*

---

## 3. The Magic of httputil.ReverseProxy
Go's standard library includes a production-ready reverse proxy out of the box in the `net/http/httputil` package. It automatically handles transferring headers, streaming large bodies without blowing up memory, and connection pooling.

### Basic API Gateway Example
Let's build a gateway that intercepts requests on port `8080` and forwards them to a microservice running on port `9090`.

```go
package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// 1. Define the destination microservice
	targetURL, _ := url.Parse("http://localhost:9090")

	// 2. Create the Reverse Proxy
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 3. Customize the Proxy (Optional but powerful)
	// The Director allows you to modify the request BEFORE it is sent to the target
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		
		// Example: Inject a secret header that the microservice requires
		req.Header.Set("X-Gateway-Secret", "super-secret-key")
		
		// Example: Log the incoming IP
		log.Printf("Forwarding request from %s to %s", req.RemoteAddr, req.URL.Path)
	}

	// 4. Start the Gateway server
	log.Println("API Gateway listening on :8080...")
	http.ListenAndServe(":8080", proxy)
}
```

---

## 4. Modifying Responses
You can also intercept the response coming *back* from the microservice before sending it to the client using `ModifyResponse`.

```go
proxy.ModifyResponse = func(resp *http.Response) error {
    // Check if the microservice returned a 500 error
    if resp.StatusCode == http.StatusInternalServerError {
        log.Println("Microservice crashed! Sending a friendly error instead.")
        // You could theoretically rewrite the response body here
    }
    
    // Add a custom security header to all responses
    resp.Header.Set("X-Frame-Options", "DENY")
    return nil
}
```

---

## 5. Quiz

1. **MCQ**: What is the primary benefit of using `httputil.ReverseProxy` instead of just writing `http.Get()` inside your handler and copying the response?
   * (A) It bypasses CORS restrictions.
   * (B) It automatically handles header translation, connection pooling, and stream copying (preventing memory exhaustion on large file uploads). *(Answer: B)*
   * (C) It compiles down to WebAssembly.

2. **System Design Follow-up**: If you were building a Load Balancer in Go, how would you modify the `proxy.Director` to route traffic to 5 different servers?
   * *(You wouldn't use `NewSingleHostReverseProxy`. You would create a custom `httputil.ReverseProxy` where your `Director` function implements an algorithm (like Round Robin or Least Connections) to dynamically select one of the 5 Target URLs and rewrite the `req.URL.Host` to match it.)*
