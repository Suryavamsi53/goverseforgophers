# Why Concurrency? (I/O Bound vs CPU Bound)

To master concurrency, you must first understand the fundamental bottleneck of the application you are building. All software bottlenecks fall into two categories: **I/O Bound** and **CPU Bound**.

## 1. CPU Bound Workloads

A task is CPU Bound if the speed of the application is purely limited by the speed of the processor.

* **Examples**: Cryptographic hashing (Bcrypt), Video encoding, Machine Learning model training, calculating the digits of Pi.
* **The Reality**: If you have a 4-core processor, and you spawn 100 Goroutines to calculate prime numbers, your program will NOT run faster. In fact, it will run *slower* due to Context Switching overhead. You can only do 4 things at exactly the same time.

## 2. I/O Bound Workloads

A task is I/O (Input/Output) Bound if the speed of the application is limited by the network or the hard drive.

* **Examples**: HTTP Requests, Database Queries, reading files from an SSD, calling the Stripe API.
* **The Reality**: 99% of modern backend web development is I/O Bound.

When an application queries a database, the CPU sends the request over the network interface and then completely halts. It waits. This waiting time is measured in milliseconds (which is an eternity for a CPU capable of billions of instructions per second).

## 3. The Power of Goroutines in I/O

This is exactly why Go was created for the modern cloud.

When a Goroutine makes an HTTP request, the Go Runtime intercepts it. The Runtime instantly puts that Goroutine to sleep and says to the CPU, *"Don't wait! Grab the next Goroutine and execute it!"*

Because of this, a single 4-core Go web server can easily juggle 100,000 concurrent Goroutines. 

* 99,996 Goroutines are "asleep" waiting for the Database/Network to reply.
* 4 Goroutines are actively running on the 4 CPU cores, formatting JSON or parsing headers.

As soon as a database replies, the Go Runtime wakes up the sleeping Goroutine and schedules it back onto a CPU core. This is called **Multiplexing**, and it is the secret to Go's massive performance in microservices.
