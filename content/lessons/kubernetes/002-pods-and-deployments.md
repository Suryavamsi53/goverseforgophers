# Pods, ReplicaSets, and Deployments

In Kubernetes, you do not deploy a "Container". The smallest deployable unit in Kubernetes is a **Pod**.

## 1. The Pod

A Pod is a logical wrapper around one or more containers. 
* Usually, a Pod contains exactly 1 container (your Go application).
* Sometimes, a Pod contains multiple containers (e.g., your Go app + an Envoy Sidecar Proxy). 

**The Rule of the Pod:**
All containers inside the same Pod share the exact same Local Network and the exact same Hard Drive. They can talk to each other using `localhost`. They are physically guaranteed to run on the exact same Worker Node.

```yaml
# pod.yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-go-app
spec:
  containers:
  - name: go-server
    image: my-registry.com/my-go-app:v1
```

**The Flaw:** If you deploy a Pod directly, and the physical server it is running on catches fire... the Pod dies, and Kubernetes does **not** bring it back. Pods are mortal.

## 2. The ReplicaSet (High Availability)

To ensure High Availability, you never deploy a Pod directly. You deploy a **ReplicaSet**.

A ReplicaSet's only job is to guarantee that a specific number of identical Pods are running at all times.

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: go-app-replicas
spec:
  replicas: 3
  # The ReplicaSet uses 'Labels' to find its Pods!
  selector:
    matchLabels:
      app: billing
  template: # The Pod blueprint!
    metadata:
      labels:
        app: billing
    spec:
      containers:
      - name: go-server
        image: my-registry.com/my-go-app:v1
```

If you delete one of those 3 Pods manually, the ReplicaSet Controller instantly detects the mismatch (Current=2, Desired=3) and spins up a new one in milliseconds!

## 3. The Deployment (Rolling Updates)

ReplicaSets are great for High Availability, but they cannot handle Version Upgrades.
If you want to upgrade your Go app from `v1` to `v2`, you have a problem. You don't want to delete all 3 `v1` Pods simultaneously, because your API will go offline!

You use a **Deployment**.

A Deployment manages ReplicaSets. It provides **Zero-Downtime Rolling Updates**.

```yaml
# deployment.yaml (The Enterprise Standard)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: billing-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: billing
  template:
    metadata:
      labels:
        app: billing
    spec:
      containers:
      - name: go-server
        image: my-registry.com/my-go-app:v2 # We upgraded to v2!
```

### The Rolling Update Strategy
When you apply this YAML, the Deployment does something magical:
1. It creates a brand new ReplicaSet for `v2`.
2. It scales the `v2` ReplicaSet up to 1 Pod.
3. It scales the old `v1` ReplicaSet down from 3 to 2.
4. It waits for the `v2` Pod to boot successfully.
5. It repeats the process one-by-one until `v2` is at 3, and `v1` is at 0!

Your users experience absolutely zero downtime during the entire migration!

If `v2` contains a fatal panic bug and crashes immediately, the Deployment detects it, halts the rollout automatically, and allows you to instantly rollback to `v1`!
