# Helm (Kubernetes Package Manager)

## 1. Learning Objectives
* **What you'll learn**: How to templatize complex Kubernetes YAML files using Helm Charts, allowing dynamic deployments across multiple environments (Dev, Staging, Prod).
* **Why it matters**: Managing 50 hardcoded Kubernetes YAML files is a nightmare. If you want to deploy to "Staging" instead of "Production", you have to manually copy-paste and change the database URLs in all 50 files! Helm solves this dynamically.
* **Where it's used**: Almost every Kubernetes production cluster uses Helm to package and distribute applications.

---

## 2. Real-world Story
Imagine writing a standard cover letter for a job application. 
Without Helm, you have to write a completely new letter from scratch for every company, manually typing the company name 10 times.
With Helm, you create a Template: "Dear {{ .CompanyName }}, I would love to work at {{ .CompanyName }}." 
When applying to Google, you just provide a tiny values file: `CompanyName: Google`. Helm injects the variable into the template and generates the perfect final letter instantly!

---

## 3. Visual Learning (Execution Flow & Architecture)
```mermaid
graph TD
    A[Helm Chart / Templates] -->|Combines with| C{Helm Engine}
    B[values-staging.yaml] -->|Combines with| C
    
    C -->|Generates Final YAML| D[deployment.yaml]
    C -->|Generates Final YAML| E[service.yaml]
    
    D -->|helm install| F(Kubernetes Cluster)
    
    style C fill:#0f172a,color:#fff
```

---

## 4. Internal Working (Under the Hood)
Helm is simply a powerful text templating engine built directly on top of Go's native `text/template` standard library!
1. **Templates**: Your K8s YAML files containing Go template injection points (e.g., `{{ .Values.replicaCount }}`).
2. **Values**: A simple `values.yaml` file containing the actual variables (e.g., `replicaCount: 3`).
3. **Release**: When you run `helm install`, Helm renders the templates, sends them to the Kubernetes API, and tracks the exact version (Release) in a secret inside K8s, allowing easy rollbacks.

---

## 5. Compiler Behavior
* **Go Templates**: Because Helm uses pure Go templating, you have access to logic operators directly inside your YAML! You can use `{{ if eq .Values.env "prod" }}`, `{{ range .Values.ports }}`, and advanced string manipulation functions provided by the `sprig` library.

---

## 6. Memory Management
* **Client-Side Rendering**: In Helm 3, the entire rendering engine runs locally on your laptop or CI/CD runner. It requires virtually zero memory overhead on the actual Kubernetes cluster, acting just like a smart wrapper around `kubectl apply`.

---

## 7. Code Examples

### 🔹 Example 1: Simple (The Values File)
```yaml
# values.yaml
replicaCount: 3
image:
  repository: myrepo/goverse-api
  tag: "v1.2"
databaseURL: "postgres://user:pass@prod-db:5432"
```

### 🔹 Example 2: Intermediate (The Template File)
```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: goverse-api
spec:
  # Injected directly from values.yaml!
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      app: goverse
  template:
    metadata:
      labels:
        app: goverse
    spec:
      containers:
      - name: api
        # Dynamically constructing the image tag!
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        env:
        - name: DB_URL
          value: {{ .Values.databaseURL | quote }}
```

### 🔹 Example 3: Advanced (Logic and Flow Control)
```yaml
# Only deploy the Prometheus ServiceMonitor if monitoring is enabled in values!
{{- if .Values.monitoring.enabled }}
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: goverse-monitor
spec:
  endpoints:
  - port: metrics
{{- end }}
```

### 🔹 Example 4: Production (Environment Overrides)
```bash
# In CI/CD, you maintain a base values.yaml, but override specific files per environment!
helm upgrade --install goverse-api ./my-chart \
  -f values.yaml \
  -f values-production.yaml \
  --set image.tag=$GITHUB_SHA
```

### 🔹 Example 5: Interview
```yaml
# Q: What is the purpose of the Chart.yaml file?
# A: It acts like a package.json for K8s. It contains the name, version, 
# and description of the Helm chart, and allows you to declare dependencies 
# on OTHER Helm charts (like automatically spinning up a Redis chart).
```

---

## 8. Production Examples
1. **Third-Party Installations**: Installing massive, complex software like Prometheus or Kafka manually on Kubernetes is insanely difficult. With Helm, it's one line: `helm install prometheus prometheus-community/kube-prometheus-stack`.
2. **Standardization**: A platform team writes one single "Golden" Helm chart. 50 different microservice teams use the exact same chart, just providing their own `values.yaml`, ensuring every app follows company security standards!

