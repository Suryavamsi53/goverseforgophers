# Redis Cluster and Sentinel (High Availability)

A single Redis instance becomes a Single Point of Failure. If the physical server burns down, your Go application's Cache is destroyed. 

To achieve High Availability, Redis provides two different architectural patterns: **Sentinel** and **Cluster**.

## 1. Master-Replica Replication

The foundation of both patterns is Replication.
You configure one Redis node as the **Master** (Leader). All Go applications write data here.
You configure two Redis nodes as **Replicas** (Followers). They establish a permanent TCP connection to the Master, and the Master streams every single command to them asynchronously.

* **Read Scaling**: Because the Replicas hold an exact copy of the data, your Go application can send all `GET` requests to the Replicas, distributing the CPU load!

## 2. Redis Sentinel (Automatic Failover)

What happens if the Master crashes? The Go application cannot write any data, because the Replicas are strictly Read-Only! A human has to SSH into a server and manually promote a Replica to Master. 

**Redis Sentinel** is an independent, distributed monitoring system that automates this.

1. You run 3 lightweight Sentinel processes alongside your Redis servers.
2. The Sentinels constantly ping the Master.
3. If the Master dies, the 3 Sentinels vote (using a Quorum).
4. If the majority agree the Master is dead, Sentinel automatically reconfigures a Replica to become the new Master!
5. **The Go Integration**: In your Go code, you do NOT connect to the Redis Master directly. You connect to the Sentinel cluster! Sentinel tells your Go application the IP address of the new Master instantly!

```go
// Connecting to Sentinel in Go!
rdb := redis.NewFailoverClient(&redis.FailoverOptions{
    MasterName:    "mymaster",
    SentinelAddrs: []string{"sentinel1:26379", "sentinel2:26379", "sentinel3:26379"},
})
```

## 3. Redis Cluster (Horizontal Sharding)

Sentinel solves High Availability, but it does NOT solve RAM limitations. 
If you need 500GB of RAM, and your physical servers only have 128GB of RAM, Sentinel cannot help you. Every Replica in a Sentinel setup holds a 100% complete copy of the data.

To scale RAM horizontally across multiple servers, you must use **Redis Cluster**.

Redis Cluster shards your data automatically.
1. Redis maps your keys to **16,384 Hash Slots**.
2. If you have 3 Master Nodes, Node A gets slots 0-5500, Node B gets 5501-11000, etc.
3. When your Go application runs `SET user:42 "data"`, the Go Redis client runs a CRC16 hash on the key, determines it belongs in Slot 8000, and routes the TCP request directly to Master Node B!

Redis Cluster also has High Availability built-in (every Master has its own Replicas). It completely eliminates the need for Sentinel. 

### The Hash Tag Problem
If you use Redis Cluster, multi-key operations (like `MGET` or Transactions) will crash if the keys belong to different Hash Slots on different physical servers!
To fix this, use **Hash Tags**: `SET {user:42}:name "Bob"` and `SET {user:42}:age 30`. Redis only hashes the substring inside the `{}` braces, mathematically guaranteeing both keys are written to the exact same physical server!
