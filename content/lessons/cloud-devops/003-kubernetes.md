# Kubernetes (K8s) Architecture for Go

## 1. Learning Objectives
* **What you'll learn**: The absolute basics of Kubernetes architecture (Pods, Deployments, Services) and how it orchestrates thousands of Go microservices in production.
* **Why it matters**: Docker Compose is for running apps on a single laptop. Kubernetes is for running apps across a cluster of 1,000 servers automatically healing them if they crash.
* **Where it's used**: The defacto operating system of the modern cloud. Over 90% of enterprises use K8s.

---

## 2. Real-world Story
Imagine hiring a fleet of 100 Uber drivers (Servers). If a driver's car breaks down, the passenger is stranded. You have to manually notice the failure, call a new driver, and redirect them. 
Kubernetes is the ultimate Dispatch AI. You tell the AI: "I always want exactly 10 cars on the road." The AI monitors them constantly. If a car explodes, the AI instantly dispatches a replacement car to take its place without you lifting a finger. It is the Master Orchestrator.

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Public Internet] -->|HTTP Request| B(LoadBalancer / Ingress)
    B -->|Routes to| C(K8s Service)
    
    subgraph Kubernetes Cluster (Worker Nodes)
        C -->|Load Balances to| D[Pod 1: Go App]
        C -->|Load Balances to| E[Pod 2: Go App]
        C -->|Load Balances to| F[Pod 3: Go App]
    end
    
    G[Control Plane / Master Node] -.->|Monitors Health| D
    G -.->|Monitors Health| E
    G -.->|Monitors Health| F
    
    style G fill:#8b5cf6,color:#fff
```

---

## 4. Internal Working (Under the Hood)
Kubernetes is divided into two parts:
1. **The Control Plane (Master)**: The brain. It contains the API Server (which accepts your YAML files), the Scheduler (which decides which server runs your Go app), and `etcd` (the database storing the cluster state).
2. **The Worker Nodes**: The physical servers. They run the `kubelet` agent, which physically starts your Docker containers and reports their health back to the Master.

---

## 5. Compiler Behavior
* **Native Integration**: Kubernetes itself is written 100% in Go! Because your backend is also written in Go, you can use the official `k8s.io/client-go` library to write Go applications that natively talk to the Kubernetes API, allowing your app to spin up other pods dynamically!

---

## 6. Memory Management
* **OOMKilled**: If you tell Kubernetes your Go Pod requires `100Mi` of memory, but your Go code suffers a memory leak and exceeds 100MB, the Linux kernel instantly terminates the process. You will see the dreaded `OOMKilled` status in Kubernetes. K8s will restart it, but it will keep dying until you fix the leak.

---

## 7. Code Examples

### 🔹 Example 1: The Deployment (Managing Pods)
```yaml
# deployment.yaml
# A Deployment guarantees exactly 3 instances (replicas) of your Go app are always running!
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goverse-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: goverse
  template:
    metadata:
      labels:
        app: goverse
    spec:
      containers:
      - name: go-backend
        image: myrepo/goverse-api:v1.2
        ports:
        - containerPort: 8080
```

### 🔹 Example 2: The Service (Networking)
```yaml
# service.yaml
# The Service acts as a stable IP address and Load Balancer for your 3 Pods.
apiVersion: v1
kind: Service
metadata:
  name: goverse-service
spec:
  selector:
    app: goverse
  ports:
    - protocol: TCP
      port: 80         # Port exposed internally in the cluster
      targetPort: 8080 # Port your Go app is listening on
```

### 🔹 Example 3: Advanced (Liveness Probes)
```yaml
# How does K8s know if your Go app froze due to a Deadlock?
# It pings the /healthz endpoint every 10 seconds. If it fails, K8s murders the Pod and replaces it!
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
```

### 🔹 Example 4: Production (Graceful Shutdown in Go)
```go
// When K8s scales down, it sends a SIGTERM signal to your Go app.
// You MUST catch this signal and cleanly finish ongoing HTTP requests!
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

