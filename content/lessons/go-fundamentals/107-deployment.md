# Deployment & Kubernetes

Because Go compiles to a single static binary, deploying it is incredibly simple. However, in modern enterprise architectures, simply copying a binary to a server via SCP is not enough. You need orchestration.

## 1. Kubernetes Integration

Go is the native language of the Cloud Native ecosystem (Kubernetes itself is written in Go). Deploying a Go web server into Kubernetes requires setting up **Probes**.

Kubernetes needs to know two things about your Go app:
1. **Liveness**: Is the app completely dead/deadlocked? (If yes, restart it).
2. **Readiness**: Is the app temporarily busy or still booting up? (If yes, don't send it web traffic yet).

To support this, your Go server should expose a dedicated health endpoint.

```go
func healthzHandler(w http.ResponseWriter, r *http.Request) {
    // Check database connection pool
    if err := db.Ping(); err != nil {
        // App is alive, but database is down. Return 503 Service Unavailable.
        // Kubernetes will temporarily stop sending traffic to this pod.
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("Database unreachable"))
        return
    }
    
    // Everything is healthy!
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}

func main() {
    mux := http.NewServeMux()
    // Typically mounted on a separate, internal port so public users can't hit it!
    mux.HandleFunc("/healthz", healthzHandler)
}
```

## 2. Configuration via ConfigMaps

In the Configuration lesson, we used Viper to read environment variables. In Kubernetes, you deploy a `ConfigMap` and a `Secret` which inject those variables directly into your Docker container's environment.

```yaml
# kubernetes-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-api
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: go-api
        image: myregistry/myapp:v1.0.0
        env:
        - name: PORT
          value: "8080"
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url
```

## 3. The Deployment Rollout

Because Go apps boot instantly (usually in under 20 milliseconds) and have tiny 10MB footprints, Kubernetes can perform rolling updates flawlessly. 

When you deploy `v2.0.0`, Kubernetes spins up the new Go containers, waits for the `/healthz` endpoint to return `200 OK` (which happens instantly), and then routes traffic over. Meanwhile, your `v1.0.0` containers receive the `SIGTERM` signal and execute their Graceful Shutdown (as covered in Lesson 102), draining active users. 

**Result: True Zero-Downtime Deployments.**
