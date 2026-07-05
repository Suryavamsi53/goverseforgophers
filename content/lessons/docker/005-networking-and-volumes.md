# Networking and Volumes (Persistence)

Containers are designed to be **Ephemeral** (temporary and disposable). 
If you delete a Postgres container, the Linux kernel instantly destroys the isolated Cgroups and Namespaces, wiping the 8KB Postgres data pages from the hard drive forever. 

If you want data to survive the death of a container, you must use **Volumes**.

## 1. Docker Volumes (Bypassing the Sandbox)

A Volume is a folder on your physical Host machine (your laptop or the EC2 server) that is mathematically "mounted" straight through the container's sandbox wall.

If the container writes data to `/var/lib/postgresql/data`, it isn't actually writing to its isolated filesystem. It is writing directly to the host's physical hard drive.

When you run `docker rm -f postgres`, the container is destroyed, but the Volume remains safely on your physical hard drive!

### Configuring Volumes in Compose

```yaml
services:
  postgres:
    image: postgres:15-alpine
    # Mount the named volume 'pg_data' to the internal Postgres folder
    volumes:
      - pg_data:/var/lib/postgresql/data

# Define the named volumes at the bottom of the file
volumes:
  pg_data:
```

Next time you run `docker-compose up`, Docker will re-attach the exact same `pg_data` volume to the new Postgres container, and all your SQL tables will miraculously still be there!

## 2. Bind Mounts (Local Development)

A standard Volume is managed by Docker (you don't know exactly where the files are stored on your Mac).
A **Bind Mount** allows you to explicitly link a specific folder on your laptop to a folder in the container.

This is how Hot-Reloading works!
```yaml
services:
  api:
    image: golang:1.21
    # Maps the current directory (.) on your laptop to /app in the container!
    volumes:
      - .:/app
    command: go run main.go
```
If you edit `main.go` in VS Code on your laptop, the file instantly changes inside the container's `/app` folder!

## 3. Advanced Networking (Bridged vs Host)

By default, Docker uses a **Bridge Network**. 
Docker creates a virtual router inside your computer. All containers get a fake IP address (e.g., `172.17.0.2`). They can talk to each other, but your laptop's web browser cannot talk to them. 
To fix this, we use **Port Mapping** (`-p 8080:8080`), which tells the virtual router to punch a hole from your laptop's port 8080 into the container's port 8080.

### Host Networking (Extreme Performance)
When your Go API handles 100,000 requests per second, the NAT (Network Address Translation) required by Port Mapping consumes massive amounts of CPU. The virtual router becomes a bottleneck!

For extreme performance, you can use **Host Networking**:
`docker run --network host my-go-api`

In Host mode, the container does NOT get an isolated IP address. The sandbox wall is torn down, and the container binds directly to the physical server's actual Network Interface Card (NIC). You get bare-metal, native network speeds, completely bypassing the Docker Bridge bottleneck! 
*(Warning: Host mode only works on Linux, it does not work on Docker Desktop for Mac/Windows).*