<-quit // Block until K8s sends the kill signal!
log.Println("Shutting down gracefully...")

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
server.Shutdown(ctx) // Drain connections cleanly!
```

### 🔹 Example 5: Interview
```yaml
# Q: What is the difference between a Liveness Probe and a Readiness Probe?
# A: Liveness checks if the app is frozen (Restarts the pod). 
# Readiness checks if the app is fully booted and connected to the DB. 
# If Readiness fails, K8s stops sending HTTP traffic to the Pod, but doesn't kill it.
```

---

## 8. Production Examples
1. **Horizontal Pod Autoscaling (HPA)**: You tell K8s: "If average CPU usage exceeds 70%, automatically scale my Go API from 3 Pods to 50 Pods." When the viral traffic spike ends, K8s scales it back down to 3 to save AWS money.
2. **Rolling Updates**: When you deploy `v2.0` of your Go API, K8s doesn't shut down all `v1.0` pods at once. It shuts down one `v1.0` pod, boots one `v2.0` pod, and repeats. Zero-Downtime deployments!

---

## 9. Performance & Benchmarking
* **Startup Speed matters**: Java Spring Boot apps take 30 seconds to boot. When a massive traffic spike hits, K8s tries to scale up Java, but by the time the pods boot, the users have already left. Go binaries boot in 0.05 seconds. By the time the traffic hits, the Go pods are already running!

---

## 10. Best Practices
* ✅ **Do**: Set strict `resources.requests` and `resources.limits` for CPU and Memory in your YAML. This allows the Scheduler to intelligently pack Pods onto physical servers without overloading them.
* ❌ **Don't**: Rely on local file storage. Pods are ephemeral! If a pod crashes and moves to another server, anything saved to `os.WriteFile` is deleted forever. Always use S3 or PostgreSQL.
* 🏢 **Google / Uber / Netflix Style**: Use `GitOps` (e.g., ArgoCD). You never run `kubectl apply` manually. You push YAML files to a GitHub repository, and ArgoCD automatically syncs the cluster to match the Git repository.

---

## 11. Common Mistakes
1. **Ignoring SIGTERM**: If your Go app ignores the shutdown signal, K8s waits 30 seconds and then violently murders the process (`SIGKILL`). Any user currently paying via Stripe will have their transaction severed instantly.
2. **CrashLoopBackOff**: The most infamous K8s error. It means your Go app is crashing instantly on boot (usually a missing Environment Variable or DB connection failure). K8s keeps trying to restart it, backing off exponentially.

---

## 12. Debugging
How to troubleshoot K8s in production:
* **The Holy Trinity of debugging**:
  1. `kubectl get pods` (Are they running or crashing?)
  2. `kubectl describe pod <name>` (Why did K8s kill it? Look at the 'Events' at the bottom).
  3. `kubectl logs <name>` (What was the exact Go panic stack trace?)

---

## 13. Exercises
1. **Easy**: Write a basic `Deployment` YAML for a Go HTTP server.
2. **Medium**: Add a `Service` YAML to expose the Deployment internally.
3. **Hard**: Modify your Go code to implement Graceful Shutdown catching `SIGTERM`.
4. **Expert**: Use `minikube` or `kind` to boot a local K8s cluster, apply your YAML files, and test a rolling update by changing the image version!

---

## 14. Quiz
1. **MCQ**: What Kubernetes object is responsible for ensuring exactly 5 replicas of a Pod are always running?
   * (A) The Service (B) The Ingress (C) The Deployment / ReplicaSet. *(Answer: C)*
2. **System Design Follow-up**: Why shouldn't you run PostgreSQL directly inside Kubernetes? *(Because Databases require persistent, highly-optimized block storage. While possible with StatefulSets, a pod crashing and moving to a new node makes data replication extremely fragile. Use managed databases like AWS RDS instead!)*

---

## 15. FAANG Interview Questions
* **Beginner**: What is a Pod? (Hint: It is the smallest deployable unit, containing one or more tightly coupled containers).
* **Intermediate**: Explain how a Service routes traffic to Pods dynamically using Labels and Selectors.
* **Senior (Google/Meta)**: Explain the internal architecture of the Kubernetes Control Plane. How do the API Server, Scheduler, and Controller Manager interact with `etcd`?

---

## 16. Mini Project
**Zero-Downtime Deployment**
* Write a Go HTTP server that returns `{"version": "v1"}`. Dockerize it.
* Write a K8s Deployment YAML for 3 replicas.
* Write a bash script that curls the K8s Service IP every 0.1 seconds in an infinite loop.
* Update the image to `v2`. Run `kubectl apply`.
* Watch the curl output gracefully transition from `v1` to `v2` without dropping a single HTTP request!

---

## 17. Enterprise Features & Observability
* **Custom Resource Definitions (CRDs)**: Because K8s is built in Go, you can extend the K8s API! You can write a Go "Operator" that watches for a custom YAML file (e.g., `kind: PostgresDatabase`), and your Go code automatically talks to the AWS API to provision a real database!

---

## 18. Source Code Reading
Walkthrough of `kubernetes/kubernetes`.
* **The Informer Pattern**: Study how K8s controllers are written in Go. They don't constantly poll `etcd` (which would crash it). They use `Informers`, which establish long-running WebSockets to the API server and react to streaming events in real-time.

---

## 19. Architecture
* **Namespaces**: In a large company, different teams share the same K8s cluster. You use Namespaces (e.g., `billing-team`, `auth-team`) to logically isolate their deployments and apply strict CPU quotas per team.

---

## 20. Summary & Cheat Sheet
* **Pod**: The running Go container.
* **Deployment**: Keeps X Pods running and handles updates.
* **Service**: Internal Load Balancer / DNS.
* **Ingress**: External Load Balancer (Public IP).
