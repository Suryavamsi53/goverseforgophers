# Consensus Algorithms

If you have a primary Database Server (Leader) and two backup servers (Followers), how do the servers know if the Leader has crashed? 

If the network cable between the Leader and Follower A is cut, Follower A might think the Leader is dead. But Follower B can still talk to the Leader just fine! If Follower A panics and declares itself the new Leader, you now have **Two Leaders** accepting writes. This is called a **Split-Brain**, and it will instantly corrupt your entire database.

Distributed systems prevent Split-Brains using **Consensus Algorithms**.

## 1. The Quorum (Majority Rules)

Consensus means getting multiple unreliable nodes to agree on a single source of truth. The fundamental mathematical rule of consensus is the **Quorum**.

To make a decision (like electing a new Leader, or committing a database write), a strict majority of nodes must agree.
* Formula: `(N / 2) + 1`
* If you have 3 nodes, a Quorum is `2`.
* If you have 5 nodes, a Quorum is `3`.

**Why 3 or 5? Why never 4?**
If you have 4 nodes, a Quorum is 3. If a network partition splits the 4 nodes down the middle (2 on left, 2 on right), neither side has 3 nodes. The entire cluster goes offline. You gain zero extra fault tolerance by adding a 4th node. Consensus clusters must always be odd numbers!

## 2. The Raft Algorithm

The industry standard consensus algorithm used by Kubernetes (etcd), HashiCorp Consul, and CockroachDB is **Raft**. (The reference implementation of Raft is written in Go: `github.com/hashicorp/raft`).

Raft manages consensus through **Leader Election**.

1. **Heartbeats**: The Leader constantly sends "Heartbeat" pings to all Followers to say "I am alive".
2. **Election Timeout**: Every Follower has a randomized countdown timer (e.g., 150ms to 300ms).
3. **The Crash**: If the Leader crashes, the heartbeats stop.
4. **The Election**: Whichever Follower's random timer hits 0 first becomes a Candidate. It votes for itself and asks the other nodes for votes.
5. **The Quorum**: Because the other nodes haven't timed out yet, they grant their vote to the Candidate. The Candidate reaches Quorum, declares itself the new Leader, and immediately starts sending out new heartbeats to suppress further elections.

Because the timeouts are randomized, the system mathematically guarantees that a new Leader is chosen in under half a second, without a split-brain.

## 3. Log Replication

Once a Leader is elected, how do we write data?

1. The client sends a `SET X=5` command to the Leader.
2. The Leader does NOT commit it yet. It writes it to its local log, and forwards the command to the Followers.
3. The Followers write it to their logs and reply "Acknowledge".
4. Once the Leader receives Acknowledgements from a **Quorum** (majority) of nodes, it formally Commits the data to its state machine and tells the user `200 OK`.

Even if a Follower crashes, as long as a Quorum is alive, the database continues to accept writes!

## 4. Paxos vs Raft

Before Raft (created in 2013), the dominant consensus algorithm was **Paxos** (used by Google Spanner and AWS DynamoDB). 
Paxos is mathematically brilliant but incredibly difficult to understand and implement correctly. Raft was explicitly designed to be human-readable and modular, which is why it dominates the open-source ecosystem today.
