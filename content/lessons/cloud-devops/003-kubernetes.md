# Kubernetes (K8s)

## 1️⃣ Learning Objectives
* **What you'll learn**: Master the core concepts of Kubernetes: Pods, Deployments, Services, ConfigMaps, Secrets, and how to properly deploy and scale Go applications.
* **Why it matters**: Kubernetes is the industry standard for container orchestration. It automates deployment, scaling, and management of containerized applications.
* **Where it's used**: Almost every modern enterprise uses Kubernetes to run microservices in production.

---

## 2️⃣ Real-world Story
Imagine a massive shipping port. 
Before containers (Docker), you had to pack your goods (app) into specific shapes to fit onto a specific truck (server). 
With Docker, everything goes into standard-sized steel shipping containers. 

But who coordinates thousands of shipping containers? Who decides which ship they go on, what happens if a ship sinks, and how to route trucks to pick them up? That's **Kubernetes**. It is the port authority, the crane operators, and the logistics engine all rolled into one. You just give Kubernetes a blueprint ("I need 5 containers of my Go app running at all times"), and Kubernetes makes it happen and *keeps* it happening.

---

## 3️⃣ Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    subgraph Control Plane (Master)
        API[API Server]
        SCHED[Scheduler]
        CM[Controller Manager]
        ETCD[(etcd)]
    end

    subgraph Worker Node 1
        KLET1[Kubelet]
        PROXY1[Kube-Proxy]
        POD1((Pod: Go App))
        POD2((Pod: Go App))
    end

    subgraph Worker Node 2
        KLET2[Kubelet]
        PROXY2[Kube-Proxy]
        POD3((Pod: Go App))
    end

    API --- KLET1
    API --- KLET2
    KLET1 --- POD1
    KLET1 --- POD2
    KLET2 --- POD3
    USER(User) -->|kubectl apply| API
```

---

## 4️⃣ Core Concepts (The Kubernetes Vocabulary)
* **Pod**: The smallest deployable unit. Usually encapsulates one container (e.g., your Go app), but can contain multiple tightly coupled containers.
* **Deployment**: A declarative way to manage Pods. It ensures a specified number of identical Pods are running (ReplicaSet) and handles rolling updates.
* **Service**: An abstract way to expose an application running on a set of Pods as a network service. Since Pod IPs change, Services provide a stable IP/DNS name.
* **Ingress**: Manages external access to the services in a cluster, typically HTTP/HTTPS. It acts as an API gateway or reverse proxy.
* **ConfigMap & Secret**: Ways to inject configuration data and sensitive credentials into your Pods without hardcoding them in the Docker image.

---

## 5️⃣ Code Examples: Deploying a Go App

### 🔹 Step 1: The Go Application
A standard Go server that gracefully handles termination (crucial for Kubernetes!).
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{Addr: ":8080"}

	go func() {
		fmt.Println("Server listening on :8080")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Kubernetes sends SIGTERM
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
```

### 🔹 Step 2: The Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/myapp .
EXPOSE 8080
CMD ["./myapp"]
```

### 🔹 Step 3: Kubernetes Deployment Manifest (`deployment.yaml`)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-app-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-app
  template:
    metadata:
      labels:
        app: go-app
    spec:
      containers:
      - name: go-app
        image: my-docker-repo/go-app:v1.0.0
        ports:
        - containerPort: 8080
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 10
```

---

## 6️⃣ Production Features: Probes & Resources
* **Liveness Probe**: Tells Kubernetes if your app is alive. If it fails, Kubernetes kills the Pod and restarts it.
* **Readiness Probe**: Tells Kubernetes if your app is ready to accept traffic. If it fails, Kubernetes removes the Pod from the Service load balancer (but doesn't kill it).
* **Resource Requests**: The minimum CPU/Memory guaranteed for your Pod. The Scheduler uses this to find a Node with enough space.
* **Resource Limits**: The maximum CPU/Memory your Pod can use. If it exceeds Memory limits, the Linux kernel will **OOMKill** (Out Of Memory Kill) it.

---

## 7️⃣ Kubernetes Service Manifest (`service.yaml`)
To allow other pods (or the outside world via Ingress) to talk to your deployment.
```yaml
apiVersion: v1
kind: Service
metadata:
  name: go-app-service
spec:
  selector:
    app: go-app
  ports:
    - protocol: TCP
      port: 80       # Port exposed by the service
      targetPort: 8080 # Port your Go app is listening on
  type: ClusterIP    # Internal only. Use LoadBalancer for cloud external IP.
```

---

## 8️⃣ Performance & Scaling
* **Horizontal Pod Autoscaler (HPA)**: Kubernetes can automatically scale the number of Pods up or down based on CPU utilization or custom metrics (like HTTP request rate).
  ```bash
  kubectl autoscale deployment go-app-deployment --cpu-percent=80 --min=3 --max=10
  ```
* **Go Advantage**: Go applications start up almost instantly and use very little memory. This makes them perfect for Kubernetes HPA, as they can scale from 1 to 100 pods in seconds to handle traffic spikes.

---

## 9️⃣ Best Practices
* ✅ **Do**: Listen for `SIGTERM` in your Go app to gracefully close database connections and finish inflight HTTP requests before shutting down.
* ✅ **Do**: Use multi-stage Docker builds to keep your Go image sizes under 20MB. This drastically speeds up Pod startup times.
* ❌ **Don't**: Write logs to local files. Write to `stdout`/`stderr` (using structured JSON logging) so Kubernetes can capture and ship them.

---

## 🔟 Common Mistakes & Debugging
### **Common Errors**
* **CrashLoopBackOff**: Your pod starts, crashes immediately, and Kubernetes keeps trying to restart it. Check `kubectl logs`.
* **OOMKilled**: Your Go app used more memory than the `limits.memory` defined in your YAML. Fix memory leaks or increase the limit.
* **ImagePullBackOff**: Kubernetes cannot download your Docker image (usually a typo in the name, or missing authentication for a private registry).

### **Essential Commands**
```bash
# See what is running
kubectl get pods
kubectl get services

# Read the logs of a specific pod
kubectl logs pod/go-app-deployment-xxxxxx

# Get detailed events to see why a pod is failing
kubectl describe pod go-app-deployment-xxxxxx

# Forward a port to your local machine for testing
kubectl port-forward svc/go-app-service 8080:80
```

---

## 11️⃣ FAANG Interview Questions
* **Beginner**: What is the difference between a Pod and a Container?
* **Intermediate**: What happens exactly when you delete a Pod that is managed by a Deployment?
* **Senior (Google/Meta)**: Explain how Kube-Proxy implements Services using iptables or IPVS. How would you design a zero-downtime deployment pipeline using Kubernetes Rolling Updates and Readiness Probes?
