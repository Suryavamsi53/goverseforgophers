# Docker Compose (Local Environments)

A modern Go microservice does not run in isolation. It needs a PostgreSQL database, a Redis cache, and maybe a Kafka broker.

If a new developer joins your team, forcing them to manually run `docker run -p 5432:5432 postgres` and `docker run -p 6379:6379 redis` with 15 different environment variables is a nightmare. 

**Docker Compose** solves this. It allows you to define an entire multi-container architecture in a single declarative YAML file.

## 1. The `docker-compose.yml`

This file lives at the root of your Go project.

```yaml
version: '3.8'

# A 'Service' is a container we want to run
services:
  # 1. Our Go API
  api:
    build: . # Tell compose to build the Dockerfile in this directory!
    ports:
      - "8080:8080" # Map localhost:8080 to container:8080
    environment:
      - DB_HOST=postgres # The hostname is the name of the service below!
      - DB_USER=user
      - DB_PASS=pass
      - REDIS_URL=redis:6379
    depends_on:
      - postgres
      - redis

  # 2. The Database
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"

  # 3. The Cache
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
```

## 2. One Command to Rule Them All

Now, the new developer simply runs:
`docker-compose up -d`

Docker Compose will:
1. Download the Postgres and Redis images from Docker Hub.
2. Build the Go Dockerfile locally.
3. Create a private, isolated virtual Network.
4. Boot all three containers and attach them to the network.

## 3. Internal DNS Resolution (The Magic)

How does the Go API know the IP address of the Postgres container? 

Docker Compose provides **Automatic DNS Resolution**.
Inside the virtual network, Docker maps the Service Name directly to the container's IP address. 

In your Go code, you do NOT connect to `192.168.1.5:5432`. You simply connect to `postgres:5432`. Docker intercepts the DNS lookup for the word "postgres" and routes it to the correct container!

## 4. Dependency Ordering (`depends_on`)

If the Go API boots up in 0.1 seconds, but Postgres takes 4.0 seconds to initialize its database tables... the Go API will try to connect to Postgres, fail, and crash!

The `depends_on` block tells Docker Compose to boot Postgres *first*, and then boot the API. 
*Warning: Standard `depends_on` only waits for the container to START, not for the database to be "Ready". Your Go code MUST still implement retry-loops (e.g., using `cenkalti/backoff`) when opening database connections!*

## 5. Hot Reloading (Air)

Running `docker-compose up` is great, but if you change a line of Go code, you have to tear everything down, rebuild the image, and boot it up again. That destroys developer velocity.

To achieve Hot-Reloading inside Docker, Go developers use **Air** (`github.com/cosmtrek/air`).
Instead of compiling a static binary in your Dockerfile, you mount your live laptop source code directly into the container using a Volume, and run `air`. When you hit `CTRL+S` in VS Code, `air` detects the change, recompiles the Go binary inside the container, and restarts it in 1 second!
