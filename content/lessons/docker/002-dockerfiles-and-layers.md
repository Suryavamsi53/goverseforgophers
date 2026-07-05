# Dockerfiles and Layer Caching

A Docker Image is a read-only template containing your Go application, its dependencies, and a stripped-down Linux filesystem.

You define how to build this image using a **Dockerfile**.

## 1. The Anatomy of a Dockerfile

```dockerfile
# 1. Base Image: We start with a pre-built image that contains Go
FROM golang:1.21-alpine

# 2. Set the working directory inside the container
WORKDIR /app

# 3. Copy our source code from our laptop into the container
COPY . .

# 4. Compile the Go application
RUN go build -o my-server main.go

# 5. Define the default command to run when the container starts!
CMD ["./my-server"]
```

You build this image using: `docker build -t my-go-app:v1 .`

## 2. The Union File System (Layers)

If you look closely at the `docker build` output, you will notice that Docker processes the file step-by-step.
Every single command (`FROM`, `WORKDIR`, `COPY`, `RUN`) creates a new, immutable **Layer**.

Docker uses a specialized storage driver (like OverlayFS). It stacks these read-only layers on top of each other. If Layer 1 adds `fileA.txt`, and Layer 2 modifies `fileA.txt`, the final container seamlessly merges them together so the application only sees the modified version.

**Why is this important?**
Because Layers are **Cached**.

If you change `main.go` and run `docker build` again:
1. Docker checks the cache for `FROM golang:1.21-alpine` (Match! Instantly loads).
2. Docker checks the cache for `WORKDIR /app` (Match! Instantly loads).
3. Docker checks `COPY . .`. It notices the checksum of `main.go` has changed. The cache is **BUSTED**.
4. Because this layer changed, Docker must physically re-execute this step, AND **every single step below it** (`RUN go build`).

## 3. Optimizing the Build Cache

If you look at the Dockerfile above, there is a massive flaw. 
If we change a single line in our `main.go` file, the `COPY . .` step busts the cache.
This means Docker will execute `RUN go build`. The Go compiler will have to re-download all 50 dependencies from GitHub (`go mod download`) every single time we build the image! This adds 3 minutes to every build!

**The Enterprise Fix:** We must copy the `go.mod` files FIRST, download the dependencies, and THEN copy the source code.

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app

# 1. Copy ONLY the dependency files
COPY go.mod go.sum ./

# 2. Download the dependencies!
# This layer will ONLY bust if we actually change our go.mod file!
RUN go mod download

# 3. Now copy the rest of the source code
COPY . .

# 4. Compile
RUN go build -o my-server main.go

CMD ["./my-server"]
```

Now, if you change `main.go`, Docker uses the cached dependencies (Step 2) and instantly proceeds to compilation (Step 4). Your build time drops from 3 minutes to 3 seconds!

## 4. Image Bloat (The Alpine Advantage)

The base image you choose is critical.
If you use `FROM ubuntu`, your final Docker image will be ~500 Megabytes.
If you use `FROM alpine` (a tiny, security-focused Linux distribution built on musl libc), your final image will be ~15 Megabytes.

Smaller images pull faster from Docker Hub, boot faster in Kubernetes, and drastically reduce the surface area for security vulnerabilities!
