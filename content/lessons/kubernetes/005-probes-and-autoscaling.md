# Probes and Autoscaling

If your Go application panics and crashes, the Docker process terminates, and Kubernetes instantly detects the failure and restarts the container. 
But what if the Go application suffers a Deadlock? The container is technically still running, the process hasn't crashed, but the application is frozen and refusing to answer HTTP requests!

Kubernetes cannot read minds. You must configure **Probes**.

## 1. Health Checks (Liveness and Readiness)

You configure Probes in your Deployment YAML to teach Kubernetes how to monitor the health of your Go application.

1. **Liveness Probe**: "Is the application frozen?"
   * Kubernetes makes an HTTP GET request to `/healthz` every 10 seconds.
   * If your Go app fails to return a `200 OK` (maybe because it's deadlocked), Kubernetes mercilessly assassinates the Pod and restarts it.
2. **Readiness Probe**: "Is the application ready to receive traffic?"
   * When a Pod boots up, it takes time to connect to PostgreSQL. During this 2-second window, it cannot process user requests.
   * Kubernetes hits `/readyz`. If it returns `500`, Kubernetes leaves the Pod alive, but removes it from the Load Balancer (Service). No user traffic is routed to it until it returns `200 OK`!

```yaml
containers:
- name: go-server
  # ...
  livenessProbe:
    httpGet:
      path: /healthz
      port: 8080
    periodSeconds: 10
  readinessProbe:
    httpGet:
      path: /readyz
      port: 8080
    initialDelaySeconds: 2 # Give it 2 seconds to boot before checking!
```

## 2. Resource Requests and Limits

Before Kubernetes can intelligently place a Pod on a Worker Node, it needs to know how heavy the Pod is.
If you deploy a Pod without specifying limits, it might consume 100% of the server's CPU, choking out all the other Pods!

* **Requests**: The guaranteed baseline. If you request 250m (a quarter of a CPU core), Kubernetes guarantees this Pod will be placed on a server that has at least that much free space.
* **Limits**: The hard cap. If the Go app has a memory leak and tries to exceed the memory Limit, the Linux OOM-Killer instantly terminates it.

```yaml
resources:
  requests:
    memory: "128Mi"
    cpu: "250m"
  limits:
    memory: "512Mi"
    cpu: "1000m" # Max 1 full CPU core
```

## 3. Horizontal Pod Autoscaler (HPA)

If your website gets mentioned on the news, traffic spikes from 1,000 requests/sec to 50,000 requests/sec. Your 3 Go Pods will max out their CPUs and crash.

You do not want to wake up at 3:00 AM to manually edit the Deployment YAML from `replicas: 3` to `replicas: 50`.

The **HPA** automates this.

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: billing-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: billing-deployment
  minReplicas: 3
  maxReplicas: 50
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

**The Magic:** The HPA controller constantly monitors the CPU usage of your 3 Pods. If the average CPU crosses 70%, the HPA dynamically modifies your Deployment to add more Pods! It scales up to 50, absorbing the massive traffic spike. When the news event is over and traffic drops, the HPA slowly kills off the extra Pods to save you AWS server costs!
