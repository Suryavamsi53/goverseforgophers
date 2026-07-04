# Helm (The K8s Package Manager)

In the previous lesson, we wrote three different YAML files (Deployment, Service, Ingress) to deploy a single Go application. 

If you have 50 microservices across 3 environments (Dev, Staging, Prod), you would have to maintain hundreds of raw YAML files. If you need to change a single port, you have to find and edit 15 different files.

**Helm** is the package manager for Kubernetes. It solves this by replacing raw YAML files with **Templates**.

## 1. The Helm Chart Structure

A Helm Chart is simply a folder containing templates and a `values.yaml` file.

```text
my-go-api/
├── Chart.yaml          # Metadata (Name, Version)
├── values.yaml         # The variables! (Environment specific)
└── templates/
    ├── deployment.yaml # The Go template
    ├── service.yaml
    └── ingress.yaml
```

## 2. Templating YAML

Instead of hardcoding the number of replicas or the Docker image tag in your `deployment.yaml`, you inject variables using the Go text templating language `{{ ... }}`.

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Release.Name }}-api
spec:
  # Injects the replica count dynamically!
  replicas: {{ .Values.replicaCount }} 
  selector:
    matchLabels:
      app: {{ .Release.Name }}
  template:
    metadata:
      labels:
        app: {{ .Release.Name }}
    spec:
      containers:
      - name: api
        # Injects the image and version tag dynamically!
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        ports:
        - containerPort: 8080
```

## 3. The `values.yaml` File

The `values.yaml` file acts as the configuration hub for the entire chart.

```yaml
# values.yaml
replicaCount: 3

image:
  repository: myregistry.com/go-api
  tag: "v1.2.0" # You only update the version here!

environment: "production"
```

## 4. Environment Overrides

The true power of Helm is environment overrides. You maintain a single, pristine Helm chart in your repository. When deploying to Staging, you pass in a `values-staging.yaml` file. When deploying to Production, you pass in a `values-prod.yaml` file.

```bash
# Deploy to Staging (1 replica, v1.2.0-rc1)
helm upgrade --install my-api ./my-go-api -f values-staging.yaml

# Deploy to Production (10 replicas, v1.2.0)
helm upgrade --install my-api ./my-go-api -f values-prod.yaml
```

By using Helm, DevOps teams can standardize the architecture of every Go microservice in the company, ensuring they all use the exact same Liveness Probes, Resource Limits, and Ingress rules automatically!
