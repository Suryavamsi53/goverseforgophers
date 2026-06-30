# Databases: SQL vs NoSQL, ACID vs BASE, Sharding & Replication

---

# Table of Contents

* Introduction
* Learning Objectives
* Prerequisites
* Why This Topic Exists
* SQL vs NoSQL
* ACID vs BASE Transactions
* Database Scaling: Replication
* Database Scaling: Sharding
* Code Examples & Good Principles
* Architecture Diagram
* Real-World Analogy
* Interview Questions
* Quiz
* Exercises
* Summary
* Key Takeaways
* Further Reading
* Next Chapter

---

# Introduction

At the heart of every application is its data. If your web servers go down, users get annoyed, but they come back later. If your database goes down and loses data, your company might go out of business.

Choosing the right database, understanding how it guarantees data safety, and knowing how to scale it when traffic spikes is the most critical and complex part of system design.

---

# Learning Objectives

After completing this chapter you will be able to:

* Choose between SQL (Relational) and NoSQL (Non-Relational) based on data access patterns.
* Differentiate between ACID (strict consistency) and BASE (eventual consistency) paradigms.
* Explain how Read Replicas scale read-heavy workloads.
* Understand the immense complexity and necessity of Database Sharding for write-heavy workloads.
* Implement safe, transaction-bound database code in Go.

---

# Prerequisites

Before reading this chapter you should know:

* System Design Goals & CAP Theorem (`01-Introduction.md`)
* Vertical vs. Horizontal Scaling (`02-Scalability.md`)

---

# Why This Topic Exists

During a system design interview, if you are asked to design Twitter, the interviewer knows your Go web server will just be a stateless API. What they *really* care about is the database layer. How do you store 500 million tweets a day? A single PostgreSQL instance will melt. Do you use Cassandra? MongoDB? How do you partition the data? 

Mastering database scaling is what separates Junior engineers from Senior architects.

---

# SQL vs NoSQL

### SQL (Relational Databases)
* **Examples**: PostgreSQL, MySQL, Oracle.
* **Structure**: Data is highly structured into rigid Tables, Rows, and Columns.
* **Relationships**: Excels at complex queries joining multiple tables together (e.g., "Find all users who bought a product in the last 30 days and live in New York").
* **Scaling**: Historically very difficult to scale horizontally (though modern NewSQL systems are changing this). Usually scaled vertically (bigger servers).

### NoSQL (Non-Relational Databases)
* **Examples**: MongoDB (Document), Cassandra (Column-family), Neo4j (Graph).
* **Structure**: Flexible schema. You can store nested JSON documents. A user record might have an array of phone numbers, while another user doesn't.
* **Relationships**: Terrible at Joins. If you need to join data, you often have to do it in your application code, or "denormalize" the data (store duplicates).
* **Scaling**: Designed from day one to scale horizontally across hundreds of cheap commodity servers.

---

# ACID vs BASE Transactions

How does the database handle data mutations?

### ACID (Usually SQL)
* **A**tomicity: All parts of a transaction succeed, or none do. No half-finished states.
* **C**onsistency: Data always follows defined rules and constraints.
* **I**solation: Concurrent transactions don't interfere with each other.
* **D**urability: Once committed, data is permanently saved, even if the power goes out.
* **Use Case**: Financial systems. If you transfer $100 from Account A to Account B, it *must* be strictly ACID.

### BASE (Usually NoSQL)
* **B**asically **A**vailable: The system guarantees availability, even if some nodes fail.
* **S**oft state: The state of the system may change over time, even without input.
* **E**ventually consistent: The system will eventually become consistent once it stops receiving input.
* **Use Case**: Social media likes. If 1,000 people like a YouTube video at the exact same time, it's okay if some users see 998 likes for a few seconds. The strict overhead of ACID isn't worth the performance hit.

---

# Database Scaling: Replication

Most web applications are **Read-Heavy** (e.g., 100 reads for every 1 write).
To scale this, we use **Master-Slave (Leader-Follower) Replication**.

1. **Master (Leader)**: There is only ONE Master database. All `INSERT`, `UPDATE`, and `DELETE` commands go here. 
2. **Slaves (Followers/Read Replicas)**: There are multiple Replica databases. They constantly sync data from the Master. All `SELECT` (read) queries go to the Replicas.

**The Catch (Replication Lag)**: Because syncing takes a few milliseconds, a user might update their profile picture (goes to Master), hit refresh (reads from Replica), and see their old picture. This is eventual consistency in action.

---

# Database Scaling: Sharding

What if your application is **Write-Heavy**, and the single Master database can no longer handle the sheer volume of `INSERT` statements? You must use **Sharding (Data Partitioning)**.

Sharding means splitting your massive database into multiple smaller databases (Shards) based on a **Shard Key**. 
* Example: Users with IDs 1-1,000,000 go to Database A. Users with IDs 1,000,001-2,000,000 go to Database B.

**The Catch (The Sharding Nightmare)**: 
* What if Database A gets all the highly active users, while Database B is empty? (Hotspotting).
* You can no longer perform SQL JOINs across Shard A and Shard B.
* Sharding requires massive changes to your application code to route queries to the correct database. Avoid sharding until absolutely necessary.

---

# Code Examples & Good Principles

### 1. Connection Pooling in Go (Good Principle)
**Bad Practice**: Opening a new `sql.Open` connection for every incoming HTTP request. This will instantly exhaust the database's maximum connections.
**Good Practice**: Open the database connection *once* globally, and configure the underlying Connection Pool.

