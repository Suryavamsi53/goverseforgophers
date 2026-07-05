# Multi-Stage Builds (Go's Superpower)

In the previous lesson, we optimized our Go Dockerfile to cache dependencies. However, if you inspect the final Docker image size (`docker images`), it will still be ~300 Megabytes!

Why? Because our base image is `golang:1.21-alpine`. This image includes the entire Go compiler toolchain, git, and debugging tools.
Once our Go binary is compiled, we don't need the Go compiler anymore! Shipping the compiler to production is a massive waste of disk space and a severe security risk.

## 1. The Multi-Stage Solution

Docker allows you to use multiple `FROM` statements in a single Dockerfile.
You can compile your code in a massive, heavy image, and then copy *only* the compiled binary into a brand new, empty image!

```dockerfile
# ==========================================
# STAGE 1: The Builder (Heavy, contains Go Compiler)
# ==========================================
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Compile a static binary!
# CGO_ENABLED=0 ensures we don't dynamically link against C libraries.
RUN CGO_ENABLED=0 GOOS=linux go build -o my-server main.go

# ==========================================
# STAGE 2: The Final Image (Tiny, No Compiler!)
# ==========================================
# 'scratch' is a completely empty, 0-byte Docker image. 
# It doesn't even have a Linux shell (sh/bash)!
FROM scratch

# Copy ONLY the compiled binary from the 'builder' stage!
COPY --from=builder /app/my-server /my-server

# Expose port 8080 (Documentation only)
EXPOSE 8080

CMD ["/my-server"]
```

## 2. The Result: 10 Megabyte Images!

If you build this image, it will be incredibly tiny—literally just the size of the Go binary itself (~10 to 15 Megabytes). 

Because Go produces statically linked binaries (unlike Python or Node.js, which require a heavy runtime interpreter to be installed on the production server), Go is the undisputed king of Containerization.

## 3. The `scratch` Trade-offs

Using `FROM scratch` is the ultimate goal, but it has some drawbacks you must be aware of:

1. **No Shell (`/bin/sh`)**: You cannot `docker exec -it <container> sh` to debug it. There is no bash, no `ls`, no `cat`. The container *only* runs your Go binary. (This is fantastic for security!).
2. **No CA Certificates**: If your Go application needs to make an HTTPS request to Stripe or Google, it will fail! `scratch` doesn't contain the public Root Certificates required to verify SSL/TLS connections.
3. **No Timezones**: `time.LoadLocation("America/New_York")` will panic, because the timezone database doesn't exist.

### The Enterprise Compromise (Distroless / Alpine)

If `scratch` is too strict, teams usually fall back to two options for Stage 2:

* **`FROM alpine:latest`**: Adds ~5MB. Gives you a shell, a package manager (`apk add ca-certificates`), and basic Linux tools.
* **`FROM gcr.io/distroless/static-debian11`** (Built by Google): Adds ~2MB. It contains CA certificates and Timezones, but strictly NO shell, providing the perfect balance of security and functionality.
