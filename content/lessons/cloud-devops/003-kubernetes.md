# Kubernetes (K8s)

Docker Compose is for your laptop. **Kubernetes** is for production. 

Kubernetes is a container orchestration engine that automatically scales your Go containers, restarts them if they crash, and load-balances traffic across hundreds of servers.

## 1. The Deployment (Pods)

The fundamental unit of Kubernetes is a **Pod** (a wrapper around your Docker container). You never deploy a Pod directly; you create a **Deployment**. 

A Deployment tells Kubernetes: *"Make sure exactly 3 copies of my Go app are running at all times."*

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-api-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-api
  template:
    metadata:
      labels:
        app: go-api
    spec:
      containers:
      - name: go-api
        image: myregistry.com/my-go-api:v1.0.0
        ports:
        - containerPort: 8080
        
        # CPU/RAM Limits prevent a single Pod from crashing the whole Node
        resources:
          requests:
            memory: "64Mi"
            cpu: "100m"
          limits:
            memory: "128Mi"
            cpu: "500m"
            
        # Liveness Probe (Restarts the pod if it deadlocks)
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 3
          periodSeconds: 10
```

## 2. The Service (Internal Networking)

Your 3 Pods are given random IP addresses (like `10.0.1.5` and `10.0.2.14`). When a Pod dies and a new one replaces it, the IP address changes. 

To solve this, we create a **Service**. A Service provides a single, permanent internal DNS name and automatically load-balances traffic to the underlying Pods.

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: go-api-service
spec:
  selector:
    app: go-api # This matches the labels in the Deployment!
  ports:
    - protocol: TCP
      port: 80       # Port the Service listens on
      targetPort: 8080 # Port the Go container listens on
```
Now, any other microservice in the cluster can simply make an HTTP request to `http://go-api-service`!

## 3. The Ingress (Public Routing)

The Service is strictly internal. If you want users on the public internet to access your API, you create an **Ingress**. 

The Ingress connects a public domain name (like `api.myapp.com`) to your internal Service, and handles SSL/TLS termination automatically.

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: go-api-ingress
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod" # Auto-generates SSL certificates!
spec:
  rules:
  - host: api.myapp.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: go-api-service
            port:
              number: 80
```