---

## 9. Performance & Benchmarking
* **Rollbacks**: Because Helm tracks the complete state of every deployment as a "Release", if a new deployment crashes your Go app, you can instantly revert the entire cluster state to exactly how it was 5 minutes ago by running `helm rollback goverse-api 1`.

---

## 10. Best Practices
* ✅ **Do**: Use the `helm template` command before installing. It prints the final, rendered YAML to your terminal so you can verify the variables injected correctly without actually deploying anything!
* ❌ **Don't**: Put sensitive secrets (like API Keys) in plain text inside `values.yaml` and push it to Git. Use tools like `SealedSecrets` or `ExternalSecrets` to inject them safely at runtime.
* 🏢 **Google / Uber / Netflix Style**: Package your final Helm Chart into an OCI Artifact and push it to a Docker Registry (like AWS ECR), treating your Kubernetes YAML templates with the same immutability as your Docker images.

---

## 11. Common Mistakes
1. **Whitespace Errors**: YAML is strictly sensitive to indentation. A very common error is `{{ .Values.variable }}` injecting a string that misaligns the YAML indentation by 2 spaces, corrupting the entire file. Use the `indent` function: `{{ .Values.configBlock | indent 4 }}`.
2. **Missing Quotes**: If `.Values.version` is `1.20`, YAML parses it as a float. You must use `value: {{ .Values.version | quote }}` to ensure Kubernetes receives it as a string `"1.20"`.

---

## 12. Debugging
How to troubleshoot Helm in production:
* **The Dry Run**: Always run `helm upgrade --install ... --dry-run --debug`. This simulates the entire deployment process against the active Kubernetes API, catching schema validation errors without mutating the cluster.

---

## 13. Exercises
1. **Easy**: Run `helm create my-chart` to auto-generate a boilerplate chart. Explore the folder structure.
2. **Medium**: Add a custom variable `contactEmail` to `values.yaml` and inject it as an Environment Variable into the Go container in `deployment.yaml`.
3. **Hard**: Use an `{{ if }}` statement to conditionally add a Redis container to the Pod only if `.Values.redis.enabled` is true.
4. **Expert**: Use `helm template` to debug a forced indentation error in your YAML.

---

## 14. Quiz
1. **MCQ**: What core Go technology powers the Helm templating engine?
   * (A) `html/template` (B) `text/template` (C) AST Parsing. *(Answer: B)*
2. **Code Review**: `{{ if .Values.enableFeature }} ... {{ end }}`. Why is the hyphen modifier `{{- if }}` usually preferred over standard `{{ if }}` in Helm? *(The hyphen automatically strips the newline character. Standard `if` will leave a blank, empty line in your YAML file, which can look messy or break strict YAML parsers).*

---

## 15. FAANG Interview Questions
* **Beginner**: Why use Helm instead of pure `kubectl apply`?
* **Intermediate**: Explain Helm Release history and how Rollbacks are mechanically achieved.
* **Senior (Google/Meta)**: Explain the difference between Helm and Kustomize. When would you choose Kustomize (patching) over Helm (templating)?

---

## 16. Mini Project
**The Multi-Environment Deployer**
* Build a basic Helm chart for a Go Web Server.
* Create a `values-dev.yaml` (Replicas: 1, LogLevel: DEBUG).
* Create a `values-prod.yaml` (Replicas: 5, LogLevel: ERROR).
* Write a bash script that uses `helm template` to output the final YAML for both environments and verifies the replica counts differ!

---

## 17. Enterprise Features & Observability
* **Helm Tests**: You can write a K8s Job template in a special `/tests` directory of your chart. When you run `helm test`, it boots a temporary Pod that runs integration tests against your newly deployed Go app. If it fails, the release is marked as broken!

---

## 18. Source Code Reading
Walkthrough of `helm/helm`.
* **The 3-Way Merge Patch**: Study how Helm 3 applies updates. It compares the Old Chart, the New Chart, and the *Live Cluster State*. If a human manually changed a replica count using `kubectl scale`, Helm's 3-way merge is smart enough to handle the conflict without completely destroying the manual patch!

---

## 19. Architecture
* **Chart Dependencies**: Your Microservice Chart can declare `postgresql` as a dependency in `Chart.yaml`. When you install your app, Helm automatically downloads the official Bitnami Postgres chart and boots a database alongside your app!

---

## 20. Summary & Cheat Sheet
* **Goal**: Templatize Kubernetes YAML.
* **Engine**: Go `text/template`.
* **Install**: `helm upgrade --install my-app ./chart`
* **Test Rendering**: `helm template ./chart`