```go
package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func initDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://user:pass@localhost/db")
	if err != nil {
		log.Fatal(err)
	}

	// Principle: Always explicitly configure your connection pool.
	// Default settings will often lead to exhausted connections under load.
	db.SetMaxOpenConns(25)                 // Max total connections
	db.SetMaxIdleConns(25)                 // Keep all connections open to avoid reconnect overhead
	db.SetConnMaxLifetime(5 * time.Minute) // Retire old connections to clear memory leaks

	return db
}
```

### 2. ACID Transactions in Go (Good Principle)
When performing multiple related writes (like moving money), always use `db.BeginTx`.

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Principle: Group related mutations into a single ACID transaction.
func TransferFunds(ctx context.Context, db *sql.DB, fromAccount, toAccount int, amount float64) error {
	// 1. Begin the Transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	// 2. Ensure Rollback on error or panic (Defer is crucial here)
	defer tx.Rollback()

	// 3. Deduct from Sender
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromAccount)
	if err != nil {
		return fmt.Errorf("failed to deduct: %w", err) // Rollback happens via defer
	}

	// 4. Add to Receiver
	_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toAccount)
	if err != nil {
		return fmt.Errorf("failed to add: %w", err) // Rollback happens via defer
	}

	// 5. Commit the Transaction (Success!)
	if err = tx.Commit(); err != nil {
		return err
	}

	log.Println("Transfer successful!")
	return nil
}
```

---

# Architecture Diagram

```mermaid
flowchart TD
    App[Go Application]
    
    subgraph Database Cluster
        Master[(Master DB<br>Writes Only)]
        Rep1[(Read Replica 1)]
        Rep2[(Read Replica 2)]
        
        Master -- Async Sync --> Rep1
        Master -- Async Sync --> Rep2
    end
    
    App -- "INSERT / UPDATE" --> Master
    App -- "SELECT" --> Rep1
    App -- "SELECT" --> Rep2
```

---

# Real-World Analogy

* **SQL**: A highly organized Excel spreadsheet. Every column has a strict type. It's easy to cross-reference data using VLOOKUPs (JOINs).
* **NoSQL**: A massive filing cabinet filled with manila folders. Every folder (Document) belongs to one person. Inside the folder, they can put whatever papers they want in whatever format they want. It's very fast to grab a single person's folder, but finding "everyone who owns a red car" requires opening every single folder manually.
* **Replication**: The author writes the original book (Master). The publisher prints thousands of copies (Replicas) so millions of people can read it simultaneously. 
* **Sharding**: The encyclopedia is too heavy to carry as one book, so it is split into 26 smaller books (A-Z). You must know what letter you are looking for (Shard Key) to know which book to open.

---

# Interview Questions

## Beginner
**Q**: What is the primary difference between a Master database and a Read Replica?
*Answer*: The Master database accepts all write operations (`INSERT`, `UPDATE`, `DELETE`) and acts as the source of truth. Read Replicas are read-only copies that sync data from the Master and serve `SELECT` queries to distribute the read load.

## Intermediate
**Q**: If your startup is building a new e-commerce checkout system, would you choose MongoDB (NoSQL) or PostgreSQL (SQL)?
*Answer*: PostgreSQL (SQL). E-commerce checkouts deal with inventory, payments, and strict states. It requires ACID compliance to guarantee that money is transferred reliably and inventory isn't oversold. The flexible schema of NoSQL is a liability here.

## Advanced
**Q**: Explain "Celebrity Hotspotting" in a sharded database.
*Answer*: Hotspotting occurs when the chosen Shard Key distributes traffic unevenly. If you shard a social network by `UserID`, and Justin Bieber is on Shard 4, Shard 4 will receive millions of reads/writes a second while Shards 1-3 sit idle. Shard 4 will crash. To fix this, you must choose a better shard key, or implement complex caching layers specifically for celebrity accounts.

---

# Quiz

## Multiple Choice Questions
**1. Which ACID property ensures that a transaction is either completely executed or entirely aborted?**
A) Consistency
B) Atomicity
C) Isolation
*Answer*: B

## True or False
**Sharding is the process of adding more RAM and CPU to your existing Master database.**
*Answer*: False. That is Vertical Scaling. Sharding is Horizontal Scaling by partitioning the data across multiple separate database servers.

---

# Exercises

## Beginner
Research "Consistent Hashing". Why is standard modulo hashing (e.g., `UserID % TotalShards`) a terrible idea for sharding when you eventually need to add a new shard to the cluster?

## Intermediate
Review the `TransferFunds` Go example. Why is `defer tx.Rollback()` completely safe to use even if `tx.Commit()` succeeds at the end of the function? (Hint: Read the documentation for Go's `sql.Tx.Rollback`).

---

# Summary

The database is the anchor of your architecture. Use SQL for structured, relational data that requires strict ACID transactions. Use NoSQL for unstructured, massively scaled data that can tolerate eventual consistency. When traffic grows, scale reads using Replicas, but avoid Sharding for writes until you have absolutely exhausted all vertical scaling and caching options, as it introduces immense application-level complexity.

---

# Key Takeaways

* ✔ SQL = Relational, Strict Schema, ACID, Good for Joins.
* ✔ NoSQL = Non-Relational, Flexible Schema, BASE, Good for Horizontal Scaling.
* ✔ Replication scales reads (at the cost of replication lag).
* ✔ Sharding scales writes (but introduces immense complexity and hotspot risks).
* ✔ In Go, always configure your database connection pool limits.

---

# Further Reading
* [Go `database/sql` Connection Pooling Tutorial](https://www.alexedwards.net/blog/configuring-sqldb)
* [Designing Data-Intensive Applications (Book) - Martin Kleppmann](https://dataintensive.net/)

---

# Next Chapter
➡️ **Next:** `08-Message-Queues.md`
