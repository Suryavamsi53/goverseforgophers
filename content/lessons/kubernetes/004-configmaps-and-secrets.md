# ConfigMaps and Secrets

A foundational rule of the "12-Factor App" methodology is that you must completely separate Configuration from Code.

If your Go application hardcodes the PostgreSQL database URL in `main.go`, you must rebuild the entire Docker Image just to switch from the Staging database to the Production database. This is a massive anti-pattern.

Your Docker image should be completely agnostic. The environment variables should be injected by the environment (Kubernetes) at runtime!

## 1. ConfigMaps (Non-Sensitive Data)

A `ConfigMap` is a dictionary of key-value pairs used to store non-sensitive configuration data (like URLs, feature flags, or timeout durations).

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: billing-config
data:
  LOG_LEVEL: "debug"
  PAYMENT_GATEWAY_URL: "https://api.stripe.com"
```

You apply this to the cluster. Then, in your Deployment YAML, you inject the ConfigMap into your Go Pod as Environment Variables!

```yaml
# Inside your Deployment YAML
containers:
- name: go-server
  image: my-go-app:v1
  envFrom:
  - configMapRef:
      name: billing-config # Pulls in all the keys as ENV vars!
```
Now, in your Go code, `os.Getenv("PAYMENT_GATEWAY_URL")` returns the Stripe URL perfectly!

## 2. Secrets (Sensitive Data)

You must never put API Keys, Database Passwords, or TLS Certificates in a ConfigMap, because anyone with read-access to the cluster can see them in plain text.

Instead, you use a **Secret**.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: billing-secrets
type: Opaque
data:
  # Values MUST be Base64 encoded in the YAML!
  # echo -n "super_secret_password" | base64
  DB_PASSWORD: c3VwZXJfc2VjcmV0X3Bhc3N3b3Jk
```

Secrets are injected into the Deployment exactly the same way as ConfigMaps (via `envFrom: secretRef`).

### The Danger of Kubernetes Secrets
Kubernetes Secrets are **NOT** encrypted by default! 
The Base64 encoding in the YAML file is just an encoding to handle special characters, not encryption. Anyone can run `base64 --decode` to read your password! Furthermore, the Secret is stored in plain text in the `etcd` database on the Master Node!

**Enterprise Best Practices:**
1. **Enable etcd Encryption**: You must configure the Kubernetes API server to encrypt Secrets at rest on the hard drive.
2. **External Secret Managers**: Instead of using native K8s Secrets, Enterprise teams use tools like **HashiCorp Vault** or **AWS Secrets Manager**. A sidecar container pulls the secret dynamically from Vault and injects it directly into the Go application's memory, bypassing Kubernetes entirely!

## 3. Volume Mounting Configuration

Sometimes your Go application expects a physical configuration file (like a `.json` or `.toml` file) rather than environment variables.

You can mount a ConfigMap as a physical file inside the container's hard drive!

```yaml
# Inside your Deployment YAML
volumes:
- name: config-volume
  configMap:
    name: billing-config
containers:
- name: go-server
  volumeMounts:
  - name: config-volume
    mountPath: /etc/config # K8s creates files here based on the ConfigMap keys!
```
If your Go application reads `/etc/config/LOG_LEVEL`, it will instantly get the string `"debug"`.
