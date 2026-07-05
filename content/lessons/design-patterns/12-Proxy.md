# Proxy Pattern

The Proxy Pattern is a structural design pattern that provides a surrogate or placeholder for another object to control access to it.

A Proxy intercepts calls to the original object, allowing you to perform actions *before* or *after* the request reaches the original object.

*(If you are thinking: "Wait, isn't that exactly what the Decorator pattern does?" You are correct. Structurally, they are almost identical in Go. The difference is entirely in **Intent**).*

## 1. Intent: Proxy vs Decorator

* **Decorator**: Adds new responsibilities to an object dynamically (e.g., Caching, Logging). The caller *knows* they are using a decorated object, and the Decorator is usually passed in via the constructor.
* **Proxy**: Controls *access* to an object. The caller has no idea the Proxy exists; they think they are talking to the real object. The Proxy often handles lifecycle management (lazy loading) or security (access control).

## 2. Use Case: Lazy Initialization (Virtual Proxy)

Imagine an application that connects to an incredibly slow, legacy SAP database. 
If we connect to it during `main()` startup, the application will take 30 seconds to boot, even if the user never actually queries the database!

We can use a Virtual Proxy to delay the initialization of the database until the exact millisecond the user actually requests data.

```go
type Database interface {
    Query(sql string) string
}

// 1. The Heavy Object
type SAPDatabase struct{}
func (s *SAPDatabase) Query(sql string) string { return "SAP Data" }

// 2. The Proxy Object
type DatabaseProxy struct {
    realDB *SAPDatabase // Starts as nil!
}

// 3. The Interceptor Method
func (p *DatabaseProxy) Query(sql string) string {
    // Lazy Initialization: Only connect if we haven't already!
    if p.realDB == nil {
        fmt.Println("Connecting to SAP (takes 30 seconds)...")
        // Simulate heavy connection load
        time.Sleep(2 * time.Second)
        p.realDB = &SAPDatabase{}
    }
    
    // Forward the request to the real object
    return p.realDB.Query(sql)
}
```

## 3. Use Case: Access Control (Protection Proxy)

A Protection Proxy intercepts method calls and checks if the current user has the correct permissions (JWT roles) to execute the action.

```go
type Document interface {
    Read() string
}

type SecretDocument struct{}
func (s *SecretDocument) Read() string { return "Top Secret Classified Data" }

type SecurityProxy struct {
    doc      Document
    userRole string
}

func (p *SecurityProxy) Read() string {
    if p.userRole != "ADMIN" {
        return "ACCESS DENIED"
    }
    return p.doc.Read()
}
```

## 4. Proxies in Microservices

The most famous real-world example of this pattern is the **API Gateway** or the **Envoy Sidecar Proxy** in a Kubernetes Service Mesh. 

When your Go application makes a network request to `http://billing-service`, it doesn't actually talk to the Billing Service. The Envoy Proxy intercepts the request, handles mTLS encryption, enforces rate limits, and then forwards it to the real Billing Service. The Go application is completely unaware of the Proxy's existence!
