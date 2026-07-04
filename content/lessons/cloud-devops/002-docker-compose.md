# Docker Compose

While a `Dockerfile` builds a single image, most modern Go applications require a database (PostgreSQL), a cache (Redis), and maybe a message broker (Kafka). 

Running all of these manually on your laptop is a nightmare. **Docker Compose** orchestrates your entire local development environment with a single command: `docker-compose up`.

## 1. The Anatomy of a Go Stack

Here is a standard `docker-compose.yml` for a Go Web Server that relies on Postgres and Redis.

```yaml
version: '3.8'

services:
  # 1. The Go Application
  api:
    build: 
      context: .
      dockerfile: Dockerfile.dev # Use a dev-specific dockerfile for hot-reloading!
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/myapp?sslmode=disable
      - REDIS_URL=redis:6379
    depends_on:
      db:
        condition: service_healthy # Wait until Postgres is fully booted!
      redis:
        condition: service_started
    volumes:
      - .:/app # Mount local code into the container for live-reloading

  # 2. PostgreSQL Database
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data # Persist data across restarts
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d myapp"]
      interval: 5s
      timeout: 5s
      retries: 5

  # 3. Redis Cache
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

# Define the persistent volume for the database
volumes:
  pgdata:
```

## 2. The Internal Docker Network

Notice the `DATABASE_URL` for the Go API: `postgres://user:pass@db:5432`.
It does not say `localhost`. 

When you run Docker Compose, it automatically creates a secure internal DNS network. Your Go container can talk to the Postgres container simply by using its service name (`db`) as the hostname!

## 3. Hot-Reloading with `air`

If you compile your Go binary directly into the Docker image, you have to rebuild the image every time you change a line of code. This slows down development.

Instead, we use a tool called **Air** (`github.com/cosmtrek/air`) for live-reloading.

Create a `Dockerfile.dev`:
```dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Install Air for hot-reloading
RUN go install github.com/cosmtrek/air@latest

# We don't copy the code here! 
# We use the 'volumes' mount in docker-compose to mount our laptop's code into the container!

CMD ["air"]
```
Now, when you run `docker-compose up`, Air watches your mounted code. Whenever you hit `Ctrl+S` to save a `.go` file on your laptop, Air instantly recompiles and restarts the server inside the Docker container!
