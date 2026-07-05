# Process vs Thread vs Goroutine

To understand why Goroutines are revolutionary, we must understand how the Operating System (Linux, Windows, macOS) handles execution.

## 1. The Process

A Process is a running instance of an application (like Google Chrome, or a Go Web Server).
* **Isolation**: Every Process has its own private memory space. Process A cannot read Process B's memory.
* **Heavyweight**: Creating a Process takes significant time and RAM.
* **Context Switching**: If the OS CPU switches from running Process A to Process B, it has to swap out the entire memory map. This is incredibly slow (microseconds).

## 2. The OS Thread

Inside a Process, you can spawn multiple Threads (e.g., Java `Thread`, C++ `pthread`).
* **Shared Memory**: All Threads inside a Process share the same memory heap. (This is why Data Races occur!).
* **Medium Weight**: An OS Thread typically consumes **1 Megabyte to 8 Megabytes** of RAM just for its stack.
* **The Limit**: If you run a Java Tomcat server and try to spawn 10,000 OS threads to handle 10,000 concurrent users, it will consume 10 Gigabytes of RAM instantly. The server will crash.

## 3. The Goroutine

A Goroutine is a "User-Space Thread" or a "Green Thread" managed entirely by the Go Runtime, not the OS.

* **Featherweight**: A Goroutine starts with an initial stack size of just **2 Kilobytes** (2KB).
* **Dynamic Stacks**: If the Goroutine needs more memory, the Go Runtime dynamically grows the stack. If it needs less, it shrinks it.
* **The Limit**: Because they only cost 2KB, a standard Go web server can easily spawn **1,000,000 Goroutines** using only 2 Gigabytes of RAM.

## 4. Context Switching Overhead

The real magic of Goroutines is not just memory; it's speed.

If the OS CPU switches from OS Thread A to OS Thread B, it must trap into the OS Kernel, save 16 hardware registers, flush CPU caches, and restore the new thread. This takes ~1000-2000 nanoseconds.

When the Go Runtime switches from Goroutine A to Goroutine B, it happens entirely in "User Space" (no Kernel traps). It only has to save 3 registers (Program Counter, Stack Pointer, DX). This takes ~200 nanoseconds.

**Goroutines are 10x faster to switch and 500x cheaper to create than OS Threads.**
