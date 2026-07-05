# Services and Ingress (Networking)

In the previous lesson, we used a Deployment to spin up 3 Pods for our Go API.
Every time a Pod boots up, it is assigned a random IP address (e.g., `10.244.1.5`, `10.244.2.9`). 

Because Pods are constantly dying and being recreated by the Deployment (with brand new IP addresses), you can never hardcode a Pod's IP address into another microservice. 
If the `Web-Frontend` needs to talk to the `Billing-API`, how does it find it?

## 1. The ClusterIP Service (Internal Load Balancing)

To solve this, Kubernetes introduces the **Service** object.

A Service creates a permanent, static IP address and a permanent DNS name that acts as an internal Load Balancer.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: billing-service
spec:
  type: ClusterIP # Only accessible from INSIDE the cluster!
  selector:
    app: billing # Magic: It automatically finds all Pods with this label!
  ports:
    - port: 80 # The port the Service listens on
      targetPort: 8080 # The port your Go app is running on!
```

Now, the `Web-Frontend` Go application simply makes an HTTP request to `http://billing-service:80`.
Kubernetes intercepts this DNS name, resolves it to the Service's permanent IP, and load-balances the request across the 3 ephemeral Pods!

## 2. Exposing to the Internet (NodePort & LoadBalancer)

A `ClusterIP` Service is internal only. You cannot access it from your laptop or from the public internet.

To expose a Service to the public, you have two basic options:
1. **NodePort**: Opens a high port (e.g., `30005`) on every physical Worker Node in the cluster. (Terrible for production, great for debugging).
2. **LoadBalancer**: If you are running on AWS or GCP, Kubernetes will physically talk to the AWS API, provision an AWS Elastic Load Balancer (ELB), and map it to your Service. (Great, but expensive! If you have 50 microservices, you will pay for 50 AWS Load Balancers!).

## 3. The Ingress Controller (The API Gateway)

To avoid paying for 50 AWS Load Balancers, you create exactly 1 Load Balancer, and point it at a Kubernetes **Ingress Controller** (like NGINX or Traefik).

The Ingress Controller is a specialized Layer 7 HTTP router running inside your cluster. You define routing rules using an `Ingress` YAML object.

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: main-routing
spec:
  rules:
  - host: api.mycompany.com
    http:
      paths:
      # Route /billing to the billing-service
      - path: /billing
        pathType: Prefix
        backend:
          service:
            name: billing-service
            port:
              number: 80
      # Route /users to the user-service
      - path: /users
        pathType: Prefix
        backend:
          service:
            name: user-service
            port:
              number: 80
```

This acts as your cluster's API Gateway. It terminates SSL/TLS (HTTPS), parses the URL paths, and routes the traffic to the correct internal `ClusterIP` services. You get infinite routing complexity while only paying Amazon for 1 single Load Balancer!
