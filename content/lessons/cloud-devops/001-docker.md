# Docker and Go

We briefly touched on Docker in the Go Fundamentals track, but in the Cloud DevOps track, we need to master it. Go and Docker are a match made in heaven—in fact, Docker itself is written entirely in Go!

## 1. The Anatomy of a Go Dockerfile

A production-grade Dockerfile for Go focuses on two things: **Build Caching** and **Security**.

```dockerfile
# STAGE 1: The Builder
# Use the official Golang alpine image to compile the code
FROM golang:1.22-alpine AS builder

# Install git (required for fetching some dependencies)
RUN apk update && apk add --no-cache git

WORKDIR /app

# CACHE OPTIMIZATION: Copy ONLY the go.mod and go.sum first!
# Docker caches layers. If we only copy the mod files, Docker will cache the 
# downloaded dependencies. When we change our Go source code later, 
# Docker won't redownload the entire internet!
COPY go.mod go.sum ./
RUN go mod download

# Now copy the rest of the source code
COPY . .

# SECURITY & SIZE: Build a statically linked binary.
# CGO_ENABLED=0 completely disables C dependencies, ensuring the binary
# can run on a literally empty operating system.
# -ldflags="-w -s" strips debugging information, shrinking the binary size by 30%!
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o myapp ./cmd/api

# STAGE 2: The Final Image
# "scratch" is an empty, 0-byte image provided by Docker.
FROM scratch

# (Optional) Copy CA Certificates so your Go app can make outbound HTTPS requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/myapp /myapp

# Run as a non-root user for security
USER 1000:1000

ENTRYPOINT ["/myapp"]
```

## 2. Why "Scratch" over "Alpine"?

Many tutorials tell you to use `alpine` as your final image because it is small (~5MB). 
**Enterprise teams use `scratch`.**

* **Alpine**: Contains a shell (`/bin/sh`), package manager (`apk`), and basic GNU utilities (`wget`, `grep`).
* **Scratch**: Contains literally nothing.

If a hacker discovers a Remote Code Execution (RCE) vulnerability in your Go app, they will try to pop a reverse shell or download malware using `wget`. 
If your container is `scratch`, **their attack fails instantly**. There is no shell to spawn, no `wget` to execute, and no file system to explore. They are trapped in a vacuum with only your compiled Go binary.

## 3. The .dockerignore File

Just like `.gitignore`, you must have a `.dockerignore` file. If you don't, `COPY . .` will accidentally copy your local `.git` folder, your `.env` files (leaking secrets into the image!), and your local `vendor/` directories.

```text
.git
.env
vendor/
**/*_test.go
```
