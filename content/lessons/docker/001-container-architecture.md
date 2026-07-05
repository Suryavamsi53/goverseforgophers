# Container Architecture (VMs vs Containers)

Before Docker, if you wanted to deploy a Go application, you rented a physical server. 
If you wanted to deploy three different applications, you had a massive problem: Dependency Hell. App 1 required Postgres 10, and App 2 required Postgres 14. Installing both on the same Linux machine would break everything.

The industry solved this by introducing **Virtual Machines (VMs)**.

## 1. The Virtual Machine (Heavyweight)

A hypervisor (like VMware or VirtualBox) divides the physical hardware (CPU/RAM) into isolated chunks.
You install a completely independent operating system (Guest OS) on every single chunk.
* App 1 runs on VM A (Ubuntu, Postgres 10).
* App 2 runs on VM B (Alpine, Postgres 14).

**The Problem:** VMs are incredibly heavy. An Ubuntu VM takes 5 Gigabytes of hard drive space and 2 Gigabytes of RAM just to boot up, even if the Go application inside it only uses 10 Megabytes of RAM! Booting a VM takes several minutes.

## 2. The Container (Lightweight)

In 2013, Docker popularized **Containers**.

A Container is NOT a Virtual Machine. It does not boot a Guest OS.
A Container is just a regular Linux process that has been heavily sandboxed using two Linux kernel features:

1. **Namespaces**: Isolates what the process can *see*. (e.g., A containerized Go app cannot see the processes running on the host machine, and cannot see the host's file system. It thinks it is alone on a brand new hard drive).
2. **Cgroups (Control Groups)**: Isolates what the process can *use*. (e.g., You restrict the container to exactly 2 CPU cores and 500MB of RAM. If the Go app tries to use 600MB, the Linux kernel instantly kills it with an OOM error).

**The Superpower:** Because a container is just a sandboxed Linux process (using the host's exact same Linux kernel), booting a container takes milliseconds. A container only uses the RAM required by the Go application (e.g., 15 Megabytes). You can pack thousands of containers onto a single physical server!

## 3. The Docker Daemon Architecture

Docker operates on a Client-Server model.

1. **Docker CLI**: The terminal command you type (e.g., `docker run ubuntu`).
2. **Docker Daemon (dockerd)**: The background server that actually creates the namespaces, configures the cgroups, and manages the network bridging.
3. **Containerd / Runc**: The ultra-low-level binaries that physically spawn the Linux processes. (You rarely interact with these directly).

When you run a container on your Mac, it feels like magic. But macOS does not have Linux kernel features (Namespaces/Cgroups)! Under the hood, Docker Desktop silently boots a tiny, hidden Linux VM on your Mac, and runs the containers inside *that* VM.

## 4. The Orchestration Leap (Kubernetes)

Docker is amazing for running 5 containers on a single server.
But what if your E-Commerce platform has 5,000 containers across 100 physical servers?

If Server #14 loses power, how do you know which containers died? How do you automatically restart them on Server #15?
Docker cannot do this. You need a **Container Orchestrator**. 

This is the entire purpose of **Kubernetes (K8s)**. Kubernetes acts as the master brain, communicating with the Docker Daemon on every single physical server to ensure the correct number of containers are running at all times.
